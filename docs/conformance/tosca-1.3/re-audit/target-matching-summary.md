# Concrete requirement target remediation — TOSCA 1.3

Concrete assignments now reject target nodes that cannot fulfill the effective
requirement's capability/node/relationship tuple. This closes
REAUDIT-TARGET-MATCHING, not the overall TOSCA 1.3 conformance audit.

Normative authority: the local pinned TOSCA Simple Profile in YAML 1.3 OS HTML.
Sections 3.7.2 (Capability definition), 3.7.3 (Requirement definition and tuple),
3.7.7/3.7.9 (Capability/Node Types and derivation), 3.7.10 (Relationship Type),
and 3.8.2 (Requirement assignment) govern this change. Section 3.8.2.2.1
explicitly requires a fulfilling capability on a concrete short-form target;
extended assignment selectors cannot discard the definition's restrictions.

`work-items/re-audit-target-matching.yaml` contains all fourteen normative
answers, exact excerpts, the 34-case matrix reconciliation, baseline hashes,
execution path and direct evidence. It serves as this remediation's matrix.

## Reproducer and behavior

The follow-up observation in remediation-summary.md is reproducible at
6ce702b7: changing TestReAuditContainerAssignments' Endpoint.network capability
from Endpoint to Root is accepted without diagnostics, and app.network still
normalizes to that endpoint. A minimal equivalent is:

```yaml
tosca_definitions_version: tosca_simple_yaml_1_3
node_types:
  Source:
    derived_from: tosca.nodes.Root
    requirements:
      - link: { capability: tosca.capabilities.Endpoint }
topology_template:
  node_templates:
    source:
      type: Source
      requirements: [{ link: target }]
    target: { type: tosca.nodes.Root }
```

The target only has its inherited Node capability, not Endpoint. The new
rendering check reports `[TOSCA13-TARGET-MATCH]` with source, requirement,
target, explicit selector when present, expected type and sorted actual types.

The former positive container fixture also bound a Storage requirement to
BlockStorage's Attachment capability. Its positive sibling now adds a Storage
capability through a custom derived type. No normative profile was changed;
the network negative now reaches its intended edge without a storage error.

## Implementation and interpretation

Read/namespace/hierarchy/inheritance resolve the entities first.
RequirementAssignments.Render then applies effective defaults and renders the
relationship before calling the versioned validateRequirementTarget policy.
The policy uses hierarchy objects and inherited capability definitions, not
string equality or target rendering order. The default shared dispatch remains
unchanged when the hook is absent; only the 1.3 grammar enables it.

Required base plus actual derived is valid; the reverse and siblings are
invalid. Explicit capability names cannot fall back to another capability.
Relationship valid_target_types is checked against actual capability types,
including the case where the relationship accepts a derived type while the
requirement requests its base. Requirement node types and capability
valid_source_types also constrain the edge.

Normalization records a unique compatible capability. Several valid candidates
remain unbound, with the node, type and relationship restrictions retained.
The sorted admissible set is serialized as capabilityNames; the Clout resolver
honors it, preventing a later selection of a capability already excluded by
rendering. Direct Parse → Compile → Resolve tests prove that behavior. There
is no normative first-candidate rule. Decision 0015 documents this policy
and local symbolic-name precedence over a colliding type selector.

Two requested matrix expectations need normative qualification:

- Section 3.8.2.2.3 prohibits assignment node_filter with a concrete Node
  Template, even when its properties satisfy the filter. The tests reject
  these combinations and preserve filters for abstract runtime selection.
- An omitted required assignment remains an abstract obligation under existing
  semantics; it is not equivalent to an invalid concrete binding. Occurrence
  interval containment, excess assignment counts and both normalized bounds
  remain tested, including the F09 regression suite.

## Verification

Permanent tests cover exact/derived/sibling/reverse compatibility, explicit
names and types, absent and case-mismatched names, type/name collisions,
unique/multiple/no matches, imports with and without namespace prefixes,
inherited/refined capabilities and requirements, node/relationship/source-type
restrictions, relationship templates, filter restrictions, occurrence bounds,
node order and repeated normalization. Earlier inheritance failures are checked
at inheritance; new concrete matching failures at rendering. The ConnectsTo
negative now checks its phase and complete edge diagnostic. The local phase
helper is reused; the general corpus runner was not redesigned.

All eight independent mutations were caught and restored byte-for-byte; exact
patches and failure explanations are in target-matching-mutations.yaml. The
compatibility mutation has an isolated negative assertion showing incorrect
acceptance. Two additional mutations drop the normalized candidate set and ignore it in
the Clout resolver; both direct tests fail. The map-order mutation violates the deferred-selection invariant,
so catching it does not depend on randomly observing two different orders.

Passed:

```text
go test -count=1 ./tests/conformance/tosca_1_3 -run 'TestReAudit(TargetMatching|ContainerAssignments)|TestVerificationConnectsTo'
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
git diff --check
```

Full runs include the TOSCA 1.3 suite, template corpus, F02/F03/F05/F06/F07/F08/
F09/F10/F16 permanent regressions, and the TOSCA 2.0 suite. The added 2.0
isolation test verifies that shared dispatch does not introduce 1.3 capability
binding, candidate lists or occurrence normalization; ordinary 2.0 resolution
also retains its existing target capability. This is an observable regression test,
not an independent normative certification of 2.0 behavior.

All 65 original community inputs were rerun at pinned upstream commit
e5b0a3ee46488921ff409bbe0feba37af2a3233e with source hashes and relative
imports preserved. Expectations were unchanged:

| Result | Count |
| --- | ---: |
| Expected pass | 21 |
| Expected reject | 18 |
| Unexpected pass | 2 |
| Unexpected reject | 0 |
| Ambiguous | 24 |

The two unexpected passes remain the duplicate TMForum apis.yaml files already
explained by valid F02 fixed refinements in remediation-summary.md. The former
storage/container alias unexpected rejects remain resolved. The new report is
target-matching-oasis-results.yaml; repeat generation was byte-identical.

## Catalog and remaining audit

No exact frozen atom covers the complete corrected behavior. Related semantic
records TOSCA13-3.8.2.2.1-003/-004 are outside the frozen denominator. Proposed
atomic additions are in re-audit/catalog-review-queue.yaml; no admission,
reclassification or denominator change occurred. requirements.yaml and
coverage.yaml retain their exact baseline hashes and all 223 historical status
members. Existing ConnectsTo evidence is strengthened without claiming it
covers the entire matching obligation.

Remaining: F11; F12–F15; F04; unfinished independent denominator; 218
unreconciled mappings; strength disagreements; existing ambiguities and the
explicit interpretation in decision 0015. Full processor conformance remains
unconfirmed.
