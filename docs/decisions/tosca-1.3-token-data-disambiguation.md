# TOSCA 1.3 token data and intrinsic function syntax

The normative Credential definition (§5.3.6.1–2) permits a value containing
only `token: secret`: token is a string, token_type has default password, and
all remaining fields are optional. Intrinsic `token` (§4.3.3.1) instead takes
a YAML sequence containing the source string, separators, and an index.

A single map key matching an intrinsic function name is insufficient to
identify a function. A string-valued `token` entry cannot have the function's
argument shape and is valid normative Credential data. Preserve such maps as
data, including nested collections, and let the declared value schema validate
them. A sequence-valued token remains function syntax and receives the existing
function argument checks. A scalar property assigned a token map still fails
its own type validation during rendering.

This resolves the overlap by using the normative shapes, without changing
function semantics or importing TOSCA 2.0 rules. Other hypothetical overlaps
are outside this focused defect. The function/data precedence algorithm is an
implementation interpretation, not an explicitly specified general algorithm.
