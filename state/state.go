package state

import (
	"context"

	"github.com/dapr/go-sdk/client"
)

type state struct {
	client *client.Client
	ctx    *context.Context
	store  string
}

func NewState(client *client.Client, ctx *context.Context, store string) *state {
	return &state{
		client: client,
		ctx:    ctx,
		store:  store,
	}
}

func (s *state) Get(key string) (*client.StateItem, error) {
	return (*s.client).GetState(*s.ctx, s.store, key)
}

func (s *state) Set(key string, value []byte) error {
	if err := (*s.client).SaveState(*s.ctx, s.store, key, value); err != nil {
		return err
	}

	return nil
}

func (s *state) Delete(key string) error {
	return (*s.client).DeleteState(*s.ctx, s.store, key)
}
