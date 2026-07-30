# TOSCA 1.3 Group Root interface contradiction

## Status

Accepted.

## Applicable specification

TOSCA Simple Profile in YAML Version 1.3.

## Context

Section 5.10.1.1, “Definition”, shows `tosca.groups.Root` with an
`interfaces.Standard` definition whose type is
`tosca.interfaces.node.lifecycle.Standard`.

This declaration conflicts with the Group Type grammar:

- §3.7.11.1, “Keynames”, exhaustively lists the additional recognized
  keynames as `attributes`, `properties`, and `members`;
- §3.7.11.2, “Grammar”, likewise has no `interfaces` keyname;
- §5.10.1.2, “Notes”, says that future specification versions will create
  Group Type subtypes describing how Group Type operations are orchestrated.

The pinned OASIS community profile and the Puccini profile both omit the
interface.

## Decision

Strict TOSCA 1.3 processing follows the explicit Group Type keyname table and
grammar. `interfaces` is not accepted on a TOSCA 1.3 Group Type, and
`tosca.groups.Root` does not acquire the contradictory interface declaration
from §5.10.1.1.

The §5.10.1.1 declaration is recorded as a prose specification ambiguity and
an obsolete declaration, not as a Puccini or OASIS community-profile defect.

## Consequences

- No production behavior or bundled profile declaration is changed.
- The three-way inventory retains the conflicting source declaration and
  links it to this decision.
- A future compatibility mode could preserve the §5.10.1.1 form, but strict
  TOSCA 1.3 behavior must not silently admit a keyname prohibited by the
  grammar.
- TOSCA 2.0 behavior is unaffected.
