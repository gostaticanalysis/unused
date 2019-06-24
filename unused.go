package unused

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "unused",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "unused find unused identifyers"

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	objects := map[string][]*ast.Ident{}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if !ast.IsExported(n.Name) && n.Name != "_" {
				if o := pass.TypesInfo.ObjectOf(n); !skip(o) {
					objects[o.Id()] = append(objects[o.Id()], n)
				}
			}
		}
	})

	for id := range objects {
		if len(objects[id]) == 1 {
			n := objects[id][0]
			pass.Reportf(n.Pos(), "%s is unused", n.Name)
		}
	}

	return nil, nil
}

func skip(o types.Object) bool {
	if o == nil {
		return true
	}

	switch o := o.(type) {
	case *types.Var:
		if o.Anonymous() {
			return true
		}
	case *types.Func:
		// main
		if o.Name() == "main" && o.Pkg().Name() == "main" {
			return true
		}

		// init
		if o.Name() == "init" && o.Pkg().Scope() == o.Parent() {
			return true
		}

		// method
		sig, ok := o.Type().(*types.Signature)
		if ok && sig.Recv() != nil {
			return true
		}
	}

	return false
}
