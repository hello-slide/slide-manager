package main

import (
	"context"
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/handler"
)

func init() {
	ctx := context.Background()

	if err := handler.InitClient(); err != nil {
		panic(err)
	}
	if err := handler.InitStorage(ctx); err != nil {
		panic(err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.RootHandler)

	mux.HandleFunc("/slide/create", handler.CreateHandler)
	mux.HandleFunc("/slide/createpage", handler.CreatePageHandler)

	mux.HandleFunc("/slide/list", handler.ListHandler)
	mux.HandleFunc("/slide/details", handler.DetailsHandler)
	mux.HandleFunc("/slide/rename", handler.RenameHandler)
	mux.HandleFunc("/slide/swap", handler.SwapHandler)

	mux.HandleFunc("/slide/setpage", handler.SetPageHandler)
	mux.HandleFunc("/slide/getpage", handler.GetPageHandler)

	mux.HandleFunc("/slide/delete", handler.DeleteSlideHandler)
	mux.HandleFunc("/slide/deleteall", handler.DeleteAllHandler)
	mux.HandleFunc("/slide/deletepage", handler.DeletePageHandler)

	handler := networkUtils.CorsConfig.Handler(mux)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		panic(err)
	}
}
