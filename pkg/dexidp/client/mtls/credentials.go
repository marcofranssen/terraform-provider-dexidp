package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"google.golang.org/grpc/credentials"
)

type Config struct {
	CA   []byte
	Cert []byte
	Key  []byte
}

func NewCredentials(config Config) (credentials.TransportCredentials, error) {
	cPool := x509.NewCertPool()

	if !cPool.AppendCertsFromPEM(config.CA) {
		return nil, fmt.Errorf("failed to parse CA crt")
	}

	certificate, err := tls.X509KeyPair(config.Cert, config.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key pair: %w", err)
	}

	clientTLSConfig := &tls.Config{
		RootCAs:      cPool,
		Certificates: []tls.Certificate{certificate},
	}

	return credentials.NewTLS(clientTLSConfig), nil
}
