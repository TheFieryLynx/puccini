#!/usr/bin/env python3
"""
Generate the TOSCA 1.3 normative-profile three-way comparison.

The prose model is extracted from the pinned OASIS Standard HTML. A small set
of Definition tables in that Word-generated HTML has damaged indentation.
SPEC_OVERRIDES restates those tables as YAML and cites the exact section; it
does not use either machine-readable profile as a normative fallback.
"""

from __future__ import annotations

import argparse
import copy
import hashlib
import importlib.util
import pathlib
import re
import sys
import textwrap
from typing import Any

import yaml


ROOT = pathlib.Path(__file__).resolve().parents[2]
SPEC_PATH = ROOT / "docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html"
OASIS_ROOT = ROOT / "third_party/oasis/tosca-simple-profile/1.3"
PUCCINI_ROOT = ROOT / "assets/tosca/profiles/simple/1.3"
PUCCINI_IMPLICIT_DATA = ROOT / "assets/tosca/profiles/implicit/1.3/data.yaml"
CORPUS_MANIFEST = ROOT / "tests/corpus/tosca_1_3/oasis_community/manifest.yaml"
COMPARISON_PATH = ROOT / "docs/conformance/tosca-1.3/profile-three-way-comparison.yaml"
DISCREPANCIES_YAML_PATH = ROOT / "docs/conformance/tosca-1.3/profile-discrepancies.yaml"
DISCREPANCIES_MD_PATH = ROOT / "docs/conformance/tosca-1.3/profile-discrepancies.md"

SOURCE_KINDS = {
    "data_types": "data_type",
    "artifact_types": "artifact_type",
    "capability_types": "capability_type",
    "relationship_types": "relationship_type",
    "interface_types": "interface_type",
    "node_types": "node_type",
    "group_types": "group_type",
    "policy_types": "policy_type",
}

SEMANTIC_KEYS = {
    "derived_from",
    "properties",
    "attributes",
    "capabilities",
    "requirements",
    "valid_source_types",
    "valid_target_types",
    "interfaces",
    "operations",
    "artifacts",
    "members",
    "targets",
    "mime_type",
    "file_ext",
    "constraints",
    "entry_schema",
}

TYPE_HEADING_RE = re.compile(r"^tosca\.[A-Za-z0-9_.]+$")

# These are transcription repairs or explicit interpretations of normative
# tables. Each entry is independently checkable against the cited local HTML.
SPEC_OVERRIDES_YAML = r"""
5.3.2:
  tosca.datatypes.json:
    derived_from: string
5.3.4:
  tosca.datatypes.xml:
    derived_from: string
5.3.6:
  tosca.datatypes.Credential:
    derived_from: tosca.datatypes.Root
    properties:
      protocol: {type: string, required: false}
      token_type: {type: string, default: password}
      token: {type: string}
      keys:
        type: map
        required: false
        entry_schema: {type: string}
      user: {type: string, required: false}
5.3.8:
  tosca.datatypes.network.NetworkInfo:
    derived_from: tosca.datatypes.Root
    properties:
      network_name: {type: string, required: false}
      network_id: {type: string, required: false}
      addresses:
        type: list
        required: false
        entry_schema: {type: string}
5.3.9:
  tosca.datatypes.network.PortInfo:
    derived_from: tosca.datatypes.Root
    properties:
      port_name: {type: string, required: false}
      port_id: {type: string, required: false}
      network_id: {type: string, required: false}
      mac_address: {type: string, required: false}
      addresses:
        type: list
        required: false
        entry_schema: {type: string}
5.3.10:
  tosca.datatypes.network.PortDef:
    derived_from: integer
    constraints:
      - in_range: [1, 65535]
5.4.3.1:
  tosca.artifacts.Deployment:
    derived_from: tosca.artifacts.Root
5.4.4.1:
  tosca.artifacts.Implementation:
    derived_from: tosca.artifacts.Root
5.4.5.1:
  tosca.artifacts.template:
    derived_from: tosca.artifacts.Root
5.5.7:
  tosca.capabilities.Endpoint:
    derived_from: tosca.capabilities.Root
    properties:
      protocol: {type: string, default: tcp}
      port: {type: tosca.datatypes.network.PortDef, required: false}
      secure: {type: boolean, required: false, default: false}
      url_path: {type: string, required: false}
      port_name: {type: string, required: false}
      network_name: {type: string, required: false, default: PRIVATE}
      initiator:
        type: string
        required: false
        default: source
        constraints:
          - valid_values: [source, target, peer]
      ports:
        type: map
        required: false
        constraints:
          - min_length: 1
        entry_schema: {type: tosca.datatypes.network.PortSpec}
    attributes:
      ip_address: {type: string}
5.5.13:
  tosca.capabilities.Scalable:
    derived_from: tosca.capabilities.Root
    properties:
      min_instances: {type: integer, default: 1}
      max_instances: {type: integer, default: 1}
      default_instances: {type: integer, required: false}
5.7.5:
  tosca.relationships.AttachesTo:
    derived_from: tosca.relationships.Root
    valid_target_types: [tosca.capabilities.Attachment]
    properties:
      location:
        type: string
        constraints:
          - min_length: 1
      device: {type: string, required: false}
    attributes:
      device: {type: string}
5.7.1:
  tosca.relationships.Root:
    attributes:
      tosca_id: {type: string}
      tosca_name: {type: string}
      state: {type: string, default: initial}
    interfaces:
      Configure: {type: tosca.interfaces.relationship.Configure}
5.8.4:
  tosca.interfaces.node.lifecycle.Standard:
    derived_from: tosca.interfaces.Root
    operations:
      create: {}
      configure: {}
      start: {}
      stop: {}
      delete: {}
5.8.5:
  tosca.interfaces.relationship.Configure:
    derived_from: tosca.interfaces.Root
    operations:
      pre_configure_source: {}
      pre_configure_target: {}
      post_configure_source: {}
      post_configure_target: {}
      add_target: {}
      add_source: {}
      target_changed: {}
      remove_target: {}
5.9.1:
  tosca.nodes.Root:
    attributes:
      tosca_id: {type: string}
      tosca_name: {type: string}
      state: {type: string, default: initial}
    capabilities:
      feature: {type: tosca.capabilities.Node}
    requirements:
      - dependency:
          capability: tosca.capabilities.Node
          node: tosca.nodes.Root
          relationship: tosca.relationships.DependsOn
          occurrences: [0, UNBOUNDED]
    interfaces:
      Standard: {type: tosca.interfaces.node.lifecycle.Standard}
5.9.3:
  tosca.nodes.Compute:
    derived_from: tosca.nodes.Abstract.Compute
    attributes:
      private_address: {type: string}
      public_address: {type: string}
      networks:
        type: map
        entry_schema: {type: tosca.datatypes.network.NetworkInfo}
      ports:
        type: map
        entry_schema: {type: tosca.datatypes.network.PortInfo}
    requirements:
      - local_storage:
          capability: tosca.capabilities.Attachment
          node: tosca.nodes.Storage.BlockStorage
          relationship: tosca.relationships.AttachesTo
          occurrences: [0, UNBOUNDED]
    capabilities:
      host:
        type: tosca.capabilities.Compute
        valid_source_types: [tosca.nodes.SoftwareComponent]
      endpoint: {type: tosca.capabilities.Endpoint.Admin}
      os: {type: tosca.capabilities.OperatingSystem}
      scalable: {type: tosca.capabilities.Scalable}
      binding: {type: tosca.capabilities.network.Bindable}
5.9.6:
  tosca.nodes.WebApplication:
    derived_from: tosca.nodes.Root
    properties:
      context_root: {type: string, required: false}
    capabilities:
      app_endpoint: {type: tosca.capabilities.Endpoint}
    requirements:
      - host:
          capability: tosca.capabilities.Compute
          node: tosca.nodes.WebServer
          relationship: tosca.relationships.HostedOn
5.9.8:
  tosca.nodes.Database:
    derived_from: tosca.nodes.Root
    properties:
      name: {type: string}
      port: {type: integer, required: false}
      user: {type: string, required: false}
      password: {type: string, required: false}
    requirements:
      - host:
          capability: tosca.capabilities.Compute
          node: tosca.nodes.DBMS
          relationship: tosca.relationships.HostedOn
    capabilities:
      database_endpoint: {type: tosca.capabilities.Endpoint.Database}
5.9.9:
  tosca.nodes.Abstract.Storage:
    derived_from: tosca.nodes.Root
    properties:
      name: {type: string}
      size:
        type: scalar-unit.size
        default: 0 MB
        constraints:
          - greater_or_equal: 0 MB
5.9.10:
  tosca.nodes.Storage.ObjectStorage:
    derived_from: tosca.nodes.Abstract.Storage
    properties:
      maxsize:
        type: scalar-unit.size
        required: false
        constraints:
          - greater_or_equal: 0 GB
    capabilities:
      storage_endpoint: {type: tosca.capabilities.Endpoint}
5.9.12:
  tosca.nodes.Container.Runtime:
    derived_from: tosca.nodes.SoftwareComponent
    capabilities:
      host:
        type: tosca.capabilities.Compute
        valid_source_types: [tosca.nodes.Container.Application]
      scalable: {type: tosca.capabilities.Scalable}
5.9.13:
  tosca.nodes.Container.Application:
    derived_from: tosca.nodes.Root
    requirements:
      - host:
          capability: tosca.capabilities.Compute
          node: tosca.nodes.Container.Runtime
          relationship: tosca.relationships.HostedOn
      - storage:
          capability: tosca.capabilities.Storage
      - network:
          capability: tosca.capabilities.Endpoint
5.10.1:
  tosca.groups.Root:
    interfaces:
      Standard: {type: tosca.interfaces.node.lifecycle.Standard}
8.5.1:
  tosca.nodes.network.Network:
    derived_from: tosca.nodes.Root
    properties:
      ip_version:
        type: integer
        required: false
        default: 4
        constraints:
          - valid_values: [4, 6]
      cidr: {type: string, required: false}
      start_ip: {type: string, required: false}
      end_ip: {type: string, required: false}
      gateway_ip: {type: string, required: false}
      network_name: {type: string, required: false}
      network_id: {type: string, required: false}
      segmentation_id: {type: string, required: false}
      network_type: {type: string, required: false}
      physical_network: {type: string, required: false}
      dhcp_enabled: {type: boolean, required: false, default: true}
    attributes:
      segmentation_id: {type: string}
    capabilities:
      link: {type: tosca.capabilities.network.Linkable}
8.5.5:
  tosca.relationships.network.BindsTo:
    derived_from: tosca.relationships.DependsOn
    valid_target_types: [tosca.capabilities.network.Bindable]
8.5.2:
  tosca.nodes.network.Port:
    derived_from: tosca.nodes.Root
    properties:
      ip_address: {type: string, required: false}
      order:
        type: integer
        required: true
        default: 0
        constraints:
          - greater_or_equal: 0
      is_default: {type: boolean, required: false, default: false}
      ip_range_start: {type: string, required: false}
      ip_range_end: {type: string, required: false}
    attributes:
      ip_address: {type: string}
    requirements:
      - link:
          capability: tosca.capabilities.network.Linkable
          relationship: tosca.relationships.network.LinksTo
      - binding:
          capability: tosca.capabilities.network.Bindable
          relationship: tosca.relationships.network.BindsTo
"""

SPEC_OVERRIDES = yaml.safe_load(SPEC_OVERRIDES_YAML)

DECISIONS = {
    "tosca.datatypes.network.NetworkInfo": (
        "docs/decisions/tosca-1.3-network-info-requiredness.md"
    ),
    "tosca.datatypes.network.PortInfo": (
        "docs/decisions/tosca-1.3-network-info-requiredness.md"
    ),
    "tosca.nodes.Storage.ObjectStorage": (
        "docs/conformance/tosca-1.3/ambiguities.md#tosca13-amb-012-objectstorage-maxsize-lower-bound-conflicts"
    ),
    "tosca.nodes.Storage.BlockStorage": (
        "docs/conformance/tosca-1.3/ambiguities.md#tosca13-amb-013-blockstorage-size-has-conditional-requiredness-and-volume_id-precedence"
    ),
    "tosca.groups.Root": (
        "docs/decisions/tosca-1.3-group-root-interface-contradiction.md"
    ),
    "tosca.nodes.network.Port": (
        "docs/decisions/0008-tosca-1.3-network-port-order-requiredness.md"
    ),
}

# Classification is intentionally explicit. Unknown discrepancies remain
# "unresolved" and make the later completeness gate fail.
CLASSIFICATIONS = {
    ("tosca.artifacts.Deployment.Image.VM", "derived_from"): "Puccini defect",
    ("tosca.capabilities.Scalable", "properties.default_instances.default"): (
        "community-profile defect"
    ),
    ("tosca.datatypes.json", "constraints"): "extension",
    ("tosca.datatypes.xml", "constraints"): "extension",
    (
        "tosca.datatypes.network.NetworkInfo",
        "properties.addresses.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.NetworkInfo",
        "properties.network_id.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.NetworkInfo",
        "properties.network_name.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.PortInfo",
        "properties.addresses.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.PortInfo",
        "properties.mac_address.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.PortInfo",
        "properties.network_id.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.PortInfo",
        "properties.port_id.required",
    ): "prose specification ambiguity",
    (
        "tosca.datatypes.network.PortInfo",
        "properties.port_name.required",
    ): "prose specification ambiguity",
    (
        "tosca.interfaces.relationship.Configure",
        "operations.remove_source",
    ): "extension",
    ("tosca.nodes.Abstract.Storage", "derived_from"): "Puccini defect",
    (
        "tosca.nodes.Abstract.Storage",
        "properties.size.constraints",
    ): "Puccini defect",
    (
        "tosca.nodes.Abstract.Storage",
        "properties.size.default",
    ): "Puccini defect",
    (
        "tosca.nodes.Abstract.Storage",
        "properties.size.required",
    ): "Puccini defect",
    (
        "tosca.nodes.Container.Application",
        "requirements.network.capability",
    ): "Puccini defect",
    (
        "tosca.nodes.Container.Runtime",
        "capabilities.host.type",
    ): "Puccini defect",
    (
        "tosca.nodes.Storage.BlockStorage",
        "properties.size",
    ): "prose specification ambiguity",
    ("tosca.groups.Root", "interfaces.Standard"): "prose specification ambiguity",
    (
        "tosca.groups.Root",
        "interfaces.Standard.type",
    ): "prose specification ambiguity",
    ("tosca.groups.Root", "interfaces"): "prose specification ambiguity",
    ("tosca.interfaces.Root", "derived_from"): "obsolete declaration",
    (
        "tosca.nodes.Abstract.Compute",
        "valid_source_types",
    ): "obsolete declaration",
    (
        "tosca.relationships.AttachesTo",
        "attributes.device",
    ): "community-profile defect",
    (
        "tosca.relationships.AttachesTo",
        "attributes",
    ): "community-profile defect",
}

INTRINSIC_SPEC_DISCREPANCIES = {
    "tosca.nodes.Storage.ObjectStorage": [
        {
            "path": "properties.maxsize.constraints.greater_or_equal",
            "specification": {
                "properties_table": "1 GB",
                "definition": "0 GB",
            },
            "oasis_community_profile": "0 GB",
            "puccini_profile": "0 GB",
            "classification": "prose specification ambiguity",
            "resolution_status": "documented",
        }
    ],
    "tosca.nodes.network.Port": [
        {
            "path": "properties.order.required",
            "specification": {
                "properties_table": False,
                "definition": True,
            },
            "oasis_community_profile": True,
            "puccini_profile": True,
            "classification": "prose specification ambiguity",
            "resolution_status": "interpreted",
        }
    ],
}

SECONDARY_CLASSIFICATIONS = {
    ("tosca.groups.Root", "interfaces"): ["obsolete declaration"],
}

REMEDIATED_FIELDS = {
    ("tosca.artifacts.Deployment.Image.VM", "derived_from"),
    ("tosca.nodes.Abstract.Storage", "derived_from"),
    ("tosca.nodes.Abstract.Storage", "properties.size.constraints"),
    ("tosca.nodes.Abstract.Storage", "properties.size.default"),
    ("tosca.nodes.Abstract.Storage", "properties.size.required"),
    (
        "tosca.nodes.Container.Application",
        "requirements.network.capability",
    ),
    ("tosca.nodes.Container.Runtime", "capabilities.host.type"),
}


def load_catalog_parser() -> tuple[Any, Any]:
    path = ROOT / "scripts/conformance/catalog_tosca13.py"
    spec = importlib.util.spec_from_file_location("tosca13_catalog_for_profile", path)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"cannot load {path}")
    module = importlib.util.module_from_spec(spec)
    sys.modules[spec.name] = module
    spec.loader.exec_module(module)
    return module.SpecificationParser, module.SOURCE


def category_for_name(name: str) -> str:
    prefixes = {
        "tosca.datatypes.": "data_type",
        "tosca.artifacts.": "artifact_type",
        "tosca.capabilities.": "capability_type",
        "tosca.relationships.": "relationship_type",
        "tosca.interfaces.": "interface_type",
        "tosca.nodes.": "node_type",
        "tosca.groups.": "group_type",
        "tosca.policies.": "policy_type",
    }
    for prefix, category in prefixes.items():
        if name.startswith(prefix):
            return category
    raise ValueError(f"cannot determine category for {name}")


def load_profile(root: pathlib.Path, implicit_data: pathlib.Path | None = None) -> dict[str, dict[str, Any]]:
    result: dict[str, dict[str, Any]] = {}
    files = sorted(root.glob("*.yaml"))
    if implicit_data is not None:
        files.append(implicit_data)
    for path in files:
        document = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
        for key, category in SOURCE_KINDS.items():
            for name, definition in (document.get(key) or {}).items():
                if implicit_data is not None and path == implicit_data and not name.startswith(
                    "tosca.datatypes."
                ):
                    continue
                result[name] = {
                    "category": category,
                    "file": str(path.relative_to(ROOT)),
                    "path": f"{key}.{name}",
                    "definition": definition or {},
                    "placement": "implicit import" if path == implicit_data else "direct",
                }
    return result


def extract_specification() -> dict[str, dict[str, Any]]:
    parser_class, source = load_catalog_parser()
    parser = parser_class()
    parser.feed(source.read_text(encoding="windows-1252"))
    headings = [
        heading
        for heading in parser.headings
        if TYPE_HEADING_RE.fullmatch(heading.title)
        and heading.section.startswith(("5.", "8.5."))
    ]
    result: dict[str, dict[str, Any]] = {}
    for heading in headings:
        override = SPEC_OVERRIDES.get(heading.section)
        if override is not None:
            definition = override[heading.title]
            extraction = "curated transcription of malformed normative table"
        else:
            tables = [
                table
                for table in parser.tables
                if table.section.startswith(f"{heading.section}.")
                and table.section_title == "Definition"
                and table.rows
                and table.rows[0]
            ]
            if not tables:
                raise ValueError(
                    f"{heading.section} {heading.title}: Definition table was not extracted"
                )
            snippet = textwrap.dedent(tables[-1].rows[0][0]).strip()
            parsed = yaml.safe_load(snippet)
            if not isinstance(parsed, dict) or heading.title not in parsed:
                raise ValueError(
                    f"{heading.section} {heading.title}: malformed Definition table "
                    "requires an explicit SPEC_OVERRIDES entry"
                )
            definition = parsed[heading.title] or {}
            extraction = "normative Definition table"
        result[heading.title] = {
            "category": category_for_name(heading.title),
            "section": heading.section,
            "title": heading.title,
            "line": heading.line,
            "definition": definition,
            "extraction": extraction,
        }
    return result


def normalize_requirement_list(value: Any) -> Any:
    if not isinstance(value, list):
        return value
    result: dict[str, Any] = {}
    for item in value:
        if isinstance(item, dict):
            result.update(item)
    return result


def normalize_value(value: Any, path: tuple[str, ...] = ()) -> Any:
    if value is None:
        return {}
    if isinstance(value, list):
        if path and path[-1] == "requirements":
            return normalize_value(normalize_requirement_list(value), path)
        return [normalize_value(item, path) for item in value]
    if not isinstance(value, dict):
        if path and path[-1] == "type":
            return {
                "PortDef": "tosca.datatypes.network.PortDef",
                "PortSpec": "tosca.datatypes.network.PortSpec",
                "NetworkInfo": "tosca.datatypes.network.NetworkInfo",
                "PortInfo": "tosca.datatypes.network.PortInfo",
            }.get(value, value)
        return value
    result: dict[str, Any] = {}
    for key, item in value.items():
        if key in {"description", "metadata"}:
            continue
        result[key] = normalize_value(item, path + (key,))

    # Capability definitions allow a type-name shorthand.
    if path and path[-1] == "capabilities":
        for name, definition in list(result.items()):
            if isinstance(definition, str):
                result[name] = {"type": definition}

    # TOSCA 1.3 property definitions default required to true.
    if path and path[-1] == "properties":
        for definition in result.values():
            if isinstance(definition, dict) and "required" not in definition:
                definition["required"] = True

    # Requirement occurrences default to [1, 1].
    if path and path[-1] == "requirements":
        for definition in result.values():
            if isinstance(definition, dict) and "occurrences" not in definition:
                definition["occurrences"] = [1, 1]
    return result


def semantic_definition(definition: dict[str, Any]) -> dict[str, Any]:
    normalized = normalize_value(definition)
    return {key: value for key, value in normalized.items() if key in SEMANTIC_KEYS}


def raw_semantic_value(value: Any) -> Any:
    if value is None:
        return {}
    if isinstance(value, list):
        return [raw_semantic_value(item) for item in value]
    if not isinstance(value, dict):
        return value
    return {
        key: raw_semantic_value(item)
        for key, item in value.items()
        if key not in {"description", "metadata"}
    }


def raw_semantic_definition(definition: dict[str, Any]) -> dict[str, Any]:
    normalized = raw_semantic_value(definition)
    return {key: value for key, value in normalized.items() if key in SEMANTIC_KEYS}


def merge_effective(parent: Any, child: Any) -> Any:
    if isinstance(parent, dict) and isinstance(child, dict):
        result = copy.deepcopy(parent)
        for key, value in child.items():
            result[key] = (
                merge_effective(result[key], value)
                if key in result
                else copy.deepcopy(value)
            )
        return result
    return copy.deepcopy(child)


def effective_definition(
    name: str, source: dict[str, dict[str, Any]], stack: tuple[str, ...] = ()
) -> dict[str, Any]:
    definition = semantic_definition(source[name]["definition"])
    parent_name = definition.get("derived_from")
    if parent_name in source and parent_name not in stack:
        parent = effective_definition(parent_name, source, stack + (name,))
        definition = merge_effective(parent, definition)
        definition["derived_from"] = parent_name
    return definition


def flatten_differences(
    left: Any, right: Any, path: tuple[str, ...] = ()
) -> list[tuple[str, Any, Any]]:
    if isinstance(left, dict) and isinstance(right, dict):
        differences: list[tuple[str, Any, Any]] = []
        for key in sorted(set(left) | set(right)):
            child_path = path + (key,)
            if key not in left:
                differences.append((".".join(child_path), None, right[key]))
            elif key not in right:
                differences.append((".".join(child_path), left[key], None))
            else:
                differences.extend(
                    flatten_differences(left[key], right[key], child_path)
                )
        return differences
    if left != right:
        return [(".".join(path), left, right)]
    return []


def comparison_status(
    name: str,
    left: dict[str, dict[str, Any]],
    right: dict[str, dict[str, Any]],
) -> str:
    if name not in left or name not in right:
        return "missing"
    if raw_semantic_definition(left[name]["definition"]) == raw_semantic_definition(
        right[name]["definition"]
    ):
        return "exact"
    left_direct = semantic_definition(left[name]["definition"])
    right_direct = semantic_definition(right[name]["definition"])
    if left_direct == right_direct:
        return "equivalent"
    if effective_definition(name, left) == effective_definition(name, right):
        return "equivalent"
    return "conflicting"


def sha256(path: pathlib.Path) -> str:
    digest = hashlib.sha256()
    digest.update(path.read_bytes())
    return digest.hexdigest()


def verify_pinned_integrity() -> tuple[str, list[str]]:
    commit = (OASIS_ROOT / "SOURCE_COMMIT").read_text(encoding="utf-8").strip()
    if not re.fullmatch(r"[0-9a-f]{40}", commit):
        raise ValueError(f"invalid SOURCE_COMMIT: {commit!r}")
    checked: list[str] = []
    for line in (OASIS_ROOT / "SHA256SUMS").read_text(encoding="utf-8").splitlines():
        expected, relative = line.split(maxsplit=1)
        relative = relative.lstrip("*")
        path = ROOT / relative
        actual = sha256(path)
        if actual != expected:
            raise ValueError(
                f"pinned-profile integrity failure for {relative}: "
                f"{actual} != {expected}"
            )
        checked.append(relative)
    return commit, checked


def corpus_evidence() -> dict[str, list[str]]:
    evidence: dict[str, list[str]] = {}
    if not CORPUS_MANIFEST.exists():
        return evidence
    manifest = yaml.safe_load(CORPUS_MANIFEST.read_text(encoding="utf-8")) or {}
    for case in manifest.get("cases", []):
        for name in case.get("types_covered", []):
            evidence.setdefault(name, []).append(case["id"])
    return evidence


def normalized_source(
    source: dict[str, dict[str, Any]]
) -> dict[str, dict[str, Any]]:
    return {
        name: {
            "category": item["category"],
            "definition": item["definition"],
        }
        for name, item in source.items()
    }


def classify_difference(name: str, path: str) -> str:
    if (name, path) in CLASSIFICATIONS:
        return CLASSIFICATIONS[(name, path)]
    # A parent path classification also covers its leaf values.
    candidates = [
        (prefix, classification)
        for (type_name, prefix), classification in CLASSIFICATIONS.items()
        if type_name == name and (path == prefix or path.startswith(f"{prefix}."))
    ]
    if candidates:
        return max(candidates, key=lambda item: len(item[0]))[1]
    return "unresolved"


def resolution_status(classification: str) -> str:
    if classification in {
        "semantically equivalent representation",
        "extension",
        "obsolete declaration",
    }:
        return "interpreted"
    if classification == "prose specification ambiguity":
        return "documented"
    if classification in {"Puccini defect", "community-profile defect"}:
        return "confirmed"
    return "unresolved"


def ancestor_names(
    name: str, source: dict[str, dict[str, Any]]
) -> list[str]:
    result: list[str] = []
    current = name
    seen: set[str] = set()
    while current in source and current not in seen:
        result.append(current)
        seen.add(current)
        parent = semantic_definition(source[current]["definition"]).get("derived_from")
        if not isinstance(parent, str):
            break
        current = parent
    return result


def classify_effective_difference(
    name: str,
    path: str,
    sources: tuple[dict[str, dict[str, Any]], ...],
) -> tuple[str, str | None]:
    candidates: list[str] = []
    for source in sources:
        for ancestor in ancestor_names(name, source):
            if ancestor not in candidates:
                candidates.append(ancestor)
    for ancestor in candidates:
        classification = classify_difference(ancestor, path)
        if classification != "unresolved":
            return classification, ancestor
    # A wrong or missing parent removes otherwise unrelated inherited fields.
    # Attribute the resulting effective-definition delta to that hierarchy
    # defect when no field-specific classification exists.
    for ancestor in candidates:
        classification = classify_difference(ancestor, "derived_from")
        if classification in {"Puccini defect", "community-profile defect"}:
            return classification, ancestor
    return "unresolved", None


def build_outputs() -> tuple[dict[str, Any], dict[str, Any], str]:
    commit, checked = verify_pinned_integrity()
    specification = extract_specification()
    oasis = load_profile(OASIS_ROOT)
    puccini = load_profile(PUCCINI_ROOT, PUCCINI_IMPLICIT_DATA)
    evidence = corpus_evidence()
    all_names = sorted(set(specification) | set(oasis) | set(puccini))

    types: list[dict[str, Any]] = []
    discrepancies: list[dict[str, Any]] = []
    for name in all_names:
        spec_item = specification.get(name)
        oasis_item = oasis.get(name)
        puccini_item = puccini.get(name)
        categories = {
            item["category"]
            for item in (spec_item, oasis_item, puccini_item)
            if item is not None
        }
        category = next(iter(categories)) if len(categories) == 1 else "unclassified"

        spec_semantic = (
            semantic_definition(spec_item["definition"]) if spec_item else None
        )
        oasis_semantic = (
            semantic_definition(oasis_item["definition"]) if oasis_item else None
        )
        puccini_semantic = (
            semantic_definition(puccini_item["definition"]) if puccini_item else None
        )
        field_paths: set[str] = set()
        for left, right in (
            (spec_semantic, oasis_semantic),
            (spec_semantic, puccini_semantic),
            (oasis_semantic, puccini_semantic),
        ):
            if left is not None and right is not None:
                field_paths.update(path for path, _, _ in flatten_differences(left, right))
        field_paths.update(
            path
            for type_name, path in REMEDIATED_FIELDS
            if type_name == name
        )

        type_differences: list[dict[str, Any]] = []
        for path in sorted(field_paths):
            def value_at(value: Any, dotted: str) -> Any:
                current = value
                for component in dotted.split("."):
                    if not isinstance(current, dict) or component not in current:
                        return None
                    current = current[component]
                return current

            classification = classify_difference(name, path)
            specification_value = value_at(spec_semantic, path)
            oasis_value = value_at(oasis_semantic, path)
            puccini_value = value_at(puccini_semantic, path)
            puccini_remediated = (
                (name, path) in REMEDIATED_FIELDS
                and classification == "Puccini defect"
                and specification_value == puccini_value
            )
            difference = {
                "path": path,
                "specification": specification_value,
                "oasis_community_profile": oasis_value,
                "puccini_profile": puccini_value,
                "classification": classification,
                "secondary_classifications": SECONDARY_CLASSIFICATIONS.get(
                    (name, path), []
                ),
                "resolution_status": (
                    "remediated"
                    if puccini_remediated
                    else resolution_status(classification)
                ),
                "puccini_remediated": puccini_remediated,
            }
            type_differences.append(difference)
            discrepancies.append(
                {
                    "id": f"PROFILE13-{len(discrepancies) + 1:03d}",
                    "canonical_name": name,
                    **difference,
                    "specification_section": (
                        spec_item["section"] if spec_item else None
                    ),
                    "decision_record": DECISIONS.get(name),
                }
            )

        for intrinsic in INTRINSIC_SPEC_DISCREPANCIES.get(name, []):
            path = intrinsic["path"]
            def value_at(value: Any, dotted: str) -> Any:
                current = value
                for component in dotted.split("."):
                    if not isinstance(current, dict) or component not in current:
                        return None
                    current = current[component]
                return current

            difference = {
                **intrinsic,
                "oasis_community_profile": intrinsic.get(
                    "oasis_community_profile", value_at(oasis_semantic, path)
                ),
                "puccini_profile": intrinsic.get(
                    "puccini_profile", value_at(puccini_semantic, path)
                ),
                "secondary_classifications": [],
            }
            type_differences.append(difference)
            discrepancies.append(
                {
                    "id": f"PROFILE13-{len(discrepancies) + 1:03d}",
                    "canonical_name": name,
                    **difference,
                    "specification_section": (
                        spec_item["section"] if spec_item else None
                    ),
                    "decision_record": DECISIONS.get(name),
                }
            )

        def effective_hash(
            source: dict[str, dict[str, Any]], source_name: str
        ) -> str | None:
            if source_name not in source:
                return None
            payload = yaml.safe_dump(
                effective_definition(source_name, source),
                sort_keys=True,
                allow_unicode=True,
            ).encode("utf-8")
            return hashlib.sha256(payload).hexdigest()

        effective_spec = (
            effective_definition(name, specification) if spec_item else None
        )
        effective_oasis = effective_definition(name, oasis) if oasis_item else None
        effective_puccini = (
            effective_definition(name, puccini) if puccini_item else None
        )
        effective_paths: set[str] = set()
        for left, right in (
            (effective_spec, effective_oasis),
            (effective_spec, effective_puccini),
            (effective_oasis, effective_puccini),
        ):
            if left is not None and right is not None:
                effective_paths.update(
                    path for path, _, _ in flatten_differences(left, right)
                )

        def effective_value_at(value: Any, dotted: str) -> Any:
            current = value
            for component in dotted.split("."):
                if not isinstance(current, dict) or component not in current:
                    return None
                current = current[component]
            return current

        effective_differences: list[dict[str, Any]] = []
        for path in sorted(effective_paths):
            classification, inherited_from = classify_effective_difference(
                name, path, (specification, oasis, puccini)
            )
            effective_differences.append(
                {
                    "path": path,
                    "specification": effective_value_at(effective_spec, path),
                    "oasis_community_profile": effective_value_at(
                        effective_oasis, path
                    ),
                    "puccini_profile": effective_value_at(
                        effective_puccini, path
                    ),
                    "classification": classification,
                    "resolution_status": resolution_status(classification),
                    "inherited_from": (
                        inherited_from if inherited_from != name else None
                    ),
                }
            )

        types.append(
            {
                "canonical_name": name,
                "category": category,
                "specification": {
                    "section": spec_item["section"] if spec_item else "",
                    "title": spec_item["title"] if spec_item else "",
                    "present": spec_item is not None,
                    "extraction": spec_item["extraction"] if spec_item else None,
                    "effective_definition_sha256": effective_hash(
                        specification, name
                    ),
                },
                "oasis_community_profile": {
                    "present": oasis_item is not None,
                    "file": oasis_item["file"] if oasis_item else "",
                    "path": oasis_item["path"] if oasis_item else "",
                    "effective_definition_sha256": effective_hash(oasis, name),
                },
                "puccini_profile": {
                    "present": puccini_item is not None,
                    "file": puccini_item["file"] if puccini_item else "",
                    "path": puccini_item["path"] if puccini_item else "",
                    "placement": puccini_item["placement"] if puccini_item else "",
                    "effective_definition_sha256": effective_hash(puccini, name),
                },
                "comparison": {
                    "specification_vs_oasis": comparison_status(
                        name, specification, oasis
                    ),
                    "specification_vs_puccini": comparison_status(
                        name, specification, puccini
                    ),
                    "oasis_vs_puccini": comparison_status(name, oasis, puccini),
                },
                "effective_comparison": {
                    "specification_vs_oasis": (
                        "exact"
                        if name in specification
                        and name in oasis
                        and effective_definition(name, specification)
                        == effective_definition(name, oasis)
                        else "conflicting"
                        if name in specification and name in oasis
                        else "missing"
                    ),
                    "specification_vs_puccini": (
                        "exact"
                        if name in specification
                        and name in puccini
                        and effective_definition(name, specification)
                        == effective_definition(name, puccini)
                        else "conflicting"
                        if name in specification and name in puccini
                        else "missing"
                    ),
                    "oasis_vs_puccini": (
                        "exact"
                        if name in oasis
                        and name in puccini
                        and effective_definition(name, oasis)
                        == effective_definition(name, puccini)
                        else "conflicting"
                        if name in oasis and name in puccini
                        else "missing"
                    ),
                },
                "field_differences": type_differences,
                "effective_field_differences": effective_differences,
                "corpus_cases": sorted(evidence.get(name, [])),
                "decision_record": DECISIONS.get(name),
            }
        )

    def counts(source: dict[str, dict[str, Any]]) -> dict[str, int]:
        result: dict[str, int] = {}
        for item in source.values():
            result[item["category"]] = result.get(item["category"], 0) + 1
        return dict(sorted(result.items()))

    pair_counts = {
        pair: {
            status: sum(
                1
                for item in types
                if item["comparison"][pair] == status
            )
            for status in ("exact", "equivalent", "conflicting", "missing")
        }
        for pair in (
            "specification_vs_oasis",
            "specification_vs_puccini",
            "oasis_vs_puccini",
        )
    }
    comparison_document = {
        "schema_version": 1,
        "audit": {
            "applicable_specification": "TOSCA Simple Profile in YAML Version 1.3",
            "normative_source": str(SPEC_PATH.relative_to(ROOT)),
            "oasis_community_upstream_repository": (
                "https://github.com/oasis-open/tosca-community-contributions"
            ),
            "oasis_community_upstream_commit": commit,
            "oasis_community_profile_root": str(OASIS_ROOT.relative_to(ROOT)),
            "puccini_profile_root": str(PUCCINI_ROOT.relative_to(ROOT)),
            "integrity_verified_files": checked,
            "comparison_scope": [
                "exact case-sensitive type name",
                "type category",
                "derived_from",
                "properties and property types",
                "required flags and defaults",
                "constraints and entry schemas",
                "attributes",
                "capabilities",
                "requirements and occurrences",
                "valid source and target types",
                "interfaces and operations",
                "artifacts",
                "group members",
                "policy targets",
                "workflow-related definitions",
                "effective inherited definition",
            ],
        },
        "summary": {
            "type_counts": {
                "specification": counts(specification),
                "oasis_community_profile": counts(oasis),
                "puccini_effective_profile": counts(puccini),
            },
            "direct_puccini_type_count": sum(
                1 for item in puccini.values() if item["placement"] == "direct"
            ),
            "effective_puccini_type_count": len(puccini),
            "comparison_counts": pair_counts,
            "discrepancy_count": len(discrepancies),
        },
        "types": types,
    }
    discrepancy_document = {
        "schema_version": 1,
        "applicable_specification": "TOSCA Simple Profile in YAML Version 1.3",
        "normative_source": str(SPEC_PATH.relative_to(ROOT)),
        "oasis_community_upstream_commit": commit,
        "allowed_classifications": [
            "Puccini defect",
            "community-profile defect",
            "prose specification ambiguity",
            "semantically equivalent representation",
            "version mismatch",
            "extension",
            "obsolete declaration",
            "unresolved",
        ],
        "discrepancies": discrepancies,
    }
    markdown = render_markdown(comparison_document, discrepancy_document)
    return comparison_document, discrepancy_document, markdown


def render_markdown(
    comparison: dict[str, Any], discrepancy_document: dict[str, Any]
) -> str:
    summary = comparison["summary"]
    discrepancies = discrepancy_document["discrepancies"]
    classification_counts: dict[str, int] = {}
    for discrepancy in discrepancies:
        classification = discrepancy["classification"]
        classification_counts[classification] = (
            classification_counts.get(classification, 0) + 1
        )
    lines = [
        "# TOSCA 1.3 normative profile discrepancies",
        "",
        "This report compares the OASIS Standard prose, the pinned OASIS community",
        "machine-readable profile, and Puccini's effective bundled TOSCA 1.3 profile.",
        "The prose Standard is the conformance authority.",
        "",
        f"- Pinned upstream commit: `{comparison['audit']['oasis_community_upstream_commit']}`",
        f"- Types in the union: **{len(comparison['types'])}**",
        f"- Field discrepancies: **{len(discrepancies)}**",
        "",
        "## Classification summary",
        "",
        "| Classification | Fields |",
        "|---|---:|",
    ]
    for classification, count in sorted(classification_counts.items()):
        lines.append(f"| {classification} | {count} |")
    lines.extend(
        [
            "",
            "## Pairwise type summary",
            "",
            "| Pair | Exact | Equivalent | Conflicting | Missing |",
            "|---|---:|---:|---:|---:|",
        ]
    )
    for pair, counts in summary["comparison_counts"].items():
        lines.append(
            f"| {pair} | {counts['exact']} | {counts['equivalent']} | "
            f"{counts['conflicting']} | {counts['missing']} |"
        )
    lines.extend(["", "## Field discrepancies", ""])
    for discrepancy in discrepancies:
        lines.extend(
            [
                f"### {discrepancy['id']}: `{discrepancy['canonical_name']}` "
                f"`{discrepancy['path']}`",
                "",
                f"- Specification section: §{discrepancy['specification_section']}",
                f"- Classification: **{discrepancy['classification']}**",
                f"- Resolution: **{discrepancy['resolution_status']}**",
                f"- Specification: `{discrepancy['specification']}`",
                f"- OASIS community: `{discrepancy['oasis_community_profile']}`",
                f"- Puccini: `{discrepancy['puccini_profile']}`",
            ]
        )
        if discrepancy["decision_record"]:
            lines.append(f"- Decision/ambiguity record: `{discrepancy['decision_record']}`")
        if discrepancy["secondary_classifications"]:
            lines.append(
                "- Secondary classification: "
                + ", ".join(
                    f"**{classification}**"
                    for classification in discrepancy["secondary_classifications"]
                )
            )
        lines.append("")
    lines.extend(
        [
            "## Effective-definition method",
            "",
            "The audit canonicalizes capability shorthand, property requiredness defaults,",
            "requirement occurrence defaults, and then recursively overlays inherited",
            "definitions. Descriptions and implementation metadata are intentionally excluded",
            "from semantic equality. Source placement remains recorded, including Puccini's",
            "implicit `json` and `xml` declarations.",
            "",
        ]
    )
    return "\n".join(lines)


def dump_yaml(document: dict[str, Any]) -> str:
    return yaml.safe_dump(
        document,
        sort_keys=False,
        allow_unicode=True,
        width=100,
    )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--check",
        action="store_true",
        help="fail if tracked outputs differ from freshly generated content",
    )
    args = parser.parse_args()
    comparison, discrepancies, markdown = build_outputs()
    outputs = {
        COMPARISON_PATH: dump_yaml(comparison),
        DISCREPANCIES_YAML_PATH: dump_yaml(discrepancies),
        DISCREPANCIES_MD_PATH: markdown,
    }
    if args.check:
        stale = [
            str(path.relative_to(ROOT))
            for path, content in outputs.items()
            if not path.exists() or path.read_text(encoding="utf-8") != content
        ]
        if stale:
            print("stale TOSCA 1.3 profile audit outputs:", file=sys.stderr)
            for path in stale:
                print(f"  {path}", file=sys.stderr)
            return 1
        print(
            f"TOSCA 1.3 profile audit: {len(comparison['types'])} types, "
            f"{len(discrepancies['discrepancies'])} field discrepancies"
        )
        return 0
    for path, content in outputs.items():
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding="utf-8")
    print(
        f"generated {len(comparison['types'])} type comparisons and "
        f"{len(discrepancies['discrepancies'])} field discrepancies"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
