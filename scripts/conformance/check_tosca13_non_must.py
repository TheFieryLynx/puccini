#!/usr/bin/env python3
"""Fail closed on TOSCA 1.3 non-MUST and full-section coverage drift."""

from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "docs/conformance/tosca-1.3"
CORPUS = ROOT / "tests/corpus/tosca_1_3"

CLASSIFICATIONS = {
    "processor-must-covered",
    "processor-non-must-applicable",
    "service-template-author",
    "orchestrator-only",
    "generator-only",
    "archive-target",
    "informative",
    "example",
    "definition-only",
    "cross-reference-only",
    "duplicate-obligation",
    "not-applicable",
}


def load(path):
    return yaml.safe_load(path.read_text(encoding="utf-8"))


def fail(message):
    raise SystemExit(message)


coverage = load(BASE / "coverage.yaml")["requirements"]
frozen = [record for record in coverage if record.get("included_in_must_denominator")]
if len(frozen) != 223:
    fail(f"frozen denominator is {len(frozen)}, want 223")

catalog = load(BASE / "non-must-requirements.yaml")["requirements"]
manifest = load(CORPUS / "non_must/manifest.yaml")["cases"]
case_ids = {case["id"] for case in manifest}
if len(case_ids) != len(manifest):
    fail("duplicate non-MUST case ID")
for requirement in catalog:
    if not requirement["fixture_ids"] or not set(requirement["fixture_ids"]) <= case_ids:
        fail(f"{requirement['id']} lacks a manifest fixture")
    if requirement["source_strength"] == "SHOULD" and not requirement["implementation_policy"]:
        fail(f"{requirement['id']} lacks a SHOULD decision")
    if requirement["source_strength"] in {"MAY", "OPTIONAL"} and requirement["implementation_policy"] not in {
        "supported",
        "unsupported",
        "implementation-defined",
    }:
        fail(f"{requirement['id']} lacks a MAY/OPTIONAL support policy")

sections = load(BASE / "template-corpus-sections.yaml")["sections"]
if len(sections) != 590 or len({str(section["section"]) for section in sections}) != 590:
    fail("full-section inventory is not exactly 590 unique identifiers")
for section in sections:
    classification = section.get("primary_classification")
    if classification not in CLASSIFICATIONS:
        fail(f"section {section['section']} has unknown classification {classification!r}")
    for layer in ("must_coverage", "non_must_coverage", "example_coverage"):
        if layer not in section or "status" not in section[layer]:
            fail(f"section {section['section']} lacks {layer}")
    if classification == "processor-non-must-applicable" and not section["non_must_coverage"]["requirement_ids"]:
        fail(f"section {section['section']} has no non-MUST requirement")
    if classification not in {"processor-must-covered", "processor-non-must-applicable"}:
        exclusion = section["exclusion"]
        if not exclusion["category"] or not exclusion["reason"]:
            fail(f"section {section['section']} lacks a target-specific exclusion")
    if classification == "example" and section["example_coverage"]["status"] != "covered":
        fail(f"example section {section['section']} lacks compatibility coverage")

defaults = [record for record in catalog if record["source_strength"] == "DEFAULT"]
if not defaults or any(not record["fixture_ids"] for record in defaults):
    fail("normative defaults lack historical fixture labels")
implementation_defined = [
    record for record in catalog if record["implementation_policy"] == "implementation-defined"
]
if not implementation_defined:
    fail("implementation-defined processor choices are undocumented")

print(
    f"TOSCA 1.3 gates: historical-frozen=223; verification=re-audit/evidence-recheck.yaml, non-must={len(catalog)}, "
    f"sections={len(sections)}, cases={len(manifest)}"
)
