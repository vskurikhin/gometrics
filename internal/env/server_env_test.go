/*
 * This file was last modified at 2024-02-11 17:36 by Victor N. Skurikhin.
 * server_env_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerEnvServerAddress(t *testing.T) {
	urlHost := "localhost:8080"
	var tests = []struct {
		name  string
		input serverEnv
		want  string
	}{
		{name: "Test URLHost() positive #0",
			input: serverEnv{
				serverAddress: &urlHost,
			},
			want: "localhost:8080",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.ServerAddress()
			assert.Equal(t, test.want, got)
		})
	}
}
