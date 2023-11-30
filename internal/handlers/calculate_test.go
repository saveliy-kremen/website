package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestCalculateHandle(t *testing.T) {
	testCases := []struct {
		testName string
		testData map[string]int
	}{
		{"testOk", map[string]int{"a": 10, "b": 20}},
		{"testNeg", map[string]int{"a": -10, "b": 20}},
		{"testLoss", map[string]int{"a": 10}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			jsonValue, _ := json.Marshal(testCase.testData)
			req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(jsonValue))
			if err != nil {
				t.Fatal(err)
			}
			rr := newRequestRecorder(req, "POST", "/calculate", MiddlewareCalculate(CalculateHandler))
			if testCase.testName == "testOk" {
				if rr.Code != 200 {
					t.Error("Expected response code to be 200")
				}
				if rr.Body.String() != string(jsonValue)+"\n" {
					t.Error("Response body does not match")
				}
			} else if rr.Code != 400 {
				t.Error("Expected response code to be 400")
			}
		})
	}
}

func newRequestRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, strPath, fnHandler)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
