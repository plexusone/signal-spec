package cmd

import (
	"fmt"
	"os"

	"github.com/plexusone/signal-spec/pkg/export"
	"github.com/plexusone/signal-spec/pkg/rootcause"
	"github.com/spf13/cobra"
)

var (
	reportInputFile  string
	reportInputDir   string
	reportLeaderFile string
	reportOutput     string
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate XLSX summary reports from root cause data",
	Long: `Generate Excel reports summarizing root causes by domain and subdomain.

The report includes two sheets:
  - Domain Summary: Aggregated metrics by domain/subdomain for prioritization
  - Root Causes: Detailed list of all root causes

Examples:
  # Generate report from a JSON file
  signal-spec report -i rootcauses.json -o summary.xlsx

  # Generate report from a directory of JSON files
  signal-spec report -d ./rootcauses/ -o summary.xlsx

  # Include leader mappings
  signal-spec report -i rootcauses.json --leaders leaders.json -o summary.xlsx`,
	RunE: runReport,
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.Flags().StringVarP(&reportInputFile, "input", "i", "", "Input JSON file (single or array of root causes)")
	reportCmd.Flags().StringVarP(&reportInputDir, "dir", "d", "", "Input directory containing root cause JSON files")
	reportCmd.Flags().StringVar(&reportLeaderFile, "leaders", "", "JSON file mapping domains to area/execution leaders")
	reportCmd.Flags().StringVarP(&reportOutput, "output", "o", "domain-summary.xlsx", "Output XLSX file")

	reportCmd.MarkFlagsMutuallyExclusive("input", "dir")
}

func runReport(cmd *cobra.Command, args []string) error {
	if reportInputFile == "" && reportInputDir == "" {
		return fmt.Errorf("specify --input <file> or --dir <directory>")
	}

	var rootCauses []rootcause.RootCause
	var err error

	if reportInputFile != "" {
		rootCauses, err = export.LoadRootCausesFromFile(reportInputFile)
	} else {
		rootCauses, err = export.LoadRootCausesFromDir(reportInputDir)
	}

	if err != nil {
		return fmt.Errorf("loading root causes: %w", err)
	}

	if len(rootCauses) == 0 {
		return fmt.Errorf("no root causes found")
	}

	fmt.Fprintf(os.Stderr, "Loaded %d root causes\n", len(rootCauses))

	report := export.BuildSummaryReport(rootCauses)

	// Apply leader mappings if provided
	if reportLeaderFile != "" {
		mapping, err := export.LoadLeaderMapping(reportLeaderFile)
		if err != nil {
			return fmt.Errorf("loading leader mapping: %w", err)
		}
		mapping.ApplyLeaders(report.Summaries)
		fmt.Fprintf(os.Stderr, "Applied leader mappings from %s\n", reportLeaderFile)
	}

	if err := report.WriteXLSX(reportOutput); err != nil {
		return fmt.Errorf("writing XLSX: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Generated %s with %d domain summaries\n", reportOutput, len(report.Summaries))
	return nil
}
