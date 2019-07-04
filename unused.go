package unused

import (
	"go/types"

	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/ident"
	"golang.org/x/tools/go/analysis"
)

// Analyzer find unused identifyers.
var Analyzer = &analysis.Analyzer{
	Name: "unused",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		ident.Analyzer,
	},
}

const doc = "unused find unused identifyers"

func run(pass *analysis.Pass) (interface{}, error) {
	m := pass.ResultOf[ident.Analyzer].(ident.Map)
	for o := range m {
		if !skip(o) && len(m[o]) == 1 {
			n := m[o][0]
			pass.Reportf(n.Pos(), "%s is unused", n.Name)
		}
	}
	return nil, nil
}

func skip(o types.Object) bool {

	if o == nil || o.Parent() == types.Universe || o.Exported() {
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
				for _, i := range analysisutil.Interfaces(o.Pkg()) {
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
