---
name: tosca-conformance
description: Implement and audit TOSCA 1.3 and TOSCA 2.0 compliance using the pinned OASIS specifications, version-specific tests, and conformance tracking.
---

# TOSCA conformance workflow

Use this skill for any work involving TOSCA parsing, validation, inheritance,
imports, namespaces, intrinsic functions, requirements, capabilities,
interfaces, artifacts, substitution mappings, workflows, CSAR processing,
normalization, or standards-compliance review.

## Normative sources

Use only these local specifications as normative sources:

- TOSCA 1.3:
  `docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html`
- TOSCA 2.0:
  `docs/specifications/tosca/2.0/TOSCA-v2.0-os.md`

Do not use model memory, third-party implementations, examples, blog posts, or
existing Puccini behavior as proof of standards compliance.

## Version selection

Before changing code:

1. Determine the applicable TOSCA version.
2. Read `tosca_definitions_version` when analyzing a service template.
3. Use only the corresponding specification.
4. Do not apply TOSCA 2.0 rules to TOSCA 1.3.
5. Do not apply TOSCA 1.3 rules to TOSCA 2.0.
6. Share implementation only when semantics are confirmed identical in both
   specifications.

Recognized versions:

- `tosca_simple_yaml_1_3`
- `tosca_2_0`

## Required workflow

For every standards-related task:

1. Find the exact relevant specification sections.
2. Record section numbers and titles.
3. Read surrounding normative text and cross-references.
4. Classify each rule as required, prohibited, recommended, optional,
   deprecated, informative, unspecified, or ambiguous.
5. Inspect the current implementation and tests.
6. Identify the responsible parser phase.
7. Add positive and negative conformance tests.
8. Implement the smallest standards-compliant change.
9. Run focused and full tests.
10. Update the conformance matrix.
11. Report cross-version impact.

Do not change production behavior before identifying the normative basis.

## Parser phases

Use the earliest phase where all required information is available:

- Structural YAML and grammar checks: read phase.
- Namespace and name resolution: namespace phases.
- Type hierarchy validation: hierarchy phases.
- Inheritance and refinement: inheritance phases.
- Assignments, defaults, and compatibility: rendering phases.
- Version-neutral representation: normalization.

Do not move validation to a later phase only because it is easier.

## Tests

Keep conformance tests separate by version:

- `tests/conformance/tosca_1_3/`
- `tests/conformance/tosca_2_0/`

Each test must identify:

- specification version;
- section number;
- section title;
- expected result;
- test category.

Add applicable tests for:

- minimal valid input;
- complete valid input;
- missing required keynames;
- unknown keynames;
- incorrect YAML types;
- invalid short notation;
- invalid long notation;
- inheritance and refinement;
- imports and namespaces;
- boundary values;
- intrinsic function arguments;
- cross-version rejection;
- normalization;
- regressions.

A feature is not implemented merely because one valid example parses.

Negative tests must fail for the intended reason.

## Shared and version-specific code

Use shared code only for behavior proven semantically identical in TOSCA 1.3
and TOSCA 2.0.

Use version-specific code when grammar, defaults, refinement rules, imports,
namespaces, functions, assignments, substitution mappings, workflows, CSAR
behavior, or normalization semantics differ.

When shared code changes, add regression coverage for both versions.

## Ambiguities

When specification text is ambiguous or contradictory:

1. Collect all relevant sections.
2. Describe competing interpretations.
3. Distinguish normative text from implementation decisions.
4. Choose the most internally consistent interpretation.
5. Record it in `docs/decisions/`.
6. Add tests that lock the selected interpretation.

Do not present an implementation interpretation as an explicit specification
requirement.

## Quirks and extensions

Existing Puccini behavior is not automatically normative.

Non-standard behavior must be:

- explicitly documented;
- isolated behind a compatibility or quirk mode;
- tested separately from strict mode;
- excluded from strict standards-compliant behavior.

Always distinguish between:

- TOSCA processor requirements;
- orchestrator behavior;
- Puccini normalization;
- deployment behavior;
- project extensions;
- compatibility quirks.

## Verification

Run available checks, including:

```bash
go test ./...
go test -race ./...
go vet ./...
```

Do not claim a command passed unless it was actually executed successfully.

## Required final report

Finish standards-related work with:

### Applicable specification

State TOSCA 1.3, TOSCA 2.0, or both independently.

### Normative basis

List section numbers, titles, and applied normative requirements.

### Implementation

List changed files and affected parser phases.

### Tests

List positive, negative, boundary, cross-version, and regression tests.

### Cross-version impact

Describe TOSCA 1.3 and TOSCA 2.0 impact separately.

### Ambiguities

List interpretations, unresolved questions, and decision records.

### Conformance status

State one of:

- implemented;
- partial;
- non-compliant;
- unimplemented;
- ambiguous;
- blocked;
- not-applicable.
