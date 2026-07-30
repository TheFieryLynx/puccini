# TOSCA 1.3 normative profile discrepancies

This report compares the OASIS Standard prose, the pinned OASIS community
machine-readable profile, and Puccini's effective bundled TOSCA 1.3 profile.
The prose Standard is the conformance authority.

- Pinned upstream commit: `e5b0a3ee46488921ff409bbe0feba37af2a3233e`
- Types in the union: **66**
- Field discrepancies: **26**

## Classification summary

| Classification | Fields |
|---|---:|
| Puccini defect | 8 |
| community-profile defect | 2 |
| extension | 3 |
| obsolete declaration | 2 |
| prose specification ambiguity | 11 |

## Pairwise type summary

| Pair | Exact | Equivalent | Conflicting | Missing |
|---|---:|---:|---:|---:|
| specification_vs_oasis | 56 | 2 | 8 | 0 |
| specification_vs_puccini | 53 | 3 | 10 | 0 |
| oasis_vs_puccini | 50 | 4 | 12 | 0 |

## Field discrepancies

### PROFILE13-001: `tosca.artifacts.Deployment.Image.VM` `derived_from`

- Specification section: §5.4.3.4
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `tosca.artifacts.Deployment.Image`
- OASIS community: `tosca.artifacts.Deployment.Image`
- Puccini: `tosca.artifacts.Deployment`

### PROFILE13-002: `tosca.capabilities.Scalable` `properties.default_instances.default`

- Specification section: §5.5.13
- Classification: **community-profile defect**
- Resolution: **confirmed**
- Specification: `None`
- OASIS community: `1`
- Puccini: `None`

### PROFILE13-003: `tosca.datatypes.json` `constraints`

- Specification section: §5.3.2
- Classification: **extension**
- Resolution: **interpreted**
- Specification: `None`
- OASIS community: `None`
- Puccini: `[{'_format': 'json'}]`

### PROFILE13-004: `tosca.datatypes.network.NetworkInfo` `properties.addresses.required`

- Specification section: §5.3.8
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-005: `tosca.datatypes.network.NetworkInfo` `properties.network_id.required`

- Specification section: §5.3.8
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-006: `tosca.datatypes.network.NetworkInfo` `properties.network_name.required`

- Specification section: §5.3.8
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-007: `tosca.datatypes.network.PortInfo` `properties.addresses.required`

- Specification section: §5.3.9
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-008: `tosca.datatypes.network.PortInfo` `properties.mac_address.required`

- Specification section: §5.3.9
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-009: `tosca.datatypes.network.PortInfo` `properties.network_id.required`

- Specification section: §5.3.9
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-010: `tosca.datatypes.network.PortInfo` `properties.port_id.required`

- Specification section: §5.3.9
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-011: `tosca.datatypes.network.PortInfo` `properties.port_name.required`

- Specification section: §5.3.9
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `False`
- OASIS community: `True`
- Puccini: `False`
- Decision/ambiguity record: `docs/decisions/tosca-1.3-network-info-requiredness.md`

### PROFILE13-012: `tosca.datatypes.xml` `constraints`

- Specification section: §5.3.4
- Classification: **extension**
- Resolution: **interpreted**
- Specification: `None`
- OASIS community: `None`
- Puccini: `[{'_format': 'xml'}]`

### PROFILE13-013: `tosca.groups.Root` `interfaces`

- Specification section: §5.10.1
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `{'Standard': {'type': 'tosca.interfaces.node.lifecycle.Standard'}}`
- OASIS community: `None`
- Puccini: `None`
- Secondary classification: **community-profile defect**

### PROFILE13-014: `tosca.interfaces.Root` `derived_from`

- Specification section: §5.8.3
- Classification: **obsolete declaration**
- Resolution: **interpreted**
- Specification: `tosca.entity.Root`
- OASIS community: `None`
- Puccini: `None`

### PROFILE13-015: `tosca.interfaces.relationship.Configure` `operations.remove_source`

- Specification section: §5.8.5
- Classification: **extension**
- Resolution: **interpreted**
- Specification: `None`
- OASIS community: `None`
- Puccini: `{}`

### PROFILE13-016: `tosca.nodes.Abstract.Compute` `valid_source_types`

- Specification section: §5.9.2
- Classification: **obsolete declaration**
- Resolution: **interpreted**
- Specification: `[]`
- OASIS community: `None`
- Puccini: `None`

### PROFILE13-017: `tosca.nodes.Abstract.Storage` `derived_from`

- Specification section: §5.9.9
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `tosca.nodes.Root`
- OASIS community: `tosca.nodes.Root`
- Puccini: `None`

### PROFILE13-018: `tosca.nodes.Abstract.Storage` `properties.size.constraints`

- Specification section: §5.9.9
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `[{'greater_or_equal': '0 MB'}]`
- OASIS community: `[{'greater_or_equal': '0 MB'}]`
- Puccini: `[{'greater_or_equal': '0 GB'}]`

### PROFILE13-019: `tosca.nodes.Abstract.Storage` `properties.size.default`

- Specification section: §5.9.9
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `0 MB`
- OASIS community: `0 MB`
- Puccini: `None`

### PROFILE13-020: `tosca.nodes.Abstract.Storage` `properties.size.required`

- Specification section: §5.9.9
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `True`
- OASIS community: `True`
- Puccini: `False`

### PROFILE13-021: `tosca.nodes.Container.Application` `requirements.network.capability`

- Specification section: §5.9.13
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `tosca.capabilities.Endpoint`
- OASIS community: `tosca.capabilities.Endpoint`
- Puccini: `tosca.capabilities.Network`

### PROFILE13-022: `tosca.nodes.Container.Runtime` `capabilities.host.type`

- Specification section: §5.9.12
- Classification: **Puccini defect**
- Resolution: **confirmed**
- Specification: `tosca.capabilities.Compute`
- OASIS community: `tosca.capabilities.Compute`
- Puccini: `tosca.capabilities.Container`

### PROFILE13-023: `tosca.nodes.Storage.BlockStorage` `properties.size`

- Specification section: §5.9.11
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `None`
- OASIS community: `{'default': '1 MB', 'constraints': [{'greater_or_equal': '1 MB'}], 'required': True}`
- Puccini: `None`
- Decision/ambiguity record: `docs/conformance/tosca-1.3/ambiguities.md#tosca13-amb-013-blockstorage-size-has-conditional-requiredness-and-volume_id-precedence`

### PROFILE13-024: `tosca.nodes.Storage.ObjectStorage` `properties.maxsize.constraints.greater_or_equal`

- Specification section: §5.9.10
- Classification: **prose specification ambiguity**
- Resolution: **documented**
- Specification: `{'properties_table': '1 GB', 'definition': '0 GB'}`
- OASIS community: `0 GB`
- Puccini: `0 GB`
- Decision/ambiguity record: `docs/conformance/tosca-1.3/ambiguities.md#tosca13-amb-012-objectstorage-maxsize-lower-bound-conflicts`

### PROFILE13-025: `tosca.nodes.network.Port` `properties.order.required`

- Specification section: §8.5.2
- Classification: **prose specification ambiguity**
- Resolution: **interpreted**
- Specification: `{'properties_table': False, 'definition': True}`
- OASIS community: `True`
- Puccini: `True`
- Decision/ambiguity record: `docs/decisions/0008-tosca-1.3-network-port-order-requiredness.md`

### PROFILE13-026: `tosca.relationships.AttachesTo` `attributes`

- Specification section: §5.7.5
- Classification: **community-profile defect**
- Resolution: **confirmed**
- Specification: `{'device': {'type': 'string'}}`
- OASIS community: `None`
- Puccini: `{'device': {'type': 'string'}}`

## Effective-definition method

The audit canonicalizes capability shorthand, property requiredness defaults,
requirement occurrence defaults, and then recursively overlays inherited
definitions. Descriptions and implementation metadata are intentionally excluded
from semantic equality. Source placement remains recorded, including Puccini's
implicit `json` and `xml` declarations.
