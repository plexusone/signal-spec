# CLI Reference

The `signal-spec` CLI provides commands for validation, schema generation, and report generation.

## Installation

```bash
go install github.com/plexusone/signal-spec/cmd/signal-spec@latest
```

Or build from source:

```bash
git clone https://github.com/plexusone/signal-spec.git
cd signal-spec
go build -o signal-spec ./cmd/signal-spec
```

## Commands

### report

Generate XLSX summary reports from root cause data.

```bash
signal-spec report [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input JSON file containing root causes |
| `--dir` | `-d` | Directory containing root cause JSON files |
| `--output` | `-o` | Output XLSX filename (required) |
| `--leaders` | | JSON file with leader mappings |

**Examples:**

```bash
# From single file
signal-spec report -i rootcauses.json -o summary.xlsx

# From directory
signal-spec report -d ./rootcauses/ -o summary.xlsx

# With leader mappings
signal-spec report -i rootcauses.json --leaders leaders.json -o summary.xlsx
```

**Output:**

The generated XLSX contains two sheets:

1. **Domain Summary** - Aggregated metrics by domain/subdomain
    - Domain (Category)
    - Subdomain (Subcategory)
    - Issue Count
    - Total Signals
    - Avg Priority
    - Max Severity
    - Area Leader
    - Execution Leader
    - Top Issues

2. **Root Causes** - Detailed root cause listing
    - ID
    - Title
    - Domain
    - Subdomain
    - Status
    - Severity
    - Signal Count
    - Priority Score
    - First Seen
    - Last Seen
    - Owner Team
    - Tags

---

### schema generate

Generate JSON schemas from Go types.

```bash
signal-spec schema generate [flags]
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output directory for schema files (default: `schema`) |

**Examples:**

```bash
# Generate to default directory
signal-spec schema generate

# Generate to custom directory
signal-spec schema generate -o ./schemas/
```

**Output:**

Generates the following schema files:

- `signal.schema.json`
- `rootcause.schema.json`
- `remediation.schema.json`
- `validation_signal.schema.json`

---

### validate

Validate JSON files against signal-spec types.

```bash
signal-spec validate -t <type> <file>
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--type` | `-t` | Type to validate: `signal`, `rootcause`, `remediation` (required) |

**Examples:**

```bash
# Validate a signal
signal-spec validate -t signal signal.json

# Validate a root cause
signal-spec validate -t rootcause rootcause.json

# Validate a remediation
signal-spec validate -t remediation remediation.json
```

**Output:**

Valid:
```
Valid signal: signal.json
```

Invalid:
```
Validation failed for signal.json:
  - id is required
  - domain.name is required
  - tag "Invalid_Tag" is not valid lower-kebab-case
```

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--help` | Help for any command |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (validation failed, file not found, etc.) |

## Leader Mappings File

The `--leaders` flag accepts a JSON file mapping domains to organizational leaders:

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
    }
  ]
}
```

When applied, the Area Leader and Execution Leader columns in the XLSX report are populated based on matching domain/subdomain pairs.
