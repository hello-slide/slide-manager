package handler

import (
	"context"
	"encoding/json"
	"net/http"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	"github.com/hello-slide/slide-manager/utils"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	headerData, err := networkUtils.GetHeader(w, r)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	title, err := networkUtils.PickValue("Title", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := utils.GetSessonToken(ctx, client, w, r, tokenManagerName, url, "/slide/create")
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	if len(userId) == 0 {
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
