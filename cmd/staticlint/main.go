/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/vskurikhin/gometrics/internal/multichecker"
	"os"
	"regexp"
)

var _ = func() int {
	if osArgsNotExistTest(os.Args) {
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

func osArgsNotExistTest(args []string) bool {
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
