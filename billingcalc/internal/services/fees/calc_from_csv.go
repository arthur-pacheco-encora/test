package fees

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/file"
)

func CalculateFromCsv(ctx context.Context, mfrFile, rewardsFile, unclaimedFile, balanceAdjustmentsFile io.Reader, opStatusesFile io.Reader, dailyBalancesFile io.Reader, firstExternalId int, invoiceDate time.Time) (*CalculatedFees, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var stakingSummary, custodySummary StakingSummary
	var combinedWarnings []Warning
	var err1, err2 error

	mfr, err := mfr.ProcessMfr(ctx, mfrFile, "", "", "")
	if err != nil {
		return nil, errors.New(err.Error())
	}

	rewardsParams := file.GetDatabindFromDataParams[*rewards.Rewards]{
		ReportName: "Delegation and Staking Rewards Activity Report",
		HeaderRow:  8,
		File:       rewardsFile,
		Parser:     rewards.NewRewards,
	}
	rewards, err := file.GetDatabindFromData[*rewards.Rewards](rewardsParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	balancesParams := file.GetDatabindFromDataParams[*ubalances.UnclaimedBalances]{
		ReportName: "Unclaimed Balances Report",
		HeaderRow:  1,
		File:       unclaimedFile,
		Parser:     ubalances.NewUnclaimedBalances,
	}
	balances, err := file.GetDatabindFromData[*ubalances.UnclaimedBalances](balancesParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	balanceadjustmentsParams := file.GetDatabindFromDataParams[*balanceadjustments.BalanceAdjustments]{
		ReportName: "Balance Adjustments Report",
		HeaderRow:  1,
		File:       balanceAdjustmentsFile,
		Parser:     balanceadjustments.NewBalanceAdjustments,
	}
	balanceadjustments, err := file.GetDatabindFromData[*balanceadjustments.BalanceAdjustments](balanceadjustmentsParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	opStatusesParams := file.GetDatabindFromDataParams[*operationsstatuses.OperationsStatuses]{
		ReportName: "Client Operations Statuses Report",
		HeaderRow:  1,
		File:       opStatusesFile,
		Parser:     operationsstatuses.NewOperationsStatuses,
	}
	opStatuses, err := file.GetDatabindFromData[*operationsstatuses.OperationsStatuses](opStatusesParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	dailyBalancesParams := file.GetDatabindFromDataParams[*dailybalances.DailyBalance]{
		ReportName: "Daily Balances Report",
		HeaderRow:  1,
		File:       dailyBalancesFile,
		Parser:     dailybalances.NewDailyBalance,
	}
	dailyBalances, err := file.GetDatabindFromData[*dailybalances.DailyBalance](dailyBalancesParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	wg.Add(2)

	go func() {
		defer wg.Done()
		summary, stakingWarn := CalculateStakingFees(mfr, rewards, balances, balanceadjustments, opStatuses, dailyBalances, firstExternalId, invoiceDate)
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
		summary, custodyWarn := CalculateCustodyFees(mfr, rewards, balances, balanceadjustments, opStatuses, dailyBalances, firstExternalId, invoiceDate)
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
