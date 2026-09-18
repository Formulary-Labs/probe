// Package validator implements the core gemara artifact validation logic
// for probe. It validates a single artifact file against the gemara schema
// using go-gemara's typed loaders, returning structured results.
package validator

import (
	"context"
	"fmt"
	"os"
	"time"

	gemara "github.com/gemaraproj/go-gemara"
	"github.com/gemaraproj/go-gemara/fetcher"
)

// Status is the validation outcome for a single artifact.
type Status string

const (
	// Valid means the artifact passed all validation checks.
	Valid Status = "valid"
	// Invalid means the artifact failed one or more validation checks.
	Invalid Status = "invalid"
	// Unreadable means the artifact file could not be read or parsed.
	Unreadable Status = "unreadable"
)

// Result is the structured output for a single artifact validation run.
type Result struct {
	// Path is the path to the artifact that was validated.
	Path string `json:"path"`
	// Status is the validation outcome.
	Status Status `json:"status"`
	// ArtifactType is the detected gemara artifact type, if identifiable.
	ArtifactType string `json:"artifact_type,omitempty"`
	// Errors contains any validation errors. Empty when Status is Valid.
	Errors []ValidationError `json:"errors,omitempty"`
	// DurationMs is the time taken to validate this artifact in milliseconds.
	DurationMs int64 `json:"duration_ms"`
}

// ValidationError is a single validation error from an artifact.
type ValidationError struct {
	// Message is the human-readable error description.
	Message string `json:"message"`
	// Field is the field path where the error occurred, if applicable.
	Field string `json:"field,omitempty"`
}

// Validate validates the gemara artifact at the given path.
// It detects the artifact type, loads it with the appropriate typed loader,
// and returns a structured Result.
func Validate(path string) Result {
	start := time.Now()

	data, err := os.ReadFile(path)
	if err != nil {
		return Result{
			Path:       path,
			Status:     Unreadable,
			Errors:     []ValidationError{{Message: fmt.Sprintf("cannot read file: %v", err)}},
			DurationMs: time.Since(start).Milliseconds(),
		}
	}

	artifactType, err := gemara.DetectType(data)
	if err != nil {
		return Result{
			Path:       path,
			Status:     Unreadable,
			Errors:     []ValidationError{{Message: fmt.Sprintf("cannot detect artifact type: %v", err)}},
			DurationMs: time.Since(start).Milliseconds(),
		}
	}

	typeName := artifactType.String()
	f := &fetcher.File{}
	ctx := context.Background()

	var validationErrs []ValidationError
	var loadErr error

	switch artifactType {
	case gemara.AuditLogArtifact:
		_, loadErr = gemara.Load[gemara.AuditLog](ctx, f, path)
	case gemara.CapabilityCatalogArtifact:
		_, loadErr = gemara.Load[gemara.CapabilityCatalog](ctx, f, path)
	case gemara.ControlCatalogArtifact:
		_, loadErr = gemara.Load[gemara.ControlCatalog](ctx, f, path)
	case gemara.EnforcementLogArtifact:
		_, loadErr = gemara.Load[gemara.EnforcementLog](ctx, f, path)
	case gemara.EvaluationLogArtifact:
		_, loadErr = gemara.Load[gemara.EvaluationLog](ctx, f, path)
	case gemara.GuidanceCatalogArtifact:
		_, loadErr = gemara.Load[gemara.GuidanceCatalog](ctx, f, path)
	case gemara.LexiconArtifact:
		_, loadErr = gemara.Load[gemara.Lexicon](ctx, f, path)
	case gemara.MappingDocumentArtifact:
		_, loadErr = gemara.Load[gemara.MappingDocument](ctx, f, path)
	case gemara.PolicyArtifact:
		_, loadErr = gemara.Load[gemara.Policy](ctx, f, path)
	case gemara.PrincipleCatalogArtifact:
		_, loadErr = gemara.Load[gemara.PrincipleCatalog](ctx, f, path)
	case gemara.RiskCatalogArtifact:
		_, loadErr = gemara.Load[gemara.RiskCatalog](ctx, f, path)
	case gemara.ThreatCatalogArtifact:
		_, loadErr = gemara.Load[gemara.ThreatCatalog](ctx, f, path)
	case gemara.VectorCatalogArtifact:
		_, loadErr = gemara.Load[gemara.VectorCatalog](ctx, f, path)
	default:
		validationErrs = append(validationErrs, ValidationError{
			Message: fmt.Sprintf("unknown or unsupported artifact type %q", typeName),
		})
	}

	if loadErr != nil {
		validationErrs = append(validationErrs, ValidationError{Message: loadErr.Error()})
	}

	status := Valid
	if len(validationErrs) > 0 {
		status = Invalid
	}

	return Result{
		Path:         path,
		Status:       status,
		ArtifactType: typeName,
		Errors:       validationErrs,
		DurationMs:   time.Since(start).Milliseconds(),
	}
}

// ValidateAll validates multiple artifact files and returns all results.
// The returned slice has the same length and order as paths.
func ValidateAll(paths []string) []Result {
	results := make([]Result, len(paths))
	for i, p := range paths {
		results[i] = Validate(p)
	}
	return results
}

// AnyInvalid returns true if any result has a non-Valid status.
func AnyInvalid(results []Result) bool {
	for _, r := range results {
		if r.Status != Valid {
			return true
		}
	}
	return false
}
