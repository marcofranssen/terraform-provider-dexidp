package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/marcofranssen/terraform-provider-dexidp/pkg/dexidp/client"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
		host string
		err  error
	}{
		{
			name: "unspecified host",
			host: "",
			err:  fmt.Errorf("dial: %w", errors.New("failed to build resolver: passthrough: received empty target in Build()")),
		},
		{
			name: "localhost with port",
			host: "localhost:5557",
			err:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			assert := assert.New(tt)

			c, err := client.New(tc.host)
			if tc.err != nil {
				assert.Nil(c)
				assert.ErrorContains(err, tc.err.Error())
			} else {
				assert.NotNil(c)
				assert.NoError(err)
			}
		})
	}
}
