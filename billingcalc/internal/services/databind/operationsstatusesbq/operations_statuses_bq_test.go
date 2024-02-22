//go:build !selectTest || unitTest

package operationsstatusesbq_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatusesbq"
)

type MockResultIterator struct {
	Data  []operationsstatusesbq.OperationsStatuses
	Index int
}

// Mock for the Big Querry Iterator
func (m *MockResultIterator) Next(dst interface{}) error {
	if m.Index >= len(m.Data) {
		return iterator.Done
	}

	result, ok := dst.(*operationsstatusesbq.OperationsStatuses)
	if !ok {
		return errors.New("type assertion to *operationsstatusesbq.OperationsStatuses failed")
	}

	*result = m.Data[m.Index]
	m.Index++
	return nil
}

func TestGetDataQuery(t *testing.T) {
	start := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 30, 0, 0, 0, 0, time.UTC)

	got := operationsstatusesbq.GetDataQuery(start, end)

	assert.Contains(t, got, "client_operations.historic_delegation_statuses_confidential AS hds")
	assert.Contains(t, got, "client_operations.accounts_confidential AS accounts")
	assert.Contains(t, got, "client_operations.organizations_confidential as orgs")

	assert.Contains(t, got, start.Format(time.RFC3339))
	assert.Contains(t, got, end.Format(time.RFC3339))

	assert.Contains(t, got, "CAST(delegation_statuses.active_delegated_quantity AS STRING) AS delegation_statuses_active_delegated_quantity")
	assert.Contains(t, got, "CAST(DATE(delegation_statuses.last_updated_at) as STRING) AS delegation_statuses_last_updated_date")

	assert.Contains(t, got, "LEFT JOIN vaults_metadata ON delegation_statuses.vault_unique_id = vaults_metadata.vault_unique_id")
}

func TestStructToSlice(t *testing.T) {
	mockData := []operationsstatusesbq.OperationsStatuses{
		{
			DELEGATION_STATUSES_ACCOUNT_NAME:               bigquery.NullString{StringVal: "Account1", Valid: true},
			DELEGATION_STATUSES_ACTIVE_DELEGATED_VALUE_USD: bigquery.NullString{StringVal: "1000.00", Valid: true},
			DELEGATION_STATUSES_ASSET_TYPE:                 bigquery.NullString{StringVal: "AssetType1", Valid: true},
			DELEGATION_STATUSES_DATE_DATE:                  bigquery.NullString{StringVal: "2024-01-29", Valid: true},
			DELEGATION_STATUSES_IS_ANCHORAGE_VALIDATOR:     bigquery.NullString{StringVal: "Yes", Valid: true},
		},
		{
			DELEGATION_STATUSES_ACCOUNT_NAME:               bigquery.NullString{StringVal: "Account2", Valid: true},
			DELEGATION_STATUSES_ACTIVE_DELEGATED_VALUE_USD: bigquery.NullString{StringVal: "31231.00", Valid: true},
			DELEGATION_STATUSES_ASSET_TYPE:                 bigquery.NullString{StringVal: "AssetType1", Valid: true},
			DELEGATION_STATUSES_DATE_DATE:                  bigquery.NullString{StringVal: "2024-01-29", Valid: true},
			DELEGATION_STATUSES_IS_ANCHORAGE_VALIDATOR:     bigquery.NullString{StringVal: "Yes", Valid: true},
		},
	}

	iter := &MockResultIterator{Data: mockData}

	got, err := operationsstatusesbq.StructToSlice(iter)
	if err != nil {
		t.Fatalf("StructToSlice returned an error: %v", err)
	}

	expected := [][]string{
		{"Account1", "", "1000.00", "", "", "AssetType1", "", "2024-01-29", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "0"},
		{"Account2", "", "31231.00", "", "", "AssetType1", "", "2024-01-29", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "0"},
	}

	assert.NoError(t, err, "StructToSlice should not return an error")
	assert.Equal(t, expected, got, "Expected and actual outcomes should match")
}
