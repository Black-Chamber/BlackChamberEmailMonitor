package main

import (
	"bcem/internal/config"
	"bcem/internal/m365"
	"bcem/internal/processor"
	"bcem/internal/storage"
	"flag"
	"log"
	"time"
)

func main() {
	configPath := flag.String("cfg", "../../config/tenant1.yml", "Path to the configuration file")
	flag.Parse()
	AppConfig, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := storage.Initialize(AppConfig.Database.Path)

	instance, err := m365.NewM365Instance(*AppConfig)
	if err != nil {
		log.Fatalf("Failed to initialize M365 instance: %v", err)
	}

	// Load rules dynamically from the "services" folder
	allRules, err := processor.LoadAllRules(AppConfig.ServiceDefinitions.Path, db)
	if err != nil {
		log.Fatalf("Failed to load detection rules: %v", err)
	}

	// Define the initial trace start time as interval duration ago from now
	var initialTraceStartTime = (time.Now().UTC().Add(-AppConfig.Scan.Interval))

	// Instant results for the first run
	scheduledStartTime, err := processInstance(instance, allRules, db, initialTraceStartTime)
	if err != nil {
		log.Fatalf("Failed to inital process instance: %v", err)
	}

	// Scheduled processing for continuous monitoring
	scheduledProcessing(instance, allRules, db, AppConfig.Scan.Interval, scheduledStartTime)
}

func processInstance(instance *m365.M365Instance, allRules []processor.DetectionRules, db *storage.Storage, startTime time.Time) (time.Time, error) {
	// Define the time range for lookup
	endTime := time.Now().UTC()
	startTime = startTime.UTC()

	// Perform the lookup
	result, err := instance.PerformLookup(startTime, endTime)
	if err != nil {
		log.Printf("Error performing lookup for tenant %s: %v", instance.TenantID, err)
		return startTime, err
	}

	// Process the messages against detection rules
	matches := processor.ProcessMessages(result, allRules, db)

	// Log the matches
	log.Printf("Processing complete for tenant %s. Matches found: %d", instance.TenantID, len(matches))
	for _, match := range matches {
		log.Println(match)
	}

	//End time for on a successful run, will be used as start time for the next run
	return endTime, nil
}

func scheduledProcessing(instance *m365.M365Instance, allRules []processor.DetectionRules, db *storage.Storage, interval time.Duration, startTime time.Time) {
	// Create a ticker that triggers at the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			endTime, err := processInstance(instance, allRules, db, startTime)
			if err != nil {
				log.Printf("Error processing instance %s: %v", instance.TenantID, err)
				continue
			}
			startTime = endTime

		}
	}

}
