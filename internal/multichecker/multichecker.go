/*
 * This file was last modified at 2024-06-10 12:50 by Victor N. Skurikhin.
 * multichecker.go
 * $Id$
 */

// Package multichecker Модуль Статический анализ кода
package multichecker

import (
	"github.com/kisielk/errcheck/errcheck"
	"github.com/vskurikhin/gometrics/internal/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"honnef.co/go/tools/analysis/facts/deprecated"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"regexp"
)

// Main статический анализ кода.
// Механизм запуска multichecker: перейти в коммандной оболочке в каталог проекта и запустить multichecker.
// Используются следующие анализаторы:
//   - appends.Analyzer который определяет, есть ли в добавлении только одна переменная.
//   - assign.Analyzer обнаруживающий бесполезные назначения.
//   - atomic.Analyzer проверяющий распространенные ошибки с помощью пакета sync/atomic.
//   - copylock.Analyzer анализатор, который проверяет блокировки, ошибочно переданные по значению.
//   - ctrlflow.Analyzer анализатор, который предоставляет синтаксический граф потока управления для тела функции.
//   - defers.Analyzer анализатор, проверяющий распространенные ошибки в операторах defer.
//   - inspect.Analyzer анализатор, который предоставляет синтаксические деревья пакета.
//   - loopclosure.Analyzer проверяющий наличие ссылок на включающие переменные цикла внутри вложенных функций.
//   - printf.Analyzer анализатор, проверяющий согласованность строк и аргументов формата Printf.
//   - shadow.Analyzer анализатор проверяет затененные переменные.
//   - structtag.Analyzer анализатор, проверяющий правильность формирования тегов полей структуры.
//   - unreachable.Analyzer проверяет недоступный код.
func Main() {
	// определяем map подключаемых правил
	var mychecks []*analysis.Analyzer
	mychecks = append(
		mychecks,
		appends.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		defers.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		unreachable.Analyzer,
	)
	re, _ := regexp.Compile(`^SA\d+$`)
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		if re.MatchString(v.Analyzer.Name) {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, deprecated.Analyzer)
	mychecks = append(mychecks, errcheck.Analyzer)
	mychecks = append(mychecks, &analysis.Analyzer{
		Name: "check_package_comment",
		Doc:  "Incorrect or missing package comment",
		Run:  stylecheck.CheckPackageComment,
	})
	mychecks = append(mychecks, &analysis.Analyzer{
		Name: "check_os_exit",
		Doc:  "Check for call os.Exit error",
		Run:  analyzer.CheckCallOsExit,
	})
	multichecker.Main(mychecks...)
}
