package llm

import (
	"strings"

	jsonrepair "github.com/RealAlexandreAI/json-repair"
)

func extractJSONContent(content string) (string, error) {
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	res, err := jsonrepair.RepairJSON(content)
	if err != nil {
		return "", err
	}

	if res == "null" {
		return "", nil
	}

	return res, nil
}
