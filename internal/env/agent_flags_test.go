/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * agent_flags_test.go
 * $Id$
 */

package env

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgentFlagsURLHost(t *testing.T) {
	urlHost := "localhost:8080"
	reportInterval := time.Duration(10)
	pollInterval := time.Duration(2)
	var tests = []struct {
		name  string
		input agentFlags
		want  string
		want1 time.Duration
		want2 time.Duration
	}{
		{name: "Test URLHost() positive #0",
			input: agentFlags{
				serverFlags:    serverFlags{serverAddress: &urlHost},
				reportInterval: &reportInterval,
				pollInterval:   &pollInterval,
			},
			want:  "http://localhost:8080",
			want1: time.Duration(10),
			want2: time.Duration(2),
		},
		{name: "Test URLHost() positive #1",
			input: agentFlags{
				serverFlags:    serverFlags{serverAddress: &urlHost},
				urlHost:        &urlHost,
				reportInterval: &reportInterval,
				pollInterval:   &pollInterval,
			},
			want:  "localhost:8080",
			want1: time.Duration(10),
			want2: time.Duration(2),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.URLHost()
			fmt.Printf("got: %v\n", *got)
			assert.Equal(t, test.want, *got)
			got1 := test.input.ReportInterval()
			assert.Equal(t, test.want1, got1)
			got2 := test.input.PollInterval()
			assert.Equal(t, test.want2, got2)
		})
	}
}
