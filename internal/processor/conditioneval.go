package processor

import (
	"bcem/internal/m365"
	"fmt"
	"regexp"
)

func checkConditions(message m365.MessageTrace, conditions []Condition) bool {
	for _, condition := range conditions {
		// Retrieve the value of the specified field from the message
		fieldValue := getFieldValue(message, condition.Field)
		if fieldValue == "" {
			// If the field value is empty and a regex is defined, it cannot match
			return false
		}

		// If a regex is defined, ensure the field value matches it
		if condition.Regex != "" {
			matched, err := regexp.MatchString(condition.Regex, fieldValue)
			if err != nil {
				fmt.Printf("Error matching regex '%s': %v\n", condition.Regex, err)
				return false // Treat regex errors as non-matching
			}
			if !matched {
				return false // If the regex does not match, the condition fails
			}
		}
	}

	// All conditions passed
	return true
}

func getFieldValue(message m365.MessageTrace, field string) string {
	switch field {
	case "SenderAddress":
		return message.SenderAddress
	case "RecipientAddress":
		return message.RecipientAddress
	case "Subject":
		return message.Subject
	case "Status":
		return message.Status
	case "ToIP":
		if message.ToIP == nil {
			return "undefined"
		}
		return *message.ToIP
	case "FromIP":
		if message.FromIP == nil {
			return "undefined"
		}
		return *message.FromIP
	default:
		return ""
	}
}
