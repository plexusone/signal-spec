// Package common provides shared types used across signal-spec entities.
package common

import (
	"fmt"
	"time"

	"github.com/grokify/mogo/text/stringcase"
	"github.com/invopop/jsonschema"
)

// Severity indicates the impact level of a signal or root cause.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityInfo     Severity = "info"
)

// JSONSchema implements jsonschema.Schema for enum generation.
func (Severity) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			SeverityCritical,
			SeverityHigh,
			SeverityMedium,
			SeverityLow,
			SeverityInfo,
		},
	}
}

// Domain represents a functional area or product domain.
// This serves as category/subcategory classification for root causes.
type Domain struct {
	// Name is the category (e.g., "authentication", "payments", "infrastructure").
	Name string `json:"name" jsonschema:"description=Category - the primary functional area"`

	// Subdomain is the subcategory within the domain (e.g., "oauth", "checkout", "kubernetes").
	Subdomain string `json:"subdomain,omitempty" jsonschema:"description=Subcategory within the category"`

	// Team is the owning team for this domain.
	Team string `json:"team,omitempty"`
}

// Tag is a user-defined label in lower-kebab-case format.
// Examples: "enterprise", "auth-related", "p0-incident", "redis-cluster"
type Tag string

// Validate checks if the tag is valid lower-kebab-case.
func (t Tag) Validate() error {
	if !stringcase.IsKebabCase(string(t)) {
		return fmt.Errorf("tag %q is not valid lower-kebab-case", t)
	}
	return nil
}

// ValidateTags checks all tags are valid lower-kebab-case.
func ValidateTags(tags []Tag) error {
	for _, t := range tags {
		if err := t.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TimeRange represents a time interval.
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Entity represents a system component referenced by signals.
type Entity struct {
	// Type is the entity type (e.g., "service", "endpoint", "database").
	Type string `json:"type"`

	// Name is the entity identifier.
	Name string `json:"name"`

	// Attributes contains additional entity metadata.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// SourceSystem identifies where a signal originated.
type SourceSystem struct {
	// Type is the source category (e.g., "ticketing", "alerting", "security").
	Type string `json:"type"`

	// Name is the specific system name (e.g., "zendesk", "pagerduty", "wiz").
	Name string `json:"name"`

	// ExternalID is the original identifier in the source system.
	ExternalID string `json:"external_id,omitempty"`

	// URL is a link to the original item in the source system.
	URL string `json:"url,omitempty"`
}
