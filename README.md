# Signal Spec

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/signal-spec/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/signal-spec/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/signal-spec/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/signal-spec/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/signal-spec/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/signal-spec/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/signal-spec
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/signal-spec
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/signal-spec
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/signal-spec
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fsignal-spec
 [loc-svg]: https://tokei.rs/b1/github/plexusone/signal-spec
 [repo-url]: https://github.com/plexusone/signal-spec
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/signal-spec/blob/main/LICENSE
 
Canonical data model for operational intelligence.

## Overview

Signal Spec defines the schemas and types for:

- **Signals** - Normalized operational observations (tickets, alerts, incidents, findings)
- **Root Causes** - Persistent clustered issues with lifecycle tracking
- **Remediations** - Corrective actions with efficacy measurement
- **Validation Signals** - Evidence of fix effectiveness

## Structure

```
signal-spec/
├── cmd/signal-spec/  # CLI application
│   ├── main.go
│   └── cmd/
│       ├── root.go
│       ├── report.go
│       ├── schema.go
│       └── validate.go
├── pkg/
│   ├── common/       # Shared types (severity, domain, entity)
│   ├── signal/       # Raw signal type
│   ├── rootcause/    # Root cause type
│   ├── remediation/  # Remediation and validation types
│   └── export/       # XLSX report generation
├── schema/           # Generated JSON schemas
├── examples/         # Example payloads
└── docs/             # Architecture documentation
```

## CLI

```bash
# Build
go build -o signal-spec ./cmd/signal-spec

# Generate XLSX report from root causes
signal-spec report -i rootcauses.json -o summary.xlsx
signal-spec report -d ./rootcauses/ --leaders leaders.json -o summary.xlsx

# Generate JSON schemas
signal-spec schema generate -o schema/

# Validate a JSON file
signal-spec validate -t signal signal.json
signal-spec validate -t rootcause rootcause.json
```

## Usage

### Go Types

```go
import (
    "github.com/plexusone/signal-spec/pkg/signal"
    "github.com/plexusone/signal-spec/pkg/rootcause"
)

// Create a signal
sig := signal.Signal{
    ID:       "sig-001",
    Type:     signal.TypeSupportTicket,
    Severity: common.SeverityHigh,
    Summary:  "OAuth token refresh failures",
    // ...
}

// Create a root cause
rc := rootcause.RootCause{
    ID:     "rc-001",
    Title:  "Redis session replication instability",
    Status: rootcause.StatusActive,
    // ...
}
```

## Core Concepts

### Signal → Root Cause Mapping

Signals are raw input. Root causes are interpretations. The mapping is performed by LLM analysis with access to:

- Signal content and history
- Codebase context (via [graphize](https://github.com/plexusone/graphize))
- System documentation

### Lifecycle States

Root causes progress through states:

```
NEW → ACTIVE → MITIGATING → VALIDATING → STABLE/REGRESSED → RESOLVED
```

This enables tracking of:

- Issue persistence
- Remediation effectiveness
- Regression detection
- Operational debt accumulation

## Documentation

See [docs/architecture.md](docs/architecture.md) for detailed design.

## License

MIT
