// package dailybalances provides an abstraction layer for reading values from a
// "Daily Balances Report" spreadsheet.
//
// Notice that even though it has "daily" in its name, the spreadsheet actually
// contains only the average daily balance from a given period of time.
package dailybalances

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

const (
	colMsaID                      databind.Column = 0
	colOrgName                    databind.Column = 1
	colAccountName                databind.Column = 2
	colAnchorEntity               databind.Column = 3
	colOrgId                      databind.Column = 4
	colAssetName                  databind.Column = 5
	colAccountId                  databind.Column = 6
	colDailyAssetTotal            databind.Column = 7
	colDailyAssetPrice            databind.Column = 8
	colDailyUsdTotal              databind.Column = 9
	colDailyTotalFromAddresses    databind.Column = 10
	colUnclaimedRewardsBalanceUsd databind.Column = 11
	colTotalAucUsd                databind.Column = 12
	colIsGraduated                databind.Column = 13
)

type DailyBalance struct {
	organizations map[string]organization
}

type organization struct {
	accounts map[string]account
}

type account struct {
	balances map[string]Balance
}

type Balance struct {
	AssetBalance               decimal.Decimal
	UsdBalance                 decimal.Decimal
	UnclaimedRewardsBalanceUsd decimal.Decimal
	TotalAucUsd                decimal.Decimal
}

func NewDailyBalance(table [][]string) *DailyBalance {
	db := DailyBalance{
		organizations: make(map[string]organization),
	}

	for _, row := range table {
		msaID := row[colMsaID]
		accountName := sanitization.SanitizeName(row[colAccountName])
		assetName := row[colAssetName]

		if _, ok := db.organizations[msaID]; !ok {
			db.organizations[msaID] = organization{
				accounts: make(map[string]account),
			}
		}
		if _, ok := db.organizations[msaID].accounts[accountName]; !ok {
			db.organizations[msaID].accounts[accountName] = account{
				balances: make(map[string]Balance),
			}
		}

		if _, ok := db.organizations[msaID].accounts[accountName].balances[assetName]; !ok {
			db.organizations[msaID].accounts[accountName].balances[assetName] = Balance{}
		}

		prevAssetBalance := db.organizations[msaID].accounts[accountName].balances[assetName].AssetBalance
		prevUsdBalance := db.organizations[msaID].accounts[accountName].balances[assetName].UsdBalance
		prevUnclaimedRewardsBalanceUsd := db.organizations[msaID].accounts[accountName].balances[assetName].UnclaimedRewardsBalanceUsd
		prevTotalAucUsd := db.organizations[msaID].accounts[accountName].balances[assetName].TotalAucUsd

		db.organizations[msaID].accounts[accountName].balances[assetName] = Balance{
			AssetBalance:               prevAssetBalance.Add(converter.FromStrToDecimal(row[colDailyAssetTotal])),
			UsdBalance:                 prevUsdBalance.Add(converter.FromStrToDecimal(row[colDailyUsdTotal])),
			UnclaimedRewardsBalanceUsd: prevUnclaimedRewardsBalanceUsd.Add(converter.FromStrToDecimal(row[colUnclaimedRewardsBalanceUsd])),
			TotalAucUsd:                prevTotalAucUsd.Add(converter.FromStrToDecimal(row[colTotalAucUsd])),
		}
	}

	return &db
}

func (db *DailyBalance) GetAssetBalance(msaID, accountName, assetName string) (decimal.Decimal, error) {
	if err := checkBalanceEntry(db, msaID, accountName, assetName); err != nil {
		return decimal.New(0, 1), err
	}
	return db.organizations[msaID].accounts[accountName].balances[assetName].AssetBalance, nil
}

func (db *DailyBalance) GetUsdBalance(msaID, accountName, assetName string) (decimal.Decimal, error) {
	if err := checkBalanceEntry(db, msaID, accountName, assetName); err != nil {
		return decimal.New(0, 1), err
	}
	return db.organizations[msaID].accounts[accountName].balances[assetName].UsdBalance, nil
}

func (db *DailyBalance) GetUnclaimedRewardsBalanceUsd(msaID, accountName, assetName string) (decimal.Decimal, error) {
	if err := checkBalanceEntry(db, msaID, accountName, assetName); err != nil {
		return decimal.New(0, 1), err
	}
	return db.organizations[msaID].accounts[accountName].balances[assetName].UnclaimedRewardsBalanceUsd, nil
}

func (db *DailyBalance) GetTotalAucUsd(msaID, accountName, assetName string) (decimal.Decimal, error) {
	if err := checkBalanceEntry(db, msaID, accountName, assetName); err != nil {
		return decimal.New(0, 1), err
	}
	return db.organizations[msaID].accounts[accountName].balances[assetName].TotalAucUsd, nil
}

func (db *DailyBalance) GetAccountBalances(msaID, accountName string) (map[string]Balance, error) {
	if err := checkOrgAndAccountEntry(db, msaID, accountName); err != nil {
		return nil, err
	}
	return db.organizations[msaID].accounts[accountName].balances, nil
}

func (db *DailyBalance) GetAverageUsdBalanceByOrg(msaID string, daysInMonth int) (decimal.Decimal, error) {
	if daysInMonth <= 0 {
		return decimal.Zero, errors.New("days in month must be greater than 0")
	}

	org, ok := db.organizations[msaID]
	if !ok {
		return decimal.Zero, errors.New(fmt.Sprintf("organization not found - org:%s", msaID))
	}

	var totalUsdBalance decimal.Decimal
	for _, acc := range org.accounts {
		for _, bal := range acc.balances {
			totalUsdBalance = totalUsdBalance.Add(bal.UsdBalance)
		}
	}

	if totalUsdBalance.IsZero() {
		return decimal.Zero, nil
	}
	averageUsdBalance := totalUsdBalance.DivRound(decimal.NewFromInt(int64(daysInMonth)), 16)
	return averageUsdBalance, nil
}

func checkOrgAndAccountEntry(db *DailyBalance, msaID, accountName string) error {
	if _, ok := db.organizations[msaID]; !ok {
		errMsg := fmt.Sprintf("Organization not found - org:%s", msaID)
		return errors.New(errMsg)
	}
	if _, ok := db.organizations[msaID].accounts[accountName]; !ok {
		errMsg := fmt.Sprintf("Account not found - msaID:%s acc:%s", msaID, accountName)
		return errors.New(errMsg)
	}

	return nil
}

func checkBalanceEntry(db *DailyBalance, msaID, accountName, assetName string) error {
	if err := checkOrgAndAccountEntry(db, msaID, accountName); err != nil {
		return err
	}

	if _, ok := db.organizations[msaID].accounts[accountName].balances[assetName]; !ok {
		errMsg := fmt.Sprintf("Account not found - org:%s acc:%s asset:%s", msaID, accountName, assetName)
		return errors.New(errMsg)
	}

	return nil
}
