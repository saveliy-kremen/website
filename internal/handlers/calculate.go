package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"website/internal/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
)

type Data struct {
	A *int `json:"a"`
	B *int `json:"b"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data Data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		sendInputError(w)
		return
	}

	a := uint(*data.A)
	b := uint(*data.B)
	chA := make(chan uint)
	chB := make(chan uint)

	go utils.FactorialNumber(a, chA)
	go utils.FactorialNumber(b, chB)
	factA := <-chA
	factB := <-chB
	fmt.Printf("Factorial of a is: %d and b is: %d\n", factA, factB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func MiddlewareCalculate(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var data Data
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		err := json.Unmarshal(bodyBytes, &data)
		if err != nil {
			sendInputError(w)
			return
		}
		spew.Dump(data)
		if data.A == nil || *data.A < 0 || data.B == nil || *data.B < 0 {
			sendInputError(w)
			return
		}

		h(w, r, params)
	}
}

func sendInputError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponse{Error: "Incorrect input"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
