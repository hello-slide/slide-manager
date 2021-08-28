package handler

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	dapr "github.com/dapr/go-sdk/client"
	_storage "github.com/hello-slide/slide-manager/storage"
)

var client dapr.Client
var storageClient *storage.Client
var tokenManagerName string = os.Getenv("TOKEN_MANAGER")

// Initialize dapr client.
func InitClient() error {
	_client, err := dapr.NewClient()
	if err != nil {
		return err
	}
	client = _client

	return nil
}

// Initialize storage client.
func InitStorage(ctx context.Context) error {
	_storageClient, err := _storage.CreateClient(ctx)
	if err != nil {
		return err
	}

	storageClient = _storageClient
	return nil
}
