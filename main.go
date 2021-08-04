package main

import (
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {

}

func listHandler(w http.ResponseWriter, r *http.Request) {

}

func editHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
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
