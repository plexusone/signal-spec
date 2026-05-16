package export

import (
	"encoding/json"
	"os"
)

// LeaderMapping maps domains/subdomains to leaders.
type LeaderMapping struct {
	// AreaLeaders maps domain names to area leaders.
	AreaLeaders map[string]string `json:"area_leaders"`

	// ExecutionLeaders maps "domain|subdomain" to execution leaders.
	ExecutionLeaders map[string]string `json:"execution_leaders"`
}

// LoadLeaderMapping loads leader assignments from a JSON file.
func LoadLeaderMapping(filename string) (*LeaderMapping, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var mapping LeaderMapping
	if err := json.Unmarshal(data, &mapping); err != nil {
		return nil, err
	}

	return &mapping, nil
}

// ApplyLeaders applies leader assignments to summaries.
func (m *LeaderMapping) ApplyLeaders(summaries []DomainSummary) {
	for i := range summaries {
		s := &summaries[i]

		// Apply area leader from domain mapping
		if leader, ok := m.AreaLeaders[s.Domain]; ok {
			s.AreaLeader = leader
		}

		// Apply execution leader from domain|subdomain mapping
		key := s.Domain + "|" + s.Subdomain
		if leader, ok := m.ExecutionLeaders[key]; ok {
			s.ExecutionLeader = leader
		}
	}
}
