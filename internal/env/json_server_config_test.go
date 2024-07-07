/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * json_server_config_test.go
 * $Id$
 */

package env

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	c0env "github.com/caarlos0/env"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
	"github.com/stretchr/testify/assert"
)

func TestGetServerConfig(t *testing.T) {
	flag = new(flags)
	flag.restore = new(bool)
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
    "restore": false,
    "store_interval": "1s",
    "store_file": "/path/to/file.db",
    "database_dsn": "",
    "crypto_key": "/path/to/key.pem"
}`,
			want: `
	dataBaseDSN     : 
	fileStoragePath : /path/to/file.db
	key             : <nil>
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : false
	serverAddress   : localhost:8080
	storeInterval   : 1s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test configFileName #2",
			fCfg: `{
  "address": "localhost:8082",
  "restore": false,
  "store_interval": "13s",
  "store_file": "/path/to/file.db",
  "database_dsn": "",
  "crypto_key": "/path/to/key.pem"
}`,
			want: `
	dataBaseDSN     : 
	fileStoragePath : /path/to/file.db
	key             : <nil>
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : false
	serverAddress   : localhost:8082
	storeInterval   : 13s
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
			agentConfig := getServerConfig()
			initServerConfig()
			assert.Nil(t, err)
			assert.Equal(t, "/path/to/key.pem", agentConfig.CryptoKey)
			got := cfg.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestServerConfigEasyJSON(t *testing.T) {
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

			c := new(serverConfig)

			body, _ := json.Marshal(test.input)
			_ = c.UnmarshalJSON(body)
			_, _ = c.MarshalJSON()
			w := new(jwriter.Writer)
			c.MarshalEasyJSON(w)
			_ = easyjson.UnmarshalFromReader(bytes.NewReader(body), c)
		})
	}
}
