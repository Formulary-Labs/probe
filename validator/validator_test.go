package validator_test

import (
	"testing"

	"github.com/Formulary-Labs/probe/validator"
)

func TestValidate_validControlCatalog(t *testing.T) {
	r := validator.Validate("../testdata/good-control-catalog.yaml")
	if r.Status != validator.Valid {
		t.Errorf("expected Valid, got %s — errors: %v", r.Status, r.Errors)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
	if r.ArtifactType == "" {
		t.Error("ArtifactType must not be empty for a valid artifact")
	}
	if r.DurationMs < 0 {
		t.Errorf("DurationMs must be non-negative, got %d", r.DurationMs)
	}
}

func TestValidate_validEvaluationLog(t *testing.T) {
	r := validator.Validate("../testdata/good-evaluation-log.yaml")
	if r.Status != validator.Valid {
		t.Errorf("expected Valid, got %s — errors: %v", r.Status, r.Errors)
	}
}

func TestValidate_validLexicon(t *testing.T) {
	r := validator.Validate("../testdata/good-lexicon.yaml")
	if r.Status != validator.Valid {
		t.Errorf("expected Valid, got %s — errors: %v", r.Status, r.Errors)
	}
}

func TestValidate_invalid(t *testing.T) {
	r := validator.Validate("../testdata/bad.yaml")
	if r.Status == validator.Valid {
		t.Error("expected Invalid or Unreadable for bad.yaml, got Valid")
	}
}

func TestValidate_missing(t *testing.T) {
	r := validator.Validate("../testdata/does-not-exist.yaml")
	if r.Status != validator.Unreadable {
		t.Errorf("expected Unreadable for missing file, got %s", r.Status)
	}
	if len(r.Errors) == 0 {
		t.Error("expected errors for missing file")
	}
}

func TestValidateAll_multipleFiles(t *testing.T) {
	paths := []string{
		"../testdata/good-control-catalog.yaml",
		"../testdata/good-evaluation-log.yaml",
	}
	results := validator.ValidateAll(paths)
	if len(results) != len(paths) {
		t.Fatalf("expected %d results, got %d", len(paths), len(results))
	}
	for i, r := range results {
		if r.Status != validator.Valid {
			t.Errorf("results[%d] expected Valid, got %s: %v", i, r.Status, r.Errors)
		}
	}
}

func TestAnyInvalid_false(t *testing.T) {
	results := []validator.Result{
		{Status: validator.Valid},
		{Status: validator.Valid},
	}
	if validator.AnyInvalid(results) {
		t.Error("AnyInvalid should be false when all results are Valid")
	}
}

func TestAnyInvalid_true(t *testing.T) {
	results := []validator.Result{
		{Status: validator.Valid},
		{Status: validator.Invalid},
	}
	if !validator.AnyInvalid(results) {
		t.Error("AnyInvalid should be true when any result is Invalid")
	}
}
