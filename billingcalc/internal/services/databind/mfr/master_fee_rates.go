package mfr

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/googlesheetsutils"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

// MFR Columns
const (
	colMinimumFeeType     databind.Column = 3
	colAccountFromRDB     databind.Column = 4
	colEntityID           databind.Column = 7
	colOrgName            databind.Column = 8
	colLegalName          databind.Column = 9
	colBillingTerms       databind.Column = 11
	colAssetType          databind.Column = 19
	colMinimumCharge      databind.Column = 29
	col1stTierFloor       databind.Column = 30
	col1stTierRate        databind.Column = 31
	col2ndTierFloor       databind.Column = 32
	col2ndTierRate        databind.Column = 33
	col3rdTierFloor       databind.Column = 34
	col3rdTierRate        databind.Column = 35
	col4thTierFloor       databind.Column = 36
	col4thTierRate        databind.Column = 37
	col5thTierFloor       databind.Column = 38
	col5thTierRate        databind.Column = 39
	col6thTierFloor       databind.Column = 40
	col6thTierRate        databind.Column = 41
	col7thTierFloor       databind.Column = 42
	col7thTierRate        databind.Column = 43
	col8thTierFloor       databind.Column = 44
	col8thTierRate        databind.Column = 45
	col9thTierFloor       databind.Column = 46
	col9thTierRate        databind.Column = 47
	col10thTierFloor      databind.Column = 48
	col10thTierRate       databind.Column = 49
	colCeloFeeAnchorage   databind.Column = 53
	colCeloFeeThirdParty  databind.Column = 54
	colFlowFeeAnchorage   databind.Column = 56
	colFlowFeeThirdParty  databind.Column = 57
	colOsmoFeeThirdParty  databind.Column = 58
	colRoseFeeAnchorage   databind.Column = 60
	colRoseFeeThirdParty  databind.Column = 61
	colEthFeeAnchorage    databind.Column = 62
	colAxlFeeThirdParty   databind.Column = 64
	colAptFeeThirdParty   databind.Column = 65
	colAtomFeeThirdParty  databind.Column = 67
	colHashFeeThirdParty  databind.Column = 69
	colEvmosFeeThirdParty databind.Column = 71
	colAptFeeAnchorage    databind.Column = 72
	colSolFeeThirdParty   databind.Column = 73
	colSuiFeeAnchorage    databind.Column = 75
	colSuiFeeThirdParty   databind.Column = 76
	colAssetID            databind.Column = 87
	colMSAID              databind.Column = 88
	colGraduatedTier      databind.Column = 89
	colCustomerID         databind.Column = 90 // NetSuite Account ID
	colBillingID          databind.Column = 91
	colRDBAccountID       databind.Column = 92

	mfrHeaderRow = 3
)

type (
	MSAID   string
	AssetID int
)

type MasterFeeRates struct {
	organizations map[MSAID]Organization
}

type Organization struct {
	Id          MSAID
	Name        string
	DisplayName string
	EntityId    string
	accounts    map[databind.AccountID]Account
}

type Account struct {
	Id           databind.AccountID
	Name         string
	DisplayName  string
	CustomerId   string
	BillingTerms string
	assetTypes   map[AssetID]AssetType
}

type AssetType struct {
	Id            AssetID
	Description   string
	stakingFees   map[string]StakingFee
	GraduatedTier string
	TierData      []TierData
	MinimumFee    MinimumFee
}

type TierData struct {
	Floor decimal.Decimal
	Rate  decimal.Decimal
}

type StakingFee struct {
	AssetName     string
	AnchorageFee  decimal.Decimal
	ThirdPartyFee decimal.Decimal
}

type MinimumFee struct {
	MinimumFeeType string
	MinimumCharge  decimal.Decimal
}

func (r *MasterFeeRates) GetOrganizations() map[MSAID]Organization {
	return r.organizations
}

func (r *MasterFeeRates) GetAccounts(msaId MSAID) map[databind.AccountID]Account {
	org := r.organizations[msaId]
	return org.accounts
}

func (r *MasterFeeRates) GetAssetTypes(msaId MSAID, account databind.AccountID) map[AssetID]AssetType {
	org := r.organizations[msaId]
	acc := org.accounts[account]
	return acc.assetTypes
}

func (a *Account) GetAssetStakingFees(assetName string) StakingFee {
	// For staking we don't look for different assetTypes
	for _, atypes := range a.assetTypes {
		for _, s := range atypes.stakingFees {
			if s.AssetName == assetName {
				return s
			}
		}
	}

	return StakingFee{}
}

func (r *MasterFeeRates) GetStakingFees(msaId MSAID, account databind.AccountID, assetID AssetID) map[string]StakingFee {
	orgs := r.organizations[msaId]
	acc := orgs.accounts[account]
	a := acc.assetTypes[assetID]
	return a.stakingFees
}

func (r *MasterFeeRates) IsEmpty() bool {
	return len(r.organizations) == 0
}

func (at *AssetType) FindAllTiers(usdBalance decimal.Decimal) ([]TierData, error) {
	filteredTiers, err := at.getTierAndFloor()
	if err != nil {
		return nil, err
	}

	suitableTiers := make([]TierData, 0)

	for _, tier := range filteredTiers {
		if usdBalance.GreaterThanOrEqual(tier.Floor) {
			suitableTiers = append(suitableTiers, tier)
		}
	}

	if len(suitableTiers) == 0 {
		return nil, errors.New("no suitable tiers found")
	}

	isGraduatedTier, err := at.getGraduatedTier()
	if err != nil {
		return nil, err
	}

	if isGraduatedTier {
		return suitableTiers, nil
	} else {
		return []TierData{suitableTiers[len(suitableTiers)-1]}, nil
	}
}

func (at *AssetType) getGraduatedTier() (bool, error) {
	return strings.ToUpper(at.GraduatedTier) == "GRADUATED", nil
}

func (at *AssetType) getTierAndFloor() ([]TierData, error) {
	var filteredTierData []TierData
	for _, tierData := range at.TierData {
		if !(tierData.Floor.String() == "0" && tierData.Rate.String() == "0") {
			filteredTierData = append(filteredTierData, tierData)
		}
	}

	if len(filteredTierData) == 0 {
		return nil, errors.New("No valid tier data found for asset")
	}

	return filteredTierData, nil
}

func (mf *MinimumFee) IsAucBased() bool {
	return strings.ToUpper(mf.MinimumFeeType) == "AUCBASED"
}

func (mf *MinimumFee) IsGreaterOf() bool {
	return strings.ToUpper(mf.MinimumFeeType) == "GREATEROF"
}

func NewMasterFeeRates(table [][]string) *MasterFeeRates {
	mfr := &MasterFeeRates{
		organizations: make(map[MSAID]Organization),
	}

	for _, row := range table {
		if strings.ToUpper(row[colRDBAccountID]) == "TERMINATED" {
			continue
		}
		msaId := MSAID(sanitization.SanitizeName(row[colMSAID]))
		orgName := sanitization.SanitizeName(row[colOrgName])
		org, exists := mfr.organizations[msaId]
		if !exists {
			org = Organization{
				Id:          msaId,
				Name:        orgName,
				DisplayName: row[colOrgName],
				EntityId:    row[colEntityID],
				accounts:    make(map[databind.AccountID]Account),
			}
		}

		legalName := sanitization.SanitizeName(row[colLegalName])
		accountId := databind.AccountID(row[colRDBAccountID])
		billingTerms := sanitization.SanitizeIntegerString(row[colBillingTerms])
		if billingTerms == "" || billingTerms == "0" {
			log.Print("Billing Terms is empty/invalid for orgName: " + orgName + " and legalName: " + legalName + " .\nSetting it to 0")
		}
		acc, accExists := org.accounts[accountId]
		if !accExists {
			acc = Account{
				Id:           accountId,
				Name:         legalName,
				DisplayName:  row[colLegalName],
				CustomerId:   row[colCustomerID],
				BillingTerms: billingTerms,
				assetTypes:   make(map[AssetID]AssetType),
			}
		}

		assetId, err := strconv.Atoi(sanitization.SanitizeFloatString(row[colAssetID]))
		if err != nil {
			panic(errors.New(fmt.Sprintf("Column AssetID has invalid value: %s", err)))
		}
		graduatedTier := row[colGraduatedTier]
		tierData := parseTierData(row)

		minimumFee := MinimumFee{
			MinimumFeeType: sanitization.SanitizeName(row[colMinimumFeeType]),
			MinimumCharge:  converter.FromStrToDecimal(row[colMinimumCharge]),
		}

		assetType, exists := acc.assetTypes[AssetID(assetId)]
		if !exists {
			assetType = AssetType{
				Id:            AssetID(assetId),
				Description:   row[colAssetType],
				stakingFees:   make(map[string]StakingFee),
				GraduatedTier: graduatedTier,
				TierData:      tierData,
				MinimumFee:    minimumFee,
			}
		}

		stakingFee := parseStakingFees(row)

		assetType.stakingFees = stakingFee
		acc.assetTypes[AssetID(assetId)] = assetType
		org.accounts[accountId] = acc
		mfr.organizations[msaId] = org
	}

	return mfr
}

func ProcessMfr(ctx context.Context, mfrFile io.Reader, sheetId string, mfrTab string, token string) (*MasterFeeRates, error) {
	// Processing CSV File
	if mfrFile != nil {
		mfr, err := parseMfrFile(mfrFile)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return mfr, nil
	}

	// Processing Google Sheet file
	errorMsg := "Missing parameters: ["
	missingParamError := false
	if len(sheetId) == 0 {
		errorMsg = fmt.Sprintf("%s %s", errorMsg, "SheetId")
		missingParamError = true
	}

	if len(mfrTab) == 0 {
		errorMsg = fmt.Sprintf("%s %s", errorMsg, "MfrTab")
		missingParamError = true
	}

	if len(token) == 0 {
		errorMsg = fmt.Sprintf("%s %s", errorMsg, "Token")
		missingParamError = true
	}
	errorMsg = fmt.Sprintf("%s]", errorMsg)

	if missingParamError {
		return nil, errors.New(errorMsg)
	}

	gSheeetRequest := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)
	mfrValues, err := gSheeetRequest.ReadGoogleSheet(ctx, mfrTab)
	if err != nil {
		msg := fmt.Sprintf("Failed to fetch data from mfr endpoint: %v", err)
		return nil, errors.New(msg)
	}

	mfrArrString, err := googlesheetsutils.GetDataFromSheet(mfrValues, mfrHeaderRow)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	mfr, err := parseMfrGSheet(mfrArrString)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return mfr, nil
}

func parseStakingFees(row []string) map[string]StakingFee {
	return map[string]StakingFee{
		"CELO": {
			AssetName:     "CELO",
			AnchorageFee:  parseFee(row[colCeloFeeAnchorage]),
			ThirdPartyFee: parseFee(row[colCeloFeeThirdParty]),
		},
		"FLOW": {
			AssetName:     "FLOW",
			AnchorageFee:  parseFee(row[colFlowFeeAnchorage]),
			ThirdPartyFee: parseFee(row[colFlowFeeThirdParty]),
		},
		"ROSE": {
			AssetName:     "ROSE",
			AnchorageFee:  parseFee(row[colRoseFeeAnchorage]),
			ThirdPartyFee: parseFee(row[colRoseFeeThirdParty]),
		},
		"APT": {
			AssetName:     "APT",
			AnchorageFee:  parseFee(row[colAptFeeAnchorage]),
			ThirdPartyFee: parseFee(row[colAptFeeThirdParty]),
		},
		"SOL": {
			AssetName:     "SOL",
			ThirdPartyFee: parseFee(row[colSolFeeThirdParty]),
		},
		"OSMO": {
			AssetName:     "OSMO",
			ThirdPartyFee: parseFee(row[colOsmoFeeThirdParty]),
		},
		"HASH": {
			AssetName:     "HASH",
			ThirdPartyFee: parseFee(row[colHashFeeThirdParty]),
		},
		"ATOM": {
			AssetName:     "ATOM",
			ThirdPartyFee: parseFee(row[colAtomFeeThirdParty]),
		},
		"AXL": {
			AssetName:     "AXL",
			ThirdPartyFee: parseFee(row[colAxlFeeThirdParty]),
		},
		"EVMOS": {
			AssetName:     "EVMOS",
			ThirdPartyFee: parseFee(row[colEvmosFeeThirdParty]),
		},
		"SUI": {
			AssetName:     "SUI",
			AnchorageFee:  parseFee(row[colSuiFeeAnchorage]),
			ThirdPartyFee: parseFee(row[colSuiFeeThirdParty]),
		},
		"ETH": {
			AssetName:    "ETH",
			AnchorageFee: parseFee(row[colEthFeeAnchorage]),
		},
	}
}

func parseTierData(row []string) []TierData {
	return []TierData{
		{Floor: getDecimalValue(row[col1stTierFloor]), Rate: getDecimalValue(row[col1stTierRate])},
		{Floor: getDecimalValue(row[col2ndTierFloor]), Rate: getDecimalValue(row[col2ndTierRate])},
		{Floor: getDecimalValue(row[col3rdTierFloor]), Rate: getDecimalValue(row[col3rdTierRate])},
		{Floor: getDecimalValue(row[col4thTierFloor]), Rate: getDecimalValue(row[col4thTierRate])},
		{Floor: getDecimalValue(row[col5thTierFloor]), Rate: getDecimalValue(row[col5thTierRate])},
		{Floor: getDecimalValue(row[col6thTierFloor]), Rate: getDecimalValue(row[col6thTierRate])},
		{Floor: getDecimalValue(row[col7thTierFloor]), Rate: getDecimalValue(row[col7thTierRate])},
		{Floor: getDecimalValue(row[col8thTierFloor]), Rate: getDecimalValue(row[col8thTierRate])},
		{Floor: getDecimalValue(row[col9thTierFloor]), Rate: getDecimalValue(row[col9thTierRate])},
		{Floor: getDecimalValue(row[col10thTierFloor]), Rate: getDecimalValue(row[col10thTierRate])},
	}
}

func getDecimalValue(value string) decimal.Decimal {
	return converter.FromStrToDecimal(value)
}

func parseFee(fee string) decimal.Decimal {
	fee = sanitization.SanitizeFloatString(fee)
	f, err := decimal.NewFromString(fee)
	if err != nil {
		panic(err)
	}
	return f
}

func parseMfrFile(mfrFile io.Reader) (*MasterFeeRates, error) {
	mfrFileData, err := csv.NewReader(mfrFile).ReadAll()
	if err != nil {
		return nil, errors.New("Error while reading MFR data from file.")
	}

	mfr := NewMasterFeeRates(mfrFileData[mfrHeaderRow:])

	return mfr, nil
}

func parseMfrGSheet(mfrGSheetData [][]string) (*MasterFeeRates, error) {
	mfr := NewMasterFeeRates(mfrGSheetData[mfrHeaderRow:])

	return mfr, nil
}
