# TOSCA 1.3 normative catalog extraction report

## Source and method

- Source: `docs/specifications/tosca/1.3/TOSCA-Simple-Profile-YAML-v1.3-os.html`
- SHA-256: `5c44638e1f9c532eb88fdc0cda1dce03c8f923a68148d26938797cca7b953df6`
- Source encoding: Windows-1252, as declared by the HTML.
- First pass: structural HTML parse of numbered headings, prose/list blocks, and tables; atomic records retain source line/table/row locators.
- Second pass: an independent coverage inventory checks all numbered headings, grammar rows, Additional Requirements blocks, exact uppercase normative markers, normative type headings, and clause 14.3 blocks.
- Three nested Word-export tables in section 5 are recovered as independent normative type-definition tables; their empty outer wrapper tables remain in the raw structural count but are not grammar.
- Explicit Example scopes are excluded under section 1.6.1.

## Second-pass completeness checks

| Check | Expected | Missing |
|---|---:|---:|
| Numbered headings classified | 1121 | 0 |
| Grammar tables inventoried | 388 | 0 |
| Grammar data rows represented | 763 | 0 |
| Additional Requirements prose blocks represented | 105 | 0 |
| Normative type headings represented | 66 | 0 |
| Included prose blocks with uppercase markers represented | 119 | 0 |
| Included table rows with uppercase markers represented | 7 | 0 |
| Clause 14.3 prose blocks represented | 6 | 0 |
| Required/optional table markers represented | 345 | 0 |
| Explicit `default:` occurrences represented | 35 | 0 |
| Short/long notation sections represented | 15 | 0 |
| Error-bearing prose blocks represented | 25 | 0 |

Raw-source checks use direct tag and keyword scans rather than the extraction traversal:

| Raw-source measure | Raw count | Parsed count |
|---|---:|---:|
| Heading tags | 1121 | 1121 |
| Table tags | 633 | 633 |
| Table-row tags | 1294 | 1287 |
| Uppercase normative markers | 144 | 144 |

The 1,294 raw `<tr>` tags include seven empty rows from layout or outer-wrapper tables; the parsed count contains 1,287 non-empty rows. All 633 table tags are nevertheless inventoried.

## Exact uppercase normative-word inventory

| Word | Included occurrences | Excluded occurrences |
|---|---:|---:|
| MAY | 22 | 1 |
| MUST | 23 | 0 |
| MUST NOT | 6 | 0 |
| OPTIONAL | 1 | 0 |
| RECOMMENDED | 1 | 0 |
| REQUIRED | 1 | 0 |
| SHALL | 50 | 0 |
| SHALL NOT | 5 | 0 |
| SHOULD | 29 | 0 |
| SHOULD NOT | 5 | 0 |

Excluded marker-bearing source locations are retained here so the exclusion is reviewable:

| Source | Section | Markers | Exclusion reason |
|---|---|---|---|
| line 47788 | 13.4.1.1.1 | MAY | explicit Example section or child of one |

## Additional Requirements inventory

| Section | Title | Source line |
|---|---|---:|
| 3.1.3.1 | Additional Requirements | 7013 |
| 3.3.2.5 | Additional Requirements | 8001 |
| 3.3.6.2 | Additional requirements | 8564 |
| 3.6.3.3 | Additional Requirements | 10475 |
| 3.6.4.2 | Additional Requirements | 10716 |
| 3.6.5.4 | Additional requirements | 10956 |
| 3.6.10.5 | Additional Requirements | 12796 |
| 3.6.12.4 | Additional Requirements | 13378 |
| 3.6.13.3 | Additional requirements | 13510 |
| 3.6.14.3 | Additional Requirements | 13780 |
| 3.6.17.3 | Additional requirements | 14614 |
| 3.6.25.4 | Additional Requirement | 16930 |
| 3.7.1.3 | Additional Requirements | 17755 |
| 3.7.2.4 | Additional requirements | 18143 |
| 3.7.3.3 | Additional Requirements | 18559 |
| 3.7.4.4 | Additional Requirements | 18882 |
| 3.7.5.4 | Additional Requirements | 19149 |
| 3.7.6.3 | Additional Requirements | 19413 |
| 3.7.9.3 | Additional Requirements | 20040 |
| 3.7.11.4 | Additional Requirements | 20611 |
| 3.8.3.3 | Additional requirements | 22208 |
| 3.8.4.3 | Additional requirements | 22547 |
| 3.8.5.4 | Additional Requirements | 22828 |
| 3.8.10.3 | Additional requirements | 23876 |
| 3.8.11.3 | Additional requirements | 24082 |
| 3.8.13.4 | Additional requirements | 24428 |
| 5.2.1 | Additional requirements | 29328 |
| 5.3.3 | Additional Requirements | 29520 |
| 5.3.5 | Additional Requirements | 29674 |
| 5.3.6.3 | Additional requirements | 29936 |
| 5.3.8.4 | Additional Requirements | 30457 |
| 5.3.9.4 | Additional Requirements | 30721 |
| 5.3.11.3 | Additional requirements | 31102 |
| 5.4.3.2 | Additional Requirements | 31324 |
| 5.4.4.2 | Additional Requirements | 31463 |
| 5.5.7.4 | Additional requirements | 32831 |
| 5.5.8.2 | Additional requirements | 32943 |
| 5.5.9.3 | Additional requirements | 33085 |
| 5.5.12.3 | Additional Requirements | 33555 |
| 5.8.1 | Additional Requirements | 34599 |
| 5.9.1.4 | Additional Requirements | 35284 |
| 5.9.3.4 | Additional Requirements | 35805 |
| 5.9.4.4 | Additional Requirements | 36043 |
| 5.9.5.3 | Additional Requirements | 36186 |
| 5.9.11.4 | Additional Requirements | 37290 |

## Grammar table inventory

| Table | Section | Section title | Source line | Data rows | Covered rows |
|---:|---|---|---:|---:|---:|
| 41 | 3.1 | TOSCA Namespace URI and alias | 6899 | 1 | 1 |
| 42 | 3.1.1 | TOSCA Namespace prefix | 6944 | 1 | 1 |
| 47 | 3.3.1 | Referenced YAML Types | 7762 | 6 | 6 |
| 48 | 3.3.2 | TOSCA version | 7870 | 1 | 1 |
| 49 | 3.3.2.1 | Grammar | 7899 | 1 | 1 |
| 51 | 3.3.3 | TOSCA range type | 8020 | 1 | 1 |
| 52 | 3.3.3.1 | Grammar | 8049 | 1 | 1 |
| 53 | 3.3.3.2 | Keywords | 8078 | 1 | 1 |
| 55 | 3.3.4 | TOSCA list type | 8154 | 1 | 1 |
| 56 | 3.3.4.1.1 | Square bracket notation | 8186 | 1 | 1 |
| 57 | 3.3.4.1.2 | Bulleted list notation | 8198 | 1 | 1 |
| 62 | 3.3.5 | TOSCA map type | 8347 | 1 | 1 |
| 63 | 3.3.5.1.1 | Single-line grammar | 8379 | 1 | 1 |
| 64 | 3.3.5.1.2 | Multi-line grammar | 8391 | 1 | 1 |
| 69 | 3.3.6.1 | Grammar | 8541 | 1 | 1 |
| 70 | 3.3.6.3 | Concrete Types | 8603 | 1 | 1 |
| 71 | 3.3.6.4.1 | Recognized Units | 8661 | 9 | 9 |
| 73 | 3.3.6.5.1 | Recognized Units | 8875 | 7 | 7 |
| 75 | 3.3.6.6.1 | Recognized Units | 9045 | 4 | 4 |
| 77 | 3.3.6.7.1 | Recognized Units | 9163 | 18 | 18 |
| 79 | 3.4.1 | Node States | 9544 | 1 | 1 |
| 80 | 3.4.2 | Relationship States | 9767 | 1 | 1 |
| 81 | 3.4.3 | Directives | 9824 | 4 | 4 |
| 82 | 3.4.4 | Network Name aliases | 9905 | 2 | 2 |
| 83 | 3.6.1.1 | Keyname | 10002 | 1 | 1 |
| 84 | 3.6.1.2 | Grammar | 10015 | 1 | 1 |
| 87 | 3.6.2.1 | Keyname | 10086 | 1 | 1 |
| 88 | 3.6.2.2 | Grammar | 10099 | 1 | 1 |
| 90 | 3.6.3.1 | Operator keynames | 10150 | 12 | 12 |
| 91 | 3.6.3.4 | Grammar | 10514 | 1 | 1 |
| 93 | 3.6.4.1.1 | Short notation: | 10666 | 1 | 1 |
| 94 | 3.6.4.1.2 | Extended notation: | 10683 | 1 | 1 |
| 95 | 3.6.5.1 | Keynames | 10736 | 2 | 2 |
| 96 | 3.6.5.2 | Additional filtering on named Capability properties | 10817 | 1 | 1 |
| 97 | 3.6.5.3 | Grammar | 10877 | 1 | 1 |
| 99 | 3.6.6.1 | Keynames | 11018 | 3 | 3 |
| 100 | 3.6.6.2.1 | Single-line grammar (no credential): | 11135 | 1 | 1 |
| 101 | 3.6.6.2.2 | Multi-line grammar | 11148 | 1 | 1 |
| 103 | 3.6.7.1 | Keynames | 11233 | 9 | 9 |
| 104 | 3.6.7.2.1 | Short notation | 11479 | 1 | 1 |
| 105 | 3.6.7.2.2 | Extended notation: | 11497 | 1 | 1 |
| 108 | 3.6.8.1 | Keynames | 11690 | 4 | 4 |
| 109 | 3.6.8.2.1 | Single-line grammar: | 11837 | 1 | 1 |
| 110 | 3.6.8.2.2 | Multi-line grammar | 11852 | 1 | 1 |
| 112 | 3.6.9.1 | Keynames | 12029 | 5 | 5 |
| 113 | 3.6.9.2 | Grammar | 12201 | 1 | 1 |
| 114 | 3.6.10.2 | Keynames | 12293 | 10 | 10 |
| 115 | 3.6.10.3 | Status values | 12613 | 4 | 4 |
| 116 | 3.6.10.4 | Grammar | 12682 | 1 | 1 |
| 120 | 3.6.11.2.1 | Short notation: | 13035 | 1 | 1 |
| 121 | 3.6.12.2 | Keynames | 13091 | 6 | 6 |
| 122 | 3.6.12.3 | Grammar | 13295 | 1 | 1 |
| 124 | 3.6.13.2.1 | Short notation: | 13451 | 1 | 1 |
| 125 | 3.6.13.2.2 | Extended notation: | 13468 | 1 | 1 |
| 126 | 3.6.14.1 | Keynames | 13534 | 2 | 2 |
| 127 | 3.6.14.2 | Grammar | 13629 | 1 | 1 |
| 128 | 3.6.14.2 | Grammar | 13673 | 1 | 1 |
| 129 | 3.6.14.2 | Grammar | 13689 | 1 | 1 |
| 132 | 3.6.15.1 | Grammar | 13864 | 1 | 1 |
| 133 | 3.6.15.1 | Grammar | 13881 | 4 | 4 |
| 134 | 3.6.16.1 | Keynames | 14046 | 4 | 4 |
| 135 | 3.6.16.2.1 | Short notation for use with single artifact | 14186 | 1 | 1 |
| 136 | 3.6.16.2.2 | Short notation for use with multiple artifact | 14212 | 1 | 1 |
| 137 | 3.6.16.2.3 | Extended notation for use with single artifact | 14244 | 1 | 1 |
| 138 | 3.6.16.2.4 | Extended notation for use with multiple artifacts | 14271 | 1 | 1 |
| 139 | 3.6.17.1 | Keynames | 14340 | 5 | 5 |
| 140 | 3.6.17.2.1 | Short notation | 14494 | 1 | 1 |
| 141 | 3.6.17.2.2 | Extended notation for use in Type definitions | 14514 | 1 | 1 |
| 142 | 3.6.17.2.3 | Extended notation for use in Template definitions | 14548 | 1 | 1 |
| 146 | 3.6.18.1 | Keynames | 14727 | 2 | 2 |
| 147 | 3.6.18.2.1 | Short notation for use with single artifact | 14808 | 1 | 1 |
| 148 | 3.6.18.2.2 | Short notation for use with multiple artifact | 14834 | 1 | 1 |
| 149 | 3.6.19.1 | Keynames | 14877 | 3 | 3 |
| 150 | 3.6.19.2 | Grammar | 14977 | 1 | 1 |
| 151 | 3.6.20.1 | Keynames | 15031 | 4 | 4 |
| 152 | 3.6.20.2.1 | Extended notation for use in Type definitions | 15159 | 1 | 1 |
| 153 | 3.6.20.2.2 | Extended notation for use in Template definitions | 15195 | 1 | 1 |
| 154 | 3.6.21.1 | Keynames | 15296 | 3 | 3 |
| 155 | 3.6.21.2 | Grammar | 15394 | 1 | 1 |
| 156 | 3.6.22.1 | Keynames | 15454 | 6 | 6 |
| 157 | 3.6.22.2 | Additional keynames for the extended condition notation | 15622 | 4 | 4 |
| 158 | 3.6.22.3.1 | Short notation | 15746 | 1 | 1 |
| 159 | 3.6.22.3.2 | Extended notation: | 15788 | 1 | 1 |
| 160 | 3.6.23.1.1 | Keynames | 15945 | 3 | 3 |
| 161 | 3.6.23.1.2.1 | Short notation | 16069 | 1 | 1 |
| 162 | 3.6.23.1.2.2 | Extended notation | 16086 | 1 | 1 |
| 163 | 3.6.23.2.1 | Keynames | 16132 | 1 | 1 |
| 164 | 3.6.23.2.2 | Grammar | 16190 | 1 | 1 |
| 165 | 3.6.23.3.1 | Keynames | 16223 | 3 | 3 |
| 166 | 3.6.23.3.2.1 | Short notation | 16348 | 1 | 1 |
| 167 | 3.6.23.3.2.2 | Extended notation | 16363 | 1 | 1 |
| 168 | 3.6.23.4.1 | Keynames | 16416 | 3 | 3 |
| 169 | 3.6.23.4.2.1 | Short notation | 16534 | 1 | 1 |
| 170 | 3.6.23.4.2.2 | Extended notation | 16549 | 1 | 1 |
| 172 | 3.6.24.2 | Grammar | 16624 | 1 | 1 |
| 175 | 3.6.25.1 | Keynames | 16695 | 4 | 4 |
| 176 | 3.6.25.2.1 | And clause | 16836 | 1 | 1 |
| 177 | 3.6.25.2.2 | Or clause | 16858 | 1 | 1 |
| 178 | 3.6.25.2.3 | Not clause | 16880 | 1 | 1 |
| 179 | 3.6.25.3 | Direct assertion definition | 16903 | 1 | 1 |
| 188 | 3.6.26.1 | Keynames | 17116 | 3 | 3 |
| 189 | 3.6.26.2 | Grammar | 17216 | 1 | 1 |
| 190 | 3.6.27.1 | Keynames | 17268 | 7 | 7 |
| 191 | 3.6.27.2 | Grammar | 17473 | 1 | 1 |
| 192 | 3.7.1.1 | Keynames | 17568 | 4 | 4 |
| 193 | 3.7.1.2 | Grammar | 17711 | 1 | 1 |
| 194 | 3.7.2.1 | Keynames | 17780 | 6 | 6 |
| 195 | 3.7.2.2.1 | Short notation | 17996 | 1 | 1 |
| 196 | 3.7.2.2.2 | Extended notation | 18015 | 1 | 1 |
| 199 | 3.7.3.1 | Keynames | 18191 | 4 | 4 |
| 200 | 3.7.3.1.1 | Additional Keynames for multi-line relationship grammar | 18349 | 2 | 2 |
| 201 | 3.7.3.2.1 | Simple grammar (Capability Type only) | 18444 | 1 | 1 |
| 202 | 3.7.3.2.2 | Extended grammar (with Node and Relationship Types) | 18458 | 1 | 1 |
| 203 | 3.7.3.2.3 | Extended grammar for declaring Property Definitions on the relationship’s Interfaces | 18491 | 1 | 1 |
| 204 | 3.7.4.1 | Keynames | 18641 | 3 | 3 |
| 205 | 3.7.4.2 | Grammar | 18756 | 1 | 1 |
| 207 | 3.7.5.1 | Keynames | 18930 | 3 | 3 |
| 208 | 3.7.5.2 | Grammar | 19027 | 1 | 1 |
| 210 | 3.7.6.1 | Keynames | 19188 | 4 | 4 |
| 211 | 3.7.6.2 | Grammar | 19315 | 1 | 1 |
| 214 | 3.7.7.1 | Keynames | 19510 | 3 | 3 |
| 215 | 3.7.7.2 | Grammar | 19609 | 1 | 1 |
| 217 | 3.7.9.1 | Keynames | 19746 | 6 | 6 |
| 218 | 3.7.9.2 | Grammar | 19916 | 1 | 1 |
| 220 | 3.7.10.1 | Keynames | 20131 | 4 | 4 |
| 221 | 3.7.10.2 | Grammar | 20253 | 1 | 1 |
| 223 | 3.7.11.1 | Keynames | 20409 | 3 | 3 |
| 224 | 3.7.11.2 | Grammar | 20512 | 1 | 1 |
| 226 | 3.7.12.1 | Keynames | 20662 | 3 | 3 |
| 227 | 3.7.12.2 | Grammar | 20767 | 1 | 1 |
| 229 | 3.8.1.1 | Keynames | 20894 | 3 | 3 |
| 230 | 3.8.1.2 | Grammar | 20996 | 1 | 1 |
| 232 | 3.8.2.1 | Keynames | 21094 | 5 | 5 |
| 233 | 3.8.2.1 | Keynames | 21268 | 3 | 3 |
| 234 | 3.8.2.2.1 | Short notation: | 21378 | 1 | 1 |
| 235 | 3.8.2.2.2 | Extended notation: | 21405 | 1 | 1 |
| 236 | 3.8.2.2.3 | Extended grammar with Property Assignments for the relationship’s Interfaces | 21453 | 1 | 1 |
| 240 | 3.8.3.1 | Keynames | 21751 | 12 | 12 |
| 241 | 3.8.3.2 | Grammar | 22057 | 1 | 1 |
| 243 | 3.8.4.1 | Keynames | 22263 | 7 | 7 |
| 244 | 3.8.4.2 | Grammar | 22451 | 1 | 1 |
| 246 | 3.8.5.1 | Keynames | 22600 | 5 | 5 |
| 247 | 3.8.5.2 | Grammar | 22740 | 1 | 1 |
| 249 | 3.8.6.1 | Keynames | 22870 | 6 | 6 |
| 250 | 3.8.6.2 | Grammar | 23037 | 1 | 1 |
| 252 | 3.8.7.1 | Keynames | 23151 | 7 | 7 |
| 253 | 3.8.7.2 | Grammar | 23345 | 1 | 1 |
| 254 | 3.8.8.1 | Keynames | 23444 | 2 | 2 |
| 255 | 3.8.8.2 | Grammar | 23519 | 1 | 1 |
| 256 | 3.8.8.2 | Grammar | 23536 | 1 | 1 |
| 257 | 3.8.9.1 | Keynames | 23616 | 1 | 1 |
| 258 | 3.8.9.2 | Grammar | 23667 | 1 | 1 |
| 259 | 3.8.10.1 | Keynames | 23691 | 3 | 3 |
| 260 | 3.8.10.2 | Grammar | 23795 | 1 | 1 |
| 261 | 3.8.10.2 | Grammar | 23809 | 1 | 1 |
| 262 | 3.8.11.1 | Keynames | 23897 | 3 | 3 |
| 263 | 3.8.11.2 | Grammar | 24001 | 1 | 1 |
| 264 | 3.8.11.2 | Grammar | 24016 | 1 | 1 |
| 265 | 3.8.12.1 | Grammar | 24103 | 1 | 1 |
| 266 | 3.8.13.1 | Keynames | 24161 | 7 | 7 |
| 267 | 3.8.13.2 | Grammar | 24351 | 1 | 1 |
| 268 | 3.9.1 | Keynames | 24459 | 9 | 9 |
| 269 | 3.9.2 | Grammar | 24704 | 1 | 1 |
| 270 | 3.9.2.1.1 | Grammar | 24856 | 1 | 1 |
| 273 | 3.9.2.2.1 | grammar | 24932 | 1 | 1 |
| 275 | 3.9.2.3.1 | Grammar | 24992 | 1 | 1 |
| 277 | 3.9.2.4.1 | Grammar | 25049 | 1 | 1 |
| 279 | 3.9.2.5.1 | Grammar | 25095 | 1 | 1 |
| 281 | 3.9.2.6.1 | Grammar | 25172 | 1 | 1 |
| 283 | 3.9.2.7.1 | requirement_mapping | 25218 | 1 | 1 |
| 284 | 3.9.2.7.1 | requirement_mapping | 25233 | 1 | 1 |
| 286 | 3.10.1 | Keynames | 25435 | 16 | 16 |
| 287 | 3.10.1.1 | Metadata keynames | 25841 | 3 | 3 |
| 288 | 3.10.2 | Grammar | 25937 | 1 | 1 |
| 289 | 3.10.3.1.1 | Keyname | 26074 | 1 | 1 |
| 290 | 3.10.3.1.2 | Grammar | 26088 | 1 | 1 |
| 293 | 3.10.3.2.1 | Keyname | 26135 | 1 | 1 |
| 294 | 3.10.3.2.2 | Grammar | 26147 | 1 | 1 |
| 296 | 3.10.3.3.1 | Keyname | 26184 | 1 | 1 |
| 297 | 3.10.3.3.2 | Grammar | 26196 | 1 | 1 |
| 299 | 3.10.3.4.1 | Keyname | 26232 | 1 | 1 |
| 300 | 3.10.3.4.2 | Grammar | 26244 | 1 | 1 |
| 302 | 3.10.3.5.1 | Keyname | 26273 | 1 | 1 |
| 303 | 3.10.3.5.2 | Grammar | 26285 | 1 | 1 |
| 305 | 3.10.3.6.1 | Keyname | 26327 | 1 | 1 |
| 306 | 3.10.3.7.1 | Keyname | 26344 | 1 | 1 |
| 307 | 3.10.3.7.2 | Grammar | 26356 | 1 | 1 |
| 309 | 3.10.3.8.1 | Keyname | 26416 | 1 | 1 |
| 310 | 3.10.3.8.2 | Grammar | 26428 | 1 | 1 |
| 312 | 3.10.3.9.1 | Keyname | 26477 | 1 | 1 |
| 313 | 3.10.3.9.2 | Grammar | 26489 | 1 | 1 |
| 315 | 3.10.3.10.1 | Keyname | 26540 | 1 | 1 |
| 316 | 3.10.3.10.2 | Grammar | 26552 | 1 | 1 |
| 318 | 3.10.3.11.1 | Keyname | 26593 | 1 | 1 |
| 319 | 3.10.3.11.2 | Grammar | 26605 | 1 | 1 |
| 321 | 3.10.3.12.1 | Keyname | 26686 | 1 | 1 |
| 322 | 3.10.3.12.2 | Grammar | 26698 | 1 | 1 |
| 324 | 3.10.3.13.1 | Keyname | 26753 | 1 | 1 |
| 325 | 3.10.3.13.2 | Grammar | 26765 | 1 | 1 |
| 327 | 3.10.3.14.1 | Keyname | 26813 | 1 | 1 |
| 328 | 3.10.3.14.2 | Grammar | 26825 | 1 | 1 |
| 330 | 3.10.3.15.1 | Keyname | 26881 | 1 | 1 |
| 331 | 3.10.3.15.2 | Grammar | 26893 | 1 | 1 |
| 333 | 3.10.3.16.1 | Keyname | 26958 | 1 | 1 |
| 334 | 3.10.3.16.2 | Grammar | 26970 | 1 | 1 |
| 336 | 3.10.3.17.1 | Keyname | 27011 | 1 | 1 |
| 337 | 3.10.3.17.2 | Grammar | 27023 | 1 | 1 |
| 339 | 4.1 | Reserved Function Keywords | 27085 | 4 | 4 |
| 340 | 4.2.1 | Reserved Environment Variable Names and Usage | 27210 | 4 | 4 |
| 344 | 4.3.1.1 | Grammar | 27535 | 1 | 1 |
| 345 | 4.3.1.2 | Parameters | 27548 | 1 | 1 |
| 347 | 4.3.2.1 | Grammar | 27631 | 1 | 1 |
| 348 | 4.3.2.2 | Parameters | 27644 | 2 | 2 |
| 350 | 4.3.3.1 | Grammar | 27752 | 1 | 1 |
| 351 | 4.3.3.2 | Parameters | 27766 | 3 | 3 |
| 353 | 4.4.1.1 | Grammar | 27912 | 1 | 1 |
| 354 | 4.4.1.1 | Grammar | 27925 | 1 | 1 |
| 355 | 4.4.1.2 | Parameters | 27939 | 2 | 2 |
| 358 | 4.4.2.1 | Grammar | 28150 | 1 | 1 |
| 359 | 4.4.2.2 | Parameters | 28166 | 4 | 4 |
| 363 | 4.5.1.1 | Grammar | 28441 | 1 | 1 |
| 364 | 4.5.1.2 | Parameters | 28456 | 4 | 4 |
| 365 | 4.6.1.1 | Grammar | 28631 | 1 | 1 |
| 366 | 4.6.1.2 | Parameters | 28645 | 4 | 4 |
| 367 | 4.7.1.1 | Grammar | 28791 | 1 | 1 |
| 368 | 4.7.1.2 | Parameters | 28804 | 1 | 1 |
| 369 | 4.7.1.3 | Returns | 28854 | 1 | 1 |
| 370 | 4.8.1.1 | Grammar | 28909 | 1 | 1 |
| 371 | 4.8.1.2 | Parameters | 28923 | 4 | 4 |
| 375 | 5.3.1.1 | Definition | 29357 | 1 | 1 |
| 376 | 5.3.2 | tosca.datatypes.json | 29377 | 2 | 2 |
| 377 | 5.3.2.1 | Definition | 29417 | 1 | 1 |
| 380 | 5.3.4 | tosca.datatypes.xml | 29535 | 2 | 2 |
| 381 | 5.3.4.1 | Definition | 29575 | 1 | 1 |
| 384 | 5.3.6 | tosca.datatypes.Credential | 29687 | 2 | 2 |
| 385 | 5.3.6.1 | Properties | 29725 | 5 | 5 |
| 386 | 5.3.6.2 | Definition | 29891 | 1 | 1 |
| 391 | 5.3.6.6 | OpenStack SSH Keypair | 30057 | 1 | 1 |
| 392 | 5.3.7 | tosca.datatypes.TimeInterval | 30087 | 2 | 2 |
| 393 | 5.3.7.1 | Properties | 30125 | 2 | 2 |
| 394 | 5.3.7.2 | Definition | 30211 | 1 | 1 |
| 396 | 5.3.8 | tosca.datatypes.network.NetworkInfo | 30265 | 2 | 2 |
| 397 | 5.3.8.1 | Properties | 30303 | 3 | 3 |
| 398 | 5.3.8.2 | Definition | 30400 | 1 | 1 |
| 400 | 5.3.9 | tosca.datatypes.network.PortInfo | 30476 | 2 | 2 |
| 401 | 5.3.9.1 | Properties | 30514 | 5 | 5 |
| 402 | 5.3.9.2 | Definition | 30652 | 1 | 1 |
| 404 | 5.3.10 | tosca.datatypes.network.PortDef | 30740 | 2 | 2 |
| 405 | 5.3.10.1 | Definition | 30780 | 1 | 1 |
| 408 | 5.3.11 | tosca.datatypes.network.PortSpec | 30839 | 2 | 2 |
| 409 | 5.3.11.1 | Properties | 30877 | 5 | 5 |
| 410 | 5.3.11.2 | Definition | 31042 | 1 | 1 |
| 412 | 5.4.1.1 | Definition | 31215 | 1 | 1 |
| 413 | 5.4.2 | tosca.artifacts.File | 31235 | 2 | 2 |
| 414 | 5.4.2.1 | Definition | 31273 | 1 | 1 |
| 416 | 5.4.3.1.1 | Definition | 31302 | 1 | 1 |
| 417 | 5.4.3.3 | tosca.artifacts.Deployment.Image | 31338 | 2 | 2 |
| 418 | 5.4.3.3.1 | Definition | 31376 | 1 | 1 |
| 419 | 5.4.3.4.1 | Definition | 31399 | 1 | 1 |
| 421 | 5.4.4.1.1 | Definition | 31441 | 1 | 1 |
| 422 | 5.4.4.3 | tosca.artifacts.Implementation.Bash | 31476 | 2 | 2 |
| 423 | 5.4.4.3.1 | Definition | 31514 | 1 | 1 |
| 424 | 5.4.4.4 | tosca.artifacts.Implementation.Python | 31539 | 2 | 2 |
| 425 | 5.4.4.4.1 | Definition | 31577 | 1 | 1 |
| 427 | 5.4.5.1.1 | Definition | 31623 | 1 | 1 |
| 428 | 5.5.1.1 | Definition | 31656 | 1 | 1 |
| 429 | 5.5.2 | tosca.capabilities.Node | 31674 | 2 | 2 |
| 430 | 5.5.2.1 | Definition | 31712 | 1 | 1 |
| 431 | 5.5.3 | tosca.capabilities.Compute | 31731 | 2 | 2 |
| 432 | 5.5.3.1 | Properties | 31769 | 5 | 5 |
| 433 | 5.5.3.2 | Definition | 31937 | 1 | 1 |
| 434 | 5.5.4 | tosca.capabilities.Network | 32007 | 2 | 2 |
| 435 | 5.5.4.1 | Properties | 32045 | 1 | 1 |
| 436 | 5.5.4.2 | Definition | 32102 | 1 | 1 |
| 437 | 5.5.5 | tosca.capabilities.Storage | 32130 | 2 | 2 |
| 438 | 5.5.5.1 | Properties | 32168 | 1 | 1 |
| 439 | 5.5.5.2 | Definition | 32225 | 1 | 1 |
| 440 | 5.5.6 | tosca.capabilities.Container | 32253 | 2 | 2 |
| 441 | 5.5.6.1 | Properties | 32291 | 1 | 1 |
| 442 | 5.5.6.2 | Definition | 32347 | 1 | 1 |
| 443 | 5.5.7 | tosca.capabilities.Endpoint | 32371 | 2 | 2 |
| 444 | 5.5.7.1 | Properties | 32409 | 8 | 8 |
| 445 | 5.5.7.2 | Attributes | 32678 | 1 | 1 |
| 446 | 5.5.7.3 | Definition | 32736 | 1 | 1 |
| 447 | 5.5.8 | tosca.capabilities.Endpoint.Public | 32852 | 2 | 2 |
| 448 | 5.5.8.1 | Definition | 32890 | 1 | 1 |
| 449 | 5.5.9 | tosca.capabilities.Endpoint.Admin | 32965 | 2 | 2 |
| 450 | 5.5.9.1 | Properties | 33003 | 1 | 1 |
| 451 | 5.5.9.2 | Definition | 33059 | 1 | 1 |
| 452 | 5.5.10 | tosca.capabilities.Endpoint.Database | 33099 | 2 | 2 |
| 453 | 5.5.10.1 | Properties | 33137 | 1 | 1 |
| 454 | 5.5.10.2 | Definition | 33193 | 1 | 1 |
| 455 | 5.5.11 | tosca.capabilities.Attachment | 33214 | 2 | 2 |
| 456 | 5.5.11.1 | Properties | 33252 | 1 | 1 |
| 457 | 5.5.11.2 | Definition | 33308 | 1 | 1 |
| 458 | 5.5.12 | tosca.capabilities.OperatingSystem | 33326 | 2 | 2 |
| 459 | 5.5.12.1 | Properties | 33364 | 4 | 4 |
| 460 | 5.5.12.2 | Definition | 33513 | 1 | 1 |
| 461 | 5.5.13 | tosca.capabilities.Scalable | 33574 | 2 | 2 |
| 462 | 5.5.13.1 | Properties | 33612 | 3 | 3 |
| 463 | 5.5.13.2 | Definition | 33730 | 1 | 1 |
| 464 | 5.5.14 | tosca.capabilities.network.Bindable | 33779 | 2 | 2 |
| 465 | 5.5.14.1 | Properties | 33817 | 1 | 1 |
| 466 | 5.5.14.2 | Definition | 33873 | 1 | 1 |
| 467 | 5.7.1.1 | Attributes | 33907 | 3 | 3 |
| 468 | 5.7.1.2 | Definition | 34023 | 1 | 1 |
| 469 | 5.7.2 | tosca.relationships.DependsOn | 34058 | 2 | 2 |
| 470 | 5.7.2.1 | Definition | 34096 | 1 | 1 |
| 471 | 5.7.3 | tosca.relationships.HostedOn | 34118 | 2 | 2 |
| 472 | 5.7.3.1 | Definition | 34156 | 1 | 1 |
| 473 | 5.7.4 | tosca.relationships.ConnectsTo | 34178 | 2 | 2 |
| 474 | 5.7.4.1 | Definition | 34216 | 1 | 1 |
| 475 | 5.7.4.2 | Properties | 34243 | 1 | 1 |
| 476 | 5.7.5 | tosca.relationships.AttachesTo | 34306 | 2 | 2 |
| 477 | 5.7.5.1 | Properties | 34344 | 2 | 2 |
| 478 | 5.7.5.2 | Attributes | 34436 | 1 | 1 |
| 479 | 5.7.5.3 | Definition | 34495 | 1 | 1 |
| 480 | 5.7.6 | tosca.relationships.RoutesTo | 34534 | 2 | 2 |
| 481 | 5.7.6.1 | Definition | 34572 | 1 | 1 |
| 482 | 5.8.3.1 | Definition | 34632 | 1 | 1 |
| 483 | 5.8.4 | tosca.interfaces.node.lifecycle.Standard | 34652 | 2 | 2 |
| 484 | 5.8.4.1 | Definition | 34690 | 1 | 1 |
| 485 | 5.8.5 | tosca.interfaces.relationship.Configure | 34784 | 2 | 2 |
| 486 | 5.8.5.1 | Definition | 34822 | 1 | 1 |
| 488 | 5.9.1 | tosca.nodes.Root | 35011 | 2 | 2 |
| 489 | 5.9.1.1 | Properties | 35049 | 1 | 1 |
| 490 | 5.9.1.2 | Attributes | 35106 | 3 | 3 |
| 491 | 5.9.1.3 | Definition | 35222 | 1 | 1 |
| 492 | 5.9.2 | tosca.nodes.Abstract.Compute | 35304 | 2 | 2 |
| 493 | 5.9.2.1 | Properties | 35342 | 1 | 1 |
| 494 | 5.9.2.2 | Attributes | 35398 | 1 | 1 |
| 495 | 5.9.2.3 | Definition | 35454 | 1 | 1 |
| 496 | 5.9.3 | tosca.nodes.Compute | 35483 | 2 | 2 |
| 497 | 5.9.3.1 | Properties | 35521 | 1 | 1 |
| 498 | 5.9.3.2 | Attributes | 35577 | 4 | 4 |
| 499 | 5.9.3.3 | Definition | 35719 | 1 | 1 |
| 500 | 5.9.4 | tosca.nodes.SoftwareComponent | 35822 | 2 | 2 |
| 501 | 5.9.4.1 | Properties | 35860 | 2 | 2 |
| 502 | 5.9.4.2 | Attributes | 35944 | 1 | 1 |
| 503 | 5.9.4.3 | Definition | 36000 | 1 | 1 |
| 504 | 5.9.5 | tosca.nodes.WebServer | 36060 | 2 | 2 |
| 505 | 5.9.5.1 | Properties | 36098 | 1 | 1 |
| 506 | 5.9.5.2 | Definition | 36154 | 1 | 1 |
| 507 | 5.9.6 | tosca.nodes.WebApplication | 36204 | 2 | 2 |
| 508 | 5.9.6.1 | Properties | 36242 | 1 | 1 |
| 509 | 5.9.6.2 | Definition | 36300 | 1 | 1 |
| 510 | 5.9.7.1 | Properties | 36349 | 2 | 2 |
| 511 | 5.9.7.2 | Definition | 36432 | 1 | 1 |
| 512 | 5.9.8 | tosca.nodes.Database | 36481 | 2 | 2 |
| 513 | 5.9.8.1 | Properties | 36519 | 4 | 4 |
| 514 | 5.9.8.2 | Definition | 36656 | 1 | 1 |
| 515 | 5.9.9 | tosca.nodes.Abstract.Storage | 36723 | 2 | 2 |
| 516 | 5.9.9.1 | Properties | 36761 | 2 | 2 |
| 517 | 5.9.9.2 | Definition | 36846 | 1 | 1 |
| 518 | 5.9.10 | tosca.nodes.Storage.ObjectStorage | 36886 | 2 | 2 |
| 519 | 5.9.10.1 | Properties | 36924 | 1 | 1 |
| 520 | 5.9.10.2 | Definition | 36982 | 1 | 1 |
| 521 | 5.9.11 | tosca.nodes.Storage.BlockStorage | 37034 | 2 | 2 |
| 522 | 5.9.11.1 | Properties | 37072 | 3 | 3 |
| 523 | 5.9.11.2 | Attributes | 37200 | 1 | 1 |
| 524 | 5.9.11.3 | Definition | 37256 | 1 | 1 |
| 525 | 5.9.12 | tosca.nodes.Container.Runtime | 37327 | 2 | 2 |
| 526 | 5.9.12.1 | Definition | 37365 | 1 | 1 |
| 527 | 5.9.13 | tosca.nodes.Container.Application | 37400 | 2 | 2 |
| 528 | 5.9.13.1 | Definition | 37438 | 1 | 1 |
| 529 | 5.9.14 | tosca.nodes.LoadBalancer | 37490 | 2 | 2 |
| 530 | 5.9.14.1 | Definition | 37528 | 1 | 1 |
| 531 | 5.10.1.1 | Definition | 37610 | 1 | 1 |
| 532 | 5.11.1.1 | Definition | 37660 | 1 | 1 |
| 533 | 5.11.2.1 | Definition | 37679 | 1 | 1 |
| 534 | 5.11.3.1 | Definition | 37702 | 1 | 1 |
| 535 | 5.11.4.1 | Definition | 37725 | 1 | 1 |
| 536 | 5.11.5.1 | Definition | 37749 | 1 | 1 |
| 553 | 8.5.1 | tosca.nodes.network.Network | 39748 | 2 | 2 |
| 554 | 8.5.1.1 | Properties | 39786 | 11 | 11 |
| 555 | 8.5.1.2 | Attributes | 40170 | 1 | 1 |
| 556 | 8.5.1.3 | Definition | 40227 | 1 | 1 |
| 557 | 8.5.2 | tosca.nodes.network.Port | 40333 | 2 | 2 |
| 558 | 8.5.2.1 | Properties | 40371 | 5 | 5 |
| 559 | 8.5.2.2 | Attributes | 40558 | 1 | 1 |
| 560 | 8.5.2.3 | Definition | 40615 | 1 | 1 |
| 561 | 8.5.3 | tosca.capabilities.network.Linkable | 40697 | 2 | 2 |
| 562 | 8.5.3.1 | Properties | 40735 | 1 | 1 |
| 563 | 8.5.3.2 | Definition | 40791 | 1 | 1 |
| 564 | 8.5.4 | tosca.relationships.network.LinksTo | 40809 | 2 | 2 |
| 565 | 8.5.4.1 | Definition | 40847 | 1 | 1 |
| 566 | 8.5.5 | tosca.relationships.network.BindsTo | 40870 | 2 | 2 |
| 567 | 8.5.5.1 | Definition | 40908 | 1 | 1 |

## Review disposition for every numbered section

Classification totals: `excluded`=363, `included`=658, `included-orchestrator-only`=89, `reviewed-informative`=11.

| Section | Title | Source line | Disposition | Basis |
|---|---|---:|---|---|
| 1 | Introduction | 2326 | reviewed-informative | introductory or convention material |
| 1.1 | IPR Policy | 2331 | reviewed-informative | introductory or convention material |
| 1.2 | Objective | 2348 | reviewed-informative | introductory or convention material |
| 1.3 | Summary of key TOSCA concepts | 2373 | reviewed-informative | introductory or convention material |
| 1.4 | Implementations | 2416 | reviewed-informative | introductory or convention material |
| 1.5 | Terminology | 2463 | reviewed-informative | introductory or convention material |
| 1.6 | Notational Conventions | 2481 | reviewed-informative | introductory or convention material |
| 1.6.1 | Notes | 2492 | reviewed-informative | introductory or convention material |
| 1.7 | Normative References | 2505 | reviewed-informative | introductory or convention material |
| 1.8 | Non-Normative References | 2576 | reviewed-informative | introductory or convention material |
| 1.9 | Glossary | 2773 | reviewed-informative | introductory or convention material |
| 2 | TOSCA by example | 2922 | excluded | explicit Example section or child of one |
| 2.1 | A “hello world” template for TOSCA Simple Profile in YAML | 2936 | excluded | explicit Example section or child of one |
| 2.1.1 | Requesting input parameters and providing output | 3063 | excluded | explicit Example section or child of one |
| 2.2 | TOSCA template for a simple software installation | 3167 | excluded | explicit Example section or child of one |
| 2.3 | Overriding behavior of predefined node types | 3298 | excluded | explicit Example section or child of one |
| 2.4 | TOSCA template for database content deployment | 3397 | excluded | explicit Example section or child of one |
| 2.5 | TOSCA template for a two-tier application | 3560 | excluded | explicit Example section or child of one |
| 2.6 | Using a custom script to establish a relationship in a template | 3734 | excluded | explicit Example section or child of one |
| 2.7 | Using custom relationship types in a TOSCA template | 3910 | excluded | explicit Example section or child of one |
| 2.7.1 | Definition of a custom relationship type | 3999 | excluded | explicit Example section or child of one |
| 2.8 | Defining generic dependencies between nodes in a template | 4043 | excluded | explicit Example section or child of one |
| 2.9 | Describing abstract requirements for nodes and capabilities in a TOSCA template | 4119 | excluded | explicit Example section or child of one |
| 2.9.1 | Using a node_filter to define hosting infrastructure requirements for a software | 4190 | excluded | explicit Example section or child of one |
| 2.9.2 | Using an abstract node template to define infrastructure requirements for software | 4302 | excluded | explicit Example section or child of one |
| 2.9.3 | Using a node_filter to define requirements on a database for an application | 4417 | excluded | explicit Example section or child of one |
| 2.10 | Using node template substitution for model composition | 4561 | excluded | explicit Example section or child of one |
| 2.10.1 | Understanding node template instantiation through a TOSCA Orchestrator | 4580 | excluded | explicit Example section or child of one |
| 2.10.2 | Definition of the top-level service template | 4612 | excluded | explicit Example section or child of one |
| 2.10.3 | Definition of the database stack in a service template | 4742 | excluded | explicit Example section or child of one |
| 2.11 | Using node template substitution for chaining subsystems | 4895 | excluded | explicit Example section or child of one |
| 2.11.1 | Defining the overall subsystem chain | 4912 | excluded | explicit Example section or child of one |
| 2.11.2 | Defining a subsystem (node) type | 5089 | excluded | explicit Example section or child of one |
| 2.11.3 | Defining the details of a subsystem | 5177 | excluded | explicit Example section or child of one |
| 2.12 | Using node template substitution to provide product choice | 5378 | excluded | explicit Example section or child of one |
| 2.12.1 | Defining a service template with vendor-independent component | 5391 | excluded | explicit Example section or child of one |
| 2.12.2 | Defining vendor-specific component options | 5503 | excluded | explicit Example section or child of one |
| 2.12.3 | Substitution matching using substitution filters | 5636 | excluded | explicit Example section or child of one |
| 2.13 | Grouping node templates | 5819 | excluded | explicit Example section or child of one |
| 2.14 | Using YAML Macros to simplify templates | 6038 | excluded | explicit Example section or child of one |
| 2.15 | Passing information as inputs to Interfaces and Operations | 6126 | excluded | explicit Example section or child of one |
| 2.15.1 | Example: declaring input variables for all operations on a single interface | 6144 | excluded | explicit Example section or child of one |
| 2.15.2 | Example: declaring input variables for a single operation | 6176 | excluded | explicit Example section or child of one |
| 2.16 | Returning output values from operations | 6220 | excluded | explicit Example section or child of one |
| 2.16.1 | Example: setting output values to a node attribute | 6230 | excluded | explicit Example section or child of one |
| 2.16.2 | Example: setting output values to a capability attribute | 6277 | excluded | explicit Example section or child of one |
| 2.17 | Receiving asynchronous notifications | 6329 | excluded | explicit Example section or child of one |
| 2.18 | Topology Template Model versus Instance Model | 6404 | excluded | explicit Example section or child of one |
| 2.19 | Using attributes implicitly reflected from properties | 6422 | excluded | explicit Example section or child of one |
| 2.20 | Creating Multiple Node Instances from the Same Node Template | 6589 | excluded | explicit Example section or child of one |
| 2.20.1 | Specifying Number of Occurrences | 6705 | excluded | explicit Example section or child of one |
| 2.20.2 | Specifying Inputs | 6792 | excluded | explicit Example section or child of one |
| 3 | TOSCA Simple Profile definitions in YAML | 6873 | included | normative baseline scope |
| 3.1 | TOSCA Namespace URI and alias | 6887 | included | normative baseline scope |
| 3.1.1 | TOSCA Namespace prefix | 6936 | included | normative baseline scope |
| 3.1.2 | TOSCA Namespacing in TOSCA Service Templates | 6971 | included | normative baseline scope |
| 3.1.3 | Rules to avoid namespace collisions | 7002 | included | normative baseline scope |
| 3.1.3.1 | Additional Requirements | 7013 | included | normative baseline scope |
| 3.2 | Using Namespaces | 7137 | included | normative baseline scope |
| 3.2.1 | Example - Importing a Service Template and Namespaces | 7166 | excluded | explicit Example section or child of one |
| 3.2.1.1 | Conceptual Global Namespace URI and Namespace Prefix tracking | 7275 | excluded | explicit Example section or child of one |
| 3.2.1.2 | Conceptual Global Namespace and Type tracking | 7430 | excluded | explicit Example section or child of one |
| 3.3 | Parameter and property types | 7741 | included | normative baseline scope |
| 3.3.1 | Referenced YAML Types | 7748 | included | normative baseline scope |
| 3.3.1.1 | Notes | 7846 | included | normative baseline scope |
| 3.3.2 | TOSCA version | 7859 | included | normative baseline scope |
| 3.3.2.1 | Grammar | 7895 | included | normative baseline scope |
| 3.3.2.2 | Version Comparison | 7944 | included | normative baseline scope |
| 3.3.2.3 | Examples | 7970 | excluded | explicit Example section or child of one |
| 3.3.2.4 | Notes | 7994 | included | normative baseline scope |
| 3.3.2.5 | Additional Requirements | 8001 | included | normative baseline scope |
| 3.3.3 | TOSCA range type | 8013 | included | normative baseline scope |
| 3.3.3.1 | Grammar | 8045 | included | normative baseline scope |
| 3.3.3.2 | Keywords | 8073 | included | normative baseline scope |
| 3.3.3.3 | Examples | 8115 | excluded | explicit Example section or child of one |
| 3.3.4 | TOSCA list type | 8139 | included | normative baseline scope |
| 3.3.4.1 | Grammar | 8179 | included | normative baseline scope |
| 3.3.4.1.1 | Square bracket notation | 8184 | included | normative baseline scope |
| 3.3.4.1.2 | Bulleted list notation | 8196 | included | normative baseline scope |
| 3.3.4.2 | Declaration Examples | 8218 | excluded | explicit Example section or child of one |
| 3.3.4.2.1 | List declaration using a simple type | 8220 | excluded | explicit Example section or child of one |
| 3.3.4.2.2 | List declaration using a complex type | 8255 | excluded | explicit Example section or child of one |
| 3.3.4.3 | Definition Examples | 8285 | excluded | explicit Example section or child of one |
| 3.3.4.3.1 | Square bracket notation | 8301 | excluded | explicit Example section or child of one |
| 3.3.4.3.2 | Bulleted list notation | 8313 | excluded | explicit Example section or child of one |
| 3.3.5 | TOSCA map type | 8327 | included | normative baseline scope |
| 3.3.5.1 | Grammar | 8372 | included | normative baseline scope |
| 3.3.5.1.1 | Single-line grammar | 8377 | included | normative baseline scope |
| 3.3.5.1.2 | Multi-line grammar | 8389 | included | normative baseline scope |
| 3.3.5.2 | Declaration Examples | 8417 | excluded | explicit Example section or child of one |
| 3.3.5.2.1 | Map declaration using a simple type | 8419 | excluded | explicit Example section or child of one |
| 3.3.5.2.2 | Map declaration using a complex type | 8454 | excluded | explicit Example section or child of one |
| 3.3.5.3 | Definition Examples | 8483 | excluded | explicit Example section or child of one |
| 3.3.5.3.1 | Single-line notation | 8499 | excluded | explicit Example section or child of one |
| 3.3.5.3.2 | Multi-line notation | 8513 | excluded | explicit Example section or child of one |
| 3.3.6 | TOSCA scalar-unit type | 8530 | included | normative baseline scope |
| 3.3.6.1 | Grammar | 8536 | included | normative baseline scope |
| 3.3.6.2 | Additional requirements | 8564 | included | normative baseline scope |
| 3.3.6.3 | Concrete Types | 8601 | included | normative baseline scope |
| 3.3.6.4 | scalar-unit.size | 8657 | included | normative baseline scope |
| 3.3.6.4.1 | Recognized Units | 8659 | included | normative baseline scope |
| 3.3.6.4.2 | Examples | 8833 | excluded | explicit Example section or child of one |
| 3.3.6.4.3 | Notes | 8848 | included | normative baseline scope |
| 3.3.6.5 | scalar-unit.time | 8870 | included | normative baseline scope |
| 3.3.6.5.1 | Recognized Units | 8873 | included | normative baseline scope |
| 3.3.6.5.2 | Examples | 9007 | excluded | explicit Example section or child of one |
| 3.3.6.5.3 | Notes | 9022 | included | normative baseline scope |
| 3.3.6.6 | scalar-unit.frequency | 9041 | included | normative baseline scope |
| 3.3.6.6.1 | Recognized Units | 9043 | included | normative baseline scope |
| 3.3.6.6.2 | Examples | 9133 | excluded | explicit Example section or child of one |
| 3.3.6.6.3 | Notes | 9149 | included | normative baseline scope |
| 3.3.6.7 | scalar-unit.bitrate | 9159 | included | normative baseline scope |
| 3.3.6.7.1 | Recognized Units | 9161 | included | normative baseline scope |
| 3.3.6.7.2 | Examples | 9485 | excluded | explicit Example section or child of one |
| 3.3.6.7.3 | Notes | 9511 | included | normative baseline scope |
| 3.4 | Normative values | 9525 | included | normative baseline scope |
| 3.4.1 | Node States | 9530 | included | normative baseline scope |
| 3.4.2 | Relationship States | 9756 | included | normative baseline scope |
| 3.4.2.1 | Notes | 9812 | included | normative baseline scope |
| 3.4.3 | Directives | 9819 | included | normative baseline scope |
| 3.4.4 | Network Name aliases | 9897 | included | normative baseline scope |
| 3.4.4.1 | Usage | 9958 | included | normative baseline scope |
| 3.5 | TOSCA Metamodel | 9966 | included | normative baseline scope |
| 3.5.1 | Required Keynames | 9975 | included | normative baseline scope |
| 3.6 | Reusable modeling definitions | 9984 | included | normative baseline scope |
| 3.6.1 | Description definition | 9989 | included | normative baseline scope |
| 3.6.1.1 | Keyname | 9997 | included | normative baseline scope |
| 3.6.1.2 | Grammar | 10011 | included | normative baseline scope |
| 3.6.1.3 | Examples | 10026 | excluded | explicit Example section or child of one |
| 3.6.1.4 | Notes | 10066 | included | normative baseline scope |
| 3.6.2 | Metadata | 10075 | included | normative baseline scope |
| 3.6.2.1 | Keyname | 10081 | included | normative baseline scope |
| 3.6.2.2 | Grammar | 10095 | included | normative baseline scope |
| 3.6.2.3 | Examples | 10113 | excluded | explicit Example section or child of one |
| 3.6.2.4 | Notes | 10129 | included | normative baseline scope |
| 3.6.3 | Constraint clause | 10137 | included | normative baseline scope |
| 3.6.3.1 | Operator keynames | 10144 | included | normative baseline scope |
| 3.6.3.1.1 | Comparable value types | 10454 | included | normative baseline scope |
| 3.6.3.2 | Schema Constraint purpose | 10463 | included | normative baseline scope |
| 3.6.3.3 | Additional Requirements | 10475 | included | normative baseline scope |
| 3.6.3.4 | Grammar | 10509 | included | normative baseline scope |
| 3.6.3.5 | Examples | 10582 | excluded | explicit Example section or child of one |
| 3.6.4 | Property Filter definition | 10649 | included | normative baseline scope |
| 3.6.4.1 | Grammar | 10656 | included | normative baseline scope |
| 3.6.4.1.1 | Short notation: | 10661 | included | normative baseline scope |
| 3.6.4.1.2 | Extended notation: | 10678 | included | normative baseline scope |
| 3.6.4.2 | Additional Requirements | 10716 | included | normative baseline scope |
| 3.6.5 | Node Filter definition | 10724 | included | normative baseline scope |
| 3.6.5.1 | Keynames | 10731 | included | normative baseline scope |
| 3.6.5.2 | Additional filtering on named Capability properties | 10812 | included | normative baseline scope |
| 3.6.5.3 | Grammar | 10873 | included | normative baseline scope |
| 3.6.5.4 | Additional requirements | 10956 | included | normative baseline scope |
| 3.6.5.5 | Example | 10964 | excluded | explicit Example section or child of one |
| 3.6.6 | Repository definition | 11006 | included | normative baseline scope |
| 3.6.6.1 | Keynames | 11013 | included | normative baseline scope |
| 3.6.6.2 | Grammar | 11128 | included | normative baseline scope |
| 3.6.6.2.1 | Single-line grammar (no credential): | 11133 | included | normative baseline scope |
| 3.6.6.2.2 | Multi-line grammar | 11146 | included | normative baseline scope |
| 3.6.6.3 | Example | 11200 | excluded | explicit Example section or child of one |
| 3.6.7 | Artifact definition | 11221 | included | normative baseline scope |
| 3.6.7.1 | Keynames | 11228 | included | normative baseline scope |
| 3.6.7.2 | Grammar | 11469 | included | normative baseline scope |
| 3.6.7.2.1 | Short notation | 11474 | included | normative baseline scope |
| 3.6.7.2.2 | Extended notation: | 11492 | included | normative baseline scope |
| 3.6.7.3 | Examples | 11605 | excluded | explicit Example section or child of one |
| 3.6.8 | Import definition | 11678 | included | normative baseline scope |
| 3.6.8.1 | Keynames | 11685 | included | normative baseline scope |
| 3.6.8.2 | Grammar | 11831 | included | normative baseline scope |
| 3.6.8.2.1 | Single-line grammar: | 11835 | included | normative baseline scope |
| 3.6.8.2.2 | Multi-line grammar | 11850 | included | normative baseline scope |
| 3.6.8.2.3 | Requirements | 11902 | included | normative baseline scope |
| 3.6.8.2.4 | Import URI processing requirements | 11931 | included | normative baseline scope |
| 3.6.8.3 | Example | 11989 | excluded | explicit Example section or child of one |
| 3.6.9 | Schema Definition | 12014 | included | normative baseline scope |
| 3.6.9.1 | Keynames | 12024 | included | normative baseline scope |
| 3.6.9.2 | Grammar | 12197 | included | normative baseline scope |
| 3.6.10 | Property definition | 12269 | included | normative baseline scope |
| 3.6.10.1 | Attribute and Property reflection | 12279 | included | normative baseline scope |
| 3.6.10.2 | Keynames | 12288 | included | normative baseline scope |
| 3.6.10.3 | Status values | 12608 | included | normative baseline scope |
| 3.6.10.4 | Grammar | 12677 | included | normative baseline scope |
| 3.6.10.5 | Additional Requirements | 12796 | included | normative baseline scope |
| 3.6.10.6 | Refining Property Definitions | 12835 | included | normative baseline scope |
| 3.6.10.7 | Notes | 12894 | included | normative baseline scope |
| 3.6.10.8 | Examples | 12911 | excluded | explicit Example section or child of one |
| 3.6.11 | Property assignment | 13016 | included | normative baseline scope |
| 3.6.11.1 | Keynames | 13022 | included | normative baseline scope |
| 3.6.11.2 | Grammar | 13026 | included | normative baseline scope |
| 3.6.11.2.1 | Short notation: | 13030 | included | normative baseline scope |
| 3.6.12 | Attribute definition | 13065 | included | normative baseline scope |
| 3.6.12.1 | Attribute and Property reflection | 13077 | included | normative baseline scope |
| 3.6.12.2 | Keynames | 13086 | included | normative baseline scope |
| 3.6.12.3 | Grammar | 13291 | included | normative baseline scope |
| 3.6.12.4 | Additional Requirements | 13378 | included | normative baseline scope |
| 3.6.12.5 | Notes | 13395 | included | normative baseline scope |
| 3.6.12.6 | Example | 13413 | excluded | explicit Example section or child of one |
| 3.6.13 | Attribute assignment | 13431 | included | normative baseline scope |
| 3.6.13.1 | Keynames | 13438 | included | normative baseline scope |
| 3.6.13.2 | Grammar | 13442 | included | normative baseline scope |
| 3.6.13.2.1 | Short notation: | 13446 | included | normative baseline scope |
| 3.6.13.2.2 | Extended notation: | 13463 | included | normative baseline scope |
| 3.6.13.3 | Additional requirements | 13510 | included | normative baseline scope |
| 3.6.14 | Parameter definition | 13519 | included | normative baseline scope |
| 3.6.14.1 | Keynames | 13528 | included | normative baseline scope |
| 3.6.14.2 | Grammar | 13624 | included | normative baseline scope |
| 3.6.14.3 | Additional Requirements | 13780 | included | normative baseline scope |
| 3.6.14.4 | Example | 13804 | excluded | explicit Example section or child of one |
| 3.6.15 | Attribute Mapping definition | 13849 | included | normative baseline scope |
| 3.6.15.1 | Grammar | 13856 | included | normative baseline scope |
| 3.6.16 | Operation implementation definition | 14032 | included | normative baseline scope |
| 3.6.16.1 | Keynames | 14041 | included | normative baseline scope |
| 3.6.16.2 | Grammar | 14176 | included | normative baseline scope |
| 3.6.16.2.1 | Short notation for use with single artifact | 14181 | included | normative baseline scope |
| 3.6.16.2.2 | Short notation for use with multiple artifact | 14206 | included | normative baseline scope |
| 3.6.16.2.3 | Extended notation for use with single artifact | 14236 | included | normative baseline scope |
| 3.6.16.2.4 | Extended notation for use with multiple artifacts | 14264 | included | normative baseline scope |
| 3.6.17 | Operation definition | 14330 | included | normative baseline scope |
| 3.6.17.1 | Keynames | 14335 | included | normative baseline scope |
| 3.6.17.2 | Grammar | 14483 | included | normative baseline scope |
| 3.6.17.2.1 | Short notation | 14487 | included | normative baseline scope |
| 3.6.17.2.2 | Extended notation for use in Type definitions | 14508 | included | normative baseline scope |
| 3.6.17.2.3 | Extended notation for use in Template definitions | 14542 | included | normative baseline scope |
| 3.6.17.3 | Additional requirements | 14614 | included | normative baseline scope |
| 3.6.17.4 | Examples | 14634 | excluded | explicit Example section or child of one |
| 3.6.17.4.1 | Single-line example | 14636 | excluded | explicit Example section or child of one |
| 3.6.17.4.2 | Multi-line example with shorthand implementation definitions | 14652 | excluded | explicit Example section or child of one |
| 3.6.17.4.3 | Multi-line example with extended implementation definitions | 14680 | excluded | explicit Example section or child of one |
| 3.6.18 | Notification implementation definition | 14714 | included | normative baseline scope |
| 3.6.18.1 | Keynames | 14722 | included | normative baseline scope |
| 3.6.18.2 | Grammar | 14798 | included | normative baseline scope |
| 3.6.18.2.1 | Short notation for use with single artifact | 14803 | included | normative baseline scope |
| 3.6.18.2.2 | Short notation for use with multiple artifact | 14828 | included | normative baseline scope |
| 3.6.19 | Notification definition | 14853 | included | normative baseline scope |
| 3.6.19.1 | Keynames | 14872 | included | normative baseline scope |
| 3.6.19.2 | Grammar | 14972 | included | normative baseline scope |
| 3.6.20 | Interface definition | 15021 | included | normative baseline scope |
| 3.6.20.1 | Keynames | 15026 | included | normative baseline scope |
| 3.6.20.2 | Grammar | 15150 | included | normative baseline scope |
| 3.6.20.2.1 | Extended notation for use in Type definitions | 15154 | included | normative baseline scope |
| 3.6.20.2.2 | Extended notation for use in Template definitions | 15190 | included | normative baseline scope |
| 3.6.20.3 | Notes | 15272 | included | normative baseline scope |
| 3.6.21 | Event Filter definition | 15284 | included | normative baseline scope |
| 3.6.21.1 | Keynames | 15291 | included | normative baseline scope |
| 3.6.21.2 | Grammar | 15390 | included | normative baseline scope |
| 3.6.22 | Trigger definition | 15442 | included | normative baseline scope |
| 3.6.22.1 | Keynames | 15449 | included | normative baseline scope |
| 3.6.22.2 | Additional keynames for the extended condition notation | 15620 | included | normative baseline scope |
| 3.6.22.3 | Grammar | 15738 | included | normative baseline scope |
| 3.6.22.3.1 | Short notation | 15742 | included | normative baseline scope |
| 3.6.22.3.2 | Extended notation: | 15784 | included | normative baseline scope |
| 3.6.23 | Activity definitions | 15885 | included | normative baseline scope |
| 3.6.23.1 | Delegate workflow activity definition | 15936 | included | normative baseline scope |
| 3.6.23.1.1 | Keynames | 15938 | included | normative baseline scope |
| 3.6.23.1.2 | Grammar | 16058 | included | normative baseline scope |
| 3.6.23.1.2.1 | Short notation | 16063 | included | normative baseline scope |
| 3.6.23.1.2.2 | Extended notation | 16082 | included | normative baseline scope |
| 3.6.23.2 | Set state activity definition | 16121 | included | normative baseline scope |
| 3.6.23.2.1 | Keynames | 16125 | included | normative baseline scope |
| 3.6.23.2.2 | Grammar | 16184 | included | normative baseline scope |
| 3.6.23.3 | Call operation activity definition | 16211 | included | normative baseline scope |
| 3.6.23.3.1 | Keynames | 16216 | included | normative baseline scope |
| 3.6.23.3.2 | Grammar | 16339 | included | normative baseline scope |
| 3.6.23.3.2.1 | Short notation | 16344 | included | normative baseline scope |
| 3.6.23.3.2.2 | Extended notation | 16359 | included | normative baseline scope |
| 3.6.23.4 | Inline workflow activity definition | 16403 | included | normative baseline scope |
| 3.6.23.4.1 | Keynames | 16409 | included | normative baseline scope |
| 3.6.23.4.2 | Grammar | 16525 | included | normative baseline scope |
| 3.6.23.4.2.1 | Short notation | 16530 | included | normative baseline scope |
| 3.6.23.4.2.2 | Extended notation | 16545 | included | normative baseline scope |
| 3.6.23.5 | Example | 16583 | excluded | explicit Example section or child of one |
| 3.6.24 | Assertion definition | 16607 | included | normative baseline scope |
| 3.6.24.1 | Keynames | 16614 | included | normative baseline scope |
| 3.6.24.2 | Grammar | 16619 | included | normative baseline scope |
| 3.6.24.3 | Example | 16650 | excluded | explicit Example section or child of one |
| 3.6.25 | Condition clause definition | 16683 | included | normative baseline scope |
| 3.6.25.1 | Keynames | 16690 | included | normative baseline scope |
| 3.6.25.2 | Grammar | 16830 | included | normative baseline scope |
| 3.6.25.2.1 | And clause | 16834 | included | normative baseline scope |
| 3.6.25.2.2 | Or clause | 16856 | included | normative baseline scope |
| 3.6.25.2.3 | Not clause | 16878 | included | normative baseline scope |
| 3.6.25.3 | Direct assertion definition | 16901 | included | normative baseline scope |
| 3.6.25.4 | Additional Requirement | 16930 | included | normative baseline scope |
| 3.6.25.5 | Notes | 16937 | included | normative baseline scope |
| 3.6.25.6 | Example | 16944 | excluded | explicit Example section or child of one |
| 3.6.26 | Workflow precondition definition | 17103 | included | normative baseline scope |
| 3.6.26.1 | Keynames | 17111 | included | normative baseline scope |
| 3.6.26.2 | Grammar | 17211 | included | normative baseline scope |
| 3.6.27 | Workflow step definition | 17256 | included | normative baseline scope |
| 3.6.27.1 | Keynames | 17263 | included | normative baseline scope |
| 3.6.27.2 | Grammar | 17468 | included | normative baseline scope |
| 3.7 | Type-specific definitions | 17550 | included | normative baseline scope |
| 3.7.1 | Entity Type Schema | 17554 | included | normative baseline scope |
| 3.7.1.1 | Keynames | 17563 | included | normative baseline scope |
| 3.7.1.2 | Grammar | 17707 | included | normative baseline scope |
| 3.7.1.3 | Additional Requirements | 17755 | included | normative baseline scope |
| 3.7.2 | Capability definition | 17768 | included | normative baseline scope |
| 3.7.2.1 | Keynames | 17775 | included | normative baseline scope |
| 3.7.2.2 | Grammar | 17985 | included | normative baseline scope |
| 3.7.2.2.1 | Short notation | 17990 | included | normative baseline scope |
| 3.7.2.2.2 | Extended notation | 18010 | included | normative baseline scope |
| 3.7.2.3 | Examples | 18098 | excluded | explicit Example section or child of one |
| 3.7.2.3.1 | Simple notation example | 18103 | excluded | explicit Example section or child of one |
| 3.7.2.3.2 | Full notation example | 18118 | excluded | explicit Example section or child of one |
| 3.7.2.4 | Additional requirements | 18143 | included | normative baseline scope |
| 3.7.2.5 | Notes | 18157 | included | normative baseline scope |
| 3.7.3 | Requirement definition | 18174 | included | normative baseline scope |
| 3.7.3.1 | Keynames | 18186 | included | normative baseline scope |
| 3.7.3.1.1 | Additional Keynames for multi-line relationship grammar | 18338 | included | normative baseline scope |
| 3.7.3.2 | Grammar | 18437 | included | normative baseline scope |
| 3.7.3.2.1 | Simple grammar (Capability Type only) | 18442 | included | normative baseline scope |
| 3.7.3.2.2 | Extended grammar (with Node and Relationship Types) | 18456 | included | normative baseline scope |
| 3.7.3.2.3 | Extended grammar for declaring Property Definitions on the relationship’s Interfaces | 18483 | included | normative baseline scope |
| 3.7.3.3 | Additional Requirements | 18559 | included | normative baseline scope |
| 3.7.3.4 | Notes | 18577 | included | normative baseline scope |
| 3.7.3.5 | Requirement Type definition is a tuple | 18591 | included | normative baseline scope |
| 3.7.3.5.1 | Property filter | 18617 | included | normative baseline scope |
| 3.7.4 | Artifact Type | 18626 | included | normative baseline scope |
| 3.7.4.1 | Keynames | 18633 | included | normative baseline scope |
| 3.7.4.2 | Grammar | 18752 | included | normative baseline scope |
| 3.7.4.3 | Examples | 18848 | excluded | explicit Example section or child of one |
| 3.7.4.4 | Additional Requirements | 18882 | included | normative baseline scope |
| 3.7.4.5 | Notes | 18890 | included | normative baseline scope |
| 3.7.5 | Interface Type | 18915 | included | normative baseline scope |
| 3.7.5.1 | Keynames | 18922 | included | normative baseline scope |
| 3.7.5.2 | Grammar | 19023 | included | normative baseline scope |
| 3.7.5.3 | Example | 19116 | excluded | explicit Example section or child of one |
| 3.7.5.4 | Additional Requirements | 19149 | included | normative baseline scope |
| 3.7.5.5 | Notes | 19162 | included | normative baseline scope |
| 3.7.6 | Data Type | 19174 | included | normative baseline scope |
| 3.7.6.1 | Keynames | 19180 | included | normative baseline scope |
| 3.7.6.2 | Grammar | 19310 | included | normative baseline scope |
| 3.7.6.3 | Additional Requirements | 19413 | included | normative baseline scope |
| 3.7.6.4 | Examples | 19432 | excluded | explicit Example section or child of one |
| 3.7.6.4.1 | Defining a complex datatype | 19438 | excluded | explicit Example section or child of one |
| 3.7.6.4.2 | Defining a datatype derived from an existing datatype | 19466 | excluded | explicit Example section or child of one |
| 3.7.7 | Capability Type | 19494 | included | normative baseline scope |
| 3.7.7.1 | Keynames | 19502 | included | normative baseline scope |
| 3.7.7.2 | Grammar | 19605 | included | normative baseline scope |
| 3.7.7.3 | Example | 19692 | excluded | explicit Example section or child of one |
| 3.7.8 | Requirement Type | 19716 | included | normative baseline scope |
| 3.7.9 | Node Type | 19730 | included | normative baseline scope |
| 3.7.9.1 | Keynames | 19738 | included | normative baseline scope |
| 3.7.9.2 | Grammar | 19912 | included | normative baseline scope |
| 3.7.9.3 | Additional Requirements | 20040 | included | normative baseline scope |
| 3.7.9.4 | Best Practices | 20048 | included | normative baseline scope |
| 3.7.9.5 | Example | 20068 | excluded | explicit Example section or child of one |
| 3.7.10 | Relationship Type | 20117 | included | normative baseline scope |
| 3.7.10.1 | Keynames | 20123 | included | normative baseline scope |
| 3.7.10.2 | Grammar | 20249 | included | normative baseline scope |
| 3.7.10.3 | Best Practices | 20356 | included | normative baseline scope |
| 3.7.10.4 | Examples | 20373 | excluded | explicit Example section or child of one |
| 3.7.11 | Group Type | 20388 | included | normative baseline scope |
| 3.7.11.1 | Keynames | 20401 | included | normative baseline scope |
| 3.7.11.2 | Grammar | 20508 | included | normative baseline scope |
| 3.7.11.3 | Notes | 20602 | included | normative baseline scope |
| 3.7.11.4 | Additional Requirements | 20611 | included | normative baseline scope |
| 3.7.11.5 | Example | 20626 | excluded | explicit Example section or child of one |
| 3.7.12 | Policy Type | 20646 | included | normative baseline scope |
| 3.7.12.1 | Keynames | 20654 | included | normative baseline scope |
| 3.7.12.2 | Grammar | 20763 | included | normative baseline scope |
| 3.7.12.3 | Example | 20854 | excluded | explicit Example section or child of one |
| 3.8 | Template-specific definitions | 20874 | included | normative baseline scope |
| 3.8.1 | Capability assignment | 20882 | included | normative baseline scope |
| 3.8.1.1 | Keynames | 20889 | included | normative baseline scope |
| 3.8.1.2 | Grammar | 20991 | included | normative baseline scope |
| 3.8.1.3 | Example | 21054 | excluded | explicit Example section or child of one |
| 3.8.1.3.1 | Notation example | 21059 | excluded | explicit Example section or child of one |
| 3.8.2 | Requirement assignment | 21081 | included | normative baseline scope |
| 3.8.2.1 | Keynames | 21089 | included | normative baseline scope |
| 3.8.2.2 | Grammar | 21367 | included | normative baseline scope |
| 3.8.2.2.1 | Short notation: | 21372 | included | normative baseline scope |
| 3.8.2.2.2 | Extended notation: | 21399 | included | normative baseline scope |
| 3.8.2.2.3 | Extended grammar with Property Assignments for the relationship’s Interfaces | 21445 | included | normative baseline scope |
| 3.8.2.3 | min_occurrences, max_occurrences: lower and upper bounds of the range that further refines the minimum and maximum occurrences for this requirement specified in the corresponding requirement definition. The range specified here must fall completely within the occurrences range specified in the corresponding requirement definitionExamples | 21582 | included | normative baseline scope |
| 3.8.2.3.1 | Example 1 - Abstract hosting requirement on a Node Type | 21589 | excluded | explicit Example section or child of one |
| 3.8.2.3.2 | Example 2 - Requirement with Node Template and a custom Relationship Type | 21631 | excluded | explicit Example section or child of one |
| 3.8.2.3.3 | Example 3 - Requirement for a Compute node with additional selection criteria (filter) | 21667 | excluded | explicit Example section or child of one |
| 3.8.3 | Node Template | 21737 | included | normative baseline scope |
| 3.8.3.1 | Keynames | 21746 | included | normative baseline scope |
| 3.8.3.2 | Grammar | 22055 | included | normative baseline scope |
| 3.8.3.3 | Additional requirements | 22208 | included | normative baseline scope |
| 3.8.3.4 | Example | 22218 | excluded | explicit Example section or child of one |
| 3.8.4 | Relationship Template | 22248 | included | normative baseline scope |
| 3.8.4.1 | Keynames | 22258 | included | normative baseline scope |
| 3.8.4.2 | Grammar | 22449 | included | normative baseline scope |
| 3.8.4.3 | Additional requirements | 22547 | included | normative baseline scope |
| 3.8.4.4 | Example | 22557 | excluded | explicit Example section or child of one |
| 3.8.5 | Group definition | 22576 | included | normative baseline scope |
| 3.8.5.1 | Keynames | 22595 | included | normative baseline scope |
| 3.8.5.2 | Grammar | 22736 | included | normative baseline scope |
| 3.8.5.3 | Notes | 22819 | included | normative baseline scope |
| 3.8.5.4 | Additional Requirements | 22828 | included | normative baseline scope |
| 3.8.5.5 | Example | 22836 | excluded | explicit Example section or child of one |
| 3.8.6 | Policy definition | 22858 | included | normative baseline scope |
| 3.8.6.1 | Keynames | 22865 | included | normative baseline scope |
| 3.8.6.2 | Grammar | 23033 | included | normative baseline scope |
| 3.8.6.3 | Example | 23112 | excluded | explicit Example section or child of one |
| 3.8.7 | Imperative Workflow definition | 23136 | included | normative baseline scope |
| 3.8.7.1 | Keynames | 23146 | included | normative baseline scope |
| 3.8.7.2 | Grammar | 23340 | included | normative baseline scope |
| 3.8.8 | Property mapping | 23431 | included | normative baseline scope |
| 3.8.8.1 | Keynames | 23437 | included | normative baseline scope |
| 3.8.8.2 | Grammar | 23514 | included | normative baseline scope |
| 3.8.8.3 | Notes | 23555 | included | normative baseline scope |
| 3.8.8.4 | Additional constraints | 23594 | included | normative baseline scope |
| 3.8.9 | Attribute mapping | 23604 | included | normative baseline scope |
| 3.8.9.1 | Keynames | 23609 | included | normative baseline scope |
| 3.8.9.2 | Grammar | 23662 | included | normative baseline scope |
| 3.8.10 | Capability mapping | 23678 | included | normative baseline scope |
| 3.8.10.1 | Keynames | 23684 | included | normative baseline scope |
| 3.8.10.2 | Grammar | 23788 | included | normative baseline scope |
| 3.8.10.3 | Additional requirements | 23876 | included | normative baseline scope |
| 3.8.11 | Requirement mapping | 23884 | included | normative baseline scope |
| 3.8.11.1 | Keynames | 23890 | included | normative baseline scope |
| 3.8.11.2 | Grammar | 23994 | included | normative baseline scope |
| 3.8.11.3 | Additional requirements | 24082 | included | normative baseline scope |
| 3.8.12 | Interface mapping | 24090 | included | normative baseline scope |
| 3.8.12.1 | Grammar | 24096 | included | normative baseline scope |
| 3.8.12.2 | Notes | 24136 | included | normative baseline scope |
| 3.8.13 | Substitution mapping | 24153 | included | normative baseline scope |
| 3.8.13.1 | Keynames | 24159 | included | normative baseline scope |
| 3.8.13.2 | Grammar | 24346 | included | normative baseline scope |
| 3.8.13.3 | Examples | 24424 | excluded | explicit Example section or child of one |
| 3.8.13.4 | Additional requirements | 24428 | included | normative baseline scope |
| 3.8.13.5 | Notes | 24435 | included | normative baseline scope |
| 3.9 | Topology Template definition | 24442 | included | normative baseline scope |
| 3.9.1 | Keynames | 24454 | included | normative baseline scope |
| 3.9.2 | Grammar | 24697 | included | normative baseline scope |
| 3.9.2.1 | inputs | 24822 | included | normative baseline scope |
| 3.9.2.1.1 | Grammar | 24851 | included | normative baseline scope |
| 3.9.2.1.2 | Examples | 24869 | excluded | explicit Example section or child of one |
| 3.9.2.2 | node_templates | 24921 | included | normative baseline scope |
| 3.9.2.2.1 | grammar | 24927 | included | normative baseline scope |
| 3.9.2.2.2 | Example | 24948 | excluded | explicit Example section or child of one |
| 3.9.2.3 | relationship_templates | 24973 | included | normative baseline scope |
| 3.9.2.3.1 | Grammar | 24987 | included | normative baseline scope |
| 3.9.2.3.2 | Example | 25008 | excluded | explicit Example section or child of one |
| 3.9.2.4 | outputs | 25036 | included | normative baseline scope |
| 3.9.2.4.1 | Grammar | 25044 | included | normative baseline scope |
| 3.9.2.4.2 | Example | 25062 | excluded | explicit Example section or child of one |
| 3.9.2.5 | groups | 25084 | included | normative baseline scope |
| 3.9.2.5.1 | Grammar | 25090 | included | normative baseline scope |
| 3.9.2.5.2 | Example | 25113 | excluded | explicit Example section or child of one |
| 3.9.2.6 | policies | 25161 | included | normative baseline scope |
| 3.9.2.6.1 | Grammar | 25167 | included | normative baseline scope |
| 3.9.2.6.2 | Example | 25190 | excluded | explicit Example section or child of one |
| 3.9.2.7 | substitution_mapping | 25209 | included | normative baseline scope |
| 3.9.2.7.1 | requirement_mapping | 25213 | included | normative baseline scope |
| 3.9.2.7.2 | Example | 25270 | excluded | explicit Example section or child of one |
| 3.9.2.8 | Notes | 25373 | included | normative baseline scope |
| 3.10 | Service Template definition | 25418 | included | normative baseline scope |
| 3.10.1 | Keynames | 25429 | included | normative baseline scope |
| 3.10.1.1 | Metadata keynames | 25836 | included | normative baseline scope |
| 3.10.2 | Grammar | 25931 | included | normative baseline scope |
| 3.10.2.1 | Requirements | 26038 | included | normative baseline scope |
| 3.10.2.2 | Notes | 26053 | included | normative baseline scope |
| 3.10.3 | Top-level keyname definitions | 26062 | included | normative baseline scope |
| 3.10.3.1 | tosca_definitions_version | 26065 | included | normative baseline scope |
| 3.10.3.1.1 | Keyname | 26072 | included | normative baseline scope |
| 3.10.3.1.2 | Grammar | 26084 | included | normative baseline scope |
| 3.10.3.1.3 | Examples: | 26098 | excluded | explicit Example section or child of one |
| 3.10.3.2 | metadata | 26127 | included | normative baseline scope |
| 3.10.3.2.1 | Keyname | 26133 | included | normative baseline scope |
| 3.10.3.2.2 | Grammar | 26145 | included | normative baseline scope |
| 3.10.3.2.3 | Example | 26158 | excluded | explicit Example section or child of one |
| 3.10.3.3 | template_name | 26177 | included | normative baseline scope |
| 3.10.3.3.1 | Keyname | 26182 | included | normative baseline scope |
| 3.10.3.3.2 | Grammar | 26194 | included | normative baseline scope |
| 3.10.3.3.3 | Example | 26205 | excluded | explicit Example section or child of one |
| 3.10.3.3.4 | Notes | 26216 | included | normative baseline scope |
| 3.10.3.4 | template_author | 26225 | included | normative baseline scope |
| 3.10.3.4.1 | Keyname | 26230 | included | normative baseline scope |
| 3.10.3.4.2 | Grammar | 26242 | included | normative baseline scope |
| 3.10.3.4.3 | Example | 26254 | excluded | explicit Example section or child of one |
| 3.10.3.5 | template_version | 26266 | included | normative baseline scope |
| 3.10.3.5.1 | Keyname | 26271 | included | normative baseline scope |
| 3.10.3.5.2 | Grammar | 26283 | included | normative baseline scope |
| 3.10.3.5.3 | Example | 26297 | excluded | explicit Example section or child of one |
| 3.10.3.5.4 | Notes: | 26308 | included | normative baseline scope |
| 3.10.3.6 | description | 26319 | included | normative baseline scope |
| 3.10.3.6.1 | Keyname | 26325 | included | normative baseline scope |
| 3.10.3.7 | dsl_definitions | 26336 | included | normative baseline scope |
| 3.10.3.7.1 | Keyname | 26342 | included | normative baseline scope |
| 3.10.3.7.2 | Grammar | 26354 | included | normative baseline scope |
| 3.10.3.7.3 | Example | 26373 | excluded | explicit Example section or child of one |
| 3.10.3.8 | repositories | 26407 | included | normative baseline scope |
| 3.10.3.8.1 | Keyname | 26414 | included | normative baseline scope |
| 3.10.3.8.2 | Grammar | 26426 | included | normative baseline scope |
| 3.10.3.8.3 | Example | 26445 | excluded | explicit Example section or child of one |
| 3.10.3.9 | imports | 26465 | included | normative baseline scope |
| 3.10.3.9.1 | Keyname | 26475 | included | normative baseline scope |
| 3.10.3.9.2 | Grammar | 26487 | included | normative baseline scope |
| 3.10.3.9.3 | Example | 26508 | excluded | explicit Example section or child of one |
| 3.10.3.10 | artifact_types | 26533 | included | normative baseline scope |
| 3.10.3.10.1 | Keyname | 26538 | included | normative baseline scope |
| 3.10.3.10.2 | Grammar | 26550 | included | normative baseline scope |
| 3.10.3.10.3 | Example | 26569 | excluded | explicit Example section or child of one |
| 3.10.3.11 | data_types | 26586 | included | normative baseline scope |
| 3.10.3.11.1 | Keyname | 26591 | included | normative baseline scope |
| 3.10.3.11.2 | Grammar | 26603 | included | normative baseline scope |
| 3.10.3.11.3 | Example | 26622 | excluded | explicit Example section or child of one |
| 3.10.3.12 | capability_types | 26677 | included | normative baseline scope |
| 3.10.3.12.1 | Keyname | 26684 | included | normative baseline scope |
| 3.10.3.12.2 | Grammar | 26696 | included | normative baseline scope |
| 3.10.3.12.3 | Example | 26715 | excluded | explicit Example section or child of one |
| 3.10.3.13 | interface_types | 26745 | included | normative baseline scope |
| 3.10.3.13.1 | Keyname | 26751 | included | normative baseline scope |
| 3.10.3.13.2 | Grammar | 26763 | included | normative baseline scope |
| 3.10.3.13.3 | Example | 26782 | excluded | explicit Example section or child of one |
| 3.10.3.14 | relationship_types | 26805 | included | normative baseline scope |
| 3.10.3.14.1 | Keyname | 26811 | included | normative baseline scope |
| 3.10.3.14.2 | Grammar | 26823 | included | normative baseline scope |
| 3.10.3.14.3 | Example | 26842 | excluded | explicit Example section or child of one |
| 3.10.3.15 | node_types | 26873 | included | normative baseline scope |
| 3.10.3.15.1 | Keyname | 26879 | included | normative baseline scope |
| 3.10.3.15.2 | Grammar | 26891 | included | normative baseline scope |
| 3.10.3.15.3 | Example | 26910 | excluded | explicit Example section or child of one |
| 3.10.3.15.4 | Notes | 26941 | included | normative baseline scope |
| 3.10.3.16 | group_types | 26950 | included | normative baseline scope |
| 3.10.3.16.1 | Keyname | 26956 | included | normative baseline scope |
| 3.10.3.16.2 | Grammar | 26968 | included | normative baseline scope |
| 3.10.3.16.3 | Example | 26987 | excluded | explicit Example section or child of one |
| 3.10.3.17 | policy_types | 27004 | included | normative baseline scope |
| 3.10.3.17.1 | Keyname | 27009 | included | normative baseline scope |
| 3.10.3.17.2 | Grammar | 27021 | included | normative baseline scope |
| 3.10.3.17.3 | Example | 27040 | excluded | explicit Example section or child of one |
| 4 | TOSCA functions | 27059 | included | normative baseline scope |
| 4.1 | Reserved Function Keywords | 27071 | included | normative baseline scope |
| 4.2 | Environment Variable Conventions | 27191 | included | normative baseline scope |
| 4.2.1 | Reserved Environment Variable Names and Usage | 27196 | included | normative baseline scope |
| 4.2.2 | Prefixed vs. Unprefixed TARGET names | 27497 | included | normative baseline scope |
| 4.2.2.1 | Notes | 27505 | included | normative baseline scope |
| 4.3 | Intrinsic functions | 27519 | included | normative baseline scope |
| 4.3.1 | concat | 27527 | included | normative baseline scope |
| 4.3.1.1 | Grammar | 27533 | included | normative baseline scope |
| 4.3.1.2 | Parameters | 27546 | included | normative baseline scope |
| 4.3.1.3 | Examples | 27600 | excluded | explicit Example section or child of one |
| 4.3.2 | join | 27624 | included | normative baseline scope |
| 4.3.2.1 | Grammar | 27629 | included | normative baseline scope |
| 4.3.2.2 | Parameters | 27642 | included | normative baseline scope |
| 4.3.2.3 | Examples | 27719 | excluded | explicit Example section or child of one |
| 4.3.3 | token | 27744 | included | normative baseline scope |
| 4.3.3.1 | Grammar | 27750 | included | normative baseline scope |
| 4.3.3.2 | Parameters | 27764 | included | normative baseline scope |
| 4.3.3.3 | Examples | 27860 | excluded | explicit Example section or child of one |
| 4.4 | Property functions | 27882 | included | normative baseline scope |
| 4.4.1 | get_input | 27904 | included | normative baseline scope |
| 4.4.1.1 | Grammar | 27910 | included | normative baseline scope |
| 4.4.1.2 | Parameters | 27937 | included | normative baseline scope |
| 4.4.1.3 | Examples | 28017 | excluded | explicit Example section or child of one |
| 4.4.2 | get_property | 28142 | included | normative baseline scope |
| 4.4.2.1 | Grammar | 28148 | included | normative baseline scope |
| 4.4.2.2 | Parameters | 28164 | included | normative baseline scope |
| 4.4.2.3 | Examples | 28297 | excluded | explicit Example section or child of one |
| 4.5 | Attribute functions | 28421 | included | normative baseline scope |
| 4.5.1 | get_attribute | 28433 | included | normative baseline scope |
| 4.5.1.1 | Grammar | 28439 | included | normative baseline scope |
| 4.5.1.2 | Parameters | 28454 | included | normative baseline scope |
| 4.5.1.3 | Examples: | 28590 | excluded | explicit Example section or child of one |
| 4.5.1.4 | Notes | 28596 | included | normative baseline scope |
| 4.6 | Operation functions | 28612 | included | normative baseline scope |
| 4.6.1 | get_operation_output | 28622 | included | normative baseline scope |
| 4.6.1.1 | Grammar | 28629 | included | normative baseline scope |
| 4.6.1.2 | Parameters | 28643 | included | normative baseline scope |
| 4.6.1.3 | Notes | 28763 | included | normative baseline scope |
| 4.7 | Navigation functions | 28773 | included | normative baseline scope |
| 4.7.1 | get_nodes_of_type | 28783 | included | normative baseline scope |
| 4.7.1.1 | Grammar | 28789 | included | normative baseline scope |
| 4.7.1.2 | Parameters | 28802 | included | normative baseline scope |
| 4.7.1.3 | Returns | 28852 | included | normative baseline scope |
| 4.8 | Artifact functions | 28896 | included | normative baseline scope |
| 4.8.1 | get_artifact | 28901 | included | normative baseline scope |
| 4.8.1.1 | Grammar | 28907 | included | normative baseline scope |
| 4.8.1.2 | Parameters | 28921 | included | normative baseline scope |
| 4.8.1.3 | Examples | 29068 | excluded | explicit Example section or child of one |
| 4.8.1.3.1 | Example: Retrieving artifact without specified location | 29075 | excluded | explicit Example section or child of one |
| 4.8.1.3.2 | Example: Retrieving artifact as a local path | 29126 | excluded | explicit Example section or child of one |
| 4.8.1.3.3 | Example: Retrieving artifact in a specified location | 29173 | excluded | explicit Example section or child of one |
| 4.9 | Context-based Entity names (global) | 29220 | included | normative baseline scope |
| 4.9.1.1 | Goals | 29229 | included | normative baseline scope |
| 5 | TOSCA normative type definitions | 29239 | included | normative baseline scope |
| 5.1 | Assumptions | 29260 | included | normative baseline scope |
| 5.2 | TOSCA normative type names | 29281 | included | normative baseline scope |
| 5.2.1 | Additional requirements | 29328 | included | normative baseline scope |
| 5.3 | Data Types | 29343 | included | normative baseline scope |
| 5.3.1 | tosca.datatypes.Root | 29347 | included | normative baseline scope |
| 5.3.1.1 | Definition | 29353 | included | normative baseline scope |
| 5.3.2 | tosca.datatypes.json | 29369 | included | normative baseline scope |
| 5.3.2.1 | Definition | 29413 | included | normative baseline scope |
| 5.3.2.2 | Examples | 29428 | excluded | explicit Example section or child of one |
| 5.3.2.2.1 | Type declaration example | 29430 | excluded | explicit Example section or child of one |
| 5.3.2.2.2 | Template definition example | 29495 | excluded | explicit Example section or child of one |
| 5.3.3 | Additional Requirements | 29520 | included | normative baseline scope |
| 5.3.4 | tosca.datatypes.xml | 29530 | included | normative baseline scope |
| 5.3.4.1 | Definition | 29571 | included | normative baseline scope |
| 5.3.4.2 | Examples | 29586 | excluded | explicit Example section or child of one |
| 5.3.4.2.1 | Type declaration example | 29588 | excluded | explicit Example section or child of one |
| 5.3.4.2.2 | Template definition example | 29651 | excluded | explicit Example section or child of one |
| 5.3.5 | Additional Requirements | 29674 | included | normative baseline scope |
| 5.3.6 | tosca.datatypes.Credential | 29681 | included | normative baseline scope |
| 5.3.6.1 | Properties | 29723 | included | normative baseline scope |
| 5.3.6.2 | Definition | 29887 | included | normative baseline scope |
| 5.3.6.3 | Additional requirements | 29936 | included | normative baseline scope |
| 5.3.6.4 | Notes | 29944 | included | normative baseline scope |
| 5.3.6.5 | Examples | 29956 | excluded | explicit Example section or child of one |
| 5.3.6.5.1 | Provide a simple user name and password without a protocol or standardized token format | 29958 | excluded | explicit Example section or child of one |
| 5.3.6.5.2 | HTTP Basic access authentication credential | 29982 | excluded | explicit Example section or child of one |
| 5.3.6.5.3 | X-Auth-Token credential | 30008 | excluded | explicit Example section or child of one |
| 5.3.6.5.4 | OAuth bearer token credential | 30032 | excluded | explicit Example section or child of one |
| 5.3.6.6 | OpenStack SSH Keypair | 30055 | included | normative baseline scope |
| 5.3.7 | tosca.datatypes.TimeInterval | 30080 | included | normative baseline scope |
| 5.3.7.1 | Properties | 30123 | included | normative baseline scope |
| 5.3.7.2 | Definition | 30207 | included | normative baseline scope |
| 5.3.7.3 | Examples | 30236 | excluded | explicit Example section or child of one |
| 5.3.7.3.1 | Multi-day evaluation time period | 30238 | excluded | explicit Example section or child of one |
| 5.3.8 | tosca.datatypes.network.NetworkInfo | 30259 | included | normative baseline scope |
| 5.3.8.1 | Properties | 30301 | included | normative baseline scope |
| 5.3.8.2 | Definition | 30395 | included | normative baseline scope |
| 5.3.8.3 | Examples | 30433 | excluded | explicit Example section or child of one |
| 5.3.8.4 | Additional Requirements | 30457 | included | normative baseline scope |
| 5.3.9 | tosca.datatypes.network.PortInfo | 30470 | included | normative baseline scope |
| 5.3.9.1 | Properties | 30512 | included | normative baseline scope |
| 5.3.9.2 | Definition | 30648 | included | normative baseline scope |
| 5.3.9.3 | Examples | 30694 | excluded | explicit Example section or child of one |
| 5.3.9.4 | Additional Requirements | 30721 | included | normative baseline scope |
| 5.3.10 | tosca.datatypes.network.PortDef | 30734 | included | normative baseline scope |
| 5.3.10.1 | Definition | 30776 | included | normative baseline scope |
| 5.3.10.2 | Examples | 30795 | excluded | explicit Example section or child of one |
| 5.3.11 | tosca.datatypes.network.PortSpec | 30833 | included | normative baseline scope |
| 5.3.11.1 | Properties | 30875 | included | normative baseline scope |
| 5.3.11.2 | Definition | 31038 | included | normative baseline scope |
| 5.3.11.3 | Additional requirements | 31102 | included | normative baseline scope |
| 5.3.11.4 | Examples | 31128 | excluded | explicit Example section or child of one |
| 5.4 | Artifact Types | 31160 | included | normative baseline scope |
| 5.4.1 | tosca.artifacts.Root | 31206 | included | normative baseline scope |
| 5.4.1.1 | Definition | 31213 | included | normative baseline scope |
| 5.4.2 | tosca.artifacts.File | 31228 | included | normative baseline scope |
| 5.4.2.1 | Definition | 31271 | included | normative baseline scope |
| 5.4.3 | Deployment Types | 31285 | included | normative baseline scope |
| 5.4.3.1 | tosca.artifacts.Deployment | 31287 | included | normative baseline scope |
| 5.4.3.1.1 | Definition | 31296 | included | normative baseline scope |
| 5.4.3.2 | Additional Requirements | 31324 | included | normative baseline scope |
| 5.4.3.3 | tosca.artifacts.Deployment.Image | 31331 | included | normative baseline scope |
| 5.4.3.3.1 | Definition | 31374 | included | normative baseline scope |
| 5.4.3.4 | tosca.artifacts.Deployment.Image.VM | 31388 | included | normative baseline scope |
| 5.4.3.4.1 | Definition | 31397 | included | normative baseline scope |
| 5.4.3.4.2 | Notes | 31414 | included | normative baseline scope |
| 5.4.4 | Implementation Types | 31423 | included | normative baseline scope |
| 5.4.4.1 | tosca.artifacts.Implementation | 31426 | included | normative baseline scope |
| 5.4.4.1.1 | Definition | 31435 | included | normative baseline scope |
| 5.4.4.2 | Additional Requirements | 31463 | included | normative baseline scope |
| 5.4.4.3 | tosca.artifacts.Implementation.Bash | 31470 | included | normative baseline scope |
| 5.4.4.3.1 | Definition | 31512 | included | normative baseline scope |
| 5.4.4.4 | tosca.artifacts.Implementation.Python | 31532 | included | normative baseline scope |
| 5.4.4.4.1 | Definition | 31575 | included | normative baseline scope |
| 5.4.5 | Template Types | 31595 | included | normative baseline scope |
| 5.4.5.1 | tosca.artifacts.template | 31599 | included | normative baseline scope |
| 5.4.5.1.1 | Definition | 31617 | included | normative baseline scope |
| 5.5 | Capabilities Types | 31646 | included | normative baseline scope |
| 5.5.1 | tosca.capabilities.Root | 31648 | included | normative baseline scope |
| 5.5.1.1 | Definition | 31654 | included | normative baseline scope |
| 5.5.2 | tosca.capabilities.Node | 31668 | included | normative baseline scope |
| 5.5.2.1 | Definition | 31710 | included | normative baseline scope |
| 5.5.3 | tosca.capabilities.Compute | 31724 | included | normative baseline scope |
| 5.5.3.1 | Properties | 31767 | included | normative baseline scope |
| 5.5.3.2 | Definition | 31935 | included | normative baseline scope |
| 5.5.4 | tosca.capabilities.Network | 32000 | included | normative baseline scope |
| 5.5.4.1 | Properties | 32043 | included | normative baseline scope |
| 5.5.4.2 | Definition | 32100 | included | normative baseline scope |
| 5.5.5 | tosca.capabilities.Storage | 32123 | included | normative baseline scope |
| 5.5.5.1 | Properties | 32166 | included | normative baseline scope |
| 5.5.5.2 | Definition | 32223 | included | normative baseline scope |
| 5.5.6 | tosca.capabilities.Container | 32246 | included | normative baseline scope |
| 5.5.6.1 | Properties | 32289 | included | normative baseline scope |
| 5.5.6.2 | Definition | 32345 | included | normative baseline scope |
| 5.5.7 | tosca.capabilities.Endpoint | 32360 | included | normative baseline scope |
| 5.5.7.1 | Properties | 32407 | included | normative baseline scope |
| 5.5.7.2 | Attributes | 32676 | included | normative baseline scope |
| 5.5.7.3 | Definition | 32734 | included | normative baseline scope |
| 5.5.7.4 | Additional requirements | 32831 | included | normative baseline scope |
| 5.5.8 | tosca.capabilities.Endpoint.Public | 32840 | included | normative baseline scope |
| 5.5.8.1 | Definition | 32888 | included | normative baseline scope |
| 5.5.8.2 | Additional requirements | 32943 | included | normative baseline scope |
| 5.5.9 | tosca.capabilities.Endpoint.Admin | 32959 | included | normative baseline scope |
| 5.5.9.1 | Properties | 33001 | included | normative baseline scope |
| 5.5.9.2 | Definition | 33057 | included | normative baseline scope |
| 5.5.9.3 | Additional requirements | 33085 | included | normative baseline scope |
| 5.5.10 | tosca.capabilities.Endpoint.Database | 33093 | included | normative baseline scope |
| 5.5.10.1 | Properties | 33135 | included | normative baseline scope |
| 5.5.10.2 | Definition | 33191 | included | normative baseline scope |
| 5.5.11 | tosca.capabilities.Attachment | 33206 | included | normative baseline scope |
| 5.5.11.1 | Properties | 33250 | included | normative baseline scope |
| 5.5.11.2 | Definition | 33306 | included | normative baseline scope |
| 5.5.12 | tosca.capabilities.OperatingSystem | 33320 | included | normative baseline scope |
| 5.5.12.1 | Properties | 33362 | included | normative baseline scope |
| 5.5.12.2 | Definition | 33511 | included | normative baseline scope |
| 5.5.12.3 | Additional Requirements | 33555 | included | normative baseline scope |
| 5.5.13 | tosca.capabilities.Scalable | 33568 | included | normative baseline scope |
| 5.5.13.1 | Properties | 33610 | included | normative baseline scope |
| 5.5.13.2 | Definition | 33728 | included | normative baseline scope |
| 5.5.13.3 | Notes | 33762 | included | normative baseline scope |
| 5.5.14 | tosca.capabilities.network.Bindable | 33772 | included | normative baseline scope |
| 5.5.14.1 | Properties | 33815 | included | normative baseline scope |
| 5.5.14.2 | Definition | 33871 | included | normative baseline scope |
| 5.6 | Requirement Types | 33885 | included | normative baseline scope |
| 5.7 | Relationship Types | 33894 | included | normative baseline scope |
| 5.7.1 | tosca.relationships.Root | 33899 | included | normative baseline scope |
| 5.7.1.1 | Attributes | 33905 | included | normative baseline scope |
| 5.7.1.2 | Definition | 34021 | included | normative baseline scope |
| 5.7.2 | tosca.relationships.DependsOn | 34052 | included | normative baseline scope |
| 5.7.2.1 | Definition | 34094 | included | normative baseline scope |
| 5.7.3 | tosca.relationships.HostedOn | 34112 | included | normative baseline scope |
| 5.7.3.1 | Definition | 34154 | included | normative baseline scope |
| 5.7.4 | tosca.relationships.ConnectsTo | 34172 | included | normative baseline scope |
| 5.7.4.1 | Definition | 34214 | included | normative baseline scope |
| 5.7.4.2 | Properties | 34241 | included | normative baseline scope |
| 5.7.5 | tosca.relationships.AttachesTo | 34299 | included | normative baseline scope |
| 5.7.5.1 | Properties | 34342 | included | normative baseline scope |
| 5.7.5.2 | Attributes | 34434 | included | normative baseline scope |
| 5.7.5.3 | Definition | 34493 | included | normative baseline scope |
| 5.7.6 | tosca.relationships.RoutesTo | 34528 | included | normative baseline scope |
| 5.7.6.1 | Definition | 34570 | included | normative baseline scope |
| 5.8 | Interface Types | 34588 | included | normative baseline scope |
| 5.8.1 | Additional Requirements | 34599 | included | normative baseline scope |
| 5.8.2 | Best Practices | 34616 | included | normative baseline scope |
| 5.8.3 | tosca.interfaces.Root | 34624 | included | normative baseline scope |
| 5.8.3.1 | Definition | 34630 | included | normative baseline scope |
| 5.8.4 | tosca.interfaces.node.lifecycle.Standard | 34646 | included | normative baseline scope |
| 5.8.4.1 | Definition | 34688 | included | normative baseline scope |
| 5.8.4.2 | Create operation | 34720 | included | normative baseline scope |
| 5.8.4.3 | TOSCA Orchestrator processing of Deployment artifacts | 34730 | included | normative baseline scope |
| 5.8.4.4 | Operation sequencing and node state | 34747 | included | normative baseline scope |
| 5.8.4.4.1 | Normal node startup sequence diagram | 34760 | included | normative baseline scope |
| 5.8.4.4.2 | Normal node shutdown sequence diagram | 34769 | included | normative baseline scope |
| 5.8.5 | tosca.interfaces.relationship.Configure | 34777 | included | normative baseline scope |
| 5.8.5.1 | Definition | 34820 | included | normative baseline scope |
| 5.8.5.2 | Invocation Conventions | 34876 | included | normative baseline scope |
| 5.8.5.3 | Normal node start sequence with Configure relationship operations | 34899 | included | normative baseline scope |
| 5.8.5.4 | Node-Relationship configuration sequence | 34909 | included | normative baseline scope |
| 5.8.5.4.1 | Node-Relationship add, remove and changed sequence | 34936 | included | normative baseline scope |
| 5.8.5.5 | Notes | 34956 | included | normative baseline scope |
| 5.9 | Node Types | 34994 | included | normative baseline scope |
| 5.9.1 | tosca.nodes.Root | 34999 | included | normative baseline scope |
| 5.9.1.1 | Properties | 35047 | included | normative baseline scope |
| 5.9.1.2 | Attributes | 35104 | included | normative baseline scope |
| 5.9.1.3 | Definition | 35220 | included | normative baseline scope |
| 5.9.1.4 | Additional Requirements | 35284 | included | normative baseline scope |
| 5.9.2 | tosca.nodes.Abstract.Compute | 35293 | included | normative baseline scope |
| 5.9.2.1 | Properties | 35340 | included | normative baseline scope |
| 5.9.2.2 | Attributes | 35396 | included | normative baseline scope |
| 5.9.2.3 | Definition | 35452 | included | normative baseline scope |
| 5.9.3 | tosca.nodes.Compute | 35475 | included | normative baseline scope |
| 5.9.3.1 | Properties | 35519 | included | normative baseline scope |
| 5.9.3.2 | Attributes | 35575 | included | normative baseline scope |
| 5.9.3.3 | Definition | 35717 | included | normative baseline scope |
| 5.9.3.4 | Additional Requirements | 35805 | included | normative baseline scope |
| 5.9.4 | tosca.nodes.SoftwareComponent | 35815 | included | normative baseline scope |
| 5.9.4.1 | Properties | 35858 | included | normative baseline scope |
| 5.9.4.2 | Attributes | 35942 | included | normative baseline scope |
| 5.9.4.3 | Definition | 35998 | included | normative baseline scope |
| 5.9.4.4 | Additional Requirements | 36043 | included | normative baseline scope |
| 5.9.5 | tosca.nodes.WebServer | 36052 | included | normative baseline scope |
| 5.9.5.1 | Properties | 36096 | included | normative baseline scope |
| 5.9.5.2 | Definition | 36152 | included | normative baseline scope |
| 5.9.5.3 | Additional Requirements | 36186 | included | normative baseline scope |
| 5.9.6 | tosca.nodes.WebApplication | 36196 | included | normative baseline scope |
| 5.9.6.1 | Properties | 36240 | included | normative baseline scope |
| 5.9.6.2 | Definition | 36298 | included | normative baseline scope |
| 5.9.7 | tosca.nodes.DBMS | 36339 | included | normative baseline scope |
| 5.9.7.1 | Properties | 36347 | included | normative baseline scope |
| 5.9.7.2 | Definition | 36430 | included | normative baseline scope |
| 5.9.8 | tosca.nodes.Database | 36474 | included | normative baseline scope |
| 5.9.8.1 | Properties | 36517 | included | normative baseline scope |
| 5.9.8.2 | Definition | 36654 | included | normative baseline scope |
| 5.9.9 | tosca.nodes.Abstract.Storage | 36716 | included | normative baseline scope |
| 5.9.9.1 | Properties | 36759 | included | normative baseline scope |
| 5.9.9.2 | Definition | 36844 | included | normative baseline scope |
| 5.9.10 | tosca.nodes.Storage.ObjectStorage | 36880 | included | normative baseline scope |
| 5.9.10.1 | Properties | 36922 | included | normative baseline scope |
| 5.9.10.2 | Definition | 36980 | included | normative baseline scope |
| 5.9.10.3 | Notes: | 37010 | included | normative baseline scope |
| 5.9.11 | tosca.nodes.Storage.BlockStorage | 37021 | included | normative baseline scope |
| 5.9.11.1 | Properties | 37070 | included | normative baseline scope |
| 5.9.11.2 | Attributes | 37198 | included | normative baseline scope |
| 5.9.11.3 | Definition | 37254 | included | normative baseline scope |
| 5.9.11.4 | Additional Requirements | 37290 | included | normative baseline scope |
| 5.9.11.5 | Notes | 37300 | included | normative baseline scope |
| 5.9.12 | tosca.nodes.Container.Runtime | 37318 | included | normative baseline scope |
| 5.9.12.1 | Definition | 37363 | included | normative baseline scope |
| 5.9.13 | tosca.nodes.Container.Application | 37393 | included | normative baseline scope |
| 5.9.13.1 | Definition | 37436 | included | normative baseline scope |
| 5.9.14 | tosca.nodes.LoadBalancer | 37473 | included | normative baseline scope |
| 5.9.14.1 | Definition | 37526 | included | normative baseline scope |
| 5.9.14.2 | Notes: | 37576 | included | normative baseline scope |
| 5.10 | Group Types | 37585 | included | normative baseline scope |
| 5.10.1 | tosca.groups.Root | 37602 | included | normative baseline scope |
| 5.10.1.1 | Definition | 37608 | included | normative baseline scope |
| 5.10.1.2 | Notes: | 37628 | included | normative baseline scope |
| 5.11 | Policy Types | 37641 | included | normative baseline scope |
| 5.11.1 | tosca.policies.Root | 37652 | included | normative baseline scope |
| 5.11.1.1 | Definition | 37658 | included | normative baseline scope |
| 5.11.2 | tosca.policies.Placement | 37671 | included | normative baseline scope |
| 5.11.2.1 | Definition | 37677 | included | normative baseline scope |
| 5.11.3 | tosca.policies.Scaling | 37694 | included | normative baseline scope |
| 5.11.3.1 | Definition | 37700 | included | normative baseline scope |
| 5.11.4 | tosca.policies.Update | 37717 | included | normative baseline scope |
| 5.11.4.1 | Definition | 37723 | included | normative baseline scope |
| 5.11.5 | tosca.policies.Performance | 37740 | included | normative baseline scope |
| 5.11.5.1 | Definition | 37747 | included | normative baseline scope |
| 6 | TOSCA Cloud Service Archive (CSAR) format | 37768 | included | normative baseline scope |
| 6.1 | Overall Structure of a CSAR | 37791 | included | normative baseline scope |
| 6.2 | TOSCA Meta File | 37818 | included | normative baseline scope |
| 6.2.1 | Custom keynames in the TOSCA.meta file | 37881 | included | normative baseline scope |
| 6.2.2 | Example | 37909 | excluded | explicit Example section or child of one |
| 6.3 | Archive without TOSCA-Metadata | 37949 | included | normative baseline scope |
| 6.3.1 | Example | 37969 | excluded | explicit Example section or child of one |
| 7 | TOSCA workflows | 37998 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.1 | Normative workflows | 38028 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.1.1 | Notes | 38046 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2 | Declarative workflows | 38062 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.1 | Notes | 38079 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.1.1 | Orchestrator provided nodes lifecycle and weaving | 38084 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2 | Relationship impacts on topology weaving | 38101 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2.1 | tosca.relationships.DependsOn | 38107 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2.2 | Note | 38113 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2.2.1 | Example DependsOn | 38120 | excluded | explicit Example section or child of one |
| 7.2.2.3 | tosca.relationships.ConnectsTo | 38134 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2.4 | tosca.relationships.HostedOn | 38145 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.2.4.1 | Example Software Component HostedOn Compute | 38167 | excluded | explicit Example section or child of one |
| 7.2.2.4.2 | Example Software Component HostedOn Software Component | 38184 | excluded | explicit Example section or child of one |
| 7.2.2.4.3 | Example 2 Software Components HostedOn Compute | 38198 | excluded | explicit Example section or child of one |
| 7.2.3 | Limitations | 38203 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.3.1 | Hosted nodes concurrency | 38205 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.3.2 | Dependent nodes concurrency | 38215 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.2.3.3 | Target operations and get_attribute on source | 38225 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3 | Imperative workflows | 38238 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.1 | Defining sequence of operations in an imperative workflow | 38248 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.1.1 | Using on_success to define steps ordering | 38263 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.1.1.1 | Example | 38293 | excluded | explicit Example section or child of one |
| 7.3.1.2 | Define a sequence of activity on the same element | 38358 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.2 | Definition of a simple workflow | 38453 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.2.1 | Example: Adding a non-normative custom workflow | 38467 | excluded | explicit Example section or child of one |
| 7.3.2.2 | Example: Creating two nodes hosted on the same compute in parallel | 38527 | excluded | explicit Example section or child of one |
| 7.3.3 | Specifying preconditions to a workflow | 38669 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.3.1 | Example : adding precondition to custom backup workflow | 38676 | excluded | explicit Example section or child of one |
| 7.3.4 | Workflow reusability | 38753 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.4.1 | Reusing a workflow to build multiple workflows | 38759 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.4.2 | Inlining a complex workflow | 38893 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.5 | Defining conditional logic on some part of the workflow | 39013 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.6 | Define inputs for a workflow | 39110 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.6.1 | Example | 39120 | excluded | explicit Example section or child of one |
| 7.3.7 | Handle operation failure | 39199 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 7.3.7.1 | Example | 39221 | excluded | explicit Example section or child of one |
| 8 | TOSCA networking | 39371 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.1 | Networking and Service Template Portability | 39389 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.2 | Connectivity semantics | 39447 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3 | Expressing connectivity semantics | 39500 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3.1 | Connection initiation semantics | 39509 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3.1.1 | Source to Target | 39532 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3.1.2 | Target to Source | 39560 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3.1.3 | Peer-to-Peer | 39587 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.3.2 | Specifying layer 4 ports | 39618 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4 | Network provisioning | 39637 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.1 | Declarative network provisioning | 39642 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.2 | Implicit network fulfillment | 39656 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.3 | Controlling network fulfillment | 39668 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.3.1 | Use case: OAM Network | 39692 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.3.2 | Use case: Data Traffic network | 39713 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.4.3.3 | Use case: Bring my own DHCP | 39722 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5 | Network Types | 39738 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.1 | tosca.nodes.network.Network | 39743 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.1.1 | Properties | 39784 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.1.2 | Attributes | 40168 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.1.3 | Definition | 40225 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.2 | tosca.nodes.network.Port | 40324 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.2.1 | Properties | 40369 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.2.2 | Attributes | 40556 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.2.3 | Definition | 40613 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.3 | tosca.capabilities.network.Linkable | 40689 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.3.1 | Properties | 40733 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.3.2 | Definition | 40789 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.4 | tosca.relationships.network.LinksTo | 40803 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.4.1 | Definition | 40845 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.5 | tosca.relationships.network.BindsTo | 40864 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.5.5.1 | Definition | 40906 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.6 | Network modeling approaches | 40925 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.6.1 | Option 1: Specifying a network outside the application’s Service Template | 40930 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 8.6.2 | Option 2: Specifying network requirements within the application’s Service Template | 41131 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 9 | Non-normative type definitions | 41225 | excluded | section is explicitly non-normative |
| 9.1 | Artifact Types | 41238 | excluded | section is explicitly non-normative |
| 9.1.1 | tosca.artifacts.Deployment.Image.Container.Docker | 41253 | excluded | section is explicitly non-normative |
| 9.1.1.1 | Definition | 41260 | excluded | section is explicitly non-normative |
| 9.1.2 | tosca.artifacts.Deployment.Image.VM.ISO | 41276 | excluded | section is explicitly non-normative |
| 9.1.2.1 | Definition | 41282 | excluded | section is explicitly non-normative |
| 9.1.3 | tosca.artifacts.Deployment.Image.VM.QCOW2 | 41303 | excluded | section is explicitly non-normative |
| 9.1.3.1 | Definition | 41309 | excluded | section is explicitly non-normative |
| 9.1.4 | tosca.artifacts.template.Jinja2 | 41329 | excluded | section is explicitly non-normative |
| 9.1.4.1 | Definition | 41374 | excluded | section is explicitly non-normative |
| 9.1.4.2 | Example | 41391 | excluded | explicit Example section or child of one |
| 9.1.5 | tosca.artifacts.template.Twig | 41432 | excluded | section is explicitly non-normative |
| 9.1.5.1 | Definition | 41475 | excluded | section is explicitly non-normative |
| 9.2 | Capability Types | 41492 | excluded | section is explicitly non-normative |
| 9.2.1 | tosca.capabilities.Container.Docker | 41497 | excluded | section is explicitly non-normative |
| 9.2.1.1 | Properties | 41538 | excluded | section is explicitly non-normative |
| 9.2.1.2 | Definition | 41767 | excluded | section is explicitly non-normative |
| 9.2.1.3 | Notes | 41828 | excluded | section is explicitly non-normative |
| 9.3 | Node Types | 41845 | excluded | section is explicitly non-normative |
| 9.3.1 | tosca.nodes.Database.MySQL | 41855 | excluded | section is explicitly non-normative |
| 9.3.1.1 | Properties | 41858 | excluded | section is explicitly non-normative |
| 9.3.1.2 | Definition | 41914 | excluded | section is explicitly non-normative |
| 9.3.2 | tosca.nodes.DBMS.MySQL | 41935 | excluded | section is explicitly non-normative |
| 9.3.2.1 | Properties | 41938 | excluded | section is explicitly non-normative |
| 9.3.2.2 | Definition | 41994 | excluded | section is explicitly non-normative |
| 9.3.3 | tosca.nodes.WebServer.Apache | 42034 | excluded | section is explicitly non-normative |
| 9.3.3.1 | Properties | 42036 | excluded | section is explicitly non-normative |
| 9.3.3.2 | Definition | 42092 | excluded | section is explicitly non-normative |
| 9.3.4 | tosca.nodes.WebApplication.WordPress | 42106 | excluded | section is explicitly non-normative |
| 9.3.4.1 | Properties | 42111 | excluded | section is explicitly non-normative |
| 9.3.4.2 | Definition | 42167 | excluded | section is explicitly non-normative |
| 9.3.5 | tosca.nodes.WebServer.Nodejs | 42211 | excluded | section is explicitly non-normative |
| 9.3.5.1 | Properties | 42216 | excluded | section is explicitly non-normative |
| 9.3.5.2 | Definition | 42272 | excluded | section is explicitly non-normative |
| 9.3.6 | tosca.nodes.Container.Application.Docker | 42310 | excluded | section is explicitly non-normative |
| 9.3.6.1 | Properties | 42312 | excluded | section is explicitly non-normative |
| 9.3.6.2 | Definition | 42368 | excluded | section is explicitly non-normative |
| 10 | Component Modeling Use Cases | 42393 | excluded | component modeling use cases |
| 10.1.1 | Use Case: Exploring the HostedOn relationship using WebApplication and WebServer | 42404 | excluded | component modeling use cases |
| 10.1.1.1 | WebServer declares its “host” capability | 42413 | excluded | component modeling use cases |
| 10.1.1.2 | WebApplication declares its “host” requirement | 42459 | excluded | component modeling use cases |
| 10.1.1.2.1 | Notes | 42499 | excluded | component modeling use cases |
| 10.1.2 | Use Case: Establishing a ConnectsTo relationship to WebServer | 42510 | excluded | component modeling use cases |
| 10.1.2.1 | Best practice | 42597 | excluded | component modeling use cases |
| 10.1.3 | Use Case: Attaching (local) BlockStorage to a Compute node | 42604 | excluded | component modeling use cases |
| 10.1.4 | Use Case: Reusing a BlockStorage Relationship using Relationship Type or Relationship Template | 42666 | excluded | component modeling use cases |
| 10.1.4.1 | Simple Profile Rationale | 42680 | excluded | component modeling use cases |
| 10.1.4.2 | Notation Style #1: Augment AttachesTo Relationship Type directly in each Node Template | 42691 | excluded | component modeling use cases |
| 10.1.4.3 | Notation Style #2: Use the ‘template’ keyword on the Node Templates to specify which named Relationship Template to use | 42785 | excluded | component modeling use cases |
| 10.1.4.4 | Notation Style #3: Using the “copy” keyname to define a similar Relationship Template | 42876 | excluded | component modeling use cases |
| 11 | Application Modeling Use Cases | 42996 | excluded | application modeling use cases |
| 11.1 | Use cases | 43010 | excluded | application modeling use cases |
| 11.1.1 | Overview | 43021 | excluded | application modeling use cases |
| 11.1.2 | Compute: Create a single Compute instance with a host Operating System | 43386 | excluded | application modeling use cases |
| 11.1.2.1 | Description | 43389 | excluded | application modeling use cases |
| 11.1.2.2 | Features | 43401 | excluded | application modeling use cases |
| 11.1.2.3 | Logical Diagram | 43442 | excluded | application modeling use cases |
| 11.1.2.4 | Sample YAML | 43447 | excluded | application modeling use cases |
| 11.1.2.5 | Notes | 43521 | excluded | application modeling use cases |
| 11.1.3 | Software Component 1: Automatic deployment of a Virtual Machine (VM) image artifact | 43528 | excluded | application modeling use cases |
| 11.1.3.1 | Description | 43533 | excluded | application modeling use cases |
| 11.1.3.2 | Features | 43544 | excluded | application modeling use cases |
| 11.1.3.3 | Assumptions | 43569 | excluded | application modeling use cases |
| 11.1.3.4 | Logical Diagram | 43587 | excluded | application modeling use cases |
| 11.1.3.5 | Sample YAML | 43592 | excluded | application modeling use cases |
| 11.1.3.6 | Notes | 43673 | excluded | application modeling use cases |
| 11.1.4 | Block Storage 1: Using the normative AttachesTo Relationship Type | 43691 | excluded | application modeling use cases |
| 11.1.4.1 | Description | 43694 | excluded | application modeling use cases |
| 11.1.4.2 | Logical Diagram | 43702 | excluded | application modeling use cases |
| 11.1.4.3 | Sample YAML | 43707 | excluded | application modeling use cases |
| 11.1.5 | Block Storage 2: Using a custom AttachesTo Relationship Type | 43838 | excluded | application modeling use cases |
| 11.1.5.1 | Description | 43842 | excluded | application modeling use cases |
| 11.1.5.2 | Logical Diagram | 43851 | excluded | application modeling use cases |
| 11.1.5.3 | Sample YAML | 43856 | excluded | application modeling use cases |
| 11.1.6 | Block Storage 3: Using a Relationship Template of type AttachesTo | 43994 | excluded | application modeling use cases |
| 11.1.6.1 | Description | 43998 | excluded | application modeling use cases |
| 11.1.6.2 | Logical Diagram | 44007 | excluded | application modeling use cases |
| 11.1.6.3 | Sample YAML | 44012 | excluded | application modeling use cases |
| 11.1.7 | Block Storage 4: Single Block Storage shared by 2-Tier Application with custom AttachesTo Type and implied relationships | 44139 | excluded | application modeling use cases |
| 11.1.7.1 | Description | 44144 | excluded | application modeling use cases |
| 11.1.7.2 | Logical Diagram | 44159 | excluded | application modeling use cases |
| 11.1.7.3 | Sample YAML | 44164 | excluded | application modeling use cases |
| 11.1.8 | Block Storage 5: Single Block Storage shared by 2-Tier Application with custom AttachesTo Type and explicit Relationship Templates | 44342 | excluded | application modeling use cases |
| 11.1.8.1 | Description | 44347 | excluded | application modeling use cases |
| 11.1.8.2 | Logical Diagram | 44361 | excluded | application modeling use cases |
| 11.1.8.3 | Sample YAML | 44366 | excluded | application modeling use cases |
| 11.1.9 | Block Storage 6: Multiple Block Storage attached to different Servers | 44564 | excluded | application modeling use cases |
| 11.1.9.1 | Description | 44568 | excluded | application modeling use cases |
| 11.1.9.2 | Logical Diagram | 44577 | excluded | application modeling use cases |
| 11.1.9.3 | Sample YAML | 44582 | excluded | application modeling use cases |
| 11.1.10 | Object Storage 1: Creating an Object Storage service | 44768 | excluded | application modeling use cases |
| 11.1.10.1 | Description | 44771 | excluded | application modeling use cases |
| 11.1.10.2 | Logical Diagram | 44773 | excluded | application modeling use cases |
| 11.1.10.3 | Sample YAML | 44778 | excluded | application modeling use cases |
| 11.1.11 | Network 1: Server bound to a new network | 44818 | excluded | application modeling use cases |
| 11.1.11.1 | Description | 44822 | excluded | application modeling use cases |
| 11.1.11.2 | Logical Diagram | 44835 | excluded | application modeling use cases |
| 11.1.11.3 | Sample YAML | 44840 | excluded | application modeling use cases |
| 11.1.12 | Network 2: Server bound to an existing network | 44928 | excluded | application modeling use cases |
| 11.1.12.1 | Description | 44931 | excluded | application modeling use cases |
| 11.1.12.2 | Logical Diagram | 44938 | excluded | application modeling use cases |
| 11.1.12.3 | Sample YAML | 44943 | excluded | application modeling use cases |
| 11.1.13 | Network 3: Two servers bound to a single network | 45025 | excluded | application modeling use cases |
| 11.1.13.1 | Description | 45028 | excluded | application modeling use cases |
| 11.1.13.2 | Logical Diagram | 45035 | excluded | application modeling use cases |
| 11.1.13.3 | Sample YAML | 45040 | excluded | application modeling use cases |
| 11.1.14 | Network 4: Server bound to three networks | 45191 | excluded | application modeling use cases |
| 11.1.14.1 | Description | 45194 | excluded | application modeling use cases |
| 11.1.14.2 | Logical Diagram | 45200 | excluded | application modeling use cases |
| 11.1.14.3 | Sample YAML | 45205 | excluded | application modeling use cases |
| 11.1.15 | WebServer-DBMS 1: WordPress + MySQL, single instance | 45331 | excluded | application modeling use cases |
| 11.1.15.1 | Description | 45335 | excluded | application modeling use cases |
| 11.1.15.2 | Logical Diagram | 45340 | excluded | application modeling use cases |
| 11.1.15.3 | Sample YAML | 45345 | excluded | application modeling use cases |
| 11.1.15.4 | Sample scripts | 45572 | excluded | application modeling use cases |
| 11.1.15.4.1 | wordpress_install.sh | 45577 | excluded | application modeling use cases |
| 11.1.15.4.2 | wordpress_configure.sh | 45588 | excluded | application modeling use cases |
| 11.1.15.4.3 | mysql_database_configure.sh | 45611 | excluded | application modeling use cases |
| 11.1.15.4.4 | mysql_dbms_install.sh | 45635 | excluded | application modeling use cases |
| 11.1.15.4.5 | mysql_dbms_start.sh | 45650 | excluded | application modeling use cases |
| 11.1.15.4.6 | mysql_dbms_configure | 45665 | excluded | application modeling use cases |
| 11.1.15.4.7 | webserver_install.sh | 45679 | excluded | application modeling use cases |
| 11.1.15.4.8 | webserver_start.sh | 45692 | excluded | application modeling use cases |
| 11.1.16 | WebServer-DBMS 2: Nodejs with PayPal Sample App and MongoDB on separate instances | 45706 | excluded | application modeling use cases |
| 11.1.16.1 | Description | 45712 | excluded | application modeling use cases |
| 11.1.16.2 | Logical Diagram | 45718 | excluded | application modeling use cases |
| 11.1.16.3 | Sample YAML | 45723 | excluded | application modeling use cases |
| 11.1.16.4 | Notes: | 45925 | excluded | application modeling use cases |
| 11.1.17 | Multi-Tier-1: Elasticsearch, Logstash, Kibana (ELK) use case with multiple instances | 45934 | excluded | application modeling use cases |
| 11.1.17.1 | Description | 45937 | excluded | application modeling use cases |
| 11.1.17.2 | Logical Diagram | 45961 | excluded | application modeling use cases |
| 11.1.17.3 | Sample YAML | 45966 | excluded | application modeling use cases |
| 11.1.17.3.1 | Master Service Template application (Entry-Definitions) | 45968 | excluded | application modeling use cases |
| 11.1.17.4 | Sample scripts | 46395 | excluded | application modeling use cases |
| 11.1.18 | Container-1: Containers using Docker single Compute instance (Containers only) | 46400 | excluded | application modeling use cases |
| 11.1.18.1 | Description | 46407 | excluded | application modeling use cases |
| 11.1.18.2 | Logical Diagram | 46430 | excluded | application modeling use cases |
| 11.1.18.3 | Sample YAML | 46435 | excluded | application modeling use cases |
| 11.1.18.3.1 | Two Docker “Container” nodes (Only) with Docker Requirements | 46437 | excluded | application modeling use cases |
| 11.1.19 | Artifacts: Compute Node with Multiple Artifacts | 46550 | excluded | application modeling use cases |
| 11.1.19.1 | Description | 46552 | excluded | application modeling use cases |
| 11.1.19.2 | Sample YAML | 46558 | excluded | application modeling use cases |
| 12 | TOSCA Policies | 46639 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.1 | A declarative approach | 46652 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.1.1 | Declarative considerations | 46671 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.2 | Consideration of Event, Condition and Action | 46707 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3 | Types of policies | 46712 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3.1 | Access control policies | 46742 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3.2 | Placement policies | 46751 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3.2.1 | Placement for governance concerns | 46771 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3.2.2 | Placement for failover | 46788 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.3.3 | Quality-of-Service (QoS) policies | 46802 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.4 | Policy relationship considerations | 46819 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 12.5 | Use Cases | 46847 | excluded | policy use cases |
| 12.5.1 | Placement | 46859 | excluded | policy use cases |
| 12.5.1.1 | Use Case 1: Simple placement for failover | 46861 | excluded | policy use cases |
| 12.5.1.1.1 | Description | 46863 | excluded | policy use cases |
| 12.5.1.1.2 | Features | 46870 | excluded | policy use cases |
| 12.5.1.1.3 | Logical Diagram | 46889 | excluded | policy use cases |
| 12.5.1.1.4 | Notes | 46915 | excluded | policy use cases |
| 12.5.1.2 | Use Case 2: Controlled placement by region | 46927 | excluded | policy use cases |
| 12.5.1.2.1 | Description | 46929 | excluded | policy use cases |
| 12.5.1.2.2 | Features | 46948 | excluded | policy use cases |
| 12.5.1.2.3 | Sample YAML: Region separation amongst named set of regions | 46957 | excluded | policy use cases |
| 12.5.1.3 | Use Case 3: Co-locate based upon Compute affinity | 46986 | excluded | policy use cases |
| 12.5.1.3.1 | Description | 46988 | excluded | policy use cases |
| 12.5.1.3.2 | Features | 47000 | excluded | policy use cases |
| 12.5.1.4 | Notes | 47008 | excluded | policy use cases |
| 12.5.1.4.1 | Sample YAML: Region separation amongst named set of regions | 47019 | excluded | policy use cases |
| 12.5.2 | Scaling | 47039 | excluded | policy use cases |
| 12.5.2.1 | Use Case 1: Simple node autoscale | 47041 | excluded | policy use cases |
| 12.5.2.1.1 | Description | 47043 | excluded | policy use cases |
| 12.5.2.1.2 | Features | 47048 | excluded | policy use cases |
| 12.5.2.1.3 | Sample YAML | 47056 | excluded | policy use cases |
| 12.5.2.1.4 | Notes | 47082 | excluded | policy use cases |
| 13 | Artifact Processing and creating Portable Service Templates | 47125 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.1 | CSAR Onboarding | 47182 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.1.1 | How is on-boarding done | 47206 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.1.2 | Artifact Distribution | 47241 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.1.3 | Why is on-boarding needed | 47268 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2 | Artifacts Processing | 47274 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.1 | Identify Artifact Processor | 47294 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.2 | Establish an Execution Environment | 47334 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.3 | Configure Artifact Processor User Account | 47391 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.4 | Deploy Artifact Processor | 47407 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.5 | Deploy Dependencies | 47438 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.6 | Identify Target | 47491 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.7 | Pass Inputs and Retrieve Results or Errors | 47526 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.2.8 | Cleanup | 47542 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.3 | Dynamic Artifacts | 47566 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 13.4 | Discussion of Examples | 47572 | excluded | explicit Example section or child of one |
| 13.4.1 | Shell Scripts | 47610 | excluded | explicit Example section or child of one |
| 13.4.1.1 | Progression of Examples | 47779 | excluded | explicit Example section or child of one |
| 13.4.1.1.1 | Simple install script that can run on all flavors for Unix. | 47786 | excluded | explicit Example section or child of one |
| 13.4.1.1.1.1 | Notes | 47797 | excluded | explicit Example section or child of one |
| 13.4.1.1.1.2 | Variants | 47815 | excluded | explicit Example section or child of one |
| 13.4.1.1.2 | Script that needs to be run as specific user | 47850 | excluded | explicit Example section or child of one |
| 13.4.1.1.3 | Simple script with dependencies | 47854 | excluded | explicit Example section or child of one |
| 13.4.1.1.4 | Different scripts for different Linux flavors | 47871 | excluded | explicit Example section or child of one |
| 13.4.1.1.4.1 | Variants | 47891 | excluded | explicit Example section or child of one |
| 13.4.1.1.5 | Scripts with environment variables | 47900 | excluded | explicit Example section or child of one |
| 13.4.1.1.6 | Scripts that require certain configuration files | 47915 | excluded | explicit Example section or child of one |
| 13.4.2 | Python Scripts | 47929 | excluded | explicit Example section or child of one |
| 13.4.2.1 | Python Scripts Executed in Orchestrator | 47937 | excluded | explicit Example section or child of one |
| 13.4.2.2 | Python Scripts Executed in Topology | 47942 | excluded | explicit Example section or child of one |
| 13.4.2.3 | Specifying Python Version | 47956 | excluded | explicit Example section or child of one |
| 13.4.2.3.1.1 | Assumptions/Questions | 47966 | excluded | explicit Example section or child of one |
| 13.4.2.4 | Deploying Dependencies | 47974 | excluded | explicit Example section or child of one |
| 13.4.2.4.1.1 | Assumptions/Questions | 47993 | excluded | explicit Example section or child of one |
| 13.4.2.4.1.2 | Notes | 48012 | excluded | explicit Example section or child of one |
| 13.4.3 | Package Artifacts | 48020 | excluded | explicit Example section or child of one |
| 13.4.3.1 | RPM Packages | 48044 | excluded | explicit Example section or child of one |
| 13.4.4 | Debian Packages | 48051 | excluded | explicit Example section or child of one |
| 13.4.4.1.1.1 | Notes | 48060 | excluded | explicit Example section or child of one |
| 13.4.4.2 | Distro-Independent Service Templates | 48074 | excluded | explicit Example section or child of one |
| 13.4.4.2.1.1 | Assumptions/Questions | 48097 | excluded | explicit Example section or child of one |
| 13.4.5 | VM Images | 48104 | excluded | explicit Example section or child of one |
| 13.4.5.1.1.1 | Premises | 48106 | excluded | explicit Example section or child of one |
| 13.4.5.1.1.2 | Notes | 48113 | excluded | explicit Example section or child of one |
| 13.4.5.1.1.3 | Assumptions/Questions | 48123 | excluded | explicit Example section or child of one |
| 13.4.5.2 | Image Onboarding - Uploading image to image repository | 48202 | excluded | explicit Example section or child of one |
| 13.4.6 | Container Images | 48231 | excluded | explicit Example section or child of one |
| 13.4.7 | API Artifacts | 48233 | excluded | explicit Example section or child of one |
| 13.4.7.1 | Examples | 48278 | excluded | explicit Example section or child of one |
| 13.4.8 | Non-Standard Artifacts with Execution Wrappers | 48305 | excluded | explicit Example section or child of one |
| 13.5 | Artifact Types and Metadata | 48309 | included-orchestrator-only | normative/runtime material kept outside processor conformance |
| 14 | Conformance | 48335 | included | normative baseline scope |
| 14.1 | Conformance Targets | 48385 | included | normative baseline scope |
| 14.2 | Conformance Clause 1: TOSCA YAML service template | 48414 | included | normative baseline scope |
| 14.3 | Conformance Clause 2: TOSCA processor | 48439 | included | normative baseline scope |
| 14.4 | Conformance Clause 3: TOSCA orchestrator | 48473 | included | normative baseline scope |
| 14.5 | Conformance Clause 4: TOSCA generator | 48516 | included | normative baseline scope |
| 14.6 | Conformance Clause 5: TOSCA archive | 48534 | included | normative baseline scope |

## Conformance clause 14.3 trace

All non-empty prose/list blocks in 14.3 have requirement records. The catalog separately inventories section 3 definitions and additional requirements, section 4 functions, section 5 normative types, imports, error rules, and the ambiguous string-normalization citation.

## Project CSAR scope

A read-only scope check found CSAR parsing, metadata validation, entry-definition selection, and archive creation under `tosca/csar` and `executables/puccini-csar`. Therefore all explicit section 6 requirements carry the separate `archive` and `csar` targets. Project behavior was not used as normative evidence.

## Normative type-definition coverage

All 61 `tosca.*` normative type headings in section 5 are represented. Five additional network type headings from section 8.5 are retained only as a separate orchestrator inventory and do not contribute to processor conformance.

## Sections requiring manual classification or external clarification

The source sections are structurally classified, but the issues listed below cannot be resolved from the permitted source alone:

- TOSCA13-AMB-001 (14.3, 3.1, 3.2, 3.3, 3.6): Conformance clause 14.3(d) cites the wrong section number.
- TOSCA13-AMB-002 (14.3, 5.4.9.3): Conformance clause 14.3(e) cites a non-existent section.
- TOSCA13-AMB-004 (3.7.3.1.1): Requirement relationship type is both required and optional.
- TOSCA13-AMB-005 (3.7.4.1): Artifact type properties conflict between Required column and prose.
- TOSCA13-AMB-007 (3.10.3.1): Version selector short name and URI forms.
- TOSCA13-AMB-008 (6.2): TOSCA.meta syntax is incorporated from an unavailable normative source.
- TOSCA13-AMB-009 (6.1, 6.3): CSAR root alternatives do not state whether they are exclusive.
- TOSCA13-AMB-010 (6.3, 3.10.1.1): Archive-without-metadata derives CSAR version from template version.
- TOSCA13-AMB-011 (6.2): Other-Definitions filename tokenization is incomplete.
- TOSCA13-AMB-012 (5.9.10.1, 5.9.10.3): ObjectStorage `maxsize` lower bound conflicts.
- TOSCA13-AMB-013 (5.9.11.1, 5.9.11.3, 5.9.11.4): BlockStorage `size` has conditional requiredness and `volume_id` precedence.
- TOSCA13-AMB-014 (3.7.1.1): Entity Type Schema derived_from constraint is malformed.

The catalog is complete as an extraction baseline: no structural coverage check is outstanding. “Complete” does not mean that the source ambiguities above have been resolved or that Puccini conforms.
