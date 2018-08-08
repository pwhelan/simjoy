package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var fd *os.File

// genDecl processes one declaration clause.
func genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	injoystick := false
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		name := vspec.Names[0].Name
		value := vspec.Values[0].(*ast.BasicLit).Value

		if name == "BTN_DIGI" || name == "ABS_PRESSURE" {
			injoystick = false
		}
		if injoystick == true {
			fmt.Fprintf(fd, "\t%s = %s\n", name, value)
		}
		if name == "BTN_JOYSTICK" || name == "REL_CNT" {
			injoystick = true
		}
	}
	return true
}

// genDecl processes one declaration clause.
func genLuaTable(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// We only care about const declarations.
		return true
	}
	injoystick := false
	for _, spec := range decl.Specs {
		vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
		name := vspec.Names[0].Name
		value := vspec.Values[0].(*ast.BasicLit).Value

		if name == "BTN_DIGI" || name == "ABS_PRESSURE" {
			injoystick = false
		}
		if injoystick == true {
			fmt.Fprintf(fd, "\tt.RawSetString(\"%s\", lua.LNumber(%s))\n", name, value)
		}
		if name == "BTN_JOYSTICK" || name == "REL_CNT" {
			injoystick = true
		}
	}
	return true
}

func main() {
	fset := token.NewFileSet()

	// Parse src but stop after processing the imports.
	cwd, _ := os.Getwd()
	fname := strings.Join([]string{
		cwd,
		"../../vendor/src/github.com/ynsta/uinput/const.go",
	},
		"/")
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	fd, err = os.OpenFile("const_linux.go", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintln(fd, "package vjoy")
	fmt.Fprintln(fd, "import (")
	fmt.Fprintln(fd, `	lua "github.com/yuin/gopher-lua"`)
	fmt.Fprintln(fd, ")")
	fmt.Fprintln(fd, "const (")
	ast.Inspect(f, genDecl)
	fmt.Fprintln(fd, ")")
	fmt.Fprintln(fd, "func Lua(L *lua.LState) {")
	fmt.Fprintln(fd, "	t := L.NewTable()")
	ast.Inspect(f, genLuaTable)
	fmt.Fprintln(fd, "	L.SetGlobal(\"joystick\", t)")
	fmt.Fprintln(fd, "}")
}
