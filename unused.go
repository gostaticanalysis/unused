package unused

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer find unused identifyers.
var Analyzer = &analysis.Analyzer{
	Name: "unused",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const doc = "unused find unused identifyers"

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	objects := map[types.Object][]*ast.Ident{}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if !ast.IsExported(n.Name) && n.Name != "_" {
				if o := pass.TypesInfo.ObjectOf(n); !skip(o) {
					objects[o] = append(objects[o], n)
				}
			}
		}
	})

	for o := range objects {
		if len(objects[o]) == 1 {
			n := objects[o][0]
			pass.Reportf(n.Pos(), "%s is unused", n.Name)
		}
	}

	return nil, nil
}

func skip(o types.Object) bool {

	if o == nil || o.Parent() == types.Universe {
		return true
	}

	switch o := o.(type) {
	case *types.PkgName:
		return true
	case *types.Var:
		if o.Pkg().Scope() != o.Parent() &&
			!(o.IsField() && !o.Anonymous() && isFieldInNamedStruct(o)) {
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
		if ok {
			if recv := sig.Recv(); recv != nil {
				for _, i := range interfaces(o.Pkg()) {
					if i == recv.Type() ||
						(types.Implements(recv.Type(), i) && has(i, o)) {
						return true
					}
				}
			}
		}
	}

	return false
}

func interfaces(pkg *types.Package) []*types.Interface {
	var ifs []*types.Interface

	for _, n := range pkg.Scope().Names() {
		o := pkg.Scope().Lookup(n)
		if o != nil {
			i, ok := o.Type().Underlying().(*types.Interface)
			if ok {
				ifs = append(ifs, i)
			}
		}
	}

	return ifs
}

func has(intf *types.Interface, m *types.Func) bool {
	for i := 0; i < intf.NumMethods(); i++ {
		if intf.Method(i).Name() == m.Name() {
			return true
		}
	}
	return false
}

func isFieldInNamedStruct(v *types.Var) bool {
	structs := allNamedStructs(v.Pkg())
	for _, s := range structs {
		for i := 0; i < s.NumFields(); i++ {
			if s.Field(i) == v {
				return true
			}
		}
	}
	return false
}

func allNamedStructs(pkg *types.Package) []*types.Struct {
	var structs []*types.Struct

	for _, n := range pkg.Scope().Names() {
		o := pkg.Scope().Lookup(n)
		if o != nil {
			switch t := o.Type().(type) {
			case *types.Named:
				switch u := t.Underlying().(type) {
				case *types.Struct:
					structs = append(structs, u)
				}
			}
		}
	}

	return structs
}
