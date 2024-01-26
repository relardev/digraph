package globals

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"github.com/maja42/goval"
)

func Get(packages map[string]*ast.Package) (map[string]interface{}, error) {
	g := map[string]interface{}{}
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					if decl.Tok == token.CONST || decl.Tok == token.VAR {
						for _, spec := range decl.Specs {
							valueSpec := spec.(*ast.ValueSpec)
							for i, name := range valueSpec.Names {
								if valueSpec.Values[i] != nil {
									var buf bytes.Buffer
									err := printer.Fprint(&buf, token.NewFileSet(), valueSpec.Values[i])
									if err != nil {
										return nil, fmt.Errorf("error printing constant: %w", err)
									}
									evaluator := goval.NewEvaluator()
									helperMap := map[string]interface{}{
										"time": map[string]interface{}{
											"Second": 1,
											"Minute": 60,
											"Hour":   3600,
										},
									}
									result, err := evaluator.Evaluate(buf.String(), helperMap, nil)
									if err != nil {
										g[name.Name] = buf.String()
									}
									g[name.Name] = result
								}
							}
						}
					}
				}
			}
		}
	}
	return g, nil
}
