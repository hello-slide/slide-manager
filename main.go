package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	dapr "github.com/dapr/go-sdk/client"
	networkUtils "github.com/hello-slide/network-util"
	networkutil "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	_storage "github.com/hello-slide/slide-manager/storage"
	"github.com/hello-slide/slide-manager/token"
)

var client dapr.Client
var storageClient *storage.Client
var tokenManagerName string = os.Getenv("TOKEN_MANAGER")

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	title, err := networkUtils.PickValue("Title", headerData, w)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}

	slideManager := slide.NewSlideManager(ctx, *storageClient, &client, "slide-data", userId)
	slideId, err := slideManager.Create(title)
	if err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}

	tokenJson, err := json.Marshal(map[string]string{
		"slide_id": slideId,
	})
	if err != nil {
		networkutil.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, *storageClient, &client, "slide-data", userId)
	slideConfig, err := slideManager.GetInfo()
	if err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}

	tokenJson, err := json.Marshal(slideConfig)
	if err != nil {
		networkutil.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, *storageClient, &client, "slide-data", userId)
	if err := slideManager.Delete(slideId); err != nil {
		networkUtils.ErrorStatus(w)
		fmt.Fprintln(w, err)
		return
	}
}

func init() {
	ctx := context.Background()

	_client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	client = _client

	storageClient, err = _storage.CreateClient(ctx)
	if err != nil {
		panic(err)
	}
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
