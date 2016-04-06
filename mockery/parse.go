package mockery

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
)

type Parser struct {
	file *ast.File
	path string

	pkg *types.Package
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(path string) error {
	fset := token.NewFileSet()

	// Parse the file containing this very example
	// but stop after processing the imports.
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Type-check a package consisting of this file.
	// Type information for the imported packages
	// comes from $GOROOT/pkg/$GOOS_$GOOARCH/fmt.a.
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(abs, fset, []*ast.File{f}, nil)
	if err != nil {
		return err
	}

	p.path = abs
	p.file = f
	p.pkg = pkg
	return nil
}

func (p *Parser) Find(name string) (*Interface, error) {
	typ := p.pkg.Scope().Lookup(name).Type().(*types.Named)

	name = typ.Obj().Name()

	iface := typ.Underlying().(*types.Interface).Complete()

	return &Interface{name, p.path, p.file, iface}, nil
}

/*
func (p *Parser) FindOld(name string) (*Interface, error) {
	for _, decl := range p.file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range gen.Specs {
				if typespec, ok := spec.(*ast.TypeSpec); ok {
					if typespec.Name.Name == name {
						if iface, ok := typespec.Type.(*ast.InterfaceType); ok {
							return &Interface{name, p.path, p.file, iface}, nil
						} else {
							return nil, ErrNotInterface
						}
					}
				}
			}
		}
	}
	return nil, nil
}
*/

type Interface struct {
	Name string
	Path string
	File *ast.File
	Type *types.Interface
}

func (p *Parser) Interfaces() []*Interface {
	return nil
	/*
		var ifaces []*Interface

		for _, decl := range p.file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range gen.Specs {
					if typespec, ok := spec.(*ast.TypeSpec); ok {
						if iface, ok := typespec.Type.(*ast.InterfaceType); ok {
							ifaces = append(ifaces, &Interface{typespec.Name.Name, p.path, p.file, nil})
						}
					}
				}
			}
		}

		return ifaces
	*/
}
