/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * multichecker.go
 * $Id$
 */

// Package multichecker Модуль Статический анализ кода
package multichecker

import (
	"regexp"

	"github.com/kisielk/errcheck/errcheck"
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

	"github.com/vskurikhin/gometrics/internal/analyzer"
)

// Main статический анализ кода.
// Механизм запуска multichecker: перейти в коммандной оболочке в каталог проекта и запустить multichecker ./...
// Используются следующие анализаторы:
//   - appends.Analyzer который определяет, есть ли в добавлении только одна переменная.
//   - asmdecl.Analyzer сообщает о несоответствиях между файлами сборки и объявлениями Go.
//   - atomic.Analyzer проверяющий распространенные ошибки с помощью пакета sync/atomic.
//   - assign.Analyzer обнаруживающий бесполезные назначения.
//   - bools.Analyzer обнаруживает распространенные ошибки, связанные с логическими операторами.
//   - buildtag.Analyzer проверяет теги сборки.
//   - cgocall.Analyzer обнаруживает некоторые нарушения правил передачи указателей cgo.
//   - composite.Analyzer проверяет наличие неключевых составных литералов.
//   - copylock.Analyzer анализатор, который проверяет блокировки, ошибочно переданные по значению.
//   - ctrlflow.Analyzer анализатор, который предоставляет синтаксический граф потока управления для тела функции.
//   - defers.Analyzer анализатор, проверяющий распространенные ошибки в операторах defer.
//   - directive.Analyzer проверяет известные директивы toolchain Go.
//   - errorsas.Analyzer проверяет, что второй аргумент error является указателем на ошибку реализации типа.
//   - framepointer.Analyzer проверяет если ассемблерный код затирает указатель кадра перед его сохранением.
//   - httpresponse.Analyzer помогает обнаружить скрытые ошибки разыменования нуля, сообщая о диагностика подобных ошибок.
//   - ifaceassert.Analyzer помечает неверные утверждения типа интерфейса.
//   - inspect.Analyzer анализатор, который предоставляет синтаксические деревья пакета.
//   - loopclosure.Analyzer проверяющий наличие ссылок на включающие переменные цикла внутри вложенных функций.
//   - lostcancel.Analyzer сообщает об ошибке вызова функции cancel(), возвращаемой context.WithCancel, либо, что переменной был присвоен пустой идентификатор, либо, что существует путь control-flow от вызова до оператора возврата и этот путь не «используется».
//   - nilfunc.Analyzer проверяет бесполезные сравнения с нулем.
//   - printf.Analyzer анализатор, проверяющий согласованность строк и аргументов формата Printf.
//   - shadow.Analyzer анализатор проверяет затененные переменные.
//   - shift.Analyzer проверяет сдвиги, превышающие ширину целого числа.
//   - sigchanyzer.Analyzer обнаруживает неправильное использование небуферизованного сигнала в качестве аргумента signal.Notify.
//   - stdmethods.Analyzer проверяет наличие орфографических ошибок в сигнатурах методов, аналогичных известным интерфейсам.
//   - stringintconv.Analyzer проверяет флаги конвертации целых чисел в строки.
//   - structtag.Analyzer анализатор, проверяющий правильность формирования тегов полей структуры.
//   - timeformat.Analyzer проверяет использование time.Format или time.Parse вызывает неправильный формат.
//   - unmarshal.Analyzer проверяет передачу типов, не являющихся указателями или неинтерфейсными, для функций демаршалинга и декодирования.
//   - unreachable.Analyzer проверяет недоступный код.
//   - unsafeptr.Analyzer проверяет недопустимые преобразования uintptr в unsafe.Pointer.
//   - unusedresult.Analyzer проверяет наличие неиспользуемых результатов вызовов определенных чистых функций.
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
