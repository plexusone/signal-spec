package export

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/plexusone/signal-spec/pkg/rootcause"
)

// LoadRootCausesFromFile loads root causes from a JSON file.
// Supports both single object and array formats.
func LoadRootCausesFromFile(filename string) ([]rootcause.RootCause, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	// Try array first
	var rootCauses []rootcause.RootCause
	if err := json.Unmarshal(data, &rootCauses); err == nil {
		return rootCauses, nil
	}

	// Try single object
	var rc rootcause.RootCause
	if err := json.Unmarshal(data, &rc); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return []rootcause.RootCause{rc}, nil
}

// LoadRootCausesFromDir loads all root cause JSON files from a directory.
func LoadRootCausesFromDir(dir string) ([]rootcause.RootCause, error) {
	var all []rootcause.RootCause

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		rcs, err := LoadRootCausesFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", entry.Name(), err)
		}
		all = append(all, rcs...)
	}

	return all, nil
}
