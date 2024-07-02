/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * parameters_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestParameters(t *testing.T) {

	onceEnv = new(sync.Once)
	env = nil
	onceParameters = new(sync.Once)
	parameters = nil

	var tests = []struct {
		name string
		fCfg func() Parameters
		want string
	}{
		{
			name: "Test parameters #1",
			fCfg: GetParameters,
			want: `
	OutboundIP : <nil>
	PrivateKey : <nil>
	PublicKey  : <nil>
`,
		},
		{
			name: "Test parameters #2",
			fCfg: getTestAgentParameters,
			want: `
	OutboundIP : 127.0.0.1
	PrivateKey : <nil>
	PublicKey  : &{102720181439843039897263195445772230954649205681423018572611629656486843115981 65537}
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			onceCfg = new(sync.Once)
			c := test.fCfg()
			got := c.String()
			assert.Equal(t, test.want, got)
		})
	}
	p := getTestServerParameters()
	got := p.String()
	assert.NotNil(t, got)
}

func getTestAgentParameters() Parameters {
	parameters = nil
	onceParameters = new(sync.Once)
	getTestConfigAgent()
	loadPublicKey()
	return GetParameters()
}

func getTestServerParameters() Parameters {
	parameters = nil
	onceParameters = new(sync.Once)
	getTestConfigServer()
	loadPrivateKey()
	return GetParameters()
}
