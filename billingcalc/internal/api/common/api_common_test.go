//go:build !selectTest || unitTest

package common_test

import (
	"testing"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/common"
)

func TestApiCommon(t *testing.T) {
	formData := map[string][]string{
		"firstExternalId": {"123"},
		"invoiceDate":     {"2023-11-23T12:00:00Z"},
		"debug":           {"true"},
	}

	apiParams := common.DefaultAPIParams{}
	err := common.SetValuesFromForm(&apiParams, formData)
	if err != nil {
		t.Fatal(err)
	}
}
