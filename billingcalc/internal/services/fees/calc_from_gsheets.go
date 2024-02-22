package fees

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/googlesheetsutils"
)

func CalculateFromGSheets(ctx context.Context, sheetId string, mfrTab string, rewardsTab string, uBalancesTab string, balanceAdjustmentsTab string, opStatusesTab string, dailyBalancesTab string, token string, firstExternalId int, invoiceDate time.Time) (*CalculatedFees, error) {
	gSheeetRequest := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var stakingSummary, custodySummary, combinedSummary StakingSummary
	var combinedWarnings []Warning
	var err1, err2 error

	// Processing Mfr
	mfrObj, err := mfr.ProcessMfr(ctx, nil, sheetId, mfrTab, token)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Processing Rewards
	rewardsParams := googlesheetsutils.GetDatabindFromSheetTabParams[*rewards.Rewards]{
		SheetRequest: gSheeetRequest,
		TabName:      rewardsTab,
		ReportName:   "Delegation and Staking Rewards",
		HeaderRow:    8,
		Parser:       rewards.NewRewards,
	}
	rewards, err := googlesheetsutils.GetDatabindFromSheetTab[*rewards.Rewards](ctx, rewardsParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Processing Unclaimed Balances
	balancesParams := googlesheetsutils.GetDatabindFromSheetTabParams[*ubalances.UnclaimedBalances]{
		SheetRequest: gSheeetRequest,
		TabName:      uBalancesTab,
		ReportName:   "Unclaimed Balances Report",
		HeaderRow:    1,
		Parser:       ubalances.NewUnclaimedBalances,
	}
	balances, err := googlesheetsutils.GetDatabindFromSheetTab[*ubalances.UnclaimedBalances](ctx, balancesParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Processing Balance Adjustments
	balanceAdjustmentsParams := googlesheetsutils.GetDatabindFromSheetTabParams[*balanceadjustments.BalanceAdjustments]{
		SheetRequest: gSheeetRequest,
		TabName:      balanceAdjustmentsTab,
		ReportName:   "Balance Adjustments Report",
		HeaderRow:    1,
		Parser:       balanceadjustments.NewBalanceAdjustments,
	}
	balanceAdjustments, err := googlesheetsutils.GetDatabindFromSheetTab[*balanceadjustments.BalanceAdjustments](ctx, balanceAdjustmentsParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Processing Operations Statuses
	opStatusesParams := googlesheetsutils.GetDatabindFromSheetTabParams[*operationsstatuses.OperationsStatuses]{
		SheetRequest: gSheeetRequest,
		TabName:      opStatusesTab,
		ReportName:   "Operations Statuses Report",
		HeaderRow:    1,
		Parser:       operationsstatuses.NewOperationsStatuses,
	}
	opStatuses, err := googlesheetsutils.GetDatabindFromSheetTab[*operationsstatuses.OperationsStatuses](ctx, opStatusesParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	dailyBalParams := googlesheetsutils.GetDatabindFromSheetTabParams[*dailybalances.DailyBalance]{
		SheetRequest: gSheeetRequest,
		TabName:      dailyBalancesTab,
		ReportName:   "Daily Balances Report",
		HeaderRow:    1,
		Parser:       dailybalances.NewDailyBalance,
	}
	dailyBalances, err := googlesheetsutils.GetDatabindFromSheetTab[*dailybalances.DailyBalance](ctx, dailyBalParams)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		summary, warn := CalculateStakingFees(mfrObj, rewards, balances, balanceAdjustments, opStatuses, dailyBalances, firstExternalId, invoiceDate)
		mu.Lock()
		defer mu.Unlock()
		if summary != nil {
			combinedSummary = append(combinedSummary, summary...)
			combinedWarnings = append(combinedWarnings, warn...)
		} else {
			err1 = errors.New("error in CalculateStakingFees")
		}
	}()

	go func() {
		defer wg.Done()
		custodySummary, custodyWarn := CalculateCustodyFees(mfrObj, rewards, balances, balanceAdjustments, opStatuses, dailyBalances, firstExternalId, invoiceDate)
		mu.Lock()
		defer mu.Unlock()
		if custodySummary != nil {
			combinedSummary = append(combinedSummary, custodySummary...)
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
