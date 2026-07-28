package parser

import (
	"sort"

	"github.com/tliron/go-kutil/terminal"
	"github.com/tliron/go-puccini/tosca/parsing"
)

func (self *Context) AddNamespaces() {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	self.Root.mergeNamespaces()
	if validator, ok := self.Root.EntityPtr.(parsing.NamespaceValidator); ok {
		validator.ValidateNamespace()
	}
}

func (self *File) mergeNamespaces() {
	context := self.GetContext()

	self.importsLock.RLock()
	imports := append(Files(nil), self.Imports...)
	self.importsLock.RUnlock()
	sort.Sort(imports)

	for _, import_ := range imports {
		import_.mergeNamespaces()
		context.Namespace.Merge(import_.GetContext().Namespace, import_.NameTransformer)
		context.ScriptletNamespace.Merge(import_.GetContext().ScriptletNamespace)
	}

	logNamespaces.Debugf("create: %s", context.URL.String())
	namespace := parsing.NewNamespaceFor(self.EntityPtr)
	context.Namespace.Merge(namespace, nil)
}

// Print

func (self *Context) PrintNamespaces(indent int) {
	self.filesLock.RLock()
	defer self.filesLock.RUnlock()

	childIndent := indent + 1
	for _, file := range self.Files {
		context := file.GetContext()
		if !context.Namespace.Empty() {
			terminal.PrintIndent(indent)
			terminal.Printf("%s\n", self.Stylist.Value(context.URL.String()))
			context.Namespace.Print(childIndent)
		}
	}
}
