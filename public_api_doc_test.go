package run9

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestExportedIdentifiersHaveDocComments(t *testing.T) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(info os.FileInfo) bool {
		name := info.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse package: %v", err)
	}

	pkg := pkgs["run9"]
	if pkg == nil {
		t.Fatal("run9 package not found")
	}

	var missing []string
	for filename, file := range pkg.Files {
		base := filepath.Base(filename)
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range node.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						if spec.Name.IsExported() && !hasDocComment(firstDoc(spec.Doc, node.Doc), spec.Name.Name) {
							missing = append(missing, base+":"+spec.Name.Name)
						}
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if name.IsExported() && !hasDocComment(firstDoc(spec.Doc, node.Doc), name.Name) {
								missing = append(missing, base+":"+name.Name)
							}
						}
					}
				}
			case *ast.FuncDecl:
				if node.Name.IsExported() && !hasDocComment(node.Doc, node.Name.Name) {
					missing = append(missing, base+":"+node.Name.Name)
				}
			}
		}
	}

	if len(missing) == 0 {
		return
	}
	sort.Strings(missing)
	t.Fatalf("missing doc comments for exported identifiers:\n%s", strings.Join(missing, "\n"))
}

func TestPackageDocCommentExists(t *testing.T) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	file, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, "doc.go"), nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse doc.go: %v", err)
	}
	if !hasDocComment(file.Doc, "Package run9") {
		t.Fatal("doc.go must define a package doc comment starting with `Package run9`")
	}
}

func TestReadmeMentionsGodoc(t *testing.T) {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	text := string(data)
	for _, needle := range []string{
		"pkg.go.dev/github.com/sys9-ai/run9-sdk-go",
		"go doc github.com/sys9-ai/run9-sdk-go",
	} {
		if strings.Contains(text, needle) {
			continue
		}
		t.Fatalf("README.md must mention %q", needle)
	}
}

func hasDocComment(doc *ast.CommentGroup, name string) bool {
	if doc == nil {
		return false
	}
	text := strings.TrimSpace(doc.Text())
	if text == "" {
		return false
	}
	return strings.HasPrefix(text, name+" ") || text == name
}

func firstDoc(primary *ast.CommentGroup, fallback *ast.CommentGroup) *ast.CommentGroup {
	if primary != nil {
		return primary
	}
	return fallback
}
