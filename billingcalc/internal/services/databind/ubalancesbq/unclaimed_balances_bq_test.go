//go:build !selectTest || unitTest

package ubalancesbq_test

import (
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalancesbq"
)

type MockResultIterator struct {
	Data  []ubalancesbq.UnclaimedBalances
	Index int
}

// Mock for the Big Querry Iterator
func (m *MockResultIterator) Next(dst interface{}) error {
	if m.Index >= len(m.Data) {
		return iterator.Done
	}
	result, ok := dst.(*ubalancesbq.UnclaimedBalances)
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
	got := ubalancesbq.GetDataQuery(start, end)

	assert.Contains(t, got, "historic_daily_blockchain_balances")
	assert.Contains(t, got, "vaults_metadata")
	assert.Contains(t, got, "historic_daily_blockchain_balances.business_day")
	assert.Contains(t, got, "organizations.org_name")

	assert.Contains(t, got, start.Format(time.RFC3339))
	assert.Contains(t, got, end.Format(time.RFC3339))

	assert.Contains(t, got, "AND (balance_type = 'DELEGATION_REWARDS')")
}

func TestStructToSlice(t *testing.T) {
	mockData := []ubalancesbq.UnclaimedBalances{
		{
			HISTORIC_DAILY_BLOCKCHAIN_BALANCES_ASSET_TYPE_ID:        bigquery.NullString{StringVal: "AssetType1", Valid: true},
			HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_TYPE:         bigquery.NullString{StringVal: "BalanceType1", Valid: true},
			HISTORIC_DAILY_BLOCKCHAIN_BALANCES_BALANCE_STR:          bigquery.NullString{StringVal: "1000.00", Valid: true},
			HISTORIC_DAILY_BLOCKCHAIN_BALANCES_LAST_UPDATED_AT_TIME: bigquery.NullString{StringVal: "2024-01-01", Valid: true},
			CUSTODY_ACCOUNTS_ACCOUNT_NAME:                           bigquery.NullString{StringVal: "Account1", Valid: true},
			USD_PRICE:                                               bigquery.NullString{StringVal: "50.00", Valid: true},
			USD_VALUE:                                               bigquery.NullString{StringVal: "50000.00", Valid: true},
		},
	}

	iter := &MockResultIterator{Data: mockData}

	got, err := ubalancesbq.StructToSlice(iter)
	if err != nil {
		t.Fatalf("StructToSlice returned an error: %v", err)
	}

	expected := [][]string{
		{
			"",
			"",
			"",
			"AssetType1",
			"",
			"BalanceType1",
			"1000.00",
			"2024-01-01",
			"Account1",
			"50.00",
			"50000.00",
		},
	}

	assert.NoError(t, err, "StructToSlice should not return an error")
	assert.Equal(t, expected, got, "Expected and actual outcomes should match")
}
