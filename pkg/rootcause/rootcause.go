// Package rootcause defines the persistent root cause signal type.
//
// A RootCause is a durable entity representing a systemic operational issue
// that multiple signals map into. Root causes are the primary analytical
// asset - they aggregate evidence, track lifecycle, and enable prioritization.
package rootcause

import (
	"time"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/signal-spec/pkg/common"
)

// Status represents the lifecycle state of a root cause.
type Status string

const (
	StatusNew        Status = "new"
	StatusActive     Status = "active"
	StatusMitigating Status = "mitigating"
	StatusValidating Status = "validating"
	StatusStable     Status = "stable"
	StatusRegressed  Status = "regressed"
	StatusResolved   Status = "resolved"
	StatusArchived   Status = "archived"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Status) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			StatusNew,
			StatusActive,
			StatusMitigating,
			StatusValidating,
			StatusStable,
			StatusRegressed,
			StatusResolved,
			StatusArchived,
		},
	}
}

// TrendDirection indicates whether signal volume is increasing or decreasing.
type TrendDirection string

const (
	TrendIncreasing TrendDirection = "increasing"
	TrendStable     TrendDirection = "stable"
	TrendDecreasing TrendDirection = "decreasing"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (TrendDirection) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			TrendIncreasing,
			TrendStable,
			TrendDecreasing,
		},
	}
}

// Trend captures the temporal behavior of a root cause.
type Trend struct {
	// Direction is the overall trend.
	Direction TrendDirection `json:"direction"`

	// Velocity is the rate of change (signals per day).
	Velocity float64 `json:"velocity"`

	// Period is the time window for trend calculation.
	Period common.TimeRange `json:"period"`
}

// ImpactMetrics quantifies the operational impact of a root cause.
type ImpactMetrics struct {
	// SignalCount is the total number of correlated signals.
	SignalCount int `json:"signal_count"`

	// AffectedCustomers is the count of unique customers impacted.
	AffectedCustomers int `json:"affected_customers,omitempty"`

	// AffectedEntities lists impacted system components.
	AffectedEntities []common.Entity `json:"affected_entities,omitempty"`

	// EscalationRate is the percentage of signals that escalated.
	EscalationRate float64 `json:"escalation_rate,omitempty"`

	// EstimatedRevenueLoss is the projected financial impact.
	EstimatedRevenueLoss float64 `json:"estimated_revenue_loss,omitempty"`
}

// RootCause represents a persistent clustered operational issue.
//
// This is the core analytical entity. Root causes aggregate multiple signals,
// track remediation lifecycle, and enable prioritization based on impact.
// They are managed entities with state transitions, not ad-hoc groupings.
type RootCause struct {
	// ID is the unique root cause identifier.
	ID string `json:"id"`

	// Title is the normalized summary of the issue.
	Title string `json:"title"`

	// Description provides detailed context.
	Description string `json:"description,omitempty"`

	// Status is the current lifecycle state.
	Status Status `json:"status"`

	// Domain is the functional area.
	Domain common.Domain `json:"domain"`

	// Severity indicates current impact level.
	Severity common.Severity `json:"severity"`

	// SymptomPatterns are common manifestations of this root cause.
	SymptomPatterns []string `json:"symptom_patterns,omitempty"`

	// SignalIDs are the correlated signal identifiers.
	SignalIDs []string `json:"signal_ids,omitempty"`

	// Impact contains quantified metrics.
	Impact ImpactMetrics `json:"impact"`

	// Trend captures temporal behavior.
	Trend Trend `json:"trend,omitempty"`

	// PriorityScore is the computed remediation priority (0-100).
	PriorityScore int `json:"priority_score"`

	// FirstSeen is when this root cause was first identified.
	FirstSeen time.Time `json:"first_seen"`

	// LastSeen is when the most recent signal was observed.
	LastSeen time.Time `json:"last_seen"`

	// OwnerTeam is the team responsible for remediation.
	OwnerTeam string `json:"owner_team,omitempty"`

	// RemediationID links to the active remediation effort, if any.
	RemediationID string `json:"remediation_id,omitempty"`

	// RecurrenceCount tracks how many times this issue has regressed.
	RecurrenceCount int `json:"recurrence_count"`

	// Embedding is the vector representation for semantic similarity.
	Embedding []float32 `json:"embedding,omitempty"`

	// Metadata contains additional context.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Tags are user-defined labels in lower-kebab-case format.
	// Examples: "redis", "auth", "enterprise-impact"
	Tags []common.Tag `json:"tags,omitempty" jsonschema:"pattern=^[a-z][a-z0-9]*(-[a-z0-9]+)*$"`
}
