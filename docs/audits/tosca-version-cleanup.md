# TOSCA grammar and profile cleanup audit

Date: 2026-07-27

This inventory was completed before deletion. The dependency checks used Go
imports, reader registrations, embedded-asset references, implicit-import
paths, example entry points, and the baseline `go test ./...` result. The
normative basis was limited to the pinned TOSCA 1.3 and TOSCA 2.0 standards.

## Normative scope

| Classification | Normative source |
| --- | --- |
| `keep-1.3` | TOSCA Simple Profile in YAML 1.3, sections 3.10, 3.10.1, 3.10.2, 3.10.3.1, and 14.2-14.3 |
| `keep-2.0` | TOSCA 2.0, sections 6, 6.1, 6.2, 6.9, and 6.9.1 |
| `keep-shared` | Infrastructure reached by either retained grammar |
| `remove` | A different language version, dialect, or profile, or code reached only by one |
| `uncertain` | A name suggested legacy ownership but reachability required inspection |

TOSCA 1.3 section 3.10.1 makes `tosca_definitions_version` required and
section 3.10.3.1 identifies it as the grammar selector. TOSCA 2.0 sections 6.1
and 6.2 make it mandatory, require it to select the grammar for the remainder
of the file, and identify `tosca_2_0` as this specification's version.

## Registered grammars and version values

| Package | Registered selector and values | Classification | Evidence and action |
| --- | --- | --- | --- |
| `tosca/grammars/tosca_v1_3` | `tosca_definitions_version: tosca_simple_yaml_1_3` | `keep-1.3` | Retained as an independent grammar. |
| `tosca/grammars/tosca_v2_0` | `tosca_definitions_version: tosca_2_0` | `keep-2.0` | Retained as an independent grammar. Readers reused by 1.3 remain shared implementation, not a version alias. |
| `tosca/grammars/tosca_v1_0` | `tosca_simple_yaml_1_0` | `remove` | Separate obsolete grammar package (12 files). |
| `tosca/grammars/tosca_v1_1` | `tosca_simple_yaml_1_1` | `remove` | Separate obsolete grammar package (10 files). |
| `tosca/grammars/tosca_v1_2` | `tosca_simple_yaml_1_2`, `tosca_simple_profile_for_nfv_1_0` | `remove` | Separate obsolete grammar/profile package (23 files). |
| `tosca/grammars/cloudify_v1_3` | `cloudify_dsl_1_3` | `remove` | Non-standard dialect package (34 files). |
| `tosca/grammars/hot` | `wallaby`, `train`, `stein`, `rocky`, `queens`, `pike`, `newton`, `ocata`, `2021-04-16`, `2018-08-31`, `2018-03-02`, `2017-09-01`, `2017-02-24`, `2016-10-14`, `2016-04-08`, `2015-10-15`, `2015-04-30`, `2014-10-16`, `2013-05-23` under `heat_template_version` | `remove` | Non-TOSCA HOT dialect package (14 files). |

No other version URI, legacy alias, or alternative language selector was
registered.

## Selection, detection, registry, and public API

| Component | Classification | Reachability and decision |
| --- | --- | --- |
| `tosca/grammars/init.go` | `keep-shared` | Central registry initializer; reduce imports and registrations to 1.3 and 2.0. |
| `tosca/grammars/common.go` (`Grammars`, `ImplicitProfilePaths`) | `keep-shared` | Used by read-phase dispatch and implicit imports; replace mutable exported version registry with a closed two-version registry. |
| `tosca/grammars/parse.go` (`DetectGrammar`, `GetGrammar`, `DetectGrammarVersion`, `GetImplicitImportSpec`, `CompatibleGrammars`) | `keep-shared` | Read/import phase infrastructure. Remove keyword scanning and the HOT timestamp hack; require the exact TOSCA selector. |
| `tosca/parsing/grammars.go` (`GrammarVersion`, `GrammarVersions`, `RegisterVersion`) | `remove` in part | The public version-registration API exists only to make the registry extensible. Remove it; retain `Grammar`, `NewGrammar`, and reader registration. |
| `tosca/parser/phase1-read.go` | `keep-shared` | YAML decoding, grammar detection, implicit imports, and read-phase validation are needed by both versions. |
| `tosca/parser/context.go` grammar compatibility check | `keep-shared` | Needed to reject cross-version imports. |
| CLI compile/parse/validate `--quirk` option | `keep-shared` | No CLI option selects a grammar directly; the document selector does. Retain only quirks meaningful for 1.3/2.0. |
| CSAR `--tosca-meta-file-version` and `--csar-version` | `keep-shared` | These select CSAR metadata format fields, not a TOSCA YAML grammar. They remain necessary for TOSCA 1.3 archive handling. |

## Built-in profiles and definitions

| Path | Classification | Evidence and action |
| --- | --- | --- |
| `assets/tosca/profiles/simple/1.3` | `keep-1.3` | Implicit normative type profile for the 1.3 grammar. |
| `assets/tosca/profiles/implicit/1.3` | `keep-1.3` | Imported by the retained 1.3 profile. |
| `assets/tosca/profiles/simple/2.0` | `keep-2.0` | Backing data for the retained `org.oasis-open.simple:2.0` well-known profile. Its stale 1.x `valid_target_types` and `scalar-unit.*` spellings were adapted to the 2.0 relationship and scalar type grammar. |
| `assets/tosca/profiles/implicit/2.0` | `keep-2.0` | Implicit 2.0 data, functions, validations, and the 2.0 implicit profile. |
| `assets/tosca/profiles/common/1.0` | `uncertain` -> `keep-shared` | Despite its directory name, both retained profile files import its Clout resolution/coercion/output scriptlets. |
| `assets/tosca/profiles/simple/1.0`, `1.1`, `1.2` | `remove` | Reached only by removed grammars. |
| `assets/tosca/profiles/implicit/1.0`, `1.1`, `1.2` | `remove` | Reached only by removed profiles. |
| `assets/tosca/profiles/simple-for-nfv/1.0` | `remove` | Reached only by the removed NFV selector. |
| `assets/tosca/profiles/cloudify/5.0.5` | `remove` | Reached only by the Cloudify grammar. |
| `assets/tosca/profiles/hot/1.0` | `remove` | Reached only by the HOT grammar. |
| `assets/tosca/profiles/profiles.go` | `keep-shared` | Embedded internal URL provider; narrow its embed patterns to retained assets. |

## Tests, examples, fixtures, and benchmarks

| Component | Classification | Decision |
| --- | --- | --- |
| `puccini_test.go` 1.3, 2.0, current JavaScript, CSAR-related and custom-type examples | `keep-shared` | Retain coverage and normalization/Clout execution for both target versions. |
| `puccini_test.go` legacy, Cloudify, HOT, and NFV invocations | `remove` | They prove only removed selectors/dialects. Replace legacy successes with rejection tests. |
| `examples/1.3` except `simple-for-nfv.yaml` | `keep-1.3` | Full TOSCA 1.3 examples. |
| `examples/2.0` | `keep-2.0` | Full TOSCA 2.0 examples; remove stale NFV documentation link. |
| `examples/legacy` | `remove` | TOSCA 1.0/1.1/1.2 fixtures only. |
| `examples/cloudify` | `remove` | Cloudify dialect fixtures only. |
| `examples/hot` | `remove` | HOT fixtures only. |
| `examples/openstack`, `examples/bpmn`, `examples/ansible`, `examples/csar`, `examples/javascript` | `uncertain` -> `keep-shared` | Their entry templates declare 1.3 or 2.0 and exercise standard custom types, imports, CSAR, functions, wrappers, or Clout processing; none registers another grammar. |
| `examples/dgraph`, `examples/neo4j`, language wrapper examples | `keep-shared` | Consumers of normalized/Clout output, independent of removed grammars. |
| `BenchmarkParse` | `keep-shared` | Runs the retained `compileAll` corpus after obsolete entries are removed. |
| Generated grammar code | `not applicable` | No generated grammar files or grammar code-generation pipeline was found. |

## Documentation

| Component | Classification | Decision |
| --- | --- | --- |
| `README.md` supported-version and dialect claims | `remove`/`keep-shared` | Replace with an exact two-version support statement and no dialect claims. |
| `examples/README.md` HOT/Cloudify links | `remove` | Remove stale dialect categories; list target-version examples. |
| `examples/1.3/README.md`, `examples/2.0/README.md` NFV links | `remove` | Remove the deleted NFV example references. |
| `scripts/csar` references to older TOSCA CSAR specifications | `uncertain` -> `remove` documentation only | Runtime is shared, but comments must cite only the retained 1.3/2.0 sources. |
| TOSCA parsing architecture/package comments | `keep-shared` | Update to describe closed two-version dispatch. |
| Historical cleanup changelog entry | `keep-shared` | Historical mentions are retained only in the changelog and this audit as an explicit record of removal. |

## Compatibility and quirks

| Component | Classification | Decision |
| --- | --- | --- |
| HOT YAML timestamp branch in `DetectGrammarVersion` | `remove` | Dialect-only detection hack. |
| `interfaces.operations.permissive` | `remove` | Explicitly accepts the removed 1.2 grammar in 1.3/2.0. |
| 1.3 extended trigger-condition reader registered in the 2.0 grammar | `remove` from 2.0 / `keep-1.3` | Move the reader to the 1.3 package and register it only there; 2.0 section 16.5 accepts a condition clause without the 1.3 `constraint`, `period`, `evaluations`, or `method` wrapper. |
| Scalar-unit zero-value bridge in the 2.0 value implementation | `uncertain` -> `keep-shared` | The 1.3 grammar delegates common value rendering to the 2.0 implementation; the bridge is reached by retained 1.3 scalar-unit types and is not a selector or alias accepted by 2.0. |
| Shared mapping and node-filter representation | `uncertain` -> `keep-shared` | Reached by both retained grammars through version-specific root/readers. It represents 1.3 and 2.0 normalized concepts; exhaustive feature isolation remains tracked as partial in the conformance matrix rather than being mistaken for removed-version code. |
| Other primitive, namespace, import, occurrence, and substitution quirks | `keep-shared` | They are executed from retained 1.3/2.0 code paths and are not selectors or packages for a removed version. |
| `etsinfv`/`onap` combination names | `uncertain` -> `remove` | These convenience aliases represented non-standard dialect compatibility. The underlying generic parser-policy toggles remain available where they operate on retained grammar code paths. |

## Dependencies and dead code

The direct and transitive import graphs for removed packages and retained
parser/normalization packages were compared before deletion. No `go.mod`
dependency was exclusive to a removed grammar: removed packages import the
same ARD, parsing, logging, URL, normalization, and JavaScript infrastructure
used by retained code. `go mod tidy` after deletion is therefore expected to
leave the module dependency set unchanged; any actual change will be recorded
in the final verification.

The version-registration types, mutable registry API, HOT timestamp branch,
obsolete embed alternatives, dialect combination aliases, the legacy
interface-operation quirk, and the 2.0 registration of the 1.3 extended
trigger-condition reader become dead or invalid after package removal and are
included in the shared cleanup.

## Verification finding: retained 2.0 profile

Replacing external community-profile URLs in the 2.0 examples with the
built-in well-known profile exposed stale syntax in that retained asset.
TOSCA 2.0 section 7.3 recognizes `valid_capability_types`,
`valid_target_node_types`, and `valid_source_node_types`, not the 1.x
`valid_target_types` keyname. Section 9.1.2.2.2 gives informative examples of
`Bitrate`, `Frequency`, `Size`, and `Time` as `scalar`-derived replacements for
the former `scalar-unit.*` types. The built-in profile now uses those 2.0
forms. The grammar was not relaxed.
