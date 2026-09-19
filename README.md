# probe

Pass it an artifact. Get a structured result and a meaningful exit code. `0` means valid. `1` means fix it. `2` means something broke.

```bash
go install github.com/Formulary-Labs/probe/cmd/probe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/Formulary-Labs/probe/releases) for `linux/amd64`, `darwin/arm64`, `darwin/amd64`, or `windows/amd64`.

## What it does

`probe` validates a gemara artifact file against the gemara schema using the same typed loaders as `go-gemara`. It detects the artifact type automatically, runs schema validation, and returns a structured result.

Exit codes are meaningful:

| Code | Meaning |
|---|---|
| `0` | All artifacts valid |
| `1` | One or more artifacts failed validation |
| `2` | Tool error — unreadable file, missing argument, or internal failure |

This makes `probe` usable as a CI gate: `probe artifact.yaml && proceed` is a complete smoke check.

## Usage

```bash
probe [flags] <artifact.yaml> [<artifact2.yaml> ...]
```

### Flags

| Flag | Default | Description |
|---|---|---|
| `--format` | `json` | Output format: `json` or `md` |
| `--program` | `""` | Program slug for provenance logging |
| `--quiet` | `false` | Suppress per-artifact output; exit code only |
| `--version` | — | Print version and exit |

### Examples

```bash
# Validate a single artifact
probe catalog.yaml

# Validate multiple artifacts
probe catalog.yaml guidance.yaml risk-register.yaml

# Markdown output — readable in terminal or CI log
probe --format md *.yaml

# Gate a CI step on artifact validity
probe --program iso42001 catalog.yaml && echo "valid"

# Exit code only — for scripts that handle output themselves
probe --quiet artifact.yaml
echo $?
```

## Output

### JSON

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

### Markdown

```
# probe results

**2/2 valid**

| Path | Status | Type | Duration |
|---|---|---|---|
| `catalog.yaml` | valid | ControlCatalog | 12ms |
| `guidance.yaml` | valid | GuidanceCatalog | 8ms |
```

## Supported artifact types

| Type | `metadata.type` value |
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

Type detection is automatic — `probe` reads the `metadata.type` field and selects the appropriate loader.

## CI integration

Add a `probe` validation step before any pipeline step that consumes a gemara artifact.

### GitHub Actions

```yaml
- name: Validate gemara artifacts
  run: probe --program ${{ env.PROGRAM }} catalog.yaml

- name: Validate all YAML artifacts in directory
  run: probe artifacts/*.yaml
```

### General CI pattern

```bash
# Validate before running assay
probe catalog.yaml || exit 1
assay --framework iso27001 --catalog catalog.yaml --product-source docs/
```

## Relationship to gemara-mcp

`probe` and [`gemara-mcp`](https://github.com/gemaraproj/gemara-mcp) use the same underlying validation logic from [`go-gemara`](https://github.com/gemaraproj/go-gemara). They serve different callsites: `probe` is for terminals and CI pipelines where clean exit codes matter; `gemara-mcp`'s `validate_gemara_artifact` tool is for AI agents working in context. Use `probe` in automated pipelines, `gemara-mcp` in agent sessions.

## License

Apache License 2.0
