// Package grammars dispatches service templates to one of two closed,
// version-specific grammars: TOSCA Simple Profile in YAML 1.3 or TOSCA 2.0.
//
// Dispatch uses the mandatory tosca_definitions_version keyname. Unknown
// values, aliases, and other selector keynames are rejected in the read phase
// without fallback.
package grammars
