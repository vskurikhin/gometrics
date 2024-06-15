/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
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
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/analysis/facts/deprecated"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"regexp"
)

// Main статический анализ кода.
// Механизм запуска multichecker: перейти в коммандной оболочке в каталог проекта и запустить multichecker ./...
// Используются следующие анализаторы:
//   - appends.Analyzer который определяет, есть ли в добавлении только одна переменная.
//   - asmdecl.Analyzer ...
//   - atomic.Analyzer проверяющий распространенные ошибки с помощью пакета sync/atomic.
//   - assign.Analyzer обнаруживающий бесполезные назначения.
//   - bools.Analyzer ...
//   - buildtag.Analyzer ...
//   - cgocall.Analyzer ...
//   - composite.Analyzer ...
//   - copylock.Analyzer анализатор, который проверяет блокировки, ошибочно переданные по значению.
//   - ctrlflow.Analyzer анализатор, который предоставляет синтаксический граф потока управления для тела функции.
//   - defers.Analyzer анализатор, проверяющий распространенные ошибки в операторах defer.
//   - directive.Analyzer ...
//   - errorsas.Analyzer ...
//   - framepointer.Analyzer ...
//   - httpresponse.Analyzer ...
//   - ifaceassert.Analyzer ...
//   - inspect.Analyzer анализатор, который предоставляет синтаксические деревья пакета.
//   - loopclosure.Analyzer проверяющий наличие ссылок на включающие переменные цикла внутри вложенных функций.
//   - lostcancel.Analyzer ...
//   - nilfunc.Analyzer ...
//   - printf.Analyzer анализатор, проверяющий согласованность строк и аргументов формата Printf.
//   - shadow.Analyzer анализатор проверяет затененные переменные.
//   - shift.Analyzer ...
//   - sigchanyzer.Analyzer ...
//   - stdmethods.Analyzer ...
//   - stringintconv.Analyzer ...
//   - structtag.Analyzer анализатор, проверяющий правильность формирования тегов полей структуры.
//   - tests.Analyzer ...
//   - testinggoroutine.Analyzer ...
//   - timeformat.Analyzer ...
//   - unmarshal.Analyzer ...
//   - unreachable.Analyzer проверяет недоступный код.
//   - unsafeptr.Analyzer ...
//   - unusedresult.Analyzer ...
func Main() {

	// определяем map подключаемых правил
	var mychecks []*analysis.Analyzer
	mychecks = append(
		mychecks,
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		testinggoroutine.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
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
