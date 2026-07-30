#!/usr/bin/env python3
"""Completeness gate for the TOSCA 1.3 profile three-way comparison."""

from __future__ import annotations

import argparse
import copy
import hashlib
import pathlib
import sys
from collections.abc import Callable
from typing import Any

import yaml


ROOT = pathlib.Path(__file__).resolve().parents[2]
BASE = ROOT / "docs/conformance/tosca-1.3"
COMPARISON_PATH = BASE / "profile-three-way-comparison.yaml"
DISCREPANCIES_PATH = BASE / "profile-discrepancies.yaml"
MANIFEST_PATH = ROOT / "tests/corpus/tosca_1_3/oasis_community/manifest.yaml"
PINNED_ROOT = ROOT / "third_party/oasis/tosca-simple-profile/1.3"

ALLOWED_COMPARISONS = {"exact", "equivalent", "conflicting", "missing"}
ALLOWED_CLASSIFICATIONS = {
    "Puccini defect",
    "community-profile defect",
    "prose specification ambiguity",
    "semantically equivalent representation",
    "version mismatch",
    "extension",
    "obsolete declaration",
}
ALLOWED_RESOLUTIONS = {"confirmed", "documented", "interpreted", "remediated"}


def load_yaml(path: pathlib.Path) -> Any:
    return yaml.safe_load(path.read_text(encoding="utf-8"))


def validate_integrity(
    reader: Callable[[pathlib.Path], bytes] | None = None,
) -> list[str]:
    if reader is None:
        reader = pathlib.Path.read_bytes
    errors: list[str] = []
    sums_path = PINNED_ROOT / "SHA256SUMS"
    for line in sums_path.read_text(encoding="utf-8").splitlines():
        expected, relative = line.split(maxsplit=1)
        relative = relative.lstrip("*")
        path = ROOT / relative
        actual = hashlib.sha256(reader(path)).hexdigest()
        if actual != expected:
            errors.append(
                f"modified third-party file: {relative}: {actual} != {expected}"
            )
    return errors


def validate_documents(
    comparison: dict[str, Any],
    discrepancy_document: dict[str, Any],
    manifest: dict[str, Any],
) -> list[str]:
    errors: list[str] = []
    source_commit = (PINNED_ROOT / "SOURCE_COMMIT").read_text(
        encoding="utf-8"
    ).strip()
    commits = {
        "comparison audit": comparison.get("audit", {}).get(
            "oasis_community_upstream_commit"
        ),
        "discrepancy document": discrepancy_document.get(
            "oasis_community_upstream_commit"
        ),
        "corpus manifest": manifest.get("upstream_commit"),
    }
    for label, commit in commits.items():
        if not commit:
            errors.append(f"missing upstream commit: {label}")
        elif commit != source_commit:
            errors.append(
                f"upstream commit mismatch: {label}: {commit} != {source_commit}"
            )

    cases = manifest.get("cases") or []
    case_index: dict[str, dict[str, Any]] = {}
    corpus_type_index: dict[str, list[str]] = {}
    for case in cases:
        case_id = case.get("id")
        if not case_id:
            errors.append("missing provenance: corpus case ID")
            continue
        if case_id in case_index:
            errors.append(f"duplicate corpus case ID: {case_id}")
        case_index[case_id] = case
        upstream = case.get("upstream") or {}
        required_provenance = (
            "repository",
            "commit",
            "original_path",
            "pinned_path",
            "original_license",
            "license_file",
        )
        for key in required_provenance:
            if not upstream.get(key):
                errors.append(f"missing provenance: {case_id}.upstream.{key}")
        if upstream.get("commit") != source_commit:
            errors.append(
                f"upstream commit mismatch: {case_id}: "
                f"{upstream.get('commit')} != {source_commit}"
            )
        adaptation = case.get("adaptation") or {}
        if adaptation.get("status") not in {"adapted", "unchanged"}:
            errors.append(f"missing provenance: {case_id}.adaptation.status")
        if not adaptation.get("explanation") or not adaptation.get("diff_summary"):
            errors.append(f"missing provenance: {case_id}.adaptation explanation/diff")
        specification = case.get("specification") or {}
        for key in (
            "sections",
            "frozen_requirement_ids",
            "frozen_requirement_rationale",
            "non_must_requirement_ids",
        ):
            if key not in specification or specification.get(key) is None:
                errors.append(f"missing provenance: {case_id}.specification.{key}")
        for name in case.get("types_covered") or []:
            corpus_type_index.setdefault(name, []).append(case_id)

    type_records = comparison.get("types") or []
    type_index: dict[str, dict[str, Any]] = {}
    for record in type_records:
        name = record.get("canonical_name")
        if not name:
            errors.append("unclassified type: missing canonical_name")
            continue
        if name in type_index:
            errors.append(f"unclassified type: duplicate {name}")
        type_index[name] = record
        if not record.get("category") or record.get("category") == "unclassified":
            errors.append(f"unclassified type: {name}")

        specification = record.get("specification") or {}
        oasis = record.get("oasis_community_profile") or {}
        puccini = record.get("puccini_profile") or {}
        if specification.get("present") and not specification.get("section"):
            errors.append(f"missing provenance: {name}.specification.section")
        for label, source in (
            ("oasis_community_profile", oasis),
            ("puccini_profile", puccini),
        ):
            if source.get("present"):
                for key in ("file", "path", "effective_definition_sha256"):
                    if not source.get(key):
                        errors.append(f"missing provenance: {name}.{label}.{key}")
        if specification.get("present") and not specification.get(
            "effective_definition_sha256"
        ):
            errors.append(
                f"missing provenance: {name}.specification.effective_definition_sha256"
            )

        for comparison_kind in ("comparison", "effective_comparison"):
            statuses = record.get(comparison_kind) or {}
            for pair in (
                "specification_vs_oasis",
                "specification_vs_puccini",
                "oasis_vs_puccini",
            ):
                if statuses.get(pair) not in ALLOWED_COMPARISONS:
                    errors.append(
                        f"unclassified type: {name}.{comparison_kind}.{pair}"
                    )

        for list_name in ("field_differences", "effective_field_differences"):
            for difference in record.get(list_name) or []:
                path = difference.get("path") or "<missing>"
                classification = difference.get("classification")
                if classification not in ALLOWED_CLASSIFICATIONS:
                    errors.append(
                        f"unclassified field difference: {name}.{list_name}.{path}"
                    )
                if difference.get("resolution_status") not in ALLOWED_RESOLUTIONS:
                    errors.append(
                        f"discrepancy without resolution status: "
                        f"{name}.{list_name}.{path}"
                    )

        evidence = record.get("corpus_cases") or []
        if not evidence:
            errors.append(f"type without corpus evidence: {name}")
        for case_id in evidence:
            case = case_index.get(case_id)
            if case is None or name not in (case.get("types_covered") or []):
                errors.append(
                    f"type without corpus evidence: {name}: invalid case {case_id}"
                )
        if set(evidence) != set(corpus_type_index.get(name, [])):
            errors.append(
                f"type without corpus evidence: {name}: comparison/manifest mismatch"
            )

        if (
            puccini.get("present")
            and not specification.get("present")
            and not oasis.get("present")
        ):
            classifications = {
                difference.get("classification")
                for difference in record.get("field_differences") or []
            }
            if "extension" not in classifications:
                errors.append(
                    "Puccini normative declaration absent from specification and "
                    f"community profile without extension classification: {name}"
                )

    for name in corpus_type_index:
        if name not in type_index:
            errors.append(f"unclassified type: corpus-only {name}")

    discrepancy_ids: set[str] = set()
    for discrepancy in discrepancy_document.get("discrepancies") or []:
        discrepancy_id = discrepancy.get("id")
        if not discrepancy_id or discrepancy_id in discrepancy_ids:
            errors.append(f"discrepancy without resolution status: ID {discrepancy_id}")
        discrepancy_ids.add(discrepancy_id)
        if discrepancy.get("canonical_name") not in type_index:
            errors.append(
                f"unclassified type: discrepancy {discrepancy_id} references "
                f"{discrepancy.get('canonical_name')}"
            )
        if discrepancy.get("classification") not in ALLOWED_CLASSIFICATIONS:
            errors.append(f"unclassified field difference: {discrepancy_id}")
        if discrepancy.get("resolution_status") not in ALLOWED_RESOLUTIONS:
            errors.append(
                f"discrepancy without resolution status: {discrepancy_id}"
            )
        if not discrepancy.get("specification_section"):
            errors.append(f"missing provenance: {discrepancy_id}.specification_section")

    integrity_inventory = set(
        comparison.get("audit", {}).get("integrity_verified_files") or []
    )
    expected_inventory = {
        line.split(maxsplit=1)[1].lstrip("*")
        for line in (PINNED_ROOT / "SHA256SUMS")
        .read_text(encoding="utf-8")
        .splitlines()
    }
    if integrity_inventory != expected_inventory:
        errors.append("missing provenance: pinned integrity inventory mismatch")

    return errors


def validate_current() -> list[str]:
    comparison = load_yaml(COMPARISON_PATH)
    discrepancies = load_yaml(DISCREPANCIES_PATH)
    manifest = load_yaml(MANIFEST_PATH)
    return validate_integrity() + validate_documents(
        comparison, discrepancies, manifest
    )


def self_test() -> list[str]:
    comparison = load_yaml(COMPARISON_PATH)
    discrepancies = load_yaml(DISCREPANCIES_PATH)
    manifest = load_yaml(MANIFEST_PATH)
    failures: list[str] = []

    def add_unclassified_extension(
        changed_comparison: dict[str, Any],
        _changed_discrepancies: dict[str, Any],
        changed_manifest: dict[str, Any],
    ) -> None:
        changed_comparison["types"].append(
            {
                "canonical_name": "tosca.nodes.UnclassifiedExtension",
                "category": "node_type",
                "specification": {"present": False},
                "oasis_community_profile": {"present": False},
                "puccini_profile": {
                    "present": True,
                    "file": "probe.yaml",
                    "path": "node_types.probe",
                    "effective_definition_sha256": "probe",
                },
                "comparison": {
                    "specification_vs_oasis": "missing",
                    "specification_vs_puccini": "missing",
                    "oasis_vs_puccini": "missing",
                },
                "effective_comparison": {
                    "specification_vs_oasis": "missing",
                    "specification_vs_puccini": "missing",
                    "oasis_vs_puccini": "missing",
                },
                "field_differences": [],
                "effective_field_differences": [],
                "corpus_cases": [changed_manifest["cases"][0]["id"]],
            }
        )

    mutation_type = Callable[
        [dict[str, Any], dict[str, Any], dict[str, Any]], None
    ]
    mutations: list[tuple[str, mutation_type]] = [
        (
            "unclassified type",
            lambda c, _d, _m: c["types"][0].update(category="unclassified"),
        ),
        (
            "unclassified field difference",
            lambda c, _d, _m: c["types"][0].update(
                field_differences=[
                    {
                        "path": "probe",
                        "classification": "unresolved",
                        "resolution_status": "unresolved",
                    }
                ]
            ),
        ),
        (
            "missing provenance",
            lambda _c, _d, m: m["cases"][0]["upstream"].update(original_path=""),
        ),
        (
            "missing upstream commit",
            lambda c, _d, _m: c["audit"].update(
                oasis_community_upstream_commit=""
            ),
        ),
        (
            "type without corpus evidence",
            lambda c, _d, _m: c["types"][0].update(corpus_cases=[]),
        ),
        (
            "discrepancy without resolution status",
            lambda _c, d, _m: d["discrepancies"][0].update(
                resolution_status="unresolved"
            ),
        ),
        (
            "Puccini normative declaration absent",
            add_unclassified_extension,
        ),
    ]
    for expected, mutate in mutations:
        changed_comparison = copy.deepcopy(comparison)
        changed_discrepancies = copy.deepcopy(discrepancies)
        changed_manifest = copy.deepcopy(manifest)
        mutate(changed_comparison, changed_discrepancies, changed_manifest)
        errors = validate_documents(
            changed_comparison, changed_discrepancies, changed_manifest
        )
        if not any(expected in error for error in errors):
            failures.append(
                f"self-test mutation did not trigger {expected!r}: {errors}"
            )

    first_line = (PINNED_ROOT / "SHA256SUMS").read_text(
        encoding="utf-8"
    ).splitlines()[0]
    _, relative = first_line.split(maxsplit=1)
    target = ROOT / relative.lstrip("*")

    def altered_reader(path: pathlib.Path) -> bytes:
        content = path.read_bytes()
        return content + b"\nmodified" if path == target else content

    integrity_errors = validate_integrity(altered_reader)
    if not any("modified third-party file" in error for error in integrity_errors):
        failures.append(
            "self-test mutation did not trigger 'modified third-party file'"
        )
    return failures


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--self-test",
        action="store_true",
        help="also prove every required failure mode with in-memory mutations",
    )
    args = parser.parse_args()
    errors = validate_current()
    if args.self_test:
        errors.extend(self_test())
    if errors:
        for error in errors:
            print(f"ERROR: {error}", file=sys.stderr)
        return 1
    comparison = load_yaml(COMPARISON_PATH)
    discrepancies = load_yaml(DISCREPANCIES_PATH)
    print(
        "TOSCA 1.3 profile completeness gate: "
        f"{len(comparison['types'])} types, "
        f"{len(discrepancies['discrepancies'])} discrepancies, "
        "all classified and evidenced"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
