# TOSCA 1.3 audit validation

## Statistical result

- Deterministic stratified sample seed: **20260727**.
- Sample size: **146** records.
- Strata: 50 legacy `partial`; all 36 `non-compliant`; all 3 `unimplemented`; all 25 `unverified`; 20 `not-applicable`; 20 `ambiguous`.
- Misclassified sampled records: **95** (65.07%).
- False-positive rate: **43/81 (53.09%)**.
- False-negative rate: **21/65 (32.31%)**.
- Duplicate clusters found: **96**.
- Catalog records classified as informative/non-normative: **35**.
- Erroneous `not-applicable` records in the sample: **14**.
- Confidence status: **original aggregate metrics rejected** because sampled classification error exceeds 5%. The regenerated atomic MUST metrics now include a complete requirement-by-requirement implementation trace; non-MUST unknown records remain outside that claim.

False positive means the old audit asserted a concrete implementation state from insufficient code-path evidence. False negative means it failed to credit traced implementation or incorrectly excluded an applicable static grammar/profile obligation.

## Record-by-record manual review

| Legacy | Requirement | Section reference accurate | Normative | Atomic | Applicability correct | Implementation evidence correct | Test evidence correct | Status correct | Duplicate-free | Revised axes |
|---|---|---:|---|---|---|---|---|---|---|---|
| partial | `TOSCA13-3.9-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.16.2.4-005` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.3.6.5.1-011` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.3.7.1-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.9.10-002` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.8.1.2-002` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.8.4.1-007` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.8.2.2-006` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.9.2.7.1-005` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.9.3.3-003` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.3.6.2-005` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.8.2.2-004` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.8.1.1-003` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.6.4.1.2-002` | yes | yes | no | no | no | yes | no | yes | unknown / untested / not-applicable |
| partial | `TOSCA13-3.10.1-046` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.10.4-002` | yes | yes | no | no | no | yes | no | yes | unknown / untested / not-applicable |
| partial | `TOSCA13-3.6.9.2-007` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.7.2.2.2-003` | yes | yes | no | yes | yes | yes | no | no | partial / untested / applicable |
| partial | `TOSCA13-3.6.23.4.1-002` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.4.2-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.5.7.1-009` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.8.3.1-029` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.21.2-004` | yes | yes | no | no | no | yes | no | yes | unknown / untested / not-applicable |
| partial | `TOSCA13-5.9.8.1-008` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.14.1-002` | yes | yes | no | no | no | yes | no | yes | unknown / untested / not-applicable |
| partial | `TOSCA13-3.7.3.2.3-007` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.7.10.1-005` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.17.1-004` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.9.9.2-006` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.6.17.1-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.3.6.6.1-005` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.5.3.2-021` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.10.3.11-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.10.1-033` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.5.7.3-005` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.8.13.2-003` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-5.9.12.1-002` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.8.2.2.3-009` | yes | no | no | no | no | yes | no | yes | unknown / untested / not-applicable |
| partial | `TOSCA13-5.5.7-004` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.8.3.3-001` | yes | yes | yes | yes | yes | yes | yes | yes | partial / untested / applicable |
| partial | `TOSCA13-5.7.5.1-007` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.11.2.1-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.16.1-012` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.7.2.1-001` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.23.4.1-010` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.22.2-008` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| partial | `TOSCA13-3.6.22.2-003` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.7.9.1-017` | yes | yes | yes | yes | yes | no | no | yes | implemented / indirectly-tested / applicable |
| partial | `TOSCA13-3.6.10.4-006` | yes | yes | no | yes | yes | yes | no | no | partial / untested / applicable |
| partial | `TOSCA13-4.1-004` | yes | yes | yes | yes | no | yes | no | yes | unknown / untested / applicable |
| non-compliant | `TOSCA13-3.6.3.1-004` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-007` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-010` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-013` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-016` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-020` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-021` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-024` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-025` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-026` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-029` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-032` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.6.3.1-035` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-3.10.2.1-003` | yes | yes | no | yes | yes | yes | yes | no | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-4.3.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-4.3.1.2-001` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-4.4.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-4.4.1.2-001` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-4.4.1.2-002` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-5.4.3.4.1-002` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-5.9.9.2-002` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-5.9.9.2-007` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-5.9.12.1-004` | no | yes | no | no | yes | yes | no | yes | unknown / untested / not-applicable |
| non-compliant | `TOSCA13-5.10.1.1-003` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-5.10.1.1-004` | yes | yes | yes | yes | yes | yes | yes | yes | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-6.3-003` | yes | yes | no | yes | yes | yes | yes | no | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-6.3-004` | yes | yes | no | yes | yes | yes | yes | no | non-compliant / verified / applicable |
| non-compliant | `TOSCA13-6.3-006` | yes | yes | no | yes | yes | yes | yes | no | non-compliant / verified / applicable |
| unimplemented | `TOSCA13-3.10.2.1-001` | yes | yes | no | yes | yes | yes | yes | no | missing / untested / applicable |
| unimplemented | `TOSCA13-3.10.2.1-002` | yes | yes | no | yes | yes | yes | yes | no | missing / untested / applicable |
| unimplemented | `TOSCA13-6.2.1-005` | yes | yes | yes | yes | yes | yes | no | yes | unknown / unverifiable / ambiguous |
| unverified | `TOSCA13-14.1-001` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.1-002` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.1-003` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.1-004` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.1-006` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.1-007` | yes | no | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-001` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-002` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-003` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-004` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-005` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-006` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-007` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-008` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-009` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-010` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.2-011` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-001` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-002` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-003` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-005` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-006` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.3-007` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.6-001` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| unverified | `TOSCA13-14.6-002` | yes | yes | no | no | yes | yes | no | yes | unknown / unverifiable / not-applicable |
| not-applicable | `TOSCA13-8.4.3-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / untested / not-applicable |
| not-applicable | `TOSCA13-8.5.2-004` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-8.5.1.3-002` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-8.5.4-004` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-8.5.1.1-028` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-8.5.2.3-018` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.5.3.2-001` | yes | yes | no | no | yes | yes | no | no | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.3.1-002` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-14.5-003` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / not-applicable |
| not-applicable | `TOSCA13-8.5.2.2-003` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / untested / not-applicable |
| not-applicable | `TOSCA13-8.5.1.3-003` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.3.1.2-003` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-13.2.7-003` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / untested / not-applicable |
| not-applicable | `TOSCA13-8.5.1.3-006` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.5.2.3-004` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.5.1.1-003` | yes | yes | yes | no | yes | yes | no | yes | unknown / untested / applicable |
| not-applicable | `TOSCA13-8.5.1.1-033` | yes | no | no | yes | yes | yes | yes | yes | unknown / untested / not-applicable |
| not-applicable | `TOSCA13-8.5.1.3-005` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-8.5.1.3-024` | yes | yes | yes | no | yes | yes | no | yes | implemented / indirectly-tested / applicable |
| not-applicable | `TOSCA13-4.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / untested / not-applicable |
| ambiguous | `TOSCA13-3.7.1.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.10.3.1-002` | yes | yes | yes | yes | yes | yes | no | yes | unknown / untested / applicable |
| ambiguous | `TOSCA13-3.7.4.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.11.1-004` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.10.3.1.1-001` | yes | yes | no | no | yes | yes | no | yes | unknown / untested / not-applicable |
| ambiguous | `TOSCA13-3.7.1.1-005` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.11.4-002` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.10.3.1-001` | yes | yes | no | yes | yes | yes | no | no | partial / untested / applicable |
| ambiguous | `TOSCA13-3.7.3.1.1-003` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-6.2-002` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.7.1.1-002` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.11.4-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.10.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.11.1-002` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-6.3-008` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.7.4.1-006` | yes | yes | yes | yes | yes | yes | yes | yes | implemented / indirectly-tested / ambiguous |
| ambiguous | `TOSCA13-6.1-001` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-3.7.3.1.1-005` | yes | yes | yes | yes | yes | yes | yes | yes | unknown / unverifiable / ambiguous |
| ambiguous | `TOSCA13-5.9.10.1-004` | yes | yes | yes | yes | yes | yes | no | yes | unknown / untested / applicable |
| ambiguous | `TOSCA13-5.9.10.1-003` | yes | yes | yes | yes | yes | yes | no | yes | unknown / untested / applicable |

## Methodological corrections

- `partial` is no longer used as a proxy for missing tests.
- Keyword-to-filename matches are discarded unless the exact reader field and generic validation route are traced.
- Positive examples produce `indirectly-tested`, never `verified`.
- Duplicate, umbrella, informative, and malformed records remain visible but are excluded from metric denominators.
- Static section 8 grammar and normative types are processor-applicable; runtime fulfillment remains orchestrator-only.

## Remaining manual review

- Counted SHOULD/MAY and non-atomic records whose `implementation_status` remains `unknown`; the applicable atomic MUST scope has no remaining unknown implementation status.
- All prose semantics that cannot be reduced to an exact grammar-field or bundled-profile comparison.
- Ambiguities linked from `ambiguities.md`, including the missing §5.4.9.3 normalization algorithm and incorporated TOSCA 1.0 CSAR syntax.
- Catalog candidates outside the sample: Notes prose, split list constraints, author-only obligations, and normative-profile description text.
