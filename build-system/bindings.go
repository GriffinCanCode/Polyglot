package builder

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// BindingGenerator creates language bindings from type definitions
type BindingGenerator struct {
	sourceDir string
	outputDir string
}

// NewBindingGenerator creates a binding generator
func NewBindingGenerator(sourceDir, outputDir string) *BindingGenerator {
	return &BindingGenerator{
		sourceDir: sourceDir,
		outputDir: outputDir,
	}
}

// Generate creates bindings for all configured languages
func (g *BindingGenerator) Generate(languages []string) error {
	// Parse Go source files for type definitions
	types, err := g.parseTypes()
	if err != nil {
		return fmt.Errorf("failed to parse types: %w", err)
	}

	// Generate bindings for each language
	for _, lang := range languages {
		if err := g.generateForLanguage(lang, types); err != nil {
			return fmt.Errorf("failed to generate %s bindings: %w", lang, err)
		}
	}

	return nil
}

// parseTypes extracts type definitions from Go source
func (g *BindingGenerator) parseTypes() ([]TypeDef, error) {
	fset := token.NewFileSet()

	var types []TypeDef

	err := filepath.Walk(g.sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// Extract type definitions
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				if genDecl.Tok == token.TYPE {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							types = append(types, g.extractType(typeSpec))
						}
					}
				}
			}
		}

		return nil
	})

	return types, err
}

// extractType converts AST type to TypeDef
func (g *BindingGenerator) extractType(spec *ast.TypeSpec) TypeDef {
	def := TypeDef{
		Name:   spec.Name.Name,
		Fields: []Field{},
	}

	if structType, ok := spec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			for _, name := range field.Names {
				def.Fields = append(def.Fields, Field{
					Name: name.Name,
					Type: g.exprToString(field.Type),
				})
			}
		}
	}

	return def
}

// exprToString converts AST expression to string
func (g *BindingGenerator) exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + g.exprToString(t.X)
	case *ast.ArrayType:
		return "[]" + g.exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + g.exprToString(t.Key) + "]" + g.exprToString(t.Value)
	default:
		return "interface{}"
	}
}

// generateForLanguage creates bindings for a specific language
func (g *BindingGenerator) generateForLanguage(lang string, types []TypeDef) error {
	switch lang {
	case "typescript":
		return g.generateTypeScript(types)
	case "python":
		return g.generatePython(types)
	case "rust":
		return g.generateRust(types)
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}
}

// generateTypeScript creates TypeScript definitions
func (g *BindingGenerator) generateTypeScript(types []TypeDef) error {
	tmpl := template.Must(template.New("typescript").Parse(tsTemplate))

	outputPath := filepath.Join(g.outputDir, "bindings.d.ts")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, types)
}

// generatePython creates Python type stubs
func (g *BindingGenerator) generatePython(types []TypeDef) error {
	tmpl := template.Must(template.New("python").Parse(pyTemplate))

	outputPath := filepath.Join(g.outputDir, "bindings.pyi")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, types)
}

// generateRust creates Rust bindings
func (g *BindingGenerator) generateRust(types []TypeDef) error {
	tmpl := template.Must(template.New("rust").Parse(rustTemplate))

	outputPath := filepath.Join(g.outputDir, "bindings.rs")
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, types)
}

// TypeDef represents a type definition
type TypeDef struct {
	Name   string
	Fields []Field
}

// Field represents a struct field
type Field struct {
	Name string
	Type string
}

// mapGoTypeToTS maps Go types to TypeScript
func mapGoTypeToTS(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		if strings.HasPrefix(goType, "[]") {
			return mapGoTypeToTS(goType[2:]) + "[]"
		}
		return "any"
	}
}

// Templates for code generation
const tsTemplate = `// Auto-generated TypeScript bindings
{{range .}}
export interface {{.Name}} {
{{- range .Fields}}
  {{.Name}}: {{mapGoTypeToTS .Type}};
{{- end}}
}
{{end}}
`

const pyTemplate = `# Auto-generated Python type stubs
from typing import Any, List, Dict

{{range .}}
class {{.Name}}:
{{- range .Fields}}
    {{.Name}}: {{mapGoTypeToPy .Type}}
{{- end}}
{{end}}
`

const rustTemplate = `// Auto-generated Rust bindings
use serde::{Deserialize, Serialize};

{{range .}}
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct {{.Name}} {
{{- range .Fields}}
    pub {{.Name}}: {{mapGoTypeToRust .Type}},
{{- end}}
}
{{end}}
`
