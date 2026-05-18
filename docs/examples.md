# Examples

Complete examples demonstrating signal-spec usage patterns.

## Signal Examples

### Support Ticket Signal

A customer support ticket normalized as a signal:

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
  "description": "Customer reports being logged out repeatedly. Error shows 'invalid_grant' when refreshing OAuth token. Started after recent mobile app update.",
  "entities": [
    {
      "type": "service",
      "name": "oauth-service",
      "attributes": {
        "environment": "production"
      }
    },
    {
      "type": "application",
      "name": "mobile-app-ios",
      "attributes": {
        "version": "4.2.1"
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
    "affected_users": 150,
    "escalated": true
  }
}
```

### Security Finding Signal

A security scan finding as a signal:

```json
{
  "id": "sig-2024-002345",
  "type": "security_finding",
  "status": "new",
  "source": {
    "type": "security",
    "name": "wiz",
    "external_id": "WIZ-CVE-2024-1234",
    "url": "https://app.wiz.io/findings/CVE-2024-1234"
  },
  "domain": {
    "name": "security",
    "subdomain": "vulnerabilities",
    "team": "security-ops"
  },
  "severity": "critical",
  "summary": "Critical RCE vulnerability in log4j dependency",
  "description": "CVE-2024-1234 detected in payment-service. Remote code execution possible via crafted log messages.",
  "entities": [
    {
      "type": "service",
      "name": "payment-service",
      "attributes": {
        "environment": "production",
        "region": "us-east-1"
      }
    }
  ],
  "observed_at": "2024-03-16T09:00:00Z",
  "received_at": "2024-03-16T09:05:00Z",
  "tags": ["cve", "critical", "rce", "java"],
  "metadata": {
    "cve_id": "CVE-2024-1234",
    "cvss_score": 9.8,
    "affected_package": "log4j:2.14.0"
  }
}
```

---

## Root Cause Examples

### Authentication Root Cause

A root cause aggregating authentication-related signals:

```json
{
  "id": "rc-auth-001",
  "title": "Redis session replication instability causing OAuth token validation failures",
  "description": "Intermittent Redis cluster replication lag causes session state inconsistency. When tokens are refreshed, the validation service sometimes reads stale data, resulting in invalid_grant errors.",
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
    "invalid_grant errors",
    "Session expired unexpectedly",
    "Token validation timeout"
  ],
  "signal_ids": [
    "sig-2024-001234",
    "sig-2024-001235",
    "sig-2024-001240",
    "sig-2024-001241",
    "sig-2024-001245"
  ],
  "impact": {
    "signal_count": 487,
    "affected_customers": 2341,
    "affected_entities": [
      {
        "type": "service",
        "name": "oauth-service"
      },
      {
        "type": "service",
        "name": "session-store"
      }
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
  "tags": ["redis", "auth", "session", "enterprise-impact"],
  "metadata": {
    "related_incidents": ["INC-2024-0234"],
    "affected_regions": ["us-east-1", "us-west-2"]
  }
}
```

---

## Remediation Examples

### Redis Fix Remediation

A remediation targeting the Redis replication issue:

```json
{
  "id": "rem-001",
  "title": "Implement Redis read-after-write consistency for session validation",
  "description": "Modify session validation to use WAIT command ensuring replication before read. Add circuit breaker for Redis cluster failover scenarios. Update health checks to detect replication lag.",
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

### Validated Remediation

A remediation with efficacy measurements:

```json
{
  "id": "rem-002",
  "title": "Add rate limiting to user search API",
  "description": "Implement token bucket rate limiting to prevent abuse of user search endpoint.",
  "status": "effective",
  "root_cause_ids": ["rc-perf-002"],
  "owner_team": "api-platform",
  "assignee": "alee",
  "created_at": "2024-03-01T10:00:00Z",
  "deployed_at": "2024-03-05T14:00:00Z",
  "validated_at": "2024-03-12T09:00:00Z",
  "efficacy": {
    "signal_reduction": 0.94,
    "validation_period": {
      "start": "2024-03-05T14:00:00Z",
      "end": "2024-03-12T09:00:00Z"
    },
    "confidence_level": 0.95,
    "notes": "7-day rolling average shows 94% reduction in timeout errors"
  },
  "external_links": [
    {
      "type": "code_change",
      "name": "github",
      "external_id": "PR-4498",
      "url": "https://github.com/company/api-gateway/pull/4498"
    }
  ],
  "tags": ["rate-limiting", "api", "performance"]
}
```

---

## Leader Mappings

Map domains to organizational leaders for reporting:

```json
{
  "mappings": [
    {
      "domain": "authentication",
      "subdomain": "oauth",
      "area_leader": "Jane Smith",
      "execution_leader": "Bob Johnson"
    },
    {
      "domain": "authentication",
      "subdomain": "sso",
      "area_leader": "Jane Smith",
      "execution_leader": "Alice Chen"
    },
    {
      "domain": "infrastructure",
      "subdomain": "kubernetes",
      "area_leader": "Mike Lee",
      "execution_leader": "David Park"
    },
    {
      "domain": "payments",
      "subdomain": "checkout",
      "area_leader": "Sarah Wong",
      "execution_leader": "Chris Taylor"
    }
  ]
}
```

---

## Workflow Example

### End-to-End Flow

1. **Ingest signal from ticketing system:**

```bash
curl -X POST https://signal-api/signals \
  -H "Content-Type: application/json" \
  -d @signal_support_ticket.json
```

2. **Validate the signal:**

```bash
signal-spec validate -t signal signal_support_ticket.json
# Valid signal: signal_support_ticket.json
```

3. **After LLM analysis, create root cause:**

```bash
signal-spec validate -t rootcause rootcause_auth_failure.json
# Valid rootcause: rootcause_auth_failure.json
```

4. **Generate summary report:**

```bash
signal-spec report \
  -d ./rootcauses/ \
  --leaders leaders.json \
  -o summary.xlsx
# Loaded 15 root causes from ./rootcauses/
# Applied leader mappings
# Generated summary.xlsx
```

5. **Review XLSX report** with Domain Summary and Root Causes sheets

---

## Sample Files

Example files are available in the repository:

- [`examples/signal_support_ticket.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_support_ticket.json)
- [`examples/signal_security_finding.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_security_finding.json)
- [`examples/rootcause_auth_failure.json`](https://github.com/plexusone/signal-spec/blob/main/examples/rootcause_auth_failure.json)
- [`examples/rootcauses_sample.json`](https://github.com/plexusone/signal-spec/blob/main/examples/rootcauses_sample.json)
- [`examples/remediation_redis_fix.json`](https://github.com/plexusone/signal-spec/blob/main/examples/remediation_redis_fix.json)
- [`examples/leaders.json`](https://github.com/plexusone/signal-spec/blob/main/examples/leaders.json)
