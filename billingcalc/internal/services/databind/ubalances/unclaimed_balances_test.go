//go:build !selectTest || unitTest

package ubalances_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
)

func initBalance() ubalances.UnclaimedBalances {
	table := [][]string{
		{"", "", "Address1A", "AXL", "", "DELEGATION_REWARDS", "1.00", "2023-05-31 23:59:31", "Acc1", "0.50", "0.50"},
		{"", "", "Address1A", "AXL", "", "DELEGATION_REWARDS", "2.00", "2023-06-01 23:59:36", "Acc1", "2.50", "5.00"},
		{"", "", "Address1B", "AXL", "", "DELEGATION_REWARDS", "3.00", "2023-05-31 23:59:32", "Acc1", "0.50", "1.50"},
		{"", "", "Address1C", "ETH", "", "DELEGATION_REWARDS", "4.00", "2023-05-31 23:59:31", "Acc1", "1.00", "4.00"},
		{"", "", "Address2A", "HASH", "", "DELEGATION_REWARDS", "5.00", "2023-05-31 23:58:17", "Acc2", "2.00", "10.00"},
	}

	return *ubalances.NewUnclaimedBalances(table)
}

func TestLookups(t *testing.T) {
	balance := initBalance()

	t.Run("Test_GetDailyBalances", func(t *testing.T) {
		time1, _ := time.Parse("2006-01-02", "2023-05-31")
		time2, _ := time.Parse("2006-01-02", "2023-06-01")

		expected := map[time.Time]ubalances.DailyBalance{
			time1: {BalanceType: "DELEGATION_REWARDS", DailyBalanceStr: newDecimalFromString("4.00"), DailyBalanceDateTime: time1, UsdPrice: newDecimalFromString("0.50"), UsdValue: newDecimalFromString("2.0000")},
			time2: {BalanceType: "DELEGATION_REWARDS", DailyBalanceStr: newDecimalFromString("2.00"), DailyBalanceDateTime: time2, UsdPrice: newDecimalFromString("2.50"), UsdValue: newDecimalFromString("5.00")},
		}

		actual := balance.GetDailyBalances("Acc1", "AXL")
		assert.Equal(t, expected, actual)
		assert.EqualValues(t, expected[time1].UsdValue, actual[time1].UsdValue)
	})

	t.Run("Test_GetSortedDailyBalances", func(t *testing.T) {
		date1, _ := time.Parse("2006-01-02", "2023-05-31")
		date2, _ := time.Parse("2006-01-02", "2023-06-01")

		testCases := []struct {
			index       int
			balanceDate time.Time
		}{
			{0, date1},
			{1, date2},
		}

		actual := balance.GetSortedDailyBalances("Acc1", "AXL")

		for _, test := range testCases {
			assert.Equal(t, test.balanceDate, actual[test.index].DailyBalanceDateTime)
		}
	})

	t.Run("Test_Aggregate_Balances_Per_Day", func(t *testing.T) {
		date1, _ := time.Parse("2006-01-02", "2023-05-31")
		date2, _ := time.Parse("2006-01-02", "2023-06-01")

		testCases := []struct {
			index          int
			balanceDate    time.Time
			totalqtyAssets decimal.Decimal
			totalUsd       decimal.Decimal
		}{
			{0, date1, newDecimalFromString("4.00"), newDecimalFromString("2.00")},
			{1, date2, newDecimalFromString("2.00"), newDecimalFromString(".00")},
		}

		actual := balance.GetSortedDailyBalances("Acc1", "AXL")

		for _, test := range testCases {
			assert.Equal(t, test.balanceDate, actual[test.index].DailyBalanceDateTime)
			assert.Equal(t, test.totalqtyAssets, actual[test.index].DailyBalanceStr)
		}
	})

	t.Run("Test_IsEmpty", func(t *testing.T) {
		expected := false

		actual := balance.IsEmpty()

		assert.Equal(t, expected, actual)
	})

	t.Run("Test_EmptyDailyBalances", func(t *testing.T) {
		var expected map[time.Time]ubalances.DailyBalance
		expected = nil

		actual := balance.GetDailyBalances("Acc1", "OSMO")
		assert.Equal(t, expected, actual)

		expectedSorted := []ubalances.DailyBalance{}

		actualSorted := balance.GetSortedDailyBalances("Acc1", "OSMO")
		assert.EqualValues(t, expectedSorted, actualSorted)
	})
}

func newDecimalFromString(value string) decimal.Decimal {
	str, _ := decimal.NewFromString(value)
	return str
}
