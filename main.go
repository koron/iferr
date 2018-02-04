package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

const noname = "(no name)"

var isNum = map[string]struct{}{
	"int":     struct{}{},
	"int16":   struct{}{},
	"int32":   struct{}{},
	"int64":   struct{}{},
	"uint":    struct{}{},
	"uint16":  struct{}{},
	"uint32":  struct{}{},
	"uint64":  struct{}{},
	"float":   struct{}{},
	"float32": struct{}{},
	"float64": struct{}{},
}

type visitor struct {
	pos token.Pos
	err error
	fd  *ast.FuncDecl
	fn  string
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	fd, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}
	fname := noname
	if fd.Name != nil {
		fname = fd.Name.Name
	}
	fname = fd.Name.Name
	if v.pos < fd.Pos() || v.pos > fd.End() {
		return nil
	}
	if fd.Type == nil {
		return v
	}
	v.fn = fname
	v.fd = fd
	return v
}

type field struct {
	name string
}

func toTypes(fl *ast.FieldList) []ast.Expr {
	if fl == nil || len(fl.List) == 0 {
		return nil
	}
	types := make([]ast.Expr, 0, len(fl.List))
	for _, f := range fl.List {
		types = append(types, f.Type)
	}
	return types
}

func typeString(x ast.Expr) string {
	switch t := x.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if _, ok := t.X.(*ast.Ident); ok {
			return typeString(t.X) + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.ArrayType:
		return "[]" + typeString(t.Elt)
	}
	return ""
}

func writeIferr(w io.Writer, types []ast.Expr) error {
	if len(types) == 0 {
		_, err := fmt.Fprint(w, "if err != nil {\n\treturn\n}\n")
		return err
	}
	bb := &bytes.Buffer{}
	bb.WriteString("if err != nil {\n\treturn ")
	for i, t := range types {
		if i > 0 {
			bb.WriteString(", ")
		}
		ts := typeString(t)
		if ts == "error" {
			bb.WriteString("err")
			continue
		}
		if ts == "string" {
			bb.WriteString(`""`)
			continue
		}
		if _, ok := isNum[ts]; ok {
			bb.WriteString("0")
			continue
		}
		if strings.HasPrefix(ts, "[]") {
			bb.WriteString("nil")
			continue
		}
		// treat it as an interface when type name has "."
		if strings.Index(ts, ".") >= 0 {
			bb.WriteString("nil")
			continue
		}
		// TODO: support more types.
		bb.WriteString(ts)
		bb.WriteString("{}")
	}
	bb.WriteString("\n}\n")
	io.Copy(w, bb)
	return nil
}

func iferr(w io.Writer, r io.Reader, pos int) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "iferr.go", r, 0)
	if err != nil {
		return err
	}
	v := &visitor{pos: token.Pos(pos)}
	ast.Walk(v, file)
	if v.err != nil {
		return err
	}
	if v.fd == nil {
		return fmt.Errorf("no functions at %d", pos)
	}
	types := toTypes(v.fd.Type.Results)
	return writeIferr(w, types)
}

func main() {
	var (
		pos int
	)
	flag.IntVar(&pos, "pos", 0, "position of cursor")
	flag.Parse()
	err := iferr(os.Stdout, os.Stdin, pos)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
