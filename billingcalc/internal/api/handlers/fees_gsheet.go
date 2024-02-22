package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/common"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/fees"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
)

type GSheetAPIParams struct {
	common.DefaultAPIParams
	MfrTab                string `schema:"mfrTab"`
	RewardsTab            string `schema:"rewardsTab"`
	UnclaimedBalancesTab  string `schema:"unclaimedTab"`
	BalanceAdjustmentsTab string `schema:"balanceAdjustmentTab"`
	OperationsStatusesTab string `schema:"operationsStatusesTab,required"`
	DailyBalancesTab      string `schema:"dailyBalancesTab,required"`
	SheetId               string `schema:"sheetID"`
	Token                 string `schema:"token"`
}

func CalcFeesFromGSheets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	debug.NewMessage("Parsing form")
	err := r.ParseForm()
	if err != nil {
		debug.NewMessage(err.Error())
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	params := &GSheetAPIParams{}
	err = common.SetValuesFromForm(params, r.PostForm)
	if err != nil {
		debug.NewMessage(err.Error())
		common.WriteErr(w, errors.New(err.Error()))
		return
	}

	debug.Init(params.Debug)
	debug.NewMessage("Starting CalculateFromGSheets")

	debug.NewMessage("Parameters: " + fmt.Sprintf("%#v", params))

	result, err := fees.CalculateFromGSheets(r.Context(), params.SheetId, params.MfrTab, params.RewardsTab, params.UnclaimedBalancesTab, params.BalanceAdjustmentsTab, params.OperationsStatusesTab, params.DailyBalancesTab, params.Token, params.FirstExternalId, params.InvoiceDate)
	if err != nil {
		debug.NewMessage("Error CalculateFromGSheets: " + err.Error())
		common.WriteErr(w, errors.New("Failed to calculate fees."))
		debug.NewMessage("Finishing CalculateFromGSheets")
		return
	}

	debug.NewMessage("Success CalculateFromGSheets")
	debug.NewMessage("Finishing CalculateFromGSheets")

	resp := &common.Response{
		Data:  result.Summary,
		Warn:  result.Warns,
		Debug: debug.GetAllMessages(),
		Err:   "",
	}
	resp.Write(w)
}
