package custody

import (
	"errors"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/assettypes"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
)

func AvgAucByAsset(balances map[string]dailybalances.Balance, assetType int, numberDays int64) (map[string]decimal.Decimal, error) {
	result := make(map[string]decimal.Decimal)
	assetTypeList, err := assettypes.NewAssetTypeList()
	if err != nil {
		return result, errors.New(err.Error())
	}

	for key, balance := range balances {
		assetType, err := assetTypeList.GetTypeByAssetName(key, assetType)
		if err != nil {
			continue
		}

		if assetType != nil {
			result[key] = balance.TotalAucUsd.DivRound(decimal.NewFromInt(numberDays), 16)
		}
	}
	return result, nil
}
