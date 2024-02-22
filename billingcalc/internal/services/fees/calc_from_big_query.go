package fees

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/bigqueryutils"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatusesbq"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewardsbq"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalancesbq"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

const env_project_id = "PROJECT_ID"

func CalculateFromBigQuery(ctx context.Context, mfrFile io.Reader, sheetId string, mfrTab string, token string, periodBegin time.Time, periodEnd time.Time, firstExternalId int, invoiceDate time.Time) (*CalculatedFees, error) {
	var balanceAdjustments *balanceadjustments.BalanceAdjustments
	var dailyBalances *dailybalances.DailyBalance

	var wg sync.WaitGroup
	var mu sync.Mutex
	var stakingSummary, custodySummary StakingSummary
	var combinedWarnings []Warning
	var err1, err2 error

	projectId := getProjectId()

	bq, err := bigqueryutils.NewBigQueryWrapper(ctx, projectId)
	if err != nil {
		msg := fmt.Sprintf("Error calling NewBigQueryWrapper %v", err)
		return nil, errors.New(msg)
	}
	defer func() {
		if clientErr := bq.Close(); clientErr != nil {
			log.Fatalf("Failed to close the BigQuery connection: %v", clientErr)
		}
	}()

	// Processing Mfr
	mfr, err := mfr.ProcessMfr(ctx, mfrFile, sheetId, mfrTab, token)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	balances, err := getUnclaimedBalances(ctx, bq, periodBegin, periodEnd)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	rewards, err := getRewards(ctx, bq, periodBegin, periodEnd)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if rewards.IsEmpty() {
		errorMessage := fmt.Sprintf("Not found any Rewards for the period between %s and %s", periodBegin, periodEnd)
		return nil, errors.New(errorMessage)
	}

	operationsStatuses, err := getStatuses(ctx, bq, periodBegin, periodEnd)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	if operationsStatuses.IsEmpty() {
		errorMessage := fmt.Sprintf("Not found any Operations Statuses for the period between %s and %s", periodBegin, periodEnd)
		return nil, errors.New(errorMessage)
	}
	wg.Add(2)

	go func() {
		defer wg.Done()
		summary, stakingWarn := CalculateStakingFees(mfr, rewards, balances, balanceAdjustments, operationsStatuses, dailyBalances, firstExternalId, invoiceDate)
		mu.Lock()
		defer mu.Unlock()
		if summary != nil {
			stakingSummary = append(stakingSummary, summary...)
			combinedWarnings = append(combinedWarnings, stakingWarn...)
		} else {
			err1 = errors.New("error in CalculateStakingFees")
		}
	}()

	go func() {
		defer wg.Done()
		summary, custodyWarn := CalculateCustodyFees(mfr, rewards, balances, balanceAdjustments, operationsStatuses, dailyBalances, firstExternalId, invoiceDate)
		mu.Lock()
		defer mu.Unlock()
		if summary != nil {
			custodySummary = append(custodySummary, summary...)
			combinedWarnings = append(combinedWarnings, custodyWarn...)
		} else {
			err2 = errors.New("error in CalculateCustodyFees")
		}
	}()

	wg.Wait()

	mergedResults := make(map[string]*OrgResult)
	for _, summary := range []StakingSummary{stakingSummary, custodySummary} {
		for _, orgResult := range summary {
			if existingOrg, ok := mergedResults[orgResult.OrgName]; ok {
				existingOrg.Accounts = MergeAccounts(existingOrg.Accounts, orgResult.Accounts)
			} else {
				orgCopy := orgResult
				mergedResults[orgResult.OrgName] = &orgCopy
			}
		}
	}

	var combinedSummary StakingSummary
	for _, org := range mergedResults {
		combinedSummary = append(combinedSummary, *org)
	}

	if err1 != nil || err2 != nil {
		var errMsg string
		if err1 != nil {
			errMsg += err1.Error()
		}
		if err2 != nil {
			if errMsg != "" {
				errMsg += " | "
			}
			errMsg += err2.Error()
		}
		return nil, errors.New(errMsg)
	}

	return &CalculatedFees{
		Summary: combinedSummary,
		Warns:   combinedWarnings,
	}, nil
}

func getUnclaimedBalances(ctx context.Context, bq bigqueryutils.BigQueryWrapper, begin time.Time, end time.Time) (*ubalances.UnclaimedBalances, error) {
	debug.NewMessage("Start querying \"Unclaimed Balances\".")
	var unclaimed *ubalances.UnclaimedBalances

	queryUnclaimed := ubalancesbq.GetDataQuery(begin, end)
	resultsUnclaimed, err := bq.ExecuteQuery(ctx, queryUnclaimed)
	if err != nil {
		msg := fmt.Sprintf("Error querying Unclaimed Balances: %v", err)
		return unclaimed, errors.New(msg)
	}
	balancesSlices, errSlice := ubalancesbq.StructToSlice(resultsUnclaimed)
	if errSlice != nil {
		return unclaimed, errors.New(errSlice.Error())
	}

	unclaimed = ubalances.NewUnclaimedBalances(balancesSlices)

	if unclaimed.IsEmpty() {
		msg := fmt.Sprintf("Not found any Unclaimed Balances for the period between %s and %s", begin, end)
		debug.NewMessage(msg)
	}

	debug.NewMessage("Finish querying \"Unclaimed Balances\".")

	return unclaimed, nil
}

func getRewards(ctx context.Context, bq bigqueryutils.BigQueryWrapper, begin time.Time, end time.Time) (*rewards.Rewards, error) {
	debug.NewMessage("Start querying \"Delegation Rewards\".")

	var result *rewards.Rewards

	queryrewards := rewardsbq.GetDataQuery(begin, end)
	resultsrewards, err := bq.ExecuteQuery(ctx, queryrewards)
	if err != nil {
		msg := fmt.Sprintf("Error querying Delegation Rewards: %v", err)
		return result, errors.New(msg)
	}
	rewardsSlices, errSlice := rewardsbq.StructToSlice(resultsrewards)
	if errSlice != nil {
		return result, errors.New(errSlice.Error())
	}

	result = rewards.NewRewards(rewardsSlices)

	debug.NewMessage("Finish querying \"Delegation Rewards\".")

	return result, nil
}

func getStatuses(ctx context.Context, bq bigqueryutils.BigQueryWrapper, begin time.Time, end time.Time) (*operationsstatuses.OperationsStatuses, error) {
	name := "Operations Statuses"
	debug.NewMessage(fmt.Sprintf("Start querying \"%s\".", name))

	var result *operationsstatuses.OperationsStatuses

	queryOpStatuses := operationsstatusesbq.GetDataQuery(begin, end)
	resultsStatuses, err := bq.ExecuteQuery(ctx, queryOpStatuses)
	if err != nil {
		msg := fmt.Sprintf("Error querying %s: %v", name, err)
		return result, errors.New(msg)
	}
	opStatusesSlices, errSlice := operationsstatusesbq.StructToSlice(resultsStatuses)
	if errSlice != nil {
		return result, errors.New(errSlice.Error())
	}

	result = operationsstatuses.NewOperationsStatuses(opStatusesSlices)

	debug.NewMessage(fmt.Sprintf("Finish querying \"%s\".", name))

	return result, nil
}

func getProjectId() string {
	projectId := os.Getenv(env_project_id)
	if projectId == "" {
		msg := fmt.Sprintf("%s env variable not found or empty", env_project_id)
		debug.NewMessage(msg)
		projectId = "anchorage-playground"
	}

	return projectId
}
