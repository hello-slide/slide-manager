package handler

import (
	"context"
	"net/http"
	"strconv"

	networkUtils "github.com/hello-slide/network-util"
	"github.com/hello-slide/slide-manager/slide"
	"github.com/hello-slide/slide-manager/utils"
)

func SwapHandler(w http.ResponseWriter, r *http.Request) {
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
	origin, err := networkUtils.PickValue("Origin", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	target, err := networkUtils.PickValue("Target", headerData, w)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	originInt, err := strconv.Atoi(origin)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
	targetInt, err := strconv.Atoi(target)
	if err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}

	userId, err := utils.VerifySessionToken(ctx, client, sessionToken, tokenManagerName)
	if err != nil {
		networkUtils.ErrorResponse(w, 2, err)
		return
	}
	slideManager := slide.NewSlideManager(ctx, &client, userId)
	if err := slideManager.SwapPage(slideId, originInt, targetInt); err != nil {
		networkUtils.ErrorResponse(w, 1, err)
		return
	}
}
