package mtls_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp/client/mtls"
)

func TestNewCredentials(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ca, err := os.ReadFile("../../../../certs/ca.crt")
	require.NoError(err)
	clientCert, err := os.ReadFile("../../../../certs/client.crt")
	require.NoError(err)
	clientKey, err := os.ReadFile("../../../../certs/client.key")
	require.NoError(err)

	creds, err := mtls.NewCredentials(mtls.Config{
		CA:   ca,
		Cert: clientCert,
		Key:  clientKey,
	})

	assert.NoError(err)
	assert.NotNil(creds)
}

func TestNewCredentialsFailing(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ca, err := os.ReadFile("../../../../certs/ca.crt")
	require.NoError(err)
	clientCert, err := os.ReadFile("../../../../certs/client.crt")
	require.NoError(err)

	creds, err := mtls.NewCredentials(mtls.Config{
		CA:   ca,
		Cert: clientCert,
		Key:  nil,
	})

	assert.Nil(creds)
	assert.EqualError(err, "failed to load client cert/key pair: tls: failed to find any PEM data in key input")
}
