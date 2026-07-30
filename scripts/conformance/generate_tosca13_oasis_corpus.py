#!/usr/bin/env python3
"""Generate provenance-preserving OASIS community profile corpus cases."""

from __future__ import annotations

import argparse
import pathlib
import re
import sys
from typing import Any

import yaml


ROOT = pathlib.Path(__file__).resolve().parents[2]
UPSTREAM_ROOT = ROOT / "third_party/oasis/tosca-simple-profile/1.3"
LICENSE_PATH = ROOT / "third_party/oasis/tosca-simple-profile/LICENSE"
COMPARISON_PATH = ROOT / "docs/conformance/tosca-1.3/profile-three-way-comparison.yaml"
NON_MUST_PATH = ROOT / "docs/conformance/tosca-1.3/non-must-requirements.yaml"
OUTPUT_ROOT = ROOT / "tests/corpus/tosca_1_3/oasis_community"
MANIFEST_PATH = OUTPUT_ROOT / "manifest.yaml"

SOURCE_KINDS = {
    "data.yaml": ("data_types", "data-types-valid.yaml"),
    "artifact.yaml": ("artifact_types", "artifact-types-valid.yaml"),
    "capability.yaml": ("capability_types", "capability-types-valid.yaml"),
    "relationship.yaml": ("relationship_types", "relationship-types-valid.yaml"),
    "interface.yaml": ("interface_types", "interface-types-valid.yaml"),
    "node.yaml": ("node_types", "node-types-valid.yaml"),
    "group.yaml": ("group_types", "group-types-valid.yaml"),
    "policy.yaml": ("policy_types", "policy-types-valid.yaml"),
}


def load_yaml(path: pathlib.Path) -> Any:
    return yaml.safe_load(path.read_text(encoding="utf-8"))


def dump_yaml(value: Any) -> str:
    return yaml.safe_dump(
        value,
        sort_keys=False,
        allow_unicode=True,
        width=100,
    )


def declaration_requirement_index() -> dict[str, str]:
    index: dict[str, str] = {}
    catalog = load_yaml(NON_MUST_PATH)
    pattern = re.compile(
        r"^The normative type `([^`]+)` is declared by this specification\.$"
    )
    for requirement in catalog["requirements"]:
        match = pattern.match(requirement.get("normative_rule", ""))
        if match:
            index[match.group(1)] = requirement["id"]
    return index


def generate_case(
    source_name: str,
    section_key: str,
    fixture_name: str,
    comparison_index: dict[str, dict[str, Any]],
    requirement_index: dict[str, str],
    upstream_commit: str,
) -> tuple[dict[str, Any], str]:
    source_path = UPSTREAM_ROOT / source_name
    source = load_yaml(source_path)
    declarations = source.get(section_key) or {}
    if not declarations:
        raise ValueError(f"{source_path} has no {section_key}")

    derived: dict[str, dict[str, str]] = {}
    type_names = list(declarations)
    for index, type_name in enumerate(type_names, start=1):
        derived[f"corpus.oasis_community.{section_key}.{index:02d}"] = {
            "derived_from": type_name
        }
    fixture = {
        "tosca_definitions_version": "tosca_simple_yaml_1_3",
        section_key: derived,
        "topology_template": {},
    }

    sections = [comparison_index[name]["specification"]["section"] for name in type_names]
    requirement_ids = [requirement_index[name] for name in type_names]
    case_id = "TOSCA13-OASIS-COMMUNITY-" + section_key.removesuffix("_types").upper()
    case = {
        "id": case_id,
        "file": fixture_name,
        "upstream": {
            "repository": (
                "https://github.com/oasis-open/tosca-community-contributions"
            ),
            "commit": upstream_commit,
            "original_path": f"profiles/org/oasis-open/simple/1.3/{source_name}",
            "pinned_path": str(source_path.relative_to(ROOT)),
            "original_license": "Apache-2.0",
            "license_file": str(LICENSE_PATH.relative_to(ROOT)),
        },
        "adaptation": {
            "status": "adapted",
            "explanation": (
                "Each upstream normative type name is used as derived_from by a "
                "minimal local declaration. Descriptions and field bodies are not "
                "copied; the case isolates profile loading, exact name resolution, "
                "category resolution, hierarchy construction, and inheritance."
            ),
            "diff_summary": (
                f"Replaced {len(type_names)} upstream declarations with "
                f"{len(type_names)} minimal derived declarations and added an empty "
                "topology_template."
            ),
        },
        "specification": {
            "version": "TOSCA Simple Profile in YAML Version 1.3",
            "sections": sections,
            "frozen_requirement_ids": [],
            "frozen_requirement_rationale": (
                "Normative type-presence declarations are tracked in the independent "
                "non-MUST/full-section catalog and do not alter the frozen 223 MUST "
                "denominator."
            ),
            "non_must_requirement_ids": requirement_ids,
        },
        "types_covered": type_names,
        "expected": {
            "processor_result": "accept",
            "phase": "hierarchy",
            "assertions": [
                "all exact case-sensitive parent type names resolve",
                "all derived declarations inherit from the requested normative type",
            ],
        },
    }
    return case, dump_yaml(fixture)


def generated_outputs() -> dict[pathlib.Path, str]:
    comparison = load_yaml(COMPARISON_PATH)
    comparison_index = {
        item["canonical_name"]: item for item in comparison["types"]
    }
    requirement_index = declaration_requirement_index()
    upstream_commit = (UPSTREAM_ROOT / "SOURCE_COMMIT").read_text(
        encoding="utf-8"
    ).strip()

    cases: list[dict[str, Any]] = []
    outputs: dict[pathlib.Path, str] = {}
    for source_name, (section_key, fixture_name) in SOURCE_KINDS.items():
        case, fixture = generate_case(
            source_name,
            section_key,
            fixture_name,
            comparison_index,
            requirement_index,
            upstream_commit,
        )
        cases.append(case)
        outputs[OUTPUT_ROOT / fixture_name] = fixture

    manifest = {
        "schema_version": 1,
        "layer": "OASIS community independent comparison corpus",
        "upstream_repository": (
            "https://github.com/oasis-open/tosca-community-contributions"
        ),
        "upstream_commit": upstream_commit,
        "pinned_profile_root": str(UPSTREAM_ROOT.relative_to(ROOT)),
        "cases": cases,
    }
    outputs[MANIFEST_PATH] = dump_yaml(manifest)
    return outputs


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    outputs = generated_outputs()
    if args.check:
        stale = [
            str(path.relative_to(ROOT))
            for path, content in outputs.items()
            if not path.exists() or path.read_text(encoding="utf-8") != content
        ]
        if stale:
            print("stale OASIS community corpus files:", file=sys.stderr)
            for path in stale:
                print(f"  {path}", file=sys.stderr)
            return 1
        print(f"OASIS community corpus: {len(outputs) - 1} adapted cases")
        return 0
    for path, content in outputs.items():
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")
    print(f"generated {len(outputs) - 1} OASIS community corpus cases")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
