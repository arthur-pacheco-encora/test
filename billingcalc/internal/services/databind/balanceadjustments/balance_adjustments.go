package balanceadjustments

import (
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

const (
	ColOrganization      databind.Column = 0
	ColAccount           databind.Column = 1
	ColBusinessDay       databind.Column = 4
	ColOperation         databind.Column = 5
	ColAsset             databind.Column = 6
	ColTotalUSD          databind.Column = 16
	ColAssetQuantity     databind.Column = 17
	ColStakingAdjustment databind.Column = 18
)

type BalanceAdjustments struct {
	Organizations map[string]Organization
}

type Organization struct {
	Name     string
	Accounts map[string]Account
}

type Account struct {
	Name   string
	Assets map[string]Asset
}

type Asset struct {
	Name                    string
	DailyBalanceAdjustments map[time.Time]DailyBalanceAdjustments
}

type DailyBalanceAdjustments struct {
	OperationType     string
	UsdValue          decimal.Decimal
	BusinessDay       time.Time
	StakingAdjustment string
}

func (d *DailyBalanceAdjustments) HasAdjustment() bool {
	return d.StakingAdjustment == "Y"
}

func isEmptyField(row []string, requiredFields []databind.Column) bool {
	for _, field := range requiredFields {
		if row[field] == "" {
			return true
		}
	}
	return false
}

func NewBalanceAdjustments(table [][]string) *BalanceAdjustments {
	balanceAdjustments := &BalanceAdjustments{
		Organizations: make(map[string]Organization),
	}

	requiredFields := []databind.Column{ColOrganization, ColAccount, ColAsset, ColBusinessDay}
	for _, row := range table {
		if isEmptyField(row, requiredFields) {
			continue
		}
		organizationName := sanitization.SanitizeName(row[ColOrganization])
		accountName := sanitization.SanitizeName(row[ColAccount])
		assetName := row[ColAsset]
		businessDay := row[ColBusinessDay]

		org, exists := balanceAdjustments.Organizations[organizationName]
		if !exists {
			org = Organization{
				Name:     organizationName,
				Accounts: make(map[string]Account),
			}
		}

		acc, accExists := org.Accounts[accountName]
		if !accExists {
			acc = Account{
				Name:   accountName,
				Assets: make(map[string]Asset),
			}
		}

		asset, assetExists := acc.Assets[assetName]
		if !assetExists {
			asset = Asset{
				Name:                    assetName,
				DailyBalanceAdjustments: make(map[time.Time]DailyBalanceAdjustments),
			}
		}

		DailyBalanceAdjustments := DailyBalanceAdjustments{
			UsdValue:          converter.FromStrToDecimal(row[ColTotalUSD]),
			BusinessDay:       converter.FromStrToTimeStamp(businessDay),
			OperationType:     row[ColOperation],
			StakingAdjustment: row[ColStakingAdjustment],
		}

		asset.DailyBalanceAdjustments[DailyBalanceAdjustments.BusinessDay] = DailyBalanceAdjustments
		acc.Assets[assetName] = asset
		org.Accounts[accountName] = acc
		balanceAdjustments.Organizations[organizationName] = org
	}
	return balanceAdjustments
}

func (b *Asset) GetSortedDailyBalanceAdjustments() []DailyBalanceAdjustments {
	keys := make([]time.Time, 0, len(b.DailyBalanceAdjustments))
	for date := range b.DailyBalanceAdjustments {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	sortedBalances := make([]DailyBalanceAdjustments, 0, len(b.DailyBalanceAdjustments))
	for _, key := range keys {
		sortedBalances = append(sortedBalances, b.DailyBalanceAdjustments[key])
	}
	return sortedBalances
}

func (b *BalanceAdjustments) SumDailyBalanceAdjustments(orgName string, accName string, assetName string) decimal.Decimal {
	var ok bool

	if len(b.Organizations) == 0 {
		return decimal.Zero
	}

	var orgAdjustment Organization
	if orgAdjustment, ok = b.Organizations[orgName]; !ok {
		return decimal.Zero
	}

	var accAdjustment Account
	if accAdjustment, ok = orgAdjustment.Accounts[accName]; !ok {
		return decimal.Zero
	}

	var asset Asset
	if asset, ok = accAdjustment.Assets[assetName]; !ok {
		return decimal.Zero
	}

	totalAdjustment := decimal.Zero
	for _, adjustment := range asset.DailyBalanceAdjustments {
		if adjustment.HasAdjustment() {
			totalAdjustment = totalAdjustment.Add(adjustment.UsdValue)
		}
	}

	return totalAdjustment
}
