# TOSCA 1.3 production remediation after c8d711f3

All nine scoped production findings have direct permanent tests and caught
mutations. This does **not** restore the old 223/223 conformance claim. The 223
records remain a frozen historical baseline; completeness has not been proved.

The authoritative source is the pinned TOSCA Simple Profile in YAML 1.3 HTML.
`work-items/re-audit-production-defects.yaml` records the normative excerpts,
real execution paths, previous statuses, test-first baselines, exact mutation
replay fragments, test commands and remediation statuses. It is the remediation
matrix; the historical requirements.yaml and coverage.yaml are unchanged.

| Finding | Normative section | Previous behavior → corrected behavior | Main production symbol | Permanent direct test | Commit |
| --- | --- | --- | --- | --- | --- |
| F02 | 3.6.10.6 Refining Property Definitions; 3.6.14.2 Grammar | Rejected fixed refinement → retain final fixed values, reject later changes | ReadPropertyDefinition; inheritPropertyDefinition; validateFixedPropertyAssignment | TestReAuditFixedPropertyRefinement | 93476846 |
| F03 | 3.6.10.6 Refining Property Definitions | Child constraints replaced parent → conjoin constraints through multiple levels | inheritPropertyDefinition | TestReAuditInheritedPropertyConstraints | 93476846 |
| F05 | 3.6.10.2 Keynames; 3.6.10.3 Status values | Accepted arbitrary status → four exact values, inherited status and supported default | ReadPropertyDefinition; validateFilePropertyRefinements | TestReAuditPropertyStatus | 9b84d6f1 |
| F06 | 3.5.1 Required keynames; 3.6.7.1 Keynames | Accepted checksum without algorithm → reject missing effective algorithm after inheritance | validateArtifactChecksums | TestReAuditArtifactChecksum | 183cc717 |
| F07 | 5.9.12.1 Definition | Container.Runtime.host Container → Compute | bundled nodes.yaml; ordinary implicit-profile inheritance | TestReAuditContainerProfile | a6d34ec4 |
| F08 | 5.9.13.1 Definition | Container.Application.network Network → Endpoint | bundled nodes.yaml; ordinary implicit-profile inheritance | TestReAuditContainerProfile; TestReAuditContainerAssignments | a6d34ec4 |
| F16 | 5.2 Naming Conventions; 5.4.4.3/4, 5.5.14, 5.9.9–11, 8.5.5 name tables | Twelve missing aliases → exact case-sensitive aliases identify the same parent object as the full URI | additionalNormativeNames; newImportNameTransformer | TestReAuditNormativeAliases | f5030091 |
| F09 | 3.7.3.1 Keynames; 3.8.2.1 and 3.8.2.2.2 requirement-assignment grammar | Lost upper bound and zero-lower declarations → validate containment and normalize the entire interval | validateRequirementOccurrences; normalizeRequirementOccurrences; Requirement.Marshalable | TestReAuditRequirementOccurrences | 69f6286b |
| F10 | 3.3.3.1 Grammar; 3.3.3.2 Keywords | Rejected signed numeric bounds → signed endpoints, ordered ranges, distinct UNBOUNDED | ReadRange; rangeConstraintBounds; Range.Within | TestReAuditSignedRange | containing commit: fix(tosca-1.3): support signed range bounds |

For each finding the permanent test failed against the pre-fix implementation
for the target reason. All nine old-behavior mutations were caught and restored.
Positive tests assert effective or normalized semantics; negative tests assert
the earliest observed phase and target diagnostics through a local phase helper.
The general corpus runner was not redesigned.

For **each of seven production groups**, focused tests and these commands passed:

```text
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
git diff --check
```

The full runs include the TOSCA 1.3 suite, template corpus and TOSCA 2.0 suite.
Shared changes have permanent observable 2.0 regression tests for property
precedence, import alias isolation and exact requirement-count expansion.
TOSCA 1.3 policies are enabled only in its grammar. F10 leaves the shared
unsigned reader and 2.0 native range behavior unchanged.

## Catalog mapping

F10 maps exactly to frozen IDs TOSCA13-3.3.3.1-004, -005 and -006.
F02, F03, F05, F06, F07, F08, F09 and F16 do not have exact frozen mappings
for the complete corrected obligations. Related non-frozen records are listed
in the work item, not misrepresented as frozen coverage. Eight proposal groups
are recorded in `catalog-review-queue.yaml`; atomic admission and denominator
changes require the separate catalog review. No denominator change is proposed
as a certified total here.

## Profile and community results

`remediation-profile-three-way.yaml` records exact specification/OASIS/Puccini
agreement for F07 and F08. All 12 F16 aliases now resolve; 21 full/short/qualified
spellings and 21 wrong-case siblings are tested. OASIS declares the full names;
its YAML is not an executable alias resolver, so no fictitious OASIS alias test
result is credited. Other profile conflicts and ambiguities remain open.

All 65 original OASIS inputs were rerun at pinned commit
`e5b0a3ee46488921ff409bbe0feba37af2a3233e`, preserving original bytes, relative
layout and historical expectations. Their SHA-256 identities were checked.
The archive hash also matches the audit archive. Results remained the same in
a final repeat after F09/F10:

| Result against unchanged expectations | Count |
| --- | ---: |
| Expected pass | 21 |
| Expected reject | 18 |
| Unexpected pass | 2 |
| Unexpected reject | 0 |
| Ambiguous | 24 |

The two former unexpected rejects, `policies-and-groups.yaml` and
`source-and-target.yaml`, now pass. The two unexpected passes are the duplicate
TMForum `apis.yaml` inputs: their id/name/swagger_file properties refine
properties declared by TMFAPI. Section 3.6.10.6 permits these fixed refinements;
F02 correctly changes acceptance. The old audit incorrectly called them new
property definitions. Expectations were retained and the discrepancy explained,
not changed to make the report green. Per-input evidence is in
`remediation-oasis-community-results.yaml`.

## Normalization and remaining work

Consumers of normalized TOSCA 1.3 requirements must honor the new complete
`occurrences` interval. Zero lower bounds retain a declaration; null upper
bounds represent UNBOUNDED. Numeric range values also use a signed lower bound
and nullable upper bound. Decisions 0012–0014 document notation and representation
choices. These output formats are Puccini decisions, not specification grammar.

The following remain separate, unfinished re-audit work:

- F11: complete independent extraction and denominator reconciliation.
- F12–F15: catalog atomicity/applicability, corpus evidence, phase-runner and
  surviving-mutation methodology gaps.
- F04: misleading empty-constraints diagnostic.
- 218 unreconciled mappings, strength disagreements and recorded ambiguities,
  including unresolved profile/community cases.
- Full TOSCA 1.3 processor conformance remains unconfirmed.

An additional follow-up observation is preserved separately from F07/F08:
changing the `Endpoint.network` capability in TestReAuditContainerAssignments
from tosca.capabilities.Endpoint to tosca.capabilities.Root still accepts with
empty diagnostics in ordinary Parse. The normalized app.network target remains
that endpoint node although its requirement definition is Endpoint. This needs
separate target-capability matching review (3.7.3/3.8.2); it was not remediated or
credited as a negative assignment test. F07/F08 direct evidence proves the
correct effective profile definitions and incompatible definition refinement,
not exhaustive target-node matching.
