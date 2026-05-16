package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/plexusone/signal-spec/pkg/common"
	"github.com/plexusone/signal-spec/pkg/remediation"
	"github.com/plexusone/signal-spec/pkg/rootcause"
	"github.com/plexusone/signal-spec/pkg/signal"
	"github.com/spf13/cobra"
)

var validateType string

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a JSON file against signal-spec types",
	Long: `Validate that a JSON file conforms to signal-spec types.

Supported types:
  - signal: Validate as a Signal
  - rootcause: Validate as a RootCause
  - remediation: Validate as a Remediation

Examples:
  signal-spec validate -t signal signal.json
  signal-spec validate -t rootcause rootcause.json
  signal-spec validate -t remediation remediation.json`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&validateType, "type", "t", "", "Type to validate against (signal, rootcause, remediation)")
	validateCmd.MarkFlagRequired("type")
}

func runValidate(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var validationErrors []string

	switch validateType {
	case "signal":
		var s signal.Signal
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("parsing as Signal: %w", err)
		}
		validationErrors = validateSignal(s)

	case "rootcause":
		var rc rootcause.RootCause
		if err := json.Unmarshal(data, &rc); err != nil {
			return fmt.Errorf("parsing as RootCause: %w", err)
		}
		validationErrors = validateRootCause(rc)

	case "remediation":
		var r remediation.Remediation
		if err := json.Unmarshal(data, &r); err != nil {
			return fmt.Errorf("parsing as Remediation: %w", err)
		}
		validationErrors = validateRemediation(r)

	default:
		return fmt.Errorf("unknown type: %s (use signal, rootcause, or remediation)", validateType)
	}

	if len(validationErrors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", filename)
		for _, e := range validationErrors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	}

	fmt.Fprintf(os.Stderr, "Valid %s: %s\n", validateType, filename)
	return nil
}

func validateSignal(s signal.Signal) []string {
	var errors []string

	if s.ID == "" {
		errors = append(errors, "id is required")
	}
	if s.Type == "" {
		errors = append(errors, "type is required")
	}
	if s.Summary == "" {
		errors = append(errors, "summary is required")
	}
	if s.Domain.Name == "" {
		errors = append(errors, "domain.name is required")
	}

	// Validate tags are kebab-case
	if err := common.ValidateTags(s.Tags); err != nil {
		errors = append(errors, err.Error())
	}

	return errors
}

func validateRootCause(rc rootcause.RootCause) []string {
	var errors []string

	if rc.ID == "" {
		errors = append(errors, "id is required")
	}
	if rc.Title == "" {
		errors = append(errors, "title is required")
	}
	if rc.Domain.Name == "" {
		errors = append(errors, "domain.name is required")
	}

	// Validate tags are kebab-case
	if err := common.ValidateTags(rc.Tags); err != nil {
		errors = append(errors, err.Error())
	}

	return errors
}

func validateRemediation(r remediation.Remediation) []string {
	var errors []string

	if r.ID == "" {
		errors = append(errors, "id is required")
	}
	if r.Title == "" {
		errors = append(errors, "title is required")
	}
	if len(r.RootCauseIDs) == 0 {
		errors = append(errors, "root_cause_ids is required")
	}

	// Validate tags are kebab-case
	if err := common.ValidateTags(r.Tags); err != nil {
		errors = append(errors, err.Error())
	}

	return errors
}
