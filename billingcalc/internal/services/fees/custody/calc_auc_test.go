//go:build !selectTest || unitTest

package custody_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	fees "github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees/custody"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
)

func TestAvgAucByAsset(t *testing.T) {
	balances := make(map[string]dailybalances.Balance)

	balances["BTC"] = dailybalances.Balance{
		AssetBalance:               converter.FromStrToDecimal("1_000_000"),
		UsdBalance:                 converter.FromStrToDecimal("1_000_000"),
		UnclaimedRewardsBalanceUsd: converter.FromStrToDecimal("1_000_000"),
		TotalAucUsd:                converter.FromStrToDecimal("1_000_000"),
	}

	balances["WPUNKS"] = dailybalances.Balance{
		AssetBalance:               converter.FromStrToDecimal("1_000_000"),
		UsdBalance:                 converter.FromStrToDecimal("1_000_000"),
		UnclaimedRewardsBalanceUsd: converter.FromStrToDecimal("1_000_000"),
		TotalAucUsd:                converter.FromStrToDecimal("1_000_000"),
	}

	t.Run("Should return the average if the asset is found and exclusive is false", func(t *testing.T) {
		assetType := 10
		numberOFDays := int64(31)
		result, err := fees.AvgAucByAsset(balances, assetType, numberOFDays)
		expected := map[string]decimal.Decimal{
			"WPUNKS": converter.FromStrToDecimal("32258.0645161290322581"),
		}
		assert.Nil(t, err)
		assert.NotContains(t, result, "BTC", "Expected to not find the value")
		assert.Equal(t, true, result["WPUNKS"].Equal(expected["WPUNKS"]))
	})

	t.Run("Should not return the average if the asset is found and exclusive is true", func(t *testing.T) {
		assetType := 11
		numberOFDays := int64(31)
		result, err := fees.AvgAucByAsset(balances, assetType, numberOFDays)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(result))
	})

	t.Run("Should return the average if the asset is not found and exclusive is false", func(t *testing.T) {
		assetType := 10
		numberOFDays := int64(31)
		result, err := fees.AvgAucByAsset(balances, assetType, numberOFDays)
		expected := map[string]decimal.Decimal{
			"WPUNKS": converter.FromStrToDecimal("32258.0645161290322581"),
		}
		assert.Nil(t, err)
		assert.NotContains(t, result, "BTC", "Expected to not find the value")
		assert.Equal(t, true, result["WPUNKS"].Equal(expected["WPUNKS"]))
	})

	t.Run("Should not return the average if the asset is not found and exclusive is true", func(t *testing.T) {
		assetType := 11
		numberOFDays := int64(31)
		result, err := fees.AvgAucByAsset(balances, assetType, numberOFDays)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(result))
	})
}
