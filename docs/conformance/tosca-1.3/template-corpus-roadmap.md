# TOSCA 1.3 template conformance corpus roadmap

## Scope and normative source

The corpus applies only to service templates whose
`tosca_definitions_version` is `tosca_simple_yaml_1_3`. Its sole normative
source is the pinned OASIS Standard:

`docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html`

TOSCA 2.0 cases are version-isolation regressions and are never normative
evidence for TOSCA 1.3. Orchestrator, generator, and deployment behavior is
excluded unless an applicable frozen processor requirement requires the
processor to validate or preserve the construct.

## Frozen baseline

- Atomic MUST denominator: 223.
- Applicable section identifiers represented by that denominator: 111.
- Implemented and directly verified: 223/223.
- Partial, missing, and non-compliant: 0.
- The requirements catalog, denominator membership, specification, and
  applicability classifications are immutable for this work.

## Completeness model

“Exhaustive” means exhaustive over identified normatively distinct equivalence
classes, not over arbitrary user names, scalar values, or unbounded nesting.
Every applicable construct is partitioned by:

1. grammar notation;
2. required, optional, unknown, misspelled, and cross-version keynames;
3. accepted YAML kind and distinct rejected reader paths;
4. cardinality boundaries;
5. defaults, explicit values, null, empty, and boundary values;
6. inheritance and refinement contracts;
7. namespace, import, and reference contexts;
8. intrinsic-function contexts and argument shapes;
9. constraint operand and assignment boundaries;
10. normalization equivalence, determinism, and semantic preservation;
11. CSAR entry-selection and metadata forms;
12. normatively significant cross-feature interactions.

Independent axes use documented pairwise combinations. A fixture is not
duplicated merely to vary an arbitrary name or a value within the same
processor path.

## Corpus contract

Every executable fixture has one stable manifest entry with:

- a unique `TOSCA13-CORPUS-*` identifier;
- at least one pinned specification section;
- at least one frozen requirement identifier;
- valid or invalid classification and the expected processor phase;
- at least one variation axis;
- diagnostic path/category evidence for invalid cases;
- semantic or normalization assertions when acceptance alone is insufficient.

Invalid cases have the nearest valid sibling that differs only in the target
characteristic. Supporting imported files are kept beneath `_support`
directories and cannot be mistaken for independently executable cases.

## Verification groups and commit order

| Order | Group | Primary processor path | Planned corpus areas |
|---:|---|---|---|
| 0 | roadmap | documentation | scope, completeness, grouping |
| 1 | runner and integrity | corpus infrastructure | manifest schema, case filters, integrity and frozen gates |
| 2 | namespaces and imports | read, namespaces, lookup | §3.1, local/imported/qualified/transitive/collision/cycle |
| 3 | data and common definitions | read, rendering | §§3.3, 3.5, scalar kinds, collection entries, constraints |
| 4 | grammar definitions | read, hierarchy, inheritance | §3.6 definitions, required/optional keys, notation, refinement |
| 5 | type definitions | hierarchy, inheritance | §3.7 effective types and refinement contracts |
| 6 | template definitions | rendering, resolution | §§3.8, 3.10 assignments, references, substitution, workflows |
| 7 | intrinsic functions | read, resolution, evaluation | §4 functions, contexts, arguments, nesting |
| 8 | normative profile | implicit profile, rendering | §§5, 8 normative data/capability/relationship/node semantics |
| 9 | CSAR | archive selection and read | §§6.1–6.3 metadata, root fallback, paths, versions |
| 10 | interactions | cross-phase | required pairwise feature combinations |
| 11 | normalization and isolation | normalization | equivalent forms, determinism, functions, TOSCA 2.0 rejection |
| 12 | fuzz seeds | parser safety | malformed YAML, nesting, functions, imports, constraints |
| 13 | CI gate | all phases | manifest, section, pairwise, frozen matrix, offline full run |

Production changes are prohibited in corpus commits. A discovered discrepancy
is first committed as a failing fixture with a matrix correction, then fixed
in a separate remediation commit before corpus construction continues.

## Coverage artifacts

The completed corpus owns these deterministic artifacts:

- `tests/corpus/tosca_1_3/manifest.yaml`;
- `docs/conformance/tosca-1.3/template-corpus-coverage.yaml`;
- `docs/conformance/tosca-1.3/template-corpus-sections.yaml`;
- `docs/conformance/tosca-1.3/template-corpus-pairwise.yaml`.

The section inventory classifies every normative catalog section. Each
processor-applicable frozen section must be complete and reference executable
fixtures. Excluded sections must state one of: definition-only,
duplicate-rule, orchestrator-only, author-only, covered-by-referenced-grammar,
or informative/example.

## Required gates

Each group runs its focused corpus selection and the manifest integrity tests.
Before every group commit:

```text
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
git diff --check
```

The final gate additionally checks deterministic corpus generation, all
section and pairwise classifications, offline execution, the frozen
223/223 implemented/verified matrix, and TOSCA 2.0 regression behavior.
