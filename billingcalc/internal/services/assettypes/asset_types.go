package assettypes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/static"
)

type AssetTypeList struct {
	assetTypes []AssetType
}

// Asset Type define a group of assets
type AssetType struct {
	AssetId   int      `json:"assetId"`   // ID of the asset type
	Assets    []string `json:"assets"`    // array of assets that are in the group
	Exclusive bool     `json:"exclusive"` // true: Include only this group | false: Not include this group.
}

func NewAssetTypeList() (*AssetTypeList, error) {
	allAssetTypes, err := listAssetTypes()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	assetTypesList := &AssetTypeList{
		assetTypes: make([]AssetType, 0),
	}

	assetTypesList.assetTypes = allAssetTypes

	return assetTypesList, nil
}

// Returns the AssetType of the assetName and assetId combination.
// Returns error if the assetId is not found or if the assetName is not in assetType.Assets array.
func (atl *AssetTypeList) GetTypeByAssetName(assetName string, assetId int) (*AssetType, error) {
	typeFound := searchForId(assetId, atl.assetTypes)
	ErrAssetNameNotFound := errors.New("asset name not found in asset type")

	if typeFound == nil {
		return nil, errors.New(fmt.Sprintf("AssetID %v not found in Asset Types", assetId))
	}

	if assetInAssetType(assetName, typeFound.Assets) && typeFound.Exclusive == false {
		return nil, ErrAssetNameNotFound
	}

	if assetInAssetType(assetName, typeFound.Assets) && typeFound.Exclusive == true {
		return typeFound, nil
	}

	if !assetInAssetType(assetName, typeFound.Assets) && typeFound.Exclusive == false {
		return typeFound, nil
	}

	return nil, errors.New(fmt.Sprintf("Asset Name '%s' not found in Asset Types ID '%v'", assetName, assetId))
}

func listAssetTypes() ([]AssetType, error) {
	fileContent, err := static.Files.ReadFile("asset_types.json")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading file asset_types.json: %v", err))
	}

	assetTypes := make([]AssetType, 0)
	if err := json.Unmarshal(fileContent, &assetTypes); err != nil {
		return nil, errors.New(fmt.Sprintf("Error in json unmarshal: %v", err))
	}

	return assetTypes, nil
}

func searchForId(assetId int, assetTypes []AssetType) *AssetType {
	for _, ele := range assetTypes {
		if ele.AssetId == assetId {
			return &ele
		}
	}

	return nil
}

func assetInAssetType(assetName string, assets []string) bool {
	for _, asset := range assets {
		if asset == assetName {
			return true
		}
	}
	return false
}
