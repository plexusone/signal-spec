package signal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/plexusone/signal-spec/pkg/common"
)

func TestTypeJSONSchema(t *testing.T) {
	schema := Type("").JSONSchema()
	if schema.Type != "string" {
		t.Errorf("expected type string, got %s", schema.Type)
	}
	if len(schema.Enum) != 8 {
		t.Errorf("expected 8 enum values, got %d", len(schema.Enum))
	}
}

func TestTypeValues(t *testing.T) {
	values := TypeValues()
	if len(values) != 8 {
		t.Errorf("expected 8 type values, got %d", len(values))
	}

	expected := []Type{
		TypeSupportTicket,
		TypeCloudIncident,
		TypeSecurityFinding,
		TypePostureDrift,
		TypeAlert,
		TypeOutage,
		TypeVulnerability,
		TypeFeedback,
	}

	for i, v := range expected {
		if values[i] != v {
			t.Errorf("expected values[%d] = %s, got %s", i, v, values[i])
		}
	}
}

func TestStatusJSONSchema(t *testing.T) {
	schema := Status("").JSONSchema()
	if schema.Type != "string" {
		t.Errorf("expected type string, got %s", schema.Type)
	}
	if len(schema.Enum) != 5 {
		t.Errorf("expected 5 enum values, got %d", len(schema.Enum))
	}
}

func TestSignalJSONMarshal(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	sig := Signal{
		ID:     "sig-001",
		Type:   TypeSupportTicket,
		Status: StatusNew,
		Source: common.SourceSystem{
			Type: "ticketing",
			Name: "zendesk",
		},
		Domain: common.Domain{
			Name:      "authentication",
			Subdomain: "oauth",
		},
		Severity:   common.SeverityHigh,
		Summary:    "OAuth token refresh failures",
		ObservedAt: now,
		ReceivedAt: now,
		Tags:       []common.Tag{"auth", "oauth"},
	}

	data, err := json.Marshal(sig)
	if err != nil {
		t.Fatalf("failed to marshal signal: %v", err)
	}

	var decoded Signal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal signal: %v", err)
	}

	if decoded.ID != sig.ID {
		t.Errorf("expected ID %s, got %s", sig.ID, decoded.ID)
	}
	if decoded.Type != sig.Type {
		t.Errorf("expected Type %s, got %s", sig.Type, decoded.Type)
	}
	if decoded.Status != sig.Status {
		t.Errorf("expected Status %s, got %s", sig.Status, decoded.Status)
	}
	if decoded.Severity != sig.Severity {
		t.Errorf("expected Severity %s, got %s", sig.Severity, decoded.Severity)
	}
	if decoded.Summary != sig.Summary {
		t.Errorf("expected Summary %s, got %s", sig.Summary, decoded.Summary)
	}
	if len(decoded.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(decoded.Tags))
	}
}

func TestSignalTypeConstants(t *testing.T) {
	tests := []struct {
		typ      Type
		expected string
	}{
		{TypeSupportTicket, "support_ticket"},
		{TypeCloudIncident, "cloud_incident"},
		{TypeSecurityFinding, "security_finding"},
		{TypePostureDrift, "posture_drift"},
		{TypeAlert, "alert"},
		{TypeOutage, "outage"},
		{TypeVulnerability, "vulnerability"},
		{TypeFeedback, "feedback"},
	}

	for _, tt := range tests {
		if string(tt.typ) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.typ)
		}
	}
}

func TestSignalStatusConstants(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusNew, "new"},
		{StatusProcessing, "processing"},
		{StatusMapped, "mapped"},
		{StatusIgnored, "ignored"},
		{StatusArchived, "archived"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.status)
		}
	}
}
