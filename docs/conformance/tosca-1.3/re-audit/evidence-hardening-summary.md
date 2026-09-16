# TOSCA 1.3 evidence hardening

The hardened evidence rules verify **5 of the historical 223 records**;
**218 have insufficient evidence**. This is an evidence assessment, not an
implementation count or a newly established normative denominator. The original
`requirements.yaml` and `coverage.yaml` remain byte-for-byte unchanged.

## F12–F15 disposition

| Finding | Result |
| --- | --- |
| F12 | Evidence consequences addressed: historical membership and old status labels no longer establish verification. The actual finding also concerns catalog atomicity and processor applicability. That part remains open under the prohibition on catalog changes. |
| F13 | Primary and supporting links separated; exact reviewed fixture/assertion ownership required. Generic fixtures and supporting labels cannot confer verification. Semantic predicates require value/resolution assertions, and negative evidence requires targeted diagnostics and nearest-valid isolation. |
| F14 | Both corpus runners enforce the actual first failing parser boundary. Deliberate wrong-phase cases and a mutation of the runner itself prove enforcement. |
| F15 | All 15 original production mutations are caught, including both former survivors. Two additional mutations of newly discovered defects are also caught. Required subsystem categories have checked detecting cases and assertions. |

The original audit files remain historical evidence. This report and
`evidence-hardening-plan.yaml` record the new disposition without rewriting
the audit's original meaning.

## Phase enforcement and ownership

The executed main and non-MUST corpora contain 167 cases: 127 main and 40
non-MUST. All 50 negative cases assert exact phase identity and have an executed
nearest-valid sibling: 26 read, 4 namespaces, 1 hierarchy, 4 inheritance and
15 rendering. Thirteen negative cases additionally have reviewed primary
ownership with owned phase and targeted diagnostic assertions; the other
negative cases remain regression/supporting evidence.

Phase identity comes from test-side calls to the ordinary parser APIs, stopping
at the first failing boundary. Resolved function validation belongs to rendering
in this architecture. CSAR opening/metadata validation belongs to read. No
normative normalization-rejection category is invented; positive normalization
semantics are asserted instead. There is no parser boundary earlier than read
for an earlier-phase CSAR sibling.

The actual runner rejects five deliberately false later-phase labels on a read
failure. Disabling its phase checks makes that integrity test fail. The helper
also checks the complete phase identity matrix and independently constructed
read, namespace, hierarchy, inheritance, assignment, resolved-function and CSAR
failures, including earlier-invalid siblings where such a boundary exists.

There are 38 primary case-to-requirement links, 94 executed owned assertions and
17 independently reviewed predicates. All 25 primary positive cases have
predicate-specific result assertions where semantic evidence is required.
There are 5,998 supporting links, retained as historical navigation across MUST
and non-MUST records; they confer no verification. Of the historical 223 IDs,
13 have reviewed owned evidence and 210 do not. Owning evidence for a subset
of a bundled historical predicate does not verify the whole record.

Schema and integrity gates reject missing ownership, unrelated requirement
claims, supporting-only verification, missing or invalid nearest-valid siblings,
semantic acceptance-only evidence, wrong phases, missing mutation categories,
missing detecting fixtures/assertions, stale hashes and uncaught mutations.
Review bindings include the exact assertion, fixture bytes and negative sibling.
Automated checks preserve a reviewed normative mapping; they do not replace the
human interpretation of the specification.

## Mutation replay

`mutation-coverage.yaml` and `evidence-mutation-replay.yaml` record **17 attempted,
17 caught, 0 missed** production mutations, with original M01–M15 identifiers,
exact patches, commands, passing baselines, failing owned assertions and restored
source hashes. All 15 required subsystem categories are represented. The
additional runner-phase mutation was also caught; total including infrastructure
is **18/18 caught**.

The original operation assignment overwrite is detected by the effective
relationship operation implementation (`local.sh`, not `parent.sh`). The
Credential default mutation is detected by the inherited effective
`token_type = password` assertion. Compiler errors or unrelated failed tests
are never counted as detecting evidence. All mutations were fully restored.

## New production defects discovered by direct evidence

Both defects followed the requested separate failing-discovery/remediation
commit policy. These are the only production changes in this task.

| Defect | Normative basis and result | Discovery | Remediation |
| --- | --- | --- | --- |
| EVIDENCE-OPERATION-INHERITANCE | §3.6.17 Operation definition (grammar and additional requirements), §5.8.1 Additional Requirements: a description-only local definition lost the inherited implementation. A 1.3 inheritance policy now preserves it unless another artifact is assigned. | `f38c63a2` | `2a2d43f0` |
| EVIDENCE-CREDENTIAL-TOKEN-DATA | §5.3.6 Credential definition and §4.3.3 token function: valid `{token: secret}` data was interpreted as malformed function syntax at read. A 1.3 literal-data policy preserves the string-valued map; sequence-valued token syntax retains function validation. | `385f9fc7` | `f8d65aa7` |

The sole normative source is the pinned
`docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html`.
The absence-of-override inheritance interpretation and function/data precedence
are documented in `docs/decisions/tosca-1.3-operation-implementation-inheritance.md`
and `docs/decisions/tosca-1.3-token-data-disambiguation.md`; neither is presented
as an explicitly specified general algorithm.

Production paths are `OperationDefinition.Inherit` through a grammar policy in
the inheritance phase, and `ParseFunctionCall` through a grammar literal-data
policy during value reading. Policies are registered only in `tosca_v1_3`;
shared `tosca_v2_0` files contain neutral hooks. New 2.0 characterization tests
prove that its observable behavior is unchanged. Operation cases cover parent,
override, absent override, multiple inheritance levels and unrelated operations.
Credential cases cover omission, explicit password, an allowed alternate string,
inherited defaults, derived custom types, invalid string type, nested collections,
ordinary token functions and normalized results.

## Historical evidence recheck

`evidence-recheck.yaml` independently lists every historical ID:

| Result | Count | Reason |
| --- | ---: | --- |
| verified | 5 | Complete reviewed predicate, executed owned assertions, applicable positive/negative evidence and caught relevant mutation path. |
| insufficient-evidence | 210 | No reviewed primary predicate; historical labels and unreviewed direct Go tests are not automatically credited. |
| insufficient-evidence | 8 | Reviewed assertions cover only a subset of the historical predicate. |

The five verified records are `TOSCA13-3.6.8.2.3-002`,
`TOSCA13-3.6.17.3-002`, `TOSCA13-5.5.7.4-001`,
`TOSCA13-5.5.12.3-001` and `TOSCA13-5.8.1-004`.
The evidence result does not resolve whether any historical record is redundant,
properly atomic, processor-applicable or part of a complete denominator.

The generator records actual Go test execution and observed phases in
`evidence-execution.yaml`. A fingerprint covers production sources, bundled
assets, corpus code, manifests, reviewed bindings and imported fixtures. Changes
require a new execution capture. Go integrity tests independently recompute the
per-ID classification, enforce the frozen file hashes and verify report totals.

Reproduce with:

```sh
python3 scripts/conformance/recheck_tosca13_evidence.py --record-run
python3 scripts/conformance/recheck_tosca13_evidence.py --check
python3 scripts/conformance/replay_tosca13_evidence_mutations.py
```

Run mutation replay separately from other builds or source edits. A changed
production/fixture hash invalidates old mutation evidence. Report regeneration
from the same recorded execution is deterministic. Generated section and
variation documents now state inventory/evidence limitations instead of claiming
exhaustive coverage.

## Regression and OASIS results

The following checks passed on the completed changes:

```sh
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
git diff --check
python3 scripts/conformance/check_tosca13_non_must.py
python3 scripts/conformance/generate_tosca13_oasis_corpus.py --check
python3 scripts/conformance/recheck_tosca13_evidence.py --check
```

These include the full TOSCA 1.3 corpus, direct conformance and prior re-audit
remediation tests (including target matching), inherited types, functions,
normalization, the eight adapted OASIS profile cases, and full TOSCA 2.0 tests.
Both production remediation commits also passed normal, race and vet checks.

All 65 previously executed community inputs were replayed against the final
production behavior with unchanged expectations: 21 expected passes, 18 expected
rejects, 2 unexpected passes, 0 unexpected rejects and 24 ambiguous. The two
unexpected passes are the pre-existing TMForum `apis.yaml` expectation mismatch
following valid fixed property refinement support; they are not new evidence of
a regression. Ambiguous dependencies and old community expectations remain
supplementary evidence and are not promoted to normative proof. Detailed source
hashes and results are in `evidence-hardening-oasis-results.yaml`.

## Commits and remaining work

The work is split into infrastructure (`9d0de08b`), operation discovery/fix
(`f38c63a2`, `2a2d43f0`), Credential discovery/fix (`385f9fc7`, `f8d65aa7`),
mutation gate (`4afab1f0`) and this final evidence recheck/report commit.

Still open: F11; F04; the catalog atomicity/applicability portion of F12;
independent catalog extraction; 218 unreconciled catalog mappings; strength
disagreements; ambiguities; and establishment of a new proven denominator.
The 218 insufficient-evidence records are a separate measurement from the
218 unreconciled catalog mappings. Further reviewed evidence is also needed
before those historical statuses can be promoted. Full TOSCA 1.3 conformance
is not restored or claimed.
