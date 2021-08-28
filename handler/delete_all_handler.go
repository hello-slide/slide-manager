package handler

import (
	"context"
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	_storage "github.com/hello-slide/slide-manager/storage"
	"github.com/hello-slide/slide-manager/utils"
)

func DeleteAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userId, err := utils.GetSessonToken(ctx, client, w, r, tokenManagerName, url, "/slide/deleteall")
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	if len(userId) == 0 {
		return
	}

	slideManager := slide.NewSlideManager(ctx, &client, userId)
	storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	storageOp := _storage.NewStorageOp(ctx, *storageClient, "page-data")
	if err := slideManager.DeleteAll(*storageOp); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
}
