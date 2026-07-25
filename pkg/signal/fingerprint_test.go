package signal

import (
	"testing"
	"time"

	"github.com/plexusone/signal-spec/pkg/common"
)

func TestComputeFingerprintDeterministic(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	sig := Signal{
		ID:   "sig-001",
		Type: TypeSupportTicket,
		Source: common.SourceSystem{
			Type: "ticketing",
			Name: "jira",
		},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
	}

	fp1, err := ComputeFingerprint(sig)
	if err != nil {
		t.Fatalf("first fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(sig)
	if err != nil {
		t.Fatalf("second fingerprint: %v", err)
	}

	if fp1 != fp2 {
		t.Errorf("fingerprint not deterministic: %s != %s", fp1, fp2)
	}
	if len(fp1) != 64 {
		t.Errorf("expected 64-char hex SHA-256, got %d chars", len(fp1))
	}
}

func TestComputeFingerprintExcludesDerived(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	base := Signal{
		ID:         "sig-001",
		Type:       TypeSupportTicket,
		Source:     common.SourceSystem{Type: "ticketing", Name: "jira"},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
	}

	frustration := 42.5
	withDerived := base
	withDerived.Derived = &DerivedMetrics{
		Frustration: &frustration,
	}

	fp1, err := ComputeFingerprint(base)
	if err != nil {
		t.Fatalf("base fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(withDerived)
	if err != nil {
		t.Fatalf("derived fingerprint: %v", err)
	}

	if fp1 != fp2 {
		t.Errorf("derived metrics changed fingerprint: %s != %s", fp1, fp2)
	}
}

func TestComputeFingerprintExcludesStatus(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	sig1 := Signal{
		ID:         "sig-001",
		Type:       TypeSupportTicket,
		Status:     StatusNew,
		Source:     common.SourceSystem{Type: "ticketing", Name: "jira"},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
	}

	sig2 := sig1
	sig2.Status = StatusMapped

	fp1, err := ComputeFingerprint(sig1)
	if err != nil {
		t.Fatalf("status new fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(sig2)
	if err != nil {
		t.Fatalf("status mapped fingerprint: %v", err)
	}

	if fp1 != fp2 {
		t.Errorf("status change affected fingerprint: %s != %s", fp1, fp2)
	}
}

func TestComputeFingerprintExcludesReceivedAt(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	sig1 := Signal{
		ID:         "sig-001",
		Type:       TypeSupportTicket,
		Source:     common.SourceSystem{Type: "ticketing", Name: "jira"},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
		ReceivedAt: now,
	}

	sig2 := sig1
	sig2.ReceivedAt = now.Add(5 * time.Minute)

	fp1, err := ComputeFingerprint(sig1)
	if err != nil {
		t.Fatalf("first fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(sig2)
	if err != nil {
		t.Fatalf("second fingerprint: %v", err)
	}

	if fp1 != fp2 {
		t.Errorf("ReceivedAt change affected fingerprint: %s != %s", fp1, fp2)
	}
}

func TestComputeFingerprintExcludesEmbedding(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	sig1 := Signal{
		ID:         "sig-001",
		Type:       TypeSupportTicket,
		Source:     common.SourceSystem{Type: "ticketing", Name: "jira"},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
	}

	sig2 := sig1
	sig2.Embedding = []float32{0.1, 0.2, 0.3}

	fp1, err := ComputeFingerprint(sig1)
	if err != nil {
		t.Fatalf("first fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(sig2)
	if err != nil {
		t.Fatalf("second fingerprint: %v", err)
	}

	if fp1 != fp2 {
		t.Errorf("embedding change affected fingerprint: %s != %s", fp1, fp2)
	}
}

func TestComputeFingerprintDiffersOnContent(t *testing.T) {
	now := time.Date(2026, 7, 24, 12, 0, 0, 0, time.UTC)

	sig1 := Signal{
		ID:         "sig-001",
		Type:       TypeSupportTicket,
		Source:     common.SourceSystem{Type: "ticketing", Name: "jira"},
		Domain:     common.Domain{Name: "infra"},
		Severity:   common.SeverityHigh,
		Summary:    "DB pool exhaustion",
		ObservedAt: now,
	}

	sig2 := sig1
	sig2.Summary = "Different summary"

	fp1, err := ComputeFingerprint(sig1)
	if err != nil {
		t.Fatalf("first fingerprint: %v", err)
	}
	fp2, err := ComputeFingerprint(sig2)
	if err != nil {
		t.Fatalf("second fingerprint: %v", err)
	}

	if fp1 == fp2 {
		t.Error("different content produced same fingerprint")
	}
}
