package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/signal-spec/pkg/remediation"
	"github.com/plexusone/signal-spec/pkg/rootcause"
	"github.com/plexusone/signal-spec/pkg/signal"
	"github.com/spf13/cobra"
)

var schemaOutputDir string

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "JSON schema operations",
	Long:  `Commands for working with signal-spec JSON schemas.`,
}

var schemaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate JSON schemas from Go types",
	Long: `Generate JSON schema files from the signal-spec Go types.

Schemas are generated for:
  - signal.schema.json
  - rootcause.schema.json
  - remediation.schema.json
  - validation_signal.schema.json

Examples:
  # Generate schemas to default directory (./schema)
  signal-spec schema generate

  # Generate schemas to custom directory
  signal-spec schema generate -o ./schemas`,
	RunE: runSchemaGenerate,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.AddCommand(schemaGenerateCmd)

	schemaGenerateCmd.Flags().StringVarP(&schemaOutputDir, "output", "o", "schema", "Output directory for schema files")
}

func runSchemaGenerate(cmd *cobra.Command, args []string) error {
	schemas := []struct {
		name string
		typ  any
	}{
		{"signal", signal.Signal{}},
		{"rootcause", rootcause.RootCause{}},
		{"remediation", remediation.Remediation{}},
		{"validation_signal", remediation.ValidationSignal{}},
	}

	// Ensure output directory exists
	if err := os.MkdirAll(schemaOutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, s := range schemas {
		r := jsonschema.Reflector{
			DoNotReference: false,
			ExpandedStruct: false,
		}

		schema := r.Reflect(s.typ)
		schema.ID = jsonschema.ID(fmt.Sprintf("https://plexusone.dev/signal-spec/%s.schema.json", s.name))
		schema.Version = "https://json-schema.org/draft/2020-12/schema"

		data, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling %s: %w", s.name, err)
		}

		filename := filepath.Join(schemaOutputDir, s.name+".schema.json")
		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("writing %s: %w", filename, err)
		}

		fmt.Fprintf(os.Stderr, "Generated %s\n", filename)
	}

	return nil
}
