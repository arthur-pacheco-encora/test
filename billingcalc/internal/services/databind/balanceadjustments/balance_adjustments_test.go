//go:build !selectTest || unitTest

package balanceadjustments_test

import (
	"testing"
	"time"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewBalanceAdjustments(t *testing.T) {
	table := [][]string{
		{"Org1", "Account1", "", "", "2023-01-01", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Yes"},
		{"Org1", "Account1", "", "", "2023-01-02", "Operation2", "Asset1", "", "", "", "", "", "", "", "", "", "150.00", "10", "No"},
	}

	balanceAdjustments := balanceadjustments.NewBalanceAdjustments(table)

	assert.NotNil(t, balanceAdjustments.Organizations["Org1"])
	assert.NotNil(t, balanceAdjustments.Organizations["Org1"].Accounts["Account1"])
	assert.NotNil(t, balanceAdjustments.Organizations["Org1"].Accounts["Account1"].Assets["Asset1"])
}

func TestGetSortedDailyBalanceAdjustments(t *testing.T) {
	asset := balanceadjustments.Asset{
		Name: "Asset1",
		DailyBalanceAdjustments: map[time.Time]balanceadjustments.DailyBalanceAdjustments{
			time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC): {OperationType: "Operation2", UsdValue: decimal.NewFromFloat(150.00)},
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC): {OperationType: "Operation1", UsdValue: decimal.NewFromFloat(100.00)},
		},
	}

	sortedBalances := asset.GetSortedDailyBalanceAdjustments()

	assert.Equal(t, "Operation1", sortedBalances[0].OperationType)
	assert.Equal(t, "Operation2", sortedBalances[1].OperationType)
}

func TestNewBalanceAdjustments_DuplicateEntries(t *testing.T) {
	table := [][]string{
		{"Org1", "Account1", "", "", "2023-01-01", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Yes"},
		{"Org1", "Account1", "", "", "2023-01-01", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Yes"},
	}

	balanceAdjustments := balanceadjustments.NewBalanceAdjustments(table)

	assert.Equal(t, 1, len(balanceAdjustments.Organizations["Org1"].Accounts["Account1"].Assets["Asset1"].DailyBalanceAdjustments))
}

func TestNewBalanceAdjustments_MultipleOrgsAndAccounts(t *testing.T) {
	table := [][]string{
		{"Org1", "Account1", "", "", "2023-01-01", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Yes"},
		{"Org2", "Account2", "", "", "2023-01-02", "Operation2", "Asset2", "", "", "", "", "", "", "", "", "", "200.00", "10", "No"},
	}

	balanceAdjustments := balanceadjustments.NewBalanceAdjustments(table)

	assert.NotNil(t, balanceAdjustments.Organizations["Org1"])
	assert.NotNil(t, balanceAdjustments.Organizations["Org2"])
	assert.NotNil(t, balanceAdjustments.Organizations["Org1"].Accounts["Account1"])
	assert.NotNil(t, balanceAdjustments.Organizations["Org2"].Accounts["Account2"])
}

func TestHasAdjustment(t *testing.T) {
	asset := balanceadjustments.Asset{
		Name: "Asset1",
		DailyBalanceAdjustments: map[time.Time]balanceadjustments.DailyBalanceAdjustments{
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC): {StakingAdjustment: "Y"},
			time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC): {StakingAdjustment: "N"},
		},
	}

	hasAdjustment := false
	for _, adj := range asset.DailyBalanceAdjustments {
		if adj.HasAdjustment() {
			hasAdjustment = true
			break
		}
	}

	assert.True(t, hasAdjustment, "Asset should have staking adjustment")

	assetWithoutStaking := balanceadjustments.Asset{
		Name: "Asset2",
		DailyBalanceAdjustments: map[time.Time]balanceadjustments.DailyBalanceAdjustments{
			time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC): {StakingAdjustment: "N"},
			time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC): {StakingAdjustment: "N"},
		},
	}

	hasAdjustment = false
	for _, adj := range assetWithoutStaking.DailyBalanceAdjustments {
		if adj.HasAdjustment() {
			hasAdjustment = true
			break
		}
	}

	assert.False(t, hasAdjustment, "Asset should not have staking adjustment")
}

func TestNewBalanceAdjustmentsWithEmptyRequiredField(t *testing.T) {
	table := [][]string{
		{"Organization1", "Account1", "", "", "", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Y"},
		{"Organization2", "Account1", "", "", "2023-01-01", "Operation2", "Asset2", "", "", "", "", "", "", "", "", "", "100.00", "5", "Y"},
	}

	balanceAdjustments := balanceadjustments.NewBalanceAdjustments(table)

	_, org1Exists := balanceAdjustments.Organizations["Organization1"]
	_, org2Exists := balanceAdjustments.Organizations["Organization2"]

	assert.False(t, org1Exists, "Organization1 should not exist due to empty required field")
	assert.True(t, org2Exists, "Organization2 should exist")
}

func TestSumDailyBalanceAdjustments(t *testing.T) {
	table := [][]string{
		{"Org1", "Account1", "", "", "2023-01-01", "Operation1", "Asset1", "", "", "", "", "", "", "", "", "", "100.00", "5", "Y"},
		{"Org1", "Account1", "", "", "2023-01-02", "Operation2", "Asset1", "", "", "", "", "", "", "", "", "", "150.00", "10", "N"},
		{"Org1", "Account1", "", "", "2023-01-03", "Operation2", "Asset1", "", "", "", "", "", "", "", "", "", "200.00", "10", "Y"},
		{"Org1", "Account1", "", "", "2023-01-03", "Operation2", "Asset2", "", "", "", "", "", "", "", "", "", "150.00", "10", "N"},
	}

	balanceAdjustments := balanceadjustments.NewBalanceAdjustments(table)

	expected := decimal.NewFromFloat(300.00)
	result := balanceAdjustments.SumDailyBalanceAdjustments("Org1", "Account1", "Asset1")
	assert.True(t, expected.Equal(result), "SumDailyBalanceAdjustments should correctly sum the daily balance adjustments")
}
