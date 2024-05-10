package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// I did not write any of this code. I asked AI to write it and it works perfectly. Muuuah Chefs kiss xoxo

func convertYAMLToJS(yamlInput []byte) (string, error) {
	var intermediary interface{}
	if err := yaml.Unmarshal(yamlInput, &intermediary); err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(intermediary)
	if err != nil {
		return "", err
	}

	var jsonObject map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonObject); err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	buffer.WriteString("module.exports = ")
	writeJSObject(&buffer, jsonObject)
	buffer.WriteString(";")

	return buffer.String(), nil
}
func writeJSObject(buffer *bytes.Buffer, object map[string]interface{}) {
	buffer.WriteString("{")
	first := true
	for key, value := range object {
		if !first {
			buffer.WriteString(", ")
		}
		buffer.WriteString(key + ": ")
		switch v := value.(type) {
		case string:
			buffer.WriteString(`"` + v + `"`)
		case []interface{}:
			writeJSArray(buffer, v)
		case map[string]interface{}:
			writeJSObject(buffer, v)
		case nil:
			buffer.WriteString("null")
		case bool:
			if v {
				buffer.WriteString("true")
			} else {
				buffer.WriteString("false")
			}
		case float64:
			buffer.WriteString(fmt.Sprintf("%v", v))
		default:
			buffer.WriteString("null")
		}
		first = false
	}
	buffer.WriteString("}")
}

func writeJSArray(buffer *bytes.Buffer, array []interface{}) {
	buffer.WriteString("[")
	for i, value := range array {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch v := value.(type) {
		case string:
			buffer.WriteString(`"` + v + `"`)
		case map[string]interface{}:
			writeJSObject(buffer, v)
		case nil:
			buffer.WriteString("null")
		case bool:
			if v {
				buffer.WriteString("true")
			} else {
				buffer.WriteString("false")
			}
		case float64:
			buffer.WriteString(fmt.Sprintf("%v", v))
		default:
			buffer.WriteString("null")
		}
	}
	buffer.WriteString("]")
}

func writeJSFile(filePath string, jsData string) error {
	return os.WriteFile(filePath, []byte(jsData), 0644)
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: go run main.go <path_to_yaml_file> <output_js_file>")
	}

	filePath := os.Args[1]
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	jsResult, err := convertYAMLToJS(yamlData)
	if err != nil {
		log.Fatalf("Error converting YAML to JS: %v", err)
	}

	outputPath := os.Args[2]
	if err := writeJSFile(outputPath, jsResult); err != nil {
		log.Fatalf("Error writing JS file: %v", err)
	}

	fmt.Printf("Successfully converted YAML to JS and wrote to %s\n", outputPath)
}
