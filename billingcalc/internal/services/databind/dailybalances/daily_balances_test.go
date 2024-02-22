//go:build !selectTest || unitTest

package dailybalances_test

import (
	"encoding/csv"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/dailybalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/static"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
)

func readCsvFile() ([][]string, error) {
	file, err := static.Files.Open("gsheet/dailybalances.csv")
	if err != nil {
		return nil, errors.New("Failed to open file: " + err.Error())
	}
	dbTable, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, errors.New("Error while reading Daily Balances CSV")
	}

	return dbTable[1:], nil // skip header
}

func TestNewDailyBalance(t *testing.T) {
	table, err := readCsvFile()
	if err != nil {
		t.Fatalf("Failed to read the CSV file: %v", err)
	}

	db := dailybalances.NewDailyBalance(table)
	var msaId, acc, asset string
	var expected decimal.Decimal

	// Get from line 22, column H
	msaId = "22222" // MSAID Following mfr_test_calc.csv
	acc = "epictank5"
	asset = "OSMO"
	expected = converter.FromStrToDecimal("5587853.072")
	if v, err := db.GetAssetBalance(msaId, acc, asset); err != nil || !v.Equal(expected) {
		t.Logf("FAILED GetAssetBalance org:%s acc: %s asset:%s | err: %v", msaId, acc, asset, err)
		t.Fail()
	}

	// Get from line 157, column J
	msaId = "1000347"
	acc = "slowpanda1"
	asset = "BTC"
	expected = converter.FromStrToDecimal("166.575272")
	if v, err := db.GetUsdBalance(msaId, acc, asset); err != nil || !v.Equal(expected) {
		t.Logf("FAILED GetUsdBalance org:%s acc: %s asset:%s | err: %v", msaId, acc, asset, err)
		t.Fail()
	}

	// Get from line 40, column L
	msaId = "1000348"
	acc = "enchantedtree4"
	asset = "HASH"
	expected = converter.FromStrToDecimal("290.65")
	if v, err := db.GetUnclaimedRewardsBalanceUsd(msaId, acc, asset); err != nil || !v.Equal(expected) {
		t.Logf("FAILED GetUnclaimedRewardsBalanceUsd org:%s acc: %s asset:%s | err: %v", msaId, acc, asset, err)
		t.Fail()
	}

	// Get from line 797, column M
	msaId = "1000349"
	acc = "shinywhale2"
	asset = "BTC"
	expected = converter.FromStrToDecimal("361,264,977.78")
	if v, err := db.GetTotalAucUsd(msaId, acc, asset); err != nil || !v.Equal(expected) {
		t.Logf("FAILED GetTotalAucUsd org:%s acc: %s asset:%s | err: %v", msaId, acc, asset, err)
		t.Fail()
	}
}

func TestGetAccountBalances(t *testing.T) {
	table, err := readCsvFile()
	if err != nil {
		t.Fatalf("Failed to read the CSV file: %v", err)
	}

	db := dailybalances.NewDailyBalance(table)

	t.Run("Should return all balances for an specific account", func(t *testing.T) {
		msaId := "11111" // MSAID Following mfr_test_calc.csv
		accountName := "strongknight3"
		actual, err := db.GetAccountBalances(msaId, accountName)

		assert.Nil(t, err)

		// BTC
		assetBalance := converter.FromStrToDecimal("96.2325969")
		usdBalance := converter.FromStrToDecimal("2671661.834")
		unclaimedRewardsBalanceUsd := decimal.Zero
		totalAucUsd := converter.FromStrToDecimal("2671661.83")
		assert.Equal(t, true, actual["BTC"].AssetBalance.Equal(assetBalance))
		assert.Equal(t, true, actual["BTC"].UsdBalance.Equal(usdBalance))
		assert.Equal(t, true, actual["BTC"].UnclaimedRewardsBalanceUsd.Equal(unclaimedRewardsBalanceUsd))
		assert.Equal(t, true, actual["BTC"].TotalAucUsd.Equal(totalAucUsd))

		// ETH
		assetBalance = converter.FromStrToDecimal("6268.935414")
		usdBalance = converter.FromStrToDecimal("11400000")
		unclaimedRewardsBalanceUsd = decimal.Zero
		totalAucUsd = converter.FromStrToDecimal("11401256.42")
		assert.Equal(t, true, actual["ETH"].AssetBalance.Equal(assetBalance))
		assert.Equal(t, true, actual["ETH"].UsdBalance.Equal(usdBalance))
		assert.Equal(t, true, actual["ETH"].UnclaimedRewardsBalanceUsd.Equal(unclaimedRewardsBalanceUsd))
		assert.Equal(t, true, actual["ETH"].TotalAucUsd.Equal(totalAucUsd))
	})

	t.Run("Should return an error message if the account doesnt exist", func(t *testing.T) {
		msaId := "11111" // MSAID Following mfr_test_calc.csv
		accountName := "fake account"
		actual, err := db.GetAccountBalances(msaId, accountName)

		errorMessage := "Account not found - msaID:11111 acc:fake account"

		assert.Nil(t, actual)
		assert.Equal(t, errorMessage, err.Error())
	})
}

func TestGetAverageUsdBalanceByOrg(t *testing.T) {
	file, err := static.Files.Open("gsheet/dailybalances.csv")
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	dbTable, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatal("Error while reading Daily Balances CSV.")
	}

	db := dailybalances.NewDailyBalance(dbTable[1:]) // skip header
	daysInInvoiceDate := 30

	testCases := []struct {
		org           string
		expectedTotal decimal.Decimal
	}{
		{"11111", converter.FromStrToDecimal("14071661.834")},
		{"22222", converter.FromStrToDecimal("82300035.985947577")},
		// MSAIDs Following mfr_test_calc.csv
	}

	for _, tc := range testCases {
		expectedAverage := tc.expectedTotal.DivRound(decimal.NewFromInt(int64(daysInInvoiceDate)), 16)

		averageUsdBalance, err := db.GetAverageUsdBalanceByOrg(tc.org, daysInInvoiceDate)
		if err != nil {
			t.Errorf("Error for org %s: %v", tc.org, err)
			continue
		}
		if !averageUsdBalance.Equal(expectedAverage) {
			t.Errorf("FAILED for org %s: expected %v, got %v", tc.org, expectedAverage, averageUsdBalance)
		}
	}
}
