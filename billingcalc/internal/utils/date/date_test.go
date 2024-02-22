//go:build !selectTest || unitTest

package date_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/date"
)

func TestDateUtil(t *testing.T) {
	t.Run("Test_NumberOfDaysInMonth", func(t *testing.T) {
		actual := date.NumberOfDaysInMonth(6, 2023)
		expected := int64(30)

		assert.Equal(t, expected, actual)
	})
}

func TestNumberOfDaysInMonth(t *testing.T) {
	tests := []struct {
		date     time.Time
		expected int
	}{
		{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), 29},
		{time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC), 28},
		{time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), 31},
		{time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC), 30},
	}

	for _, test := range tests {
		got := date.NumberOfDaysInTime(test.date)
		assert.Equal(t, test.expected, got, "Number Of Days In Month(%v) should be %v, got %v", test.date, test.expected, got)
	}
}
