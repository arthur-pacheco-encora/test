package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/common"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/api/handlers"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := &common.Response{
		Data:  "Welcome! V1.0",
		Warn:  "",
		Debug: "",
		Err:   "",
	}
	resp.Write(w)
}

func Install(r *httprouter.Router) {
	r.GET("/", Index)
	r.GET("/assets", handlers.GetAssets)
	r.POST("/fees", handlers.CalcFeesFromGSheets)
	r.POST("/fees-csv", handlers.CalcFeesFromCsv)
	r.POST("/fees-bq", handlers.CalcFeesFromBigQuery)
	r.POST("/assets", handlers.UpdateAssets)
}
