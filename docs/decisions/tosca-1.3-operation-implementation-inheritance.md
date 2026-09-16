# TOSCA 1.3 operation definitions without a new implementation

TOSCA 1.3 §3.6.17.1 makes implementation optional. §3.6.17.3 states
that implementation artifacts assigned in subclasses override the parent's.
§5.8.1 likewise describes overriding when a script is provided, and describes
operations without an associated implementation as no-ops.

The text does not prescribe a complete field-by-field inheritance algorithm.
Two readings of a child operation containing only description are possible:
it replaces the whole operation with a no-op, or it preserves the inherited
implementation because no new artifact was assigned.

Use the latter interpretation: descriptions and inputs do not erase an
inherited implementation. The overriding unit described by §3.6.17.3 is the
implementation artifact, not the presence of an operation map. An operation
with no implementation anywhere in its ancestry remains a no-op. A supplied
implementation overrides the parent. This decision applies only to 1.3;
it is not a claim about TOSCA 2.0 inheritance.

The evidence-hardening fixture `operation-type-inheritance.yaml` checks parent,
local override, inherited operation, local description without an override,
multi-level inheritance, and unrelated operations. The relationship template
fixtures separately exercise `OperationAssignments.CopyUnassigned`, the path
whose overwrite mutation survived the original audit.
