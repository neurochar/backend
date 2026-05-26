//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func escapeNewlinesInJSON(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	inString := false
	escape := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escape {
			escape = false
			buf.WriteByte(ch)
			continue
		}
		if ch == '\\' && inString {
			escape = true
			buf.WriteByte(ch)
			continue
		}
		if ch == '"' {
			inString = !inString
			buf.WriteByte(ch)
			continue
		}
		if inString {
			switch ch {
			case '\n':
				buf.WriteString("\\n")
			case '\r':
				if i+1 < len(s) && s[i+1] == '\n' {
					buf.WriteString("\\n")
					i++
				} else {
					buf.WriteString("\\r")
				}
			default:
				buf.WriteByte(ch)
			}
			continue
		}
		buf.WriteByte(ch)
	}
	return buf.String()
}

func main() {
	// Test 1: literal newlines inside JSON strings
	input := "{\n\"key\": \"line1\nline2\",\n\"nested\": {\n\"inner\": \"text\nwith\nbreaks\"\n}\n}"
	fmt.Println("=== Test 1: Literal newlines in strings ===")
	cleaned := escapeNewlinesInJSON(input)
	var result map[string]any
	err := json.Unmarshal([]byte(cleaned), &result)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println("OK - parsed successfully")
		if s, ok := result["key"].(string); ok {
			fmt.Printf("key = %q\n", s)
			fmt.Printf("contains newline: %v (expected true)\n", strings.Contains(s, "\n"))
		}
		if nested, ok := result["nested"].(map[string]any); ok {
			if s, ok := nested["inner"].(string); ok {
				fmt.Printf("nested.inner = %q\n", s)
				fmt.Printf("contains newline: %v (expected true)\n", strings.Contains(s, "\n"))
			}
		}
	}

	// Test 2: already valid JSON with \n escapes
	fmt.Println("\n=== Test 2: Already valid JSON ===")
	input2 := "{\"a\": \"hello\\nworld\", \"b\": \"test\"}"
	cleaned2 := escapeNewlinesInJSON(input2)
	var result2 map[string]any
	err2 := json.Unmarshal([]byte(cleaned2), &result2)
	if err2 != nil {
		fmt.Printf("ERROR: %v\n", err2)
	} else {
		fmt.Println("OK - parsed successfully")
		if s, ok := result2["a"].(string); ok {
			fmt.Printf("a = %q\n", s)
			fmt.Printf("contains newline: %v (expected true)\n", strings.Contains(s, "\n"))
		}
	}

	// Test 3: escaped quotes preserved
	fmt.Println("\n=== Test 3: Escaped quotes ===")
	input3 := "{\"x\": \"he said \\\"hello\\\"\"}"
	cleaned3 := escapeNewlinesInJSON(input3)
	var result3 map[string]any
	err3 := json.Unmarshal([]byte(cleaned3), &result3)
	if err3 != nil {
		fmt.Printf("ERROR: %v\n", err3)
	} else {
		fmt.Println("OK - parsed successfully")
		fmt.Printf("x = %q\n", result3["x"])
	}

	// Test 4: CRLF in string
	fmt.Println("\n=== Test 4: CRLF in string ===")
	input4 := "{\"y\": \"line1\r\nline2\"}"
	cleaned4 := escapeNewlinesInJSON(input4)
	var result4 map[string]any
	err4 := json.Unmarshal([]byte(cleaned4), &result4)
	if err4 != nil {
		fmt.Printf("ERROR: %v\n", err4)
	} else {
		fmt.Println("OK - parsed successfully")
		if s, ok := result4["y"].(string); ok {
			fmt.Printf("y = %q\n", s)
			fmt.Printf("contains \\r: %v (expected false)\n", strings.Contains(s, "\r"))
		}
	}

	// Test 5: real-world-like JSON with multiline strings
	fmt.Println("\n=== Test 5: Real-world multiline ===")
	input5 := `{
"hiring_decision": "do_not_hire",
"confidence_score": 0.85,
"main_recommendation": "Кандидат имеет опыт.
Второе предложение.",
"personality_fit": {
"score": 29,
"summary": "Кандидат демонстрирует
сдержанность."
}
}`
	cleaned5 := escapeNewlinesInJSON(input5)
	var result5 map[string]any
	err5 := json.Unmarshal([]byte(cleaned5), &result5)
	if err5 != nil {
		fmt.Printf("ERROR: %v\n", err5)
	} else {
		fmt.Println("OK - parsed successfully")
		fmt.Printf("main_recommendation = %q\n", result5["main_recommendation"])
		if pf, ok := result5["personality_fit"].(map[string]any); ok {
			fmt.Printf("summary = %q\n", pf["summary"])
		}
	}
}
