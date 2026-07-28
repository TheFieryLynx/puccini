# Changelog

## Unreleased

### Breaking

* Limit service-template parsing to TOSCA Simple Profile in YAML 1.3
  (`tosca_simple_yaml_1_3`) and TOSCA 2.0 (`tosca_2_0`).
* Remove TOSCA Simple Profile 1.0, 1.1, and 1.2, the Simple Profile for NFV
  selector, Cloudify DSL, and HOT grammar packages, profiles, examples, and
  successful-selection behavior.
* Remove the runtime grammar-registration API, non-standard dialect
  combination aliases, and legacy interface-operation permissive quirk.
* Raise the minimum Go toolchain from 1.26.0 to 1.26.1 because the
  `github.com/lestrrat-go/helium` v0.7.0 XML Schema compiler requires Go
  1.26.1.

### Validation

* Reject missing, unknown, removed, alias, and wrong-typed version selectors in
  the read phase with deterministic diagnostics.
* Add version-specific positive, normalization, full-example, cross-version,
  and removed-selector conformance tests.
* Adapt the retained TOSCA 2.0 Simple Profile asset to the 2.0 relationship
  type and scalar type grammar, without weakening the 2.0 parser.
* Move the TOSCA 1.3 extended trigger-condition reader out of the 2.0 grammar
  and add positive/negative version-isolation coverage.
