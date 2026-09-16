#!/usr/bin/env python3
"""Generate deterministic TOSCA 1.3 template-corpus audit artifacts."""

from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "docs/conformance/tosca-1.3"
CORPUS = ROOT / "tests/corpus/tosca_1_3"


def load(path):
    return yaml.load(path.read_text(encoding="utf-8"), Loader=yaml.CSafeLoader)


def dump(path, value):
    path.write_text(
        yaml.safe_dump(value, sort_keys=False, allow_unicode=True, width=120),
        encoding="utf-8",
    )


catalog = {r["id"]: r for r in load(BASE / "requirements.yaml")["requirements"]}
coverage = load(BASE / "coverage.yaml")["requirements"]
frozen = [r for r in coverage if r.get("included_in_must_denominator")]
manifest = load(CORPUS / "manifest.yaml")["cases"]

by_section = {}
for case in manifest:
    for section in case["specification"]["sections"]:
        entry = by_section.setdefault(section, {"requirements": set(), "axes": set(), "fixtures": []})
        entry["requirements"].update(case["specification"]["requirement_ids"])
        entry["axes"].update(case["coverage"]["variation_axes"])
        entry["fixtures"].append(case["id"])

frozen_sections = {}
for record in frozen:
    frozen_sections.setdefault(record["section"], []).append(record["requirement_id"])

sections = []
all_sections = sorted(
    {r["section"] for r in catalog.values()},
    key=lambda value: tuple(int(part) if part.isdigit() else part for part in value.split(".")),
)
for section in all_sections:
    source = next(r for r in catalog.values() if r["section"] == section)
    evidence = by_section.get(section)
    applicable = section in frozen_sections
    sections.append(
        {
            "section": section,
            "title": source["section_title"],
            "historical_frozen_membership": applicable,
            "applicability_review": "deferred; membership is not applicability proof",
            "grammar_entities": sorted({catalog[i]["subject"] for i in frozen_sections.get(section, [])}),
            "requirement_ids": frozen_sections.get(section, []),
            "variation_axes": sorted(evidence["axes"]) if evidence else [],
            "fixtures": sorted(set(evidence["fixtures"])) if evidence else [],
            "status": "insufficient-evidence",
            "exclusions": [],
        }
    )

dump(
    BASE / "template-corpus-sections.yaml",
    {
        "schema_version": 1,
        "normative_source": "docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html",
        "scope": "historical section/fixture label inventory; neither applicability nor complete verification is established",
        "sections": sections,
    },
)

case_ids = {case["id"] for case in manifest}
pairwise = {
    "schema_version": 1,
    "method": "historical candidate axis inventory; combinations and exclusions need independent review; fixture names are not verification",
    "pairs": [
        {
            "axes": ["notation", "inheritance"],
            "values": {"notation": ["short", "long"], "inheritance": ["direct", "inherited"]},
            "fixtures": sorted(i for i in case_ids if "IMPORT" in i or "TYPE-HIERARCHY" in i),
            "uncovered_combinations": [],
        },
        {
            "axes": ["notation", "imports"],
            "values": {"notation": ["short", "long"], "imports": ["direct", "qualified", "transitive"]},
            "fixtures": sorted(i for i in case_ids if "IMPORT" in i),
            "uncovered_combinations": [],
        },
        {
            "axes": ["function", "assignment-context"],
            "values": {"function": ["get_input", "get_property", "get_attribute"], "assignment-context": ["output", "template"]},
            "fixtures": sorted(i for i in case_ids if "FUNCTION" in i or "NORMALIZATION" in i),
            "uncovered_combinations": [
                {
                    "combination": ["get_attribute", "definition-default"],
                    "reason": "attribute values are runtime values; a property-reflection fixture covers the effective-value path",
                }
            ],
        },
        {
            "axes": ["constraint", "data-type"],
            "values": {"constraint": ["range", "length", "pattern"], "data-type": ["integer", "collection", "string"]},
            "fixtures": sorted(i for i in case_ids if "CONSTRAINT" in i or "SCALAR-UNIT" in i),
            "uncovered_combinations": [
                {"combination": ["pattern", "integer"], "reason": "normatively incompatible and covered as an invalid operand case"}
            ],
        },
        {
            "axes": ["default", "refinement"],
            "values": {"default": ["inherited", "overridden"], "refinement": ["direct", "derived"]},
            "fixtures": sorted(
                i
                for i in case_ids
                if i
                in {
                    "TOSCA13-CORPUS-INTERACTION-INHERITANCE-REFINEMENT",
                    "TOSCA13-CORPUS-INTERACTION-NORMATIVE-DERIVATION",
                    "TOSCA13-CORPUS-NORMALIZATION-FUNCTION-INHERITANCE-VALID",
                    "TOSCA13-CORPUS-NORMATIVE-PROFILE-DERIVED-VALID",
                }
            ),
            "uncovered_combinations": [],
        },
        {
            "axes": ["reference-kind", "namespace"],
            "values": {"reference-kind": ["type", "template", "property"], "namespace": ["local", "qualified"]},
            "fixtures": sorted(i for i in case_ids if "IMPORT" in i or "SUBSTITUTION" in i),
            "uncovered_combinations": [],
        },
        {
            "axes": ["template-kind", "normalization"],
            "values": {"template-kind": ["node", "relationship", "group", "policy"], "normalization": ["direct", "interaction"]},
            "fixtures": sorted(i for i in case_ids if "TEMPLATE" in i or "NORMALIZATION" in i),
            "uncovered_combinations": [],
        },
    ],
}
dump(BASE / "template-corpus-pairwise.yaml", pairwise)

valid = sum(case["classification"]["kind"] == "valid" for case in manifest)
invalid = len(manifest) - valid
summary = {
    "schema_version": 1,
    "total_cases": len(manifest),
    "valid_cases": valid,
    "invalid_cases": invalid,
    "interaction_cases": sum(case["file"].startswith("interactions/") for case in manifest),
    "csar_cases": sum(case["classification"]["category"] == "csar" for case in manifest),
    "normalization_cases": sum(case["classification"]["category"] == "normalization" for case in manifest),
    "historical_requirement_ids_labelled": len(
        {i for case in manifest for i in case["specification"]["requirement_ids"]}
        & {r["requirement_id"] for r in frozen}
    ),
    "historical_denominator": len(frozen),
    "owned_primary_requirements": len({i for c in manifest for i in c["coverage"]["primary_requirements"]}),
    "verification_result": "re-audit/evidence-recheck.yaml",
    "sections_with_historical_membership": len(frozen_sections),
}
dump(BASE / "template-corpus-summary.yaml", summary)
