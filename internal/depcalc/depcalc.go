package depcalc

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/vimcki/go-di-graph/internal/deptree"
	"github.com/vimcki/go-di-graph/internal/deptree/evaluator"
	"github.com/vimcki/go-di-graph/internal/globals"
)

type BuildInfo struct {
	Imports map[string]string
	Config  Config
	Result  Result
}

type Config struct {
	Package string
	Struct  string
}

type Result struct {
	Type string
}

func Depcalc(entryPoint, path string) (string, error) {
	packages, err := parsePackages(path)
	if err != nil {
		return "", fmt.Errorf("failed to parse packages: %w", err)
	}

	var start *ast.FuncDecl

	c, err := globals.Get(packages)
	if err != nil {
		return "", fmt.Errorf("failed to get globals: %w", err)
	}

	for _, pkg := range packages {
		fnMap := make(map[string]*ast.FuncDecl)
		allImorts := make(map[string]map[string]string)
		for _, file := range pkg.Files {
			importsMap := make(map[string]string)
			ast.Inspect(file, func(node ast.Node) bool {
				switch t := node.(type) {
				case *ast.ImportSpec:
					var name string
					if t.Name == nil {
						split := strings.Split(t.Path.Value, "/")
						name = split[len(split)-1]
						name = strings.Trim(name, "\"")
					} else {
						name = t.Name.Name
					}

					importsMap[name] = strings.Trim(t.Path.Value, "\"")
				}
				return true
			})
			ast.Inspect(file, func(node ast.Node) bool {
				switch t := node.(type) {
				case *ast.FuncDecl:
					var fnIdent string
					fnIdent, err = functionIdentifier(t)
					if err != nil {
						return false
					}
					allImorts[t.Name.Name] = importsMap
					if fnIdent != entryPoint {
						fnMap[fnIdent] = t
						return true
					}
					start = t
				}
				return true
			})
			if err != nil {
				return "", fmt.Errorf("failed to inspect file: %w", err)
			}
		}

		e := evaluator.NewEvaluator(fnMap, c, allImorts)
		dep, err := e.Eval(start)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate dependencies: %w", err)
		}
		printDep(dep, 0)
		bytes, err := json.Marshal(dep)
		if err != nil {
			return "", fmt.Errorf("failed to marshal dependencies: %w", err)
		}
		return string(bytes), nil
	}
	return "", errors.New("no packages found")
}

// parsePackages parses all packages in the specified directory.
func parsePackages(dirPath string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, 0)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

func printDep(dep deptree.Dependency, indent int) {
	for _, dep := range dep.Deps {
		printDep(dep, indent+1)
	}
}

func functionIdentifier(fn *ast.FuncDecl) (string, error) {
	var receiver string
	if fn.Recv != nil {
		switch fn.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			receiver = fn.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		case *ast.Ident:
			receiver = fn.Recv.List[0].Type.(*ast.Ident).Name
		default:
			return "", fmt.Errorf("unknown receiver type: %T", fn.Recv.List[0].Type)
		}
	}
	if receiver == "" {
		return fn.Name.Name, nil
	}
	return fmt.Sprintf("%s.%s", receiver, fn.Name.Name), nil
}
