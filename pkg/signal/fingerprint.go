package signal

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// fingerprintInput is the subset of Signal fields that define its identity.
// Excluded: Fingerprint (output), Embedding (derived), Derived (computed),
// Status (mutable), ReceivedAt (varies per ingestion), RootCauseID (assigned later).
type fingerprintInput struct {
	ID          string         `json:"id"`
	Type        Type           `json:"type"`
	Source      any            `json:"source"`
	Domain      any            `json:"domain"`
	Severity    any            `json:"severity"`
	Summary     string         `json:"summary"`
	Description string         `json:"description,omitempty"`
	Entities    any            `json:"entities,omitempty"`
	ObservedAt  any            `json:"observed_at"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Tags        any            `json:"tags,omitempty"`
}

// ComputeFingerprint returns a deterministic SHA-256 hex digest of the
// signal's identity fields. The same raw input always produces the same
// fingerprint, regardless of Status, Derived, Embedding, ReceivedAt,
// RootCauseID, or the Fingerprint field itself.
func ComputeFingerprint(s Signal) (string, error) {
	input := fingerprintInput{
		ID:          s.ID,
		Type:        s.Type,
		Source:      s.Source,
		Domain:      s.Domain,
		Severity:    s.Severity,
		Summary:     s.Summary,
		Description: s.Description,
		Entities:    s.Entities,
		ObservedAt:  s.ObservedAt,
		Metadata:    s.Metadata,
		Tags:        s.Tags,
	}

	data, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("fingerprint marshal: %w", err)
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
