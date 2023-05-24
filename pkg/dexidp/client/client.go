package client

import (
	"fmt"

	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
)

// New instantiates a new api.DexClient
func New(host string) (api.DexClient, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return api.NewDexClient(conn), nil
}
