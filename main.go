package main

import (
	"net/http"
	"os"

	networkUtils "github.com/hello-slide/network-util"
)

var sample string

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(sample))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	sample = os.Getenv("KEY")

	mux := http.NewServeMux()
	mux.HandleFunc("/slide/", rootHandler)
	mux.HandleFunc("/slide/create", createHandler)
	mux.HandleFunc("/slide/list", listHandler)
	mux.HandleFunc("/slide/edit", editHandler)
	mux.HandleFunc("/slide/delete", deleteHandler)

	handler := networkUtils.CorsConfig.Handler(mux)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		panic(err)
	}
}
