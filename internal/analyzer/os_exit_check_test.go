/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * os_exit_check_test.go
 * $Id$
 */

package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for call os.Exit errors",
	Run:  CheckCallOsExit,
}

func TestMyAnalyzer(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор OsExitCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	analysistest.Run(t, analysistest.TestData(), OsExitCheckAnalyzer, "./...")
}
