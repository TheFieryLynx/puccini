# TOSCA 1.3 NetworkInfo and PortInfo property requiredness

## Status

Accepted for the profile comparison; processor behavior remains unchanged.

## Context

TOSCA Simple Profile in YAML 1.3 §§5.3.8.1 and 5.3.9.1 list the
`NetworkInfo` and `PortInfo` properties without a `Required` column. Their
Definition subsections omit `required`, which ordinarily means `true` for a
property definition under §3.6.10.3.

However, §§5.3.8.4 and 5.3.9.4 say that these properties “may or may not be
required depending on usage context.” A profile declaration with
`required: true` cannot be refined to optional under §3.6.10.6, so the literal
default prevents the usage-dependent optional case described by the same
normative type sections.

The OASIS community profile omits `required`. Puccini declares
`required: false`.

## Decision

For this audit, `required: false` is the internally consistent machine-readable
representation. It permits a usage-specific refinement to make a property
required without permitting a derived definition to relax an already-required
property.

This is an implementation interpretation, not an assertion that the
specification explicitly prints `required: false`.

## Consequences

- Puccini's representation is classified as semantically equivalent to the
  selected prose interpretation.
- The community profile's omission is recorded as part of a prose
  specification ambiguity rather than silently called authoritative.
- A future OASIS erratum or clarification can replace this decision without
  changing the frozen 223-requirement catalog.
