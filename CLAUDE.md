# CLAUDE.md — signal-spec

Signal Spec is the canonical data model for operational and product intelligence signals. Go structs in `pkg/` (signal, rootcause, remediation, common, export) are the source of truth; JSON schemas in `schema/` are generated from them via invopop/jsonschema, linted with schemago, and embedded via `go:embed`. The OmniSignal runtime (`plexusone/omnisignal`) consumes this IR and never forks it.

## PRISM Control

This repo's roadmap items are tracked in [prism-control](https://github.com/ProductBuildersHQ/prism-control). Use `prismctl work ready --repo github.com/plexusone/signal-spec` to find claimable work, and carry the `Refs: RMI-SIGNALSPEC-<NNN>` trailer on every commit.
