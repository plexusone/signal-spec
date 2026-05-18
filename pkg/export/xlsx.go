// Package export provides report generation from signal-spec data.
package export

import (
	"fmt"
	"sort"

	"github.com/plexusone/signal-spec/pkg/rootcause"
	"github.com/xuri/excelize/v2"
)

// DomainSummary aggregates root cause metrics by domain/subdomain.
type DomainSummary struct {
	Domain           string
	Subdomain        string
	IssueCount       int
	TotalSignals     int
	AvgPriorityScore float64
	MaxSeverity      string
	AreaLeader       string // Assignable - maps to domain owner
	ExecutionLeader  string // Assignable - maps to subdomain owner
	TopIssues        []string
}

// SummaryReport contains aggregated domain summaries and root cause details.
type SummaryReport struct {
	Summaries  []DomainSummary
	RootCauses []rootcause.RootCause
}

// BuildSummaryReport aggregates root causes by domain/subdomain.
func BuildSummaryReport(rootCauses []rootcause.RootCause) *SummaryReport {
	// Map key: "domain|subdomain"
	aggregates := make(map[string]*DomainSummary)

	severityRank := map[string]int{
		"critical": 5,
		"high":     4,
		"medium":   3,
		"low":      2,
		"info":     1,
	}

	for _, rc := range rootCauses {
		key := rc.Domain.Name + "|" + rc.Domain.Subdomain
		summary, exists := aggregates[key]
		if !exists {
			summary = &DomainSummary{
				Domain:      rc.Domain.Name,
				Subdomain:   rc.Domain.Subdomain,
				AreaLeader:  rc.Domain.Team, // Default to team if set
				MaxSeverity: string(rc.Severity),
			}
			aggregates[key] = summary
		}

		summary.IssueCount++
		summary.TotalSignals += rc.Impact.SignalCount
		summary.AvgPriorityScore += float64(rc.PriorityScore)

		// Track max severity
		if severityRank[string(rc.Severity)] > severityRank[summary.MaxSeverity] {
			summary.MaxSeverity = string(rc.Severity)
		}

		// Collect top issues (limit to 5)
		if len(summary.TopIssues) < 5 {
			summary.TopIssues = append(summary.TopIssues, rc.Title)
		}
	}

	// Calculate averages and build result
	var summaries []DomainSummary
	for _, summary := range aggregates {
		if summary.IssueCount > 0 {
			summary.AvgPriorityScore /= float64(summary.IssueCount)
		}
		summaries = append(summaries, *summary)
	}

	// Sort by total signals descending (highest impact first)
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].TotalSignals > summaries[j].TotalSignals
	})

	// Sort root causes by priority score descending
	sortedRCs := make([]rootcause.RootCause, len(rootCauses))
	copy(sortedRCs, rootCauses)
	sort.Slice(sortedRCs, func(i, j int) bool {
		return sortedRCs[i].PriorityScore > sortedRCs[j].PriorityScore
	})

	return &SummaryReport{
		Summaries:  summaries,
		RootCauses: sortedRCs,
	}
}

// WriteXLSX generates an Excel report from the summary.
func (r *SummaryReport) WriteXLSX(filename string) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Domain Summary"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("creating sheet: %w", err)
	}
	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1") // Default sheet, error unlikely

	// Define headers
	headers := []string{
		"Domain (Category)",
		"Subdomain (Subcategory)",
		"Issue Count",
		"Total Signals",
		"Avg Priority",
		"Max Severity",
		"Area Leader",
		"Execution Leader",
		"Top Issues",
	}

	// Write headers with styling
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})
	if err != nil {
		return fmt.Errorf("creating header style: %w", err)
	}

	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		_ = f.SetCellValue(sheetName, cell, header)
		_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Severity color mapping
	severityColors := map[string]string{
		"critical": "FF0000",
		"high":     "FFA500",
		"medium":   "FFFF00",
		"low":      "90EE90",
		"info":     "ADD8E6",
	}

	// Write data rows (SetCellValue errors ignored - valid cells from controlled input)
	for row, summary := range r.Summaries {
		rowNum := row + 2

		_ = f.SetCellValue(sheetName, cellName(1, rowNum), summary.Domain)
		_ = f.SetCellValue(sheetName, cellName(2, rowNum), summary.Subdomain)
		_ = f.SetCellValue(sheetName, cellName(3, rowNum), summary.IssueCount)
		_ = f.SetCellValue(sheetName, cellName(4, rowNum), summary.TotalSignals)
		_ = f.SetCellValue(sheetName, cellName(5, rowNum), fmt.Sprintf("%.1f", summary.AvgPriorityScore))
		_ = f.SetCellValue(sheetName, cellName(6, rowNum), summary.MaxSeverity)
		_ = f.SetCellValue(sheetName, cellName(7, rowNum), summary.AreaLeader)
		_ = f.SetCellValue(sheetName, cellName(8, rowNum), summary.ExecutionLeader)

		// Join top issues
		topIssuesStr := ""
		for i, issue := range summary.TopIssues {
			if i > 0 {
				topIssuesStr += "; "
			}
			topIssuesStr += issue
		}
		_ = f.SetCellValue(sheetName, cellName(9, rowNum), topIssuesStr)

		// Color-code severity cell
		if color, ok := severityColors[summary.MaxSeverity]; ok {
			style, _ := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
			})
			_ = f.SetCellStyle(sheetName, cellName(6, rowNum), cellName(6, rowNum), style)
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 20, // Domain
		"B": 20, // Subdomain
		"C": 12, // Issue Count
		"D": 12, // Total Signals
		"E": 12, // Avg Priority
		"F": 12, // Max Severity
		"G": 20, // Area Leader
		"H": 20, // Execution Leader
		"I": 60, // Top Issues
	}
	for col, width := range colWidths {
		_ = f.SetColWidth(sheetName, col, col, width) // Valid column, error impossible
	}

	// Freeze header row
	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Add auto-filter
	lastRow := len(r.Summaries) + 1
	_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:I%d", lastRow), nil)

	// Create Root Causes sheet
	if err := r.writeRootCausesSheet(f, severityColors); err != nil {
		return fmt.Errorf("writing root causes sheet: %w", err)
	}

	return f.SaveAs(filename)
}

// writeRootCausesSheet adds a detailed root causes sheet.
func (r *SummaryReport) writeRootCausesSheet(f *excelize.File, severityColors map[string]string) error {
	sheetName := "Root Causes"
	if _, err := f.NewSheet(sheetName); err != nil {
		return err
	}

	// Define headers
	headers := []string{
		"ID",
		"Title",
		"Domain",
		"Subdomain",
		"Status",
		"Severity",
		"Signal Count",
		"Priority Score",
		"First Seen",
		"Last Seen",
		"Owner Team",
		"Tags",
	}

	// Write headers with styling
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "000000", Style: 2},
		},
	})

	for col, header := range headers {
		cell := cellName(col+1, 1)
		_ = f.SetCellValue(sheetName, cell, header)
		_ = f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Write data rows (SetCellValue errors ignored - valid cells from controlled input)
	for row, rc := range r.RootCauses {
		rowNum := row + 2

		_ = f.SetCellValue(sheetName, cellName(1, rowNum), rc.ID)
		_ = f.SetCellValue(sheetName, cellName(2, rowNum), rc.Title)
		_ = f.SetCellValue(sheetName, cellName(3, rowNum), rc.Domain.Name)
		_ = f.SetCellValue(sheetName, cellName(4, rowNum), rc.Domain.Subdomain)
		_ = f.SetCellValue(sheetName, cellName(5, rowNum), string(rc.Status))
		_ = f.SetCellValue(sheetName, cellName(6, rowNum), string(rc.Severity))
		_ = f.SetCellValue(sheetName, cellName(7, rowNum), rc.Impact.SignalCount)
		_ = f.SetCellValue(sheetName, cellName(8, rowNum), rc.PriorityScore)

		// Format dates
		if !rc.FirstSeen.IsZero() {
			_ = f.SetCellValue(sheetName, cellName(9, rowNum), rc.FirstSeen.Format("2006-01-02"))
		}
		if !rc.LastSeen.IsZero() {
			_ = f.SetCellValue(sheetName, cellName(10, rowNum), rc.LastSeen.Format("2006-01-02"))
		}

		_ = f.SetCellValue(sheetName, cellName(11, rowNum), rc.OwnerTeam)

		// Join tags
		tagsStr := ""
		for i, tag := range rc.Tags {
			if i > 0 {
				tagsStr += ", "
			}
			tagsStr += string(tag)
		}
		_ = f.SetCellValue(sheetName, cellName(12, rowNum), tagsStr)

		// Color-code severity cell
		if color, ok := severityColors[string(rc.Severity)]; ok {
			style, _ := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{color}, Pattern: 1},
			})
			_ = f.SetCellStyle(sheetName, cellName(6, rowNum), cellName(6, rowNum), style)
		}
	}

	// Set column widths
	colWidths := map[string]float64{
		"A": 15, // ID
		"B": 50, // Title
		"C": 18, // Domain
		"D": 18, // Subdomain
		"E": 12, // Status
		"F": 10, // Severity
		"G": 12, // Signal Count
		"H": 12, // Priority Score
		"I": 12, // First Seen
		"J": 12, // Last Seen
		"K": 20, // Owner Team
		"L": 30, // Tags
	}
	for col, width := range colWidths {
		_ = f.SetColWidth(sheetName, col, col, width) // Valid column, error impossible
	}

	// Freeze header row
	_ = f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	// Add auto-filter
	lastRow := len(r.RootCauses) + 1
	_ = f.AutoFilter(sheetName, fmt.Sprintf("A1:L%d", lastRow), nil)

	return nil
}

func cellName(col, row int) string {
	name, _ := excelize.CoordinatesToCellName(col, row)
	return name
}
