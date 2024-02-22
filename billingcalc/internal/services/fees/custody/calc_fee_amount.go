package custody

import (
	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
)

// CalcEffectiveFeeAmount calculates the effective fee amount based on fee tiers,
// minimum fees, and the total average AUC amount of the organization.
//
// It returns the effective fee amount calculated according to the provided data
func CalcEffectiveFeeAmount(tiers []mfr.TierData, minimumFee mfr.MinimumFee, totalOrgAvgAuc decimal.Decimal) decimal.Decimal {
	if minimumFee.IsAucBased() && totalOrgAvgAuc.LessThan(tiers[0].Floor) {
		return minimumFee.MinimumCharge
	}

	var effectiveAmount, monthlyRate decimal.Decimal
	remainingAmount := totalOrgAvgAuc

	if len(tiers) == 1 {
		monthlyRate = tiers[0].Rate.DivRound(decimal.NewFromFloat(12), 16)
		effectiveAmount = totalOrgAvgAuc.Mul(monthlyRate)
	} else {
		// graduated tiers
		for _, tier := range tiers {
			monthlyRate = tier.Rate.DivRound(decimal.NewFromFloat(12), 16)
			if remainingAmount.GreaterThan(tier.Floor) {
				effectiveAmount = effectiveAmount.Add(tier.Floor.Mul(monthlyRate))
				remainingAmount = remainingAmount.Sub(tier.Floor)
			} else {
				effectiveAmount = effectiveAmount.Add(remainingAmount.Mul(monthlyRate))
				break
			}
		}
	}

	// Check if the effective amount is less than the minimum fee 'Greater of'
	if minimumFee.IsGreaterOf() && effectiveAmount.LessThan(minimumFee.MinimumCharge) {
		return minimumFee.MinimumCharge
	}

	return effectiveAmount
}
