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
	if len(schema.Enum) != 13 {
		t.Errorf("expected 9 enum values, got %d", len(schema.Enum))
	}
}

func TestTypeValues(t *testing.T) {
	values := TypeValues()
	if len(values) != 13 {
		t.Errorf("expected 9 type values, got %d", len(values))
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
		TypeEnhancementRequest,
		TypeCompetitiveGap,
		TypeCompetitorLaunch,
		TypeAnalystFinding,
		TypeMarketObservation,
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
		{TypeEnhancementRequest, "enhancement_request"},
		{TypeCompetitiveGap, "competitive_gap"},
		{TypeCompetitorLaunch, "competitor_launch"},
		{TypeAnalystFinding, "analyst_finding"},
		{TypeMarketObservation, "market_observation"},
	}

	for _, tt := range tests {
		if string(tt.typ) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.typ)
		}
	}
}

func TestEnhancementRequestMetadata(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	sig := Signal{
		ID:     "sig-enhancement-001",
		Type:   TypeEnhancementRequest,
		Status: StatusNew,
		Source: common.SourceSystem{
			Type: "product_management",
			Name: "aha",
		},
		Domain: common.Domain{
			Name: "identity",
		},
		Severity:   common.SeverityMedium,
		Summary:    "Support SCIM provisioning for enterprise SSO",
		ObservedAt: now,
		ReceivedAt: now,
		Metadata: map[string]any{
			MetaVotes:         42,
			MetaSubscribers:   15,
			MetaOrganizations: []string{"Acme Corp", "Globex"},
			MetaCustomers:     []string{"acme-001", "globex-002"},
			MetaOpportunities: []string{"OPP-2026-100"},
			MetaEstimatedARR:  150000_00,
		},
	}

	data, err := json.Marshal(sig)
	if err != nil {
		t.Fatalf("failed to marshal enhancement request signal: %v", err)
	}

	var decoded Signal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal enhancement request signal: %v", err)
	}

	if decoded.Type != TypeEnhancementRequest {
		t.Errorf("expected type %s, got %s", TypeEnhancementRequest, decoded.Type)
	}
	if decoded.Metadata[MetaVotes] != float64(42) {
		t.Errorf("expected votes 42, got %v", decoded.Metadata[MetaVotes])
	}
	if decoded.Metadata[MetaEstimatedARR] != float64(150000_00) {
		t.Errorf("expected estimated_arr 15000000, got %v", decoded.Metadata[MetaEstimatedARR])
	}
}

func TestDerivedMetricsExcludedFromFingerprint(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	base := Signal{
		ID:     "sig-fp-001",
		Type:   TypeEnhancementRequest,
		Status: StatusNew,
		Source: common.SourceSystem{
			Type: "product_management",
			Name: "aha",
		},
		Domain:      common.Domain{Name: "identity"},
		Severity:    common.SeverityMedium,
		Summary:     "Add SCIM provisioning",
		ObservedAt:  now,
		ReceivedAt:  now,
		Fingerprint: "abc123",
		Metadata:    map[string]any{MetaVotes: 10},
	}

	frustration := 42.5
	momentum := 7.0
	withDerived := base
	withDerived.Derived = &DerivedMetrics{
		Frustration: &frustration,
		Momentum:    &momentum,
		ComputedAt:  &now,
	}

	dataBase, err := json.Marshal(base)
	if err != nil {
		t.Fatalf("marshal base: %v", err)
	}
	dataDerived, err := json.Marshal(withDerived)
	if err != nil {
		t.Fatalf("marshal derived: %v", err)
	}

	if base.Fingerprint != withDerived.Fingerprint {
		t.Error("fingerprint should be identical regardless of derived metrics")
	}

	var decoded Signal
	if err := json.Unmarshal(dataDerived, &decoded); err != nil {
		t.Fatalf("unmarshal derived: %v", err)
	}
	if decoded.Derived == nil {
		t.Fatal("expected derived metrics to be present after round-trip")
	}
	if *decoded.Derived.Frustration != 42.5 {
		t.Errorf("expected frustration 42.5, got %v", *decoded.Derived.Frustration)
	}
	if *decoded.Derived.Momentum != 7.0 {
		t.Errorf("expected momentum 7.0, got %v", *decoded.Derived.Momentum)
	}

	var decodedBase Signal
	if err := json.Unmarshal(dataBase, &decodedBase); err != nil {
		t.Fatalf("unmarshal base: %v", err)
	}
	if decodedBase.Derived != nil {
		t.Error("base signal without derived should have nil Derived after round-trip")
	}
}

func TestDerivedMetricsExtra(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	dm := DerivedMetrics{
		ComputedAt: &now,
		Extra: map[string]float64{
			"custom_score":    0.85,
			"weighted_impact": 12.3,
		},
	}

	sig := Signal{
		ID:         "sig-extra-001",
		Type:       TypeFeedback,
		Status:     StatusNew,
		Source:     common.SourceSystem{Type: "survey", Name: "typeform"},
		Domain:     common.Domain{Name: "ux"},
		Severity:   common.SeverityLow,
		Summary:    "UI feels slow",
		ObservedAt: now,
		ReceivedAt: now,
		Derived:    &dm,
	}

	data, err := json.Marshal(sig)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Signal
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.Derived == nil || len(decoded.Derived.Extra) != 2 {
		t.Fatalf("expected 2 extra metrics, got %v", decoded.Derived)
	}
	if decoded.Derived.Extra["custom_score"] != 0.85 {
		t.Errorf("expected custom_score 0.85, got %v", decoded.Derived.Extra["custom_score"])
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
