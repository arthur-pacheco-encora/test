package slices

import (
	"errors"
)

type FilterStruct struct {
	Column    int
	Value     string
	Reference string
}

type ResponseStruct struct {
	TableSlice [][]string
	Filter     FilterStruct
	Siblings   []string
	Error      error
}

func AppendRef(ref, s string) string {
	if s == "" {
		return ref
	}
	if ref == "" {
		return s
	}
	return ref + ":" + s
}

func TriggerJobs(tbl [][]string, col int, ref string, router chan *ResponseStruct) int {
	jobCount := 0
	entries := Uniq(tbl, col)
	for _, entry := range entries {
		jobCount++
		go SliceTable(
			FilterStruct{
				col,
				entry,
				AppendRef(ref, entry),
			},
			tbl,
			router,
		)
	}
	jobCount++
	go func() {
		router <- &ResponseStruct{
			Filter: FilterStruct{
				Reference: ref,
			},
			Siblings: entries,
		}
	}()

	return jobCount
}

func Uniq(table [][]string, col int) []string {
	found := make(map[string]bool)
	unique := []string{}
	for _, row := range table {
		key := row[col]
		if key == "" {
			continue
		}
		if !found[key] {
			found[key] = true
			unique = append(unique, key)
		}
	}
	return unique
}

func SliceTable(filter FilterStruct, table [][]string, router chan *ResponseStruct) {
	startIdx := -1
	endIdx := -1
	var tSlice [][]string

	for idx, row := range table {
		if row[filter.Column] == filter.Value {
			if startIdx == -1 {
				startIdx = idx
			}
		} else {
			if startIdx != -1 {
				endIdx = idx
				break
			}
		}
	}

	if startIdx < 0 {
		router <- &ResponseStruct{
			Error: errors.New("Provided filter does not match any entry."),
		}
		return
	}

	if endIdx == -1 {
		tSlice = table[startIdx:]
	} else {
		tSlice = table[startIdx:endIdx]
	}

	router <- &ResponseStruct{
		Filter:     filter,
		TableSlice: tSlice,
		Error:      nil,
	}
}

// Insert a item in a slice at a given index
func Insert(target []string, index int, value string) []string {
	if len(target) == index {
		return append(target, value)
	}

	if len(target) < index {
		newTarget := make([]string, index+1)
		copy(newTarget, target)
		target = newTarget
	}

	target[index] = value
	return target
}
