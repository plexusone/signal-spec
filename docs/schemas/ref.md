# Cross-Repo References (`pkg/ref`)

`pkg/ref` defines `TypedRef`, a compact reference format for linking signal-spec entities — signals, `common.Entity` values, and well-known `Metadata` keys — to canonical definitions owned by other repositories in the PlexusOne/ProductBuildersHQ ecosystem (e.g., a market defined in MarketSpec, a customer defined in OrganizationSpec).

Signal Spec does not own these entities; it only carries typed pointers to them. Resolving a `TypedRef` to the full entity is the responsibility of the consuming system.

## TypedRef

`TypedRef` is a string in `{type}:{slug}` format.

```go
type TypedRef string
```

| Example | Meaning |
|---------|---------|
| `market:identity-governance` | A market defined in MarketSpec |
| `customer:acme-001` | A customer defined in OrganizationSpec |
| `capability:sso` | A product capability |
| `competitor:okta` | A tracked competitor |
| `analyst-report:gartner-mq-iam-2026` | An analyst report |

## Format

```
{type}:{slug}
```

- `{type}` is a well-known [`RefType`](#reftype) (see below).
- `{slug}` is a lowercase alphanumeric string with hyphens, matching `^[a-z0-9]+(-[a-z0-9]+)*$`. It must not be empty, and must not start or end with a hyphen.
- The first `:` separates type from slug; slugs themselves never contain `:`.

## RefType

`RefType` is a well-known entity reference type.

| Constant | Value | Description |
|----------|-------|--------------|
| `ref.TypeCustomer` | `customer` | Customer entity (OrganizationSpec) |
| `ref.TypeCapability` | `capability` | Product capability |
| `ref.TypeMarket` | `market` | Market entity (MarketSpec) |
| `ref.TypeCompetitor` | `competitor` | Competitor entity |
| `ref.TypeAnalystReport` | `analyst-report` | Analyst report entity |

```go
import "github.com/plexusone/signal-spec/pkg/ref"

types := ref.KnownTypes()
// []ref.RefType{TypeCustomer, TypeCapability, TypeMarket, TypeCompetitor, TypeAnalystReport}
```

## Functions

### New

Constructs a `TypedRef` from a type and slug.

```go
r := ref.New(ref.TypeMarket, "identity-governance")
// r == "market:identity-governance"
```

### Parse

Splits a `TypedRef` into its type and slug components. Returns an error if the format is invalid (missing `:` or empty slug). `Parse` does not check whether the type is a known type — use [`Validate`](#validate) for that.

```go
typ, slug, err := ref.Parse("market:identity-governance")
// typ == ref.TypeMarket, slug == "identity-governance", err == nil

_, _, err = ref.Parse("no-colon-here")
// err != nil
```

### Validate

Checks that a `TypedRef` is well-formed (via `Parse`) and uses a known `RefType`. Does **not** validate the slug's character set.

```go
err := ref.Validate("market:identity-governance") // nil
err  = ref.Validate("unknown:something")           // error: unknown ref type
```

### ValidateSlug

Checks that a slug contains only lowercase alphanumeric characters and hyphens, and does not start or end with a hyphen.

```go
err := ref.ValidateSlug("identity-governance") // nil
err  = ref.ValidateSlug("UPPERCASE")            // error
err  = ref.ValidateSlug("-leading")             // error
err  = ref.ValidateSlug("has space")            // error
```

### ValidateStrict

Combines format, known-type, and slug validation — the strictest check, recommended before persisting or emitting a `TypedRef`.

```go
err := ref.ValidateStrict("market:identity-governance") // nil
err  = ref.ValidateStrict("market:UPPER")                // error: invalid slug
err  = ref.ValidateStrict("unknown:valid-slug")          // error: unknown type
```

## Usage in Signal Spec

`TypedRef` values appear in two places in the signal-spec data model:

### `common.Entity.Ref`

`Entity` has an optional `Ref` field carrying a typed cross-repo reference, so an entity referenced by a signal or root cause can be linked directly to its canonical definition:

```go
entity := common.Entity{
    Type: "customer",
    Name: "Acme Corp",
    Ref:  "customer:acme-001",
}
```

```json
{
  "type": "customer",
  "name": "Acme Corp",
  "ref": "customer:acme-001"
}
```

### Well-known `Signal.Metadata` keys

Signals carry typed references via well-known `Metadata` keys, defined as `signal.Meta*` constants:

| Constant | Metadata Key | Example |
|----------|--------------|---------|
| `signal.MetaCustomerRef` | `customer_ref` | `customer:acme-001` |
| `signal.MetaCapabilityRef` | `capability_ref` | `capability:sso` |
| `signal.MetaMarketRef` | `market_ref` | `market:identity-governance` |
| `signal.MetaCompetitorRef` | `competitor_ref` | `competitor:okta` |
| `signal.MetaAnalystReportRef` | `analyst_report_ref` | `analyst-report:gartner-mq-iam-2026` |

```go
sig.Metadata = map[string]any{
    signal.MetaMarketRef:     string(ref.New(ref.TypeMarket, "identity-governance")),
    signal.MetaCompetitorRef: string(ref.New(ref.TypeCompetitor, "okta")),
}
```

See [Signal: Cross-Repo References](signal.md#cross-repo-references) for the full list of metadata keys and how they relate to the [product signal types](signal.md#type-enum).

## Full Example

```go
package main

import (
    "fmt"

    "github.com/plexusone/signal-spec/pkg/ref"
)

func main() {
    r := ref.New(ref.TypeMarket, "identity-governance")

    if err := ref.ValidateStrict(r); err != nil {
        panic(err)
    }

    typ, slug, err := ref.Parse(r)
    if err != nil {
        panic(err)
    }

    fmt.Printf("type=%s slug=%s\n", typ, slug)
    // type=market slug=identity-governance
}
```

## Validation Summary

| Check | `Parse` | `Validate` | `ValidateSlug` | `ValidateStrict` |
|-------|:-------:|:----------:|:--------------:|:-----------------:|
| Format (`{type}:{slug}`, non-empty slug) | Yes | Yes | - | Yes |
| Known `RefType` | No | Yes | - | Yes |
| Slug character set / hyphen rules | No | No | Yes | Yes |
