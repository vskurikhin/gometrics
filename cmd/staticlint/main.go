package main

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
	"os"
	"regexp"
)

var _ = func() int {
	os.Args = []string{"multichecker", "-test=false", "./..."}
	return 0
}()

func main() {
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
