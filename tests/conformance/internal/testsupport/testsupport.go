package testsupport

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tliron/exturl"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"
)

func ParseFile(t *testing.T, path string) (*normal.ServiceTemplate, string, error) {
	t.Helper()

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("resolve fixture path: %v", err)
	}

	urlContext := exturl.NewContext()
	defer urlContext.Release()

	parserContext := parser.NewParser().NewContext()
	parserContext.URL = urlContext.NewFileURL(filepath.ToSlash(absolutePath))

	serviceTemplate, parseErr := parserContext.Parse(context.Background())
	problems := ""
	if parserContext.Root != nil {
		problems = parserContext.GetProblems().ToString(false)
	}
	return serviceTemplate, problems, parseErr
}

func ParseSource(t *testing.T, source string) (*normal.ServiceTemplate, string, error) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "service-template.yaml")
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return ParseFile(t, path)
}
