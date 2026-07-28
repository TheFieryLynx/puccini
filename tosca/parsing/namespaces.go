package parsing

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tliron/go-kutil/reflection"
	"github.com/tliron/go-kutil/terminal"
)

type NameTransformer = func(string, EntityPtr) []string

// NamespaceValidator is an optional version-neutral phase hook. Grammar
// entities implement it only when their specification defines additional
// rules that require the fully merged namespace.
type NamespaceValidator interface {
	ValidateNamespace()
}

func GetCanonicalName(entityPtr EntityPtr) string {
	if metadata, ok := GetMetadata(entityPtr); ok {
		if canonicalName, ok := metadata[MetadataCanonicalName]; ok {
			return canonicalName
		}
	}

	context := GetContext(entityPtr)
	canonicalNamespace := context.GetCanonicalNamespace()
	if canonicalNamespace != nil {
		return fmt.Sprintf("%s::%s", *canonicalNamespace, context.Name)
	} else {
		return context.Name
	}
}

//
// Namespace
//

type Namespace struct {
	namespace map[reflect.Type]map[string]EntityPtr
}

func NewNamespace() *Namespace {
	return &Namespace{
		namespace: make(map[reflect.Type]map[string]EntityPtr),
	}
}

// From "namespace" tags
func NewNamespaceFor(entityPtr EntityPtr) *Namespace {
	self := NewNamespace()

	reflection.TraverseEntities(entityPtr, false, func(entityPtr EntityPtr) bool {
		for _, field := range reflection.GetTaggedFields(entityPtr, "namespace") {
			if field.Kind() != reflect.String {
				panic(fmt.Sprintf("\"namespace\" tag can only be used on \"string\" field in struct: %T", entityPtr))
			}

			name := field.String()

			if context := GetContext(entityPtr); context != nil {
				if context.HasQuirk(QuirkNamespaceNormativeIgnore) {
					// Do not add normative types to the namespace
					if metadata, ok := GetMetadata(entityPtr); ok {
						if normative, ok := metadata[MetadataNormative]; ok {
							if normative == "true" {
								continue
							}
						}
					}
				}

				// Check for invalid characters
				if context.Grammar.InvalidNamespaceCharacters != "" {
					if strings.Contains(name, context.Grammar.InvalidNamespaceCharacters) {
						context.ReportNameInvalid(entityPtr, name)
					}
				}
			}

			self.set(name, entityPtr)
		}
		return true
	})

	return self
}

func (self *Namespace) Empty() bool {
	return len(self.namespace) == 0
}

func (self *Namespace) Range(f func(EntityPtr) bool) {
	for _, forType := range self.namespace {
		for _, entityPtr := range forType {
			if !f(entityPtr) {
				return
			}
		}
	}
}

func (self *Namespace) Lookup(name string) (EntityPtr, bool) {
	for _, forType := range self.namespace {
		if entityPtr, ok := forType[name]; ok {
			return entityPtr, true
		}
	}

	return nil, false
}

func (self *Namespace) LookupForType(name string, type_ reflect.Type) (EntityPtr, bool) {
	if forType, ok := self.namespace[type_]; ok {
		entityPtr, ok := forType[name]
		return entityPtr, ok
	} else {
		return nil, false
	}
}

// If the name has already been set returns existing entityPtr, true
func (self *Namespace) Set(name string, entityPtr EntityPtr) (EntityPtr, bool) {
	return self.set(name, entityPtr)
}

func (self *Namespace) Merge(namespace *Namespace, nameTransformer NameTransformer) {
	if self == namespace {
		return
	}

	type entry struct {
		type_     reflect.Type
		name      string
		entityPtr EntityPtr
	}

	var entries []entry
	for type_, forType := range namespace.namespace {
		for name, entityPtr := range forType {
			entries = append(entries, entry{type_, name, entityPtr})
		}
	}
	sort.Slice(entries, func(i int, j int) bool {
		iType := entries[i].type_.String()
		jType := entries[j].type_.String()
		if iType != jType {
			return iType < jType
		}
		if entries[i].name != entries[j].name {
			return entries[i].name < entries[j].name
		}
		return GetContext(entries[i].entityPtr).URL.String() < GetContext(entries[j].entityPtr).URL.String()
	})

	for _, entry := range entries {
		var names []string

		if nameTransformer != nil {
			names = append(names, nameTransformer(entry.name, entry.entityPtr)...)
		} else {
			names = []string{entry.name}
		}

		for _, name := range names {
			if existing, exists := self.set(name, entry.entityPtr); exists {
				GetContext(entry.entityPtr).ReportNameAmbiguous(entry.type_.Elem(), name, entry.entityPtr, existing)
			}
		}
	}
}

func (self *Namespace) set(name string, entityPtr EntityPtr) (EntityPtr, bool) {
	type_ := reflect.TypeOf(entityPtr)
	forType, ok := self.namespace[type_]
	if !ok {
		forType = make(map[string]EntityPtr)
		self.namespace[type_] = forType
	}

	if existing, ok := forType[name]; ok {
		if existing != entityPtr {
			// We are trying to give the name to a different entity
			return existing, true
		}

		// We already have this entity at this name, and that's fine
		return nil, false
	}

	forType[name] = entityPtr

	return nil, false
}

// Print

func (self *Namespace) Print(indent int) {
	// Sort type names
	var types TypesByName
	for type_ := range self.namespace {
		types = append(types, type_)
	}
	sort.Sort(types)

	nameIndent := indent + 1
	for _, type_ := range types {
		forType := self.namespace[type_]
		terminal.PrintIndent(indent)
		terminal.Printf("%s\n", terminal.StdoutStylist.TypeName(type_.Elem().String()))

		// Sort names
		names := make([]string, len(forType))
		i := 0
		for name := range forType {
			names[i] = name
			i += 1
		}
		sort.Strings(names)

		for _, name := range names {
			terminal.PrintIndent(nameIndent)
			terminal.Printf("%s\n", name)
		}
	}
}

//
// TypesByName
//

type TypesByName []reflect.Type

// ([sort.Interface])
func (self TypesByName) Len() int {
	return len(self)
}

// ([sort.Interface])
func (self TypesByName) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// ([sort.Interface])
func (self TypesByName) Less(i, j int) bool {
	iName := self[i].Elem().String()
	jName := self[j].Elem().String()
	return strings.Compare(iName, jName) < 0
}
