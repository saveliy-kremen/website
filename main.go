package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Data struct {
	A *int `json:"a"`
	B *int `json:"b"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	router := httprouter.New()

	router.POST("/calculate", middleware(calculate))
	port := ":8989"
	log.Printf("Server started at port %v", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var data Data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		sendError(w)
		return
	}

	a := uint(*data.A)
	b := uint(*data.B)
	chA := make(chan uint)
	chB := make(chan uint)

	go factorialNumber(a, chA)
	go factorialNumber(b, chB)
	factA := <-chA
	factB := <-chB
	fmt.Printf("Factorial of a is: %d and b is: %d\n", factA, factB)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func middleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var data Data
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		err := json.Unmarshal(bodyBytes, &data)
		if err != nil {
			sendError(w)
			return
		}
		if data.A == nil || *data.A < 0 || data.B == nil || *data.B < 0 {
			sendError(w)
			return
		}

		h(w, r, params)
	}
}

func sendError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponse{Error: "Incorrect input"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

func factorialNumber(n uint, ch chan uint) {
	if n == 0 {
		ch <- 1
		return
	}
	fact := n
	for i := n - 1; i > 0; i-- {
		fact *= i
	}
	ch <- fact
}
