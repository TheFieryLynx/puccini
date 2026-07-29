# TOSCA 1.3 direct verification roadmap

This roadmap freezes the 50 applicable atomic MUST records that were implemented
but lacked direct requirement-isolating evidence at task start. The denominator is 223.

## Baseline matrix

| Implementation status | verified | indirectly-tested | untested | Total |
|---|---:|---:|---:|---:|
| implemented | 173 | 26 | 24 | 223 |
| partial | 0 | 0 | 0 | 0 |
| missing | 0 | 0 | 0 | 0 |
| non-compliant | 0 | 0 | 0 | 0 |
| **Total** | **173** | **26** | **24** | **223** |

Frozen verification backlog: **50** (26 indirectly-tested, 24 untested).

## Verification order

| Order | Group | Requirements | Count | Parser phase / execution path |
|---:|---|---|---:|---|
| 1 | `scalar-semantics` | `TOSCA13-3.3.2.1-004`, `TOSCA13-3.3.2.1-005`, `TOSCA13-3.3.3.1-004`, `TOSCA13-3.3.3.1-005`, `TOSCA13-3.3.3.1-006`, `TOSCA13-3.3.6.1-004`, `TOSCA13-3.3.6.1-005`, `TOSCA13-3.3.6.1-006`, `TOSCA13-3.3.6.2-001`, `TOSCA13-3.3.6.2-002` | 10 | scalar value rendering |
| 2 | `map-entry-grammar` | `TOSCA13-3.3.5.1.2-002` | 1 | read |
| 3 | `definition-defaults` | `TOSCA13-3.6.10.4-008`, `TOSCA13-3.6.10.5-003`, `TOSCA13-3.6.12.2-010`, `TOSCA13-3.6.14.2-015`, `TOSCA13-3.6.14.3-002` | 5 | rendering after type resolution |
| 4 | `override-workflow-occurrences` | `TOSCA13-3.6.17.3-002`, `TOSCA13-3.6.25.2.1-002`, `TOSCA13-3.7.3.3-002`, `TOSCA13-5.8.1-004` | 4 | inheritance and rendering |
| 5 | `normative-profile-fields` | `TOSCA13-5.3.11.1-002`, `TOSCA13-5.3.11.2-005`, `TOSCA13-5.3.6.1-005`, `TOSCA13-5.3.6.1-008`, `TOSCA13-5.3.7.1-002`, `TOSCA13-5.3.7.1-005`, `TOSCA13-5.3.7.2-005`, `TOSCA13-5.3.7.2-008`, `TOSCA13-5.5.13.1-003`, `TOSCA13-5.5.13.1-007`, `TOSCA13-5.5.7.1-003`, `TOSCA13-5.5.7.2-002`, `TOSCA13-5.5.7.3-005`, `TOSCA13-5.7.1.1-002`, `TOSCA13-5.7.1.1-005`, `TOSCA13-5.7.1.1-010`, `TOSCA13-5.7.5.1-003`, `TOSCA13-5.9.1.2-002`, `TOSCA13-5.9.1.2-005`, `TOSCA13-5.9.1.2-010`, `TOSCA13-5.9.5.3-001`, `TOSCA13-5.9.8.1-002`, `TOSCA13-5.9.9.1-002` | 23 | implicit profile read, hierarchy, and inheritance |
| 6 | `normative-connectivity` | `TOSCA13-8.2-002`, `TOSCA13-8.3.2-001` | 2 | implicit profile read and hierarchy |
| 7 | `csar-entry-processing` | `TOSCA13-6.1-005`, `TOSCA13-6.1-008`, `TOSCA13-6.2-001`, `TOSCA13-6.3-002`, `TOSCA13-6.3-005` | 5 | CSAR metadata/root selection and service-template read |

The groups are independent where parser phase, grammar entity, production symbols,
fixture architecture, or regression risk differs. TOSCA 2.0 tests are regression
protection only and never TOSCA 1.3 normative evidence.

## Current matrix

After verification group 1 (`scalar-semantics`), all ten scalar records have
direct positive, negative, boundary, diagnostic, and normalization evidence.
Verified increases monotonically while the denominator and implementation axis
remain fixed.

| Implementation status | verified | indirectly-tested | untested | Total |
|---|---:|---:|---:|---:|
| implemented | 183 | 26 | 14 | 223 |
| partial | 0 | 0 | 0 | 0 |
| missing | 0 | 0 | 0 | 0 |
| non-compliant | 0 | 0 | 0 | 0 |
| **Total** | **183** | **26** | **14** | **223** |
