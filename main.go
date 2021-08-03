package main

import (
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/slide/", rootHandler)

	handler := networkUtils.CorsConfig.Handler(mux)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		panic(err)
	}
}
