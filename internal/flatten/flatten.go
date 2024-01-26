package flatten

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/maja42/goval"
	"github.com/vimcki/go-di-graph/internal/evaluator"
	"github.com/vimcki/go-di-graph/internal/globals"
)

func Flatten(basePath, buildPackage, flatPackage, entryPoint, configFilePath string) error {
	conf, err := getConfiguration(configFilePath)
	if err != nil {
		return fmt.Errorf("error getting configuration: %v", err)
	}
	// Parse all packages in the given directory.
	fset := token.NewFileSet()

	packages, err := parsePackages(filepath.Join(basePath, buildPackage), fset)
	if err != nil {
		return fmt.Errorf("error parsing packages: %v", err)
	}

	entrypoint, err := findEntrypoint(entryPoint, packages)
	if err != nil {
		return fmt.Errorf("error finding entrypoint: %v", err)
	}

	cfgPath, err := findPath(entrypoint.Type.Params.List[0])
	if err != nil {
		return fmt.Errorf("error finding path: %v", err)
	}

	c, err := globals.Get(packages)
	if err != nil {
		return fmt.Errorf("error getting globals: %v", err)
	}

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					flattener, err := NewFlattener(fn, conf, cfgPath, c)
					if err != nil {
						return fmt.Errorf("error creating flattener: %v", err)
					}

					err = flattener.Flatten()
					if err != nil {
						return fmt.Errorf("error flattening: %v", err)
					}
				}
			}
		}
	}

	for _, pkg := range packages {
		for filename, file := range pkg.Files {
			// Create the file in the output directory
			outputFile, err := os.Create(
				filepath.Join(basePath, flatPackage, filepath.Base(filename)),
			)
			if err != nil {
				return fmt.Errorf("error creating file: %v", err)
			}
			defer outputFile.Close()
			// Print the AST to the file
			err = format.Node(outputFile, fset, file)
			if err != nil {
				return fmt.Errorf("error printing file: %v", err)
			}
		}
	}

	return nil
}

func findPath(entrypoint *ast.Field) (string, error) {
	switch t := entrypoint.Type.(type) {
	case *ast.Ident:
		return t.Name, nil
	case *ast.SelectorExpr:
		pkg := t.X.(*ast.Ident).Name
		paramStruct := entrypoint.Type.(*ast.SelectorExpr).Sel.Name

		return pkg + "." + paramStruct, nil
	case *ast.MapType:
		// map type doesn't have a path
		return "", nil
	case *ast.StarExpr:
		// print expr
		return findPath(&ast.Field{Type: t.X})

	default:
		return "", fmt.Errorf("unknown entrypoint type: %T", entrypoint.Type)
	}
}

func getConfiguration(s string) (map[string]interface{}, error) {
	if strings.HasSuffix(s, ".toml") {
		k := koanf.New(".")
		if err := k.Load(file.Provider(s), toml.Parser()); err != nil {
			return nil, fmt.Errorf("error loading file: %v", err)
		}

		return k.All(), nil
	}

	if strings.HasSuffix(s, ".json") {
		var config map[string]interface{}

		bytes, err := os.ReadFile(s)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			return nil, fmt.Errorf("error parsing json: %v", err)
		}

		return config, nil
	}

	return nil, fmt.Errorf("unknown config file type")
}

func parsePackages(dirPath string, fset *token.FileSet) (map[string]*ast.Package, error) {
	pkgs, err := parser.ParseDir(fset, dirPath, nil, 0)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func findEntrypoint(
	entryPoint string,
	packages map[string]*ast.Package,
) (*ast.FuncDecl, error) {
	var receiver string
	split := strings.Split(entryPoint, ".")
	if len(split) == 2 {
		receiver = split[0]
		entryPoint = split[1]
	}
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					if fn.Name.Name == entryPoint {
						if receiver != "" && fn.Recv != nil {
							var foundReceiver string
							switch fn.Recv.List[0].Type.(type) {
							case *ast.StarExpr:
								foundReceiver = fn.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
							case *ast.Ident:
								foundReceiver = fn.Recv.List[0].Type.(*ast.Ident).Name
							default:
								return nil, fmt.Errorf("unknown receiver type: %T", fn.Recv.List[0].Type)
							}
							if foundReceiver == receiver {
								return fn, nil
							} else {
								continue
							}
						} else {
							return fn, nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("entrypoint %v not found in package %v", entryPoint, packages)
}

type Flattener struct {
	cfgPath  string
	ctx      map[string]interface{}
	err      error
	entry    *ast.FuncDecl
	selector string
}

func NewFlattener(
	entry *ast.FuncDecl,
	conf map[string]interface{},
	cfgPath string,
	constants map[string]interface{},
) (*Flattener, error) {
	var cfgParamName string

	selector := ""
	if entry.Recv != nil {
		selector = entry.Recv.List[0].Names[0].Name
	}

	for _, param := range entry.Type.Params.List {
		path, err := findPath(param)
		if err != nil {
			return nil, err
		}

		if path == cfgPath {
			cfgParamName = param.Names[0].Name
			break
		}
	}

	context := map[string]interface{}{
		cfgParamName: conf,
		"err":        nil,
	}

	for k, v := range constants {
		context[k] = v
	}

	return &Flattener{
		ctx:      context,
		cfgPath:  cfgPath,
		entry:    entry,
		selector: selector,
	}, nil
}

func (f *Flattener) Flatten() error {
	ast.Inspect(f.entry, f.inspect)

	return f.err
}

func (f *Flattener) inspect(
	n ast.Node,
) bool {
	if n == nil {
		return true
	}

	switch n := n.(type) {
	case *ast.BlockStmt:
		err := f.flattenBlockStmt(n)
		if err != nil {
			f.err = fmt.Errorf("error flattening block statement: %v", err)
			return false
		}
	}

	return true
}

func (f *Flattener) flattenBlockStmt(
	blockStmt *ast.BlockStmt,
) error {
	newList := []ast.Stmt{}

	for _, stmt := range blockStmt.List {
		switch stmt := stmt.(type) {
		case *ast.AssignStmt:
			var buf bytes.Buffer

			fset := token.NewFileSet()
			printer.Fprint(&buf, fset, stmt.Rhs[0])

			eval := goval.NewEvaluator()
			result, _ := eval.Evaluate(evaluator.Prepare(buf.String()), f.ctx, evaluator.Funcs)

			err := f.fillCtxWithResult(stmt.Lhs[0], result)
			if err != nil {
				return fmt.Errorf("error filling context with result: %v", err)
			}

			newList = append(newList, stmt)
		case *ast.IfStmt:
			res, err := f.flattenIfStmt(stmt)
			if err != nil {
				return fmt.Errorf("error flattening if statement: %v", err)
			}

			if res != nil {
				newList = append(newList, res...)
			}
		case *ast.SwitchStmt:
			res, err := f.flattenSwitchStmt(stmt)
			if err != nil {
				return fmt.Errorf("error flattening switch statement: %v", err)
			}

			if res != nil {
				newList = append(newList, res...)
			}
		default:
			newList = append(newList, stmt)
		}
	}

	blockStmt.List = newList

	return nil
}

func (f *Flattener) fillCtxWithResult(
	lhs ast.Expr,
	result interface{},
) error {
	switch lhs := lhs.(type) {
	case *ast.Ident:
		f.ctx[lhs.Name] = result
		return nil
	case *ast.SelectorExpr:
		name := lhs.X.(*ast.Ident).Name
		ctxEntry, ok := f.ctx[name]

		if !ok {
			f.ctx[name] = map[string]interface{}{}
			ctxEntry = f.ctx[name]
		}

		ctxEntry.(map[string]interface{})[lhs.Sel.Name] = result

		return nil
	case *ast.IndexExpr:
		// TODO: implement
		// this is case when we write to map, like this:
		// endpoints["/"] = common.New()
		// not sure if its needed
		return nil
	default:
		return fmt.Errorf("unknown lhs type: %T", lhs)
	}
}

func (f *Flattener) flattenIfStmt(ifStmt *ast.IfStmt) ([]ast.Stmt, error) {
	ast.Inspect(ifStmt.Cond, f.fillCtxWithNil)

	var buf bytes.Buffer

	fset := token.NewFileSet()
	printer.Fprint(&buf, fset, ifStmt.Cond)

	eval := goval.NewEvaluator()

	result, err := eval.Evaluate(buf.String(), f.ctx, evaluator.Funcs)
	if err != nil {
		if strings.HasPrefix(buf.String(), f.selector+".") {
			// here we let the if statement be, because we don't have enough
			// information to evaluate it
			return []ast.Stmt{ifStmt}, nil
		} else {
			return nil, fmt.Errorf("error evaluating if statement: %v", err)
		}
	}

	if result.(bool) {
		return f.flattenListStatements(ifStmt.Body.List)
	}

	if ifStmt.Else != nil {
		switch elseStmt := ifStmt.Else.(type) {
		case *ast.IfStmt:
			return f.flattenIfStmt(elseStmt)
		case *ast.BlockStmt:
			return elseStmt.List, nil
		default:
			return nil, errors.New("Unknown else statement")
		}
	}

	return nil, nil
}

func (f *Flattener) fillCtxWithNil(
	n ast.Node,
) bool {
	if n == nil {
		return true
	}

	switch n := n.(type) {
	case *ast.SelectorExpr:
		switch n.X.(type) {
		case *ast.Ident:
			selector := n.X.(*ast.Ident).Name
			if selector != f.selector {
				return true
			}

			m, ok := f.ctx[selector]
			if !ok {
				return true
			}

			switch m := m.(type) {
			case map[string]interface{}:
				m[n.Sel.Name] = nil
			}
		}
	}

	return true
}

func (f *Flattener) flattenSwitchStmt(switchStmt *ast.SwitchStmt) ([]ast.Stmt, error) {
	var buf bytes.Buffer

	err := printer.Fprint(&buf, token.NewFileSet(), switchStmt.Tag)
	if err != nil {
		return nil, err
	}

	left := buf.String()

	for _, stmt := range switchStmt.Body.List {
		caseStmt := stmt.(*ast.CaseClause)
		if len(caseStmt.List) == 0 {
			return caseStmt.Body, nil
		}

		for _, expr := range caseStmt.List {
			buf.Reset()

			err := printer.Fprint(&buf, token.NewFileSet(), expr)
			if err != nil {
				return nil, err
			}

			right := buf.String()
			evaluator := goval.NewEvaluator()

			result, err := evaluator.Evaluate(left+" == "+right, f.ctx, nil)
			if err != nil {
				return nil, err
			}

			if result.(bool) {
				body, err := f.flattenListStatements(caseStmt.Body)
				if err != nil {
					return nil, err
				}

				return body, nil
			}
		}
	}

	return nil, nil
}

func (f *Flattener) flattenListStatements(list []ast.Stmt) ([]ast.Stmt, error) {
	newList := []ast.Stmt{}

	for _, stmt := range list {
		switch stmt := stmt.(type) {
		case *ast.IfStmt:
			res, err := f.flattenIfStmt(stmt)
			if err != nil {
				return nil, err
			}

			if res != nil {
				newList = append(newList, res...)
			}
		case *ast.SwitchStmt:
			res, err := f.flattenSwitchStmt(stmt)
			if err != nil {
				return nil, err
			}

			if res != nil {
				newList = append(newList, res...)
			}
		default:
			newList = append(newList, stmt)
		}
	}

	return newList, nil
}
