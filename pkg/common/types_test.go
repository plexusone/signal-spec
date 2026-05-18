package common

import (
	"testing"
)

func TestSeverityJSONSchema(t *testing.T) {
	schema := Severity("").JSONSchema()
	if schema.Type != "string" {
		t.Errorf("expected type string, got %s", schema.Type)
	}
	if len(schema.Enum) != 5 {
		t.Errorf("expected 5 enum values, got %d", len(schema.Enum))
	}
}

func TestTagValidate(t *testing.T) {
	tests := []struct {
		tag     Tag
		wantErr bool
	}{
		{"valid-tag", false},
		{"auth", false},
		{"p0-incident", false},
		{"redis-cluster-issue", false},
		{"a1b2c3", false},
		{"Invalid", true},       // uppercase
		{"invalid_tag", true},   // underscore
		{"invalid tag", true},   // space
		{"-invalid", true},      // starts with hyphen
		{"123invalid", true},    // starts with number
		{"UPPERCASE", true},     // all uppercase
		{"camelCase", true},     // camelCase
	}

	for _, tt := range tests {
		t.Run(string(tt.tag), func(t *testing.T) {
			err := tt.tag.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Tag(%q).Validate() error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTags(t *testing.T) {
	validTags := []Tag{"auth", "redis", "p0-incident"}
	if err := ValidateTags(validTags); err != nil {
		t.Errorf("ValidateTags(%v) unexpected error: %v", validTags, err)
	}

	invalidTags := []Tag{"valid", "Invalid", "also-valid"}
	if err := ValidateTags(invalidTags); err == nil {
		t.Errorf("ValidateTags(%v) expected error, got nil", invalidTags)
	}
}

func TestDomainFields(t *testing.T) {
	d := Domain{
		Name:      "authentication",
		Subdomain: "oauth",
		Team:      "identity-platform",
	}

	if d.Name != "authentication" {
		t.Errorf("expected name 'authentication', got %s", d.Name)
	}
	if d.Subdomain != "oauth" {
		t.Errorf("expected subdomain 'oauth', got %s", d.Subdomain)
	}
	if d.Team != "identity-platform" {
		t.Errorf("expected team 'identity-platform', got %s", d.Team)
	}
}

func TestEntityFields(t *testing.T) {
	e := Entity{
		Type: "service",
		Name: "oauth-service",
		Attributes: map[string]string{
			"environment": "production",
		},
	}

	if e.Type != "service" {
		t.Errorf("expected type 'service', got %s", e.Type)
	}
	if e.Name != "oauth-service" {
		t.Errorf("expected name 'oauth-service', got %s", e.Name)
	}
	if e.Attributes["environment"] != "production" {
		t.Errorf("expected attribute 'production', got %s", e.Attributes["environment"])
	}
}

func TestSourceSystemFields(t *testing.T) {
	s := SourceSystem{
		Type:       "ticketing",
		Name:       "zendesk",
		ExternalID: "ZD-12345",
		URL:        "https://example.zendesk.com/tickets/12345",
	}

	if s.Type != "ticketing" {
		t.Errorf("expected type 'ticketing', got %s", s.Type)
	}
	if s.Name != "zendesk" {
		t.Errorf("expected name 'zendesk', got %s", s.Name)
	}
	if s.ExternalID != "ZD-12345" {
		t.Errorf("expected external_id 'ZD-12345', got %s", s.ExternalID)
	}
}
