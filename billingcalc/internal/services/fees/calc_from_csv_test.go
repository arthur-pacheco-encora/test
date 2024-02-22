//go:build !selectTest || unitTest

package fees_test

// import (
// 	"context"
// 	"encoding/json"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// 	"google.golang.org/protobuf/types/known/timestamppb"

// 	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
// 	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
// 	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees"
// )

// var fixedTime = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

// func TestParseMfr(t *testing.T) {
// 	ctx := context.Background()

// 	req := &billingcalcproto.GetCalcFromCsvRequest{
// 		Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 		Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 		Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 		FirstexternalId: 123,
// 		InvoiceDate:     timestamppb.New(fixedTime),
// 	}

// 	resp, err := Calcul(ctx, req)
// 	assert.Nilf(t, err, "Function should not return an error, but got %v", err)
// 	assert.NotNil(t, resp)
// }

// func TestGetRewards(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("Valid Unclaimed Balances data", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}

// 		resp, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.Nilf(t, err, "Function should not return an error, but got %v", err)
// 		assert.NotNil(t, resp)

// 		var resultData []struct {
// 			Rewards *rewards.Rewards `json:"rewards"`
// 		}
// 		err = json.Unmarshal(resp.Data, &resultData)
// 		assert.Nilf(t, err, "Process should not return an error, but got %v", err)
// 		assert.NotNil(t, resultData)
// 	})
// }

// func TestGetCalcFromCsv_UnclaimedBalances(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("Valid Unclaimed Balances data", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}

// 		resp, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.Nilf(t, err, "Function should not return an error, but got %v", err)
// 		assert.NotNil(t, resp)

// 		var resultData []struct {
// 			UnclaimedBalances *unclaimedbalances.UnclaimedBalances `json:"unclaimed_balances"`
// 		}
// 		err = json.Unmarshal(resp.Data, &resultData)
// 		assert.Nilf(t, err, "Process should not return an error, but got %v", err)
// 		assert.NotNil(t, resultData)
// 	})
// }

// func TestGetCalcFromCsv(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("Success case", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}

// 		resp, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.Nil(t, err)
// 		assert.NotNil(t, resp)
// 	})

// 	t.Run("Invalid MFR data", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvInvalidMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}
// 		_, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.NotNil(t, err)
// 	})

// 	t.Run("Invalid Rewards data", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvInvalidRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}

// 		_, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.NotNil(t, err)
// 	})

// 	t.Run("Invalid Unclaimed Balances data", func(t *testing.T) {
// 		req := &billingcalcproto.GetCalcFromCsvRequest{
// 			Mfr:             []byte(getcalcfromcsv.CsvMfrString),
// 			Rewards:         []byte(getcalcfromcsv.CsvRewardsString),
// 			Unclaimed:       []byte(getcalcfromcsv.CsvInvalidUnclaimedBalancesString),
// 			FirstexternalId: 123,
// 			InvoiceDate:     timestamppb.New(fixedTime),
// 		}

// 		_, err := getcalcfromcsv.GetCalcFromCsv(ctx, req)
// 		assert.NotNil(t, err)
// 	})
// }
