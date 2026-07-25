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

## Product Signal Examples

Product signal types (`enhancement_request`, `competitive_gap`, `competitor_launch`, `analyst_finding`, `market_observation`) flow through the same pipeline as operational signals but carry structured product data via well-known `metadata` keys and typed cross-repo references. See [`signal.MetaVotes`, `signal.MetaSubscribers`, and related constants](schemas/signal.md#enhancement-signal-metadata) and the [`pkg/ref`](schemas/ref.md) reference format.

### Enhancement Request Signal

A feature request aggregated from a feedback board, carrying vote/subscriber counts and cross-repo references to the affected customers and market:

```json
{
  "id": "sig-2026-010001",
  "type": "enhancement_request",
  "status": "new",
  "source": {
    "type": "feedback",
    "name": "productboard",
    "external_id": "PB-4821",
    "url": "https://company.productboard.com/feature/4821"
  },
  "domain": {
    "name": "identity",
    "subdomain": "sso",
    "team": "identity-platform"
  },
  "severity": "medium",
  "summary": "Support SCIM group provisioning for Okta integration",
  "description": "Enterprise customers using Okta as an identity provider need automatic group membership sync via SCIM, not just user provisioning.",
  "entities": [
    {
      "type": "capability",
      "name": "scim-provisioning",
      "ref": "capability:scim-provisioning"
    },
    {
      "type": "customer",
      "name": "Acme Corp",
      "ref": "customer:acme-001"
    }
  ],
  "observed_at": "2026-06-01T00:00:00Z",
  "received_at": "2026-06-01T00:05:00Z",
  "tags": ["enterprise", "sso", "scim"],
  "metadata": {
    "votes": 142,
    "subscribers": 58,
    "organizations": ["Acme Corp", "Globex", "Initech"],
    "customers": ["acme-001", "globex-014"],
    "opportunities": ["OPP-88213"],
    "estimated_arr": 24000000,
    "market_ref": "market:identity-governance",
    "capability_ref": "capability:scim-provisioning"
  },
  "derived": {
    "frustration": 6.1,
    "momentum": 3.4,
    "reach": 12,
    "urgency": 2.8,
    "computed_at": "2026-07-20T00:00:00Z"
  }
}
```

Field notes:

- `metadata.votes`, `metadata.subscribers`, `metadata.organizations`, `metadata.customers`, `metadata.opportunities`, and `metadata.estimated_arr` correspond to the `signal.MetaVotes`, `signal.MetaSubscribers`, `signal.MetaOrganizations`, `signal.MetaCustomers`, `signal.MetaOpportunities`, and `signal.MetaEstimatedARR` Go constants. `estimated_arr` is in cents (`$240,000.00` above).
- `metadata.market_ref` and `metadata.capability_ref` are typed references (`signal.MetaMarketRef`, `signal.MetaCapabilityRef`) in `{type}:{slug}` format, validated by [`pkg/ref`](schemas/ref.md).
- `entities[].ref` links an entity directly to its canonical definition in another repo (here, MarketSpec's `capability:scim-provisioning` and OrganizationSpec's `customer:acme-001`).
- `derived` holds recomputed scores and is excluded from `ComputeFingerprint()`.

### Competitive Gap Signal

A gap identified from win/loss analysis, referencing the losing competitor and affected market:

```json
{
  "id": "sig-2026-010002",
  "type": "competitive_gap",
  "status": "new",
  "source": {
    "type": "sales",
    "name": "gong",
    "external_id": "CALL-99231",
    "url": "https://company.gong.io/calls/99231"
  },
  "domain": {
    "name": "identity",
    "subdomain": "governance",
    "team": "product-identity"
  },
  "severity": "high",
  "summary": "Lost enterprise deal to Okta over lack of access certification workflows",
  "description": "Win/loss interview with Initech (closed-lost) cites missing periodic access certification and attestation reporting as the deciding factor against Okta Identity Governance.",
  "entities": [
    {
      "type": "competitor",
      "name": "Okta",
      "ref": "competitor:okta"
    },
    {
      "type": "customer",
      "name": "Initech",
      "ref": "customer:initech-002"
    }
  ],
  "observed_at": "2026-06-10T00:00:00Z",
  "received_at": "2026-06-10T00:10:00Z",
  "tags": ["win-loss", "enterprise", "governance"],
  "metadata": {
    "organizations": ["Initech"],
    "opportunities": ["OPP-91004"],
    "estimated_arr": 18000000,
    "competitor_ref": "competitor:okta",
    "market_ref": "market:identity-governance"
  }
}
```

### Competitor Launch Signal

A tracked competitor product announcement:

```json
{
  "id": "sig-2026-010003",
  "type": "competitor_launch",
  "status": "new",
  "source": {
    "type": "market_intel",
    "name": "competitor-rss",
    "external_id": "OKTA-2026-06-15",
    "url": "https://www.okta.com/blog/2026/06/identity-governance-launch"
  },
  "domain": {
    "name": "identity",
    "subdomain": "governance",
    "team": "product-identity"
  },
  "severity": "medium",
  "summary": "Okta launches Identity Governance access certification module",
  "description": "Okta announced general availability of periodic access certification and attestation reporting as part of Okta Identity Governance, directly addressing a gap cited in recent competitive losses.",
  "entities": [
    {
      "type": "competitor",
      "name": "Okta",
      "ref": "competitor:okta"
    }
  ],
  "observed_at": "2026-06-15T00:00:00Z",
  "received_at": "2026-06-15T06:00:00Z",
  "tags": ["competitor-launch", "governance"],
  "metadata": {
    "competitor_ref": "competitor:okta",
    "market_ref": "market:identity-governance",
    "announcement_type": "general_availability"
  }
}
```

### Analyst Finding Signal

An insight extracted from an analyst report:

```json
{
  "id": "sig-2026-010004",
  "type": "analyst_finding",
  "status": "new",
  "source": {
    "type": "analyst",
    "name": "gartner",
    "external_id": "GARTNER-MQ-IAM-2026",
    "url": "https://www.gartner.com/en/documents/gartner-mq-iam-2026"
  },
  "domain": {
    "name": "identity",
    "subdomain": "governance",
    "team": "product-identity"
  },
  "severity": "medium",
  "summary": "Gartner MQ for IAM cites access certification as a required capability for Leaders quadrant",
  "description": "The 2026 Gartner Magic Quadrant for Identity and Access Management identifies periodic access certification workflows as a differentiating capability separating Leaders from Visionaries.",
  "entities": [
    {
      "type": "analyst-report",
      "name": "Gartner MQ IAM 2026",
      "ref": "analyst-report:gartner-mq-iam-2026"
    }
  ],
  "observed_at": "2026-06-20T00:00:00Z",
  "received_at": "2026-06-20T09:00:00Z",
  "tags": ["gartner", "governance", "quadrant"],
  "metadata": {
    "analyst_report_ref": "analyst-report:gartner-mq-iam-2026",
    "market_ref": "market:identity-governance",
    "quadrant_position": "visionaries"
  }
}
```

### Market Observation Signal

A general market trend observation:

```json
{
  "id": "sig-2026-010005",
  "type": "market_observation",
  "status": "new",
  "source": {
    "type": "market_intel",
    "name": "internal-research",
    "external_id": "MKT-2026-Q2-014"
  },
  "domain": {
    "name": "identity",
    "subdomain": "governance",
    "team": "product-identity"
  },
  "severity": "low",
  "summary": "Growing enterprise demand for continuous access certification over periodic reviews",
  "description": "Quarterly market scan shows increasing RFP language requiring continuous (event-driven) access certification rather than quarterly/annual review cycles, particularly in regulated industries.",
  "entities": [
    {
      "type": "market",
      "name": "Identity Governance",
      "ref": "market:identity-governance"
    }
  ],
  "observed_at": "2026-06-25T00:00:00Z",
  "received_at": "2026-06-25T00:00:00Z",
  "tags": ["market-trend", "governance", "regulated"],
  "metadata": {
    "market_ref": "market:identity-governance"
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
- [`examples/signal_enhancement_request.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_enhancement_request.json)
- [`examples/signal_competitive_gap.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_competitive_gap.json)
- [`examples/signal_competitor_launch.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_competitor_launch.json)
- [`examples/signal_analyst_finding.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_analyst_finding.json)
- [`examples/signal_market_observation.json`](https://github.com/plexusone/signal-spec/blob/main/examples/signal_market_observation.json)
- [`examples/rootcause_auth_failure.json`](https://github.com/plexusone/signal-spec/blob/main/examples/rootcause_auth_failure.json)
- [`examples/rootcauses_sample.json`](https://github.com/plexusone/signal-spec/blob/main/examples/rootcauses_sample.json)
- [`examples/remediation_redis_fix.json`](https://github.com/plexusone/signal-spec/blob/main/examples/remediation_redis_fix.json)
- [`examples/leaders.json`](https://github.com/plexusone/signal-spec/blob/main/examples/leaders.json)
