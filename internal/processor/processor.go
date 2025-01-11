package processor

import (
	"bcem/internal/m365"
	"bcem/internal/storage"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadAllRules(folder string, db *storage.Storage) ([]DetectionRules, error) {
	var allRules []DetectionRules

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}

		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		rules, err := loadServiceRules(path)
		if err != nil {
			return fmt.Errorf("error loading rules from %s: %v", path, err)
		}

		allRules = append(allRules, rules)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allRules, nil
}

func loadServiceRules(filename string) (DetectionRules, error) {
	var rules DetectionRules

	data, err := os.ReadFile(filename)
	if err != nil {
		return rules, fmt.Errorf("failed to read file %s: %v", filename, err)
	}

	if err := yaml.Unmarshal(data, &rules); err != nil {
		return rules, fmt.Errorf("failed to parse YAML in file %s: %v", filename, err)
	}

	return rules, nil
}

func ProcessMessages(response *m365.MessageTraceResponse, allRules []DetectionRules, db *storage.Storage) []string {
	matches := []string{}

	// Process each message trace record
	for _, trace := range response.Value {
		// Apply rules to the message trace
		for _, serviceRules := range allRules {
			for _, rule := range serviceRules.Rules {
				if checkConditions(trace, rule.Conditions) {
					// Add to matches list for output
					matches = append(matches, fmt.Sprintf(
						"Matched rule: %s for service: %s with confidence %d",
						rule.Name, serviceRules.Service, rule.Confidence,
					))

					// Insert match into database
					if err := db.InsertMatch(trace.MessageID, trace.RecipientAddress, serviceRules.Service, rule.Name, rule.Confidence); err != nil {
						fmt.Printf("Failed to insert match: %v\n", err)
					}
				}
			}
		}
	}

	return matches
}
