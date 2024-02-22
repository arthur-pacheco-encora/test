//go:build !selectTest || unitTest

package operationsstatuses_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
)

var newActiveDelegatedValue, _ = decimal.NewFromString("407678.72")

var tests = []struct {
	name   string
	status operationsstatuses.Status
	want   bool
}{
	{"is asset from external validator", operationsstatuses.Status{
		CosmosValidatorsRatesRate:    "100.00%",
		StatusesActiveDelegatedValue: newActiveDelegatedValue,
		StatusesAssetType:            "HASH",
		StatusesDate:                 time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
	}, true},
	{"is not asset from external validator 50%", operationsstatuses.Status{
		CosmosValidatorsRatesRate:    "50.00%",
		StatusesActiveDelegatedValue: newActiveDelegatedValue,
		StatusesAssetType:            "HASH",
		StatusesDate:                 time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
	}, false},
	{"is not asset from external validator if asset is ROSE", operationsstatuses.Status{
		CosmosValidatorsRatesRate:    "100.00%",
		StatusesActiveDelegatedValue: newActiveDelegatedValue,
		StatusesAssetType:            "ROSE",
		StatusesDate:                 time.Date(2007, 1, 1, 0, 0, 0, 0, time.UTC),
	}, false},
}

func TestIsAssetFromExternalValidator(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := operationsstatuses.IsAssetFromExternalValidator(tt.status)
			if got != tt.want {
				t.Errorf("IsAssetFromExternalValidator(%s) got %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		name     string
		input    operationsstatuses.OperationsStatuses
		expected bool
	}{
		{
			name:     "Empty OperationsStatuses",
			input:    operationsstatuses.OperationsStatuses{Accounts: map[string]operationsstatuses.Account{}},
			expected: true,
		},
		{
			name: "Non-empty OperationsStatuses",
			input: operationsstatuses.OperationsStatuses{Accounts: map[string]operationsstatuses.Account{
				"account1": {Statuses: []operationsstatuses.Status{{}}},
			}},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.input.IsEmpty()
			if result != c.expected {
				t.Errorf("TestIsEmpty - %s: expected %v, got %v", c.name, c.expected, result)
			}
		})
	}
}

func TestGetStatusByAssetAndDate(t *testing.T) {
	testDate := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)
	testStatus := operationsstatuses.Status{
		StatusesAssetType: "TestAsset",
		StatusesDate:      testDate,
	}
	operations := operationsstatuses.OperationsStatuses{
		Accounts: map[string]operationsstatuses.Account{
			"testAccount": {
				Statuses: []operationsstatuses.Status{testStatus},
			},
		},
	}

	cases := []struct {
		name          string
		account       string
		asset         string
		date          time.Time
		expectedFound bool
	}{
		{
			name:          "Status Found",
			account:       "testAccount",
			asset:         "TestAsset",
			date:          testDate,
			expectedFound: true,
		},
		{
			name:          "Status Not Found",
			account:       "testAccount",
			asset:         "NonExistentAsset",
			date:          testDate,
			expectedFound: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := operations.GetStatusByAssetAndDate(c.account, c.asset, c.date)
			found := !(result == operationsstatuses.Status{})
			if found != c.expectedFound {
				t.Errorf("TestGetStatusByAssetAndDate - %s: expected %v, got %v", c.name, c.expectedFound, found)
			}
		})
	}
}

func TestNewOperationsStatuses(t *testing.T) {
	table := [][]string{
		{"Test Alpha Account", "", "1459200.01", "", "", "HASH", "", "2023-06-29", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "100.00%", "", ""},
		{"Test Alpha Account2", "", "412229550.01", "", "", "FLOW", "", "2023-06-29", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "100.00%", "", ""},
		{"Test Alpha Account3", "", "3349200.01", "", "", "CELO", "", "2023-06-29", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "100.00%", "", ""},
	}

	result := operationsstatuses.NewOperationsStatuses(table)

	if len(result.Accounts) == 0 {
		t.Error("TestNewOperationsStatuses: expected non-empty accounts")
	}

	for accountName, account := range result.Accounts {
		for _, status := range account.Statuses {
			if status.StatusesAssetType == "" || status.StatusesDate.IsZero() {
				t.Errorf("TestNewOperationsStatuses: missing asset type or date in account %s", accountName)
			}
		}
	}
}
