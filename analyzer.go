package checkexhaustive

import (
	"go/ast"
	"go/token"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type sPos struct {
	filename string
	line     int
}

func positionToSPos(position token.Position) sPos {
	return sPos{
		filename: position.Filename,
		line:     position.Line,
	}
}

var Analyzer = &analysis.Analyzer{
	Name:     "checkexhaustive",
	Doc:      `Ensure exhaustive filling of struct literals labelled with "// check:exhaustive"`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.CompositeLit)(nil),
	}

	type comment struct {
		handled bool
		pos     token.Pos
	}
	cmts := map[sPos]comment{}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		switch stmt := node.(type) {
		case *ast.File:
			for _, c := range stmt.Comments {
				if len(c.List) == 1 && strings.HasPrefix(c.List[0].Text+" ", "// check:exhaustive ") {
					pos := c.Pos()
					sPos := positionToSPos(pass.Fset.Position(pos))
					cmts[sPos] = comment{
						handled: false,
						pos:     pos,
					}
				}
			}
		case *ast.CompositeLit:
			var ident *ast.Ident
			var pkgIdent *ast.Ident

			switch v := stmt.Type.(type) {
			case *ast.Ident:
				ident = v
			case *ast.SelectorExpr:
				ident = v.Sel
				if pIdent, ok := v.X.(*ast.Ident); ok {
					pkgIdent = pIdent
				}
			default:
				return
			}

			sPos := positionToSPos(pass.Fset.Position(ident.NamePos))
			sPos.line--
			if _, ok := cmts[sPos]; !ok {
				return
			}

			t := pass.TypesInfo.TypeOf(stmt.Type)
			if t == nil {
				return
			}
			str, ok := t.Underlying().(*types.Struct)
			if !ok {
				return
			}

			hasField := map[string]bool{}
			for i := 0; i < str.NumFields(); i++ {
				hasField[str.Field(i).Name()] = false
			}

			for _, ele := range stmt.Elts {
				if kvExpr, ok := ele.(*ast.KeyValueExpr); ok {
					hasField[kvExpr.Key.(*ast.Ident).Name] = true
				}
			}

			missingFields := []string{}
			for field, isSet := range hasField {
				if !isSet {
					missingFields = append(missingFields, field)
				}
			}
			sort.Strings(missingFields)

			if len(missingFields) > 0 {
				sName := ident.Name
				if pkgIdent != nil {
					sName = pkgIdent.Name + "." + sName
				}
				pass.Reportf(stmt.Pos(), "%s is missing fields: %s", sName, strings.Join(missingFields, ", "))
			}

			cmts[sPos] = comment{
				handled: true,
				pos:     cmts[sPos].pos,
			}
		}
	})

	for _, cmt := range cmts {
		if !cmt.handled {
			pass.Reportf(cmt.pos, "unmatched check:exhaustive comment")
		}
	}

	return nil, nil
}
