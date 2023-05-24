package dexidp

import (
	"fmt"

	"github.com/dexidp/dex/api/v2"
	"google.golang.org/grpc"
)

func newClient(host string) (api.DexClient, error) {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return api.NewDexClient(conn), nil
}
