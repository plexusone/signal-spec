# Signal

A Signal is an atomic operational observation from an external system. Signals are the **input layer** - raw events normalized from various sources that will be correlated and mapped to root causes.

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://plexusone.dev/signal-spec/signal.schema.json"
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique signal identifier |
| `type` | [Type](#type-enum) | Signal source category |
| `status` | [Status](#status-enum) | Current processing state |
| `source` | [SourceSystem](common.md#sourcesystem) | Originating system |
| `domain` | [Domain](common.md#domain) | Functional area |
| `severity` | [Severity](common.md#severity) | Impact level |
| `summary` | string | Brief description |
| `observed_at` | datetime | When signal was first observed |
| `received_at` | datetime | When signal was received by system |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Full signal content |
| `entities` | [Entity](common.md#entity)[] | Referenced system components |
| `root_cause_id` | string | Linked root cause, if mapped |
| `fingerprint` | string | Hash for deduplication |
| `embedding` | float32[] | Vector for semantic similarity |
| `metadata` | object | Source-specific additional data |
| `derived` | [DerivedMetrics](#derivedmetrics) | Computed scores excluded from fingerprinting |
| `tags` | [Tag](common.md#tag)[] | User-defined labels (kebab-case) |

## Type Enum

Categorizes the signal source.

| Value | Description |
|-------|-------------|
| `support_ticket` | Customer support ticket |
| `cloud_incident` | Cloud provider incident |
| `security_finding` | Security scan finding |
| `posture_drift` | Configuration drift detection |
| `alert` | Monitoring alert |
| `outage` | Service outage |
| `vulnerability` | Vulnerability scan result |
| `feedback` | Customer feedback |
| `enhancement_request` | Customer feature request |
| `competitive_gap` | Gap vs. a competitor identified from win/loss analysis |
| `competitor_launch` | Competitor product announcement |
| `analyst_finding` | Insight extracted from an analyst report |
| `market_observation` | General market trend |

See [Enhancement Signal Metadata](#enhancement-signal-metadata) and [Cross-Repo References](#cross-repo-references) below for the well-known `metadata` keys used with these product signal types.

## Status Enum

Processing state of the signal.

| Value | Description |
|-------|-------------|
| `new` | Received but not yet processed |
| `processing` | Currently being analyzed |
| `mapped` | Successfully mapped to root cause |
| `ignored` | Determined to be noise/duplicate |
| `archived` | Processed and archived |

## DerivedMetrics

Computed scores kept separate from a signal's identity. These values are recomputed over time and are explicitly excluded from [fingerprinting](#fingerprinting), so identical raw input always produces the same fingerprint regardless of derived values.

| Field | Type | Description |
|-------|------|-------------|
| `frustration` | float | Weighted signal count multiplied by age |
| `momentum` | float | Trailing signal count over a rolling window (e.g., 30 days) |
| `reach` | float | Count of distinct customer references contributing to a root cause |
| `urgency` | float | Case count weighted by severity |
| `computed_at` | datetime | When these metrics were last computed |
| `extra` | object | Additional derived scores not covered by the well-known fields |

```json
{
  "derived": {
    "frustration": 4.2,
    "momentum": 1.8,
    "reach": 12,
    "urgency": 3.5,
    "computed_at": "2026-07-20T00:00:00Z"
  }
}
```

## Enhancement Signal Metadata

`enhancement_request` signals carry structured product data via well-known `metadata` keys. All keys are optional; adapters populate whichever keys their source system provides.

| Metadata Key | Type | Description |
|--------------|------|-------------|
| `votes` | int | Total vote/upvote count |
| `subscribers` | int | Number of watchers/subscribers |
| `organizations` | string[] | Requesting organization names |
| `customers` | string[] | Named customer identifiers |
| `opportunities` | string[] | Sales opportunity IDs linked to this request |
| `estimated_arr` | int64 | Estimated ARR at stake, in cents |

## Cross-Repo References

Signals carry typed references to entities owned by other repositories (e.g., MarketSpec, OrganizationSpec) via well-known `metadata` keys. References use the `{type}:{slug}` format defined by [`pkg/ref`](https://github.com/plexusone/signal-spec/tree/main/pkg/ref).

| Metadata Key | Example |
|--------------|---------|
| `customer_ref` | `customer:acme-001` |
| `capability_ref` | `capability:sso` |
| `market_ref` | `market:identity-governance` |
| `competitor_ref` | `competitor:okta` |
| `analyst_report_ref` | `analyst-report:gartner-mq-iam-2026` |

## Example

```json
{
  "id": "sig-2024-001234",
  "type": "support_ticket",
  "status": "mapped",
  "source": {
    "type": "ticketing",
    "name": "zendesk",
    "external_id": "ZD-98765",
    "url": "https://company.zendesk.com/tickets/98765"
  },
  "domain": {
    "name": "authentication",
    "subdomain": "oauth",
    "team": "identity-platform"
  },
  "severity": "high",
  "summary": "OAuth refresh token failures causing repeated logouts",
  "description": "Customer reports being logged out repeatedly. Error shows 'invalid_grant' when refreshing OAuth token.",
  "entities": [
    {
      "type": "service",
      "name": "oauth-service",
      "attributes": {
        "environment": "production"
      }
    }
  ],
  "observed_at": "2024-03-15T14:30:00Z",
  "received_at": "2024-03-15T14:35:22Z",
  "root_cause_id": "rc-auth-001",
  "fingerprint": "sha256:abc123...",
  "tags": ["enterprise", "mobile", "auth"],
  "metadata": {
    "customer_tier": "enterprise",
    "affected_users": 150
  }
}
```

## Go Usage

```go
import "github.com/plexusone/signal-spec/pkg/signal"

sig := signal.Signal{
    ID:       "sig-2024-001234",
    Type:     signal.TypeSupportTicket,
    Status:   signal.StatusNew,
    Severity: common.SeverityHigh,
    Summary:  "OAuth refresh token failures",
    Domain: common.Domain{
        Name:      "authentication",
        Subdomain: "oauth",
    },
    ObservedAt: time.Now(),
    ReceivedAt: time.Now(),
}
```

## Fingerprinting

`signal.ComputeFingerprint()` returns a deterministic SHA-256 hex digest computed from a signal's identity fields (`ID`, `Type`, `Source`, `Domain`, `Severity`, `Summary`, `Description`, `Entities`, `ObservedAt`, `Metadata`, `Tags`). Mutable and derived fields — `Status`, `Derived`, `Embedding`, `ReceivedAt`, `RootCauseID`, and `Fingerprint` itself — are excluded, so the same raw input always produces the same fingerprint regardless of processing state.

```go
fp, err := signal.ComputeFingerprint(sig)
if err != nil {
    // handle error
}
sig.Fingerprint = fp
```

Use fingerprints to deduplicate signals ingested from the same underlying event across multiple adapters or retries.

## Validation

Signals are validated against these rules:

- `id` is required and must be non-empty
- `type` must be one of the valid enum values
- `summary` is required and must be non-empty
- `domain.name` is required
- `tags` must be lowercase kebab-case

```bash
signal-spec validate -t signal signal.json
```
