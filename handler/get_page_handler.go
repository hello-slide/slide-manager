package handler

import (
	"context"
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	_storage "github.com/hello-slide/slide-manager/storage"
	"github.com/hello-slide/slide-manager/token"
)

func GetPageHandler(w http.ResponseWriter, r *http.Request) {
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
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "page-data")
	data, err := slideManager.GetPage(slideId, pageId, *storageOp)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write(data)
}
