package schema

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedSchemasAreValidJSON(t *testing.T) {
	schemas := []struct {
		name string
		data []byte
	}{
		{"signal", SignalSchema},
		{"rootcause", RootCauseSchema},
		{"remediation", RemediationSchema},
		{"validation_signal", ValidationSignalSchema},
	}

	for _, s := range schemas {
		if len(s.data) == 0 {
			t.Errorf("%s: embedded schema is empty", s.name)
			continue
		}
		if !json.Valid(s.data) {
			t.Errorf("%s: embedded schema is not valid JSON", s.name)
		}
	}
}

func TestEmbeddedSignalSchemaHasEnhancementRequest(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(SignalSchema, &schema); err != nil {
		t.Fatalf("unmarshal signal schema: %v", err)
	}

	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("missing $defs in signal schema")
	}

	typeDef, ok := defs["Type"].(map[string]any)
	if !ok {
		t.Fatal("missing Type definition in signal schema")
	}

	enumValues, ok := typeDef["enum"].([]any)
	if !ok {
		t.Fatal("missing enum in Type definition")
	}

	found := false
	for _, v := range enumValues {
		if v == "enhancement_request" {
			found = true
			break
		}
	}
	if !found {
		t.Error("enhancement_request not found in Type enum")
	}
}

func TestEmbeddedSignalSchemaHasDerivedMetrics(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(SignalSchema, &schema); err != nil {
		t.Fatalf("unmarshal signal schema: %v", err)
	}

	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("missing $defs in signal schema")
	}

	if _, ok := defs["DerivedMetrics"]; !ok {
		t.Error("DerivedMetrics definition missing from signal schema")
	}
}

func TestAllFS(t *testing.T) {
	entries, err := All.ReadDir(".")
	if err != nil {
		t.Fatalf("reading All FS: %v", err)
	}
	if len(entries) != 4 {
		t.Errorf("expected 4 schema files in All FS, got %d", len(entries))
	}
}
