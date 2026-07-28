# Confirmed TOSCA 1.3 non-compliance

Confirmed non-compliant catalog records: **33**. Duplicate records are shown individually but excluded from atomic coverage denominators.
Every behavioral claim below was reproduced on 2026-07-27 or directly compared against the bundled profile. Temporary probes were removed.

## TOSCA13-3.10.2.1-003

- Normative basis: §3.10.2.1 **Requirements** — The key “tosca_definitions_version” SHOULD be the first line of each Service Template.
- Minimal reproducer: Place `description` before `tosca_definitions_version` in a minimal service template.
- Actual Puccini result: Compilation succeeded without diagnostics.
- Expected result: Reject because the selector is not the first YAML line.
- Responsible path: `YAML decode → grammars.DetectGrammar → parser phase 1 read`
- Fix phase: **read/lexical**
- Required regression test: Direct negative first-line test asserting a lexical-position diagnostic.

## TOSCA13-3.6.3.1-004

- Normative basis: §3.6.3.1 **Operator keynames** — For `equal`: Constrains a property or parameter to a value equal to (‘=’) the value declared.
- Minimal reproducer: input constraints: [{ equal: 7 }]; default: 8
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `equal`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-007

- Normative basis: §3.6.3.1 **Operator keynames** — For `greater_than`: Constrains a property or parameter to a value greater than (‘>’) the value declared.
- Minimal reproducer: input constraints: [{ greater_than: 7 }]; default: 7
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `greater_than`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-010

- Normative basis: §3.6.3.1 **Operator keynames** — For `greater_or_equal`: Constrains a property or parameter to a value greater than or equal to (‘>=’) the value declared.
- Minimal reproducer: input constraints: [{ greater_or_equal: 7 }]; default: 6
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `greater_or_equal`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-013

- Normative basis: §3.6.3.1 **Operator keynames** — For `less_than`: Constrains a property or parameter to a value less than (‘<’) the value declared.
- Minimal reproducer: input constraints: [{ less_than: 7 }]; default: 7
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `less_than`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-016

- Normative basis: §3.6.3.1 **Operator keynames** — For `less_or_equal`: Constrains a property or parameter to a value less than or equal to (‘<=’) the value declared.
- Minimal reproducer: input constraints: [{ less_or_equal: 7 }]; default: 8
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `less_or_equal`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-020

- Normative basis: §3.6.3.1 **Operator keynames** — For `in_range`: Constrains a property or parameter to a value in range of (inclusive) the two values declared.
- Minimal reproducer: input constraints: [{ in_range: [1, 3] }]; default: 4
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `in_range`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-021

- Normative basis: §3.6.3.1 **Operator keynames** — For `valid_values`: Constrains a property or parameter to a value that is in the list of declared values.
- Minimal reproducer: input constraints: [{ valid_values: [1, 2] }]; default: 3
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `valid_values`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-024

- Normative basis: §3.6.3.1 **Operator keynames** — For `length`, `value type` is `string, list, map`.
- Minimal reproducer: input type string; constraints: [{ length: 3 }]; default: abc
- Actual Puccini result: Rejected during read/rendering as unsupported operator: length.
- Expected result: Accept the valid length clause and enforce it on assignments.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `length`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-025

- Normative basis: §3.6.3.1 **Operator keynames** — The `length` entry has type/schema `scalar`.
- Minimal reproducer: input type string; constraints: [{ length: 3 }]; default: abc
- Actual Puccini result: Rejected during read/rendering as unsupported operator: length.
- Expected result: Accept the valid length clause and enforce it on assignments.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `length`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-026

- Normative basis: §3.6.3.1 **Operator keynames** — For `length`: Constrains the property or parameter to a value of a given length.
- Minimal reproducer: input type string; constraints: [{ length: 3 }]; default: abc
- Actual Puccini result: Rejected during read/rendering as unsupported operator: length.
- Expected result: Accept the valid length clause and enforce it on assignments.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `length`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-029

- Normative basis: §3.6.3.1 **Operator keynames** — For `min_length`: Constrains the property or parameter to a value to a minimum length.
- Minimal reproducer: input constraints: [{ min_length: 3 }]; default: ab
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `min_length`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-032

- Normative basis: §3.6.3.1 **Operator keynames** — For `max_length`: Constrains the property or parameter to a value to a maximum length.
- Minimal reproducer: input constraints: [{ max_length: 3 }]; default: abcd
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `max_length`; for `length`, also a valid-clause positive test.

## TOSCA13-3.6.3.1-035

- Normative basis: §3.6.3.1 **Operator keynames** — For `pattern`: Constrains the property or parameter to a value that is allowed by the provided regular expression.
- Minimal reproducer: input constraints: [{ pattern: ^a+$ }]; default: b
- Actual Puccini result: Compilation accepted the constraint-violating default without diagnostics.
- Expected result: Reject the assignment during rendering with a constraint-specific diagnostic.
- Responsible path: `ReadParameterDefinition → ReadConstraintClause/ReadValidationClause → Value.render → parser.Context.Render`
- Fix phase: **rendering**
- Required regression test: Direct negative assignment test for `pattern`; for `length`, also a valid-clause positive test.

## TOSCA13-4.3.1-001

- Normative basis: §4.3.1 **concat** — The concat function is used to concatenate two or more string values within a TOSCA service template.
- Minimal reproducer: topology output value: { concat: [] }
- Actual Puccini result: Compilation accepted and normalized the invalid function call without diagnostics.
- Expected result: Reject invalid argument count/type or unresolved input reference.
- Responsible path: `ReadValue → ParseFunctionCall → setFunctionCall → rendering/normalization`
- Fix phase: **rendering/function resolution**
- Required regression test: Direct negative `concat` grammar/resolution test asserting the specific diagnostic.

## TOSCA13-4.3.1.2-001

- Normative basis: §4.3.1.2 **Parameters** — For `<string_value_expressions_*>`: A list of one or more strings (or expressions that result in a string value) which can be concatenated together into a single string.
- Minimal reproducer: topology output value: { concat: [] }
- Actual Puccini result: Compilation accepted and normalized the invalid function call without diagnostics.
- Expected result: Reject invalid argument count/type or unresolved input reference.
- Responsible path: `ReadValue → ParseFunctionCall → setFunctionCall → rendering/normalization`
- Fix phase: **rendering/function resolution**
- Required regression test: Direct negative `concat` grammar/resolution test asserting the specific diagnostic.

## TOSCA13-4.4.1-001

- Normative basis: §4.4.1 **get_input** — The get_input function is used to retrieve the values of properties declared within the inputs section of a TOSCA Service Template.
- Minimal reproducer: topology output value: { get_input: absent }
- Actual Puccini result: Compilation accepted and normalized the invalid function call without diagnostics.
- Expected result: Reject invalid argument count/type or unresolved input reference.
- Responsible path: `ReadValue → ParseFunctionCall → setFunctionCall → rendering/normalization`
- Fix phase: **rendering/function resolution**
- Required regression test: Direct negative `get_input` grammar/resolution test asserting the specific diagnostic.

## TOSCA13-4.4.1.2-001

- Normative basis: §4.4.1.2 **Parameters** — For `<input_property_name>`: The name of the property as defined in the inputs section of the service template.
- Minimal reproducer: topology output value: { get_input: absent }
- Actual Puccini result: Compilation accepted and normalized the invalid function call without diagnostics.
- Expected result: Reject invalid argument count/type or unresolved input reference.
- Responsible path: `ReadValue → ParseFunctionCall → setFunctionCall → rendering/normalization`
- Fix phase: **rendering/function resolution**
- Required regression test: Direct negative `get_input` grammar/resolution test asserting the specific diagnostic.

## TOSCA13-4.4.1.2-002

- Normative basis: §4.4.1.2 **Parameters** — The `<input_property_name>` entry has type/schema `string`.
- Minimal reproducer: topology output value: { get_input: [] }
- Actual Puccini result: Compilation accepted and normalized the invalid function call without diagnostics.
- Expected result: Reject invalid argument count/type or unresolved input reference.
- Responsible path: `ReadValue → ParseFunctionCall → setFunctionCall → rendering/normalization`
- Fix phase: **rendering/function resolution**
- Required regression test: Direct negative `get_input` grammar/resolution test asserting the specific diagnostic.

## TOSCA13-5.10.1.1-003

- Normative basis: §5.10.1.1 **Definition** — The normative type schema declares `tosca.groups.Root.interfaces.Standard`.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/groups.yaml` and inspect `tosca.groups.root.interfaces.standard`.
- Actual Puccini result: interfaces.Standard is absent
- Expected result: interfaces.Standard is declared
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.groups.root.interfaces.standard`.

## TOSCA13-5.10.1.1-004

- Normative basis: §5.10.1.1 **Definition** — `tosca.groups.Root.interfaces.Standard` has type `tosca.interfaces.node.lifecycle.Standard`.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/groups.yaml` and inspect `tosca.groups.root.interfaces.standard.type`.
- Actual Puccini result: interfaces.Standard.type is absent
- Expected result: type: tosca.interfaces.node.lifecycle.Standard
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.groups.root.interfaces.standard.type`.

## TOSCA13-5.4.3.4.1-002

- Normative basis: §5.4.3.4.1 **Definition** — `tosca.artifacts.Deployment.Image.VM` derives from `tosca.artifacts.Deployment.Image`.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/artifacts.yaml` and inspect `tosca.artifacts.deployment.image.vm.derived_from`.
- Actual Puccini result: derived_from: tosca.artifacts.Deployment
- Expected result: derived_from: tosca.artifacts.Deployment.Image
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.artifacts.deployment.image.vm.derived_from`.

## TOSCA13-5.9.9.2-002

- Normative basis: §5.9.9.2 **Definition** — `tosca.nodes.Abstract.Storage` derives from `tosca.nodes.Root`.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/nodes.yaml` and inspect `tosca.nodes.abstract.storage.derived_from`.
- Actual Puccini result: derived_from is absent
- Expected result: derived_from: tosca.nodes.Root
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.nodes.abstract.storage.derived_from`.

## TOSCA13-5.9.9.2-007

- Normative basis: §5.9.9.2 **Definition** — `tosca.nodes.Abstract.Storage.properties.size` defaults to `0 MB`.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/nodes.yaml` and inspect `tosca.nodes.abstract.storage.properties.size.default`.
- Actual Puccini result: size.default is absent
- Expected result: size.default: 0 MB
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.nodes.abstract.storage.properties.size.default`.

## TOSCA13-6.1-004

- Normative basis: §6.1 **Overall Structure of a CSAR** — The root YAML file MUST define a metadata section.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.1-006

- Normative basis: §6.1 **Overall Structure of a CSAR** — The root YAML file metadata section requires `template_name`.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.1-007

- Normative basis: §6.1 **Overall Structure of a CSAR** — The root YAML file metadata section requires `template_version`.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.2-005

- Normative basis: §6.2 **TOSCA Meta File** — `block_0` requires the `Entry-Definitions` key.
- Minimal reproducer: TOSCA.meta 1.1 with CSAR-Version and Created-By but no Entry-Definitions
- Actual Puccini result: csar.ReadMeta returned success.
- Expected result: Reject missing Entry-Definitions.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.2-018

- Normative basis: §6.2 **TOSCA Meta File** — Due to the changes to the TOSCA.meta file compared to TOSCA 1.2 the TOSCA-Meta-File-Version keyword listed in block_0 of the the meta-file is required to denote version 1.1.
- Minimal reproducer: TOSCA.meta with TOSCA-Meta-File-Version: 1.0 and otherwise valid block_0
- Actual Puccini result: csar.ReadMeta returned success.
- Expected result: Reject; TOSCA 1.3 requires meta-file version 1.1.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.3-003

- Normative basis: §6.3 **Archive without TOSCA-Metadata** — The root entry file metadata requires `template_name`.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.3-004

- Normative basis: §6.3 **Archive without TOSCA-Metadata** — The root entry file metadata requires `template_version`.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-6.3-006

- Normative basis: §6.3 **Archive without TOSCA-Metadata** — The root entry file requires a metadata section.
- Minimal reproducer: CSAR without TOSCA-Metadata containing one root service.yaml with no metadata
- Actual Puccini result: Archive parsing succeeded with no diagnostics.
- Expected result: Reject the root template for missing metadata/template_name/template_version.
- Responsible path: `GetServiceTemplateURL → ReadMetaFromURL or GetRootPath fallback → service-template parse`
- Fix phase: **csar/read**
- Required regression test: Direct in-memory CSAR negative test asserting the missing/invalid metadata diagnostic.

## TOSCA13-8.5.2.3-008

- Normative basis: §8.5.2.3 **Definition** — `tosca.nodes.network.Port.properties.order` is required.
- Minimal reproducer: Load bundled profile `assets/tosca/profiles/simple/1.3/nodes.yaml` and inspect `tosca.nodes.network.port.properties.order.required`.
- Actual Puccini result: tosca.nodes.network.Port.properties.order.required: false
- Expected result: tosca.nodes.network.Port.properties.order.required: true
- Responsible path: `implicit/simple profile load → namespaces → hierarchy → inheritance`
- Fix phase: **implicit-profile/read or hierarchy**
- Required regression test: Direct normative-profile assertion for `tosca.nodes.network.port.properties.order.required`.

## Rejected former non-compliant claim

- `TOSCA13-5.9.12.1-004` is a catalog error: the specification declares a `host` requirement assignment with capability `tosca.capabilities.Compute`; the record incorrectly says `capabilities.host`. The bundled profile is semantically incompatible (it supplies a host capability and omits the normative requirement), but this malformed record is not retained as non-compliant.

## Normative profile classification

- The six retained profile findings above are literal bundled-profile mismatches and semantic incompatibilities.
- Omitted `derived_from` on conceptual root types is not automatically a mismatch; no finding is emitted unless the specification explicitly names a parent.
- `tosca.nodes.Storage.ObjectStorage.maxsize` remains ambiguous because §5.9.10.1 says 1 GB while §5.9.10.3's definition says 0 GB.
- `tosca.groups.Root.interfaces.Standard` is not an equivalent internal representation: both the interface and its type are absent.
