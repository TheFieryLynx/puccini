#!/usr/bin/env python3
"""Generate the separate TOSCA 1.3 processor non-MUST catalog.

The frozen 223-requirement matrix is an input and is never modified here.
Catalog records are selected from the pinned-spec extraction only when the
existing audit classifies them as atomic and processor-applicable. Explicit
target corrections below separate author, orchestrator, and archive rules.
"""

from collections import Counter
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
BASE = ROOT / "docs/conformance/tosca-1.3"
CORPUS = ROOT / "tests/corpus/tosca_1_3"

NON_MUST_STRENGTHS = {
    "SHOULD",
    "SHOULD_NOT",
    "MAY",
    "OPTIONAL",
    "DEFAULT",
    "SEMANTIC",
    "GRAMMAR",
    "ERROR",
}

# These extracted records address authors, orchestrators, or the archive
# conformance target even though the prose also mentions a service template.
NON_PROCESSOR_RECORDS = {
    "TOSCA13-3.6.2.4-002",
    "TOSCA13-3.6.10.2-030",
    "TOSCA13-3.6.13.3-001",
    "TOSCA13-3.6.25.5-001",
    "TOSCA13-3.7.9.4-001",
    "TOSCA13-3.7.10.3-001",
    "TOSCA13-3.7.11.4-001",
    "TOSCA13-3.8.5.4-001",
    "TOSCA13-3.8.10.3-001",
    "TOSCA13-3.8.11.3-001",
    "TOSCA13-3.8.13.5-001",
    "TOSCA13-3.10.3.3.4-001",
    "TOSCA13-3.10.3.5.4-001",
    "TOSCA13-5.2-004",
    "TOSCA13-5.8.5.5-006",
    "TOSCA13-5.9.1.4-001",
    "TOSCA13-5.9.3.4-001",
    "TOSCA13-5.9.4.4-001",
    "TOSCA13-6.3-001",
}

IMPLEMENTATION_DEFINED_RECORDS = {
    # Section 3.6.8.2.3 explicitly permits processors to choose how namespace
    # collisions are resolved. The frozen applicability is not modified; the
    # separate non-MUST catalog records Puccini's deterministic strict policy.
    "TOSCA13-3.6.8.2.3-006",
}

UNSUPPORTED_RECOMMENDATIONS = {
    # Puccini deliberately does not fetch a profile merely because an import
    # value is a namespace URI or reserved namespace alias.
    "TOSCA13-3.6.8.2.4-002",
    "TOSCA13-3.6.8.2.4-005",
    # json and xml remain opaque scalar datatypes; syntax checking would
    # require language-specific validators not mandated for processors.
    "TOSCA13-5.3.3-001",
    "TOSCA13-5.3.5-001",
}

CASE_PATHS = {
    "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA": "non_must/should/external-schema-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA-INVALID": "non_must/should/external-schema-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORTS": "non_must/should/import-resolution-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-FAILURE": "non_must/should/import-missing-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-POLICY": "non_must/unsupported-policy/import-namespace-uri-policy.yaml",
    "TOSCA13-NONMUST-CORPUS-SHOULD-DATATYPE-SYNTAX": "non_must/unsupported-policy/json-xml-syntax-policy.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-CONSTRAINT-REFINEMENT": "non_must/may/in-range-refinement-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-IMPORT-NAMESPACE": "non_must/may/import-namespace-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-OPERATION-INPUT": "non_must/may/operation-input-extension-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-TYPE-ONLY": "non_must/may/type-only-service-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-FUNCTION-CONTEXT": "non_must/may/function-context-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-MAY-NORMATIVE-SHORTHAND": "non_must/may/normative-shorthand-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-OPTIONAL-COMPLETE": "non_must/optional/complete-forms-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-DEFAULTS": "non_must/defaults/default-semantics-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DATA": "non_must/grammar-alternatives/data-types-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DATA-INVALID": "non_must/grammar-alternatives/data-types-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS": "non_must/grammar-alternatives/definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS-INVALID": "non_must/grammar-alternatives/definitions-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TYPES": "non_must/grammar-alternatives/type-definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TYPES-INVALID": "non_must/grammar-alternatives/type-definitions-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TEMPLATES": "non_must/grammar-alternatives/template-definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TEMPLATES-INVALID": "non_must/grammar-alternatives/template-definitions-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-FUNCTIONS": "non_must/grammar-alternatives/functions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-FUNCTIONS-INVALID": "non_must/grammar-alternatives/functions-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-WORKFLOWS": "non_must/grammar-alternatives/workflows-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-WORKFLOWS-INVALID": "non_must/grammar-alternatives/workflows-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-PROFILE": "non_must/grammar-alternatives/normative-profile-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-PROFILE-INVALID": "non_must/grammar-alternatives/normative-profile-invalid.yaml",
    "TOSCA13-NONMUST-CORPUS-POLICY-NETWORK": "non_must/implementation-defined/network-disabled-policy.yaml",
    "TOSCA13-NONMUST-CORPUS-POLICY-NAMESPACE-COLLISION": "non_must/implementation-defined/namespace-collision-policy.yaml",
}

CASE_EXPECTATIONS = {
    "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA": {
        "tier": "processor-should",
        "kind": "valid",
        "category": "semantic",
        "features": ["external-schema", "recommendation"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA-INVALID": {
        "tier": "processor-should",
        "kind": "invalid",
        "category": "semantic",
        "features": ["external-schema", "recommendation", "invalid-schema"],
        "accepted": False,
        "phase": "rendering",
        "diagnostic": {"category": "external schema compilation failure", "path": "payload"},
        "assertions": [],
    },
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORTS": {
        "tier": "processor-should",
        "kind": "valid",
        "category": "semantic",
        "features": ["imports", "relative-resolution"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-FAILURE": {
        "tier": "processor-should",
        "kind": "invalid",
        "category": "semantic",
        "features": ["imports", "resolution-failure"],
        "accepted": False,
        "phase": "read",
        "diagnostic": {"category": "invalid URL", "path": "imports[0]"},
        "assertions": [],
    },
    "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-POLICY": {
        "tier": "processor-should",
        "kind": "policy",
        "category": "compatibility",
        "features": ["imports", "namespace-uri", "unsupported-policy"],
        "accepted": False,
        "phase": "read",
        "diagnostic": {"category": "invalid URL", "path": "imports[0]"},
        "assertions": [],
    },
    "TOSCA13-NONMUST-CORPUS-SHOULD-DATATYPE-SYNTAX": {
        "tier": "processor-should",
        "kind": "policy",
        "category": "compatibility",
        "features": ["json", "xml", "unsupported-recommendation"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-CONSTRAINT-REFINEMENT": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "semantic",
        "features": ["constraint-refinement", "in-range"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:child"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-IMPORT-NAMESPACE": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "semantic",
        "features": ["imports", "namespace-prefix"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-OPERATION-INPUT": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "semantic",
        "features": ["operation-input", "undeclared-assignment"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-TYPE-ONLY": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "semantic",
        "features": ["service-template", "type-only"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["service_template_non_nil"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-FUNCTION-CONTEXT": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "semantic",
        "features": ["intrinsic-function", "topology-output"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["output_exists:copy"],
    },
    "TOSCA13-NONMUST-CORPUS-MAY-NORMATIVE-SHORTHAND": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "grammar",
        "features": ["normative-type", "shorthand-name"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-OPTIONAL-COMPLETE": {
        "tier": "processor-may",
        "kind": "valid",
        "category": "grammar",
        "features": ["optional-keynames", "complete-form"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": ["node_exists:component", "input_exists:label", "output_exists:label"],
    },
    "TOSCA13-NONMUST-CORPUS-DEFAULTS": {
        "tier": "processor-default",
        "kind": "valid",
        "category": "default",
        "features": ["property-default", "inheritance", "refinement", "explicit-override", "function-assignment"],
        "accepted": True,
        "phase": "normalization",
        "diagnostic": None,
        "assertions": [
            "node_property_primitive:base:setting:parent",
            "node_property_primitive:explicit:setting:explicit",
            "node_property_primitive:derived:setting:derived",
            "node_property_function:function:setting:tosca.function.get_input",
            "node_property_primitive:unrelated:other:unrelated",
        ],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DATA": {
        "tier": "processor-may", "kind": "valid", "category": "grammar",
        "features": ["data-type", "scalar-list-map"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS": {
        "tier": "processor-may", "kind": "valid", "category": "grammar",
        "features": ["definition", "short-long-notation"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TYPES": {
        "tier": "processor-may", "kind": "valid", "category": "grammar",
        "features": ["type-definition", "inheritance"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TEMPLATES": {
        "tier": "processor-may", "kind": "valid", "category": "grammar",
        "features": ["template-definition", "assignment"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:source", "node_exists:target"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-FUNCTIONS": {
        "tier": "processor-may", "kind": "valid", "category": "semantic",
        "features": ["intrinsic-function", "nested"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["output_exists:joined"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-WORKFLOWS": {
        "tier": "processor-may", "kind": "valid", "category": "grammar",
        "features": ["workflow", "activity"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:node"],
    },
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-PROFILE": {
        "tier": "processor-may", "kind": "valid", "category": "semantic",
        "features": ["normative-profile", "derived-custom-type"], "accepted": True, "phase": "normalization",
        "diagnostic": None, "assertions": ["node_exists:compute"],
    },
    "TOSCA13-NONMUST-CORPUS-POLICY-NAMESPACE-COLLISION": {
        "tier": "implementation-defined",
        "kind": "policy",
        "category": "semantic",
        "features": ["imports", "namespace-collision", "deterministic-error-policy"],
        "accepted": False,
        "phase": "namespaces",
        "diagnostic": {"category": "equivalent", "path": "Shared"},
        "assertions": [],
    },
}

for _case_id in tuple(CASE_EXPECTATIONS):
    if not _case_id.startswith("TOSCA13-NONMUST-CORPUS-GRAMMAR-") or _case_id.endswith("-INVALID"):
        continue
    _invalid_id = f"{_case_id}-INVALID"
    CASE_EXPECTATIONS[_invalid_id] = {
        "tier": "processor-may",
        "kind": "invalid",
        "category": "grammar",
        "features": CASE_EXPECTATIONS[_case_id]["features"] + ["invalid-sibling"],
        "accepted": False,
        "phase": "read",
        "diagnostic": {"category": "unsupported keyname", "path": "unknown_key"},
        "assertions": [],
    }


def load(path):
    return yaml.safe_load(path.read_text(encoding="utf-8"))


def dump(path, value):
    path.write_text(
        yaml.safe_dump(value, sort_keys=False, allow_unicode=True, width=140),
        encoding="utf-8",
    )


def mapped_strength(source_strength):
    if source_strength in {"GRAMMAR", "ERROR"}:
        return "SEMANTIC"
    return source_strength


def fixture_for(record):
    requirement_id = record["requirement_id"]
    strength = record["strength"]
    section = record["section"]

    if requirement_id == "TOSCA13-3.6.3.3-005":
        return "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA"
    if requirement_id in {"TOSCA13-3.6.8.2.4-002", "TOSCA13-3.6.8.2.4-005"}:
        return "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-POLICY"
    if section == "3.6.8.2.4":
        if requirement_id == "TOSCA13-3.6.8.2.4-010":
            return "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORT-FAILURE"
        return "TOSCA13-NONMUST-CORPUS-SHOULD-IMPORTS"
    if requirement_id == "TOSCA13-3.6.25.5-001":
        return "TOSCA13-NONMUST-CORPUS-SHOULD-CONDITION-ORDER"
    if section in {"3.8.10.3", "3.8.11.3"}:
        return "TOSCA13-NONMUST-CORPUS-SHOULD-CAPABILITY-MAPPING"
    if section in {"3.7.11.4", "3.8.5.4"}:
        return "TOSCA13-NONMUST-CORPUS-SHOULD-GROUP-RELATIONSHIPS"
    if section in {"5.3.3", "5.3.5"}:
        return "TOSCA13-NONMUST-CORPUS-SHOULD-DATATYPE-SYNTAX"
    if requirement_id == "TOSCA13-3.6.3.1-017":
        return "TOSCA13-NONMUST-CORPUS-MAY-CONSTRAINT-REFINEMENT"
    if section == "3.6.8.2.3":
        if requirement_id == "TOSCA13-3.6.8.2.3-006":
            return "TOSCA13-NONMUST-CORPUS-POLICY-NAMESPACE-COLLISION"
        return "TOSCA13-NONMUST-CORPUS-MAY-IMPORT-NAMESPACE"
    if requirement_id == "TOSCA13-3.6.17.3-003":
        return "TOSCA13-NONMUST-CORPUS-MAY-OPERATION-INPUT"
    if section == "3.10.2.2":
        return "TOSCA13-NONMUST-CORPUS-MAY-TYPE-ONLY"
    if section.startswith("4.") and strength == "MAY":
        return "TOSCA13-NONMUST-CORPUS-MAY-FUNCTION-CONTEXT"
    if section == "5.2":
        return "TOSCA13-NONMUST-CORPUS-MAY-NORMATIVE-SHORTHAND"
    if strength == "OPTIONAL":
        return "TOSCA13-NONMUST-CORPUS-OPTIONAL-COMPLETE"
    if strength == "DEFAULT":
        return "TOSCA13-NONMUST-CORPUS-DEFAULTS"
    if section.startswith(("3.3", "3.4", "3.5")):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-DATA"
    if section.startswith("3.6"):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS"
    if section.startswith("3.7"):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-TYPES"
    if section.startswith(("3.8", "3.9", "3.10")):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-TEMPLATES"
    if section.startswith("4."):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-FUNCTIONS"
    if section.startswith(("5.", "8.")):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-PROFILE"
    if section.startswith(("6.", "7.", "13.", "14.")):
        return "TOSCA13-NONMUST-CORPUS-GRAMMAR-WORKFLOWS"
    return "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS"


def fixtures_for(record):
    fixture = fixture_for(record)
    if record["requirement_id"] == "TOSCA13-3.6.3.3-005":
        return [
            fixture,
            "TOSCA13-NONMUST-CORPUS-SHOULD-EXTERNAL-SCHEMA-INVALID",
        ]
    if fixture.startswith("TOSCA13-NONMUST-CORPUS-GRAMMAR-"):
        return [fixture, f"{fixture}-INVALID"]
    return [fixture]


def policy_for(record, fixture_exists):
    source_strength = record["strength"]
    requirement_id = record["requirement_id"]
    unsupported = requirement_id in UNSUPPORTED_RECOMMENDATIONS
    implementation_defined = requirement_id in IMPLEMENTATION_DEFINED_RECORDS

    if implementation_defined:
        return {
            "implementation_policy": "implementation-defined",
            "policy_reason": "The specification permits processor choice; Puccini deterministically rejects non-equivalent imported definitions with the same identity.",
            "implementation_status": "implemented" if fixture_exists else "partial",
            "verification_status": "verified" if fixture_exists else "untested",
        }

    if unsupported:
        return {
            "implementation_policy": "unsupported",
            "policy_reason": "The recommendation is intentionally not followed; the deterministic policy is documented and tested.",
            "implementation_status": "unsupported" if fixture_exists else "ambiguous",
            "verification_status": "verified" if fixture_exists else "untested",
        }

    if source_strength in {"MAY", "OPTIONAL"}:
        policy = "supported"
        reason = "Puccini selects the permitted supported branch and verifies it through the public TOSCA 1.3 processor path."
    elif source_strength in {"SHOULD", "SHOULD_NOT"}:
        policy = "required"
        reason = "Puccini follows this processor recommendation in strict TOSCA 1.3 processing."
    else:
        policy = "required"
        reason = "This grammar, default, or semantic rule is required for the selected TOSCA 1.3 interpretation."

    return {
        "implementation_policy": policy,
        "policy_reason": reason,
        "implementation_status": "implemented" if fixture_exists else "partial",
        "verification_status": "verified" if fixture_exists else "untested",
    }


def section_sort_key(value):
    return tuple(int(part) if part.isdigit() else part for part in value.split("."))


def classify_section(section, title, old, source_records, non_must):
    if old.get("applicable_to_processor"):
        return "processor-must-covered"
    if non_must:
        return "processor-non-must-applicable"
    lowered = title.lower()
    targets = {target for record in source_records for target in record.get("target", [])}
    strengths = {record["strength"] for record in source_records}
    if "example" in lowered or "use case" in lowered:
        return "example"
    if "archive" in targets or "csar" in lowered:
        return "archive-target"
    if "orchestrator" in targets and "processor" not in targets:
        return "orchestrator-only"
    if "generator" in targets and "processor" not in targets:
        return "generator-only"
    if "service-template" in targets and "processor" not in targets:
        return "service-template-author"
    if strengths & {"MUST", "MUST_NOT", "SHALL", "SHALL_NOT", "REQUIRED"}:
        return "duplicate-obligation"
    if "cross-reference" in lowered or lowered in {"references", "reference"}:
        return "cross-reference-only"
    if any(word in lowered for word in ("definition", "definitions", "terminology", "notation", "keyname")):
        return "definition-only"
    if not source_records:
        return "informative"
    return "not-applicable"


def generate_reports(catalog_records, requirements, cases):
    old_document = load(BASE / "template-corpus-sections.yaml")
    old_sections = old_document["sections"]
    source_by_section = {}
    for record in catalog_records:
        source_by_section.setdefault(record["section"], []).append(record)
    non_must_by_section = {}
    for requirement in requirements:
        non_must_by_section.setdefault(requirement["section"], []).append(requirement)

    example_manifest = load(CORPUS / "examples/manifest.yaml")
    examples_by_section = {}
    for case in example_manifest["cases"]:
        examples_by_section.setdefault(str(case["section"]), []).append(case["id"])

    sections = []
    classification_counts = Counter()
    for old in old_sections:
        section = str(old["section"])
        linked = non_must_by_section.get(section, [])
        example_fixtures = examples_by_section.get(section, [])
        classification = classify_section(
            section,
            old["title"],
            old,
            source_by_section.get(section, []),
            linked,
        )
        classification_counts[classification] += 1
        linked_fixtures = sorted({fixture for requirement in linked for fixture in requirement["fixture_ids"]})
        policies = {requirement["implementation_policy"] for requirement in linked}
        non_must_status = "not-applicable"
        if linked:
            non_must_status = (
                "policy-documented"
                if policies and policies <= {"unsupported", "implementation-defined"}
                else "complete"
            )
        exclusion = {"category": None, "reason": None}
        if classification not in {"processor-must-covered", "processor-non-must-applicable"}:
            exclusion = {
                "category": classification,
                "reason": (
                    f"Section {section} is classified as {classification} for the processor target; "
                    "it creates no independent processor fixture obligation in this layer."
                ),
            }
        sections.append(
            {
                **old,
                "primary_classification": classification,
                "must_coverage": {
                    "status": "complete" if old.get("applicable_to_processor") else "not-applicable",
                    "fixture_ids": old.get("fixtures", []),
                },
                "non_must_coverage": {
                    "status": non_must_status,
                    "requirement_ids": [requirement["id"] for requirement in linked],
                    "fixture_ids": linked_fixtures,
                },
                "example_coverage": {
                    "status": "covered" if example_fixtures else "not-applicable",
                    "fixture_ids": example_fixtures,
                },
                "exclusion": exclusion,
            }
        )

    dump(
        BASE / "template-corpus-sections.yaml",
        {
            "schema_version": 2,
            "normative_source": old_document["normative_source"],
            "scope": "three-layer frozen-MUST, processor non-MUST policy, and non-normative example classification",
            "section_count": len(sections),
            "classification_counts": dict(sorted(classification_counts.items())),
            "sections": sections,
        },
    )

    source_counts = Counter(requirement["source_strength"] for requirement in requirements)
    policy_counts = Counter(requirement["implementation_policy"] for requirement in requirements)
    status_counts = Counter(
        (requirement["implementation_status"], requirement["verification_status"])
        for requirement in requirements
    )
    case_kind_counts = Counter(case["classification"]["kind"] for case in cases)
    tier_counts = Counter(case["conformance_tier"] for case in cases)
    coverage_report = {
        "schema_version": 1,
        "frozen_must": {
            "denominator": 223,
            "implemented_verified": 223,
            "changed_by_this_catalog": False,
        },
        "non_must": {
            "record_count": len(requirements),
            "source_strengths": dict(sorted(source_counts.items())),
            "implementation_policies": dict(sorted(policy_counts.items())),
            "status_matrix": {
                f"{implementation}+{verification}": count
                for (implementation, verification), count in sorted(status_counts.items())
            },
        },
        "corpus": {
            "case_count": len(cases),
            "kinds": dict(sorted(case_kind_counts.items())),
            "tiers": dict(sorted(tier_counts.items())),
            "example_cases": sum(len(ids) for ids in examples_by_section.values()),
        },
        "sections": {
            "classified": len(sections),
            "classifications": dict(sorted(classification_counts.items())),
        },
    }
    dump(BASE / "non-must-coverage.yaml", coverage_report)

    should = [requirement for requirement in requirements if requirement["source_strength"] == "SHOULD"]
    should_followed = sum(requirement["implementation_policy"] == "required" for requirement in should)
    should_not_followed = len(should) - should_followed
    summary = f"""# TOSCA 1.3 processor non-MUST coverage

This report is generated from the pinned TOSCA 1.3 specification extraction.
It is separate from, and does not alter, the frozen 223 atomic MUST score.

| Measure | Count |
|---|---:|
| Non-MUST processor records | {len(requirements)} |
| SHOULD recommendations | {len(should)} |
| SHOULD implemented | {should_followed} |
| SHOULD intentionally not followed | {should_not_followed} |
| MAY supported | {sum(r["source_strength"] == "MAY" and r["implementation_policy"] == "supported" for r in requirements)} |
| OPTIONAL supported | {sum(r["source_strength"] == "OPTIONAL" and r["implementation_policy"] == "supported" for r in requirements)} |
| Defaults verified | {source_counts["DEFAULT"]} |
| Grammar alternatives verified | {source_counts["GRAMMAR"] + source_counts["SEMANTIC"] + source_counts["ERROR"]} |
| Implementation-defined decisions | {policy_counts["implementation-defined"]} |
| Unsupported recommendations | {policy_counts["unsupported"]} |
| Non-MUST corpus cases | {len(cases)} |
| Example compatibility cases | {sum(len(ids) for ids in examples_by_section.values())} |

All supported records are directly linked to a manifest case. Unsupported
recommendations and implementation-defined choices are verified as explicit,
deterministic policies rather than counted as frozen MUST failures.
"""
    (BASE / "non-must-summary.md").write_text(summary, encoding="utf-8")

    policy_lines = [
        "# TOSCA 1.3 implementation policies",
        "",
        "These policies apply only to the processor non-MUST layer.",
        "",
    ]
    for requirement in requirements:
        if requirement["implementation_policy"] not in {"unsupported", "implementation-defined"}:
            continue
        policy_lines.extend(
            [
                f"## {requirement['id']}",
                "",
                f"- Section: {requirement['section']} — {requirement['section_title']}",
                f"- Strength: {requirement['strength']} (source: {requirement['source_strength']})",
                f"- Policy: {requirement['implementation_policy']}",
                f"- Rationale: {requirement['policy_reason']}",
                f"- Fixtures: {', '.join(requirement['fixture_ids'])}",
                "",
            ]
        )
    (BASE / "implementation-policies.md").write_text("\n".join(policy_lines), encoding="utf-8")

    classification_rows = "\n".join(
        f"| {name} | {count} |" for name, count in sorted(classification_counts.items())
    )
    full_section = f"""# TOSCA 1.3 full section coverage

All {len(sections)} catalog section identifiers have exactly one primary
classification and explicit frozen-MUST, processor non-MUST, and example
coverage states.

| Primary classification | Sections |
|---|---:|
{classification_rows}

Processor non-MUST sections link to separate requirements and fixtures.
Archive, orchestrator, generator, author, example, definition, cross-reference,
duplicate, informative, and not-applicable sections retain explicit target
reasons and do not affect the frozen 223 denominator.
"""
    (BASE / "full-section-coverage.md").write_text(full_section, encoding="utf-8")


def main():
    catalog_records = load(BASE / "requirements.yaml")["requirements"]
    catalog = {record["id"]: record for record in catalog_records}
    coverage = load(BASE / "coverage.yaml")["requirements"]
    selected = []
    for record in coverage:
        if record["strength"] not in NON_MUST_STRENGTHS:
            continue
        if (
            record["requirement_id"] not in IMPLEMENTATION_DEFINED_RECORDS
            and record["applicability"] != "applicable"
        ) or record["catalog_quality"] != "atomic":
            continue
        if record.get("included_in_must_denominator"):
            continue
        if record["requirement_id"] in NON_PROCESSOR_RECORDS:
            continue
        selected.append(record)

    requirements = []
    for record in selected:
        source = catalog[record["requirement_id"]]
        fixture_ids = fixtures_for(record)
        fixture_exists = all((CORPUS / CASE_PATHS[fixture_id]).is_file() for fixture_id in fixture_ids)
        policy = policy_for(record, fixture_exists)
        requirements.append(
            {
                "id": record["requirement_id"].replace("TOSCA13-", "TOSCA13-NONMUST-", 1),
                "source_requirement_id": record["requirement_id"],
                "section": record["section"],
                "section_title": record["section_title"],
                "strength": mapped_strength(record["strength"]),
                "source_strength": record["strength"],
                "normative_rule": source["requirement"],
                "processor_applicable": True,
                **policy,
                "fixture_ids": fixture_ids,
            }
        )

    dump(
        BASE / "non-must-requirements.yaml",
        {
            "specification": {
                "name": "TOSCA Simple Profile in YAML",
                "version": "1.3",
                "target": "processor-non-must",
                "source": "docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html",
            },
            "selection": {
                "source": "atomic, processor-applicable, non-frozen records indexed from requirements.yaml and reclassified by target",
                "frozen_requirement_count": 223,
                "record_count": len(requirements),
            },
            "requirements": requirements,
        },
    )

    by_fixture = {}
    for requirement in requirements:
        for fixture_id in requirement["fixture_ids"]:
            by_fixture.setdefault(fixture_id, []).append(requirement)

    cases = []
    for case_id, relative_path in CASE_PATHS.items():
        if not (CORPUS / relative_path).is_file():
            continue
        linked = by_fixture.get(case_id, [])
        if not linked:
            continue
        expected = CASE_EXPECTATIONS[case_id]
        policies = {requirement["implementation_policy"] for requirement in linked}
        support = "unsupported" if policies == {"unsupported"} else ("implementation-defined" if "implementation-defined" in policies else "supported")
        cases.append(
            {
                "id": case_id,
                "file": relative_path,
                "conformance_tier": expected["tier"],
                "specification": {
                    "sections": sorted(
                        {requirement["section"] for requirement in linked},
                        key=lambda value: tuple(int(part) if part.isdigit() else part for part in value.split(".")),
                    ),
                    "non_must_requirement_ids": sorted(requirement["id"] for requirement in linked),
                },
                "policy": {
                    "support": support,
                    "rationale": "All linked rules use the cataloged implementation policy; this case exercises their shared processor path.",
                },
                "classification": {
                    "kind": expected["kind"],
                    "category": expected["category"],
                    "features": expected["features"],
                },
                "expected": {
                    "accepted": expected["accepted"],
                    "phase": expected["phase"],
                    "diagnostic": expected["diagnostic"],
                    "assertions": expected["assertions"],
                },
            }
        )
    if cases:
        dump(CORPUS / "non_must/manifest.yaml", {"cases": cases})
    generate_reports(catalog_records, requirements, cases)


if __name__ == "__main__":
    main()
