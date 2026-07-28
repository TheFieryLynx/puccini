# TOSCA 1.3 specification ambiguities

These are specification-source issues, not claims about Puccini behavior. Catalog records preserve the conflicting or incomplete rules and link back here through their `notes` fields.

Recorded ambiguities: **16**.

## TOSCA13-AMB-001: Conformance clause 14.3(d) cites the wrong section number

- Related sections: 14.3, 3.1, 3.2, 3.3, 3.6
- Issue: Clause 14.3(d) calls section 3.2 “Parameter and property type”, but section 3.2 is “Using Namespaces”; “Parameter and property types” is section 3.3.
- Possible interpretations:

  - Treat the parenthetical title as controlling and read the citation as section 3.3.
  - Treat the numeric citation as controlling and include namespace errors from section 3.2.
  - Include both sections until an erratum resolves the conflict.

- Catalog treatment: The catalog includes error rules from sections 3.1, 3.2, 3.3, and 3.6.
- Manual recheck required: yes

## TOSCA13-AMB-002: Conformance clause 14.3(e) cites a non-existent section

- Related sections: 14.3, 5.4.9.3
- Issue: Clause 14.3(e) requires string normalization described in section 5.4.9.3, but the published document has no section 5.4.9.3 and no other occurrence describing string normalization.
- Possible interpretations:

  - The requirement is normative but its referenced algorithm is missing.
  - The citation is a stale reference to text removed or renumbered before publication.

- Catalog treatment: Record the processor normalization obligation as ambiguous and algorithm-unspecified.
- Manual recheck required: yes

## TOSCA13-AMB-003: Implementation definitions contain stale service-template cross-references

- Related sections: 1.4, 3.8, 3.9, 3.10
- Issue: Section 1.4 points to section 3.9 for Service Template definition and section 3.8 for Topology Template definition; the actual sections are 3.10 and 3.9 respectively.
- Possible interpretations:

  - Follow the section titles and current numbering.
  - Treat the references as historical numbering with no semantic effect.

- Catalog treatment: Use sections 3.10 and 3.9 while retaining this discrepancy.
- Manual recheck required: no; retained for traceability

## TOSCA13-AMB-004: Requirement relationship type is both required and optional

- Related sections: 3.7.3.1.1
- Issue: The grammar table marks relationship keyname `type` as Required “yes” while its description calls the keyname optional.
- Possible interpretations:

  - The Required column controls and `type` is mandatory.
  - The prose controls and the relationship may use other supported notation.

- Catalog treatment: Emit both atomic records and flag their conflict; do not select an interpretation.
- Manual recheck required: yes

## TOSCA13-AMB-005: Artifact type properties conflict between Required column and prose

- Related sections: 3.7.4.1
- Issue: The `mime_type` and `file_ext` rows are marked Required “no” but their descriptions call the properties required.
- Possible interpretations:

  - Treat both properties as optional according to the table marker.
  - Treat both properties as mandatory according to the description.

- Catalog treatment: Emit both table-marker and description records and flag the conflict.
- Manual recheck required: yes

## TOSCA13-AMB-006: Parameter `type` requirement is context-dependent

- Related sections: 3.6.10.2, 3.6.14.1
- Issue: The parameter table marks `type` optional while describing it as required for property definitions but not parameter definitions; the shared schema is reused in multiple contexts.
- Possible interpretations:

  - Require `type` only when the schema is used as a property definition.
  - Allow omission in parameter contexts and derive the type where another rule permits it.

- Catalog treatment: Record the conditional rules independently.
- Manual recheck required: no; retained for traceability

## TOSCA13-AMB-007: Version selector short name and URI forms

- Related sections: 3.10.3.1
- Issue: The section presents both `tosca_simple_yaml_1_3` and a fully qualified URI as examples but does not state an exhaustive lexical set.
- Possible interpretations:

  - Accept both forms.
  - Accept any profile identifier that unambiguously selects the grammar.
  - Limit a product support profile to the short form while documenting that restriction.

- Catalog treatment: Catalog the explicit grammar-selection duty without inventing an exhaustive alias rule.
- Manual recheck required: yes

## TOSCA13-AMB-008: TOSCA.meta syntax is incorporated from an unavailable normative source

- Related sections: 6.2
- Issue: Section 6.2 says the TOSCA.meta syntax is exactly that of TOSCA 1.0 section 16.2, but the allowed local source does not reproduce that syntax.
- Possible interpretations:

  - The local 1.3 document normatively incorporates the external syntax.
  - Only the explicit 1.3 deltas can be cataloged without consulting another normative source.

- Catalog treatment: Record the incorporation requirement and all explicit 1.3 deltas; mark the inherited syntax for manual review.
- Manual recheck required: yes

## TOSCA13-AMB-009: CSAR root alternatives do not state whether they are exclusive

- Related sections: 6.1, 6.3
- Issue: Section 6.1 requires one of two layouts, but does not explicitly say whether an archive containing both TOSCA-Metadata and a root YAML file is invalid.
- Possible interpretations:

  - “One of” is exclusive.
  - Presence of TOSCA-Metadata takes precedence and a root YAML file is merely additional content.

- Catalog treatment: Record both valid alternatives and flag exclusivity as unspecified.
- Manual recheck required: yes

## TOSCA13-AMB-010: Archive-without-metadata derives CSAR version from template version

- Related sections: 6.3, 3.10.1.1
- Issue: Section 6.3 says CSAR-Version is defined by `template_version`, although template version identifies the template rather than the archive format.
- Possible interpretations:

  - Use `template_version` literally as CSAR-Version.
  - Treat the sentence as an editorial error and use the fixed CSAR format version 1.1.

- Catalog treatment: Preserve the explicit sentence and mark the semantic mismatch.
- Manual recheck required: yes

## TOSCA13-AMB-011: Other-Definitions filename tokenization is incomplete

- Related sections: 6.2
- Issue: Space-delimited filenames and double-quoted filenames are defined, but escaping quotes, backslashes, repeated whitespace, and empty entries is unspecified.
- Possible interpretations:

  - Use a shell-like quoted token grammar.
  - Implement only the stated blank-space delimiter and simple double quoting.

- Catalog treatment: Catalog only the stated delimiter and quoting rules.
- Manual recheck required: yes

## TOSCA13-AMB-012: ObjectStorage `maxsize` lower bound conflicts

- Related sections: 5.9.10.1, 5.9.10.3
- Issue: The property table gives `greater_or_equal: 1 GB`; the normative type definition gives `greater_or_equal: 0 GB`.
- Possible interpretations:

  - Use the table constraint of 1 GB.
  - Use the definition constraint of 0 GB.

- Catalog treatment: Emit both constraints and flag the conflict.
- Manual recheck required: yes

## TOSCA13-AMB-013: BlockStorage `size` has conditional requiredness and `volume_id` precedence

- Related sections: 5.9.11.1, 5.9.11.3, 5.9.11.4
- Issue: The table marks `size` as “yes *”, then makes it conditional on `volume_id`; the inherited storage definition also supplies a default.
- Possible interpretations:

  - Require `size` only when `volume_id` is absent.
  - Allow the inherited default to satisfy the condition.

- Catalog treatment: Record conditional requiredness, precedence, and inherited default separately.
- Manual recheck required: yes

## TOSCA13-AMB-014: Entity Type Schema derived_from constraint is malformed

- Related sections: 3.7.1.1
- Issue: The `derived_from` constraints cell states “None is the only allowed value” while the description says it contains an optional parent type name.
- Possible interpretations:

  - “None” means there is no additional constraint.
  - Only a root entity may omit or null the parent.

- Catalog treatment: Record the cell and prose independently and flag the contradiction.
- Manual recheck required: yes

## TOSCA13-AMB-015: Examples are non-normative but some normative rules depend on example-only detail

- Related sections: 1.6.1, 3, 4, 5, 6
- Issue: The specification declares Example sections non-normative, while some prose and grammar descriptions refer to examples for concrete argument forms or value interpretation.
- Possible interpretations:

  - Exclude every rule that exists only in an Example section.
  - Use examples only to illustrate a rule independently stated in normative text.

- Catalog treatment: Example-scoped content is excluded and never used as the sole source for a requirement.
- Manual recheck required: no; retained for traceability

## TOSCA13-AMB-016: Lowercase modal language has uncertain normative force

- Related sections: 3, 4, 5, 6, 7, 8, 12, 13
- Issue: The document mixes uppercase RFC 2119 terms with lowercase words such as “should”, “may”, and “required” in grammar descriptions.
- Possible interpretations:

  - Only uppercase terms carry RFC 2119 strength.
  - Grammar tables and unambiguous declarative prose remain binding even without uppercase terms.

- Catalog treatment: Uppercase terms retain their exact strength; grammar markers use GRAMMAR/REQUIRED/OPTIONAL, unambiguous required/optional declarations use REQUIRED/OPTIONAL, and other declarative prose uses SEMANTIC.
- Manual recheck required: no; retained for traceability
