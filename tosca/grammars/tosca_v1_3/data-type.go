package tosca_v1_3

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// Embedding: inherits all methods from the v2.0 struct
type DataType struct {
	*tosca_v2_0.DataType
	// Optional: if old code accesses the field directly
	ConstraintClauses ConstraintClauses `json:"-" yaml:"-"`
}

// ([parsing.Reader] signature)
func ReadDataType(ctx *parsing.Context) parsing.EntityPtr {
	var (
		hasDerivedFrom bool
		hasProperties  bool
		properties     ard.Value
	)

	// Transform YAML 1.x: constraints → validation
	if m, ok := ctx.Data.(ard.Map); ok {
		_, hasDerivedFrom = m["derived_from"]
		properties, hasProperties = m["properties"]

		if c, ok := m["constraints"].(ard.List); ok && len(c) > 0 {
			c = normalizeConstraintList(c)
			if len(c) == 1 {
				m["validation"] = c[0]
			} else {
				m["validation"] = ard.Map{"$and": c}
			}
			delete(m, "constraints")
		}
	}

	// Disable TOSCA 2.0 specific fields that require UnitsReader
	ctx.SetReadTag("Units", "")         // disable Units field
	ctx.SetReadTag("CanonicalUnit", "") // disable CanonicalUnit field
	ctx.SetReadTag("Prefixes", "")      // disable Prefixes field
	ctx.SetReadTag("DataTypeName", "")  // disable DataTypeName field

	// Call TOSCA 2.0 reader
	v2dt := tosca_v2_0.ReadDataType(ctx).(*tosca_v2_0.DataType)

	// TOSCA 1.3 §3.7.6.3: a data type must declare a valid parent or at
	// least one property, and an explicitly supplied properties map cannot
	// be empty. Structural type errors and parent validity remain the
	// responsibility of the shared field reader and lookup phase.
	if propertyMap, ok := properties.(ard.Map); hasProperties && ok && len(propertyMap) == 0 {
		ctx.FieldChild("properties", properties).
			ReportValueMalformed("properties", "must contain one or more property definitions")
	}
	if !hasDerivedFrom && !hasProperties && !isBundledDataType(ctx) {
		ctx.ReportValueMalformed(
			"data type definition",
			"requires derived_from or at least one property definition",
		)
	}

	// Return the v2.0 entity directly
	return v2dt
}

func isBundledDataType(ctx *parsing.Context) bool {
	if (ctx.URL == nil) || !strings.HasPrefix(ctx.URL.String(), "internal:/profiles/") {
		return false
	}
	if ctx.Name == "tosca.datatypes.Root" {
		return true
	}
	if data, ok := ctx.Data.(ard.Map); ok {
		if metadata, ok := data["metadata"].(ard.Map); ok {
			_, isPrimitive := metadata["puccini.type"]
			return isPrimitive
		}
	}
	return false
}
