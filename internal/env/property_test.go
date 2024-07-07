/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * property_test.go
 * $Id$
 */

package env

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProperty(t *testing.T) {

	onceEnv = new(sync.Once)
	env = nil
	onceProperty = new(sync.Once)

	var tests = []struct {
		name string
		fCfg func() Property
		want string
	}{
		{
			name: "Test property #1",
			fCfg: GetProperty,
			want: `
	OutboundIP : 127.0.0.1
	PrivateKey : <nil>
	PublicKey  : <nil>
    Storage    : &{{{0 0} 0 0 {{} 0} {{} 0}} map[] 0}
`,
		},
		{
			name: "Test property #2",
			fCfg: getTestAgentProperty,
			want: `
	OutboundIP : 127.0.0.1
	PrivateKey : <nil>
	PublicKey  : &{102720181439843039897263195445772230954649205681423018572611629656486843115981 65537}
    Storage    : &{{{0 0} 0 0 {{} 0} {{} 0}} map[] 0}
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
	p := getTestServerProperty()
	got := p.String()
	assert.NotNil(t, got)
}

func getTestAgentProperty() Property {
	onceProperty = new(sync.Once)
	getTestConfigAgent()
	LoadPublicKey()
	return GetProperty()
}

func getTestServerProperty() Property {
	onceProperty = new(sync.Once)
	getTestConfigServer()
	LoadPrivateKey()
	return GetProperty()
}
