# Producing Signals

This guide covers how to emit signals from your systems into the signal-spec format.

## Overview

Signals are the entry point to the operational intelligence pipeline. They represent raw observations from various sources that need to be normalized into a canonical format.

```mermaid
flowchart LR
    subgraph Sources
        T[Ticketing]
        A[Alerting]
        S[Security]
    end

    subgraph Normalization
        N[Signal Producer]
    end

    subgraph Output
        SIG[Signal JSON]
    end

    T --> N
    A --> N
    S --> N
    N --> SIG
```

## Basic Signal Structure

At minimum, a signal requires:

```json
{
  "id": "sig-2024-001234",
  "type": "support_ticket",
  "status": "new",
  "source": {
    "type": "ticketing",
    "name": "zendesk"
  },
  "domain": {
    "name": "authentication"
  },
  "severity": "high",
  "summary": "OAuth refresh token failures",
  "observed_at": "2024-03-15T14:30:00Z",
  "received_at": "2024-03-15T14:35:22Z"
}
```

## Go Producer

```go
package main

import (
    "encoding/json"
    "time"

    "github.com/plexusone/signal-spec/pkg/signal"
    "github.com/plexusone/signal-spec/pkg/common"
)

func NewSignalFromTicket(ticket *Ticket) *signal.Signal {
    return &signal.Signal{
        ID:     generateSignalID(),
        Type:   signal.TypeSupportTicket,
        Status: signal.StatusNew,
        Source: common.SourceSystem{
            Type:       "ticketing",
            Name:       "zendesk",
            ExternalID: ticket.ID,
            URL:        ticket.URL,
        },
        Domain: common.Domain{
            Name:      mapToDomain(ticket.Category),
            Subdomain: mapToSubdomain(ticket.Subcategory),
        },
        Severity:    mapSeverity(ticket.Priority),
        Summary:     ticket.Subject,
        Description: ticket.Description,
        ObservedAt:  ticket.CreatedAt,
        ReceivedAt:  time.Now(),
        Tags:        extractTags(ticket),
    }
}
```

## Source-Specific Mapping

### Support Tickets

Map ticket fields to signal fields:

| Ticket Field | Signal Field |
|--------------|--------------|
| Ticket ID | `source.external_id` |
| Subject | `summary` |
| Description | `description` |
| Priority | `severity` |
| Category | `domain.name` |
| Created At | `observed_at` |

### Alerts

Map alert fields:

| Alert Field | Signal Field |
|-------------|--------------|
| Alert ID | `source.external_id` |
| Alert Name | `summary` |
| Description | `description` |
| Severity | `severity` |
| Service | `entities[0]` |
| Triggered At | `observed_at` |

### Security Findings

Map security finding fields:

| Finding Field | Signal Field |
|---------------|--------------|
| Finding ID | `source.external_id` |
| Title | `summary` |
| Description | `description` |
| Severity | `severity` |
| Resource | `entities[0]` |
| Detected At | `observed_at` |

## ID Generation

Signal IDs should be:

- Unique across all signals
- Deterministic (same input = same ID) for deduplication
- Human-readable prefix recommended

```go
func generateSignalID() string {
    return fmt.Sprintf("sig-%s-%s",
        time.Now().Format("2006"),
        uuid.New().String()[:8],
    )
}
```

## Deduplication

Use the `fingerprint` field to identify duplicate signals:

```go
func computeFingerprint(sig *signal.Signal) string {
    data := fmt.Sprintf("%s|%s|%s|%s",
        sig.Source.Name,
        sig.Source.ExternalID,
        sig.Summary,
        sig.ObservedAt.Format(time.RFC3339),
    )
    hash := sha256.Sum256([]byte(data))
    return fmt.Sprintf("sha256:%x", hash[:8])
}
```

## Entity Extraction

Extract referenced entities from signal content:

```go
func extractEntities(description string) []common.Entity {
    var entities []common.Entity

    // Extract service names
    servicePattern := regexp.MustCompile(`service[:\s]+([a-z0-9-]+)`)
    if matches := servicePattern.FindStringSubmatch(description); len(matches) > 1 {
        entities = append(entities, common.Entity{
            Type: "service",
            Name: matches[1],
        })
    }

    return entities
}
```

## Severity Mapping

Map source-specific severity to signal-spec severity:

| Source Severity | Signal Severity |
|-----------------|-----------------|
| P0, Critical, Sev1 | `critical` |
| P1, High, Sev2 | `high` |
| P2, Medium, Sev3 | `medium` |
| P3, Low, Sev4 | `low` |
| P4, Info | `info` |

## Batch Processing

For high-volume sources, batch signals:

```go
func ProcessBatch(tickets []*Ticket) []*signal.Signal {
    signals := make([]*signal.Signal, 0, len(tickets))

    for _, ticket := range tickets {
        sig := NewSignalFromTicket(ticket)
        sig.Fingerprint = computeFingerprint(sig)
        signals = append(signals, sig)
    }

    return signals
}
```

## Validation Before Emission

Always validate before emitting:

```go
func emitSignal(sig *signal.Signal) error {
    // Validate required fields
    if sig.ID == "" {
        return errors.New("id is required")
    }
    if sig.Summary == "" {
        return errors.New("summary is required")
    }
    if sig.Domain.Name == "" {
        return errors.New("domain.name is required")
    }

    // Validate tags
    if err := common.ValidateTags(sig.Tags); err != nil {
        return err
    }

    // Emit to queue/storage
    return publish(sig)
}
```

Or use the CLI:

```bash
signal-spec validate -t signal signal.json
```

## Best Practices

!!! tip "Preserve Original Data"
    Store source-specific data in `metadata` to preserve context:

    ```json
    {
      "metadata": {
        "zendesk_ticket_type": "incident",
        "zendesk_via": "email",
        "original_assignee": "support-team"
      }
    }
    ```

!!! tip "Use Consistent Domain Names"
    Establish a domain taxonomy and use it consistently across all signal sources.

!!! warning "Don't Over-Extract Entities"
    Only extract entities that are clearly referenced. False positives reduce signal quality.

!!! info "Set Status Appropriately"
    New signals should have `status: new`. Only set `mapped` after root cause analysis.
