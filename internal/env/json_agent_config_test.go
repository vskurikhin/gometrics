/*
 * This file was last modified at 2024-06-24 22:54 by Victor N. Skurikhin.
 * json_agent_config_test.go
 * $Id$
 */

package env

import (
	"bytes"
	"encoding/json"
	c0env "github.com/caarlos0/env"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetAgentConfig(t *testing.T) {
	flag = new(flags)
	flagFile := os.O_CREATE | os.O_RDWR | os.O_TRUNC
	var tests = []struct {
		name string
		fCfg string
		want string
	}{
		{
			name: "Test configFileName #1",
			fCfg: `{
    "address": "localhost:8080",
    "report_interval": "1s",
    "poll_interval": "1s",
    "crypto_key": "/path/to/key.pem"
}`,
			want: `
	dataBaseDSN     : 
	fileStoragePath : 
	key             : <nil>
	pollInterval    : 1s
	reportInterval  : 1s
	restore         : false
	serverAddress   : localhost:8080
	storeInterval   : 0s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test configFileName #2",
			fCfg: `{
  "address": "localhost:8082",
  "report_interval": "42s",
  "poll_interval": "13s",
  "crypto_key": "/path/to/key.pem"
}`,
			want: `
	dataBaseDSN     : 
	fileStoragePath : 
	key             : <nil>
	pollInterval    : 13s
	reportInterval  : 42s
	restore         : false
	serverAddress   : localhost:8082
	storeInterval   : 0s
	urlHost         : http://localhost:8082
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("CONFIG", testConfigFileName)
			file, err := os.OpenFile(testConfigFileName, flagFile, 0640)
			assert.Nil(t, err)
			n, err := file.Write([]byte(test.fCfg))
			assert.Nil(t, err)
			file.Close()
			assert.True(t, n > 0)
			cfg = new(config)
			env = new(environments)
			err = c0env.Parse(env)
			agentConfig := getAgentConfig()
			initAgentConfig()
			assert.Nil(t, err)
			assert.Equal(t, "/path/to/key.pem", agentConfig.CryptoKey)
			got := cfg.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestAgentConfigEasyJSON(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name: "positive test #1",
			input: map[string]interface{}{
				"address":         "localhost:8080",
				"report_interval": "1s",
				"poll_interval":   "1s",
				"crypto_key":      "crypto_key",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := new(agentConfig)

			body, _ := json.Marshal(test.input)
			c.UnmarshalJSON(body)
			c.MarshalJSON()
			w := new(jwriter.Writer)
			c.MarshalEasyJSON(w)
			easyjson.UnmarshalFromReader(bytes.NewReader(body), c)
		})
	}
}
