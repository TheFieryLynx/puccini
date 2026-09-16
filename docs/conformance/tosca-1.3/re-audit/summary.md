# TOSCA 1.3 independent red-team re-audit

**The existing complete processor conformance claim is refuted. This is not a completed independent recertification or a certified replacement denominator.** Confirmed defects and mutation results are reviewable; the full atomic extraction, semantic catalog reconciliation and per-requirement fixture verification remain incomplete. Unknown/null totals below are intentional and must not be reported as zero or as passed checks.

The audit preserved the clean starting tree and HEAD `b4d40d3ebea7a74a5a55e49020d5fcbbf015be34`. Production behavior, bundled and pinned profiles, existing tests/expectations and the frozen catalog have no persistent changes. Only this report directory is committed.

## Independent extraction

The blind extraction checkpoint contains **642 candidate obligations**, all provisionally processor-applicable under the extraction scope. It was constructed before reading old requirements, coverage or tests. **The independently validated atomic MUST total and processor MUST total are undetermined.** Table/grammar conflicts, unexpanded YAML-only fields, umbrella/atomicity reconciliation and subsequent list-item checks prevent certification of 642 as the new denominator.

The separate inventory retains **105** excluded/review candidates: **35 non-MUST**, **34 informative**, **15 target exclusions**, **21 pending review**. Overlapping target flags are author=71, orchestrator=9, archive=2, generator=0; processor=false=105. These are candidate-inventory counts, not a complete count of all excluded obligations.

All **1121** numbered HTML headings were independently enumerated, versus **590** existing section records. The extra **531 headings** are not automatically missing requirements: many are examples or structural subdivisions. Concrete problematic classification/evidence examples include §3.6.10.6 and §§5.9.12.1–5.9.13.1. A later `<li>` check recovered 13 additional text blocks; its post-comparison timing is disclosed in catalog-validation.md.

## Catalog comparison

* Exact ID matches: **0**, because a new ID namespace was used.
* Direct table predicate semantic matches: **404 candidate matches** (field/predicate keys; not completed contextual semantic reconciliation).
* Strength-disagreement candidate mappings: **20**.
* Remaining unverified candidate mappings: **218**.
* Certified wholly new obligations, old-only-valid/invalid totals, and final proposed denominator: **undetermined**.
* Eight concrete review proposals cover strength, duplicate-old, atomicity and applicability disagreements in catalog-diff.yaml. These are proposals, not edits to the frozen catalog.

The omitted explicit requirement-range sentence in the malformed §3.8.2.3 heading has a related existing table record, `TOSCA13-3.8.2.1-015`; it is not falsely called wholly new merely because the section ID is absent. The old catalog has 3101 records; omission from its 223 gate and absence from the entire catalog are different findings.

**Frozen stored denominator = 223; status sum = implemented 223.** Requirement IDs, classifications and all baseline bytes remain unchanged. This verifies preservation only.

## Profile comparison

Inventory: **66 specification types, 66 OASIS types, 64 types in Puccini's versioned folder, all 66 normative types available through ordinary implicit imports**. JSON/XML are supplied by the implicit profile, so they are not missing types. No extra OASIS or version-folder Puccini types were found in this 66-type inventory.

Declared-definition reconciliation: **52 semantically equivalent**, **2 specification/OASIS conflicts**, **2 specification/Puccini conflicts**, **10 ambiguous**. The two confirmed Puccini definition-field mismatches are:

* §5.9.12.1: `Container.Runtime.capabilities.host.type` is `Container`, but must be `Compute`.
* §5.9.13.1: `Container.Application.requirements.network.capability` is `Network`, but must be `Endpoint`.

A separate name-table audit tested **100 aliases: 85 accepted, 12 unambiguous unexpected rejects, 3 editorial ambiguities**. All 12 full-Type-URI siblings accept. Affected types are Implementation.Bash, Implementation.Python, network.Bindable, Abstract.Storage, Storage.ObjectStorage, Storage.BlockStorage, and network.BindsTo. Counting name failures as well as declared fields gives 9 types with confirmed Puccini conflicts; this does not mean nine independent implementation bugs.

Both were confirmed in effective inherited definitions obtained through ordinary `parser.Context.Parse`, and by minimal rejected valid inputs. The full untouched OASIS profile fails direct import due to namespace/collision/root special handling. A temporary loading adapter subsequently produced an effective OASIS graph through ordinary parser phases: 34/66 inherited definitions match after descriptive metadata removal and 32 differ (including inherited repetitions, defaults, and open prose conflicts). The adapter and its restrictions are fully embedded in profile-three-way.yaml; it is not credited as a strict-mode acceptance result for the original profile. The 10 ambiguous declared-definition comparisons therefore remain open, not certified matches. NetworkInfo/PortInfo requiredness is context-dependent under their Additional Requirements and was not mislabeled a defect based on raw YAML alone.

## OASIS community inputs

Upstream was fetched at pinned commit `e5b0a3ee46488921ff409bbe0feba37af2a3233e`. The inventory executes original `.yaml` sources declaring version 1.3 in tests, examples and the normative 1.3 profile, including two 1.3 support files located in a 2.0 test directory.

**65 applicable source files executed: 19 accepted, 46 rejected.** Source-based triage: **19 provisional passes, 20 expected rejects, 2 unexpected rejects, 24 ambiguous**. Unexpected accepts: **0**. The two unexpected rejects use storage aliases explicitly permitted by §5.2 and the type tables; full-URI-only siblings accept. Unresolved cases are not passed checks: unresolved dependencies, unsupplied inputs, profile collisions and source interpretation keep 24 cases unclassified. Accepted inputs are not credited as exhaustive semantic tests. Other vendor profile suites were not exhaustively executed.

## Corpus and implementation audit

The main manifest still has **63 cases**; all manifests together contain **101 cases**. Of the 642 independent candidates, **639** have section-labelled fixture candidates and **3** do not. These counts measure labels, not valid evidence. Full counts of fully evidenced requirements, missing positives/negatives/boundaries and actual wrong-phase failures remain **undetermined**.

The generic definitions fixture has no fixed refinement or additive constraints; its invalid sibling fails on `unknown_key`. The generic normative-profile fixture only derives Compute. Crediting these to the unrelated refinement/container obligations is demonstrably false confidence. The runner also never asserts the expected phase: it only checks diagnostic fragments.

Fresh source tracing located 69 live reader registrations, 35 delegated directly to `tosca_v2_0`. The per-candidate matrix has **9 verified non-compliant entries and 633 unknown** entries; it does not reuse old implementation symbols or promote route existence to implementation proof. Additional confirmed findings, such as status domain validation, are listed independently even when the blind candidate set did not isolate the relevant atom.

## Mutation audit

**15 attempted; 13 caught by combined conformance/corpus suites; 2 survived. Template corpus alone caught 4 and missed 11.** Every mutation was restored byte-for-byte.

* M04 survives: copied source operation implementation overwrites a local implementation. Direct replay changes `child-implementation-preserved=true` to `false`; the production route is RelationshipAssignment.Render → InterfaceAssignments.CopyUnassigned → OperationAssignments.CopyUnassigned.
* M13 survives: changing `Credential.token_type` default from `password` to `audit-mutant` changes the effective definition observed by ordinary Parse.

The score is for these selected mutations, not an exhaustive mutation score.

## Findings

* **F01 — critical**: The complete processor conformance and 223/223 verified claim is not sustainable.
* **F02 — high**: Valid fixed property refinements are rejected.
* **F03 — high**: Adding a child property constraint drops the parent constraint.
* **F04 — low**: Empty constraints are correctly rejected with a misleading diagnostic.
* **F05 — high**: Undefined property status values are accepted.
* **F06 — high**: Artifact checksum does not require its checksum_algorithm.
* **F07 — high**: Container.Runtime host capability has the wrong normative type.
* **F08 — high**: Container.Application network requirement has the wrong normative capability.
* **F09 — high**: Requirement assignment upper occurrence bound is neither validated nor represented.
* **F10 — high**: Generic range incorrectly forbids negative integer bounds.
* **F11 — high**: Mandatory grammar and semantics are omitted from the closed MUST claim.
* **F12 — medium**: Frozen atomicity and processor applicability require separate catalog review.
* **F13 — medium**: Large manifest label sets are not direct evidence of the labelled obligations.
* **F14 — medium**: Corpus expected-phase metadata is not asserted.
* **F15 — medium**: Mutation survivors disprove exhaustive regression sensitivity.
* **F16 — high**: Twelve declared normative aliases across seven types fail namespace resolution.

Totals: **1 critical, 10 high, 4 medium, 1 low**. F04 is a withdrawn conformance false positive: §3.6.10.4 requires one or more constraint clauses, so rejecting `constraints: []` is correct; only its misleading unsupported-key diagnostic remains low severity. Nine distinct production-behavior defects are established by the other targeted probes.

## False confidence checks

### Denominator is exactly 223

* Supporting evidence: Stored included flags sum to 223; baseline hashes unchanged.
* Against: Mandatory grammar/semantics excluded by strength, duplicate presence predicates, bundled atoms and archive applicability issues.
* Conclusion: Stored count confirmed; complete atomic processor denominator not confirmed.
* Confidence: high.

### No processor MUST obligations are missed

* Supporting evidence: Old catalog contains 3101 records beyond the frozen subset.
* Against: Refinement fixed values, additive constraints, checksum conditional presence and assignment range restrictions do not receive a complete mandatory gate.
* Conclusion: Refuted as a completeness claim about the 223 subset.
* Confidence: high.

### Every frozen ID is atomic

* Supporting evidence: Many rows have a single clear predicate.
* Against: Range integer/presence bundled; operation and notification prohibitions combined; repeated requiredness prose counted separately.
* Conclusion: Not confirmed; separate atomicity review required.
* Confidence: high.

### Every frozen ID is processor-applicable

* Supporting evidence: Most grammar rows follow 14.3(a/b).
* Against: CSAR entries are incorporated under 14.4(c)/14.6, not required of a processor by 14.3.
* Conclusion: Not confirmed for a pure processor denominator.
* Confidence: high.

### Tests verify the correct behavior for all 223 IDs

* Supporting evidence: Full Go, race and vet checks pass.
* Against: Generic signed ranges fail while range IDs remain implemented/verified. Missing boundary invalidates the inference.
* Conclusion: Refuted.
* Confidence: high.

### Every normative profile definition is correct

* Supporting evidence: All 66 normative full Type URIs appear on ordinary implicit import; 52 declared definitions reconcile under stated normalizations.
* Against: Two independently reproduced container field mismatches; 10 further declared-definition comparisons unresolved.
* Conclusion: Refuted.
* Confidence: high.

### Automatic Simple Profile import is correct

* Supporting evidence: All 66 normative full Type URIs, including json/xml, are available via ordinary Parse; this does not imply all declared aliases work.
* Against: Import loads wrong Runtime.host and Application.network definitions; twelve declared aliases do not resolve.
* Conclusion: Availability supported; full semantic correctness refuted.
* Confidence: high.

### Corpus variations exhaust normative classes

* Supporting evidence: 63 main cases plus supplemental manifests; existing tests all pass.
* Against: Tiny generic fixtures credited to many unrelated obligations; 11/15 mutations escape template corpus.
* Conclusion: Refuted.
* Confidence: high.

### Normalization hides no semantic loss

* Supporting evidence: Existing normalization assertions pass for their samples.
* Against: Requirement assignment upper bound is discarded; widened range accepted.
* Conclusion: Refuted.
* Confidence: high.

### Diagnostics are deterministic

* Supporting evidence: Existing negative corpus fixtures parse twice and compare identical diagnostic strings.
* Against: No exhaustive permutation/concurrency proof; wrong-phase metadata is not actually checked.
* Conclusion: Supported for tested cases only; universal claim not confirmed.
* Confidence: medium.

### TOSCA 1.3 behavior does not depend on 2.0 implementation details

* Supporting evidence: Version-specific grammar registry and hooks exist; cross-version test package passes.
* Against: 35/69 readers delegate directly to v2_0, plus shared model/function code; missing 1.3 refinement value support is concrete.
* Conclusion: Literal code independence refuted; arbitrary semantic coupling is not inferred from sharing alone.
* Confidence: high.

## Conformance conclusion

1. **Is denominator 223 confirmed?** Only the stored/frozen count is confirmed. Its atomicity, applicability and completeness are not. No exact replacement denominator is asserted.
2. **Is implemented/verified 223/223 confirmed?** No. Valid signed ranges fail existing frozen range obligations, and automatic normative imports expose incorrect definitions.
3. **Is corpus completeness confirmed?** No. Counterexamples, unrelated fixtures and surviving mutations refute it.
4. **Can the existing complete TOSCA 1.3 processor conformance claim remain?** No. Retain the 223 count only as the historical frozen baseline, not as proof of complete processor conformance.
5. **What is required before closure?** Separate fixes for the nine demonstrated production defects; direct positive/negative/boundary/refinement/profile tests; phase assertions; tests that kill M04/M13; a separate catalog-review change; complete independent atomic extraction and semantic reconciliation; effective OASIS/profile ambiguity resolution; resolve ambiguous community cases; then regenerate and independently validate conformance artifacts.

## Standards and verification

Applicable specification: **TOSCA Simple Profile in YAML 1.3**, pinned OASIS Standard only. Principal normative basis: §3.1.2 TOSCA Namespacing in TOSCA Service Templates; §3.3.3.1 Grammar (range); §3.6.7.1 Keynames (artifacts); §§3.6.10.2–3.6.10.6 property keynames/status/grammar/refinement; §3.6.14.2 parameter grammar; §3.8.2 requirement assignment; §§5.2–5.2.1 TOSCA normative type names and Additional requirements; §§5.3.6.1–5.3.6.2 Credential; §§5.9.12.1–5.9.13.1 container definitions; §§14.2–14.6 conformance targets. Exact normative excerpts and source lines accompany findings.

Production phases changed: **none permanently**. Findings concern read, namespace lookup, inheritance, rendering, implicit profile import and normalization. Existing test expectations and conformance matrices are unchanged; this audit has its own matrix. Temporary positive, negative and boundary probes are embedded verbatim in findings.yaml and removed as executable files before commit.

Checks completed successfully: `go test ./...`, `go test -race ./...`, `go vet ./...`, both version-specific packages for each mutation, and `sha256sum -c third_party/oasis/tosca-simple-profile/1.3/SHA256SUMS`. Final uncached restoration tests and artifact/hash checks are recorded in baseline.yaml.

Cross-version impact: **no persistent TOSCA 1.3 or TOSCA 2.0 behavior changes**. Existing cross-version tests passed. Shared implementation dependency is observed, with a concrete 1.3 refinement failure; broader semantic independence is not assumed. Ambiguities and limitations are recorded in ambiguities.yaml and catalog-validation.md; no production decision record was changed.

Audit evidence status: **partial, with confirmed non-compliance**. This commit does not claim completion of the requested exhaustive 13-phase audit. Required follow-up is stated explicitly rather than treating unknown evidence as verified.

Audit self-review also corrected an earlier false-positive community classification: BlockStorage/ObjectStorage are explicitly declared standard aliases, not obsolete invalid syntax. Final F16 and community results supersede that preliminary hypothesis. The exact validated replacement denominator and comprehensive positive/negative/boundary coverage totals remain unfinished, not zero.
