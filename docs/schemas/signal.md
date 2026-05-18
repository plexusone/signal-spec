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

## Status Enum

Processing state of the signal.

| Value | Description |
|-------|-------------|
| `new` | Received but not yet processed |
| `processing` | Currently being analyzed |
| `mapped` | Successfully mapped to root cause |
| `ignored` | Determined to be noise/duplicate |
| `archived` | Processed and archived |

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
