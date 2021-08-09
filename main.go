package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"cloud.google.com/go/storage"
	dapr "github.com/dapr/go-sdk/client"
	networkUtils "github.com/hello-slide/network-util"
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
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	title, err := networkUtils.PickValue("Title", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}

	slideManager := slide.NewSlideManager(ctx, &client, userId)
	slideId, err := slideManager.Create(title)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(map[string]string{
		"slide_id": slideId,
	})
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)
}

func createPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	pageType, err := networkUtils.PickValue("PageType", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}

	slideManager := slide.NewSlideManager(ctx, &client, userId)
	pageData, err := slideManager.CreatePage(slideId, pageType)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(pageData)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
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
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	slideConfig, err := slideManager.GetInfo()
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(slideConfig)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)

}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	slideDetails, err := slideManager.GetSlideDetails(slideId)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(slideDetails)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)
}

func renameHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	newName, err := networkUtils.PickValue("newName", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	if err := slideManager.Rename(slideId, newName); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
}

func setPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	pageId, err := networkUtils.PickValue("PageID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	data, err := networkUtils.PickValue("Data", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}

	slideManager := slide.NewSlideManager(ctx, &client, userId)
	storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "SlideData")
	if err := slideManager.SetPage([]byte(data), slideId, pageId, *storageOp); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(map[string]string{
		"slide_id": slideId,
	})
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(tokenJson)
}

func getPageHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	pageId, err := networkUtils.PickValue("PageID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}

	slideManager := slide.NewSlideManager(ctx, &client, userId)
	storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "SlideData")
	data, err := slideManager.GetPage(slideId, pageId, *storageOp)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	tokenJson, err := json.Marshal(data)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(tokenJson)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	slideId, err := networkUtils.PickValue("SlideID", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "SlideData")
	if err := slideManager.Delete(slideId, *storageOp); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
}

func deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	sessionToken, err := networkUtils.PickValue("SessionToken", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := token.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "SlideData")
	if err := slideManager.DeleteAll(*storageOp); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
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
	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/slide/create", createHandler)
	mux.HandleFunc("/slide/createpage", createPageHandler)

	mux.HandleFunc("/slide/list", listHandler)
	mux.HandleFunc("/slide/details", detailsHandler)
	mux.HandleFunc("/slide/rename", renameHandler)

	mux.HandleFunc("/slide/setpage", setPageHandler)
	mux.HandleFunc("/slide/getpage", getPageHandler)

	mux.HandleFunc("/slide/delete", deleteHandler)
	mux.HandleFunc("/slide/deleteall", deleteAllHandler)

	handler := networkUtils.CorsConfig.Handler(mux)

	if err := http.ListenAndServe(":3000", handler); err != nil {
		panic(err)
	}
}
