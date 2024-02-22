package handlers

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/common"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/converter"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

type CsvAPIParams struct {
	common.DefaultAPIParams
	MfrFile            string `schema:"mfr,required"`
	RewardsFile        string `schema:"rewards,required"`
	UnclaimedFile      string `schema:"unclaimed,required"`
	BalanceAdjustments string `schema:"balanceAdjustments,required"`
	OperationsStatuses string `schema:"operationsStatuses,required"`
	DailyBalances      string `schema:"dailyBalances,required"`
}

func CalcFeesFromCsv(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(math.MaxInt64)
	if err != nil {
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	csvParams := &CsvAPIParams{}
	err = common.SetValuesFromForm(csvParams, r.PostForm)
	if err != nil {
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	debug.Init(csvParams.Debug)
	debug.NewMessage("Starting CalculateFromCsv")

	debug.NewMessage("Values From Form: " + fmt.Sprintf("%#v", csvParams))

	mfr := converter.FromStringBase64ToIoReader(csvParams.MfrFile)

	rewards := converter.FromStringBase64ToIoReader(csvParams.RewardsFile)

	unclaimed := converter.FromStringBase64ToIoReader(csvParams.UnclaimedFile)

	balanceAdjustments := converter.FromStringBase64ToIoReader(csvParams.BalanceAdjustments)

	operationsStatuses := converter.FromStringBase64ToIoReader(csvParams.OperationsStatuses)

	dailyBalances := converter.FromStringBase64ToIoReader(csvParams.DailyBalances)

	result, err := fees.CalculateFromCsv(r.Context(), mfr, rewards, unclaimed, balanceAdjustments, operationsStatuses, dailyBalances, csvParams.FirstExternalId, csvParams.InvoiceDate)
	if err != nil {
		debug.NewMessage("Error CalculateFromCsv: " + err.Error())
		common.WriteErr(w, errors.New("Failed to calculate fees."))
		debug.NewMessage("Finishing CalculateFromCsv")
		return
	}

	debug.NewMessage("Success CalculateFromCsv")
	debug.NewMessage("Finishing CalculateFromCsv")

	resp := &common.Response{
		Data:  result.Summary,
		Warn:  result.Warns,
		Debug: debug.GetAllMessages(),
		Err:   "",
	}
	resp.Write(w)
}
