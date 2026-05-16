// Package remediation defines types for tracking corrective actions.
//
// A Remediation represents an effort to fix a root cause. It tracks
// implementation, validation, and efficacy - enabling closed-loop
// measurement of whether fixes actually reduced operational pain.
package remediation

import (
	"time"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/signal-spec/pkg/common"
)

// Status represents the lifecycle state of a remediation.
type Status string

const (
	StatusProposed    Status = "proposed"
	StatusApproved    Status = "approved"
	StatusInProgress  Status = "in_progress"
	StatusDeployed    Status = "deployed"
	StatusValidating  Status = "validating"
	StatusEffective   Status = "effective"
	StatusIneffective Status = "ineffective"
	StatusRolledBack  Status = "rolled_back"
	StatusCancelled   Status = "cancelled"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Status) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			StatusProposed,
			StatusApproved,
			StatusInProgress,
			StatusDeployed,
			StatusValidating,
			StatusEffective,
			StatusIneffective,
			StatusRolledBack,
			StatusCancelled,
		},
	}
}

// Efficacy captures the measured effectiveness of a remediation.
type Efficacy struct {
	// SignalReduction is the percentage decrease in signal volume.
	SignalReduction float64 `json:"signal_reduction"`

	// ValidationPeriod is the time window for measurement.
	ValidationPeriod common.TimeRange `json:"validation_period"`

	// ConfidenceLevel indicates statistical confidence (0-1).
	ConfidenceLevel float64 `json:"confidence_level"`

	// Notes provides context on the measurement.
	Notes string `json:"notes,omitempty"`
}

// Remediation represents a corrective action for a root cause.
//
// This enables closed-loop validation: track whether fixes actually
// reduced signal volume and detect regressions.
type Remediation struct {
	// ID is the unique remediation identifier.
	ID string `json:"id"`

	// Title summarizes the fix.
	Title string `json:"title"`

	// Description provides implementation details.
	Description string `json:"description,omitempty"`

	// Status is the current lifecycle state.
	Status Status `json:"status"`

	// RootCauseIDs are the targeted root causes.
	RootCauseIDs []string `json:"root_cause_ids"`

	// OwnerTeam is responsible for implementation.
	OwnerTeam string `json:"owner_team"`

	// Assignee is the individual owner.
	Assignee string `json:"assignee,omitempty"`

	// CreatedAt is when the remediation was proposed.
	CreatedAt time.Time `json:"created_at"`

	// DeployedAt is when the fix was deployed.
	DeployedAt *time.Time `json:"deployed_at,omitempty"`

	// ValidatedAt is when efficacy was measured.
	ValidatedAt *time.Time `json:"validated_at,omitempty"`

	// Efficacy contains measured effectiveness.
	Efficacy *Efficacy `json:"efficacy,omitempty"`

	// ExternalLinks references related items (PRs, tickets, etc.).
	ExternalLinks []common.SourceSystem `json:"external_links,omitempty"`

	// Metadata contains additional context.
	Metadata map[string]any `json:"metadata,omitempty"`

	// Tags are user-defined labels in lower-kebab-case format.
	Tags []common.Tag `json:"tags,omitempty" jsonschema:"pattern=^[a-z][a-z0-9]*(-[a-z0-9]+)*$"`
}

// ValidationSignal represents evidence of remediation effectiveness.
//
// These are generated after remediation deployment to track whether
// the fix worked. They detect both success (signal decay) and
// regression (signal resurgence).
type ValidationSignal struct {
	// ID is the unique validation signal identifier.
	ID string `json:"id"`

	// RemediationID links to the remediation being validated.
	RemediationID string `json:"remediation_id"`

	// RootCauseID links to the root cause.
	RootCauseID string `json:"root_cause_id"`

	// Type indicates the validation result.
	Type ValidationResult `json:"type"`

	// ObservedAt is when the validation occurred.
	ObservedAt time.Time `json:"observed_at"`

	// BaselineSignalRate is the pre-remediation signal rate.
	BaselineSignalRate float64 `json:"baseline_signal_rate"`

	// CurrentSignalRate is the post-remediation signal rate.
	CurrentSignalRate float64 `json:"current_signal_rate"`

	// ReductionPercent is the calculated improvement.
	ReductionPercent float64 `json:"reduction_percent"`

	// Notes provides context.
	Notes string `json:"notes,omitempty"`
}

// ValidationResult indicates the outcome of validation.
type ValidationResult string

const (
	ValidationImproved  ValidationResult = "improved"
	ValidationNoChange  ValidationResult = "no_change"
	ValidationRegressed ValidationResult = "regressed"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (ValidationResult) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			ValidationImproved,
			ValidationNoChange,
			ValidationRegressed,
		},
	}
}
