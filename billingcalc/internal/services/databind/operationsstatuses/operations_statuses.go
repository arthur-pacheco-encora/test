package operationsstatuses

import (
	"slices"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

const (
	ColStatusesAccountName          databind.Column = 0
	ColStatusesActiveDelegatedValue databind.Column = 2 // USD
	ColStatusesAssetType            databind.Column = 5
	ColStatusesDate                 databind.Column = 7
	ColCosmosValidatorsRate         databind.Column = 24
)

type OperationsStatuses struct {
	Accounts map[string]Account
}

type Account struct {
	Statuses []Status
}

type Status struct {
	CosmosValidatorsRatesRate    string
	StatusesActiveDelegatedValue decimal.Decimal
	StatusesAssetType            string
	StatusesDate                 time.Time
}

func (c *OperationsStatuses) IsEmpty() bool {
	return len(c.Accounts) == 0
}

func (c *OperationsStatuses) GetStatusByAssetAndDate(account string, asset string, date time.Time) Status {
	var status Status
	acc := c.Accounts[account]

	zeroDuration, _ := time.ParseDuration("0")

	for _, row := range acc.Statuses {
		isSameDate := row.StatusesDate.Sub(date) == zeroDuration

		if row.StatusesAssetType == asset && isSameDate {
			status = row
			break
		}
	}

	return status
}

func NewOperationsStatuses(table [][]string) *OperationsStatuses {
	operationsStatuses := &OperationsStatuses{
		Accounts: make(map[string]Account),
	}

	for _, row := range table {
		accountName := sanitization.SanitizeName(row[ColStatusesAccountName])

		acc, exists := operationsStatuses.Accounts[accountName]
		if !exists {
			acc = Account{
				Statuses: []Status{},
			}
		}

		decimalActiveDelegatedValue, _ := decimal.NewFromString(row[ColStatusesActiveDelegatedValue])

		acc.Statuses = append(acc.Statuses, Status{
			CosmosValidatorsRatesRate:    row[ColCosmosValidatorsRate],
			StatusesActiveDelegatedValue: decimalActiveDelegatedValue,
			StatusesAssetType:            row[ColStatusesAssetType],
			StatusesDate:                 converter.FromStrToTimeStamp(row[ColStatusesDate]),
		})

		operationsStatuses.Accounts[accountName] = acc
	}

	return operationsStatuses
}

func IsAssetFromExternalValidator(operationStatus Status) bool {
	cosmosAssets := []string{"OSMO", "HASH", "ATOM", "AXL", "EVMOS", "SEI", "SUI"}

	isAssetsInCosmosAssetsList := slices.IndexFunc(cosmosAssets, func(cosmoAsset string) bool {
		return cosmoAsset == strings.ToUpper(operationStatus.StatusesAssetType)
	}) != -1

	isFullRate := operationStatus.CosmosValidatorsRatesRate == "100.00%" || operationStatus.CosmosValidatorsRatesRate == "100%" || operationStatus.CosmosValidatorsRatesRate == "1"

	return isAssetsInCosmosAssetsList && isFullRate
}
