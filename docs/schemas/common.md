# Common Types

Shared types used across Signal, RootCause, and Remediation entities.

## Severity

Impact level indicator.

| Value | Description | Typical Response |
|-------|-------------|------------------|
| `critical` | Service down, major customer impact | Immediate response required |
| `high` | Significant degradation | Response within hours |
| `medium` | Noticeable impact, workarounds exist | Response within days |
| `low` | Minor impact | Normal prioritization |
| `info` | Informational, no action needed | Review only |

```json
"severity": "high"
```

## Domain

Functional area classification using category/subcategory taxonomy.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Category (e.g., "authentication") |
| `subdomain` | string | No | Subcategory (e.g., "oauth") |
| `team` | string | No | Owning team |

```json
{
  "domain": {
    "name": "authentication",
    "subdomain": "oauth",
    "team": "identity-platform"
  }
}
```

See [Domains & Taxonomy](../concepts/domains.md) for guidelines.

## Entity

A system component referenced by signals or root causes.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | Entity type |
| `name` | string | Yes | Entity identifier |
| `ref` | string | No | Typed cross-repo reference in `{type}:{slug}` format (e.g., `customer:acme-001`, `market:identity-governance`); see [`pkg/ref`](https://github.com/plexusone/signal-spec/tree/main/pkg/ref) for known types |
| `attributes` | object | No | Additional metadata |

### Common Entity Types

| Type | Examples |
|------|----------|
| `service` | oauth-service, payment-api |
| `endpoint` | /api/v1/users, /auth/token |
| `database` | users-db, analytics-cluster |
| `queue` | order-events, notifications |
| `host` | prod-web-01, k8s-node-03 |
| `application` | mobile-app-ios, web-portal |

```json
{
  "entities": [
    {
      "type": "service",
      "name": "oauth-service",
      "attributes": {
        "environment": "production",
        "region": "us-east-1"
      }
    }
  ]
}
```

## SourceSystem

Identifies the origin of a signal or external reference.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | Source category |
| `name` | string | Yes | System name |
| `external_id` | string | No | Original ID |
| `url` | string | No | Link to source |

### Common Source Types

| Type | Examples |
|------|----------|
| `ticketing` | zendesk, jira, freshdesk |
| `alerting` | pagerduty, opsgenie, datadog |
| `security` | wiz, snyk, crowdstrike |
| `monitoring` | prometheus, grafana, newrelic |
| `code_change` | github, gitlab, bitbucket |
| `incident` | statuspage, incident.io |

```json
{
  "source": {
    "type": "ticketing",
    "name": "zendesk",
    "external_id": "ZD-98765",
    "url": "https://company.zendesk.com/tickets/98765"
  }
}
```

## Tag

User-defined labels in lowercase kebab-case format.

**Format:** `^[a-z][a-z0-9]*(-[a-z0-9]+)*$`

### Valid Examples

- `enterprise`
- `auth-related`
- `p0-incident`
- `redis-cluster`
- `q4-2024`

### Invalid Examples

- `Enterprise` (uppercase)
- `auth_related` (underscore)
- `P0-incident` (uppercase)
- `-invalid` (leading hyphen)
- `123start` (starts with number)

```json
{
  "tags": ["enterprise", "mobile", "auth-failure"]
}
```

### Validation

Tags are validated using the `mogo/text/stringcase.IsKebabCase()` function.

```go
import "github.com/plexusone/signal-spec/pkg/common"

tag := common.Tag("my-tag")
if err := tag.Validate(); err != nil {
    // Handle invalid tag
}
```

## TimeRange

A time interval for trends and validation periods.

| Field | Type | Description |
|-------|------|-------------|
| `start` | datetime | Start of period |
| `end` | datetime | End of period |

```json
{
  "period": {
    "start": "2024-03-01T00:00:00Z",
    "end": "2024-03-15T23:59:59Z"
  }
}
```

## Datetime Format

All datetime fields use ISO 8601 format with timezone:

```
YYYY-MM-DDTHH:MM:SSZ
```

Examples:

- `2024-03-15T14:30:00Z` (UTC)
- `2024-03-15T10:30:00-04:00` (EDT)
