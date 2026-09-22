// probe validates gemara artifacts from the command line.
//
// Usage:
//
//	probe [flags] <artifact.yaml> [<artifact2.yaml> ...]
//
// Exits 0 if all artifacts are valid.
// Exits 1 if any artifact fails validation.
// Exits 2 on tool errors (unreadable file, missing argument, etc.).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Formulary-Labs/probe/validator"
	"github.com/Formulary-Labs/substrate/exit"
	"github.com/Formulary-Labs/substrate/format"
	"github.com/Formulary-Labs/substrate/provenance"
)

const version = "0.1.0"

func main() {
	var (
		fmtFlag     = flag.String("format", "json", "Output format: json (default), md")
		programFlag = flag.String("program", "", "Program slug for provenance logging")
		quietFlag   = flag.Bool("quiet", false, "Suppress per-artifact output; only emit summary and exit code")
		versionFlag = flag.Bool("version", false, "Print version and exit")
	)
	flag.Usage = usage
	flag.Parse()

	if *versionFlag {
		fmt.Fprintf(os.Stdout, "probe version %s\n", version)
		os.Exit(exit.OK)
	}

	paths := flag.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, `{"error": "no artifact paths provided", "code": 2}`)
		flag.Usage()
		os.Exit(exit.ToolError)
	}

	f, err := format.Parse(*fmtFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"error": %q, "code": 2}`+"\n", err.Error())
		os.Exit(exit.ToolError)
	}

	results := validator.ValidateAll(paths)
	anyFailed := validator.AnyInvalid(results)

	switch f {
	case format.JSON:
		printJSON(results, !*quietFlag)
	case format.MD:
		printMD(results, !*quietFlag)
	default:
		printJSON(results, !*quietFlag)
	}

	if *programFlag != "" {
		qg := provenance.Pass
		if anyFailed {
			qg = provenance.FailedOnceCorrected
		}
		_ = provenance.Write("logs/provenance.jsonl", provenance.Entry{
			Spec:        "functions/probe-spec.md",
			Output:      paths[0],
			OutputType:  "artifact",
			Program:     *programFlag,
			Purpose:     fmt.Sprintf("probe validated %d artifact(s)", len(paths)),
			Reusability: provenance.Instance,
			QualityGate: qg,
			Tool:        "probe",
			ToolVersion: version,
		})
	}

	if anyFailed {
		os.Exit(exit.Validation)
	}
	os.Exit(exit.OK)
}

// Summary is the top-level JSON output structure.
type Summary struct {
	Total   int                `json:"total"`
	Valid   int                `json:"valid"`
	Invalid int                `json:"invalid"`
	Results []validator.Result `json:"results"`
}

func printJSON(results []validator.Result, _ bool) {
	validCount := 0
	for _, r := range results {
		if r.Status == validator.Valid {
			validCount++
		}
	}
	s := Summary{
		Total:   len(results),
		Valid:   validCount,
		Invalid: len(results) - validCount,
		Results: results,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, `{"error": "encoding output: %v", "code": 2}`+"\n", err)
		os.Exit(exit.ToolError)
	}
}

func printMD(results []validator.Result, verbose bool) {
	validCount := 0
	for _, r := range results {
		if r.Status == validator.Valid {
			validCount++
		}
	}
	fmt.Printf("# probe results\n\n")
	fmt.Printf("**%d/%d valid**\n\n", validCount, len(results))
	fmt.Printf("| Path | Status | Type | Duration |\n")
	fmt.Printf("|---|---|---|---|\n")
	for _, r := range results {
		icon := "✓"
		if r.Status != validator.Valid {
			icon = "✗"
		}
		fmt.Printf("| `%s` | %s %s | %s | %dms |\n",
			r.Path, icon, r.Status, r.ArtifactType, r.DurationMs)
	}
	if verbose {
		for _, r := range results {
			if len(r.Errors) > 0 {
				fmt.Printf("\n## %s\n\n", r.Path)
				for _, e := range r.Errors {
					fmt.Printf("- %s\n", e.Message)
				}
			}
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `probe — validate gemara artifacts

Usage:
  probe [flags] <artifact.yaml> [<artifact2.yaml> ...]

Flags:
  --format string    Output format: json (default), md
  --program string   Program slug for provenance logging
  --quiet            Suppress per-artifact output; exit code only
  --version          Print version and exit

Exit codes:
  0  All artifacts are valid
  1  One or more artifacts failed validation
  2  Tool error (missing file, bad argument, etc.)

Examples:
  probe artifact.yaml
  probe --format md *.yaml
  probe --program iso42001 artifact.yaml && echo "valid"

GitHub Actions:
  - name: Validate gemara artifacts
    run: probe artifact.yaml`)
}
