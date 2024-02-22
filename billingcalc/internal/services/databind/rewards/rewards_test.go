//go:build !selectTest || unitTest

package rewards_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
)

func GetRewards() *rewards.Rewards {
	table := [][]string{
		{"Org1", "Acc1", "Col2", "Col3", "Col4", "Delegation Reward", "ETH", "Col7", "Col8", "Col9", "Col10", "1.00", "1.50", "0.00", "0.00", "Col15", "Col16", "01/01/2023", "AccId1-1"},
		{"Org1", "Acc1", "Col2", "Col3", "Col4", "Delegation Reward", "BTC", "Col7", "Col8", "Col9", "Col10", "0.50", "1.50", "0.00", "0.00", "Col15", "Col16", "01/01/2023", "AccId1-1"},
		{"Org1", "Acc1", "Col2", "Col3", "Col4", "Delegation Reward", "BTC", "Col7", "Col8", "Col9", "Col10", "0.30", "1.80", "0.20", "1.00", "Col15", "Col16", "01/01/2023", "AccId1-1"},
		{"Org1", "Acc2", "Col2", "Col3", "Col4", "Delegation Reward", "BTC", "Col7", "Col8", "Col9", "Col10", "Col11", "5.00", "Col13", "2.00", "Col15", "Col16", "01/01/2023", "AccId1-2"},
		{"Org2", "Acc21", "Col2", "Col3", "Col4", "Delegation Reward", "BTC", "Col7", "Col8", "Col9", "Col10", "Col11", "1.99", "Col13", "0.00", "Col15", "Col16", "01/01/2023", "AccId2-21"},
	}

	return rewards.NewRewards(table)
}

func TestLookupFunctions(t *testing.T) {
	testRewards := GetRewards()

	t.Run("Test GetOrganizations", func(t *testing.T) {
		expected := []string{
			"Org1",
			"Org2",
		}

		actual := testRewards.GetOrganizationNames()

		assert.Equal(t, len(expected), len(actual))
		for _, v := range actual {
			assert.Contains(t, expected, v.Name)
		}
	})

	t.Run("Test GetAccounts", func(t *testing.T) {
		expected := []string{
			"Acc1",
			"Acc2",
		}

		actual := testRewards.GetAccountNames("Org1")

		assert.Equal(t, len(expected), len(actual))
		for _, v := range actual {
			assert.Contains(t, expected, v.Name)
		}
	})

	t.Run("Test_GetAssets", func(t *testing.T) {
		expected := []string{
			"BTC",
			"ETH",
		}

		account := testRewards.GetAccountById(databind.AccountID("AccId1-1"))

		actual := account.GetAssets()

		assert.Equal(t, len(expected), len(actual))
		for _, v := range actual {
			assert.Contains(t, expected, v.Name)
		}
	})

	t.Run("Test_GetClaimedValues", func(t *testing.T) {
		date := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
		opeType := "Delegation Reward"

		expected := []rewards.ClaimedReward{
			{
				AnchorageAssetQty:  newDecimalFromString("0.50"),
				AnchorageUsdValue:  newDecimalFromString("1.50"),
				BusinessDay:        date,
				OperationType:      opeType,
				ThirdPartyAssetQty: newDecimalFromString("0.00"),
				ThirdPartyUsdValue: newDecimalFromString("0.00"),
			},
			{
				AnchorageAssetQty:  newDecimalFromString("0.30"),
				AnchorageUsdValue:  newDecimalFromString("1.80"),
				BusinessDay:        date,
				OperationType:      opeType,
				ThirdPartyAssetQty: newDecimalFromString("0.20"),
				ThirdPartyUsdValue: newDecimalFromString("1.00"),
			},
		}

		account := testRewards.GetAccountById(databind.AccountID("AccId1-1"))
		asset := account.GetAssets()["BTC"]

		actual := asset.GetClaimedRewards()
		assert.Equal(t, expected, actual)
	})

	t.Run("Test_IsEmpty", func(t *testing.T) {
		expected := false

		actual := testRewards.IsEmpty()

		assert.Equal(t, expected, actual)
	})
}

func newDecimalFromString(value string) decimal.Decimal {
	str, _ := decimal.NewFromString(value)
	return str
}
