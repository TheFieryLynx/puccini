# TOSCA 1.3 atomic MUST trace report

This report covers the continuation trace of every formerly-unknown record that remains in the stable applicable atomic MUST denominator.

## Stable scope

- Previous denominator: **288**.
- Explicit denominator changes after catalog/applicability review: **65**.
- Stable applicable atomic MUST denominator: **223**.
- Originally unknown records reviewed: **268**.
- Records excluded/reclassified before implementation tracing: **65**.
- Formerly unknown records fully traced in the stable denominator: **203**.
- Remaining applicable atomic MUST implementation_status=unknown: **0**.

## Stable denominator result

- implemented: **222**.
- partial: **1**.
- missing: **0**.
- non-compliant: **0**.
- verified: **172**.
- indirectly-tested: **26**.
- untested: **25**.
- unverifiable: **0**.
- implementation coverage: **222/223 (99.55%)**.
- verified conformance coverage: **172/223 (77.13%)**.
- broad test evidence coverage: **198/223 (88.79%)**.

## Continuation trace by subsystem

| Subsystem | Reviewed | implemented | partial | missing | non-compliant | verified | indirectly-tested | untested | blocked |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Service Template and tosca_definitions_version | 1 | 1 | 0 | 0 | 0 | 1 | 0 | 0 | 0 |
| Grammar and structural validation | 35 | 35 | 0 | 0 | 0 | 25 | 0 | 10 | 0 |
| Imports and namespaces | 10 | 10 | 0 | 0 | 0 | 10 | 0 | 0 | 0 |
| Type system and hierarchy | 6 | 6 | 0 | 0 | 0 | 6 | 0 | 0 | 0 |
| Inheritance and refinement | 3 | 3 | 0 | 0 | 0 | 1 | 0 | 2 | 0 |
| Properties, attributes and assignments | 20 | 20 | 0 | 0 | 0 | 15 | 1 | 4 | 0 |
| Constraints | 8 | 8 | 0 | 0 | 0 | 8 | 0 | 0 | 0 |
| Intrinsic functions | 22 | 22 | 0 | 0 | 0 | 22 | 0 | 0 | 0 |
| Requirements and capabilities | 6 | 6 | 0 | 0 | 0 | 5 | 0 | 1 | 0 |
| Interfaces, operations and artifacts | 11 | 11 | 0 | 0 | 0 | 11 | 0 | 0 | 0 |
| Groups and policies | 5 | 5 | 0 | 0 | 0 | 5 | 0 | 0 | 0 |
| Substitution mappings | 2 | 2 | 0 | 0 | 0 | 2 | 0 | 0 | 0 |
| Workflows | 15 | 15 | 0 | 0 | 0 | 14 | 0 | 1 | 0 |
| Normative profile types | 29 | 28 | 1 | 0 | 0 | 7 | 21 | 1 | 0 |
| CSAR | 6 | 6 | 0 | 0 | 0 | 1 | 0 | 5 | 0 |
| Normalization | 1 | 1 | 0 | 0 | 0 | 1 | 0 | 0 | 0 |
| Conformance and error requirements | 23 | 23 | 0 | 0 | 0 | 22 | 0 | 1 | 0 |

## Denominator changes

Each change below removes a record only because its catalog quality or processor applicability was re-established from the pinned source. No implementation result was used to shrink the denominator.

| Requirement | Change | Reason |
|---|---|---|
| `TOSCA13-14.3-002` | catalog_quality=conformance-umbrella | Restates all lower-level parsing obligations. |
| `TOSCA13-14.3-003` | catalog_quality=conformance-umbrella | Restates all lower-level recognition obligations. |
| `TOSCA13-14.3-004` | catalog_quality=conformance-umbrella | Restates all lower-level rejection obligations. |
| `TOSCA13-14.3-008` | catalog_quality=conformance-umbrella | Restates all atomic import-resolution obligations. |
| `TOSCA13-14.3-009` | catalog_quality=conformance-umbrella | Restates the atomic errors in §3.1. |
| `TOSCA13-14.3-010` | catalog_quality=conformance-umbrella | Restates the atomic errors in §3.2. |
| `TOSCA13-14.3-011` | catalog_quality=conformance-umbrella | Restates the atomic errors in §3.6. |
| `TOSCA13-3.1.3.1-006` | catalog_quality=dependent-summary | Heading for the independently testable duplicate-name scopes in records 007–013. |
| `TOSCA13-3.1.3.1-014` | catalog_quality=dependent-summary | Heading for the independently testable duplicate-template scopes in records 015–018. |
| `TOSCA13-3.1.3.1-019` | catalog_quality=dependent-summary | Heading for the independently testable nested-name scopes in records 020–027. |
| `TOSCA13-3.3.1-001` | catalog_quality=author-only | The SHALL governs names chosen by service-template/type authors; it does not impose an independent processor algorithm. |
| `TOSCA13-3.3.6.2-004` | catalog_quality=example | The generated record is the example that follows the scalar-unit comparison rule, not an additional normative obligation. |
| `TOSCA13-3.5.1-001` | catalog_quality=dependent-summary | General introduction to required keynames; concrete required-key rules are recorded separately. |
| `TOSCA13-3.6.10.2-008` | catalog_quality=informative | Field description says the key is optional; REQUIRED was extracted from the property name, not normative language. |
| `TOSCA13-3.6.14.1-001` | catalog_quality=informative | Cross-reference contrasting property and parameter definitions; it is not a parameter requirement. |
| `TOSCA13-3.6.14.1-002` | catalog_quality=misclassified | The prose calls the represented data type 'required' while the adjacent normative note explicitly makes parameter type optional. |
| `TOSCA13-3.6.14.2-010` | catalog_quality=informative | The note explicitly states that parameter type is not required; keyword extraction produced the opposite strength. |
| `TOSCA13-3.6.16.2.4-001` | catalog_quality=informative | A permission to use an extended notation, not an independently testable MUST-level condition. |
| `TOSCA13-3.6.21.2-004` | catalog_quality=dependent-grammar-fragment | One alternative value of the required node selector; it cannot be tested independently from the node key and the template alternative. |
| `TOSCA13-3.6.21.2-005` | catalog_quality=dependent-grammar-fragment | One alternative value of the required node selector; it cannot be tested independently from the node key and the type alternative. |
| `TOSCA13-3.7.1.3-001` | catalog_quality=metamodel-declaration | Declares the abstract TOSCA Entity Type in the metamodel; it is not an independently addressable YAML processor algorithm. |
| `TOSCA13-3.7.1.3-002` | catalog_quality=author-only | Prohibits profiles from creating new top-level base type families; ordinary service templates cannot derive directly from the abstract metamodel Entity Type. |
| `TOSCA13-3.7.11-001` | catalog_quality=informative | Conceptual explanation of group membership, not a processor conformance rule. |
| `TOSCA13-3.8.2.2.3-008` | catalog_quality=example | Generated from an explanatory list of extended-grammar uses, not a separate assignment obligation. |
| `TOSCA13-3.9.2.3-002` | catalog_quality=informative | Historical note that relationship_templates is optional; REQUIRED refers to TOSCA 1.0. |
| `TOSCA13-5-001` | catalog_quality=conformance-umbrella | Section-wide summary; every normative profile declaration must be compared atomically instead. |
| `TOSCA13-5.3.8.4-002` | catalog_quality=informative | States that requiredness depends on usage context and defines no independently testable condition. |
| `TOSCA13-5.3.9.4-002` | catalog_quality=informative | States that requiredness depends on usage context and defines no independently testable condition. |
| `TOSCA13-5.8.1-001` | catalog_quality=author-only | Permission for type designers to omit implementation artifacts, not a processor algorithm. |
| `TOSCA13-8.3.1.2-001` | catalog_quality=informative | Use-case rationale ('typically required'), not normative processor behavior. |
| `TOSCA13-8.4.1-001` | catalog_quality=informative | Use-case design guidance without an atomic validation condition. |
| `TOSCA13-8.4.2-002` | catalog_quality=author-only | Permission for service-template designers, not processor behavior. |
| `TOSCA13-8.5.2.1-007` | catalog_quality=malformed | The source sentence is truncated ('and.') and cannot define an independently testable rule. |
| `TOSCA13-8.6.2-002` | catalog_quality=informative | Introduction to an example use case, not a normative processor obligation. |
| `TOSCA13-3.10.1-007` | applicability=not-applicable | Obligation is delegated to domain-profile specifications/authors. |
| `TOSCA13-4.1-011` | applicability=not-applicable | The text explicitly assigns HOST traversal to a TOSCA orchestrator over running instances. |
| `TOSCA13-4.7.1.2-001` | applicability=not-applicable | The text assigns running-instance search to a TOSCA orchestrator. |
| `TOSCA13-5.2-003` | applicability=not-applicable | Name-collision avoidance is an obligation on profile authors. |
| `TOSCA13-5.8.1-002` | applicability=not-applicable | No-op treatment is orchestrator runtime behavior. |
| `TOSCA13-5.8.5.5-001` | applicability=not-applicable | Activation ordering is orchestrator runtime behavior. |
| `TOSCA13-5.8.5.5-002` | applicability=not-applicable | Capability advertisement after fulfillment is orchestrator runtime behavior. |
| `TOSCA13-3.10.3.1-001` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.10.1-002. |
| `TOSCA13-3.6.10.2-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.10.2-001. |
| `TOSCA13-3.6.10.4-006` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.10.2-001. |
| `TOSCA13-3.6.10.5-002` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.10.4-008. |
| `TOSCA13-3.6.12.2-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.12.2-001. |
| `TOSCA13-3.6.12.3-005` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.12.2-001. |
| `TOSCA13-3.6.12.4-001` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.10.5-001. |
| `TOSCA13-3.6.14.3-001` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.14.2-015. |
| `TOSCA13-3.6.22.1-006` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.22.1-004. |
| `TOSCA13-3.6.23.1.1-005` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.23.1.1-003. |
| `TOSCA13-3.6.9.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.9.1-001. |
| `TOSCA13-3.6.9.2-005` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.6.9.1-001. |
| `TOSCA13-3.7.2.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.2.1-001. |
| `TOSCA13-3.7.2.2.2-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.2.1-001. |
| `TOSCA13-3.7.3.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.3.1-001. |
| `TOSCA13-3.7.3.3-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.3.3-002. |
| `TOSCA13-3.7.3.5-001` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.3.1-001. |
| `TOSCA13-3.7.3.5-004` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.7.3.1-001. |
| `TOSCA13-3.8.3.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.8.3.1-001. |
| `TOSCA13-3.8.4.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.8.4.1-001. |
| `TOSCA13-3.8.5.1-003` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-3.8.5.1-001. |
| `TOSCA13-5.3.11.1-004` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-5.3.11.1-002. |
| `TOSCA13-5.3.6.1-007` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-5.3.6.1-005. |
| `TOSCA13-5.3.6.1-010` | catalog_quality=duplicate | Same independently testable obligation as TOSCA13-5.3.6.1-008. |

## Evidence interpretation

- `implemented` is based on a complete code-path trace, independently of direct test availability.
- `partial` means the construct reaches production parsing/semantic code but at least one normative condition is not enforced.
- `missing` is used only after the read, namespace, hierarchy, inheritance, rendering, function, normalization, and CSAR paths relevant to the record were searched.
- `verified` remains reserved for an executed direct test/probe with the intended result reason; positive examples are only `indirectly-tested`.
