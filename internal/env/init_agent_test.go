/*
 * This file was last modified at 2024-02-11 17:08 by Victor N. Skurikhin.
 * init_agent_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitAgentWithEnv(t *testing.T) {
	old := cfg
	urlHost := "localhost:8080"
	arrHost := []string{"localhost", "8080"}
	var tests = []struct {
		name  string
		input config
		want  agentEnv
	}{
		{
			name:  "Test Parse() function for PollCount with type: parser",
			input: config{Address: arrHost, ReportInterval: 10, PollInterval: 2},
			want: agentEnv{
				serverEnv:      serverEnv{serverAddress: &urlHost},
				reportInterval: time.Duration(10),
				pollInterval:   time.Duration(2),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg = test.input
			InitAgent()
			//fmt.Printf("urlHost: %s\n", *Agent.urlHost)
			assert.Equal(t, test.want, Agent)
		})
	}
	cfg = old
}
