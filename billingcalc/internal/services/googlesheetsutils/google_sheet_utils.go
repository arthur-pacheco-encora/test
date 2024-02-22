package googlesheetsutils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/debug"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/utils/sanitization"
)

type GoogleSheetRequest struct {
	baseURL string
	sheetId string
	token   string
}

type GoogleSheetResponse struct {
	Values [][]string `json:"values"`
}

const env_sheets_base_url = "SHEETS_API_BASE_URL"

func NewGoogleSheetRequest(sheetId string, token string) *GoogleSheetRequest {
	g := &GoogleSheetRequest{
		sheetId: sheetId,
		token:   token,
	}

	g.setBaseUrl()

	return g
}

func GetDataFromSheet(input [][]string, headerRow int) ([][]string, error) {
	data := sanitization.SanitizeGoogleSheet(input, headerRow)
	if len(data) == 0 {
		err := errors.New("Error while getting data Google Sheet: No data found")
		return nil, err
	}
	return data, nil
}

type GetDatabindFromSheetTabParams[T any] struct {
	SheetRequest *GoogleSheetRequest
	TabName      string
	ReportName   string
	HeaderRow    int
	Parser       func(data [][]string) T
}

func GetDatabindFromSheetTab[T any](ctx context.Context, params GetDatabindFromSheetTabParams[T]) (T, error) {
	debug.NewMessage(fmt.Sprintf("Start parse \"%s\" file from sheet.", params.ReportName))

	var emptyReturn T

	values, err := params.SheetRequest.ReadGoogleSheet(ctx, params.TabName)
	if err != nil {
		return emptyReturn, errors.New(err.Error())
	}

	table := sanitization.SanitizeGoogleSheet(values, params.HeaderRow)

	if len(table) == 0 {
		errMessage := fmt.Sprintf("Error while parsing %s: No %s found.", params.ReportName, strings.Replace(params.ReportName, " Report", "", -1))
		err := errors.New(errMessage)

		debug.NewMessage(fmt.Sprintf("End parse \"%s\" file from sheet. %v", params.ReportName, errMessage))
		return emptyReturn, err
	}

	newData := params.Parser(table)

	debug.NewMessage(fmt.Sprintf("End parse \"%s\" file.", params.ReportName))
	return newData, nil
}

func (g *GoogleSheetRequest) ReadGoogleSheet(ctx context.Context, tab string) ([][]string, error) {
	client := &http.Client{}
	endpoint := g.formatRequestUrl(tab)

	httpRequest, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		msg := fmt.Sprintf("failed to create new HTTP request: %v", err)
		return nil, errors.New(msg)
	}

	bearer := "Bearer " + g.token
	httpRequest.Header.Add("Authorization", bearer)

	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		msg := fmt.Sprintf("error making HTTP request: %v", err)
		return nil, errors.New(msg)
	}
	defer func() {
		if closeErr := httpResponse.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if httpResponse.StatusCode != 200 {
		msg := fmt.Sprintf("API responded with %d status code for endpoint: %s", httpResponse.StatusCode, endpoint)
		return nil, errors.New(msg)
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		msg := fmt.Sprintf("failed to read response body: %v", err)
		return nil, errors.New(msg)
	}

	var sheetResponse GoogleSheetResponse
	if err := json.Unmarshal(body, &sheetResponse); err != nil {
		msg := fmt.Sprintf("error unmarshalling JSON response: %v", err)
		return nil, errors.New(msg)
	}

	return sheetResponse.Values, nil
}

func (g *GoogleSheetRequest) formatRequestUrl(tab string) string {
	url := fmt.Sprintf("%s/%s/values/%s", g.baseURL, g.sheetId, tab)
	return url
}

func (g *GoogleSheetRequest) setBaseUrl() {
	baseURL := os.Getenv(env_sheets_base_url)
	if baseURL == "" {
		baseURL = "https://sheets.googleapis.com/v4/spreadsheets"
	}

	g.baseURL = baseURL
}
