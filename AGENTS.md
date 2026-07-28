# go-puccini specification compliance rules

## Project objective

The objective of this fork is to provide a complete, precise, maintainable,
and standards-compliant implementation of:

* TOSCA Simple Profile in YAML Version 1.3
* TOSCA Version 2.0

Standards compliance has priority over:

* compatibility with undocumented historical Puccini behavior;
* assumptions made by other TOSCA implementations;
* convenient but non-standard behavior;
* existing code when that code contradicts the applicable specification.

Backward compatibility may be preserved through an explicitly documented
compatibility or quirk mode, but must never silently change strict standards
behavior.

## Normative sources

The authoritative local sources are:

### TOSCA 1.3

`docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html`

This specification governs files declaring:

`tosca_definitions_version: tosca_simple_yaml_1_3`

### TOSCA 2.0

`docs/specifications/tosca/2.0/TOSCA-v2.0-os.md`

This specification governs files declaring:

`tosca_definitions_version: tosca_2_0`

The local files are pinned copies of the published OASIS Standard versions.
Do not use drafts, Committee Specification versions, blog posts, examples from
third-party implementations, or model memory as substitutes for these sources.

## Version isolation

TOSCA 1.3 and TOSCA 2.0 are separate language versions.

Do not treat TOSCA 2.0 as an incremental extension that can be applied to
TOSCA 1.3 files.

Before analyzing a service template:

1. Read `tosca_definitions_version`.
2. Select the corresponding grammar and semantic rules.
3. Apply only the selected specification.
4. Reject syntax or semantics that belong exclusively to another version,
   unless an explicit compatibility mode requires otherwise.

The implementation may share internal infrastructure, algorithms, utilities,
and normalized data structures, but version-specific grammar and semantics
must remain explicitly distinguishable.

Do not implement a behavior in shared code when doing so would silently alter
the behavior of another TOSCA version.

## Documentation-first workflow

Before changing parser, grammar, validation, inheritance, rendering,
normalization, CSAR, function, namespace, type-system, or assignment behavior:

1. Identify the exact applicable specification version.
2. Find the relevant section or sections in the local specification.
3. Read the surrounding normative text, including definitions, requirements,
   constraints, notes, and referenced sections.
4. Determine whether the requested behavior is:

   * required;
   * prohibited;
   * optional;
   * recommended;
   * deprecated;
   * informative;
   * unspecified;
   * ambiguous or contradictory.
5. Inspect the current implementation and existing tests.
6. Write or update conformance tests.
7. Implement the smallest standards-compliant change.
8. Run the relevant test suite.
9. Update the conformance matrix.
10. Report the specification sections used.

Do not modify production behavior before establishing the applicable
specification requirement and expected test result.

## Normative language

Interpret MUST, MUST NOT, SHALL, SHALL NOT, REQUIRED, SHOULD, SHOULD NOT,
RECOMMENDED, MAY, and OPTIONAL according to their normative meaning.

Examples, tutorials, explanatory diagrams, and sections explicitly marked as
examples are informative unless the specification states otherwise.

An example is not sufficient proof that syntax is mandatory, prohibited, or
exhaustive.

## Required evidence

Every standards-related production change must include:

* applicable TOSCA version;
* specification section number and title;
* normative statement or rule being implemented;
* at least one positive test when valid behavior is involved;
* at least one negative test when invalid behavior must be rejected;
* expected parser phase in which validation occurs;
* conformance matrix update.

When useful, also include:

* boundary-value tests;
* short-notation and long-notation equivalents;
* inheritance tests;
* import and namespace tests;
* cross-version rejection tests;
* normalization equivalence tests;
* CSAR tests.

## Frozen TOSCA 1.3 conformance baseline

The TOSCA 1.3 atomic MUST denominator is frozen at 223 requirements.

Implementation changes, bug fixes, refactoring, and conformance-test additions
must not change:

* requirement IDs;
* applicability;
* catalog quality;
* deduplication;
* denominator membership;
* the total atomic MUST denominator of 223.

If an error is discovered in the requirements catalog, applicability
classification, deduplication, or denominator membership:

1. Do not modify it as part of an implementation or conformance-test task.
2. Record the proposed correction in:
   `docs/conformance/tosca-1.3/catalog-review-queue.yaml`
3. Describe:

   * the affected requirement IDs;
   * the current classification;
   * the proposed classification;
   * the normative justification;
   * the expected denominator change.
4. Apply the correction only in a separate catalog-review change.
5. Regenerate and independently validate all conformance artifacts after the
   catalog-review change.

A coverage improvement must result from implementing or verifying requirements,
not from removing requirements from the denominator.

Every implementation or conformance-test task must verify that:

* the atomic MUST denominator remains exactly 223;
* no requirement IDs were added, removed, or reclassified;
* the sum of all applicable atomic MUST implementation statuses remains 223.

## Conformance test organization

Version-specific tests must be stored separately:

* `tests/conformance/tosca_1_3/`
* `tests/conformance/tosca_2_0/`

Each conformance test must identify its normative source using adjacent
metadata or comments containing:

* specification version;
* section number;
* section title;
* expected result;
* test category.

A TOSCA 1.3 test must not be reused as proof of TOSCA 2.0 compliance unless the
same requirement has been independently verified in both specifications.

Tests derived from community repositories or third-party implementations are
supplementary evidence only. Their expectations must be verified against the
applicable OASIS specification before acceptance.

## Positive and negative validation

For every supported grammar construct, verify both:

1. valid documents are accepted;
2. invalid documents are rejected at the appropriate processing phase.

Do not consider a feature implemented merely because one example parses.

Validation should cover, where applicable:

* missing required keynames;
* unknown keynames;
* incorrect YAML types;
* invalid short notation;
* invalid long notation;
* invalid names and namespace references;
* incompatible type assignments;
* inheritance contract violations;
* invalid occurrence or count ranges;
* invalid intrinsic function arguments;
* invalid substitution mappings;
* unsupported cross-version syntax.

## Parser phases

Preserve clear responsibilities between parser phases.

In general:

* YAML decoding and structural grammar validation belong in the read phase.
* Namespace construction and name resolution belong in namespace phases.
* cyclic and incomplete type hierarchy validation belongs in the hierarchy
  phase.
* inheritance and refinement contract checks belong in the inheritance phase.
* assignment, default application, and template/type compatibility checks
  belong in the rendering phase.
* conversion to version-neutral output belongs in normalization.

Do not move validation into a later phase solely because it is easier to
implement there when the error can be detected correctly and unambiguously in
an earlier phase.

When a validation rule necessarily depends on resolved types, inherited
definitions, imports, or assignments, perform it in the earliest phase where
all required information is available.

## Shared code and version-specific code

Prefer shared code only for behavior that is semantically identical across
TOSCA versions.

Use version-specific code when any of the following differ:

* recognized keynames;
* grammar shape;
* default values;
* short notation;
* type refinement rules;
* import behavior;
* namespace behavior;
* intrinsic functions;
* requirement assignment;
* substitution mappings;
* workflows;
* CSAR layout or metadata;
* normalization semantics.

Do not add checks such as `if version == ...` throughout unrelated generic
code when a version-specific grammar implementation, interface, policy object,
or strategy can isolate the difference more clearly.

## Normalization

Normalization is an internal Puccini representation and is not itself defined
by the TOSCA specifications.

Do not claim that normalized output is standards-compliant merely because the
input parser is compliant.

When both TOSCA versions describe equivalent concepts, they may normalize to a
common representation.

When the versions contain semantically distinct concepts, the normalized model
must retain enough information to avoid semantic loss.

A normalization change must include tests proving that:

* required source semantics are preserved;
* version-specific distinctions are not silently discarded;
* equivalent short and long notation normalize consistently;
* invalid source input is rejected before normalization.

## Ambiguity and contradictions

Do not silently invent behavior when the applicable specification is unclear,
incomplete, or contradictory.

When ambiguity is found:

1. document all relevant specification sections;
2. describe the competing interpretations;
3. inspect related definitions and normative terminology;
4. choose the interpretation that is most internally consistent with the
   applicable specification;
5. record the decision in `docs/decisions/`;
6. add tests that lock the selected interpretation;
7. keep alternative historical behavior behind an explicit quirk mode when
   compatibility is valuable.

A decision record must not present an interpretation as an explicit
specification requirement when it is only an implementation decision.

## Existing quirks

Existing Puccini quirk behavior is not normative.

For every existing or new quirk:

* document the behavior;
* document why strict behavior differs;
* cite the relevant specification sections;
* ensure strict mode follows the selected interpretation of the specification;
* test strict and quirk behavior separately;
* do not enable non-standard behavior silently.

## Scope distinction

Always distinguish between:

* TOSCA processor requirements;
* TOSCA orchestrator behavior;
* Puccini normalization behavior;
* deployment-engine behavior;
* project-specific extensions.

Puccini is primarily a TOSCA processor. Do not implement or infer orchestrator
runtime behavior as parser semantics unless the applicable specification
explicitly requires the processor to validate or represent it.

## Completion criteria

A standards-related task is complete only when:

* the applicable specification sections were identified;
* positive and negative tests exist;
* tests pass;
* behavior is isolated by TOSCA version where required;
* no unrelated version behavior regressed;
* the conformance matrix was updated;
* ambiguities and implementation decisions were documented;
* the final response includes a standards compliance summary.

## Required final report

For every standards-related task, finish with:

### Applicable specification

State TOSCA 1.3, TOSCA 2.0, or both.

### Normative basis

List section numbers and section titles.

### Implementation

List the changed parser phases and files.

### Tests

List positive, negative, boundary, and regression tests added or changed.

### Cross-version impact

State whether TOSCA 1.3 and TOSCA 2.0 behavior changed.

### Ambiguities

List unresolved or interpreted specification language.

### Conformance status

State whether the relevant conformance matrix entries are implemented, partial,
blocked, ambiguous, or not applicable.

