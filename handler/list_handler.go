package handler

import (
	"context"
	"encoding/json"
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	"github.com/hello-slide/slide-manager/utils"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userId, err := utils.GetSessonToken(ctx, client, w, r, tokenManagerName, url, "/slide/list")
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	if len(userId) == 0 {
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
