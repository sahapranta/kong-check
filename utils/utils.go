package utils

import (
	"encoding/json"
	"log"
	"strings"
)

func ParsePostgresArray(pgArray string) []string {
	var result []string

	if pgArray == "{}" || pgArray == "" {
		return result
	}

	trimmed := strings.Trim(pgArray, "{}")

	if strings.Contains(trimmed, "\"") {
		var parts []string
		inQuote := false
		current := ""

		for _, char := range trimmed {
			if char == '"' {
				inQuote = !inQuote
				continue
			}

			if char == ',' && !inQuote {
				parts = append(parts, strings.Trim(current, "\""))
				current = ""
				continue
			}

			current += string(char)
		}

		if current != "" {
			parts = append(parts, strings.Trim(current, "\""))
		}

		result = append(result, parts...)
	} else {
		parts := strings.Split(trimmed, ",")
		for _, p := range parts {
			result = append(result, strings.TrimSpace(p))
		}
	}

	return result
}

func ParsePostgresJSONB(jsonbStr string) map[string][]string {
	result := make(map[string][]string)

	if jsonbStr == "{}" || jsonbStr == "" {
		return result
	}

	err := json.Unmarshal([]byte(jsonbStr), &result)
	if err != nil {
		log.Printf("Error parsing JSONB: %v", err)
	}

	return result
}

func JoinStrings(strs []string, sep string) string {
	quoted := make([]string, len(strs))
	for i, s := range strs {
		quoted[i] = "\"" + strings.Replace(s, "\"", "\\\"", -1) + "\""
	}
	return strings.Join(quoted, sep)
}
