package ubalances

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

const (
	ColAsset            databind.Column = 3
	ColBalanceType      databind.Column = 5
	ColDailyBalanceStr  databind.Column = 6
	ColDailyBalanceDate databind.Column = 7
	ColAccountName      databind.Column = 8
	ColUsdPrice         databind.Column = 9
	ColUsdValue         databind.Column = 10
)

type UnclaimedBalances struct {
	Accounts map[string]Account
}

type Account struct {
	Name        string
	DisplayName string
	Assets      map[string]Asset
}

type Asset struct {
	Asset         string
	DailyBalances map[time.Time]DailyBalance
}
type DailyBalance struct {
	BalanceType          string
	DailyBalanceStr      decimal.Decimal
	DailyBalanceDateTime time.Time
	UsdPrice             decimal.Decimal
	UsdValue             decimal.Decimal
}

func (u *UnclaimedBalances) GetDailyBalances(account, asset string) map[time.Time]DailyBalance {
	acc := u.Accounts[account]
	assets := acc.Assets[asset]
	return assets.DailyBalances
}

func (u *UnclaimedBalances) GetSortedDailyBalances(account, asset string) []DailyBalance {
	acc := u.Accounts[account]
	assets := acc.Assets[asset]
	balances := assets.DailyBalances

	keys := make([]time.Time, 0, len(balances))
	for key := range balances {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	result := make([]DailyBalance, 0)
	for _, key := range keys {
		result = append(result, balances[key])
	}
	return result
}

func (u *UnclaimedBalances) IsEmpty() bool {
	return len(u.Accounts) == 0
}

func NewUnclaimedBalances(table [][]string) *UnclaimedBalances {
	uncBal := &UnclaimedBalances{
		Accounts: make(map[string]Account),
	}

	for _, row := range table {
		accountName := sanitization.SanitizeName(row[ColAccountName])
		assetName := row[ColAsset]
		balanceDay := converter.FromStrToTimeStamp(row[ColDailyBalanceDate])

		acc, exists := uncBal.Accounts[accountName]
		if !exists {
			acc = Account{
				Name:        accountName,
				DisplayName: row[ColAccountName],
				Assets:      make(map[string]Asset),
			}
		}

		asset, exists := acc.Assets[assetName]
		if !exists {
			asset = Asset{
				Asset:         assetName,
				DailyBalances: make(map[time.Time]DailyBalance, 0),
			}
		}

		dailyBalanceStr := converter.FromStrToDecimal(row[ColDailyBalanceStr])
		usdPrice := converter.FromStrToDecimal(row[ColUsdPrice])

		balance, exists := asset.DailyBalances[balanceDay]
		if exists {
			balance.DailyBalanceStr = balance.DailyBalanceStr.Add(dailyBalanceStr)
			balance.UsdValue = balance.UsdValue.Add(dailyBalanceStr.Mul(usdPrice))
		} else {
			balance = DailyBalance{
				BalanceType:          row[ColBalanceType],
				DailyBalanceStr:      dailyBalanceStr,
				DailyBalanceDateTime: balanceDay,
				UsdPrice:             usdPrice,
				UsdValue:             converter.FromStrToDecimal(row[ColUsdValue]),
			}
		}

		asset.DailyBalances[balanceDay] = balance
		acc.Assets[assetName] = asset
		uncBal.Accounts[accountName] = acc
	}

	return uncBal
}
