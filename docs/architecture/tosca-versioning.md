# TOSCA version dispatch architecture

Puccini supports two independent language grammars:

| `tosca_definitions_version` | Grammar package | Built-in profile entry point |
| --- | --- | --- |
| `tosca_simple_yaml_1_3` | `tosca/grammars/tosca_v1_3` | `/profiles/simple/1.3/profile.yaml` |
| `tosca_2_0` | `tosca/grammars/tosca_v2_0` | `/profiles/implicit/2.0/profile.yaml` |

## Read-phase dispatch

Phase 1 decodes YAML, requires a string-valued `tosca_definitions_version`, and
performs an exact lookup in a closed internal map. There is no keyword scan,
version URI alias, closest-version match, or fallback grammar. Unsupported
values report the value and both accepted values. A missing selector reports a
missing required keyname.

The supported-version API returns a copy of the two values; callers cannot
register or mutate grammar entries. Adding another grammar therefore requires
an intentional source change rather than runtime registry mutation.

## Version isolation

The 1.3 and 2.0 packages register separate root readers and version-specific
overrides. Shared reader implementations remain in the 2.0 package where the
existing 1.3 grammar imports them, but registration is independent: reuse of a
Go function does not make a 2.0 keyname valid in a 1.3 document. Cross-version
tests verify that `service_template` is rejected by 1.3 and
`topology_template` is rejected by 2.0.

Imports are parsed using their own mandatory version selector. The parser
compares the resolved grammar identities and rejects incompatible imports
unless a separately requested non-strict import quirk changes that policy.

## Shared infrastructure

YAML and URL reading, filesystem access, imports, namespaces, hierarchy
construction, inheritance, rendering, normalization, CSAR handling, Clout,
diagnostics, and test helpers are shared. The embedded `common/1.0` directory
name is historical, but its scriptlets are actively imported by both retained
profiles; it is shared runtime infrastructure rather than a selectable grammar.

## Normative boundaries

The dispatch rules are grounded in:

* TOSCA Simple Profile in YAML 1.3 sections 3.10, 3.10.1, 3.10.2, and 3.10.3.1.
* TOSCA 2.0 sections 6, 6.1, and 6.2.

Normalization and Clout are Puccini representations, not TOSCA specifications.
Their shared implementation does not imply that source-language semantics are
identical.
