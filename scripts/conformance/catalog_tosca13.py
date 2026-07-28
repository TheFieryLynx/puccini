#!/usr/bin/env python3
"""
Build and verify the TOSCA Simple Profile in YAML 1.3 normative baseline.

This tool reads only the pinned local OASIS HTML specification. It does not
inspect Puccini production code and does not infer requirements from current
processor behavior.
"""

from __future__ import annotations

import argparse
import collections
import dataclasses
import hashlib
import html
from html.parser import HTMLParser
import pathlib
import re
import sys
from typing import Any, Iterable

import yaml


REPO_ROOT = pathlib.Path(__file__).resolve().parents[2]
SOURCE = REPO_ROOT / "docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html"
OUTPUT_DIR = REPO_ROOT / "docs/conformance/tosca-1.3"
REQUIREMENTS_PATH = OUTPUT_DIR / "requirements.yaml"
SUMMARY_PATH = OUTPUT_DIR / "catalog-summary.md"
AMBIGUITIES_PATH = OUTPUT_DIR / "ambiguities.md"
EXTRACTION_REPORT_PATH = OUTPUT_DIR / "extraction-report.md"

NORMATIVE_WORDS = (
    "MUST NOT",
    "SHALL NOT",
    "SHOULD NOT",
    "MUST",
    "SHALL",
    "REQUIRED",
    "SHOULD",
    "RECOMMENDED",
    "MAY",
    "OPTIONAL",
)

SECTION_RE = re.compile(r"^\s*(\d+(?:\.\d+)*)(?:\s+|$)(.*)$")
NORMATIVE_RE = re.compile(
    r"\b(MUST\s+NOT|SHALL\s+NOT|SHOULD\s+NOT|MUST|SHALL|REQUIRED|"
    r"SHOULD|RECOMMENDED|MAY|OPTIONAL)\b"
)
SPACE_RE = re.compile(r"\s+")
SENTENCE_BOUNDARY_RE = re.compile(r"(?<=[.!?;])\s+(?=[A-Z0-9(<`\[])")
SEMANTIC_CUE_RE = re.compile(
    r"\b("
    r"required|optional|only|default|valid|invalid|allowed|not allowed|"
    r"cannot|can not|error|ignored|exclusive|mutually|one or more|"
    r"at least|at most|exactly|between|resolved|resolution|namespace|"
    r"imported|derived|inherits?|refin(?:e|ement)|compatible|constraint|"
    r"short notation|long notation|grammar|syntax|shall|must|should|may"
    r")\b",
    re.IGNORECASE,
)
FUNCTION_SEMANTIC_CUE_RE = re.compile(
    r"\b(return|retrieve|resolve|evaluate|argument|parameter|context|"
    r"value|search|concatenate|join|token|result|function)\b",
    re.IGNORECASE,
)

VALID_STRENGTHS = {
    "MUST",
    "MUST_NOT",
    "SHALL",
    "SHALL_NOT",
    "REQUIRED",
    "SHOULD",
    "SHOULD_NOT",
    "RECOMMENDED",
    "MAY",
    "OPTIONAL",
    "GRAMMAR",
    "DEFAULT",
    "ERROR",
    "SEMANTIC",
}
VALID_VALIDATION_KINDS = {
    "lexical",
    "yaml-type",
    "structural",
    "reference-resolution",
    "namespace",
    "import",
    "hierarchy",
    "inheritance",
    "refinement",
    "assignment",
    "function",
    "constraint",
    "semantic",
    "normalization",
    "csar",
    "orchestrator-only",
}
VALID_EXPECTED_BEHAVIORS = {
    "accept",
    "reject",
    "resolve",
    "normalize",
    "preserve",
    "evaluate",
    "not-applicable",
}
TEST_OBLIGATION_KEYS = (
    "positive",
    "negative",
    "boundary",
    "inheritance",
    "imports",
    "normalization",
)

PROSE_SCOPE_PREFIXES = ("3", "4", "5", "6", "7", "8", "12", "13", "14")
RUNTIME_ONLY_PREFIXES = ("7", "8", "12", "13")
EXCLUDED_TOP_LEVEL_PREFIXES = ("2", "9", "10", "11")
APPENDIX_START_LINE = 48548


MANUAL_AMBIGUITIES: list[dict[str, Any]] = [
    {
        "id": "TOSCA13-AMB-001",
        "title": "Conformance clause 14.3(d) cites the wrong section number",
        "sections": ["14.3", "3.1", "3.2", "3.3", "3.6"],
        "issue": (
            "Clause 14.3(d) calls section 3.2 “Parameter and property type”, "
            "but section 3.2 is “Using Namespaces”; “Parameter and property "
            "types” is section 3.3."
        ),
        "interpretations": [
            "Treat the parenthetical title as controlling and read the citation as section 3.3.",
            "Treat the numeric citation as controlling and include namespace errors from section 3.2.",
            "Include both sections until an erratum resolves the conflict.",
        ],
        "catalog_policy": "The catalog includes error rules from sections 3.1, 3.2, 3.3, and 3.6.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-002",
        "title": "Conformance clause 14.3(e) cites a non-existent section",
        "sections": ["14.3", "5.4.9.3"],
        "issue": (
            "Clause 14.3(e) requires string normalization described in section "
            "5.4.9.3, but the published document has no section 5.4.9.3 and no "
            "other occurrence describing string normalization."
        ),
        "interpretations": [
            "The requirement is normative but its referenced algorithm is missing.",
            "The citation is a stale reference to text removed or renumbered before publication.",
        ],
        "catalog_policy": "Record the processor normalization obligation as ambiguous and algorithm-unspecified.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-003",
        "title": "Implementation definitions contain stale service-template cross-references",
        "sections": ["1.4", "3.8", "3.9", "3.10"],
        "issue": (
            "Section 1.4 points to section 3.9 for Service Template definition "
            "and section 3.8 for Topology Template definition; the actual "
            "sections are 3.10 and 3.9 respectively."
        ),
        "interpretations": [
            "Follow the section titles and current numbering.",
            "Treat the references as historical numbering with no semantic effect.",
        ],
        "catalog_policy": "Use sections 3.10 and 3.9 while retaining this discrepancy.",
        "manual_review": False,
    },
    {
        "id": "TOSCA13-AMB-004",
        "title": "Requirement relationship type is both required and optional",
        "sections": ["3.7.3.1.1"],
        "issue": (
            "The grammar table marks relationship keyname `type` as Required "
            "“yes” while its description calls the keyname optional."
        ),
        "interpretations": [
            "The Required column controls and `type` is mandatory.",
            "The prose controls and the relationship may use other supported notation.",
        ],
        "catalog_policy": "Emit both atomic records and flag their conflict; do not select an interpretation.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-005",
        "title": "Artifact type properties conflict between Required column and prose",
        "sections": ["3.7.4.1"],
        "issue": (
            "The `mime_type` and `file_ext` rows are marked Required “no” but "
            "their descriptions call the properties required."
        ),
        "interpretations": [
            "Treat both properties as optional according to the table marker.",
            "Treat both properties as mandatory according to the description.",
        ],
        "catalog_policy": "Emit both table-marker and description records and flag the conflict.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-006",
        "title": "Parameter `type` requirement is context-dependent",
        "sections": ["3.6.10.2", "3.6.14.1"],
        "issue": (
            "The parameter table marks `type` optional while describing it as "
            "required for property definitions but not parameter definitions; "
            "the shared schema is reused in multiple contexts."
        ),
        "interpretations": [
            "Require `type` only when the schema is used as a property definition.",
            "Allow omission in parameter contexts and derive the type where another rule permits it.",
        ],
        "catalog_policy": "Record the conditional rules independently.",
        "manual_review": False,
    },
    {
        "id": "TOSCA13-AMB-007",
        "title": "Version selector short name and URI forms",
        "sections": ["3.10.3.1"],
        "issue": (
            "The section presents both `tosca_simple_yaml_1_3` and a fully "
            "qualified URI as examples but does not state an exhaustive lexical set."
        ),
        "interpretations": [
            "Accept both forms.",
            "Accept any profile identifier that unambiguously selects the grammar.",
            "Limit a product support profile to the short form while documenting that restriction.",
        ],
        "catalog_policy": "Catalog the explicit grammar-selection duty without inventing an exhaustive alias rule.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-008",
        "title": "TOSCA.meta syntax is incorporated from an unavailable normative source",
        "sections": ["6.2"],
        "issue": (
            "Section 6.2 says the TOSCA.meta syntax is exactly that of TOSCA "
            "1.0 section 16.2, but the allowed local source does not reproduce "
            "that syntax."
        ),
        "interpretations": [
            "The local 1.3 document normatively incorporates the external syntax.",
            "Only the explicit 1.3 deltas can be cataloged without consulting another normative source.",
        ],
        "catalog_policy": "Record the incorporation requirement and all explicit 1.3 deltas; mark the inherited syntax for manual review.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-009",
        "title": "CSAR root alternatives do not state whether they are exclusive",
        "sections": ["6.1", "6.3"],
        "issue": (
            "Section 6.1 requires one of two layouts, but does not explicitly "
            "say whether an archive containing both TOSCA-Metadata and a root "
            "YAML file is invalid."
        ),
        "interpretations": [
            "“One of” is exclusive.",
            "Presence of TOSCA-Metadata takes precedence and a root YAML file is merely additional content.",
        ],
        "catalog_policy": "Record both valid alternatives and flag exclusivity as unspecified.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-010",
        "title": "Archive-without-metadata derives CSAR version from template version",
        "sections": ["6.3", "3.10.1.1"],
        "issue": (
            "Section 6.3 says CSAR-Version is defined by `template_version`, "
            "although template version identifies the template rather than the archive format."
        ),
        "interpretations": [
            "Use `template_version` literally as CSAR-Version.",
            "Treat the sentence as an editorial error and use the fixed CSAR format version 1.1.",
        ],
        "catalog_policy": "Preserve the explicit sentence and mark the semantic mismatch.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-011",
        "title": "Other-Definitions filename tokenization is incomplete",
        "sections": ["6.2"],
        "issue": (
            "Space-delimited filenames and double-quoted filenames are defined, "
            "but escaping quotes, backslashes, repeated whitespace, and empty "
            "entries is unspecified."
        ),
        "interpretations": [
            "Use a shell-like quoted token grammar.",
            "Implement only the stated blank-space delimiter and simple double quoting.",
        ],
        "catalog_policy": "Catalog only the stated delimiter and quoting rules.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-012",
        "title": "ObjectStorage `maxsize` lower bound conflicts",
        "sections": ["5.9.10.1", "5.9.10.3"],
        "issue": (
            "The property table gives `greater_or_equal: 1 GB`; the normative "
            "type definition gives `greater_or_equal: 0 GB`."
        ),
        "interpretations": [
            "Use the table constraint of 1 GB.",
            "Use the definition constraint of 0 GB.",
        ],
        "catalog_policy": "Emit both constraints and flag the conflict.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-013",
        "title": "BlockStorage `size` has conditional requiredness and `volume_id` precedence",
        "sections": ["5.9.11.1", "5.9.11.3", "5.9.11.4"],
        "issue": (
            "The table marks `size` as “yes *”, then makes it conditional on "
            "`volume_id`; the inherited storage definition also supplies a default."
        ),
        "interpretations": [
            "Require `size` only when `volume_id` is absent.",
            "Allow the inherited default to satisfy the condition.",
        ],
        "catalog_policy": "Record conditional requiredness, precedence, and inherited default separately.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-014",
        "title": "Entity Type Schema derived_from constraint is malformed",
        "sections": ["3.7.1.1"],
        "issue": (
            "The `derived_from` constraints cell states “None is the only "
            "allowed value” while the description says it contains an optional parent type name."
        ),
        "interpretations": [
            "“None” means there is no additional constraint.",
            "Only a root entity may omit or null the parent.",
        ],
        "catalog_policy": "Record the cell and prose independently and flag the contradiction.",
        "manual_review": True,
    },
    {
        "id": "TOSCA13-AMB-015",
        "title": "Examples are non-normative but some normative rules depend on example-only detail",
        "sections": ["1.6.1", "3", "4", "5", "6"],
        "issue": (
            "The specification declares Example sections non-normative, while "
            "some prose and grammar descriptions refer to examples for concrete "
            "argument forms or value interpretation."
        ),
        "interpretations": [
            "Exclude every rule that exists only in an Example section.",
            "Use examples only to illustrate a rule independently stated in normative text.",
        ],
        "catalog_policy": "Example-scoped content is excluded and never used as the sole source for a requirement.",
        "manual_review": False,
    },
    {
        "id": "TOSCA13-AMB-016",
        "title": "Lowercase modal language has uncertain normative force",
        "sections": ["3", "4", "5", "6", "7", "8", "12", "13"],
        "issue": (
            "The document mixes uppercase RFC 2119 terms with lowercase words "
            "such as “should”, “may”, and “required” in grammar descriptions."
        ),
        "interpretations": [
            "Only uppercase terms carry RFC 2119 strength.",
            "Grammar tables and unambiguous declarative prose remain binding even without uppercase terms.",
        ],
        "catalog_policy": (
            "Uppercase terms retain their exact strength; grammar markers use "
            "GRAMMAR/REQUIRED/OPTIONAL, unambiguous required/optional declarations "
            "use REQUIRED/OPTIONAL, and other declarative prose uses SEMANTIC."
        ),
        "manual_review": False,
    },
]


def normalize_text(value: str) -> str:
    value = html.unescape(value)
    value = value.replace("\xa0", " ").replace("\u00ad", "")
    value = value.replace("\u2011", "-").replace("\u2013", "-")
    value = value.replace("\ufffd", "")
    return SPACE_RE.sub(" ", value).strip()


def normalize_cell_text(value: str) -> str:
    lines: list[str] = []
    for raw_line in value.split("\u241e"):
        raw_line = html.unescape(raw_line).lstrip(" \r\n\t")
        if not raw_line:
            continue
        non_breaking_spaces = len(raw_line) - len(raw_line.lstrip("\xa0"))
        raw_line = raw_line[non_breaking_spaces:]
        text = normalize_text(raw_line)
        if text:
            indent = " " * (non_breaking_spaces + 1) if non_breaking_spaces else ""
            lines.append(f"{indent}{text}")
    return "\n".join(lines)


@dataclasses.dataclass
class Heading:
    level: int
    section: str
    title: str
    text: str
    line: int
    example_scope: bool = False


@dataclasses.dataclass
class Block:
    kind: str
    text: str
    line: int
    section: str = ""
    section_title: str = ""
    example_scope: bool = False


@dataclasses.dataclass
class Table:
    index: int
    line: int
    section: str
    section_title: str
    example_scope: bool
    rows: list[list[str]]


class SpecificationParser(HTMLParser):
    BLOCK_TAGS = {"p", "li"}

    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.headings: list[Heading] = []
        self.blocks: list[Block] = []
        self.tables: list[Table] = []
        self._capture_tag: str | None = None
        self._capture_depth = 0
        self._capture_line = 0
        self._capture_chunks: list[str] = []
        self._heading_level = 0
        self._heading_stack: dict[int, Heading] = {}
        self._current_heading: Heading | None = None
        self._ignored_depth = 0
        self._table_depth = 0
        self._table_line = 0
        self._table_rows: list[list[str]] = []
        self._row: list[str] | None = None
        self._cell_depth = 0
        self._cell_chunks: list[str] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        tag = tag.lower()
        if tag in {"style", "script"}:
            self._ignored_depth += 1
            return
        if self._ignored_depth:
            return

        if tag == "table":
            self._table_depth += 1
            if self._table_depth == 1:
                self._table_line = self.getpos()[0]
                self._table_rows = []
            return

        if self._table_depth:
            if tag == "tr" and self._table_depth == 1:
                self._row = []
            elif tag in {"td", "th"} and self._row is not None:
                self._cell_depth += 1
                if self._cell_depth == 1:
                    self._cell_chunks = []
            elif tag == "p" and self._cell_depth:
                self._cell_chunks.append("\u241e")
            elif tag == "br" and self._cell_depth:
                self._cell_chunks.append("\n")
            return

        heading_match = re.fullmatch(r"h([1-6])", tag)
        if heading_match:
            self._begin_capture(tag, int(heading_match.group(1)))
        elif tag in self.BLOCK_TAGS:
            self._begin_capture(tag, 0)
        elif tag == "br" and self._capture_tag:
            self._capture_chunks.append("\n")

    def handle_endtag(self, tag: str) -> None:
        tag = tag.lower()
        if tag in {"style", "script"} and self._ignored_depth:
            self._ignored_depth -= 1
            return
        if self._ignored_depth:
            return

        if self._table_depth:
            if tag in {"td", "th"} and self._cell_depth:
                self._cell_depth -= 1
                if self._cell_depth == 0 and self._row is not None:
                    self._row.append(normalize_cell_text("".join(self._cell_chunks)))
            elif tag == "tr" and self._row is not None:
                if any(self._row):
                    self._table_rows.append(self._row)
                self._row = None
            elif tag == "table":
                self._table_depth -= 1
                if self._table_depth == 0:
                    heading = self._current_heading
                    self.tables.append(
                        Table(
                            index=len(self.tables) + 1,
                            line=self._table_line,
                            section=heading.section if heading else "",
                            section_title=heading.title if heading else "",
                            example_scope=heading.example_scope if heading else False,
                            rows=self._table_rows,
                        )
                    )
            return

        if self._capture_tag:
            if tag == self._capture_tag and self._capture_depth == 1:
                self._finish_capture()
            elif tag == self._capture_tag:
                self._capture_depth -= 1

    def handle_data(self, data: str) -> None:
        if self._ignored_depth:
            return
        if self._cell_depth:
            self._cell_chunks.append(data)
        elif self._capture_tag:
            self._capture_chunks.append(data)

    def _begin_capture(self, tag: str, heading_level: int) -> None:
        if self._capture_tag is None:
            self._capture_tag = tag
            self._capture_depth = 1
            self._capture_line = self.getpos()[0]
            self._capture_chunks = []
            self._heading_level = heading_level
        elif self._capture_tag == tag:
            self._capture_depth += 1

    def _finish_capture(self) -> None:
        text = normalize_text("".join(self._capture_chunks))
        tag = self._capture_tag or ""
        line = self._capture_line
        heading_level = self._heading_level
        self._capture_tag = None
        self._capture_depth = 0
        self._capture_chunks = []
        self._heading_level = 0
        if not text:
            return

        if heading_level:
            match = SECTION_RE.match(text)
            section = match.group(1) if match else ""
            title = normalize_text(match.group(2)) if match else text
            for level in tuple(self._heading_stack):
                if level >= heading_level:
                    del self._heading_stack[level]
            parent_examples = any(h.example_scope for h in self._heading_stack.values())
            own_example = bool(re.search(r"\bexamples?\b", title, re.IGNORECASE))
            heading = Heading(
                level=heading_level,
                section=section,
                title=title,
                text=text,
                line=line,
                example_scope=parent_examples or own_example,
            )
            self._heading_stack[heading_level] = heading
            self._current_heading = heading
            self.headings.append(heading)
            return

        heading = self._current_heading
        self.blocks.append(
            Block(
                kind=tag,
                text=text,
                line=line,
                section=heading.section if heading else "",
                section_title=heading.title if heading else "",
                example_scope=heading.example_scope if heading else False,
            )
        )


def extract_nested_tables(
    source_text: str, headings: list[Heading]
) -> list[Table]:
    tag_re = re.compile(r"</?table\b[^>]*>", re.IGNORECASE)
    stack: list[tuple[int, int]] = []
    fragments: list[tuple[int, str]] = []
    for match in tag_re.finditer(source_text):
        if match.group(0).startswith("</"):
            if not stack:
                continue
            start, depth = stack.pop()
            if depth > 0:
                line = source_text.count("\n", 0, start) + 1
                fragments.append((line, source_text[start : match.end()]))
        else:
            stack.append((match.start(), len(stack)))

    nested: list[Table] = []
    for line, fragment in fragments:
        paragraphs = re.findall(
            r"<p\b[^>]*>(.*?)</p>", fragment, flags=re.IGNORECASE | re.DOTALL
        )
        values = [
            normalize_text(re.sub(r"<[^>]+>", "", paragraph))
            for paragraph in paragraphs
        ]
        values = [value for value in values if value]
        if not values:
            continue
        schema = "\n".join(
            value if index == 0 else f"  {value}"
            for index, value in enumerate(values)
        )
        heading = max(
            (candidate for candidate in headings if candidate.line <= line),
            key=lambda candidate: candidate.line,
        )
        nested.append(
            Table(
                index=0,
                line=line,
                section=heading.section,
                section_title=heading.title,
                example_scope=heading.example_scope,
                rows=[[schema]],
            )
        )
    return nested


def parse_specification() -> SpecificationParser:
    source_text = SOURCE.read_text(encoding="windows-1252")
    parser = SpecificationParser()
    parser.feed(source_text)
    parser.close()
    coalesced: list[Block] = []
    index = 0
    while index < len(parser.blocks):
        block = parser.blocks[index]
        if index + 1 < len(parser.blocks):
            following = parser.blocks[index + 1]
            balance = block.text.count("(") - block.text.count(")")
            following_balance = following.text.count(")") - following.text.count("(")
            if (
                balance > 0
                and following_balance >= balance
                and block.section == following.section
                and block.kind == following.kind
            ):
                coalesced.append(
                    dataclasses.replace(
                        block,
                        text=normalize_text(f"{block.text} {following.text}"),
                    )
                )
                index += 2
                continue
        coalesced.append(block)
        index += 1
    parser.blocks = coalesced
    parser.tables.extend(extract_nested_tables(source_text, parser.headings))
    parser.tables.sort(key=lambda table: (table.line, table.index))
    for index, table in enumerate(parser.tables, 1):
        table.index = index
    return parser


def inspect(parser: SpecificationParser) -> None:
    keyword_blocks = [
        block for block in parser.blocks if NORMATIVE_RE.search(block.text)
    ]
    additional = [
        heading
        for heading in parser.headings
        if re.search(r"\bAdditional Requirements?\b", heading.title, re.IGNORECASE)
    ]
    print(f"headings={len(parser.headings)}")
    print(f"sections={sum(bool(h.section) for h in parser.headings)}")
    print(f"tables={len(parser.tables)}")
    print(f"table_rows={sum(len(t.rows) for t in parser.tables)}")
    print(f"keyword_blocks={len(keyword_blocks)}")
    print(f"keyword_blocks_in_examples={sum(b.example_scope for b in keyword_blocks)}")
    print(f"additional_requirement_sections={len(additional)}")
    print("top_level_sections:")
    for heading in parser.headings:
        if heading.level == 1:
            print(f"  {heading.section} {heading.title} @L{heading.line}")
    print("additional_requirement_sections:")
    for heading in additional:
        print(f"  {heading.section} {heading.title} @L{heading.line}")


def section_sort_key(section: str) -> tuple[int, ...]:
    if not section:
        return (9999,)
    try:
        return tuple(int(part) for part in section.split("."))
    except ValueError:
        return (9998,)


def top_level_section(section: str) -> str:
    return section.split(".", 1)[0] if section else ""


def slugify(value: str) -> str:
    value = value.lower().replace("<", "").replace(">", "")
    value = re.sub(r"[^a-z0-9._-]+", "-", value)
    return value.strip("-") or "tosca-document"


def split_statements(text: str) -> list[str]:
    text = normalize_text(text.lstrip("·•- "))
    text = re.sub(r"^(?:[oO]|\d+)\.\s+", "", text)
    text = re.sub(r"^[oO]\s+", "", text)
    if not text:
        return []
    abbreviations = {
        "i.e.": "i__e__",
        "e.g.": "e__g__",
        "etc.": "etc__",
        "vs.": "vs__",
    }
    for abbreviation, placeholder in abbreviations.items():
        text = text.replace(abbreviation, placeholder)
    initial: list[str] = []
    for line in text.splitlines():
        start = 0
        depth = 0
        for index, character in enumerate(line):
            if character in "([":
                depth += 1
            elif character in ")]" and depth:
                depth -= 1
            if character not in ".!?;" or depth:
                continue
            cursor = index + 1
            while cursor < len(line) and line[cursor].isspace():
                cursor += 1
            if (
                cursor > index + 1
                and cursor < len(line)
                and re.match(r"[A-Z0-9(<`\[]", line[cursor])
            ):
                initial.append(line[start : index + 1].strip())
                start = cursor
        tail = line[start:].strip()
        if tail:
            initial.append(tail)
    statements: list[str] = []
    compound = re.compile(
        r",?\s+and\s+(?=(?:it\s+)?(?:generates?|resolves?|implements?|"
        r"normalizes?|recognizes?|requires?|rejects?|accepts?|is|has|can|"
        r"MUST|SHALL|SHOULD|MAY)\b)"
    )
    repeated_marker = re.compile(
        r"\s+(?:and|or)\s+(?=(?:MUST|SHALL|SHOULD|MAY|REQUIRED|"
        r"RECOMMENDED|OPTIONAL)\b)"
    )
    for part in initial:
        pieces = compound.split(part)
        expanded: list[str] = []
        for piece in pieces:
            expanded.extend(repeated_marker.split(piece))
        expanded = [piece.strip(" ,;") for piece in expanded if piece.strip(" ,;")]
        if len(expanded) > 1:
            subject_match = re.match(
                r"^(\([a-z]\)\s+)?(It|it|The\s+[A-Za-z][A-Za-z -]*?)\s+(can\s+)?",
                expanded[0],
            )
            if subject_match:
                label = subject_match.group(1) or ""
                subject = subject_match.group(2)
                modal = subject_match.group(3) or ""
                propagated = [expanded[0]]
                for piece in expanded[1:]:
                    if re.match(
                        r"^(?:recognize|resolve|implement|normalize|reject|accept|parse)\b",
                        piece,
                        re.IGNORECASE,
                    ):
                        propagated.append(f"{label}{subject} {modal}{piece}")
                    elif re.match(
                        r"^(?:generates|resolves|implements|normalizes|recognizes|"
                        r"requires|rejects|accepts|is|has|can)\b",
                        piece,
                        re.IGNORECASE,
                    ):
                        propagated.append(f"{label}{subject} {piece}")
                    else:
                        propagated.append(piece)
                expanded = propagated
        if (
            len(expanded) >= 2
            and re.search(r"\bcan parse$", expanded[0], re.IGNORECASE)
            and (match := re.search(r"\brecognize\s+(.+)$", expanded[1], re.IGNORECASE))
        ):
            expanded[0] = f"{expanded[0]} {match.group(1)}"
        statements.extend(expanded)
    restored: list[str] = []
    for statement in statements:
        for abbreviation, placeholder in abbreviations.items():
            statement = statement.replace(placeholder, abbreviation)
        restored.append(statement)
    return restored


def source_markers(text: str) -> list[str]:
    return [match.group(1).replace(" ", "_") for match in NORMATIVE_RE.finditer(text)]


def strength_for_text(text: str, *, fallback: str = "SEMANTIC") -> str:
    markers = source_markers(text)
    if markers:
        return markers[0]
    lowered = text.lower()
    if re.search(r"\bmust not\b|\bshall not\b|\bnot allowed\b|\bcannot\b", lowered):
        return "SEMANTIC"
    if re.search(r"\brequired\b|\bmandatory\b", lowered):
        return "REQUIRED"
    if re.search(r"\boptional\b", lowered):
        return "OPTIONAL"
    if re.search(r"\bdefault(?:s|ed)?\b", lowered):
        return "DEFAULT"
    if re.search(r"\berror\b|\binvalid\b", lowered):
        return "ERROR"
    return fallback


def section_classification(heading: Heading) -> tuple[str, str]:
    section = heading.section
    top = top_level_section(section)
    if heading.line >= APPENDIX_START_LINE:
        return "excluded", "appendix outside the numbered specification body"
    if heading.example_scope:
        return "excluded", "explicit Example section or child of one"
    if top in EXCLUDED_TOP_LEVEL_PREFIXES:
        reasons = {
            "2": "top-level section is TOSCA by example",
            "9": "section is explicitly non-normative",
            "10": "component modeling use cases",
            "11": "application modeling use cases",
        }
        return "excluded", reasons[top]
    if section.startswith("12.5"):
        return "excluded", "policy use cases"
    if section.startswith("13.4"):
        return "excluded", "discussion of examples"
    if top in RUNTIME_ONLY_PREFIXES:
        return "included-orchestrator-only", "normative/runtime material kept outside processor conformance"
    if top in {"3", "4", "5", "6", "14"}:
        return "included", "normative baseline scope"
    if top == "1":
        return "reviewed-informative", "introductory or convention material"
    return "reviewed-informative", "descriptive material outside conformance target"


def targets_for(section: str, text: str = "", *, grammar: bool = False) -> list[str]:
    top = top_level_section(section)
    lowered = text.lower()
    if section.startswith("14.2"):
        return ["service-template"]
    if section.startswith("14.3"):
        return ["processor"]
    if section.startswith("14.4"):
        return ["orchestrator"]
    if section.startswith("14.5"):
        return ["generator"]
    if section.startswith("14.6") or top == "6":
        return ["archive", "csar"]
    if top in RUNTIME_ONLY_PREFIXES:
        return ["orchestrator"]
    if "tosca orchestrator" in lowered or re.search(r"\borchestrator\b", lowered):
        if grammar and top in {"3", "4", "5"}:
            return ["service-template", "processor"]
        return ["orchestrator"]
    if top == "4":
        return ["service-template", "processor"]
    if top in {"3", "5"}:
        return ["service-template", "processor"]
    if section.startswith("1.4"):
        if "archive" in lowered:
            return ["archive", "csar"]
        if "orchestrator" in lowered:
            return ["orchestrator"]
        if "generator" in lowered:
            return ["generator"]
        if "processor" in lowered:
            return ["processor"]
        return ["service-template"]
    return ["service-template", "processor"]


def category_for(section: str, title: str, text: str, strength: str) -> str:
    haystack = f"{title} {text}".lower()
    if section.startswith("14"):
        if "error" in haystack or "fail to conform" in haystack:
            return "error"
        if "import" in haystack:
            return "import"
        if "normaliz" in haystack:
            return "normalization"
        return "conformance"
    if top_level_section(section) == "6":
        return "csar"
    if strength == "DEFAULT":
        return "default"
    if strength == "ERROR" or "error" in haystack or "invalid" in haystack:
        return "error"
    if top_level_section(section) == "4" or "function" in haystack:
        return "function"
    if "namespace" in haystack:
        return "namespace"
    if "import" in haystack:
        return "import"
    if "inherit" in haystack:
        return "inheritance"
    if "refin" in haystack:
        return "refinement"
    if "hierarch" in haystack or "derived_from" in haystack or "derives from" in haystack:
        return "hierarchy"
    if "constraint" in haystack or "valid_values" in haystack:
        return "constraint"
    if "assignment" in haystack or "mapping" in haystack:
        return "assignment"
    if top_level_section(section) in {"5", "8"} and "tosca." in haystack:
        return "normative-type"
    if strength in {"GRAMMAR", "REQUIRED", "OPTIONAL"}:
        return "grammar"
    return "semantic"


def validation_kind_for(
    section: str,
    title: str,
    text: str,
    targets: list[str],
    category: str,
    *,
    type_record: bool = False,
) -> str:
    haystack = f"{title} {text}".lower()
    if targets == ["orchestrator"] or targets == ["generator"]:
        return "orchestrator-only"
    if top_level_section(section) == "6" or "csar" in targets:
        return "csar"
    if "normaliz" in haystack:
        return "normalization"
    if category == "namespace":
        return "namespace"
    if category == "import":
        return "import"
    if category == "inheritance":
        return "inheritance"
    if category == "refinement":
        return "refinement"
    if category == "hierarchy":
        return "hierarchy"
    if category == "assignment":
        return "assignment"
    if category == "function":
        return "function"
    if category == "constraint":
        return "constraint"
    if type_record:
        return "yaml-type"
    if re.search(r"\b(name|identifier|format|character|string representation)\b", haystack):
        return "lexical"
    if re.search(r"\b(resolve|reference|lookup|find|target type|source type)\b", haystack):
        return "reference-resolution"
    if category in {"semantic", "normative-type", "conformance"}:
        return "semantic"
    return "structural"


def expected_behavior_for(
    strength: str,
    validation_kind: str,
    category: str,
    targets: list[str],
    text: str,
) -> str:
    lowered = text.lower()
    if targets == ["orchestrator"] or targets == ["generator"]:
        return "not-applicable"
    if validation_kind == "normalization":
        return "normalize"
    if validation_kind in {"namespace", "import", "hierarchy", "inheritance", "refinement", "reference-resolution"}:
        if strength in {"MUST_NOT", "SHALL_NOT", "ERROR"} or "error" in lowered or "invalid" in lowered:
            return "reject"
        return "resolve"
    if validation_kind == "function":
        if strength in {"MUST_NOT", "SHALL_NOT", "ERROR"} or "invalid" in lowered:
            return "reject"
        return "evaluate"
    if strength == "DEFAULT":
        return "resolve"
    if strength in {"MUST_NOT", "SHALL_NOT", "ERROR"}:
        return "reject"
    if strength == "REQUIRED" and category == "grammar":
        return "reject"
    if "shall not" in lowered or "must not" in lowered or "cannot" in lowered or "not allowed" in lowered:
        return "reject"
    if category == "error":
        return "reject"
    if strength in {"MAY", "OPTIONAL", "SHOULD", "SHOULD_NOT", "RECOMMENDED"}:
        return "accept"
    return "accept"


def test_obligations_for(
    strength: str,
    validation_kind: str,
    expected_behavior: str,
    targets: list[str],
) -> dict[str, bool]:
    if targets == ["orchestrator"] or targets == ["generator"]:
        return {key: False for key in TEST_OBLIGATION_KEYS}
    positive = expected_behavior in {"accept", "resolve", "normalize", "evaluate", "preserve"}
    negative = expected_behavior == "reject" or strength in {
        "MUST_NOT",
        "SHALL_NOT",
        "ERROR",
    }
    return {
        "positive": positive,
        "negative": negative,
        "boundary": validation_kind in {"lexical", "yaml-type", "constraint"},
        "inheritance": validation_kind in {"hierarchy", "inheritance", "refinement"},
        "imports": validation_kind in {"import", "namespace"},
        "normalization": validation_kind == "normalization",
    }


def ambiguity_refs(section: str) -> list[str]:
    refs: list[str] = []
    for ambiguity in MANUAL_AMBIGUITIES:
        for candidate in ambiguity["sections"]:
            if section == candidate or section.startswith(f"{candidate}."):
                refs.append(ambiguity["id"])
                break
    return refs


class CatalogBuilder:
    def __init__(self, parser: SpecificationParser) -> None:
        self.parser = parser
        self.records: list[dict[str, Any]] = []
        self._dedupe: set[tuple[Any, ...]] = set()
        self.heading_by_section = {
            heading.section: heading for heading in parser.headings if heading.section
        }
        self.grammar_tables: list[Table] = []
        self.covered_grammar_rows: set[tuple[int, int]] = set()
        self.additional_blocks: set[int] = set()
        self.covered_additional_blocks: set[int] = set()
        self.included_keyword_occurrences: collections.Counter[str] = collections.Counter()
        self.excluded_keyword_occurrences: collections.Counter[str] = collections.Counter()
        self.excluded_keyword_locations: list[tuple[int, str, list[str], str]] = []
        self.normative_type_sections: set[str] = set()
        self.covered_normative_type_sections: set[str] = set()

    def section_title(self, section: str) -> str:
        heading = self.heading_by_section.get(section)
        return heading.title if heading else ""

    def ancestor_heading(self, section: str, predicate: Any) -> Heading | None:
        parts = section.split(".")
        while parts:
            candidate = self.heading_by_section.get(".".join(parts))
            if candidate and predicate(candidate):
                return candidate
            parts.pop()
        return None

    def add(
        self,
        *,
        section: str,
        section_title: str,
        requirement: str,
        subject: str,
        strength: str,
        targets: list[str],
        line: int,
        origin: str,
        category: str | None = None,
        validation_kind: str | None = None,
        expected_behavior: str | None = None,
        notes: str | None = None,
        table: int | None = None,
        row: int | None = None,
        subrow: int | None = None,
        excerpt: str | None = None,
        schema: str | None = None,
        type_record: bool = False,
    ) -> None:
        requirement = normalize_text(requirement)
        if not requirement:
            return
        if requirement.endswith(":"):
            requirement = requirement[:-1] + "."
        elif requirement[-1] not in ".!?":
            requirement += "."
        targets = list(dict.fromkeys(targets))
        category = category or category_for(section, section_title, requirement, strength)
        validation_kind = validation_kind or validation_kind_for(
            section,
            section_title,
            requirement,
            targets,
            category,
            type_record=type_record,
        )
        expected_behavior = expected_behavior or expected_behavior_for(
            strength, validation_kind, category, targets, requirement
        )
        if targets == ["orchestrator"] or targets == ["generator"]:
            validation_kind = "orchestrator-only"
            expected_behavior = "not-applicable"
        refs = ambiguity_refs(section)
        if refs:
            ref_note = f"Related ambiguities: {', '.join(refs)}."
            notes = f"{notes} {ref_note}".strip() if notes else ref_note
        locator: dict[str, Any] = {
            "line": line,
            "origin": origin,
            "excerpt": normalize_text(excerpt or requirement)[:500],
        }
        markers = source_markers(excerpt or requirement)
        if markers:
            locator["normative_markers"] = markers
        if table is not None:
            locator["table"] = table
        if row is not None:
            locator["row"] = row
        if subrow is not None:
            locator["subrow"] = subrow
        if schema is not None:
            locator["schema"] = schema
        dedupe = (
            section,
            requirement,
            strength,
            tuple(targets),
            line,
            origin,
            table,
            row,
            subrow,
        )
        if dedupe in self._dedupe:
            return
        self._dedupe.add(dedupe)
        record = {
            "id": "",
            "target": targets,
            "section": section,
            "section_title": section_title,
            "category": category,
            "strength": strength,
            "normative": True,
            "requirement": requirement,
            "subject": slugify(subject),
            "validation_kind": validation_kind,
            "expected_processor_behavior": expected_behavior,
            "test_obligations": test_obligations_for(
                strength, validation_kind, expected_behavior, targets
            ),
            "notes": notes,
            "source_locator": locator,
        }
        self.records.append(record)
        if table is not None and row is not None:
            self.covered_grammar_rows.add((table, row))
        if line in self.additional_blocks:
            self.covered_additional_blocks.add(line)
        self.included_keyword_occurrences.update(markers)

    def build(self) -> list[dict[str, Any]]:
        self._collect_normative_type_headings()
        self._collect_prose()
        self._apply_manual_atomization()
        self._collect_tables()
        self._assign_ids()
        return self.records

    def _apply_manual_atomization(self) -> None:
        replacements: list[tuple[str, str, list[dict[str, str]]]] = [
            (
                "6.1",
                "The yaml file being a valid tosca definition template that MUST define a metadata section",
                [
                    {"requirement": "The root YAML file is a valid TOSCA definition template", "strength": "REQUIRED"},
                    {"requirement": "The root YAML file MUST define a metadata section", "strength": "MUST"},
                    {"requirement": "The root YAML file metadata section requires `template_name`", "strength": "REQUIRED"},
                    {"requirement": "The root YAML file metadata section requires `template_version`", "strength": "REQUIRED"},
                ],
            ),
            (
                "6.2",
                "it is only required to include block_0",
                [
                    {"requirement": "Only `block_0` is required in `TOSCA.meta`", "strength": "REQUIRED"},
                    {"requirement": "`block_0` requires the `Entry-Definitions` key", "strength": "REQUIRED"},
                    {"requirement": "`Entry-Definitions` identifies a valid TOSCA definitions YAML file", "strength": "GRAMMAR"},
                    {"requirement": "The orchestrator uses the `Entry-Definitions` target as the entry for parsing the CSAR", "strength": "SEMANTIC"},
                ],
            ),
            (
                "6.2",
                "it is not required to explicitly list TOSCA definitions files in subsequent blocks",
                [
                    {"requirement": "TOSCA definitions files other than the entry definition need not be listed in subsequent `TOSCA.meta` blocks", "strength": "OPTIONAL"},
                    {"requirement": "Additional TOSCA definitions files are discovered through imports in the entry definition and recursively imported files", "strength": "SEMANTIC", "validation_kind": "import", "expected_behavior": "resolve"},
                ],
            ),
            (
                "6.2",
                "The value of the Other-Definitions key is a list of filenames",
                [
                    {"requirement": "`Other-Definitions` contains a list of filenames relative to the CSAR root", "strength": "GRAMMAR"},
                    {"requirement": "A blank space delimits filenames in `Other-Definitions`", "strength": "GRAMMAR"},
                ],
            ),
            (
                "6.3",
                "the archive is required to contains a single YAML file at the root",
                [
                    {"requirement": "An archive without `TOSCA-Metadata` requires exactly one YAML file at the archive root", "strength": "REQUIRED"},
                    {"requirement": "An archive without `TOSCA-Metadata` may contain other templates only in subdirectories", "strength": "MAY"},
                ],
            ),
            (
                "6.3",
                "This file must be a valid TOSCA definitions YAML file",
                [
                    {
                        "requirement": "The root entry file must be a valid TOSCA definitions YAML file",
                        "strength": "REQUIRED",
                        "expected_behavior": "reject",
                    },
                    {"requirement": "The root entry file requires a metadata section", "strength": "REQUIRED"},
                    {"requirement": "The root entry file metadata requires `template_name`", "strength": "REQUIRED"},
                    {"requirement": "The root entry file metadata requires `template_version`", "strength": "REQUIRED"},
                ],
            ),
            (
                "14.2",
                "When using or referring to data types, artifact types, capability types",
                [
                    {
                        "requirement": f"When a service template uses or refers to {type_name}, it is valid according to the section 5 definition"
                        ,
                        "strength": "REQUIRED",
                    }
                    for type_name in (
                        "data types",
                        "artifact types",
                        "capability types",
                        "interface types",
                        "node types",
                        "relationship types",
                        "group types",
                        "policy types",
                    )
                ],
            ),
            (
                "14.3",
                "(b) It implements the requirements and semantics associated with the definitions and grammar",
                [
                    {"requirement": "(b) A conforming TOSCA processor implements the section 3 definitions", "strength": "REQUIRED"},
                    {"requirement": "(b) A conforming TOSCA processor implements the section 3 grammar", "strength": "REQUIRED"},
                    {"requirement": "(b) A conforming TOSCA processor implements every section 3 Additional Requirement", "strength": "REQUIRED"},
                ],
            ),
            (
                "14.3",
                "(d) It generates errors as required in error cases described in sections",
                [
                    {
                        "requirement": f"(d) A conforming TOSCA processor generates the required errors from section {section_ref}",
                        "strength": "REQUIRED",
                        "category": "error",
                        "expected_behavior": "reject",
                    }
                    for section_ref in ("3.1", "3.2", "3.6")
                ],
            ),
        ]
        for section, needle, atoms in replacements:
            matches = [
                record
                for record in self.records
                if record["section"] == section and needle in record["requirement"]
            ]
            if len(matches) != 1:
                raise ValueError(
                    f"manual atomization expected one match for {section}: {needle!r}; "
                    f"found {len(matches)}"
                )
            source_record = matches[0]
            self.records.remove(source_record)
            locator = source_record["source_locator"]
            for atom in atoms:
                self.add(
                    section=section,
                    section_title=source_record["section_title"],
                    requirement=atom["requirement"],
                    subject=source_record["subject"],
                    strength=atom["strength"],
                    targets=source_record["target"],
                    line=locator["line"],
                    origin="manual-atomic-prose",
                    category=atom.get("category", source_record["category"]),
                    validation_kind=atom.get(
                        "validation_kind", source_record["validation_kind"]
                    ),
                    expected_behavior=atom.get(
                        "expected_behavior",
                        (
                            "reject"
                            if atom["strength"]
                            in {"MUST", "MUST_NOT", "SHALL_NOT", "REQUIRED"}
                            and top_level_section(section) == "6"
                            else source_record["expected_processor_behavior"]
                        ),
                    ),
                    notes="Atomized from one compound source sentence.",
                    excerpt=locator["excerpt"],
                )

    def _collect_normative_type_headings(self) -> None:
        for heading in self.parser.headings:
            if heading.example_scope:
                continue
            if (
                (heading.section.startswith("5.") or heading.section.startswith("8.5."))
                and heading.title.startswith("tosca.")
            ):
                self.normative_type_sections.add(heading.section)
                targets = targets_for(heading.section, heading.title, grammar=True)
                self.add(
                    section=heading.section,
                    section_title=heading.title,
                    requirement=f"The normative type `{heading.title}` is declared by this specification",
                    subject=heading.title,
                    strength="SEMANTIC",
                    targets=targets,
                    line=heading.line,
                    origin="normative-type-heading",
                    category="normative-type",
                    validation_kind=(
                        "orchestrator-only" if targets == ["orchestrator"] else "semantic"
                    ),
                    expected_behavior=(
                        "not-applicable" if targets == ["orchestrator"] else "accept"
                    ),
                )
                self.covered_normative_type_sections.add(heading.section)

    def _block_is_included(self, block: Block) -> tuple[bool, str]:
        heading = self.heading_by_section.get(block.section)
        if heading is None:
            return False, "no numbered section"
        classification, reason = section_classification(heading)
        if classification == "excluded" or block.line >= APPENDIX_START_LINE:
            return False, reason
        if NORMATIVE_RE.search(block.text):
            return True, "explicit uppercase normative keyword"
        top = top_level_section(block.section)
        title = block.section_title.lower()
        if top == "14":
            return True, "conformance clause"
        if top == "6":
            return True, "CSAR normative prose"
        if top == "4":
            return bool(
                SEMANTIC_CUE_RE.search(block.text)
                or FUNCTION_SEMANTIC_CUE_RE.search(block.text)
            ), "function grammar and semantics"
        if re.search(
            r"additional requirements?|constraints?|requirements?|validation|"
            r"error|refinement|inheritance|namespace|import resolution",
            title,
            re.IGNORECASE,
        ):
            return True, "normative subsection"
        if top in {"3", "5"}:
            return bool(NORMATIVE_RE.search(block.text) or SEMANTIC_CUE_RE.search(block.text)), "semantic cue"
        if top == "7":
            return bool(NORMATIVE_RE.search(block.text) or SEMANTIC_CUE_RE.search(block.text)), "workflow semantics"
        if block.section.startswith(("8.2", "8.3", "8.4", "8.5")):
            return bool(NORMATIVE_RE.search(block.text) or SEMANTIC_CUE_RE.search(block.text)), "network semantics"
        if top in {"12", "13"}:
            return bool(NORMATIVE_RE.search(block.text)), "runtime normative keyword"
        return False, "descriptive prose"

    def _collect_prose(self) -> None:
        list_context: dict[str, tuple[int, str]] = {}
        for block in self.parser.blocks:
            if re.search(r"additional requirements?", block.section_title, re.IGNORECASE):
                self.additional_blocks.add(block.line)
            included, reason = self._block_is_included(block)
            markers = source_markers(block.text)
            if not included:
                if markers:
                    self.excluded_keyword_occurrences.update(markers)
                    self.excluded_keyword_locations.append(
                        (block.line, block.section, markers, reason)
                    )
                continue
            is_list_item = bool(
                re.match(r"^\s*(?:·|•|-|[oO]\s+|\d+\.\s+)", block.text)
            )
            context_for_item = list_context.get(block.section)
            for statement in split_statements(block.text):
                source_statement = statement
                markers = source_markers(statement)
                if block.section == "5" and re.match(
                    r"^Except for the examples, this section is normative and "
                    r"contains normative type definitions",
                    statement,
                    re.IGNORECASE,
                ):
                    statement = (
                        "All normative type definitions in section 5 are "
                        "required for conformance"
                    )
                if block.section == "5" and (
                    statement.startswith(
                        "The declarative approach is heavily dependent"
                    )
                    or "definition of these types must be very clear" in statement
                ):
                    continue
                if not markers and re.match(
                    r"^(?:for example|e\.g\.|future versions?\b)",
                    statement,
                    re.IGNORECASE,
                ):
                    continue
                if not markers and re.search(
                    r"\bfuture versions? of (?:this|the) specification\b",
                    statement,
                    re.IGNORECASE,
                ):
                    continue
                if not markers and re.match(
                    r"^Except for the examples, this section is normative\b",
                    statement,
                    re.IGNORECASE,
                ):
                    continue
                if (
                    not markers
                    and top_level_section(block.section) != "14"
                    and statement.endswith(":")
                    and re.search(
                        r"\b(?:statements below|these include|following additional "
                        r"requirements apply)\b",
                        statement,
                        re.IGNORECASE,
                    )
                ):
                    continue
                if not markers and not (
                    block.section.startswith(("6", "14"))
                    or (
                        top_level_section(block.section) == "4"
                        and (
                            SEMANTIC_CUE_RE.search(statement)
                            or FUNCTION_SEMANTIC_CUE_RE.search(statement)
                        )
                    )
                    or re.search(
                        r"additional requirements?|constraints?|requirements?|"
                        r"validation|error|refinement|inheritance",
                        block.section_title,
                        re.IGNORECASE,
                    )
                    or SEMANTIC_CUE_RE.search(statement)
                ):
                    continue
                original_statement = source_statement
                context_note: str | None = None
                forced_strength: str | None = None
                if is_list_item and not markers:
                    required_item = re.fullmatch(
                        r"(.+?)\s*\((required|optional)\)",
                        statement,
                        re.IGNORECASE,
                    )
                    if required_item:
                        statement = (
                            f"`{required_item.group(1).strip()}` is "
                            f"{required_item.group(2).lower()} in the preceding list context"
                        )
                        forced_strength = required_item.group(2).upper()
                    elif context_for_item:
                        context_line, context = context_for_item
                        statement = (
                            f"The preceding rule `{context.rstrip(':')}` includes "
                            f"`{statement}`"
                        )
                        context_note = (
                            f"List context is the source block at line {context_line}."
                        )
                        context_markers = source_markers(context)
                        if context_markers:
                            forced_strength = context_markers[0]
                        elif re.search(
                            r"\b(?:is|are|represents?|defines?|contains?|allows?|"
                            r"provides?|supports?|requires?|can|cannot|must|shall|"
                            r"should|may|will)\b",
                            original_statement,
                            re.IGNORECASE,
                        ):
                            forced_strength = strength_for_text(original_statement)
                        else:
                            forced_strength = "SEMANTIC"
                strength = forced_strength or strength_for_text(statement)
                if block.section == "1.6":
                    strength = "SEMANTIC"
                if block.section.startswith(("14.2", "14.3", "14.4", "14.6")) and statement.startswith("("):
                    strength = "REQUIRED"
                if block.section.startswith("14.5") and statement.startswith("("):
                    strength = "SEMANTIC"
                targets = targets_for(block.section, statement)
                category = category_for(
                    block.section, block.section_title, statement, strength
                )
                validation_kind = validation_kind_for(
                    block.section,
                    block.section_title,
                    statement,
                    targets,
                    category,
                )
                self.add(
                    section=block.section,
                    section_title=block.section_title,
                    requirement=statement,
                    subject=block.section_title,
                    strength=strength,
                    targets=targets,
                    line=block.line,
                    origin="normative-prose" if markers else "semantic-prose",
                    category=category,
                    validation_kind=validation_kind,
                    notes=context_note,
                    excerpt=original_statement,
                )
            if normalize_text(block.text).endswith(":"):
                list_context[block.section] = (
                    block.line,
                    normalize_text(block.text.lstrip("·•- ")),
                )
            elif not is_list_item:
                list_context.pop(block.section, None)

    @staticmethod
    def _headers(table: Table) -> list[str]:
        if not table.rows:
            return []
        return [normalize_text(cell).lower() for cell in table.rows[0]]

    def _is_grammar_table(self, table: Table) -> bool:
        if not table.rows:
            return False
        heading = self.heading_by_section.get(table.section)
        if heading is None:
            return False
        classification, _ = section_classification(heading)
        if classification == "excluded" or table.example_scope:
            return False
        top = top_level_section(table.section)
        if top not in {"3", "4", "5", "6", "8"}:
            return False
        if top == "8" and not table.section.startswith("8.5"):
            return False
        headers = self._headers(table)
        if len(headers) > 1:
            interesting = {
                "keyname",
                "name",
                "parameter",
                "required",
                "type",
                "constraints",
                "keyword",
                "unit",
                "shorthand name",
                "node state",
                "reference tag",
                "valid contexts",
                "directive",
                "network name",
                "namespace alias",
                "namespace prefix",
                "valid aliases",
                "shorthand names",
                "alias value",
                "value",
            }
            return bool(interesting.intersection(headers))
        title = table.section_title.lower()
        if top == "3" and (
            "_mapping" in title
            or re.search(
                r"\b(keyname|states?|aliases|values|clause|mapping|"
                r"referenced yaml types|concrete types)\b",
                title,
                re.IGNORECASE,
            )
        ):
            return True
        if top in {"5", "8"} and (
            "definition" in title
            or self.ancestor_heading(
                table.section, lambda candidate: candidate.title.startswith("tosca.")
            )
        ):
            return True
        return bool(
            re.search(
                r"grammar|notation|syntax|definition",
                title,
                re.IGNORECASE,
            )
        )

    def _collect_tables(self) -> None:
        for table in self.parser.tables:
            if not self._table_is_included(table):
                continue
            grammar = self._is_grammar_table(table)
            if grammar:
                self.grammar_tables.append(table)
                if table.rows:
                    if len(table.rows[0]) == 1:
                        self._collect_one_cell_table(table)
                    else:
                        self._collect_structured_table(table)
            self._collect_table_normative_markers(table)

    def _table_is_included(self, table: Table) -> bool:
        heading = self.heading_by_section.get(table.section)
        if heading is None or table.example_scope or table.line >= APPENDIX_START_LINE:
            return False
        classification, _ = section_classification(heading)
        return classification not in {"excluded", "reviewed-informative"}

    def _collect_table_normative_markers(self, table: Table) -> None:
        existing = {
            (
                record["source_locator"].get("row"),
                record["source_locator"].get("excerpt"),
            )
            for record in self.records
            if record["source_locator"].get("table") == table.index
        }
        for row_number, row in enumerate(table.rows, 1):
            for cell in row:
                for statement in split_statements(cell):
                    if not source_markers(statement):
                        continue
                    if (row_number, normalize_text(statement)[:500]) in existing:
                        continue
                    self.add(
                        section=table.section,
                        section_title=table.section_title,
                        requirement=statement,
                        subject=table.section_title,
                        strength=strength_for_text(statement),
                        targets=targets_for(table.section, statement, grammar=True),
                        line=table.line,
                        origin="normative-table-prose",
                        table=table.index,
                        row=row_number,
                        excerpt=statement,
                    )

    def _collect_one_cell_table(self, table: Table) -> None:
        if len(table.rows) > 1 and all(
            len(row) == 1 and "\n" not in row[0] and ":" not in row[0]
            for row in table.rows
        ):
            for row_number, row in enumerate(table.rows, 1):
                value = normalize_text(row[0])
                if not value:
                    continue
                self.add(
                    section=table.section,
                    section_title=table.section_title,
                    requirement=(
                        f"`{value}` is a permitted value in "
                        f"`{table.section_title}`"
                    ),
                    subject=table.section_title,
                    strength="GRAMMAR",
                    targets=targets_for(table.section, value, grammar=True),
                    line=table.line,
                    origin="grammar-table-value",
                    category="grammar",
                    validation_kind="lexical",
                    expected_behavior="accept",
                    table=table.index,
                    row=row_number,
                    excerpt=value,
                )
            return
        schema = "\n".join(
            cell for row in table.rows for cell in row if normalize_text(cell)
        ).strip()
        if not schema:
            return
        type_heading = self.ancestor_heading(
            table.section, lambda candidate: candidate.title.startswith("tosca.")
        )
        if type_heading and top_level_section(table.section) in {"5", "8"}:
            self._collect_type_schema(table, schema, type_heading)
            return
        self.add(
            section=table.section,
            section_title=table.section_title,
            requirement=(
                f"The `{table.section_title}` form conforms to the schema "
                f"recorded in source table {table.index}"
            ),
            subject=table.section_title,
            strength="GRAMMAR",
            targets=targets_for(table.section, schema, grammar=True),
            line=table.line,
            origin="grammar-schema",
            category="grammar",
            validation_kind="structural",
            expected_behavior="reject",
            table=table.index,
            row=1,
            excerpt=schema[:500],
            schema=schema,
        )

    def _collect_type_schema(
        self, table: Table, schema: str, type_heading: Heading
    ) -> None:
        stack: list[tuple[int, str]] = []
        previous_container: tuple[int, str] | None = None
        for line_number, raw_line in enumerate(schema.splitlines(), 1):
            if not normalize_text(raw_line) or normalize_text(raw_line).startswith("#"):
                continue
            indent = len(raw_line) - len(raw_line.lstrip(" "))
            if indent % 2:
                indent += 1
            content = normalize_text(raw_line)
            list_item = content.startswith("- ")
            if list_item:
                content = content[2:].strip()
            if ":" not in content:
                if content not in {"...", "N/A", "None"}:
                    path = ".".join(item[1] for item in stack)
                    self.add(
                        section=table.section,
                        section_title=table.section_title,
                        requirement=f"The normative type schema includes `{content}` at `{path or type_heading.title}`",
                        subject=path or type_heading.title,
                        strength="GRAMMAR",
                        targets=targets_for(table.section, content, grammar=True),
                        line=table.line,
                        origin="normative-type-schema",
                        category="normative-type",
                        table=table.index,
                        row=1,
                        subrow=line_number,
                        excerpt=raw_line,
                    )
                continue
            key, value = (part.strip() for part in content.split(":", 1))
            while stack and stack[-1][0] >= indent:
                stack.pop()
            if previous_container and previous_container[0] == indent and key in {
                "type",
                "required",
                "default",
                "description",
                "entry_schema",
                "constraints",
                "capability",
                "node",
                "relationship",
                "occurrences",
                "valid_source_types",
            }:
                stack.append(previous_container)
            path_parts = [item[1] for item in stack]
            current_path = ".".join(path_parts + [key])
            if not value:
                stack.append((indent, key))
                previous_container = (indent, key)
                if key not in {
                    "properties",
                    "attributes",
                    "capabilities",
                    "requirements",
                    "interfaces",
                    "artifacts",
                    "constraints",
                    "entry_schema",
                }:
                    self.add(
                        section=table.section,
                        section_title=table.section_title,
                        requirement=f"The normative type schema declares `{current_path}`",
                        subject=current_path,
                        strength="GRAMMAR",
                        targets=targets_for(table.section, current_path, grammar=True),
                        line=table.line,
                        origin="normative-type-schema",
                        category="normative-type",
                        table=table.index,
                        row=1,
                        subrow=line_number,
                        excerpt=raw_line,
                    )
                continue
            previous_container = None
            strength = "GRAMMAR"
            requirement = f"`{current_path}` has the schema value `{value}`"
            category = "normative-type"
            validation_kind = "structural"
            expected = None
            type_record = False
            if key == "derived_from":
                strength = "SEMANTIC"
                requirement = f"`{type_heading.title}` derives from `{value}`"
                category = "hierarchy"
                validation_kind = "hierarchy"
                expected = "resolve"
            elif key == "type":
                requirement = f"`{'.'.join(path_parts)}` has type `{value}`"
                validation_kind = "yaml-type"
                expected = "reject"
                type_record = True
            elif key == "required":
                is_required = value.lower() == "true"
                strength = "REQUIRED" if is_required else "OPTIONAL"
                requirement = (
                    f"`{'.'.join(path_parts)}` is "
                    f"{'required' if is_required else 'optional'}"
                )
                expected = "reject" if is_required else "accept"
            elif key == "default":
                strength = "DEFAULT"
                requirement = f"`{'.'.join(path_parts)}` defaults to `{value}`"
                expected = "resolve"
            elif key in {
                "valid_target_types",
                "valid_source_types",
                "occurrences",
                "file_ext",
                "mime_type",
                "constraints",
            } or (stack and stack[-1][1] == "constraints"):
                requirement = f"`{current_path}` is constrained by `{value}`"
                category = "constraint"
                validation_kind = "constraint"
                expected = "reject"
            elif key == "description":
                strength = "SEMANTIC"
                requirement = f"`{'.'.join(path_parts)}` has the specified semantics: {value}"
            self.add(
                section=table.section,
                section_title=table.section_title,
                requirement=requirement,
                subject=current_path,
                strength=strength,
                targets=targets_for(table.section, requirement, grammar=True),
                line=table.line,
                origin="normative-type-schema",
                category=category,
                validation_kind=validation_kind,
                expected_behavior=expected,
                table=table.index,
                row=1,
                subrow=line_number,
                excerpt=raw_line,
                type_record=type_record,
            )

    def _collect_structured_table(self, table: Table) -> None:
        headers = self._headers(table)
        for row_number, row in enumerate(table.rows[1:], 2):
            cells = list(row) + [""] * max(0, len(headers) - len(row))
            name = normalize_text(cells[0])
            if not name:
                continue
            self._collect_structured_row(table, row_number, headers, cells, name)

    def _collect_structured_row(
        self,
        table: Table,
        row_number: int,
        headers: list[str],
        cells: list[str],
        name: str,
    ) -> None:
        targets = targets_for(table.section, name, grammar=True)
        if name.lower() in {"n/a", "none"} and all(
            normalize_text(cell).lower() in {"", "n/a", "none"} for cell in cells[1:]
        ):
            self.add(
                section=table.section,
                section_title=table.section_title,
                requirement=(
                    f"The `{table.section_title}` table declares no entries "
                    "at this type-definition level"
                ),
                subject=table.section_title,
                strength="GRAMMAR",
                targets=targets,
                line=table.line,
                origin="grammar-table-empty",
                category="grammar",
                validation_kind="structural",
                expected_behavior="accept",
                table=table.index,
                row=row_number,
                excerpt=" | ".join(cells),
            )
            return
        required_index = headers.index("required") if "required" in headers else None
        type_index = headers.index("type") if "type" in headers else None
        constraints_index = (
            headers.index("constraints") if "constraints" in headers else None
        )
        description_index = (
            headers.index("description") if "description" in headers else None
        )

        if required_index is not None:
            marker = normalize_text(cells[required_index])
            lowered = marker.lower()
            if lowered not in {"", "n/a", "none"}:
                required = lowered.startswith(("yes", "true"))
                conditional = "*" in marker or "conditional" in lowered
                strength = "REQUIRED" if required else "OPTIONAL"
                qualifier = " conditionally" if conditional else ""
                self.add(
                    section=table.section,
                    section_title=table.section_title,
                    requirement=(
                        f"The `{name}` entry is{qualifier} "
                        f"{'required' if required else 'optional'}"
                    ),
                    subject=name,
                    strength=strength,
                    targets=targets,
                    line=table.line,
                    origin="grammar-table-presence",
                    category="grammar",
                    validation_kind="structural",
                    expected_behavior="reject" if required else "accept",
                    notes=(
                        "The table uses a conditional requiredness marker."
                        if conditional
                        else None
                    ),
                    table=table.index,
                    row=row_number,
                    excerpt=f"{name} | required={marker}",
                )

        if type_index is not None:
            type_value = normalize_text(cells[type_index])
            if type_value.lower() not in {"", "n/a", "none"}:
                self.add(
                    section=table.section,
                    section_title=table.section_title,
                    requirement=f"The `{name}` entry has type/schema `{type_value}`",
                    subject=name,
                    strength="GRAMMAR",
                    targets=targets,
                    line=table.line,
                    origin="grammar-table-type",
                    category="grammar",
                    validation_kind="yaml-type",
                    expected_behavior="reject",
                    table=table.index,
                    row=row_number,
                    excerpt=f"{name} | type={type_value}",
                    type_record=True,
                )

        if constraints_index is not None:
            constraint_cell = cells[constraints_index]
            for fragment in self._constraint_fragments(constraint_cell):
                default_match = re.match(r"default\s*:\s*(.+)", fragment, re.IGNORECASE)
                if default_match:
                    self.add(
                        section=table.section,
                        section_title=table.section_title,
                        requirement=f"The `{name}` entry defaults to `{default_match.group(1).strip()}`",
                        subject=name,
                        strength="DEFAULT",
                        targets=targets,
                        line=table.line,
                        origin="grammar-table-default",
                        category="default",
                        validation_kind="semantic",
                        expected_behavior="resolve",
                        table=table.index,
                        row=row_number,
                        excerpt=fragment,
                    )
                elif fragment.lower() not in {"none", "n/a"}:
                    self.add(
                        section=table.section,
                        section_title=table.section_title,
                        requirement=f"The `{name}` entry is constrained by `{fragment}`",
                        subject=name,
                        strength="GRAMMAR",
                        targets=targets,
                        line=table.line,
                        origin="grammar-table-constraint",
                        category="constraint",
                        validation_kind="constraint",
                        expected_behavior="reject",
                        table=table.index,
                        row=row_number,
                        excerpt=fragment,
                    )

        if description_index is not None:
            description = cells[description_index]
            for statement in split_statements(description):
                if statement.lower() in {"n/a", "none"}:
                    continue
                strength = strength_for_text(statement)
                self.add(
                    section=table.section,
                    section_title=table.section_title,
                    requirement=f"For `{name}`: {statement}",
                    subject=name,
                    strength=strength,
                    targets=targets_for(table.section, statement, grammar=True),
                    line=table.line,
                    origin="grammar-table-description",
                    table=table.index,
                    row=row_number,
                    excerpt=statement,
                )

        known_indexes = {
            0,
            *(
                index
                for index in (
                    required_index,
                    type_index,
                    constraints_index,
                    description_index,
                )
                if index is not None
            ),
        }
        for index, header in enumerate(headers):
            if index in known_indexes or index >= len(cells):
                continue
            value = normalize_text(cells[index])
            if value.lower() in {"", "none", "n/a"}:
                continue
            self.add(
                section=table.section,
                section_title=table.section_title,
                requirement=f"For `{name}`, `{header}` is `{value}`",
                subject=name,
                strength="GRAMMAR",
                targets=targets,
                line=table.line,
                origin="grammar-table-field",
                category="grammar",
                table=table.index,
                row=row_number,
                excerpt=f"{name} | {header}={value}",
            )

    @staticmethod
    def _constraint_fragments(value: str) -> list[str]:
        value = value.strip()
        if not value:
            return []
        fragments: list[str] = []
        for line in value.splitlines():
            parts = re.split(
                r"\s+(?=(?:default|valid_values|equal|greater_or_equal|"
                r"greater_than|less_or_equal|less_than|in_range|min_length|"
                r"max_length|pattern|length)\s*:)",
                normalize_text(line),
                flags=re.IGNORECASE,
            )
            fragments.extend(part.strip(" ;,") for part in parts if part.strip(" ;,"))
        return fragments

    def _assign_ids(self) -> None:
        self.records.sort(
            key=lambda record: (
                section_sort_key(record["section"]),
                record["source_locator"]["line"],
                record["source_locator"].get("table", 0),
                record["source_locator"].get("row", 0),
                record["source_locator"].get("subrow", 0),
                record["category"],
                record["requirement"],
            )
        )
        counters: collections.Counter[str] = collections.Counter()
        for record in self.records:
            section = record["section"] or "UNNUMBERED"
            counters[section] += 1
            record["id"] = f"TOSCA13-{section}-{counters[section]:03d}"


def source_sha256() -> str:
    return hashlib.sha256(SOURCE.read_bytes()).hexdigest()


def build_catalog(parser: SpecificationParser) -> tuple[CatalogBuilder, dict[str, Any]]:
    builder = CatalogBuilder(parser)
    requirements = builder.build()
    catalog = {
        "specification": {
            "name": "TOSCA Simple Profile in YAML",
            "version": "1.3",
            "source": str(SOURCE.relative_to(REPO_ROOT)),
            "source_sha256": source_sha256(),
        },
        "requirements": requirements,
    }
    return builder, catalog


def catalog_stats(requirements: list[dict[str, Any]]) -> dict[str, Any]:
    counters: dict[str, collections.Counter[str]] = {
        "section": collections.Counter(),
        "category": collections.Counter(),
        "strength": collections.Counter(),
        "target": collections.Counter(),
        "validation_kind": collections.Counter(),
    }
    for requirement in requirements:
        counters["section"][requirement["section"]] += 1
        counters["category"][requirement["category"]] += 1
        counters["strength"][requirement["strength"]] += 1
        counters["validation_kind"][requirement["validation_kind"]] += 1
        counters["target"].update(requirement["target"])
    return {
        "total": len(requirements),
        "mandatory": sum(
            counters["strength"][strength]
            for strength in ("MUST", "SHALL", "REQUIRED")
        ),
        "prohibitions": sum(
            counters["strength"][strength] for strength in ("MUST_NOT", "SHALL_NOT")
        ),
        "grammar": counters["category"]["grammar"],
        "grammar_strength": counters["strength"]["GRAMMAR"],
        "errors": sum(
            1
            for requirement in requirements
            if requirement["strength"] == "ERROR"
            or requirement["category"] == "error"
        ),
        "ambiguities": len(MANUAL_AMBIGUITIES),
        "processor": counters["target"]["processor"],
        "service_template": counters["target"]["service-template"],
        "archive_csar": sum(
            1
            for requirement in requirements
            if {"archive", "csar"}.intersection(requirement["target"])
        ),
        "counters": counters,
    }


def markdown_cell(value: Any) -> str:
    return str(value).replace("|", "\\|").replace("\n", " ")


def counter_table(counter: collections.Counter[str], *, section: bool = False) -> str:
    keys = sorted(counter, key=section_sort_key if section else lambda value: value)
    lines = ["| Value | Requirements |", "|---|---:|"]
    lines.extend(f"| {markdown_cell(key)} | {counter[key]} |" for key in keys)
    return "\n".join(lines)


def render_summary(requirements: list[dict[str, Any]]) -> str:
    stats = catalog_stats(requirements)
    counters = stats["counters"]
    return f"""# TOSCA 1.3 normative requirement catalog summary

This summary is generated from `requirements.yaml`. Counts are requirement-record
counts; a record can contribute to more than one conformance-target count.

## Headline counts

| Measure | Count |
|---|---:|
| Atomic requirements | {stats["total"]} |
| Mandatory (`MUST`, `SHALL`, `REQUIRED`) | {stats["mandatory"]} |
| Explicit prohibitions (`MUST_NOT`, `SHALL_NOT`) | {stats["prohibitions"]} |
| Grammar-category records | {stats["grammar"]} |
| `GRAMMAR`-strength records | {stats["grammar_strength"]} |
| Error requirements | {stats["errors"]} |
| Recorded ambiguities | {stats["ambiguities"]} |
| Processor target | {stats["processor"]} |
| Service-template target | {stats["service_template"]} |
| Archive/CSAR target (unique records) | {stats["archive_csar"]} |

## By section

{counter_table(counters["section"], section=True)}

## By category

{counter_table(counters["category"])}

## By strength

{counter_table(counters["strength"])}

## By conformance target

{counter_table(counters["target"])}

The `orchestrator` and `generator` rows are separate inventories. They are not
counted as processor conformance.

## By validation kind

{counter_table(counters["validation_kind"])}
"""


def render_ambiguities() -> str:
    sections = [
        "# TOSCA 1.3 specification ambiguities",
        "",
        (
            "These are specification-source issues, not claims about Puccini "
            "behavior. Catalog records preserve the conflicting or incomplete "
            "rules and link back here through their `notes` fields."
        ),
        "",
        f"Recorded ambiguities: **{len(MANUAL_AMBIGUITIES)}**.",
        "",
    ]
    for ambiguity in MANUAL_AMBIGUITIES:
        sections.extend(
            [
                f"## {ambiguity['id']}: {ambiguity['title']}",
                "",
                f"- Related sections: {', '.join(ambiguity['sections'])}",
                f"- Issue: {ambiguity['issue']}",
                "- Possible interpretations:",
                "",
                *[f"  - {item}" for item in ambiguity["interpretations"]],
                "",
                f"- Catalog treatment: {ambiguity['catalog_policy']}",
                (
                    "- Manual recheck required: yes"
                    if ambiguity["manual_review"]
                    else "- Manual recheck required: no; retained for traceability"
                ),
                "",
            ]
        )
    return "\n".join(sections)


def grammar_expected_rows(table: Table) -> set[tuple[int, int]]:
    if not table.rows:
        return set()
    if len(table.rows[0]) == 1:
        if len(table.rows) > 1 and all(
            len(row) == 1 and "\n" not in row[0] and ":" not in row[0]
            for row in table.rows
        ):
            return {
                (table.index, row_number)
                for row_number, row in enumerate(table.rows, 1)
                if normalize_text(row[0])
            }
        return {(table.index, 1)}
    return {
        (table.index, row_number)
        for row_number, row in enumerate(table.rows[1:], 2)
        if row and normalize_text(row[0])
    }


def keyword_inventory(
    parser: SpecificationParser, builder: CatalogBuilder
) -> dict[str, Any]:
    included_blocks: list[tuple[int, str, list[str]]] = []
    excluded_blocks: list[tuple[int, str, list[str], str]] = []
    included_tables: list[tuple[int, int, str, list[str]]] = []
    excluded_tables: list[tuple[int, int, str, list[str], str]] = []
    included_counts: collections.Counter[str] = collections.Counter()
    excluded_counts: collections.Counter[str] = collections.Counter()

    for block in parser.blocks:
        markers = source_markers(block.text)
        if not markers:
            continue
        included, reason = builder._block_is_included(block)
        if included:
            included_blocks.append((block.line, block.section, markers))
            included_counts.update(markers)
        else:
            excluded_blocks.append((block.line, block.section, markers, reason))
            excluded_counts.update(markers)

    for table in parser.tables:
        table_included = builder._table_is_included(table)
        heading = builder.heading_by_section.get(table.section)
        reason = (
            section_classification(heading)[1]
            if heading is not None
            else "no numbered section"
        )
        for row_number, row in enumerate(table.rows, 1):
            markers = source_markers(" ".join(row))
            if not markers:
                continue
            if table_included:
                included_tables.append(
                    (table.index, row_number, table.section, markers)
                )
                included_counts.update(markers)
            else:
                excluded_tables.append(
                    (table.index, row_number, table.section, markers, reason)
                )
                excluded_counts.update(markers)
    return {
        "included_blocks": included_blocks,
        "excluded_blocks": excluded_blocks,
        "included_tables": included_tables,
        "excluded_tables": excluded_tables,
        "included_counts": included_counts,
        "excluded_counts": excluded_counts,
    }


def coverage_results(
    parser: SpecificationParser,
    builder: CatalogBuilder,
    requirements: list[dict[str, Any]],
) -> dict[str, Any]:
    expected_rows = set().union(
        *(grammar_expected_rows(table) for table in builder.grammar_tables)
    ) if builder.grammar_tables else set()
    record_lines = {record["source_locator"]["line"] for record in requirements}
    record_table_rows = {
        (
            record["source_locator"].get("table"),
            record["source_locator"].get("row"),
        )
        for record in requirements
        if "table" in record["source_locator"]
    }
    inventory = keyword_inventory(parser, builder)
    raw_source = SOURCE.read_text(encoding="windows-1252")
    raw_marker_counts = collections.Counter(source_markers(raw_source))
    parsed_marker_counts = (
        inventory["included_counts"] + inventory["excluded_counts"]
    )
    presence_rows: set[tuple[int, int]] = set()
    default_occurrences = 0
    for table in builder.grammar_tables:
        if not table.rows:
            continue
        headers = builder._headers(table)
        if len(table.rows[0]) > 1 and "required" in headers:
            required_index = headers.index("required")
            for row_number, row in enumerate(table.rows[1:], 2):
                marker = normalize_text(
                    row[required_index] if required_index < len(row) else ""
                ).lower()
                if marker not in {"", "n/a", "none"}:
                    presence_rows.add((table.index, row_number))
        default_occurrences += sum(
            len(re.findall(r"\bdefault\s*:", cell, re.IGNORECASE))
            for row in table.rows
            for cell in row
        )
    covered_presence_rows = {
        (
            record["source_locator"].get("table"),
            record["source_locator"].get("row"),
        )
        for record in requirements
        if record["source_locator"]["origin"] == "grammar-table-presence"
    }
    default_records = [
        record
        for record in requirements
        if record["strength"] == "DEFAULT"
        and "table" in record["source_locator"]
    ]
    notation_headings = [
        heading
        for heading in parser.headings
        if heading.section
        and not heading.example_scope
        and re.search(r"\b(?:short|long)(?:-form)? notation\b", heading.title, re.IGNORECASE)
        and section_classification(heading)[0] != "excluded"
    ]
    sections_with_records = {record["section"] for record in requirements}
    error_blocks = [
        block
        for block in parser.blocks
        if builder._block_is_included(block)[0]
        and re.search(
            r"\b(?:error|invalid|not allowed|cannot|fail(?:s|ed)? to conform)\b",
            block.text,
            re.IGNORECASE,
        )
    ]
    clause_143_lines = {
        block.line
        for block in parser.blocks
        if block.section.startswith("14.3") and normalize_text(block.text)
    }
    numbered_headings = [heading for heading in parser.headings if heading.section]
    return {
        "numbered_headings": len(numbered_headings),
        "unclassified_headings": [
            heading.section
            for heading in numbered_headings
            if not section_classification(heading)[0]
        ],
        "grammar_tables": len(builder.grammar_tables),
        "grammar_rows": len(expected_rows),
        "missing_grammar_rows": sorted(expected_rows - builder.covered_grammar_rows),
        "additional_blocks": len(builder.additional_blocks),
        "missing_additional_blocks": sorted(
            builder.additional_blocks - builder.covered_additional_blocks
        ),
        "normative_type_sections": len(builder.normative_type_sections),
        "missing_normative_type_sections": sorted(
            builder.normative_type_sections - builder.covered_normative_type_sections,
            key=section_sort_key,
        ),
        "missing_keyword_blocks": [
            (line, section, markers)
            for line, section, markers in inventory["included_blocks"]
            if line not in record_lines
        ],
        "missing_keyword_table_rows": [
            (table, row, section, markers)
            for table, row, section, markers in inventory["included_tables"]
            if (table, row) not in record_table_rows
        ],
        "missing_clause_143_blocks": sorted(clause_143_lines - record_lines),
        "raw_heading_tags": len(re.findall(r"<h[1-6]\b", raw_source, re.IGNORECASE)),
        "raw_table_tags": len(re.findall(r"<table\b", raw_source, re.IGNORECASE)),
        "raw_row_tags": len(re.findall(r"<tr\b", raw_source, re.IGNORECASE)),
        "raw_marker_counts": raw_marker_counts,
        "parsed_marker_counts": parsed_marker_counts,
        "marker_count_mismatch": (
            sum(raw_marker_counts.values()) != sum(parsed_marker_counts.values())
        ),
        "presence_markers": len(presence_rows),
        "missing_presence_markers": sorted(presence_rows - covered_presence_rows),
        "default_occurrences": default_occurrences,
        "default_records": len(default_records),
        "missing_default_records": max(0, default_occurrences - len(default_records)),
        "notation_sections": len(notation_headings),
        "missing_notation_sections": sorted(
            {
                heading.section
                for heading in notation_headings
                if heading.section not in sections_with_records
            },
            key=section_sort_key,
        ),
        "error_blocks": len(error_blocks),
        "missing_error_blocks": [
            (block.line, block.section)
            for block in error_blocks
            if block.line not in record_lines
        ],
        "keyword_inventory": inventory,
    }


def validate_catalog(
    parser: SpecificationParser,
    builder: CatalogBuilder,
    catalog: dict[str, Any],
) -> dict[str, Any]:
    errors: list[str] = []
    requirements = catalog.get("requirements")
    if not isinstance(requirements, list):
        raise ValueError("requirements must be a list")
    if catalog.get("specification", {}).get("source_sha256") != source_sha256():
        errors.append("source SHA-256 does not match")

    ids: set[str] = set()
    required_fields = {
        "id",
        "target",
        "section",
        "section_title",
        "category",
        "strength",
        "normative",
        "requirement",
        "subject",
        "validation_kind",
        "expected_processor_behavior",
        "test_obligations",
        "notes",
        "source_locator",
    }
    for index, record in enumerate(requirements):
        missing = required_fields - set(record)
        if missing:
            errors.append(f"record {index} missing fields: {sorted(missing)}")
            continue
        record_id = record["id"]
        if record_id in ids:
            errors.append(f"duplicate id: {record_id}")
        ids.add(record_id)
        if not re.fullmatch(r"TOSCA13-(?:\d+(?:\.\d+)*|UNNUMBERED)-\d{3}", record_id):
            errors.append(f"invalid id: {record_id}")
        if record["strength"] not in VALID_STRENGTHS:
            errors.append(f"{record_id}: invalid strength {record['strength']}")
        if record["validation_kind"] not in VALID_VALIDATION_KINDS:
            errors.append(
                f"{record_id}: invalid validation kind {record['validation_kind']}"
            )
        if record["expected_processor_behavior"] not in VALID_EXPECTED_BEHAVIORS:
            errors.append(
                f"{record_id}: invalid processor behavior "
                f"{record['expected_processor_behavior']}"
            )
        obligations = record["test_obligations"]
        if set(obligations) != set(TEST_OBLIGATION_KEYS) or not all(
            isinstance(obligations[key], bool) for key in TEST_OBLIGATION_KEYS
        ):
            errors.append(f"{record_id}: invalid test obligations")
        if record["normative"] is not True:
            errors.append(f"{record_id}: normative must be true")

    coverage = coverage_results(parser, builder, requirements)
    for key in (
        "unclassified_headings",
        "missing_grammar_rows",
        "missing_additional_blocks",
        "missing_normative_type_sections",
        "missing_keyword_blocks",
        "missing_keyword_table_rows",
        "missing_clause_143_blocks",
        "missing_presence_markers",
        "missing_notation_sections",
        "missing_error_blocks",
    ):
        if coverage[key]:
            errors.append(f"{key}: {coverage[key][:20]}")
    if coverage["marker_count_mismatch"]:
        errors.append(
            "raw and parsed normative-marker totals differ: "
            f"{coverage['raw_marker_counts']} != {coverage['parsed_marker_counts']}"
        )
    if coverage["missing_default_records"]:
        errors.append(
            f"default records missing: {coverage['missing_default_records']}"
        )
    if errors:
        raise ValueError("catalog verification failed:\n- " + "\n- ".join(errors))
    return coverage


def render_extraction_report(
    parser: SpecificationParser,
    builder: CatalogBuilder,
    requirements: list[dict[str, Any]],
    coverage: dict[str, Any],
) -> str:
    inventory = coverage["keyword_inventory"]
    status_counts = collections.Counter(
        section_classification(heading)[0]
        for heading in parser.headings
        if heading.section
    )
    lines = [
        "# TOSCA 1.3 normative catalog extraction report",
        "",
        "## Source and method",
        "",
        f"- Source: `{SOURCE.relative_to(REPO_ROOT)}`",
        f"- SHA-256: `{source_sha256()}`",
        "- Source encoding: Windows-1252, as declared by the HTML.",
        (
            "- First pass: structural HTML parse of numbered headings, prose/list "
            "blocks, and tables; atomic records retain source line/table/row locators."
        ),
        (
            "- Second pass: an independent coverage inventory checks all numbered "
            "headings, grammar rows, Additional Requirements blocks, exact uppercase "
            "normative markers, normative type headings, and clause 14.3 blocks."
        ),
        (
            "- Three nested Word-export tables in section 5 are recovered as "
            "independent normative type-definition tables; their empty outer "
            "wrapper tables remain in the raw structural count but are not grammar."
        ),
        "- Explicit Example scopes are excluded under section 1.6.1.",
        "",
        "## Second-pass completeness checks",
        "",
        "| Check | Expected | Missing |",
        "|---|---:|---:|",
        f"| Numbered headings classified | {coverage['numbered_headings']} | {len(coverage['unclassified_headings'])} |",
        f"| Grammar tables inventoried | {coverage['grammar_tables']} | 0 |",
        f"| Grammar data rows represented | {coverage['grammar_rows']} | {len(coverage['missing_grammar_rows'])} |",
        f"| Additional Requirements prose blocks represented | {coverage['additional_blocks']} | {len(coverage['missing_additional_blocks'])} |",
        f"| Normative type headings represented | {coverage['normative_type_sections']} | {len(coverage['missing_normative_type_sections'])} |",
        f"| Included prose blocks with uppercase markers represented | {len(inventory['included_blocks'])} | {len(coverage['missing_keyword_blocks'])} |",
        f"| Included table rows with uppercase markers represented | {len(inventory['included_tables'])} | {len(coverage['missing_keyword_table_rows'])} |",
        f"| Clause 14.3 prose blocks represented | {sum(1 for block in parser.blocks if block.section.startswith('14.3') and normalize_text(block.text))} | {len(coverage['missing_clause_143_blocks'])} |",
        f"| Required/optional table markers represented | {coverage['presence_markers']} | {len(coverage['missing_presence_markers'])} |",
        f"| Explicit `default:` occurrences represented | {coverage['default_occurrences']} | {coverage['missing_default_records']} |",
        f"| Short/long notation sections represented | {coverage['notation_sections']} | {len(coverage['missing_notation_sections'])} |",
        f"| Error-bearing prose blocks represented | {coverage['error_blocks']} | {len(coverage['missing_error_blocks'])} |",
        "",
        "Raw-source checks use direct tag and keyword scans rather than the extraction traversal:",
        "",
        "| Raw-source measure | Raw count | Parsed count |",
        "|---|---:|---:|",
        f"| Heading tags | {coverage['raw_heading_tags']} | {len(parser.headings)} |",
        f"| Table tags | {coverage['raw_table_tags']} | {len(parser.tables)} |",
        f"| Table-row tags | {coverage['raw_row_tags']} | {sum(len(table.rows) for table in parser.tables)} |",
        f"| Uppercase normative markers | {sum(coverage['raw_marker_counts'].values())} | {sum(coverage['parsed_marker_counts'].values())} |",
        "",
        (
            "The 1,294 raw `<tr>` tags include seven empty rows from layout or "
            "outer-wrapper tables; the parsed count contains 1,287 non-empty rows. "
            "All 633 table tags are nevertheless inventoried."
        ),
        "",
        "## Exact uppercase normative-word inventory",
        "",
        "| Word | Included occurrences | Excluded occurrences |",
        "|---|---:|---:|",
    ]
    for word in sorted(set(inventory["included_counts"]) | set(inventory["excluded_counts"])):
        lines.append(
            f"| {word.replace('_', ' ')} | {inventory['included_counts'][word]} | "
            f"{inventory['excluded_counts'][word]} |"
        )

    lines.extend(
        [
            "",
            "Excluded marker-bearing source locations are retained here so the "
            "exclusion is reviewable:",
            "",
            "| Source | Section | Markers | Exclusion reason |",
            "|---|---|---|---|",
        ]
    )
    for line, section, markers, reason in inventory["excluded_blocks"]:
        lines.append(
            f"| line {line} | {section or 'unnumbered'} | {', '.join(markers)} | "
            f"{markdown_cell(reason)} |"
        )
    for table, row, section, markers, reason in inventory["excluded_tables"]:
        lines.append(
            f"| table {table}, row {row} | {section or 'unnumbered'} | "
            f"{', '.join(markers)} | {markdown_cell(reason)} |"
        )
    if not inventory["excluded_blocks"] and not inventory["excluded_tables"]:
        lines.append("| none | - | - | - |")

    lines.extend(
        [
            "",
            "## Additional Requirements inventory",
            "",
            "| Section | Title | Source line |",
            "|---|---|---:|",
        ]
    )
    additional_headings = [
        heading
        for heading in parser.headings
        if heading.section
        and re.search(r"\bAdditional Requirements?\b", heading.title, re.IGNORECASE)
    ]
    for heading in additional_headings:
        lines.append(
            f"| {heading.section} | {markdown_cell(heading.title)} | {heading.line} |"
        )

    lines.extend(
        [
            "",
            "## Grammar table inventory",
            "",
            "| Table | Section | Section title | Source line | Data rows | Covered rows |",
            "|---:|---|---|---:|---:|---:|",
        ]
    )
    for table in builder.grammar_tables:
        expected = grammar_expected_rows(table)
        covered = expected.intersection(builder.covered_grammar_rows)
        lines.append(
            f"| {table.index} | {table.section} | {markdown_cell(table.section_title)} | "
            f"{table.line} | {len(expected)} | {len(covered)} |"
        )

    lines.extend(
        [
            "",
            "## Review disposition for every numbered section",
            "",
            f"Classification totals: {', '.join(f'`{key}`={value}' for key, value in sorted(status_counts.items()))}.",
            "",
            "| Section | Title | Source line | Disposition | Basis |",
            "|---|---|---:|---|---|",
        ]
    )
    for heading in (h for h in parser.headings if h.section):
        status, reason = section_classification(heading)
        lines.append(
            f"| {heading.section} | {markdown_cell(heading.title)} | {heading.line} | "
            f"{status} | {markdown_cell(reason)} |"
        )

    lines.extend(
        [
            "",
            "## Conformance clause 14.3 trace",
            "",
            (
                "All non-empty prose/list blocks in 14.3 have requirement records. "
                "The catalog separately inventories section 3 definitions and "
                "additional requirements, section 4 functions, section 5 normative "
                "types, imports, error rules, and the ambiguous string-normalization "
                "citation."
            ),
            "",
            "## Project CSAR scope",
            "",
            (
                "A read-only scope check found CSAR parsing, metadata validation, "
                "entry-definition selection, and archive creation under `tosca/csar` "
                "and `executables/puccini-csar`. Therefore all explicit section 6 "
                "requirements carry the separate `archive` and `csar` targets. "
                "Project behavior was not used as normative evidence."
            ),
            "",
            "## Normative type-definition coverage",
            "",
            (
                "All 61 `tosca.*` normative type headings in section 5 are "
                "represented. Five additional network type headings from section "
                "8.5 are retained only as a separate orchestrator inventory and do "
                "not contribute to processor conformance."
            ),
            "",
            "## Sections requiring manual classification or external clarification",
            "",
            (
                "The source sections are structurally classified, but the issues "
                "listed below cannot be resolved from the permitted source alone:"
            ),
            "",
        ]
    )
    for ambiguity in MANUAL_AMBIGUITIES:
        if ambiguity["manual_review"]:
            lines.append(
                f"- {ambiguity['id']} ({', '.join(ambiguity['sections'])}): "
                f"{ambiguity['title']}."
            )
    lines.extend(
        [
            "",
            "The catalog is complete as an extraction baseline: no structural "
            "coverage check is outstanding. “Complete” does not mean that the "
            "source ambiguities above have been resolved or that Puccini conforms.",
            "",
        ]
    )
    return "\n".join(lines)


def render_outputs(
    parser: SpecificationParser,
    builder: CatalogBuilder,
    catalog: dict[str, Any],
    coverage: dict[str, Any],
) -> dict[pathlib.Path, str]:
    yaml_text = yaml.safe_dump(
        catalog,
        sort_keys=False,
        allow_unicode=True,
        width=120,
    )
    requirements = catalog["requirements"]
    return {
        REQUIREMENTS_PATH: yaml_text,
        SUMMARY_PATH: render_summary(requirements),
        AMBIGUITIES_PATH: render_ambiguities(),
        EXTRACTION_REPORT_PATH: render_extraction_report(
            parser, builder, requirements, coverage
        ),
    }


def write_outputs(outputs: dict[pathlib.Path, str]) -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    for path, content in outputs.items():
        path.write_text(content.rstrip() + "\n", encoding="utf-8")


def main() -> int:
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument(
        "command",
        nargs="?",
        choices=("inspect", "build", "verify"),
        default="inspect",
    )
    args = arg_parser.parse_args()
    parser = parse_specification()
    if args.command == "inspect":
        inspect(parser)
        return 0
    builder, catalog = build_catalog(parser)
    coverage = validate_catalog(parser, builder, catalog)
    outputs = render_outputs(parser, builder, catalog, coverage)
    if args.command == "build":
        write_outputs(outputs)
    else:
        for path, expected in outputs.items():
            if not path.exists():
                raise ValueError(f"missing generated file: {path}")
            actual = path.read_text(encoding="utf-8")
            if actual != expected.rstrip() + "\n":
                raise ValueError(f"generated file is stale: {path}")
        loaded = yaml.safe_load(REQUIREMENTS_PATH.read_text(encoding="utf-8"))
        if loaded != catalog:
            raise ValueError("requirements.yaml does not round-trip to expected data")
    stats = catalog_stats(catalog["requirements"])
    print(
        "requirements={total} mandatory={mandatory} prohibitions={prohibitions} "
        "grammar_category={grammar} grammar_strength={grammar_strength} "
        "errors={errors} ambiguities={ambiguities} "
        "processor={processor} service_template={service_template} "
        "archive_csar={archive_csar}".format(**stats)
    )
    print(
        f"grammar_tables={coverage['grammar_tables']} "
        f"grammar_rows={coverage['grammar_rows']} "
        f"additional_blocks={coverage['additional_blocks']} "
        f"normative_types={coverage['normative_type_sections']}"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
