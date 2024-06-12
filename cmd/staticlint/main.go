/*
 * This file was last modified at 2024-06-10 11:40 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/vskurikhin/gometrics/internal/multichecker"
	"os"
)

var _ = func() int {
	os.Args = []string{"multichecker", "-test=false", "./..."}
	return 0
}()

func main() {
	multichecker.Main()
}
