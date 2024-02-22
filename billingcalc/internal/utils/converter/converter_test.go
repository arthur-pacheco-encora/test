//go:build !selectTest || unitTest

package converter_test

import (
	"encoding/csv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
)

func TestFromStrToTimeStamp(t *testing.T) {
	t.Run("Test_Date_Without_Time", func(t *testing.T) {
		date := "2023-06-23"
		expected := time.Date(2023, time.June, 23, 0, 0, 0, 0, time.UTC)
		actual := converter.FromStrToTimeStamp(date)

		assert.Equal(t, expected, actual)
	})

	t.Run("Test_Date_With_Time", func(t *testing.T) {
		date := "2023-06-23 10:30:10"
		expected := time.Date(2023, time.June, 23, 0, 0, 0, 0, time.UTC)
		actual := converter.FromStrToTimeStamp(date)

		assert.Equal(t, expected, actual)
	})

	t.Run("Test_Date_With_One_Digit", func(t *testing.T) {
		date := "2023-6-1"
		expected := time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)
		actual := converter.FromStrToTimeStamp(date)

		assert.Equal(t, expected, actual)
	})

	t.Run("Test_American_Date", func(t *testing.T) {
		date := "06/23/2023"
		expected := time.Date(2023, time.June, 23, 0, 0, 0, 0, time.UTC)
		actual := converter.FromStrToTimeStamp(date)
		assert.Equal(t, expected, actual)
	})

	t.Run("Test_American_Date_With_One_Digit", func(t *testing.T) {
		date := "6/1/2023"
		expected := time.Date(2023, time.June, 1, 0, 0, 0, 0, time.UTC)
		actual := converter.FromStrToTimeStamp(date)
		assert.Equal(t, expected, actual)
	})
}

func TestFromStringBase64ToIoReader(t *testing.T) {
	t.Run("Test_with_CSV", func(t *testing.T) {
		base64 := "TGluZSAxCkxpbmUgMgo="
		file := converter.FromStringBase64ToIoReader(base64)
		content, _ := csv.NewReader(file).ReadAll()

		expectedLine1 := "Line 1"
		expectedLine2 := "Line 2"

		assert.Equal(t, expectedLine1, content[0][0])
		assert.Equal(t, expectedLine2, content[1][0])
	})

	t.Run("Test_return_nil_passing_empty_String", func(t *testing.T) {
		base64 := ""
		file := converter.FromStringBase64ToIoReader(base64)

		assert.Equal(t, file, nil)
	})
}
