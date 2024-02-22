package file

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

type GetDatabindFromDataParams[T any] struct {
	ReportName string
	HeaderRow  int
	File       io.Reader
	Parser     func(data [][]string) T
}

func GetDatabindFromData[T any](params GetDatabindFromDataParams[T]) (T, error) {
	debug.NewMessage(fmt.Sprintf("Start parse \"%s\" file.", params.ReportName))

	reader, err := csv.NewReader(params.File).ReadAll()
	if err != nil {
		var empty T
		return empty, errors.New(fmt.Sprintf("Error while reading %s sheet from file.", params.ReportName))
	}

	databinds := params.Parser(reader[params.HeaderRow:])

	debug.NewMessage(fmt.Sprintf("End parse \"%s\" file.", params.ReportName))

	return databinds, nil
}
