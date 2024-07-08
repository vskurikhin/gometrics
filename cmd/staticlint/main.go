/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"os"
	"regexp"

	"github.com/vskurikhin/gometrics/internal/multichecker"
)

var _ = func() int {
	if osArgsNotExistTest() {
		args := make([]string, 0)
		for i, arg := range os.Args {
			if i == 0 {
				args = append(args, arg)
				args = append(args, "-test=false")
			} else {
				args = append(args, arg)
			}
		}
		os.Args = args
	}
	return 0
}()

func osArgsNotExistTest() bool {
	re, _ := regexp.Compile(`^SA\d+$`)
	for _, arg := range os.Args {
		if re.MatchString(arg) {
			return false
		}
	}
	return true
}

func main() {
	multichecker.Main()
}
