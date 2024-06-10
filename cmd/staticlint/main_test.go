/*
 * This file was last modified at 2024-06-10 09:39 by Victor N. Skurikhin.
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
