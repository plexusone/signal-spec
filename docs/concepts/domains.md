# Domains & Taxonomy

Signal Spec uses a two-level domain taxonomy to categorize operational issues. This enables prioritization, routing, and reporting by functional area.

## Structure

```
Domain (Category)
└── Subdomain (Subcategory)
    └── Team (Owner)
```

### Domain

The top-level category representing a major functional area of your system.

**Examples:**

- `authentication` - Identity, login, sessions
- `payments` - Billing, transactions, subscriptions
- `infrastructure` - Cloud resources, networking, compute
- `security` - Vulnerabilities, compliance, access control
- `data` - Databases, pipelines, storage
- `integration` - APIs, webhooks, third-party services

### Subdomain

A more specific area within the domain.

**Examples for `authentication`:**

- `oauth` - OAuth/OIDC flows
- `sso` - Single sign-on
- `mfa` - Multi-factor authentication
- `session` - Session management

**Examples for `infrastructure`:**

- `kubernetes` - Container orchestration
- `networking` - Load balancers, DNS, CDN
- `compute` - VMs, serverless functions
- `storage` - Object storage, file systems

### Team

The team responsible for the domain/subdomain. Used for routing and reporting.

## Schema Definition

```json
{
  "domain": {
    "name": "authentication",
    "subdomain": "oauth",
    "team": "identity-platform"
  }
}
```

## Leader Mapping

For reporting purposes, domains and subdomains can be mapped to organizational leaders:

| Level | Leader Type | Description |
|-------|-------------|-------------|
| Domain | Area Leader | Executive or director responsible for the functional area |
| Subdomain | Execution Leader | Manager or tech lead responsible for implementation |

### Leader Mapping File

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
    }
  ]
}
```

Use with the CLI:

```bash
signal-spec report -i rootcauses.json --leaders leaders.json -o report.xlsx
```

## Taxonomy Guidelines

!!! tip "Keep Domains Stable"
    Domain names should be stable and change rarely. They represent architectural boundaries, not organizational structure.

!!! tip "Use Kebab-Case"
    All domain, subdomain, and tag values should be lowercase kebab-case:

    - ✅ `user-management`
    - ✅ `api-gateway`
    - ❌ `UserManagement`
    - ❌ `api_gateway`

!!! warning "Avoid Deep Hierarchies"
    Two levels (domain/subdomain) is usually sufficient. Deeper hierarchies become hard to maintain and report on.

## Example Taxonomy

```yaml
authentication:
  - oauth
  - sso
  - mfa
  - session
  - password

payments:
  - checkout
  - subscriptions
  - refunds
  - fraud

infrastructure:
  - kubernetes
  - networking
  - compute
  - storage
  - monitoring

security:
  - vulnerabilities
  - compliance
  - access-control
  - secrets

data:
  - databases
  - pipelines
  - analytics
  - search

integration:
  - webhooks
  - apis
  - third-party
```

## Prioritization by Domain

The XLSX report aggregates issues by domain/subdomain, enabling prioritization:

| Domain | Subdomain | Issue Count | Total Signals | Max Severity | Area Leader |
|--------|-----------|-------------|---------------|--------------|-------------|
| authentication | oauth | 3 | 487 | high | Jane Smith |
| infrastructure | kubernetes | 5 | 234 | critical | Mike Lee |
| payments | checkout | 2 | 156 | medium | Sarah Wong |

This view helps leadership understand:

- Which areas have the most issues
- Where customer impact is highest
- Who owns each problem space
