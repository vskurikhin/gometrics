/*
 * This file was last modified at 2024-06-10 11:52 by Victor N. Skurikhin.
 * multichecker_test.go
 * $Id$
 */

package multichecker

import (
	"os"
	"testing"
)

func TestUpdateMain(t *testing.T) {
	os.Args = []string{"multichecker", "-test=false", "./..."}
	defer func() {
		if r := recover(); r != nil {
			if r.(string) == "unexpected call to os.Exit(0) during test" {
				t.Log("Test passed as expected")
			} else {
				t.Fatal(r)
			}
		}
	}()
	Main()
}
