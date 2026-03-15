package service

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/GlebMoskalev/go-path-backend/internal/model"
	"go.uber.org/zap"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

type FormatContextProvider interface {
	GetFormatContext(projectSlug, stepSlug string) (*model.FormatContext, error)
}

type FormatService struct {
	log     *zap.Logger
	project FormatContextProvider
}

func NewFormatService(log *zap.Logger, project FormatContextProvider) *FormatService {
	return &FormatService{log: log, project: project}
}

func (s *FormatService) FormatCode(code string) (string, error) {
	opts := &imports.Options{
		Fragment:   true,
		Comments:   true,
		TabIndent:  true,
		TabWidth:   8,
		FormatOnly: false,
	}

	formatted, err := imports.Process("", []byte(code), opts)
	if err != nil {
		s.log.Error("failed to format code", zap.Error(err))
		return code, err
	}

	return string(formatted), nil
}

func (s *FormatService) FormatProjectCode(code, projectSlug, stepSlug string) (string, error) {
	fmtCtx, err := s.project.GetFormatContext(projectSlug, stepSlug)
	if err != nil {
		return code, err
	}

	localPkgs := buildLocalPackageMap(fmtCtx)
	code = addMissingLocalImports(code, localPkgs)

	return s.FormatCode(code)
}

func buildLocalPackageMap(ctx *model.FormatContext) map[string]string {
	moduleName := parseModuleName(ctx.GoMod)
	if moduleName == "" {
		return nil
	}

	pkgs := make(map[string]string)
	for filePath, content := range ctx.Files {
		dir := filepath.Dir(filePath)
		if dir == "." {
			continue
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", content, parser.PackageClauseOnly)
		if err != nil {
			continue
		}
		pkgs[f.Name.Name] = moduleName + "/" + dir
	}
	return pkgs
}

func addMissingLocalImports(code string, localPkgs map[string]string) string {
	if len(localPkgs) == 0 {
		return code
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return code
	}

	imported := make(map[string]bool, len(f.Imports))
	for _, imp := range f.Imports {
		imported[strings.Trim(imp.Path.Value, `"`)] = true
	}

	usedNames := make(map[string]bool)
	ast.Inspect(f, func(n ast.Node) bool {
		sel, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		usedNames[ident.Name] = true
		return true
	})

	modified := false
	for name, importPath := range localPkgs {
		if usedNames[name] && !imported[importPath] {
			astutil.AddImport(fset, f, importPath)
			modified = true
		}
	}

	if !modified {
		return code
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return code
	}
	return buf.String()
}

func parseModuleName(gomod string) string {
	for _, line := range strings.Split(gomod, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}
