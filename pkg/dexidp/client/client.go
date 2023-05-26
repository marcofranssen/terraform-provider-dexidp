package client

import (
	"fmt"

	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// New instantiates a new api.DexClient
func New(host string, creds credentials.TransportCredentials) (api.DexClient, error) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return api.NewDexClient(conn), nil
}
