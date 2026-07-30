//go:build profile_discovery

package tosca_1_3_test

import (
	"fmt"
	"testing"
)

// Specification: TOSCA Simple Profile in YAML 1.3
// Section: 5.9.9.2 Definition
// Expected: the ordinary processor path loads Abstract.Storage as a Root
// subtype whose required size has the 0 MB default and lower bound.
// Category: positive, direct, normative profile, hierarchy, boundary
func TestProfileAbstractStorageDefinition(t *testing.T) {
	parserContext, release := parseNormativeProfile(t)
	defer release()
	storage := findProfileNodeType(
		t, parserContext, "tosca.nodes.Abstract.Storage")

	size := storage.PropertyDefinitions["size"]
	if size == nil {
		t.Fatal("Abstract.Storage.size is absent")
	}

	t.Run("parent", func(t *testing.T) {
		if storage.ParentName == nil || *storage.ParentName != "tosca.nodes.Root" {
			t.Fatalf("Abstract.Storage derived_from = %v, want tosca.nodes.Root",
				storage.ParentName)
		}
	})
	t.Run("required", func(t *testing.T) {
		if !size.IsRequired() {
			t.Fatalf("Abstract.Storage.size required = %v (effective %t), want true",
				size.Required, size.IsRequired())
		}
	})
	t.Run("default", func(t *testing.T) {
		if size.Default == nil || fmt.Sprint(size.Default.Context.Data) != "0 MB" {
			t.Fatalf("Abstract.Storage.size default = %#v, want 0 MB", size.Default)
		}
	})
	t.Run("constraint", func(t *testing.T) {
		if size.ValidationClause == nil ||
			size.ValidationClause.Operator != "greater_or_equal" ||
			len(size.ValidationClause.Arguments) != 1 ||
			fmt.Sprint(size.ValidationClause.Arguments[0]) != "0 MB" {
			t.Fatalf("Abstract.Storage.size constraint = %#v, want greater_or_equal 0 MB",
				size.ValidationClause)
		}
	})
}
