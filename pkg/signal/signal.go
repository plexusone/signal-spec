// Package signal defines the raw signal type.
//
// A Signal is an atomic observation from an external system — operational
// (incidents, alerts, tickets) or product-oriented (enhancement requests,
// feedback). Signals are the input layer of the intelligence pipeline:
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
	TypeSupportTicket      Type = "support_ticket"
	TypeCloudIncident      Type = "cloud_incident"
	TypeSecurityFinding    Type = "security_finding"
	TypePostureDrift       Type = "posture_drift"
	TypeAlert              Type = "alert"
	TypeOutage             Type = "outage"
	TypeVulnerability      Type = "vulnerability"
	TypeFeedback           Type = "feedback"
	TypeEnhancementRequest Type = "enhancement_request"
	TypeCompetitiveGap     Type = "competitive_gap"
	TypeCompetitorLaunch   Type = "competitor_launch"
	TypeAnalystFinding     Type = "analyst_finding"
	TypeMarketObservation  Type = "market_observation"
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
		TypeEnhancementRequest,
		TypeCompetitiveGap,
		TypeCompetitorLaunch,
		TypeAnalystFinding,
		TypeMarketObservation,
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
			TypeEnhancementRequest,
			TypeCompetitiveGap,
			TypeCompetitorLaunch,
			TypeAnalystFinding,
			TypeMarketObservation,
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

// Enhancement request metadata key conventions.
// Signals of type enhancement_request use these well-known metadata keys
// for structured product signal data. All keys are optional; adapters
// populate whichever keys their source system provides.
const (
	// MetaVotes is the total vote/upvote count (int).
	MetaVotes = "votes"
	// MetaSubscribers is the number of watchers/subscribers (int).
	MetaSubscribers = "subscribers"
	// MetaOrganizations is a list of requesting organization names ([]string).
	MetaOrganizations = "organizations"
	// MetaCustomers is a list of named customer identifiers ([]string).
	MetaCustomers = "customers"
	// MetaOpportunities is a list of sales opportunity IDs linked to this request ([]string).
	MetaOpportunities = "opportunities"
	// MetaEstimatedARR is the estimated annual recurring revenue at stake, in cents (int64).
	MetaEstimatedARR = "estimated_arr"
)

// Cross-repo reference metadata key conventions.
// Signals use these well-known metadata keys to carry typed entity
// references in "{type}:{slug}" format (see pkg/ref for validation).
const (
	// MetaCustomerRef is a typed ref to a customer entity (e.g., "customer:acme-001").
	MetaCustomerRef = "customer_ref"
	// MetaCapabilityRef is a typed ref to a product capability (e.g., "capability:sso").
	MetaCapabilityRef = "capability_ref"
	// MetaMarketRef is a typed ref to a market (e.g., "market:identity-governance").
	MetaMarketRef = "market_ref"
	// MetaCompetitorRef is a typed ref to a competitor (e.g., "competitor:okta").
	MetaCompetitorRef = "competitor_ref"
	// MetaAnalystReportRef is a typed ref to an analyst report (e.g., "analyst-report:gartner-mq-iam-2026").
	MetaAnalystReportRef = "analyst_report_ref"
)

// DerivedMetrics holds computed scores that are NOT part of the signal's
// identity. These values are recomputed from raw data and must be excluded
// from fingerprinting — identical raw input must always produce an identical
// fingerprint regardless of derived values.
type DerivedMetrics struct {
	// Frustration is the weighted signal count multiplied by age.
	Frustration *float64 `json:"frustration,omitempty"`

	// Momentum is the trailing signal count over a rolling window (e.g., 30 days).
	Momentum *float64 `json:"momentum,omitempty"`

	// Reach is the count of distinct customer references contributing to a root cause.
	Reach *float64 `json:"reach,omitempty"`

	// Urgency is the case count weighted by severity.
	Urgency *float64 `json:"urgency,omitempty"`

	// ComputedAt is when these metrics were last computed.
	ComputedAt *time.Time `json:"computed_at,omitempty"`

	// Extra holds additional derived scores not covered by the well-known fields.
	Extra map[string]float64 `json:"extra,omitempty"`
}

// Signal represents an atomic observation from an external system.
//
// This is the raw input to the intelligence pipeline. Signals are normalized
// from various sources (tickets, alerts, incidents, enhancement requests)
// into a canonical format for correlation and root cause mapping.
//
// Fingerprinting contract: the Fingerprint field is computed from raw
// signal data only (Type, Source, Domain, Summary, Metadata, etc.).
// The Derived field is explicitly excluded — it holds recomputed scores
// that vary over time without changing the signal's identity.
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

	// Derived holds computed metrics excluded from fingerprinting.
	Derived *DerivedMetrics `json:"derived,omitempty"`

	// Tags are user-defined labels in lower-kebab-case format.
	// Examples: "enterprise", "mobile", "auth-failure"
	Tags []common.Tag `json:"tags,omitempty" jsonschema:"pattern=^[a-z][a-z0-9]*(-[a-z0-9]+)*$"`
}
