/*
 * This file was last modified at 2024-04-05 09:33 by Victor N. Skurikhin.
 * main_test.go
 * $Id$
 */

package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
