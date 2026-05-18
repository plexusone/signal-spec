package remediation

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
	if len(schema.Enum) != 9 {
		t.Errorf("expected 9 enum values, got %d", len(schema.Enum))
	}
}

func TestValidationResultJSONSchema(t *testing.T) {
	schema := ValidationResult("").JSONSchema()
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
		{StatusProposed, "proposed"},
		{StatusApproved, "approved"},
		{StatusInProgress, "in_progress"},
		{StatusDeployed, "deployed"},
		{StatusValidating, "validating"},
		{StatusEffective, "effective"},
		{StatusIneffective, "ineffective"},
		{StatusRolledBack, "rolled_back"},
		{StatusCancelled, "cancelled"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.status)
		}
	}
}

func TestValidationResultConstants(t *testing.T) {
	tests := []struct {
		result   ValidationResult
		expected string
	}{
		{ValidationImproved, "improved"},
		{ValidationNoChange, "no_change"},
		{ValidationRegressed, "regressed"},
	}

	for _, tt := range tests {
		if string(tt.result) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.result)
		}
	}
}

func TestRemediationJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	deployedAt := now.Add(-1 * time.Hour)

	rem := Remediation{
		ID:           "rem-001",
		Title:        "Implement Redis read-after-write consistency",
		Description:  "Modify session validation to use WAIT command",
		Status:       StatusDeployed,
		RootCauseIDs: []string{"rc-001"},
		OwnerTeam:    "identity-platform",
		Assignee:     "jsmith",
		CreatedAt:    now.Add(-24 * time.Hour),
		DeployedAt:   &deployedAt,
		ExternalLinks: []common.SourceSystem{
			{
				Type:       "code_change",
				Name:       "github",
				ExternalID: "PR-4521",
			},
		},
		Tags: []common.Tag{"redis", "consistency"},
	}

	data, err := json.Marshal(rem)
	if err != nil {
		t.Fatalf("failed to marshal remediation: %v", err)
	}

	var decoded Remediation
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal remediation: %v", err)
	}

	if decoded.ID != rem.ID {
		t.Errorf("expected ID %s, got %s", rem.ID, decoded.ID)
	}
	if decoded.Title != rem.Title {
		t.Errorf("expected Title %s, got %s", rem.Title, decoded.Title)
	}
	if decoded.Status != rem.Status {
		t.Errorf("expected Status %s, got %s", rem.Status, decoded.Status)
	}
	if len(decoded.RootCauseIDs) != 1 {
		t.Errorf("expected 1 root cause ID, got %d", len(decoded.RootCauseIDs))
	}
	if decoded.DeployedAt == nil {
		t.Error("expected DeployedAt to be set")
	}
}

func TestValidationSignalJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	vs := ValidationSignal{
		ID:                 "val-001",
		RemediationID:      "rem-001",
		RootCauseID:        "rc-001",
		Type:               ValidationImproved,
		ObservedAt:         now,
		BaselineSignalRate: 50.0,
		CurrentSignalRate:  10.0,
		ReductionPercent:   80.0,
		Notes:              "Signal rate dropped significantly after deployment",
	}

	data, err := json.Marshal(vs)
	if err != nil {
		t.Fatalf("failed to marshal validation signal: %v", err)
	}

	var decoded ValidationSignal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal validation signal: %v", err)
	}

	if decoded.ID != vs.ID {
		t.Errorf("expected ID %s, got %s", vs.ID, decoded.ID)
	}
	if decoded.Type != vs.Type {
		t.Errorf("expected Type %s, got %s", vs.Type, decoded.Type)
	}
	if decoded.ReductionPercent != vs.ReductionPercent {
		t.Errorf("expected ReductionPercent %f, got %f", vs.ReductionPercent, decoded.ReductionPercent)
	}
}

func TestEfficacy(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	eff := Efficacy{
		SignalReduction: 85.5,
		ValidationPeriod: common.TimeRange{
			Start: now.Add(-7 * 24 * time.Hour),
			End:   now,
		},
		ConfidenceLevel: 0.95,
		Notes:           "High confidence based on 7-day observation",
	}

	if eff.SignalReduction != 85.5 {
		t.Errorf("expected SignalReduction 85.5, got %f", eff.SignalReduction)
	}
	if eff.ConfidenceLevel != 0.95 {
		t.Errorf("expected ConfidenceLevel 0.95, got %f", eff.ConfidenceLevel)
	}
}
