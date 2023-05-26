package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"

	"google.golang.org/grpc/credentials"
)

type Config struct {
	CA   io.Reader
	Cert io.Reader
	Key  io.Reader
}

func NewCredentials(config Config) (credentials.TransportCredentials, error) {
	cPool := x509.NewCertPool()

	ca, err := io.ReadAll(config.CA)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA crt: %w", err)
	}

	if !cPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to parse CA crt")
	}

	cert, err := io.ReadAll(config.Cert)
	if err != nil {
		return nil, fmt.Errorf("failed to read client crt: %w", err)
	}

	key, err := io.ReadAll(config.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %w", err)
	}

	certificate, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key pair: %w", err)
	}

	clientTLSConfig := &tls.Config{
		RootCAs:      cPool,
		Certificates: []tls.Certificate{certificate},
	}

	return credentials.NewTLS(clientTLSConfig), nil
}
