# TOSCA 1.3 normative profile table and Definition precedence

## Status

Accepted for the profile comparison; no processor behavior is changed by this
record.

## Context

Several normative profile sections contain a structured property or attribute
table followed by a YAML Definition table. The two representations are not
always consistent:

- §5.9.6.1 marks `WebApplication.context_root` optional, while §5.9.6.2 omits
  `required: false`.
- §5.9.8.1 marks `Database.port` optional, while §5.9.8.2 omits
  `required: false`.
- §5.7.1.1 supplies the Relationship Root `state` attribute and its `initial`
  default, while §5.7.1.2 omits that attribute.
- §8.5.2.2 supplies the network Port `ip_address` attribute, while §8.5.2.3
  omits it.

The default requiredness rule in §3.6.10.3 would make the first two properties
required if the Definition snippets were read literally.

## Decision

For field inventories, an explicit row in a normative Properties or Attributes
table supplements the YAML Definition and takes precedence over an omitted
field or omitted `required: false`. A directly contradictory explicit value is
not resolved by this rule and remains an ambiguity.

This interpretation preserves all explicit declarations in the normative
section and treats the YAML blocks as machine-readable summaries rather than
as permission to discard fields from the adjacent normative tables.

## Consequences

- `WebApplication.context_root` and `Database.port` are treated as optional.
- Relationship Root `state` and network Port `ip_address` are part of their
  normative source definitions.
- The explicit `order.required` contradiction in §§8.5.2.1 and 8.5.2.3 remains
  governed by
  `docs/decisions/0008-tosca-1.3-network-port-order-requiredness.md`.
- The `ObjectStorage.maxsize` constraint conflict remains unresolved as
  `TOSCA13-AMB-012`.
