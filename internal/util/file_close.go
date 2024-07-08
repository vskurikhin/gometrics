/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * file_close.go
 * $Id$
 */

package util

import (
	"fmt"
	"os"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func FileClose(file *os.File) {
	if err := file.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error file close %v", err)
	}
}

func FileCloseAndLog(file *os.File) {
	if err := file.Close(); err != nil {
		zf := MakeZapFields()
		zf.Append("error", err)
		logger.Log.Error("in Close", zf.Slice()...)
	}
}
