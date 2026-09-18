# probe

Validate gemara artifacts from the command line.

```
go install github.com/Formulary-Labs/probe/cmd/probe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/Formulary-Labs/probe/releases).

---

## Usage

```
probe [flags] <artifact.yaml> [<artifact2.yaml> ...]
```

Exits **0** if all artifacts are valid.  
Exits **1** if any artifact fails validation.  
Exits **2** on tool errors (unreadable file, missing argument, etc.).

### Flags

| Flag | Default | Description |
|---|---|---|
| `--format` | `json` | Output format: `json` or `md` |
| `--program` | `""` | Program slug for provenance logging |
| `--quiet` | `false` | Suppress per-artifact output; only exit code |
| `--version` | — | Print version and exit |

### Examples

```bash
# Validate a single artifact
probe artifact.yaml

# Validate multiple artifacts
probe catalog.yaml guidance.yaml risk-register.yaml

# Markdown output
probe --format md *.yaml

# Validate in a CI script and fail the step on invalid artifacts
probe --program iso42001 artifact.yaml && echo "valid"

# Quiet mode — exit code only (for scripts)
probe --quiet artifact.yaml
echo $?  # 0 = valid, 1 = invalid
```

### GitHub Actions

```yaml
- name: Validate gemara artifacts
  run: probe artifact.yaml

# Or validate all YAML files in the artifacts directory
- name: Validate all gemara artifacts
  run: probe artifacts/*.yaml
```

---

## Output

### JSON (default)

```json
{
  "total": 2,
  "valid": 2,
  "invalid": 0,
  "results": [
    {
      "path": "catalog.yaml",
      "status": "valid",
      "artifact_type": "ControlCatalog",
      "errors": [],
      "duration_ms": 12
    }
  ]
}
```

### Markdown (`--format md`)

```
# probe results

**2/2 valid**

| Path | Status | Type | Duration |
|---|---|---|---|
| `catalog.yaml` | ✓ valid | ControlCatalog | 12ms |
| `guidance.yaml` | ✓ valid | GuidanceCatalog | 8ms |
```

---

## Supported artifact types

| Type | gemara metadata.type value |
|---|---|
| Control Catalog | `ControlCatalog` |
| Guidance Catalog | `GuidanceCatalog` |
| Audit Log | `AuditLog` |
| Evaluation Log | `EvaluationLog` |
| Capability Catalog | `CapabilityCatalog` |
| Enforcement Log | `EnforcementLog` |
| Lexicon | `Lexicon` |
| Mapping Document | `MappingDocument` |
| Policy | `Policy` |
| Principle Catalog | `PrincipleCatalog` |
| Risk Catalog | `RiskCatalog` |
| Threat Catalog | `ThreatCatalog` |
| Vector Catalog | `VectorCatalog` |

---

## Relationship to gemara-mcp

`probe` is for terminals and CI pipelines — it's a binary you run with clean exit codes.  
[`gemara-mcp`](https://github.com/gemaraproj/gemara-mcp)'s `validate_gemara_artifact` tool is for AI agents in context.  
Same underlying validation logic ([go-gemara](https://github.com/gemaraproj/go-gemara)), different callsite.

---

## Part of Formulary

`probe` is part of the [Formulary](https://github.com/Formulary-Labs) compliance micro-tools ecosystem.  
See [CONTRIBUTING.md](https://github.com/Formulary-Labs/.github/blob/main/CONTRIBUTING.md) to contribute.

---

## License

Apache License 2.0
