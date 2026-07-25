// Package schema provides embedded JSON Schema files for signal-spec types.
package schema

import "embed"

//go:embed signal.schema.json
var SignalSchema []byte

//go:embed rootcause.schema.json
var RootCauseSchema []byte

//go:embed remediation.schema.json
var RemediationSchema []byte

//go:embed validation_signal.schema.json
var ValidationSignalSchema []byte

//go:embed *.schema.json
var All embed.FS
