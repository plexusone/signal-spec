// Package signal defines the raw operational signal type.
//
// A Signal is an atomic operational observation from an external system.
// Signals are the input layer of the intelligence pipeline - they represent
// raw events that will be correlated and mapped to root causes.
package signal

import (
	"time"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/signal-spec/pkg/common"
)

// Type categorizes the kind of operational signal.
type Type string

const (
	TypeSupportTicket   Type = "support_ticket"
	TypeCloudIncident   Type = "cloud_incident"
	TypeSecurityFinding Type = "security_finding"
	TypePostureDrift    Type = "posture_drift"
	TypeAlert           Type = "alert"
	TypeOutage          Type = "outage"
	TypeVulnerability   Type = "vulnerability"
	TypeFeedback        Type = "feedback"
)

// TypeValues returns all valid Type values.
func TypeValues() []Type {
	return []Type{
		TypeSupportTicket,
		TypeCloudIncident,
		TypeSecurityFinding,
		TypePostureDrift,
		TypeAlert,
		TypeOutage,
		TypeVulnerability,
		TypeFeedback,
	}
}

// JSONSchema implements jsonschema.Schema for enum generation.
func (Type) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			TypeSupportTicket,
			TypeCloudIncident,
			TypeSecurityFinding,
			TypePostureDrift,
			TypeAlert,
			TypeOutage,
			TypeVulnerability,
			TypeFeedback,
		},
	}
}

// Status indicates the current state of a signal.
type Status string

const (
	StatusNew        Status = "new"
	StatusProcessing Status = "processing"
	StatusMapped     Status = "mapped"
	StatusIgnored    Status = "ignored"
	StatusArchived   Status = "archived"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Status) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			StatusNew,
			StatusProcessing,
			StatusMapped,
			StatusIgnored,
			StatusArchived,
		},
	}
}

// Signal represents an atomic operational observation from an external system.
//
// This is the raw input to the intelligence pipeline. Signals are normalized
// from various sources (tickets, alerts, incidents) into a canonical format
// for correlation and root cause mapping.
type Signal struct {
	// ID is the unique signal identifier.
	ID string `json:"id"`

	// Type categorizes the signal source.
	Type Type `json:"type"`

	// Status is the current processing state.
	Status Status `json:"status"`

	// Source identifies the originating system.
	Source common.SourceSystem `json:"source"`

	// Domain is the functional area this signal relates to.
	Domain common.Domain `json:"domain"`

	// Severity indicates impact level.
	Severity common.Severity `json:"severity"`

	// Summary is a brief description of the signal.
	Summary string `json:"summary"`

	// Description is the full signal content.
	Description string `json:"description,omitempty"`

	// Entities are system components referenced by this signal.
	Entities []common.Entity `json:"entities,omitempty"`

	// ObservedAt is when the signal was first observed.
	ObservedAt time.Time `json:"observed_at"`

	// ReceivedAt is when the signal was received by the system.
	ReceivedAt time.Time `json:"received_at"`

	// RootCauseID links to the mapped root cause, if any.
	RootCauseID string `json:"root_cause_id,omitempty"`

	// Fingerprint is a hash for deduplication.
	Fingerprint string `json:"fingerprint,omitempty"`

	// Embedding is the vector representation for semantic similarity.
	Embedding []float32 `json:"embedding,omitempty"`

	// Metadata contains source-specific additional data.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Tags are user-defined labels in lower-kebab-case format.
	// Examples: "enterprise", "mobile", "auth-failure"
	Tags []common.Tag `json:"tags,omitempty" jsonschema:"pattern=^[a-z][a-z0-9]*(-[a-z0-9]+)*$"`
}
