package ref

import (
	"testing"
)

func TestNew(t *testing.T) {
	r := New(TypeMarket, "identity-governance")
	if r != "market:identity-governance" {
		t.Errorf("expected market:identity-governance, got %s", r)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input   TypedRef
		typ     RefType
		slug    string
		wantErr bool
	}{
		{"market:identity-governance", TypeMarket, "identity-governance", false},
		{"customer:acme-001", TypeCustomer, "acme-001", false},
		{"capability:sso", TypeCapability, "sso", false},
		{"competitor:okta", TypeCompetitor, "okta", false},
		{"analyst-report:gartner-mq-iam-2026", TypeAnalystReport, "gartner-mq-iam-2026", false},
		{"nocolon", "", "", true},
		{":no-type", "", "", true},
		{"market:", "", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		typ, slug, err := Parse(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Parse(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			if typ != tt.typ {
				t.Errorf("Parse(%q): type=%s, want=%s", tt.input, typ, tt.typ)
			}
			if slug != tt.slug {
				t.Errorf("Parse(%q): slug=%s, want=%s", tt.input, slug, tt.slug)
			}
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		input   TypedRef
		wantErr bool
	}{
		{"market:identity-governance", false},
		{"customer:acme-001", false},
		{"capability:sso", false},
		{"competitor:okta", false},
		{"analyst-report:gartner-mq-2026", false},
		{"unknown:something", true},
		{"nocolon", true},
		{":no-type", true},
		{"market:", true},
	}

	for _, tt := range tests {
		err := Validate(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Validate(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidateSlug(t *testing.T) {
	tests := []struct {
		slug    string
		wantErr bool
	}{
		{"identity-governance", false},
		{"acme-001", false},
		{"sso", false},
		{"a", false},
		{"abc123", false},
		{"", true},
		{"-leading", true},
		{"trailing-", true},
		{"has space", true},
		{"UPPERCASE", true},
		{"under_score", true},
	}

	for _, tt := range tests {
		err := ValidateSlug(tt.slug)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateSlug(%q): err=%v, wantErr=%v", tt.slug, err, tt.wantErr)
		}
	}
}

func TestValidateStrict(t *testing.T) {
	tests := []struct {
		input   TypedRef
		wantErr bool
	}{
		{"market:identity-governance", false},
		{"customer:acme-001", false},
		{"market:UPPER", true},
		{"market:has space", true},
		{"market:-leading", true},
		{"unknown:valid-slug", true},
		{"nocolon", true},
	}

	for _, tt := range tests {
		err := ValidateStrict(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateStrict(%q): err=%v, wantErr=%v", tt.input, err, tt.wantErr)
		}
	}
}

func TestKnownTypes(t *testing.T) {
	types := KnownTypes()
	if len(types) != 5 {
		t.Errorf("expected 5 known types, got %d", len(types))
	}
}

func TestRoundTrip(t *testing.T) {
	for _, typ := range KnownTypes() {
		r := New(typ, "test-slug")
		parsedType, parsedSlug, err := Parse(r)
		if err != nil {
			t.Fatalf("Parse(New(%s, test-slug)): %v", typ, err)
		}
		if parsedType != typ {
			t.Errorf("round-trip type: got %s, want %s", parsedType, typ)
		}
		if parsedSlug != "test-slug" {
			t.Errorf("round-trip slug: got %s, want test-slug", parsedSlug)
		}
		if err := Validate(r); err != nil {
			t.Errorf("Validate(New(%s, test-slug)): %v", typ, err)
		}
	}
}
