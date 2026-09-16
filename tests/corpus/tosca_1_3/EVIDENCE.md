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

For normalized maps with string keys, assertions use a keyed evidence view:
`$map/token_type/$primitive` selects the entry whose `$key.$primitive` equals
`token_type`. It does not select an unstable array index. Duplicate keys fail;
non-string map keys retain their original representation.

`/effective/` assertions inspect resolved type fields after inheritance and
rendering; `/evaluated/` assertions compile and coerce the normalized template
through the ordinary function evaluator. Neither view reads raw fixture YAML
as the asserted result. `complete_predicate: false` deliberately prevents a
reviewed subset (for example one class in a bundled historical rule) from
being called verified. A caught mutation is sensitivity evidence, not a proof
that every variation of the subsystem was tested.

Replay reviewed patches with `python3 scripts/conformance/replay_tosca13_evidence_mutations.py`.
Do not run it concurrently with another build, test, or production edit. Each
patch requires a passing baseline, an observed failure of a mapped owned
assertion, and restoration of exact original bytes. Compilation errors and
unrelated test failures are not caught mutations. The map gate checks all
fifteen required categories, actual case/assertion ownership, and source and
fixture hashes against the replay report.
