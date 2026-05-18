# Remediation

A Remediation represents a corrective action targeting one or more root causes. Remediations enable **closed-loop validation** - measuring whether fixes actually reduced operational pain.

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://plexusone.dev/signal-spec/remediation.schema.json"
}
```

## Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique remediation identifier |
| `title` | string | Summary of the fix |
| `status` | [Status](#status-enum) | Current lifecycle state |
| `root_cause_ids` | string[] | Targeted root causes |
| `owner_team` | string | Responsible team |
| `created_at` | datetime | When remediation was proposed |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Implementation details |
| `assignee` | string | Individual owner |
| `deployed_at` | datetime | When fix was deployed |
| `validated_at` | datetime | When efficacy was measured |
| `efficacy` | [Efficacy](#efficacy) | Measured effectiveness |
| `external_links` | SourceSystem[] | Related PRs, tickets, etc. |
| `metadata` | object | Additional context |
| `tags` | [Tag](common.md#tag)[] | User-defined labels |

## Status Enum

Lifecycle state of the remediation.

| Value | Description |
|-------|-------------|
| `proposed` | Remediation identified |
| `approved` | Approved for implementation |
| `in_progress` | Work underway |
| `deployed` | Released to production |
| `validating` | Measuring effectiveness |
| `effective` | Confirmed working |
| `ineffective` | Did not achieve outcome |
| `rolled_back` | Reverted due to issues |
| `cancelled` | Abandoned |

## Efficacy

Measured effectiveness of the remediation.

| Field | Type | Description |
|-------|------|-------------|
| `signal_reduction` | float | Percentage decrease in signals |
| `validation_period` | TimeRange | Measurement window |
| `confidence_level` | float | Statistical confidence (0-1) |
| `notes` | string | Context on measurement |

## Example

```json
{
  "id": "rem-001",
  "title": "Implement Redis read-after-write consistency for session validation",
  "description": "Modify session validation to use WAIT command ensuring replication before read.",
  "status": "deployed",
  "root_cause_ids": ["rc-auth-001"],
  "owner_team": "identity-platform",
  "assignee": "jsmith",
  "created_at": "2024-03-10T09:00:00Z",
  "deployed_at": "2024-03-14T16:30:00Z",
  "validated_at": null,
  "efficacy": null,
  "external_links": [
    {
      "type": "code_change",
      "name": "github",
      "external_id": "PR-4521",
      "url": "https://github.com/company/oauth-service/pull/4521"
    },
    {
      "type": "incident",
      "name": "pagerduty",
      "external_id": "INC-2024-0234",
      "url": "https://company.pagerduty.com/incidents/Q1234567"
    }
  ],
  "tags": ["redis", "consistency", "auth"],
  "metadata": {
    "rollback_plan": "Revert PR-4521 and restart oauth-service pods",
    "validation_criteria": "Signal rate drops >80% within 7 days"
  }
}
```

## Go Usage

```go
import "github.com/plexusone/signal-spec/pkg/remediation"

rem := remediation.Remediation{
    ID:           "rem-001",
    Title:        "Implement Redis read-after-write consistency",
    Status:       remediation.StatusInProgress,
    RootCauseIDs: []string{"rc-auth-001"},
    OwnerTeam:    "identity-platform",
    Assignee:     "jsmith",
    CreatedAt:    time.Now(),
}
```

## Validation

Remediations are validated against these rules:

- `id` is required and must be non-empty
- `title` is required and must be non-empty
- `root_cause_ids` must contain at least one ID
- `tags` must be lowercase kebab-case

```bash
signal-spec validate -t remediation remediation.json
```

---

# ValidationSignal

A ValidationSignal represents evidence of remediation effectiveness. Generated after deployment to measure signal decay or resurgence.

## Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique validation signal ID |
| `remediation_id` | string | Linked remediation |
| `root_cause_id` | string | Linked root cause |
| `type` | ValidationResult | Outcome |
| `observed_at` | datetime | When validation occurred |
| `baseline_signal_rate` | float | Pre-remediation rate |
| `current_signal_rate` | float | Post-remediation rate |
| `reduction_percent` | float | Calculated improvement |
| `notes` | string | Context |

## ValidationResult Enum

| Value | Description |
|-------|-------------|
| `improved` | Signal rate decreased |
| `no_change` | Signal rate unchanged |
| `regressed` | Signal rate increased |

## Example

```json
{
  "id": "val-001",
  "remediation_id": "rem-001",
  "root_cause_id": "rc-auth-001",
  "type": "improved",
  "observed_at": "2024-03-21T09:00:00Z",
  "baseline_signal_rate": 15.3,
  "current_signal_rate": 2.1,
  "reduction_percent": 86.3,
  "notes": "7-day rolling average post-deployment"
}
```
