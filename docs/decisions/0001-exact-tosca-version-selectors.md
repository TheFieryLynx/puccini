# Exact TOSCA version selectors

## Status

Accepted for this fork.

## Context

TOSCA Simple Profile in YAML 1.3, section 3.10.3.1,
“tosca_definitions_version”, requires the service template to identify the
grammar and gives both the short value `tosca_simple_yaml_1_3` and a fully
qualified URI as examples.

TOSCA 2.0, section 6.2, “TOSCA Definitions Version”, defines
`tosca_2_0` as the version selector for TOSCA 2.0.

The project support boundary requires exactly two accepted selector strings
and explicitly excludes aliases, legacy names, and URI forms. This is a
project support decision; it is not presented as a TOSCA 1.3 requirement that
processors reject the URI example.

## Decision

The grammar dispatcher accepts only:

- `tosca_simple_yaml_1_3`
- `tosca_2_0`

All other strings, including the TOSCA 1.3 URI examples, are rejected during
the read phase with a deterministic unsupported-version diagnostic. Missing
or non-string selectors are reported separately as structural grammar errors.
No compatibility or fallback path is retained.

## Consequences

- Version selection is closed and deterministic.
- Removed versions and dialects cannot be selected through aliases.
- TOSCA 1.3 URI-form selectors are outside this fork's supported input
  contract even though the specification includes them as examples.
- Conformance reporting must distinguish this project support restriction
  from normative TOSCA grammar requirements.
