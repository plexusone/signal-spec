// Package ref defines typed cross-repo entity references.
//
// A TypedRef uses the format "{type}:{slug}" to reference entities across
// repository boundaries. For example, "market:identity-governance" references
// a market defined in MarketSpec, and "customer:acme-001" references a
// customer entity in OrganizationSpec.
package ref

import (
	"fmt"
	"strings"
)

// TypedRef is a cross-repo entity reference in "{type}:{slug}" format.
type TypedRef string

// RefType is a well-known entity reference type.
type RefType string

const (
	TypeCustomer      RefType = "customer"
	TypeCapability    RefType = "capability"
	TypeMarket        RefType = "market"
	TypeCompetitor    RefType = "competitor"
	TypeAnalystReport RefType = "analyst-report"
)

// KnownTypes returns all well-known reference types.
func KnownTypes() []RefType {
	return []RefType{
		TypeCustomer,
		TypeCapability,
		TypeMarket,
		TypeCompetitor,
		TypeAnalystReport,
	}
}

func isKnownType(t RefType) bool {
	for _, k := range KnownTypes() {
		if k == t {
			return true
		}
	}
	return false
}

// New creates a TypedRef from a type and slug.
func New(t RefType, slug string) TypedRef {
	return TypedRef(fmt.Sprintf("%s:%s", t, slug))
}

// Parse splits a TypedRef into its type and slug components.
// Returns an error if the format is invalid.
func Parse(r TypedRef) (RefType, string, error) {
	s := string(r)
	idx := strings.IndexByte(s, ':')
	if idx < 1 {
		return "", "", fmt.Errorf("invalid typed ref %q: missing ':'", s)
	}
	typ := s[:idx]
	slug := s[idx+1:]
	if slug == "" {
		return "", "", fmt.Errorf("invalid typed ref %q: empty slug", s)
	}
	return RefType(typ), slug, nil
}

// Validate checks that a TypedRef is well-formed and uses a known type.
func Validate(r TypedRef) error {
	t, _, err := Parse(r)
	if err != nil {
		return err
	}
	if !isKnownType(t) {
		return fmt.Errorf("unknown ref type %q in %q: expected one of %v", t, r, KnownTypes())
	}
	return nil
}

// ValidateSlug checks that a slug contains only lowercase alphanumeric
// characters and hyphens, and does not start or end with a hyphen.
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug must not be empty")
	}
	if slug[0] == '-' || slug[len(slug)-1] == '-' {
		return fmt.Errorf("slug %q must not start or end with a hyphen", slug)
	}
	for _, c := range slug {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return fmt.Errorf("slug %q contains invalid character %q: only lowercase alphanumeric and hyphens allowed", slug, string(c))
		}
	}
	return nil
}

// ValidateStrict checks format, known type, and slug validity.
func ValidateStrict(r TypedRef) error {
	t, slug, err := Parse(r)
	if err != nil {
		return err
	}
	if !isKnownType(t) {
		return fmt.Errorf("unknown ref type %q in %q: expected one of %v", t, r, KnownTypes())
	}
	if err := ValidateSlug(slug); err != nil {
		return fmt.Errorf("invalid typed ref %q: %w", r, err)
	}
	return nil
}
