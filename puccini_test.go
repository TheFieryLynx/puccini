package main

import (
	contextpkg "context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/exturl"
	cloutpkg "github.com/tliron/go-puccini/clout"
	"github.com/tliron/go-puccini/clout/js"
	"github.com/tliron/go-puccini/normal"
	"github.com/tliron/go-puccini/tosca/parser"

	_ "github.com/tliron/commonlog/simple"
)

func TestParse(t *testing.T) {
	context := NewContext(t)
	defer context.urlContext.Release()

	context.compileAll()
}

func BenchmarkParse(b *testing.B) {
	context := NewContext(b)
	defer context.urlContext.Release()

	for i := 0; i < b.N; i++ {
		context.compileAll()
	}
}

// TestHighConcurrency runs the same test 100 times in parallel to detect race conditions
// This test specifically targets URL resolution issues that occur under high concurrency
func TestHighConcurrency(t *testing.T) {
	const numConcurrent = 100

	// Test a representative sample of files that import remote profiles
	testFiles := []struct {
		url    string
		inputs map[string]any
	}{
		{"2.0/unicode.yaml", nil},
		{"2.0/descriptions.yaml", nil},
		{"2.0/functions.yaml", nil},
		{"2.0/inputs-and-outputs.yaml", map[string]any{"ram": "1 GiB"}},
	}

	for _, tf := range testFiles {
		tf := tf // capture range variable
		t.Run(tf.url, func(t *testing.T) {
			// Run the same file 100 times in parallel
			for i := 0; i < numConcurrent; i++ {
				i := i // capture loop variable
				t.Run(fmt.Sprintf("run-%d", i), func(t *testing.T) {
					t.Parallel()
					context := NewContext(t)
					defer context.urlContext.Release()
					context.compile_(t, tf.url, tf.inputs)
				})
			}
		})
	}
}

func TestJavaScriptErrors(t *testing.T) {
	context := NewContext(t)
	defer context.urlContext.Release()

	context.compileFailure("javascript/err-infinite-loop.yaml", nil)
	context.compileFailure("javascript/err-infinite-loop-v2.yaml", nil)
}

func (self *Context) compileFailure(url string, inputs map[string]any) {
	if t, ok := self.tb.(*testing.T); ok {
		t.Run(url, func(t_ *testing.T) {
			t_.Parallel()
			self.compileFailure_(t_, url, inputs)
		})
	} else {
		self.compileFailure_(self.tb, url, inputs)
	}
}

func (self *Context) compileFailure_(t testing.TB, url string, inputs map[string]any) {
	var normalServiceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var err error

	url_ := self.urlContext.NewFileURL(path.Join(filepath.ToSlash(self.root), "examples", url))

	parserContext := self.parser.NewContext()
	parserContext.URL = url_
	parserContext.Inputs = inputs
	if normalServiceTemplate, err = parserContext.Parse(contextpkg.TODO()); err != nil {
		return
	}

	problems := parserContext.GetProblems()
	if clout, err = normalServiceTemplate.Compile(); err != nil {
		return
	}

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   problems,
		URLContext: self.urlContext,
		History:    true,
		Format:     "yaml",
		Pretty:     true,
	}

	execContext.Resolve()
	execContext.Coerce()

	if problems.Empty() {
		t.Errorf("expected problems for %s", url)
	} else if url == "javascript/err-infinite-loop.yaml" {
		// Verify we got the specific error message
		found := false
		expected := "failed because I am a string, not an error object"
		for _, problem := range problems.Problems {
			if strings.Contains(problem.Message, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("did not find expected error message in problems: %s", problems.ToString(true))
		}
	}
}

//
// Context
//

type Context struct {
	tb         testing.TB
	root       string
	urlContext *exturl.Context
	parser     *parser.Parser
}

func NewContext(tb testing.TB) *Context {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("PUCCINI_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			tb.Errorf("%s", err.Error())
		}
	}

	return &Context{
		tb:         tb,
		root:       root,
		urlContext: exturl.NewContext(),
		parser:     parser.NewParser(),
	}
}

func (self *Context) compileAll() {
	self.compile("1.3/artifacts.yaml", nil)
	self.compile("1.3/attributes.yaml", nil)
	self.compile("1.3/copy.yaml", nil)
	self.compile("1.3/data-types.yaml", nil)
	self.compile("1.3/descriptions.yaml", nil)
	self.compile("1.3/dsl-definitions.yaml", nil)
	self.compile("1.3/functions.yaml", nil)
	self.compile("1.3/inputs-and-outputs.yaml", map[string]any{"ram": "1 GiB"})
	self.compile("1.3/interfaces.yaml", nil)
	self.compile("1.3/metadata.yaml", nil)
	self.compile("1.3/namespaces.yaml", nil)
	self.compile("1.3/policies-and-groups.yaml", nil)
	self.compile("1.3/requirements-and-capabilities.yaml", nil)
	self.compile("1.3/source-and-target.yaml", nil)
	self.compile("1.3/substitution-mapping-client.yaml", nil)
	self.compile("1.3/substitution-mapping.yaml", nil)
	self.compile("1.3/unicode.yaml", nil)
	self.compile("1.3/workflows.yaml", nil)

	self.compile("2.0/artifacts.yaml", nil)
	self.compile("2.0/attributes.yaml", nil)
	self.compile("2.0/copy.yaml", nil)
	self.compile("2.0/data-types.yaml", nil)
	self.compile("2.0/descriptions.yaml", nil)
	self.compile("2.0/dsl-definitions.yaml", nil)
	self.compile("2.0/functions.yaml", nil)
	self.compile("2.0/inputs-and-outputs.yaml", map[string]any{"ram": "1 GiB"})
	self.compile("2.0/interfaces.yaml", nil)
	self.compile("2.0/metadata.yaml", nil)
	self.compile("2.0/namespaces.yaml", nil)
	self.compile("2.0/policies-and-groups.yaml", nil)
	self.compile("2.0/requirements-and-capabilities.yaml", nil)
	self.compile("2.0/source-and-target.yaml", nil)
	self.compile("2.0/substitution-mapping-client.yaml", nil)
	self.compile("2.0/substitution-mapping.yaml", nil)
	self.compile("2.0/unicode.yaml", nil)
	self.compile("2.0/workflows.yaml", nil)

	self.compile("javascript/artifacts.yaml", nil)
	self.compile("javascript/constraints.yaml", nil)
	self.compile("javascript/converters.yaml", nil)
	self.compile("javascript/define.yaml", nil)
	self.compile("javascript/exec.yaml", nil)
	self.compile("javascript/functions.yaml", nil)

	self.compile("openstack/hello-world.yaml", nil)

	self.compile("bpmn/open-loop.yaml", nil)
}

func (self *Context) compile(url string, inputs map[string]any) {
	if t, ok := self.tb.(*testing.T); ok {
		t.Run(url, func(t_ *testing.T) {
			// Running the tests in parallel is not just for speed;
			// it actually helps us to find concurrency bugs
			t_.Parallel()
			self.compile_(t_, url, inputs)
		})
	} else {
		self.compile_(self.tb, url, inputs)
	}
}

func (self *Context) compile_(t testing.TB, url string, inputs map[string]any) {
	var normalServiceTemplate *normal.ServiceTemplate
	var clout *cloutpkg.Clout
	var err error

	url_ := self.urlContext.NewFileURL(path.Join(filepath.ToSlash(self.root), "examples", url))

	parserContext := self.parser.NewContext()
	parserContext.URL = url_
	parserContext.Inputs = inputs
	if normalServiceTemplate, err = parserContext.Parse(contextpkg.TODO()); err != nil {
		t.Errorf("%s\n%s", err.Error(), parserContext.GetProblems().ToString(true))
		return
	}

	problems := parserContext.GetProblems()
	if clout, err = normalServiceTemplate.Compile(); err != nil {
		t.Errorf("%s\n%s", err.Error(), problems.ToString(true))
		return
	}

	execContext := js.ExecContext{
		Clout:      clout,
		Problems:   problems,
		URLContext: self.urlContext,
		History:    true,
		Format:     "yaml",
		Pretty:     true,
	}

	execContext.Resolve()
	if !problems.Empty() {
		t.Errorf("%s", problems.ToString(true))
		return
	}

	execContext.Coerce()
	if !problems.Empty() {
		t.Errorf("%s", problems.ToString(true))
		return
	}
}
