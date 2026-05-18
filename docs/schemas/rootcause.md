# RootCause

A RootCause is a persistent clustered operational issue. Root causes are the **primary analytical asset** - durable entities that aggregate evidence, track lifecycle, and enable prioritization based on impact.

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://plexusone.dev/signal-spec/rootcause.schema.json"
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique root cause identifier |
| `title` | string | Normalized summary of the issue |
| `status` | [Status](#status-enum) | Current lifecycle state |
| `domain` | [Domain](common.md#domain) | Functional area |
| `severity` | [Severity](common.md#severity) | Current impact level |
| `impact` | [ImpactMetrics](#impactmetrics) | Quantified impact |
| `priority_score` | integer | Computed priority (0-100) |
| `first_seen` | datetime | When first identified |
| `last_seen` | datetime | Most recent signal |
| `recurrence_count` | integer | Times issue regressed |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Detailed context |
| `symptom_patterns` | string[] | Common manifestations |
| `signal_ids` | string[] | Correlated signal IDs |
| `trend` | [Trend](#trend) | Temporal behavior |
| `owner_team` | string | Responsible team |
| `remediation_id` | string | Active remediation |
| `embedding` | float32[] | Vector for similarity |
| `metadata` | object | Additional context |
| `tags` | [Tag](common.md#tag)[] | User-defined labels |

## Status Enum

Lifecycle state of the root cause.

| Value | Description |
|-------|-------------|
| `new` | Recently identified |
| `active` | Confirmed, requiring attention |
| `mitigating` | Remediation in progress |
| `validating` | Fix deployed, measuring efficacy |
| `stable` | Fix appears effective |
| `regressed` | Issue recurred after stabilization |
| `resolved` | Confirmed fixed |
| `archived` | Closed and archived |

See [Lifecycle States](../concepts/lifecycle.md) for state transitions.

## ImpactMetrics

Quantifies operational impact.

| Field | Type | Description |
|-------|------|-------------|
| `signal_count` | integer | Total correlated signals |
| `affected_customers` | integer | Unique customers impacted |
| `affected_entities` | Entity[] | Impacted components |
| `escalation_rate` | float | Percentage escalated |
| `estimated_revenue_loss` | float | Projected financial impact |

## Trend

Captures temporal behavior.

| Field | Type | Description |
|-------|------|-------------|
| `direction` | TrendDirection | Overall trend |
| `velocity` | float | Rate of change (signals/day) |
| `period` | TimeRange | Measurement window |

### TrendDirection Enum

| Value | Description |
|-------|-------------|
| `increasing` | Signal volume growing |
| `stable` | Signal volume constant |
| `decreasing` | Signal volume declining |

## Example

```json
{
  "id": "rc-auth-001",
  "title": "Redis session replication instability causing OAuth token validation failures",
  "description": "Intermittent Redis cluster replication lag causes session state inconsistency.",
  "status": "mitigating",
  "domain": {
    "name": "authentication",
    "subdomain": "oauth",
    "team": "identity-platform"
  },
  "severity": "high",
  "symptom_patterns": [
    "OAuth refresh token failures",
    "Repeated logout events",
    "invalid_grant errors"
  ],
  "signal_ids": [
    "sig-2024-001234",
    "sig-2024-001235"
  ],
  "impact": {
    "signal_count": 487,
    "affected_customers": 2341,
    "affected_entities": [
      {"type": "service", "name": "oauth-service"},
      {"type": "service", "name": "session-store"}
    ],
    "escalation_rate": 0.12,
    "estimated_revenue_loss": 45000
  },
  "trend": {
    "direction": "stable",
    "velocity": 15.3,
    "period": {
      "start": "2024-03-01T00:00:00Z",
      "end": "2024-03-15T23:59:59Z"
    }
  },
  "priority_score": 87,
  "first_seen": "2024-02-28T10:15:00Z",
  "last_seen": "2024-03-15T14:30:00Z",
  "owner_team": "identity-platform",
  "remediation_id": "rem-001",
  "recurrence_count": 1,
  "tags": ["redis", "auth", "session", "enterprise-impact"]
}
```

## Go Usage

```go
import "github.com/plexusone/signal-spec/pkg/rootcause"

rc := rootcause.RootCause{
    ID:     "rc-auth-001",
    Title:  "Redis session replication instability",
    Status: rootcause.StatusActive,
    Domain: common.Domain{
        Name:      "authentication",
        Subdomain: "oauth",
    },
    Severity: common.SeverityHigh,
    Impact: rootcause.ImpactMetrics{
        SignalCount:       487,
        AffectedCustomers: 2341,
    },
    PriorityScore:   87,
    FirstSeen:       time.Now().Add(-14 * 24 * time.Hour),
    LastSeen:        time.Now(),
    RecurrenceCount: 0,
}
```

## Validation

Root causes are validated against these rules:

- `id` is required and must be non-empty
- `title` is required and must be non-empty
- `domain.name` is required
- `tags` must be lowercase kebab-case

```bash
signal-spec validate -t rootcause rootcause.json
```
