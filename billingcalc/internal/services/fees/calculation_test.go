//go:build !selectTest || unitTest

package fees_test

import (
	"fmt"
	"testing"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees"
)

var (
	currentExternalID int
	firstAcc          bool
)

func TestGenerateInvoiceNumber(t *testing.T) {
	getInvoiceNumber := fees.GenInvoiceNumber(1)

	invoice1 := getInvoiceNumber("15", "1001")
	expected1 := "ADB-1"
	if invoice1 != expected1 {
		t.Errorf("expected %s, but got %s", expected1, invoice1)
	}

	invoice2 := getInvoiceNumber("33", "1002")
	expected2 := "ABS-2"
	if invoice2 != expected2 {
		t.Errorf("expected %s, but got %s", expected2, invoice2)
	}

	invoice3 := getInvoiceNumber("100", "1003")
	expected3 := "Acronym Not Found"
	if invoice3 != expected3 {
		t.Errorf("expected %s, but got %s", expected3, invoice3)
	}

	invoice4 := getInvoiceNumber("15", "1001")
	expected4 := "ADB-1"
	if invoice4 != expected4 {
		t.Errorf("expected %s, but got %s", invoice4, expected4)
	}
}

func TestIntegration(t *testing.T) {
	firstAcc = false
	currentExternalID := fees.GenExternalID(currentExternalID, firstAcc)
	externalID := fmt.Sprintf("%d", currentExternalID)

	getInvoiceNumber := fees.GenInvoiceNumber(currentExternalID)

	entityID := "15"
	invoice := getInvoiceNumber(entityID, externalID)
	expected := "ADB-" + externalID

	if invoice != expected {
		t.Errorf("expected %s, but got %s", expected, invoice)
	}
}
