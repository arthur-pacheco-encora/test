//go:build !selectTest || unitTest

package main_test

import (
	"encoding/csv"
	"os"
	"testing"

	main "github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/cmd/genfakerewards"
)

func TestCSVCreation(t *testing.T) {
	main.RunGenFakeRewards()

	_, err := os.Stat("rewards.csv")
	if err != nil {
		t.Fatalf("failed creating file: %s", err)
	}
}

func TestCSVStructure(t *testing.T) {
	main.RunGenFakeRewards()

	file, err := os.Open("rewards.csv")
	if err != nil {
		t.Fatalf("failed opening file: %s", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Logf("Failed to close connection: %v", closeErr)
		}
	}()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed reading file: %s", err)
	}

	expectedHeader := []string{"organization", "account", "operation", "asset", "anch_val", "non_anch_val"}
	if len(records[0]) != len(expectedHeader) {
		t.Fatalf("header length mismatch. got=%d want=%d", len(records[0]), len(expectedHeader))
	}

	for i, header := range expectedHeader {
		if records[0][i] != header {
			t.Errorf("header mismatch at %d. got=%s want=%s", i, records[0][i], header)
		}
	}

	if len(records) != main.CsvLines+1 {
		t.Errorf("number of lines mismatch. got=%d want=%d", len(records), main.CsvLines+1)
	}
}
