//go:build !selectTest || unitTest

package mfr_test

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/mfr"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/static"
)

var mfrAll [][]string

func readCsvFile() error {
	if mfrAll != nil && len(mfrAll) > 0 {
		return nil
	}

	mfrFile, err := static.Files.Open("gsheet/mfr_test_calc.csv")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open MFR sample file: %v", err))
	}

	csvData, err := csv.NewReader(mfrFile).ReadAll()
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading MFR: %v", err))
	}

	mfrAll = csvData[3:]

	return nil
}

func TestNewMasterFeeRates(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}

	mfrBind := mfr.NewMasterFeeRates(mfrAll)

	assert.Falsef(t, mfrBind.IsEmpty(), "Organization map should not be empty")
}

func TestLookupFunctions(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}

	mfrBind := mfr.NewMasterFeeRates(mfrAll)

	t.Run("Test GetOrganizations", func(t *testing.T) {
		orgs := mfrBind.GetOrganizations()
		assert.Equal(t, 2, len(orgs), "Must have 2 organizations read from csv file.")
	})

	t.Run("Test GetAccounts", func(t *testing.T) {
		msaId := mfr.MSAID("22222")
		accounts := mfrBind.GetAccounts(msaId)
		assert.Equal(t, 2, len(accounts), "Must have 2 accounts for MSAID=22222")
	})

	t.Run("Test GetAssetTypes", func(t *testing.T) {
		msaId := mfr.MSAID("22222")
		accountId := databind.AccountID("accountIdFor2222")
		assetTypes := mfrBind.GetAssetTypes(msaId, accountId)
		assert.Equal(t, 2, len(assetTypes), "Must have 2 assettypes for MSAID=22222 | AccountId=accountIdFor2222")

		t.Run("TestMinimumFee", func(t *testing.T) {
			expectedCharge, _ := decimal.NewFromString("5000")

			minimumFee := assetTypes[10].MinimumFee
			assert.Equal(t, "Greaterof", minimumFee.MinimumFeeType, "Minimum Fee Type for AssetId=10 invalid: %s", minimumFee.MinimumFeeType)
			assert.True(t, expectedCharge.Equal(minimumFee.MinimumCharge), "Minimum Charge for AssetId=10 invalid: %v", minimumFee.MinimumCharge)

			minimumFee = assetTypes[9].MinimumFee
			expectedCharge, _ = decimal.NewFromString("5000")

			assert.Equal(t, "AUCbased", minimumFee.MinimumFeeType, "Minimum Fee Type for AssetId=9 invalid: %s", minimumFee.MinimumFeeType)
			assert.True(t, expectedCharge.Equal(minimumFee.MinimumCharge), "Minimum Charge for AssetId=9 invalid: %v", minimumFee.MinimumCharge)
		})
	})

	t.Run("Test GetStakingFees", func(t *testing.T) {
		msaId := mfr.MSAID("22222")
		accountId := databind.AccountID("accountIdFor2222")
		assetId := mfr.AssetID(10)
		stakingFees := mfrBind.GetStakingFees(msaId, accountId, assetId)

		// testing CELO because it is at the beginning of the array
		asset := stakingFees["CELO"]
		expAnchorFee, _ := decimal.NewFromString("10")
		expThirdFee, _ := decimal.NewFromString("10")

		assert.Truef(t, expAnchorFee.Equal(asset.AnchorageFee), "Expected Anchorage Fee for CELO = %s, but found %s", expAnchorFee, asset.AnchorageFee)
		assert.Truef(t, expThirdFee.Equal(asset.ThirdPartyFee), "Expected Third Party Fee for CELO = %s, but found %s", expThirdFee, asset.ThirdPartyFee)

		// testing SUI because it is at the end of the array
		asset = stakingFees["SUI"]
		expAnchorFee, _ = decimal.NewFromString("6.0")
		expThirdFee, _ = decimal.NewFromString("3.0")

		assert.Truef(t, expAnchorFee.Equal(asset.AnchorageFee), "Expected Anchorage Fee for SUI = %s, but found %s", expAnchorFee, asset.AnchorageFee)
		assert.Truef(t, expThirdFee.Equal(asset.ThirdPartyFee), "Expected Third Party Fee for SUI = %s, but found %s", expThirdFee, asset.ThirdPartyFee)
	})

	t.Run("Test GetAssetStakingFees", func(t *testing.T) {
		msaId := mfr.MSAID("22222")
		accountId := databind.AccountID("accountIdFor2222")

		account := mfrBind.GetAccounts(msaId)[accountId]

		stakingFees := account.GetAssetStakingFees("CELO")

		expAnchorFee, _ := decimal.NewFromString("10")
		expThirdFee, _ := decimal.NewFromString("10")

		assert.Truef(t, expAnchorFee.Equal(stakingFees.AnchorageFee), "Expected Anchorage Fee for CELO = %s, but found %s", expAnchorFee, stakingFees.AnchorageFee)
		assert.Truef(t, expThirdFee.Equal(stakingFees.ThirdPartyFee), "Expected Third Party Fee for CELO = %s, but found %s", expThirdFee, stakingFees.ThirdPartyFee)
	})

	t.Run("Test GetAssetStakingFees NotFound", func(t *testing.T) {
		msaId := mfr.MSAID("22222")
		accountId := databind.AccountID("accountIdFor2222")

		account := mfrBind.GetAccounts(msaId)[accountId]

		stakingFees := account.GetAssetStakingFees("MY_COIN")

		assert.True(t, stakingFees.AssetName == "", "Expected empty AssetName, but got %s", stakingFees.AssetName)
	})
}

func TestFindAllTiersGraduatedClient(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}

	mfrBind := mfr.NewMasterFeeRates(mfrAll)
	msaId := mfr.MSAID("22222")
	accountId := databind.AccountID("accountIdFor2222")
	assetId := mfr.AssetID(10)

	usdBalance, err := decimal.NewFromString("350000000")
	if err != nil {
		t.Fatal(err)
	}

	asset := mfrBind.GetAssetTypes(msaId, accountId)[assetId]

	floor1, _ := decimal.NewFromString("0")
	rate1, _ := decimal.NewFromString("0.17")
	floor2, _ := decimal.NewFromString("100000000")
	rate2, _ := decimal.NewFromString("0.15")
	floor3, _ := decimal.NewFromString("250000000")
	rate3, _ := decimal.NewFromString("0.13")

	expectedTiers := []mfr.TierData{
		{Floor: floor1, Rate: rate1},
		{Floor: floor2, Rate: rate2},
		{Floor: floor3, Rate: rate3},
	}
	tiers, err := asset.FindAllTiers(usdBalance)
	assert.NoError(t, err, "FindAllTiers should not return an error for valid data")
	assert.NotNil(t, tiers, "Returned tiers should not be nil")
	assert.NotEmpty(t, tiers, "Returned tiers slice should not be empty for valid data")
	assert.Equal(t, expectedTiers, tiers, "Returned tiers should match expected tiers")
}

func TestFindAllTiersNotGraduated(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}
	mfrBind := mfr.NewMasterFeeRates(mfrAll)
	msaId := mfr.MSAID("11111")
	accountId := databind.AccountID("2d0d35f608815f0a406d9b44d4b3af141b6c2258937028dc8a0b003616afdf22")
	assetId := mfr.AssetID(10)

	usdBalance, err := decimal.NewFromString("60000000")
	if err != nil {
		t.Fatal(err)
	}

	asset := mfrBind.GetAssetTypes(msaId, accountId)[assetId]
	floor1, _ := decimal.NewFromString("50000000")
	rate1, _ := decimal.NewFromString("0.25")
	expectedTiers := []mfr.TierData{
		{Floor: floor1, Rate: rate1},
	}

	tiers, err := asset.FindAllTiers(usdBalance)
	assert.Equal(t, expectedTiers, tiers, "Returned tiers should match expected tiers")
	assert.NoError(t, err, "FindAllTiers should not return an error for valid data")
	assert.NotNil(t, tiers, "Returned tiers should not be nil")
	assert.NotEmpty(t, tiers, "Returned tiers slice should not be empty for valid data")
}

func TestFindAllTiersNegativeValue(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}

	mfrBind := mfr.NewMasterFeeRates(mfrAll)
	msaId := mfr.MSAID("22222")
	accountId := databind.AccountID("accountIdFor2222")
	assetId := mfr.AssetID(10)

	usdBalance, err := decimal.NewFromString("-350000000")
	if err != nil {
		t.Fatal(err)
	}

	asset := mfrBind.GetAssetTypes(msaId, accountId)[assetId]

	tiers, err := asset.FindAllTiers(usdBalance)

	assert.Error(t, err, "FindAllTiers should return an error for negative balance")
	assert.Empty(t, tiers, "FindAllTiers should return empty tiers for negative balance")
}

func TestSkipsTerminatedAccounts(t *testing.T) {
	err := readCsvFile()
	if err != nil {
		t.Fatal(err)
	}

	mfrBind := mfr.NewMasterFeeRates(mfrAll)
	msaId := mfr.MSAID("22222")

	accounts := mfrBind.GetAccounts(msaId)

	for accountId := range accounts {
		assert.NotEqual(t, "TERMINATED", string(accountId), "Account with ID 'Terminated' should not be present")
		assert.NotEqual(t, "Terminated", string(accountId), "Account with ID 'Terminated' should not be present")
	}
}

func TestNewMasterFeeRatesWithEmptyData(t *testing.T) {
	emptyData := [][]string{}

	mfrBind := mfr.NewMasterFeeRates(emptyData)
	assert.True(t, mfrBind.IsEmpty(), "MasterFeeRates should be empty when initialized with empty data")
}

func TestProcessMfrWithMissingParameters(t *testing.T) {
	_, err := mfr.ProcessMfr(context.Background(), nil, "", "", "")
	assert.Error(t, err, "ProcessMfr should return an error when missing parameters")
	assert.Contains(t, err.Error(), "Missing parameters", "Error message should indicate missing parameters")
}

func TestProcessMfrWithInvalidGoogleSheet(t *testing.T) {
	_, err := mfr.ProcessMfr(context.Background(), nil, "sheetId", "mfrTab", "token")
	assert.Error(t, err, "ProcessMfr should return an error for invalid Google Sheet request")
	assert.Contains(t, err.Error(), "Failed to fetch data from mfr endpoint", "Error message should indicate failure to fetch data")
}
