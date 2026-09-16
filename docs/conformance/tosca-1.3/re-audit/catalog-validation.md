# Independent catalog validation

First pass: enumerate all numbered HTML headings, inspect selected normative prose and conformance clauses, mechanically extract keyname/parameter/profile table fields independently, then add source-selected prose obligations. Existing catalogs, decisions and tests were not opened.

Second-pass checks compared selected modal prose, grammar blocks and cross-references against that extraction; the requested exhaustive second pass is not complete. The candidate catalog deliberately does not equate lexical matches, table cells, or profile type umbrella records with a final atomic denominator. YAML-only profile fields and conflicting table/grammar alternatives need reconciliation; the exact independently validated MUST total is therefore **undetermined**, not 223. This is a limitation of this audit evidence, not itself proof of a production defect.

Clause 14.3(a) imports service-template validity through 14.2, including section 4 grammar and section 5 definitions. Clause 14.3(b) imports section 3 semantics. Runtime function execution and archive deployment are distinct under 14.4. Optional fields have conditional type obligations when supplied; optional presence is never counted as required. Inherited required keynames are subject to 3.5.1.

The section inventory includes all 1121 numbered h1–h6 headings, not a preselected subset. Example ancestry is excluded except that the malformed 3.8.2.3 heading itself contains a normative range restriction. Stale references in clause 14 are recorded as ambiguities. No production interpretation was changed.

The extraction snapshot hash recorded in catalog-diff.yaml establishes the pre-comparison checkpoint. Remaining candidate atomicity and grammar reconciliation issues must be resolved before replacing any frozen catalog.

A later HTML-parser check also examined literal `<li>` elements (13 additional blocks). It recovered operand and parameter constraint compatibility sentences at HTML lines 10491 and 13799. This occurred after catalog comparison began and is therefore explicitly supplemental, not retroactively presented as part of the blind first pass. The original 642-candidate checkpoint is unchanged. The independent catalog has not completed a certified field-by-field reconciliation; no final replacement denominator is claimed.

Additional audit self-checks: the HTML inventory does not assign separate identifiers to unnumbered appendix headings; appendix content following 14.6 must not be treated as obligations of 14.6. The modal-candidate extractor skipped some paragraphs in sections already containing candidates, so section coverage is not proof of prose completeness. None of these inventory limitations is used as proof that the old implementation violates a requirement.
