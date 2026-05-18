package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plexusone/signal-spec/pkg/common"
	"github.com/plexusone/signal-spec/pkg/rootcause"
)

func TestBuildSummaryReport(t *testing.T) {
	rootCauses := []rootcause.RootCause{
		{
			ID:            "rc-001",
			Title:         "Issue 1",
			Domain:        common.Domain{Name: "auth", Subdomain: "oauth"},
			Severity:      common.SeverityHigh,
			Impact:        rootcause.ImpactMetrics{SignalCount: 100},
			PriorityScore: 80,
		},
		{
			ID:            "rc-002",
			Title:         "Issue 2",
			Domain:        common.Domain{Name: "auth", Subdomain: "oauth"},
			Severity:      common.SeverityCritical,
			Impact:        rootcause.ImpactMetrics{SignalCount: 200},
			PriorityScore: 90,
		},
		{
			ID:            "rc-003",
			Title:         "Issue 3",
			Domain:        common.Domain{Name: "payments", Subdomain: "checkout"},
			Severity:      common.SeverityMedium,
			Impact:        rootcause.ImpactMetrics{SignalCount: 50},
			PriorityScore: 60,
		},
	}

	report := BuildSummaryReport(rootCauses)

	if len(report.Summaries) != 2 {
		t.Errorf("expected 2 domain summaries, got %d", len(report.Summaries))
	}

	if len(report.RootCauses) != 3 {
		t.Errorf("expected 3 root causes, got %d", len(report.RootCauses))
	}

	// First summary should be auth|oauth with 300 total signals (highest)
	if report.Summaries[0].Domain != "auth" {
		t.Errorf("expected first domain to be 'auth', got %s", report.Summaries[0].Domain)
	}
	if report.Summaries[0].TotalSignals != 300 {
		t.Errorf("expected 300 total signals, got %d", report.Summaries[0].TotalSignals)
	}
	if report.Summaries[0].IssueCount != 2 {
		t.Errorf("expected 2 issues, got %d", report.Summaries[0].IssueCount)
	}
	if report.Summaries[0].MaxSeverity != "critical" {
		t.Errorf("expected max severity 'critical', got %s", report.Summaries[0].MaxSeverity)
	}

	// Root causes should be sorted by priority score descending
	if report.RootCauses[0].PriorityScore != 90 {
		t.Errorf("expected first root cause priority 90, got %d", report.RootCauses[0].PriorityScore)
	}
}

func TestLeaderMappingApply(t *testing.T) {
	mapping := &LeaderMapping{
		AreaLeaders: map[string]string{
			"auth": "Jane Smith",
		},
		ExecutionLeaders: map[string]string{
			"auth|oauth": "Mike Johnson",
		},
	}

	summaries := []DomainSummary{
		{Domain: "auth", Subdomain: "oauth"},
		{Domain: "auth", Subdomain: "ldap"},
		{Domain: "payments", Subdomain: "checkout"},
	}

	mapping.ApplyLeaders(summaries)

	if summaries[0].AreaLeader != "Jane Smith" {
		t.Errorf("expected area leader 'Jane Smith', got %s", summaries[0].AreaLeader)
	}
	if summaries[0].ExecutionLeader != "Mike Johnson" {
		t.Errorf("expected execution leader 'Mike Johnson', got %s", summaries[0].ExecutionLeader)
	}

	// auth|ldap should have area leader but no execution leader
	if summaries[1].AreaLeader != "Jane Smith" {
		t.Errorf("expected area leader 'Jane Smith', got %s", summaries[1].AreaLeader)
	}
	if summaries[1].ExecutionLeader != "" {
		t.Errorf("expected no execution leader, got %s", summaries[1].ExecutionLeader)
	}

	// payments should have no leaders
	if summaries[2].AreaLeader != "" {
		t.Errorf("expected no area leader, got %s", summaries[2].AreaLeader)
	}
}

func TestWriteXLSX(t *testing.T) {
	rootCauses := []rootcause.RootCause{
		{
			ID:            "rc-001",
			Title:         "Test Issue",
			Domain:        common.Domain{Name: "test", Subdomain: "unit"},
			Severity:      common.SeverityHigh,
			Impact:        rootcause.ImpactMetrics{SignalCount: 100},
			PriorityScore: 80,
		},
	}

	report := BuildSummaryReport(rootCauses)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-report.xlsx")

	err := report.WriteXLSX(outputPath)
	if err != nil {
		t.Fatalf("failed to write XLSX: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("XLSX file was not created")
	}

	// Verify file is not empty
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("failed to stat XLSX file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("XLSX file is empty")
	}
}

func TestLoadRootCausesFromFile(t *testing.T) {
	// Create a temp file with test data
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	jsonData := `[
		{"id": "rc-001", "title": "Issue 1", "domain": {"name": "auth"}},
		{"id": "rc-002", "title": "Issue 2", "domain": {"name": "payments"}}
	]`

	if err := os.WriteFile(testFile, []byte(jsonData), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	rootCauses, err := LoadRootCausesFromFile(testFile)
	if err != nil {
		t.Fatalf("failed to load root causes: %v", err)
	}

	if len(rootCauses) != 2 {
		t.Errorf("expected 2 root causes, got %d", len(rootCauses))
	}
	if rootCauses[0].ID != "rc-001" {
		t.Errorf("expected ID 'rc-001', got %s", rootCauses[0].ID)
	}
}

func TestLoadRootCausesFromFileSingleObject(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "single.json")

	jsonData := `{"id": "rc-001", "title": "Single Issue", "domain": {"name": "auth"}}`

	if err := os.WriteFile(testFile, []byte(jsonData), 0600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	rootCauses, err := LoadRootCausesFromFile(testFile)
	if err != nil {
		t.Fatalf("failed to load root causes: %v", err)
	}

	if len(rootCauses) != 1 {
		t.Errorf("expected 1 root cause, got %d", len(rootCauses))
	}
}

func TestLoadRootCausesFromDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple JSON files
	file1 := `{"id": "rc-001", "title": "Issue 1", "domain": {"name": "auth"}}`
	file2 := `{"id": "rc-002", "title": "Issue 2", "domain": {"name": "payments"}}`

	if err := os.WriteFile(filepath.Join(tmpDir, "rc1.json"), []byte(file1), 0600); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "rc2.json"), []byte(file2), 0600); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}
	// Create a non-JSON file that should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignore me"), 0600); err != nil {
		t.Fatalf("failed to write txt file: %v", err)
	}

	rootCauses, err := LoadRootCausesFromDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to load root causes from dir: %v", err)
	}

	if len(rootCauses) != 2 {
		t.Errorf("expected 2 root causes, got %d", len(rootCauses))
	}
}

func TestDomainSummaryTopIssuesLimit(t *testing.T) {
	// Create 10 root causes for the same domain
	var rootCauses []rootcause.RootCause
	for i := 0; i < 10; i++ {
		rootCauses = append(rootCauses, rootcause.RootCause{
			ID:     "rc-" + string(rune('0'+i)),
			Title:  "Issue " + string(rune('0'+i)),
			Domain: common.Domain{Name: "auth", Subdomain: "oauth"},
			Impact: rootcause.ImpactMetrics{SignalCount: 10},
		})
	}

	report := BuildSummaryReport(rootCauses)

	if len(report.Summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(report.Summaries))
	}

	// TopIssues should be limited to 5
	if len(report.Summaries[0].TopIssues) != 5 {
		t.Errorf("expected 5 top issues, got %d", len(report.Summaries[0].TopIssues))
	}
}
