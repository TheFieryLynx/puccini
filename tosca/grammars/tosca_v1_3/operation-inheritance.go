package tosca_v1_3

import (
	"github.com/tliron/go-puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/go-puccini/tosca/parsing"
)

// TOSCA 1.3 §3.6.17.3 and §5.8.1: supplying an implementation overrides
// the inherited artifact. A description/input-only operation supplies none.
// See docs/decisions/tosca-1.3-operation-implementation-inheritance.md.
func inheritOperationDefinition(childEntity, parentEntity parsing.EntityPtr) {
	child := childEntity.(*tosca_v2_0.OperationDefinition)
	parent := parentEntity.(*tosca_v2_0.OperationDefinition)
	if child.Implementation == nil {
		child.Implementation = parent.Implementation
	}
}
