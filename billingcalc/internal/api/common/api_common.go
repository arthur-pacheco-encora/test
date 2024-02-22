package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/schema"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

type DefaultAPIParams struct {
	FirstExternalId int       `schema:"firstExternalId,required"`
	InvoiceDate     time.Time `schema:"invoiceDate,required"`
	Debug           bool      `schema:"debug,required"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Warn  interface{} `json:"warn"`
	Debug interface{} `json:"debug"`
	Err   string      `json:"err"`
}

func (r *Response) Write(w http.ResponseWriter) {
	m, err := json.Marshal(r)
	if err != nil {
		WriteErr(w, errors.New("Error while encoding response."))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(m))
	debug.NewMessage(fmt.Sprintf("http response wrote: %s", string(m)))
}

func WriteErr(w http.ResponseWriter, err error) {
	resp := &Response{
		Data:  "",
		Warn:  "",
		Debug: debug.GetAllMessages(),
		Err:   err.Error(),
	}
	resp.Write(w)
}

func SetValuesFromForm(s interface{}, formData map[string][]string) error {
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]byte{}, func(s string) reflect.Value {
		return reflect.ValueOf([]byte(s))
	})
	decoder.RegisterConverter(time.Time{}, func(s string) reflect.Value {
		return reflect.ValueOf(converter.FromStrToTimeStamp(s))
	})

	err := decoder.Decode(s, formData)
	if err != nil {
		msg := "Error decoding parameters:" + err.Error()
		return errors.New(msg)
	}

	return nil
}
