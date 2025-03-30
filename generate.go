//go:build ignore
// +build ignore

// generate.go
//
// Run with: go run generate.go
// This generator scans the standard math package and produces:
//   - core_functions.go
//   - core_consts.go
//   - core_vars.go
//   - core_types.go

package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// FuncInfo holds information about a math package function.
type FuncInfo struct {
	Name          string // Function name.
	Params        string // Original (non‑generic) parameter list.
	GenericParams string // Generated generic parameters (e.g. "x N, y N").
	CastArgs      string // Arguments cast to float64 (e.g. "float64(x), float64(y)").
	ReturnType    string // Return type (expected to be "float64" for generic wrappers).
	IsGeneric     bool   // True if all parameters and the return type are float64 and there's only one return value.
	OriginalSig   string // The original function signature (for reference).
}

// ConstInfo holds information about a constant.
type ConstInfo struct {
	Name  string
	Value string
}

// VarInfo holds information about a variable.
type VarInfo struct {
	Name  string
	Value string
}

// TypeInfo holds information about a type.
type TypeInfo struct {
	Name string
	Decl string
}

func main() {
	// Locate the standard math package.
	pkg, err := build.Import("math", "", 0)
	if err != nil {
		log.Fatalf("failed to import math package: %v", err)
	}

	fset := token.NewFileSet()
	var funcs []FuncInfo
	var consts []ConstInfo
	var vars []VarInfo
	var types []TypeInfo

	// Process each Go file in the math package.
	for _, file := range pkg.GoFiles {
		filePath := filepath.Join(pkg.Dir, file)
		f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("failed to parse file %s: %v", filePath, err)
		}

		for _, decl := range f.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				// Only process exported, top‑level functions.
				if d.Name.IsExported() && d.Recv == nil {
					isGeneric := true
					var genericParams []string
					var castArgs []string
					var nonGenericParams []string

					if d.Type.Params != nil {
						for i, field := range d.Type.Params.List {
							// For a generic wrapper we require the parameter to be a float64.
							ident, ok := field.Type.(*ast.Ident)
							if !ok || ident.Name != "float64" {
								isGeneric = false
							}
							var names []string
							if len(field.Names) == 0 {
								// Generate a name if none is provided.
								name := fmt.Sprintf("arg%d", i)
								names = append(names, name)
							} else {
								for _, n := range field.Names {
									names = append(names, n.Name)
								}
							}
							// For the generic wrapper, use type N.
							genericParams = append(genericParams, strings.Join(names, ", ")+" N")
							// Cast each parameter to float64 for the call.
							for _, n := range names {
								castArgs = append(castArgs, "float64("+n+")")
							}
							// Build the original parameter string.
							var typeBuf strings.Builder
							printer.Fprint(&typeBuf, fset, field.Type)
							nonGenericParams = append(nonGenericParams, strings.Join(names, ", ")+" "+typeBuf.String())
						}
					}
					nonGenericParamsStr := strings.Join(nonGenericParams, ", ")
					genericParamsStr := strings.Join(genericParams, ", ")
					castArgsStr := strings.Join(castArgs, ", ")

					// Process the result type.
					retType := ""
					if d.Type.Results != nil && len(d.Type.Results.List) == 1 {
						// Check if the single Field has more than one name (i.e. multiple returns)
						if d.Type.Results.List[0].Names != nil && len(d.Type.Results.List[0].Names) > 1 {
							isGeneric = false
						}
						ident, ok := d.Type.Results.List[0].Type.(*ast.Ident)
						if !ok || ident.Name != "float64" {
							isGeneric = false
						}
						var retBuf strings.Builder
						printer.Fprint(&retBuf, fset, d.Type.Results.List[0].Type)
						retType = retBuf.String()
					} else {
						// Functions with zero or multiple return values are not made generic.
						isGeneric = false
					}

					var sigBuf strings.Builder
					printer.Fprint(&sigBuf, fset, d.Type)

					fi := FuncInfo{
						Name:          d.Name.Name,
						Params:        nonGenericParamsStr,
						GenericParams: genericParamsStr,
						CastArgs:      castArgsStr,
						ReturnType:    retType,
						IsGeneric:     isGeneric,
						OriginalSig:   sigBuf.String(),
					}
					funcs = append(funcs, fi)
				}
			case *ast.GenDecl:
				switch d.Tok {
				case token.CONST:
					for _, spec := range d.Specs {
						vspec := spec.(*ast.ValueSpec)
						for i, name := range vspec.Names {
							if name.IsExported() {
								var valueBuf strings.Builder
								if i < len(vspec.Values) {
									printer.Fprint(&valueBuf, fset, vspec.Values[i])
								}
								consts = append(consts, ConstInfo{
									Name:  name.Name,
									Value: valueBuf.String(),
								})
							}
						}
					}
				case token.VAR:
					for _, spec := range d.Specs {
						vspec := spec.(*ast.ValueSpec)
						for i, name := range vspec.Names {
							if name.IsExported() {
								var valueBuf strings.Builder
								if i < len(vspec.Values) {
									printer.Fprint(&valueBuf, fset, vspec.Values[i])
								}
								vars = append(vars, VarInfo{
									Name:  name.Name,
									Value: valueBuf.String(),
								})
							}
						}
					}
				case token.TYPE:
					for _, spec := range d.Specs {
						tspec := spec.(*ast.TypeSpec)
						if tspec.Name.IsExported() {
							var declBuf strings.Builder
							printer.Fprint(&declBuf, fset, d)
							types = append(types, TypeInfo{
								Name: tspec.Name.Name,
								Decl: declBuf.String(),
							})
						}
					}
				}
			}
		}
	}

	// Use a standard quoted string for the header.
	const header = "// Code generated by go:generate; DO NOT EDIT.\n\npackage math\n"

	// Generate core_functions.go
	funcFile, err := os.Create("core_functions.go")
	if err != nil {
		log.Fatalf("failed to create core_functions.go: %v", err)
	}
	defer funcFile.Close()

	funcTmplText := header + `
import "math"

// Core functions: wrappers for functions in the standard math package.
{{range .}}
// {{.Name}} {{if .IsGeneric}}wraps math.{{.Name}} in a generic function.{{else}}is a direct alias to math.{{.Name}}{{end}}.
{{if .IsGeneric}}
func {{.Name}}[N Number]({{.GenericParams}}) N {
	return N(math.{{.Name}}({{.CastArgs}}))
}
{{else}}
// Direct alias.
var {{.Name}} = math.{{.Name}}
{{end}}
{{end}}
`
	funcTmpl := template.Must(template.New("functions").Parse(funcTmplText))
	if err := funcTmpl.Execute(funcFile, funcs); err != nil {
		log.Fatalf("failed to execute template for core_functions.go: %v", err)
	}

	// Generate core_consts.go
	constFile, err := os.Create("core_consts.go")
	if err != nil {
		log.Fatalf("failed to create core_consts.go: %v", err)
	}
	defer constFile.Close()

	constTmplText := header + `
import "math"

// Core constants: re-exported from the standard math package.
const (
{{range .}}	{{.Name}} = math.{{.Name}}
{{end}})
`
	constTmpl := template.Must(template.New("consts").Parse(constTmplText))
	if err := constTmpl.Execute(constFile, consts); err != nil {
		log.Fatalf("failed to execute template for core_consts.go: %v", err)
	}

	// Generate core_vars.go
	varFile, err := os.Create("core_vars.go")
	if err != nil {
		log.Fatalf("failed to create core_vars.go: %v", err)
	}
	defer varFile.Close()

	varTmplText := header + `
import (
	"math"
)
// Core variables: re-exported from the standard math package.
var (
{{range .}}	{{.Name}} = math.{{.Name}}
{{end}})
var _ = math.Pi // dummy usage to avoid unused import error.
`
	varTmpl := template.Must(template.New("vars").Parse(varTmplText))
	if err := varTmpl.Execute(varFile, vars); err != nil {
		log.Fatalf("failed to execute template for core_vars.go: %v", err)
	}

	// Generate core_types.go
	typeFile, err := os.Create("core_types.go")
	if err != nil {
		log.Fatalf("failed to create core_types.go: %v", err)
	}
	defer typeFile.Close()

	typeTmplText := header + `
import (
	"math"
)
// Core types: re-exported from the standard math package.
{{range .}}
{{.Decl}}

{{end}}
var _ = math.Pi // dummy usage to avoid unused import error.
`
	typeTmpl := template.Must(template.New("types").Parse(typeTmplText))
	if err := typeTmpl.Execute(typeFile, types); err != nil {
		log.Fatalf("failed to execute template for core_types.go: %v", err)
	}

	log.Println("Core files generated successfully.")
}
