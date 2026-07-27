# TOSCA conformance matrix

This matrix records verified behavior; it is not a claim of complete
conformance for either specification.

| Area | TOSCA 1.3 | TOSCA 2.0 | Phase | Evidence |
| --- | --- | --- | --- | --- |
| Exact version dispatch | implemented | implemented | read | 1.3 §3.10.3.1; 2.0 §6.2; `tests/conformance/version_selection_test.go` |
| Mandatory version keyname | implemented | implemented | read | 1.3 §3.10.1; 2.0 §§6.1-6.2 |
| Non-string version value | implemented | implemented | read | 1.3 §3.10.1 string type; 2.0 §6.1 `str` type |
| Removed/unknown version rejection | implemented | implemented | read | Exact closed registry and negative selector tests |
| Root grammar isolation | implemented | implemented | read | 1.3 §§3.10.1-3.10.2; 2.0 §§6.1, 6.9.1; cross-version tests |
| Extended trigger-condition isolation | implemented | implemented | read/render | 1.3 §3.6.22.2 retained only in 1.3; 2.0 §16.5 negative regression |
| Retained profile loading | implemented | implemented | read/lookup/render | 1.3 implicit Simple Profile; 2.0 §6.7.2 explicit profile import, §7.3 relationship types, §9.1.2.2 scalar types |
| Minimal service template normalization | implemented | implemented | normalization | Version-specific minimal fixtures |
| Full existing example regression | implemented | implemented | all phases | Version-specific example tests and root corpus |
| Imports and namespace behavior | partial | partial | read/namespaces | Existing corpus; feature-by-feature conformance audit remains |
| Type hierarchy and inheritance | partial | partial | hierarchy/inheritance | Existing corpus; feature-by-feature conformance audit remains |
| Assignments and rendering | partial | partial | rendering | Existing corpus; feature-by-feature conformance audit remains |
| Intrinsic functions | partial | partial | read/render/Clout | Existing examples; exhaustive negative coverage remains |
| Substitution mappings and workflows | partial | partial | read/render/normalization | Existing examples; exhaustive version-specific coverage remains |
| CSAR | partial | partial | URL/CSAR/read | Shared CSAR implementation retained; full version-specific matrix remains |

Status vocabulary: `implemented`, `partial`, `non-compliant`, `unimplemented`,
`ambiguous`, `blocked`, and `not-applicable`.
