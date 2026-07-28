# TOSCA 1.3 frozen partial remediation roadmap

This roadmap covers only the 32 frozen atomic MUST records that had
`implementation_status: partial` at the start of the task. The frozen
denominator is 223.

## Baseline cross-status matrix

| Implementation status | verified | indirectly-tested | untested | Total |
|---|---:|---:|---:|---:|
| implemented | 141 | 26 | 24 | 191 |
| partial | 0 | 0 | 32 | 32 |
| missing | 0 | 0 | 0 | 0 |
| non-compliant | 0 | 0 | 0 | 0 |
| **Total** | **141** | **26** | **56** | **223** |

The six requested cross-tab values are therefore:

- partial + verified: 0
- partial + indirectly-tested: 0
- partial + untested: 32
- implemented + verified: 141
- implemented + indirectly-tested: 26
- implemented + untested: 24

If every scoped record obtains direct requirement-isolating evidence, the
expected final row is `implemented: verified=173, indirectly-tested=26,
untested=24, total=223`. This is an expectation derived from the baseline, not
a manually assigned target.

## Remediation order

| Order | Group | Requirements | Dependency |
|---:|---|---|---|
| 1 | version-zero-semantics | `3.3.2.5-001`, `3.3.2.5-002` | none |
| 2 | interface-reserved-operation-name | `3.7.5.4-002` | none |
| 3 | intrinsic-required-arguments | `4.3.2.2-003`, `4.7.1.2-003` | none |
| 4 | normative-name-case-sensitivity | `5.2.1-001` | none |
| 5 | inherited-required-keynames | `3.5.1-002` | none |
| 6 | constraint-semantics | `3.3.6.2-003`, `3.6.3.3-001..003`, `3.6.10.5-004`, `3.6.14.3-003`, `3.7.6.3-002` | version scalar behavior |
| 7 | datatype-shape | `3.7.6.3-001`, `3.7.6.3-003` | constraint type checks |
| 8 | capability-source-refinement | `3.7.2.4-001` | completed hierarchy |
| 9 | group-member-homogeneity | `3.7.11.4-002` | completed hierarchy |
| 10 | requirement-node-filter | `3.8.2.2.3-015` | hierarchy and assignment resolution |
| 11 | attribute-default-provenance | `3.6.12.4-002` | intrinsic syntax and resolution |
| 12 | workflow-operation-host | `3.6.27.1-008`, `3.6.27.1-010` | target and relationship resolution |
| 13 | template-copy-depth | `3.8.3.3-001`, `3.8.4.3-001` | none |
| 14 | substitution-coverage | `3.8.8.3-005`, `3.8.13.4-001` | rendered assignments and effective node type |
| 15 | profile-cross-property | `5.3.11.3-001..003`, `5.5.7.4-001`, `5.5.13.1-012`, `8.5.1.1-036` | rendered complex values and defaults |

The groups are separated where parser phase, changed symbols, or regression
risk differs. In particular, grammar-only interface/function rules are not
mixed with hierarchy checks; substitution mapping coverage is not mixed with
ordinary assignment validation; and normative profile cross-property rules
remain isolated from generic constraint processing.

## Current cross-status matrix

After remediation group 6 (`constraint-semantics`):

| Implementation status | verified | indirectly-tested | untested | Total |
|---|---:|---:|---:|---:|
| implemented | 155 | 26 | 24 | 205 |
| partial | 0 | 0 | 18 | 18 |
| missing | 0 | 0 | 0 | 0 |
| non-compliant | 0 | 0 | 0 | 0 |
| **Total** | **155** | **26** | **42** | **223** |
