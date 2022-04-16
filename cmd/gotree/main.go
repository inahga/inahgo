package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/inahga/inahgo/rick/blowdash"
)

// main packages are represented as `command` instead
// packages:
//  - constants
//  - vars
//  - funcs
//  - types
//    - methods
//  - source files
//  - directories

// run in server mode

// flags
//   - server [localhost:8080]
//   - expand-all
//   - expand [things,to,expand]
//   - pretty
//   - root (path to where go.mod is)
//   - path (relative to root to start analyzing from)

type foo = []string

type (
	Object struct {
		Name     string `json:",omitempty"`
		LongName string `json:",omitempty"`
		Position string `json:",omitempty"`
	}

	Type struct {
		Object
		Methods []Object `json:",omitempty"`
	}

	PackageKind string

	Package struct {
		Name                                            string
		Kind                                            PackageKind
		Packages                                        []*Package `json:",omitempty"`
		Consts, Vars, Funcs, Tests, Interfaces, Aliases []*Object  `json:",omitempty"`
		Types                                           []*Type    `json:",omitempty"`
		Sources                                         []string   `json:",omitempty"`
	}
)

const (
	StubPackage    PackageKind = ""
	NormalPackage  PackageKind = "package"
	CommandPackage PackageKind = "command"
	TestPackage    PackageKind = "test"
)

func printj(i any) {
	byt, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(byt))
}

type baz func(foo, bar string, baz int, bin []string, bla map[string]string)

func typeString(t any) string {
	switch t := t.(type) {
	case *ast.Ident:
		return t.String()
	case *ast.ArrayType:
		return "[]" + fmt.Sprint(t.Elt)
	case *ast.FuncType:
		ret := "func("
		if t.Params != nil && len(t.Params.List) > 0 {
			for _, param := range t.Params.List {
				for _, name := range param.Names {
					ret += name.String() + ", "
				}
				ret = ret[:len(ret)-2] + " " + fmt.Sprint(param.Type) + ", "
			}
			ret = ret[:len(ret)-2]
		}
		ret += ")"
		return ret
	case *ast.MapType:
	case *ast.ChanType:
	case *ast.InterfaceType:
	case *ast.StructType:
		fmt.Println(t)
	default:
		return ""
	}
	return ""
}

func makePackage(fset *token.FileSet, root string, pkgs map[string]*ast.Package) (normal *Package, test *Package, err error) {
	for _, pkg := range pkgs {
		newpkg := &Package{
			Name: pkg.Name,
		}

		if strings.HasSuffix(pkg.Name, "_test") {
			test = newpkg
			newpkg.Kind = TestPackage
			newpkg.Name = pkg.Name
		} else if pkg.Name == "main" {
			// A command package is mostly handled the same as a normal package, so just
			// reuse that variable.
			normal = newpkg
			newpkg.Kind = CommandPackage
			newpkg.Name = root
		} else {
			normal = newpkg
			newpkg.Kind = NormalPackage
			newpkg.Name = pkg.Name
		}

		newpkg.Sources = blowdash.Keys(pkg.Files)

		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					// fmt.Println(d)
					// need to capture functions vs methods

				case *ast.GenDecl:
					for _, spec := range d.Specs {
						if d.Tok == token.CONST || d.Tok == token.VAR {
							s := spec.(*ast.ValueSpec)
							name := blowdash.SliceStringer(s.Names, ", ")
							longName := name
							if s.Type != nil {
								ident, ok := s.Type.(*ast.Ident)
								if ok {
									longName = longName + " " + ident.String()
								}
							}
							newobj := &Object{
								Name:     name,
								LongName: longName,
								Position: fset.Position(s.Pos()).String(),
							}
							if d.Tok == token.CONST {
								newobj.LongName = "const " + newobj.LongName
								newpkg.Consts = append(newpkg.Consts, newobj)
							} else {
								newobj.LongName = "var " + newobj.LongName
								newpkg.Vars = append(newpkg.Vars, newobj)
							}
						} else if d.Tok == token.TYPE {
							s := spec.(*ast.TypeSpec)
							// TypeSpec struct {
							// 	Doc        *CommentGroup // associated documentation; or nil
							// 	Name       *Ident        // type name
							// 	TypeParams *FieldList    // type parameters; or nil
							// 	Assign     token.Pos     // position of '=', if any
							// 	Type       Expr          // *Ident, *ParenExpr, *SelectorExpr, *StarExpr, or any of the *XxxTypes
							// 	Comment    *CommentGroup // line comments; or nil
							// }
							fmt.Println(s.Name)
							fmt.Println(typeString(s.Type))

							// type aliases do not have methods
						}
					}
				}
			}
			// 	for _, object := range f.Scope.Objects {
			// 		position := fset.Position(object.Pos()).String()
			// 		newobj := &Object{
			// 			Name:     object.Name,
			// 			Position: position,
			// 		}

			// 		switch object.Kind {
			// 		case ast.Con:
			// 			newobj.LongName = object.Name
			// 			if decl, ok := object.Decl.(*ast.ValueSpec); ok {
			// 				if decl.Type != nil {
			// 					newobj.LongName = fmt.Sprintf("%s %s", object.Name, decl.Type)
			// 				}
			// 			}
			// 			newpkg.Consts = append(newpkg.Consts, newobj)

			// 		case ast.Fun:
			// 			// decl := object.Decl.(*ast.FuncDecl)
			// 			// blowdash.ForEach(decl.Type.Params.List, func(e *ast.Field) {
			// 			// 	fmt.Println(e)
			// 			// })

			// 			if strings.HasPrefix(object.Name, "Test") {
			// 				newpkg.Tests = append(newpkg.Tests, newobj)
			// 			} else {
			// 				newpkg.Funcs = append(newpkg.Funcs, newobj)
			// 			}

			// 		case ast.Typ:

			// 		case ast.Var:
			// 			newobj.LongName = object.Name
			// 			if decl, ok := object.Decl.(*ast.ValueSpec); ok {
			// 				if decl.Type != nil {
			// 					newobj.LongName = fmt.Sprintf("%s %s", object.Name, decl.Type)
			// 				}
			// 			}
			// 			newpkg.Vars = append(newpkg.Consts, newobj)
			// 		}
			// 	}
		}
	}
	return normal, test, nil
}

func scanPackageHelper(fset *token.FileSet, path string, curr int, depth int) (*Package, bool, error) {
	var (
		include      bool
		normal, test *Package
		err          error
	)
	if curr >= depth {
		return nil, false, nil
	}

	pkgs, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, false, err
	}

	if len(pkgs) > 0 {
		// Marking this path for inclusion indicates that paths below should be included
		// in the final tree.
		include = true

		root, err := filepath.Abs(path)
		if err != nil {
			return nil, false, err
		}
		normal, test, err = makePackage(fset, filepath.Base(root), pkgs)
		if err != nil {
			return nil, false, err
		}
	} else {
		// Place a stub package as the normal package, in case higher subpaths contain
		// go source.
		normal = &Package{Name: filepath.Base(path), Kind: StubPackage}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, false, err
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		add, incl, err := scanPackageHelper(fset, filepath.Join(path, file.Name()), curr+1, depth)
		if err != nil {
			return nil, false, err
		}

		// This handles the case that the folder is not a package but subpaths do
		// contain go. This is common with the `cmd` folder pattern.
		if incl {
			normal.Packages = append(normal.Packages, add)
			include = true
		}
	}

	if test != nil {
		// Consider a test package a child of the normal package, even though this is
		// not strictly true.
		normal.Packages = append(normal.Packages, test)
	}
	return normal, include, nil
}

func scanPackage(path string, depth int) (*Package, error) {
	ret, _, err := scanPackageHelper(token.NewFileSet(), path, 0, depth)
	return ret, err
}

func main() {
	root := os.Args[1]
	pkg, err := scanPackage(root, 3)
	if err != nil {
		panic(err)
	}
	printj(pkg)
	// fset := token.NewFileSet()

	// if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if d.Type().IsDir() {

	// 		pkgs, err := parser.ParseDir(fset, path, nil, 0)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		if len(pkgs) != 0 {
	// 			fmt.Println(pkgs)
	// 		}
	// 	}

	// 	return nil
	// }); err != nil {
	// 	panic(err)
	// }

	// for n, pkg := range pkgs {
	// 	fmt.Println(n)
	// 	fmt.Println(pkg)
	// 	for _, f := range pkg.Files {
	// 		fmt.Println(f.Scope)
	// 	}
	// }
	// printj(pkgs)
}
