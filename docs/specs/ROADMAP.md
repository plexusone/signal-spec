# Signal Spec — Roadmap

**Initiative:** `INIT-OMNISIGNAL-001`
**Repository:** `github.com/plexusone/signal-spec`
**Status:** Proposed — 0 of 5 items completed

> RMI IDs are stable and permanent. Commits implementing an item carry the trailer `Refs: RMI-SIGNALSPEC-NNN`. Phase status is derived from member RMIs — a phase is complete only when all its required RMIs are complete. This repository owns the canonical signal IR for initiative `INIT-OMNISIGNAL-001`; the OmniSignal runtime work lives in omnisignal's ROADMAP.md (`RMI-OMNISIGNAL-*`). Schema changes here always land before the OmniSignal code that depends on them.

## Phase 6 — Signal IR Extensions

**Theme:** Extend the IR from operational signals to product signals while keeping it deterministic.
**Status:** Proposed — 0 of 5 items completed

- [ ] `RMI-SIGNALSPEC-001` Add `enhancement_request` signal type with metadata conventions for votes, subscribers, organizations, named customers, opportunities, and estimated ARR
  - Acceptance: `signal.Type` enum (currently 8 values: `support_ticket`, `cloud_incident`, `security_finding`, `posture_drift`, `alert`, `outage`, `vulnerability`, `feedback`) gains `enhancement_request`; `TypeValues()` and `JSONSchema()` updated; metadata key conventions documented
- [ ] `RMI-SIGNALSPEC-002` Add product signal types: `competitive_gap`, `competitor_launch`, `analyst_finding`, `market_observation`
  - Depends on: `RMI-SIGNALSPEC-001`
- [ ] `RMI-SIGNALSPEC-003` Typed cross-repo entity reference conventions: `customer:`, `capability:`, `market:`, `competitor:`, `analyst-report:` ID formats for `Signal.Metadata` and `common.Entity`
  - Acceptance: reference format documented; validation helpers reject malformed typed IDs
- [ ] `RMI-SIGNALSPEC-004` Raw vs. derived metrics separation: raw source facts stay in `Metadata`; derived scores (frustration, momentum, reach) live in a distinct block excluded from fingerprinting
  - Depends on: `RMI-SIGNALSPEC-001`
  - Acceptance: IR remains deterministic — identical raw input always produces an identical fingerprint regardless of derived values
- [ ] `RMI-SIGNALSPEC-005` Regenerate and lint JSON schemas for all type and field additions
  - Depends on: `RMI-SIGNALSPEC-001`, `RMI-SIGNALSPEC-002`, `RMI-SIGNALSPEC-003`, `RMI-SIGNALSPEC-004`
  - Acceptance: `schema/*.schema.json` regenerated via invopop/jsonschema, `schemago lint` clean, embedded via `go:embed`, schema tests pass
