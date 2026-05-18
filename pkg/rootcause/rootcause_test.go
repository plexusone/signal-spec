package rootcause

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/plexusone/signal-spec/pkg/common"
)

func TestStatusJSONSchema(t *testing.T) {
	schema := Status("").JSONSchema()
	if schema.Type != "string" {
		t.Errorf("expected type string, got %s", schema.Type)
	}
	if len(schema.Enum) != 8 {
		t.Errorf("expected 8 enum values, got %d", len(schema.Enum))
	}
}

func TestTrendDirectionJSONSchema(t *testing.T) {
	schema := TrendDirection("").JSONSchema()
	if schema.Type != "string" {
		t.Errorf("expected type string, got %s", schema.Type)
	}
	if len(schema.Enum) != 3 {
		t.Errorf("expected 3 enum values, got %d", len(schema.Enum))
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusNew, "new"},
		{StatusActive, "active"},
		{StatusMitigating, "mitigating"},
		{StatusValidating, "validating"},
		{StatusStable, "stable"},
		{StatusRegressed, "regressed"},
		{StatusResolved, "resolved"},
		{StatusArchived, "archived"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.status)
		}
	}
}

func TestTrendDirectionConstants(t *testing.T) {
	tests := []struct {
		trend    TrendDirection
		expected string
	}{
		{TrendIncreasing, "increasing"},
		{TrendStable, "stable"},
		{TrendDecreasing, "decreasing"},
	}

	for _, tt := range tests {
		if string(tt.trend) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.trend)
		}
	}
}

func TestRootCauseJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	rc := RootCause{
		ID:     "rc-001",
		Title:  "Redis session replication instability",
		Status: StatusActive,
		Domain: common.Domain{
			Name:      "authentication",
			Subdomain: "oauth",
			Team:      "identity-platform",
		},
		Severity: common.SeverityHigh,
		SymptomPatterns: []string{
			"OAuth token failures",
			"Session expired",
		},
		SignalIDs: []string{"sig-001", "sig-002"},
		Impact: ImpactMetrics{
			SignalCount:       487,
			AffectedCustomers: 2341,
		},
		Trend: Trend{
			Direction: TrendStable,
			Velocity:  15.3,
		},
		PriorityScore:   87,
		FirstSeen:       now.Add(-24 * time.Hour),
		LastSeen:        now,
		OwnerTeam:       "identity-platform",
		RecurrenceCount: 1,
		Tags:            []common.Tag{"redis", "auth"},
	}

	data, err := json.Marshal(rc)
	if err != nil {
		t.Fatalf("failed to marshal root cause: %v", err)
	}

	var decoded RootCause
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal root cause: %v", err)
	}

	if decoded.ID != rc.ID {
		t.Errorf("expected ID %s, got %s", rc.ID, decoded.ID)
	}
	if decoded.Title != rc.Title {
		t.Errorf("expected Title %s, got %s", rc.Title, decoded.Title)
	}
	if decoded.Status != rc.Status {
		t.Errorf("expected Status %s, got %s", rc.Status, decoded.Status)
	}
	if decoded.Impact.SignalCount != rc.Impact.SignalCount {
		t.Errorf("expected SignalCount %d, got %d", rc.Impact.SignalCount, decoded.Impact.SignalCount)
	}
	if decoded.PriorityScore != rc.PriorityScore {
		t.Errorf("expected PriorityScore %d, got %d", rc.PriorityScore, decoded.PriorityScore)
	}
}

func TestImpactMetrics(t *testing.T) {
	impact := ImpactMetrics{
		SignalCount:          500,
		AffectedCustomers:    1000,
		EscalationRate:       0.15,
		EstimatedRevenueLoss: 50000.00,
	}

	if impact.SignalCount != 500 {
		t.Errorf("expected SignalCount 500, got %d", impact.SignalCount)
	}
	if impact.AffectedCustomers != 1000 {
		t.Errorf("expected AffectedCustomers 1000, got %d", impact.AffectedCustomers)
	}
	if impact.EscalationRate != 0.15 {
		t.Errorf("expected EscalationRate 0.15, got %f", impact.EscalationRate)
	}
}
