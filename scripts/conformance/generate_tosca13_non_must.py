#!/usr/bin/env python3
"""Generate the separate TOSCA 1.3 processor non-MUST catalog.

The frozen 223-requirement matrix is an input and is never modified here.
Catalog records are selected from the pinned-spec extraction only when the
existing audit classifies them as atomic and processor-applicable. Explicit
target corrections below separate author, orchestrator, and archive rules.
"""

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
    "TOSCA13-5.9.1.4-001",
    "TOSCA13-5.9.3.4-001",
    "TOSCA13-5.9.4.4-001",
    "TOSCA13-6.3-001",
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
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-DEFINITIONS": "non_must/grammar-alternatives/definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TYPES": "non_must/grammar-alternatives/type-definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-TEMPLATES": "non_must/grammar-alternatives/template-definitions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-FUNCTIONS": "non_must/grammar-alternatives/functions-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-WORKFLOWS": "non_must/grammar-alternatives/workflows-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-GRAMMAR-PROFILE": "non_must/grammar-alternatives/normative-profile-valid.yaml",
    "TOSCA13-NONMUST-CORPUS-POLICY-NETWORK": "non_must/implementation-defined/network-disabled-policy.yaml",
    "TOSCA13-NONMUST-CORPUS-POLICY-EXTENSIONS": "non_must/implementation-defined/unknown-extension-policy.yaml",
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
}


def load(path):
    return yaml.safe_load(path.read_text(encoding="utf-8"))


def dump(path, value):
    path.write_text(
        yaml.safe_dump(value, sort_keys=False, allow_unicode=True, width=140),
        encoding="utf-8",
    )


def mapped_strength(source_strength):
    if source_strength == "GRAMMAR":
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
        return "TOSCA13-NONMUST-CORPUS-MAY-IMPORT-NAMESPACE"
    if requirement_id == "TOSCA13-3.6.17.3-003":
        return "TOSCA13-NONMUST-CORPUS-MAY-OPERATION-INPUT"
    if section == "3.10.2.2":
        return "TOSCA13-NONMUST-CORPUS-MAY-TYPE-ONLY"
    if section.startswith("4."):
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
    return [fixture]


def policy_for(record, fixture_exists):
    source_strength = record["strength"]
    requirement_id = record["requirement_id"]
    unsupported = requirement_id in UNSUPPORTED_RECOMMENDATIONS

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


def main():
    catalog_records = load(BASE / "requirements.yaml")["requirements"]
    catalog = {record["id"]: record for record in catalog_records}
    coverage = load(BASE / "coverage.yaml")["requirements"]
    selected = []
    for record in coverage:
        if record["strength"] not in NON_MUST_STRENGTHS:
            continue
        if record["applicability"] != "applicable" or record["catalog_quality"] != "atomic":
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


if __name__ == "__main__":
    main()
