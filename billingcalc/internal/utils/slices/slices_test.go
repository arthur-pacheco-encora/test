//go:build !selectTest || unitTest

package slices_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/slices"
)

func TestAppendRef(t *testing.T) {
	cases := []struct {
		refStr      string
		strToAppend string
		expected    string
	}{
		{"a", "b", "a:b"},
		{"a", "", "a"},
		{"", "b", "b"},
	}

	for _, tt := range cases {
		actual := slices.AppendRef(tt.refStr, tt.strToAppend)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestUniq(t *testing.T) {
	table := [][]string{
		{"A1", "Value1"},
		{"A2", "Value2"},
		{"A3", "Value1"},
	}

	expected := []string{"Value1", "Value2"}
	actual := slices.Uniq(table, 1)

	assert.Equal(t, expected, actual)
}

func TestSliceTable(t *testing.T) {
	table := [][]string{
		{"A1", "Value1"},
		{"A2", "Value2"},
		{"A3", "Value1"},
	}

	resultChannel := make(chan *slices.ResponseStruct)

	t.Run("Slice table test success", func(t *testing.T) {
		filter := slices.FilterStruct{
			Column:    1,
			Value:     "Value1",
			Reference: "",
		}

		expected := slices.ResponseStruct{
			TableSlice: [][]string{{"A1", "Value1"}},
			Filter:     filter,
			Siblings:   nil,
			Error:      nil,
		}

		go slices.SliceTable(filter, table, resultChannel)

		actual := <-resultChannel

		assert.Equal(t, expected.TableSlice, actual.TableSlice)
		assert.Equal(t, expected.Filter, actual.Filter)
		assert.Equal(t, expected.Siblings, actual.Siblings)
		assert.Equal(t, expected.Error, actual.Error)
	})

	t.Run("Slice table test not match", func(t *testing.T) {
		filterError := slices.FilterStruct{
			Column:    1,
			Value:     "A1",
			Reference: "",
		}

		expected := slices.ResponseStruct{
			TableSlice: [][]string{{"A1", "Value1"}},
			Filter:     filterError,
			Siblings:   nil,
			Error:      nil,
		}

		expected.Filter = filterError
		expected.Error = errors.New("Provided filter does not match any entry.")
		go slices.SliceTable(filterError, table, resultChannel)

		actualError := <-resultChannel
		assert.Equal(t, expected.Error, actualError.Error)
	})
}

func TestInsert(t *testing.T) {
	cases := []string{
		"a", "b", "c",
	}

	newCaseValue := "e"

	newCases := slices.Insert(cases, 5, newCaseValue)

	assert.Equal(t, cases[0], newCases[0])
	assert.Equal(t, cases[1], newCases[1])
	assert.Equal(t, cases[2], newCases[2])
	assert.Equal(t, newCaseValue, newCases[5])
}
