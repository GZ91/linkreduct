/*
Пакет main содержит анализаторы для проверки кода на соответствие стандартам и соглашениям,
а также специфический анализатор для обнаружения вызовов os.Exit в функции main основного пакета.

Импорт необходимых пакетов и библиотек.
*/
package main

import (
	"go/ast"
	"strings"

	"github.com/fatih/errwrap/errwrap"
	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/ast/inspector"
	"honnef.co/go/tools/staticcheck"
)

/*
Функция main инициализирует и запускает анализаторы кода.
Она создает срез mychecks, содержащий все анализаторы staticcheck.Analyzers, а также другие
стандартные анализаторы из пакетов golang.org/x/tools/go/analysis/passes/.
Затем добавляет свой собственный анализатор MyAnalyzer в срез mychecks.
Наконец, вызывает multichecker.Main(), передавая все анализаторы в качестве аргументов.
*/
func main() {
	var mychecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if strings.HasSuffix(v.Analyzer.Name, "SA") || strings.HasSuffix(v.Analyzer.Name, "QF1") {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	mychecks = append(mychecks, asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer)

	mychecks = append(mychecks, goone.Analyzer)
	mychecks = append(mychecks, errwrap.Analyzer)

	mychecks = append(mychecks, MyAnalyzer)

	multichecker.Main(
		mychecks...,
	)
}

/*
Функция run выполняет анализ кода, проверяя вызовы os.Exit в функции main.
Она использует анализатор inspect.Analyzer для обхода AST и проверки условий.

Функция содержит вложенные функции, которые определяют, является ли файл main пакетом,
является ли функция main функцией main, и проверяют, является ли выражение вызовом os.Exit в main функции.
Если условия выполняются, генерируется сообщение об ошибке.
*/
var MyAnalyzer = &analysis.Analyzer{
	Name:     "exitinmain",
	Doc:      "проверка на наличиеos.Exit",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

/*
Функция run выполняет анализ кода, проверяя вызовы os.Exit в функции main.
Она использует анализатор inspect.Analyzer для обхода AST и проверки условий.

Функция содержит вложенные функции, которые определяют, является ли файл main пакетом,
является ли функция main функцией main, и проверяют, является ли выражение вызовом os.Exit в main функции.
Если условия выполняются, генерируется сообщение об ошибке.
*/
func run(pass *analysis.Pass) (interface{}, error) {
	isMainPkg := func(x *ast.File) bool {
		return x.Name.Name == "main"
	}

	isMainFunc := func(x *ast.FuncDecl) bool {
		return x.Name.Name == "main"
	}

	isOsExit := func(x *ast.SelectorExpr, isMain bool) bool {
		if !isMain || x.X == nil {
			return false
		}
		ident, ok := x.X.(*ast.Ident)
		if !ok {
			return false
		}
		if ident.Name == "os" && x.Sel.Name == "Exit" {
			pass.Reportf(ident.NamePos, "os.Exit called in main func in main package")
			return true
		}
		return false
	}

	i := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.SelectorExpr)(nil),
	}
	mainInspecting := false
	i.Preorder(nodeFilter, func(n ast.Node) {
		switch x := n.(type) {
		case *ast.File:
			if !isMainPkg(x) {
				return
			}
		case *ast.FuncDecl:
			f := isMainFunc(x)
			if mainInspecting && !f {
				mainInspecting = false
				return
			}
			mainInspecting = f
		case *ast.SelectorExpr:
			if isOsExit(x, mainInspecting) {
				return
			}
		}
	})

	return nil, nil
}
