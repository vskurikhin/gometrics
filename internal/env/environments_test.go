/*
 * This file was last modified at 2024-06-10 18:40 by Victor N. Skurikhin.
 * environments_test.go
 * $Id$
 */

package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnvironments(t *testing.T) {

	getEnvironments()
	env.Address = []string{"localhost", "8080"}
	assert.Equal(t, "localhost:8080", env.parseEnvAddress())
}

func TestEnvironmentsNegative(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			if r.(string) == "strconv.Atoi: parsing \"\": invalid syntax" {
				t.Log("Test passed as expected")
			} else {
				t.Fatal(r)
			}
		}
	}()

	getEnvironments()
	env.Address = []string{"localhost", ""}
}
