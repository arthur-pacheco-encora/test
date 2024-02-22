//go:build !selectTest || unitTest

package googlesheetsutils_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/balanceadjustments"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/operationsstatuses"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/rewards"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/databind/ubalances"
	"github.com/anchorlabsinc/anchorage/source/go/service/billingcalc/internal/services/googlesheetsutils"
)

func setEnv(t *testing.T, testServer *httptest.Server) {
	t.Setenv("SHEETS_API_BASE_URL", testServer.URL)
}

func TestReadGoogleSheet(t *testing.T) {
	sheetId := "sheet1"
	token := "token1"
	tab := "tab1"
	body := []byte(`{"values": [["a","b","c"],["1","2","3"],["x","y","z"]]}`)
	expectedResult := [][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
		{"x", "y", "z"},
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := fmt.Sprintf("%s/%s/values/%s", r.Host, sheetId, tab)
		actualPath := fmt.Sprintf("%s%s", r.Host, r.URL.Path)
		if actualPath != expectedPath {
			t.Errorf("Expected path was %s, got %s", expectedPath, actualPath)
		}

		expectedHeader := fmt.Sprintf("Bearer %s", token)

		if values, ok := r.Header["Authorization"]; ok {
			if values[0] != expectedHeader {
				t.Errorf("Expected Authorization header was %s, got %s", expectedHeader, values[0])
			}
		} else {
			t.Error("Authorization not found in the Header")
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(body)
		if err != nil {
			t.Errorf("Unexpected error on Write: %s", err.Error())
		}
	}))

	setEnv(t, testServer)

	g := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)
	result, err := g.ReadGoogleSheet(context.Background(), tab)
	if err != nil {
		t.Errorf("Unexpected error on ReadGoogleSheet: %s", err.Error())
	}

	assert.Equal(t, expectedResult, result, "Expected result was %v, but got: %v", expectedResult, result)
}

func TestReadGoogleSheet_UnmarshalError(t *testing.T) {
	sheetId := "sheet1"
	token := "token1"
	tab := "tab1"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{this is not a json}"))
		if err != nil {
			t.Errorf("Unexpected error on Write: %s", err.Error())
		}
	}))

	expectedMsg := "error unmarshalling JSON response: invalid character 't' looking for beginning of object key string"
	setEnv(t, testServer)
	g := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)
	_, err := g.ReadGoogleSheet(context.Background(), tab)

	assert.EqualError(t, err, expectedMsg, "Expected Unmarshall error: %s, but got: %s", expectedMsg, err.Error())
}

func TestReadGoogleSheet_Not200(t *testing.T) {
	statusCode := http.StatusBadRequest
	sheetId := "sheet1"
	token := "token1"
	tab := "tab1"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		_, err := w.Write([]byte("bad request"))
		if err != nil {
			t.Errorf("Unexpected error on Write: %s", err.Error())
		}
	}))

	expectedMsg := fmt.Sprintf("API responded with 400 status code for endpoint: %s/%s/values/%s", testServer.URL, sheetId, tab)
	setEnv(t, testServer)
	g := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)
	_, err := g.ReadGoogleSheet(context.Background(), tab)

	assert.EqualError(t, err, expectedMsg, "Expected Status Code error: %s, but got: %s", expectedMsg, err.Error())
}

func TestReadGoogleSheet_HttpRequestError(t *testing.T) {
	sheetId := "sheet1"
	token := "token1"
	tab := "tab1"
	badUrl := "http:\\url.com"
	t.Setenv("SHEETS_API_BASE_URL", badUrl)

	httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("http request error"))
		if err != nil {
			t.Errorf("Unexpected error on Write: %s", err.Error())
		}
	}))

	expectedMsg := "error making HTTP request: Get \"http:\\\\url.com/sheet1/values/tab1\": http: no Host in request URL"
	g := googlesheetsutils.NewGoogleSheetRequest(sheetId, token)
	_, err := g.ReadGoogleSheet(context.Background(), tab)

	assert.EqualError(t, err, expectedMsg, "Expected Status Code error: %s, but got: %s", expectedMsg, err.Error())
}

func TestReadGoogleSheet_GetDatabindFromTab(t *testing.T) {
	ctx := context.Background()

	t.Run("Test_Parse_Rewards", func(t *testing.T) {
		fakeDataResponse := createEmptyTable[string](8, 18)
		fakeDataResponse = append(fakeDataResponse, []string{"Commonwealth Asset Management LP", "Commonwealth Asset Management Co-Investment Fund II LP", "FLOW 2 CW2", "2023-06-28 19:03:01", "2023-06-28 19:03:01", "Delegation Reward", "FLOW", "41e8672b3ecb08bf91c6dc430fa6bd8af1e5a43547efdb59dc1814eeca278480", "0x05406290e744c813", "", "COMPLETE", "0", "0", "0", "0", "0", "0", "6/28/2023", "7f24f990eb2e6d70cdb9c30a91ad8090ca403bd5ce4981e03ac6868dc681bc83"})
		headerRows := 8
		parser := rewards.NewRewards

		params := googlesheetsutils.GetDatabindFromSheetTabParams[*rewards.Rewards]{
			SheetRequest: getMockGoogleSheetRequest(t, fakeDataResponse),
			TabName:      "rewards",
			ReportName:   "Balance Adjustments Report",
			HeaderRow:    headerRows,
			Parser:       parser,
		}
		databind, err := googlesheetsutils.GetDatabindFromSheetTab[*rewards.Rewards](ctx, params)

		// Get the Rewards struct for the first line and to be compared
		expected := parser(fakeDataResponse[headerRows:])

		assert.Nil(t, err)
		assert.Equal(t, expected, databind)
	})

	t.Run("Test_Parse_Unclaimed_Balances", func(t *testing.T) {
		fakeDataResponse := createEmptyTable[string](1, 11)
		fakeDataResponse = append(fakeDataResponse, []string{"45107AXLAxelar Foundation", "2023-06-30", "axelar1hmssmjjun20znd3nnfx462hahj4w4tej06khmy", "AXL", "1c1967712da3513afb6ebdc7d64cf9b283a53f8c9fde7705068e599161c54ac4", "DELEGATION_REWARDS", "0.147223", "2023-06-30 23:56:39", "Axelar Foundation", "0.3895", "0.0573433585"})

		headerRows := 1
		parser := ubalances.NewUnclaimedBalances

		params := googlesheetsutils.GetDatabindFromSheetTabParams[*ubalances.UnclaimedBalances]{
			SheetRequest: getMockGoogleSheetRequest(t, fakeDataResponse),
			TabName:      "unclaimed",
			ReportName:   "Unclaimed Balances Report",
			HeaderRow:    headerRows,
			Parser:       parser,
		}
		databind, err := googlesheetsutils.GetDatabindFromSheetTab[*ubalances.UnclaimedBalances](ctx, params)

		// Get the Unclaimed Balances struct for the first line and compared
		expected := parser(fakeDataResponse[headerRows:])

		assert.Nil(t, err)
		assert.Equal(t, expected, databind)
	})

	t.Run("Test_Parse_Operations_Statuses", func(t *testing.T) {
		fakeDataResponse := createEmptyTable[string](1, 25)
		fakeDataResponse = append(fakeDataResponse, []string{"Anchorage", "1", "9.29", "HOLD", "0", "ATOM", "1cb6771c24450099fc5bdc66bafe2314", "2023-06-30", "No", "2023-06-30", "06a26df605a0776e4ba08bc4248678ce8575cbea29abb46c7eed5a7e8afaedaf", "Anchorage", "0", "0", "42", "0", "0", "9.288", "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", "", "Grace Noah’s Vault", "1", "9ed4bea1540311595c544a5b4a69d5a3", "9dc5526a3dc5a62c3c704e0dc85b2cc9", "0.089"})

		headerRows := 1
		parser := operationsstatuses.NewOperationsStatuses

		params := googlesheetsutils.GetDatabindFromSheetTabParams[*operationsstatuses.OperationsStatuses]{
			SheetRequest: getMockGoogleSheetRequest(t, fakeDataResponse),
			TabName:      "operations_statuses",
			ReportName:   "Operations Statuses Report",
			HeaderRow:    headerRows,
			Parser:       parser,
		}
		databind, err := googlesheetsutils.GetDatabindFromSheetTab[*operationsstatuses.OperationsStatuses](ctx, params)

		// Get the Operations Statuses struct for the first line and compared
		expected := parser(fakeDataResponse[headerRows:])

		assert.Nil(t, err)
		assert.Equal(t, expected, databind)
	})

	t.Run("Test_Parse_Balance_Adjustments", func(t *testing.T) {
		fakeDataResponse := createEmptyTable[string](1, 25)
		fakeDataResponse = append(fakeDataResponse, []string{"Visa", "Visa International Service Association", "NFTs", "2023-06-16 15:54:49", "2023-06-16 15:54:49", "Balance Adjustment", "POLY", "0xeadf6e3781849dc0d279c9f32ae5d77f58e96ece3cdf8723abd2d6aff477a66b", "", "0x98D7Bc088Fc1E7a2f9B3536fCFBCA405f902467e", "COMPLETE", "missed the decimal point in the first balance adjustment", "0", "0", "0", "0", "0", "-1335.663", "N"})

		headerRows := 1
		parser := balanceadjustments.NewBalanceAdjustments

		params := googlesheetsutils.GetDatabindFromSheetTabParams[*balanceadjustments.BalanceAdjustments]{
			SheetRequest: getMockGoogleSheetRequest(t, fakeDataResponse),
			TabName:      "balance_adjustments",
			ReportName:   "Balance Adjustments Report",
			HeaderRow:    headerRows,
			Parser:       parser,
		}
		databind, err := googlesheetsutils.GetDatabindFromSheetTab[*balanceadjustments.BalanceAdjustments](ctx, params)

		// Get the Balance Adjustments struct for the first line and compared
		expected := parser(fakeDataResponse[headerRows:])

		assert.Nil(t, err)
		assert.Equal(t, expected, databind)
	})

	t.Run("Test_Empty_Data", func(t *testing.T) {
		var fakeDataResponse [][]string
		parser := operationsstatuses.NewOperationsStatuses
		expected := "Error while parsing Operations Statuses Report: No Operations Statuses found."

		params := googlesheetsutils.GetDatabindFromSheetTabParams[*operationsstatuses.OperationsStatuses]{
			SheetRequest: getMockGoogleSheetRequest(t, fakeDataResponse),
			TabName:      "wrong_tab_name",
			ReportName:   "Operations Statuses Report",
			HeaderRow:    1,
			Parser:       parser,
		}
		databind, err := googlesheetsutils.GetDatabindFromSheetTab[*operationsstatuses.OperationsStatuses](ctx, params)

		assert.Nil(t, databind)
		assert.Equal(t, expected, err.Error())
	})
}

func getMockGoogleSheetRequest(t *testing.T, response [][]string) *googlesheetsutils.GoogleSheetRequest {
	type MockResponse struct {
		Values [][]string `json:"Values"`
	}

	fake := &MockResponse{Values: response}
	fakeJson, err := json.Marshal(fake)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(string(fakeJson)))
		if err != nil {
			t.Errorf("Unexpected error on Write: %s", err.Error())
		}
	}))

	setEnv(t, testServer)

	return googlesheetsutils.NewGoogleSheetRequest("fake_sheet_id", "fake_token")
}

func createEmptyTable[T any](numberOfLines int, numberOfColumns int) [][]T {
	var data [][]T

	for i := 0; i < numberOfLines; i++ {
		newLine := make([]T, numberOfColumns)
		data = append(data, newLine)
	}

	return data
}
