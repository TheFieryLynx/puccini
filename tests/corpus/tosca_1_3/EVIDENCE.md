# Executable evidence contract (schema 2)

`specification.requirement_ids` and `supporting_requirements` retain historical
labels for navigation. They do not prove verification. Legacy string assertions
in `expected.assertions` are smoke checks only.

An owned assertion has exactly one primary requirement. Its predicate and exact
case/assertion binding must be independently reviewed in `evidence-predicates.yaml`.
Bindings hash the source catalog predicate and fixture bytes. This prevents
renaming a label and an assertion together to manufacture unrelated evidence;
it does not replace human normative review. New fixture bytes require review.

`value` and `normalization` assertions address the normalized model with RFC 6901
JSON pointers. `accept`/`reject` compare the accepted boolean; a reject expects
false. Diagnostics assert a stable category and entity path; phase asserts a
parser API boundary. Negative cases must also have an accepted nearest-valid
sibling and describe the minimal change. Exceptions require a reviewed primary
predicate and a conceptual justification. No exceptions are currently used.

Phases are measured after ReadRoot (including CSAR metadata), namespace/lookup,
hierarchy, inheritance, rendering (including resolved function validation), and
normalization. The harness stops at the first failing boundary. This is test-side
instrumentation, not a semantic change. A read rejection cannot prove a later
phase obligation even when its diagnostic contains a desired word. No normative
normalization error category is invented: normalized semantic values are checked.

An assertion or passing fixture alone is still insufficient for `verified`.
The recheck additionally requires reviewed positive/negative applicability,
semantic positive values, phase and targeted diagnostics, nearest-valid isolation,
and caught mutations on the relevant production paths. Historical 223 membership
is neither proof of applicability nor proof of a complete denominator. Section
and variation-axis documents are inventories, not exhaustive coverage claims.
