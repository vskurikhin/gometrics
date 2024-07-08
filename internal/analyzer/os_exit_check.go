/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * os_exit_check.go
 * $Id$
 */

// Package analyzer анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const MAIN = "main"

func CheckCallOsExit(pass *analysis.Pass) (interface{}, error) {
	if !isPkg(pass.Pkg.Name()) {
		return nil, nil
	}
	checkCall := func(x ast.Expr) {
		if isOs(x) && isExit(x) {
			pass.Reportf(x.Pos(), "call os.Exit error")
		}
	}
	for _, file := range pass.Files {
		// функцией ast.Inspect проходим по всем узлам AST
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.CallExpr:
				checkCall(x.Fun)
			case *ast.FuncDecl:
				return x.Name.Name == MAIN
			}
			return true
		})
	}
	return nil, nil
}

func isOs(call ast.Expr) bool {
	if x, ok1 := call.(*ast.SelectorExpr); ok1 {
		if ident, ok2 := x.X.(*ast.Ident); ok2 {
			return ident.Name == "os"
		}
	}
	return false
}

func isExit(call ast.Expr) bool {
	if x, ok := call.(*ast.SelectorExpr); ok {
		return x.Sel.Name == "Exit"
	}
	return false
}

func isPkg(name string) bool {
	return name == MAIN || name == "pkg1"
}
