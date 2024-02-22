//go:build !selectTest || unitTest

package rewardsbq_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewardsbq"
)

type MockResultIterator struct {
	Data  []rewardsbq.Rewards
	Index int
}

// Mock for the Big Querry Iterator
func (m *MockResultIterator) Next(dst interface{}) error {
	if m.Index >= len(m.Data) {
		return iterator.Done
	}
	result, ok := dst.(*rewardsbq.Rewards)
	if !ok {
		return errors.New("type assertion to *rewardsbq.Rewards failed")
	}
	*result = m.Data[m.Index]
	m.Index++
	return nil
}

func TestGetDataQuery(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	got := rewardsbq.GetDataQuery(start, end)

	assert.Contains(t, got, "client_operations.operations_confidential ops")
	assert.Contains(t, got, "client_operations.accounts_confidential a")
	assert.Contains(t, got, "client_operations.vaults_metadata vm")
	assert.Contains(t, got, "kyc_operations.duplicate_affiliates dp")
	assert.Contains(t, got, "client_operations.organizations_confidential orgs")
	assert.Contains(t, got, "client_operations.organizations_confidential dest_orgs")

	assert.Contains(t, got, start.Format(time.RFC3339))
	assert.Contains(t, got, end.Format(time.RFC3339))

	assert.Contains(t, got, "AND (operations.operation_state ) = 'COMPLETE'")
	assert.Contains(t, got, "AND (operations.type ) IN ('Delegation Reward', 'Staking Reward')")
}

func TestStructToSlice(t *testing.T) {
	mockData := []rewardsbq.Rewards{
		{
			OPERATIONS_ORGANIZATION_NAME:                   bigquery.NullString{StringVal: "Org1", Valid: true},
			OPERATIONS_ACCOUNT_NAME:                        bigquery.NullString{StringVal: "Account1", Valid: true},
			OPERATIONS_TYPE:                                bigquery.NullString{StringVal: "Delegation Reward", Valid: true},
			OPERATIONS_TOTAL_ASSET_QUANTITY:                bigquery.NullString{StringVal: "124.07", Valid: true},
			OPERATIONS_ASSET_TYPE:                          bigquery.NullString{StringVal: "ROSE", Valid: true},
			OPERATIONS_TOTAL_ANCHORAGE_USD_REWARD_PART:     bigquery.NullString{StringVal: "12.00", Valid: true},
			OPERATIONS_TOTAL_NON_ANCHORAGE_REWARD_PART:     bigquery.NullString{StringVal: "1.27", Valid: true},
			OPERATIONS_TOTAL_NON_ANCHORAGE_USD_REWARD_PART: bigquery.NullString{StringVal: "1.27", Valid: true},
			BUSINESS_DAY:                                   bigquery.NullString{StringVal: "2024-01-01", Valid: true},
		},
	}

	iter := &MockResultIterator{Data: mockData}

	got, err := rewardsbq.StructToSlice(iter)
	if err != nil {
		t.Fatalf("StructToSlice returned an error: %v", err)
	}

	expected := [][]string{
		{
			"Org1", "Account1", "", "", "", "Delegation Reward", "ROSE", "", "", "", "", "124.07", "12.00", "1.27", "1.27", "", "", "2024-01-01", "",
		},
	}

	assert.NoError(t, err, "StructToSlice should not return an error")
	assert.Equal(t, expected, got, "Expected and actual outcomes should match")
}
