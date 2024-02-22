package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/common"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

type BigQueryAPIParams struct {
	common.DefaultAPIParams
	MfrFile     string    `schema:"mfr"`
	MfrTab      string    `schema:"mfrTab"`
	SheetId     string    `schema:"sheetID"`
	Token       string    `schema:"token"`
	PeriodBegin time.Time `schema:"periodBegin"`
	PeriodEnd   time.Time `schema:"periodEnd"`
}

func CalcFeesFromBigQuery(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(math.MaxInt64)
	if err != nil {
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	bqParams := &BigQueryAPIParams{}
	err = common.SetValuesFromForm(bqParams, r.PostForm)
	if err != nil {
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	mfr := converter.FromStringBase64ToIoReader(bqParams.MfrFile)

	debug.Init(bqParams.Debug)
	debug.NewMessage("Starting CalculateFromBigQuery")

	debug.NewMessage("Parameters: " + fmt.Sprintf("%#v", bqParams))

	result, err := fees.CalculateFromBigQuery(r.Context(), mfr, bqParams.SheetId, bqParams.MfrTab, bqParams.Token, bqParams.PeriodBegin, bqParams.PeriodEnd, bqParams.FirstExternalId, bqParams.InvoiceDate)
	if err != nil {
		debug.NewMessage("Error CalculateFromBigQuery: " + err.Error())
		common.WriteErr(w, errors.New("Failed to calculate fees."))
		debug.NewMessage("Finishing CalculateFromBigQuery")
		return
	}

	debug.NewMessage("Success CalculateFromBigQuery")
	debug.NewMessage("Finishing CalculateFromBigQuery")

	resp := &common.Response{
		Data:  result.Summary,
		Warn:  result.Warns,
		Debug: debug.GetAllMessages(),
		Err:   "",
	}
	resp.Write(w)
}
