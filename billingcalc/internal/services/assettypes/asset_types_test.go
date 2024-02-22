//go:build !selectTest || unitTest

package assettypes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/assettypes"
)

func TestTypesList(t *testing.T) {
	assetTypeList, err := assettypes.NewAssetTypeList()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("GetTypeByAssetName_Found", func(t *testing.T) {
		assetType, err := assetTypeList.GetTypeByAssetName("ETH", 0)
		if err != nil {
			t.Fatal(err)
		}

		expectedList := []string{"BTC", "ETH"}
		expectedExclusive := true

		assert.Equal(t, expectedList, assetType.Assets, "Expected Assets were %v , but found %v", expectedList, assetType.Assets)
		assert.Equal(t, expectedExclusive, assetType.Exclusive, "Expected Exclude was %v , but found %v", expectedExclusive, assetType.Exclusive)
	})

	t.Run("GetTypeByAssetName_AssetIdNotFound", func(t *testing.T) {
		_, err := assetTypeList.GetTypeByAssetName("BTC", 99)

		expectedErrMsg := "AssetID 99 not found in Asset Types"

		assert.NotNil(t, err, "Expected thrown error")
		assert.Equal(t, expectedErrMsg, err.Error(), "Expected error message was %s , but found %s", expectedErrMsg, err.Error())
	})

	t.Run("GetTypeByAssetName_AssetNameNotFound", func(t *testing.T) {
		_, err := assetTypeList.GetTypeByAssetName("ROSE", 0)

		expectedErrMsg := "Asset Name 'ROSE' not found in Asset Types ID '0'"

		assert.NotNil(t, err, "Expected thrown error")
		assert.Equal(t, expectedErrMsg, err.Error(), "Expected error message was %s , but found %s", expectedErrMsg, err.Error())
	})
}
