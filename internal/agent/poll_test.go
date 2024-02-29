/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * poll_test.go
 * $Id$
 */

package agent

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"github.com/vskurikhin/gometrics/internal/types"
	"runtime"
	"testing"
)

func TestPoll(t *testing.T) {

	memStats := new(runtime.MemStats)
	memStorage := memory.Instance()

	var tests = []struct {
		name  string
		input []types.Name
		want  string
	}{
		{name: "positive test #0", input: []types.Name{types.Alloc}, want: "Alloc"},
		{name: "positive test #1", input: []types.Name{types.GCCPUFraction}, want: "GCCPUFraction"},
		{name: "positive test #2", input: []types.Name{types.NumForcedGC}, want: "NumForcedGC"},
		{name: "positive test #3", input: []types.Name{types.PollCount}, want: "PollCount"},
		{name: "positive test #4", input: []types.Name{types.RandomValue}, want: "RandomValue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			poll(test.input, memStats)
			got := memStorage.Get(test.want)
			assert.NotNil(t, got)
		})
	}
}
