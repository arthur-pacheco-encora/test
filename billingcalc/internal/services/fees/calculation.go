package fees

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees/custody"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/static"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/date"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

var EntityIDToAcronym = map[string]string{
	"15": "ADB",
	"33": "ABS",
}

func CalculateStakingFees(mfr *mfr.MasterFeeRates, rwd *rewards.Rewards, ubal *ubalances.UnclaimedBalances, balAdj *balanceadjustments.BalanceAdjustments, ops *operationsstatuses.OperationsStatuses, dayBal *dailybalances.DailyBalance, firstExternalId int, invoiceDate time.Time) (StakingSummary, []Warning) {
	msg := "Start of Staking Fees calculation."
	debug.NewMessage(msg)

	var calcTable CalcTable
	var summary StakingSummary
	var warnings []Warning

	firstAcc := true
	currentExternalID := firstExternalId
	getInvoiceNumber := GenInvoiceNumber(firstExternalId)

	err := setCalcTable(&calcTable)
	if err != nil {
		msg := fmt.Sprintf("Failed in Calc Table: %v", err)
		debug.NewMessage(msg)
		panic(err)
	}

	for _, organization := range mfr.GetOrganizations() {
		accResults := []AccountResult{}
		debug.NewMessage(fmt.Sprintf("Processing organization %s id:%s", organization.DisplayName, organization.Id))
		for _, mfrAccount := range mfr.GetAccounts(organization.Id) {
			debug.NewMessage(fmt.Sprintf("MFR account name: %s", mfrAccount.Name))
			var out []StakingOutput

			currentExternalID = GenExternalID(currentExternalID, firstAcc)
			firstAcc = false
			dueDate := GenDueDate(invoiceDate, mfrAccount.BillingTerms)

			rwdAccount := rwd.GetAccountById(mfrAccount.Id)
			if rwdAccount.Name == "" {
				debug.NewMessage(fmt.Sprintf("Skipping because no entry in Rewards sheet was found for ID %s", mfrAccount.Id))
			}

			for _, rwdAsset := range rwdAccount.GetAssets() {
				fee := mfrAccount.GetAssetStakingFees(rwdAsset.Name)
				if fee.AssetName == "" {
					warnMsg := fmt.Sprintf("MFR entry not found for account %v, asset %v.", rwdAccount.Id, rwdAsset.Name)
					warnings = addWarning(organization.Name, rwdAccount.Name, rwdAsset.Name, warnMsg, warnings)
					debug.NewMessage(warnMsg)
				}

				balAdjUsdValue := balAdj.SumDailyBalanceAdjustments(organization.Name, mfrAccount.Name, rwdAsset.Name)

				claimedRewards := rwdAsset.GetClaimedRewards()
				unclaimedBalances := ubal.GetSortedDailyBalances(rwdAccount.Name, rwdAsset.Name)

				operationsStatus := ops.GetStatusByAssetAndDate(rwdAccount.Name, rwdAsset.Name, invoiceDate)

				appendStakingOutput(&out, calcTable, rwdAccount.Name, rwdAsset.Name, fee, claimedRewards, unclaimedBalances, operationsStatus, invoiceDate, balAdjUsdValue)
			}

			customerID := mfrAccount.CustomerId
			entityID := organization.EntityId
			if entityID == "" {
				msg := "Failed to fetch EntityID from: " + organization.Name
				debug.NewMessage(msg)
				log.Println(msg)
			} else {
				invoiceNumber := getInvoiceNumber(entityID, string(organization.Id))
				externalID := fmt.Sprintf("%d", currentExternalID)
				accResults = append(accResults, AccountResult{
					AccName:       mfrAccount.Name,
					BillingTerms:  mfrAccount.BillingTerms,
					CustomerID:    customerID,
					DisplayName:   mfrAccount.DisplayName,
					ExternalID:    externalID,
					EntityID:      entityID,
					InvoiceNumber: invoiceNumber,
					InvoiceDate:   invoiceDate.Format("01/02/2006"),
					DueDate:       dueDate.Format("01/02/2006"),
					Assets:        out,
				})
			}

			if accResultsJson, err := json.Marshal(accResults); err != nil {
				debug.NewMessage(fmt.Sprintf("fail: accResults json.Marshal %s", err.Error()))
			} else {
				debug.NewMessage(fmt.Sprintf("success: accResults json.Marshal %s", string(accResultsJson)))
			}
		}
		summary = append(summary, OrgResult{
			OrgName:  organization.Name,
			Accounts: accResults,
		})
	}

	debug.NewMessage("End of Staking Fees calculation.")
	if summaryJson, err := json.Marshal(summary); err != nil {
		debug.NewMessage(fmt.Sprintf("Failed to json.Marshal summary %s", err.Error()))
	} else {
		debug.NewMessage(string(summaryJson))
	}

	return summary, warnings
}

func Contains(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GetResultByOrgName(orgName string, result []OrgResult) OrgResult {
	for _, o := range result {
		if o.OrgName == orgName {
			return o
		}
	}

	return OrgResult{
		OrgName:  "",
		Accounts: make([]AccountResult, 0),
	}
}

func GenDueDate(invoiceDate time.Time, billingTerms string) time.Time {
	if billingTerms == "" {
		return invoiceDate
	}
	var billingTermsToCalc int
	_, err := fmt.Sscanf(billingTerms, "%d", &billingTermsToCalc)
	if err != nil {
		debug.NewMessage(err.Error())
		debug.NewMessage(fmt.Sprintf("billingTerms string: %s", billingTerms))
		panic("Error while calculating due date")
	}
	day := 24 * time.Hour
	return invoiceDate.Add(time.Duration(billingTermsToCalc) * day)
}

func GenExternalID(currentExternalID int, firstAcc bool) int {
	if !firstAcc {
		currentExternalID++
	}
	return currentExternalID
}

func GenInvoiceNumber(initialValue int) func(entityID, msaID string) string {
	value := initialValue
	invoiceNumberTable := make(map[string]int) // { "msaID": value }

	return func(entityID, msaID string) string {
		// Acronym
		acronym, ok := EntityIDToAcronym[entityID]
		if !ok {
			return "Acronym Not Found"
		}

		currentValue := value
		// Get the value related with the MSAID
		_, exist := invoiceNumberTable[msaID]
		if !exist {
			invoiceNumberTable[msaID] = currentValue
			value++
		} else {
			currentValue = invoiceNumberTable[msaID]
		}

		return acronym + "-" + strconv.Itoa(currentValue)
	}
}

func appendStakingOutput(out *[]StakingOutput, calcTable CalcTable, account string, asset string, stakingFee mfr.StakingFee, filteredRewards []rewards.ClaimedReward, dailyBalances []ubalances.DailyBalance, opStatus operationsstatuses.Status, invoiceDate time.Time, balAdjuUsdValue decimal.Decimal) {
	for _, entry := range calcTable {
		var earnedRewards decimal.Decimal
		var fee decimal.Decimal
		var amount decimal.Decimal
		monthlyRate := ""

		ItemCategory := "Delegation Rewards Fees"

		if entry.Asset != asset {
			continue
		}

		validator := AnchorageValidator
		if entry.Validator != "anchorage" {
			validator = NonAnchorageValidator
		}

		if entry.Claimable {
			if len(dailyBalances) == 0 {
				msg := fmt.Sprintf("Account %s delegation reward for %s does not relate to Anchorage staked balances.", account, asset)
				debug.NewMessage(msg)
				return
			}
			earnedRewards = sumDiffUnclaimedInUsd(dailyBalances, filteredRewards, validator)
		} else {
			earnedRewards = sumClaimedRewards(filteredRewards, validator)
		}

		switch true {
		case operationsstatuses.IsAssetFromExternalValidator(opStatus):
			calcAmountFromExternalValidator(opStatus, &ItemCategory, &earnedRewards, invoiceDate, stakingFee.ThirdPartyFee, &fee, &amount, &monthlyRate, balAdjuUsdValue)
		default:
			calcAmountDefault(stakingFee, entry, &earnedRewards, &fee, &amount, balAdjuUsdValue)
		}

		*out = append(*out, StakingOutput{
			ServiceType:             "Staking Fee",
			Asset:                   asset,
			Amount:                  amount,
			CollectedOnChainAlready: entry.On_chain,
			EarnedRewards:           earnedRewards.Round(2),
			FeeRates:                fee,
			ItemCategory:            ItemCategory,
			ItemDescription:         "",
			ItemQuantity:            "",
			Memo:                    "",
			MonthlyRate:             monthlyRate,
		})
	}
}

// Sum asset claimed rewards for a client
func sumClaimedRewards(claimedRewards []rewards.ClaimedReward, validator string) decimal.Decimal {
	earnedRewards := decimal.Zero

	for _, entry := range claimedRewards {
		var usdValue decimal.Decimal

		if validator == AnchorageValidator {
			usdValue = entry.AnchorageUsdValue
		} else {
			usdValue = entry.ThirdPartyUsdValue
		}

		earnedRewards = earnedRewards.Add(usdValue)
	}

	return earnedRewards
}

// Sum of daily USD value change in 'Historic Daily Blockchain Balances Balance Str'
// from 'Daily Historic Blockchain Balances' + sumif('Operations Total USD Reward Value' > 0)
// from 'Staking Rewards and Delegation Rewards'
func sumDiffUnclaimedInUsd(dailyBalances []ubalances.DailyBalance, claimedValues []rewards.ClaimedReward, validator string) decimal.Decimal {
	if len(dailyBalances) == 0 {
		return decimal.Zero
	}

	previousBalance := decimal.Zero
	sumDiffUsdValue := decimal.Zero
	first := true

	for _, balance := range dailyBalances {
		balanceDay := balance.DailyBalanceDateTime

		diffBalance := balance.DailyBalanceStr.Sub(previousBalance)

		if !first {
			claimedAssets := sumClaimedAssetsByDay(claimedValues, validator, balanceDay)

			total := diffBalance.Add(claimedAssets)

			usdValue := total.Mul(balance.UsdPrice)

			sumDiffUsdValue = sumDiffUsdValue.Add(usdValue)
		}

		previousBalance = balance.DailyBalanceStr
		first = false
	}

	return sumDiffUsdValue
}

// Sum asset claimed rewards by business day for a client
func sumClaimedAssetsByDay(claimedValues []rewards.ClaimedReward, validator string, date time.Time) decimal.Decimal {
	totalClaimed := decimal.Zero

	for _, value := range claimedValues {
		if value.BusinessDay == date {
			if validator == AnchorageValidator {
				totalClaimed = totalClaimed.Add(value.AnchorageAssetQty)
				continue
			}

			totalClaimed = totalClaimed.Add(value.ThirdPartyAssetQty)
		}
	}

	return totalClaimed
}

/* When clients stake to a validator that charges 100% commission, none of their earned rewards are sent to their Anchorage account. To ensure Anchorage earns revenue from these staking arrangements, there is an alternative fee charged in these scenarios, which is a percentage of the total average balance staked to the 100% validator during the month.
 */
func calcAmountFromExternalValidator(opStatus operationsstatuses.Status, ItemCategory *string, earnedRewards *decimal.Decimal, invoiceDate time.Time, thirdPartyValidatorFee decimal.Decimal, fee *decimal.Decimal, amount *decimal.Decimal, monthlyRate *string, balAdjuUsdValue decimal.Decimal) {
	*ItemCategory = "Delegation Rewards Fees - 100% validator"
	*earnedRewards, _ = decimal.NewFromString("0")
	*monthlyRate = fmt.Sprint(thirdPartyValidatorFee.Round(2).StringFixed(2), "%")
	monthsInYear := decimal.NewFromInt(12)

	activeDelegatedValue := opStatus.StatusesActiveDelegatedValue
	invoiceMonthDays := decimal.NewFromInt(date.NumberOfDaysInMonth(invoiceDate.Month(), invoiceDate.Year()))
	averageStakedBalance := activeDelegatedValue.DivRound(invoiceMonthDays, 16).Round(2)
	alternativeFeeRate := thirdPartyValidatorFee.DivRound(decimal.NewFromInt(100), 16)

	adjustedAmount := averageStakedBalance.Add(balAdjuUsdValue)
	*amount = adjustedAmount.Mul(alternativeFeeRate).DivRound(monthsInYear, 16).Round(2)
	*fee = averageStakedBalance
}

func calcAmountDefault(mfrFees mfr.StakingFee, entry CalcTableEntry, earnedRewards *decimal.Decimal, fee *decimal.Decimal, amount *decimal.Decimal, balAdjuUsdValue decimal.Decimal) {
	*fee = mfrFees.AnchorageFee
	if entry.Validator == "non_anchorage" {
		*fee = mfrFees.ThirdPartyFee
	}

	feeRate := fee.DivRound(decimal.NewFromInt(100), 16)
	if entry.On_chain {
		one := decimal.NewFromInt(1)
		feeRate = feeRate.DivRound(one.Sub(feeRate), 16)
	}

	adjustedAmount := earnedRewards.Add(balAdjuUsdValue)
	*amount = feeRate.Mul(adjustedAmount).Round(2)
}

func addWarning(organizationName, accountName, assetName, message string, warnings []Warning) []Warning {
	return append(warnings, Warning{
		OrgName:     organizationName,
		AccName:     accountName,
		Asset:       assetName,
		Description: message,
	})
}

func setCalcTable(calcTable *CalcTable) error {
	fileContent, err := static.Files.ReadFile("calc_table.json")
	if err != nil {
		return errors.New(fmt.Sprintf("Error reading file calc_table.json: %v", err))
	}

	if err := json.Unmarshal(fileContent, &calcTable); err != nil {
		return errors.New(fmt.Sprintf("Error in json unmarshal: %v", err))
	}

	return nil
}

func CalculateCustodyFees(mfr *mfr.MasterFeeRates, rwd *rewards.Rewards, ubal *ubalances.UnclaimedBalances, balAdj *balanceadjustments.BalanceAdjustments, ops *operationsstatuses.OperationsStatuses, dayBal *dailybalances.DailyBalance, firstExternalId int, invoiceDate time.Time) (StakingSummary, []Warning) {
	debug.NewMessage("Start of Custody Fees calculation.")
	var summary StakingSummary
	var warnings []Warning

	ItemCategory := "Custody Fee by Asset"
	daysInMonth := date.NumberOfDaysInTime(invoiceDate)
	firstAcc := true
	currentExternalID := firstExternalId
	getInvoiceNumber := GenInvoiceNumber(firstExternalId)

	for _, organization := range mfr.GetOrganizations() {
		accResults := []AccountResult{}
		debug.NewMessage(fmt.Sprintf("Processing organization %s id:%s", organization.DisplayName, organization.Id))

		aucOrgValue, _ := dayBal.GetAverageUsdBalanceByOrg(string(organization.Id), daysInMonth)
		if aucOrgValue.IsZero() {
			warnMsg := fmt.Sprintf("Skipping organization %s due to zero AUC value", organization.Name)
			debug.NewMessage(warnMsg)
		}

		for _, mfrAccount := range mfr.GetAccounts(organization.Id) {
			var custodyOutputs []StakingOutput
			debug.NewMessage(fmt.Sprintf("MFR account name: %s", mfrAccount.Name))

			rwdAccount := rwd.GetAccountById(mfrAccount.Id)
			if rwdAccount.Name == "" {
				debug.NewMessage(fmt.Sprintf("Skipping because no entry in Rewards sheet was found for ID %s", mfrAccount.Id))
			}

			currentExternalID = GenExternalID(currentExternalID, firstAcc)
			firstAcc = false
			dueDate := GenDueDate(invoiceDate, mfrAccount.BillingTerms)

			accountBalances, err := dayBal.GetAccountBalances(string(organization.Id), mfrAccount.Name)
			if err != nil {
				warnMsg := fmt.Sprintf("Error getting account balances: %s", err.Error())
				debug.NewMessage(warnMsg)
			}

			if len(accountBalances) == 0 {
				msg := fmt.Sprintf("Account Balances not found in Daily Balances for this account: %s", mfrAccount.Name)
				debug.NewMessage(msg)
			}

			for _, assetType := range mfr.GetAssetTypes(organization.Id, mfrAccount.Id) {
				avgAucAssets, err := custody.AvgAucByAsset(accountBalances, int(assetType.Id), int64(daysInMonth)) // B
				if err != nil {
					warnMsg := fmt.Sprintf("Error calculating AvgAucAsset: %s", err.Error())
					debug.NewMessage(warnMsg)
				}

				if len(avgAucAssets) == 0 {
					msg := fmt.Sprintf("Assets not found in the list of assets for this AssetID: %d", assetType.Id)
					debug.NewMessage(msg)
				}

				tierRate, err := assetType.FindAllTiers(aucOrgValue)
				if err != nil {
					warnMsg := fmt.Sprintf("Error finding tiers: %s", err.Error())
					debug.NewMessage(warnMsg)
				}

				for assetName, avgAucAsset := range avgAucAssets {
					feeAmount := custody.CalcEffectiveFeeAmount(tierRate, assetType.MinimumFee, aucOrgValue).Round(2)

					billedAmount := avgAucAsset.DivRound(aucOrgValue, 2).Mul(feeAmount).Round(2)
					monthlyRate := feeAmount.DivRound(decimal.NewFromInt(12), 2).Round(2)
					custodyOutputs = append(custodyOutputs, StakingOutput{
						ServiceType:     "Custody Fee",
						Asset:           assetName,
						Amount:          billedAmount,
						EarnedRewards:   avgAucAsset,
						FeeRates:        feeAmount,
						ItemCategory:    ItemCategory,
						ItemDescription: "",
						ItemQuantity:    "",
						Memo:            "",
						MonthlyRate:     monthlyRate.StringFixed(2),
					})
				}
			}

			customerID := mfrAccount.CustomerId
			entityID := organization.EntityId
			if entityID == "" {
				msg := "Failed to fetch EntityID from: " + organization.Name
				debug.NewMessage(msg)
				log.Println(msg)
			} else {
				invoiceNumber := getInvoiceNumber(entityID, string(organization.Id))
				externalID := fmt.Sprintf("%d", currentExternalID)
				accResults = append(accResults, AccountResult{
					AccName:       mfrAccount.Name,
					BillingTerms:  mfrAccount.BillingTerms,
					CustomerID:    customerID,
					DisplayName:   mfrAccount.DisplayName,
					ExternalID:    externalID,
					EntityID:      entityID,
					InvoiceNumber: invoiceNumber,
					InvoiceDate:   invoiceDate.Format("01/02/2006"),
					DueDate:       dueDate.Format("01/02/2006"),
					Assets:        custodyOutputs,
				})
			}
		}

		summary = append(summary, OrgResult{
			OrgName:  organization.Name,
			Accounts: accResults,
		})
	}

	debug.NewMessage("End of Custody Fees calculation.")
	if summaryJson, err := json.Marshal(summary); err != nil {
		debug.NewMessage(fmt.Sprintf("Failed to json.Marshal summary %s", err.Error()))
	} else {
		debug.NewMessage(string(summaryJson))
	}

	return summary, warnings
}
