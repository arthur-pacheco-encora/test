//go:build !selectTest || unitTest

package custody_test

import (
	"testing"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees/custody"
)

func TestCalcEffectiveFeeAmount(t *testing.T) {
	tiers := []mfr.TierData{
		{Floor: decimal.NewFromFloat(10_000_000.0), Rate: decimal.NewFromFloat(0.35)},
		{Floor: decimal.NewFromFloat(100_000_000.0), Rate: decimal.NewFromFloat(0.25)},
		{Floor: decimal.NewFromFloat(1_000_000_000.0), Rate: decimal.NewFromFloat(0.15)},
	}

	minimumFee := mfr.MinimumFee{
		MinimumCharge:  decimal.NewFromFloat(1_000_000.0),
		MinimumFeeType: "AUCBased",
	}

	totalOrgAvgAuc := decimal.NewFromFloat(12_000_000.0)

	t.Run("FeeAmountGreaterThanMinimumFee", func(t *testing.T) {
		effectiveFeeAmount := custody.CalcEffectiveFeeAmount(tiers, minimumFee, totalOrgAvgAuc)
		expected := decimal.NewFromFloat(333_333.33)

		assert.True(t, expected.Equal(effectiveFeeAmount.Round(2)), "Fee Amount %v different than expected %v", effectiveFeeAmount.Round(2), expected)
	})

	t.Run("OrgAmountLesserThanAucBased", func(t *testing.T) {
		totalOrgAvgAuc = decimal.NewFromFloat(2_000_000.0)

		effectiveFeeAmount := custody.CalcEffectiveFeeAmount(tiers, minimumFee, totalOrgAvgAuc)

		expected := minimumFee.MinimumCharge

		assert.True(t, expected.Equal(effectiveFeeAmount.Round(2)), "Fee Amount %v different than expected %v", effectiveFeeAmount.Round(2), expected)
	})

	t.Run("FeeAmountLesserThanGreaterOf", func(t *testing.T) {
		totalOrgAvgAuc = decimal.NewFromFloat(2_000_000.0)

		minimumFee := mfr.MinimumFee{
			MinimumCharge:  decimal.NewFromFloat(50_000_000.0),
			MinimumFeeType: "GreaterOf",
		}

		effectiveFeeAmount := custody.CalcEffectiveFeeAmount(tiers, minimumFee, totalOrgAvgAuc)

		expected := minimumFee.MinimumCharge

		assert.True(t, expected.Equal(effectiveFeeAmount.Round(2)), "Fee Amount %v different than expected %v", effectiveFeeAmount.Round(2), expected)
	})

	t.Run("NotGraduatedTier", func(t *testing.T) {
		tiers := []mfr.TierData{
			{Floor: decimal.NewFromFloat(10_000_000.0), Rate: decimal.NewFromFloat(0.25)},
		}

		totalOrgAvgAuc = decimal.NewFromFloat(12_000_000.0)

		effectiveFeeAmount := custody.CalcEffectiveFeeAmount(tiers, minimumFee, totalOrgAvgAuc)

		expected := decimal.NewFromFloat(250_000.0)

		assert.True(t, expected.Equal(effectiveFeeAmount.Round(2)), "Fee Amount %v different than expected %v", effectiveFeeAmount.Round(2), expected)
	})
}
