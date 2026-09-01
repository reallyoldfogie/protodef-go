package protocol

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateJSONSchema_Valid(t *testing.T) {
	schemaPath := "../ProtoDef/schemas/protocol_schema.json"
	// Valid protocol JSON with types field
	jsonData := []byte(`{
		"types": {
			"myType": "u8",
			"myStruct": ["container", [{"name": "field1", "type": "u8"}]]
		}
	}`)

	// jsonData, err := os.ReadFile("../ProtoDef/test/structures.json")
	// if err != nil {
	// 	t.Fatalf("failed to read test JSON: %v", err)
	// }
	valid, err := ValidateJSONSchema(schemaPath, jsonData)
	if err != nil {
		t.Fatalf("validation error: %v", err)
	}
	if !valid {
		t.Error("expected valid JSON to pass schema validation")
	}
}

func TestValidateJSONSchema_Invalid(t *testing.T) {
	schemaPath := "../ProtoDef/schemas/protocol_schema.json"
	// Intentionally broken JSON (missing required fields)
	jsonData := []byte(`{"foo": 123}`)
	valid, err := ValidateJSONSchema(schemaPath, jsonData)
	if err != nil {
		t.Fatalf("validation error: %v", err)
	}
	if valid {
		t.Error("expected invalid JSON to fail schema validation")
	}
}

func TestValidateJSONSchemaWithSimpleSchema(t *testing.T) {
	tmpDir := t.TempDir()
	schemaFile := filepath.Join(tmpDir, "schema.json")

	schema := `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "number"}
		},
		"required": ["name"]
	}`

	data := `{"name": "test", "age": 25}`

	err := os.WriteFile(schemaFile, []byte(schema), 0644)
	if err != nil {
		t.Fatalf("Failed to create schema file: %v", err)
	}

	valid, err := ValidateJSONSchema(schemaFile, []byte(data))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v", err)
	}

	if !valid {
		t.Error("valid JSON should pass validation")
	}
}

func TestValidateJSONSchemaWithMissingRequired(t *testing.T) {
	tmpDir := t.TempDir()
	schemaFile := filepath.Join(tmpDir, "schema.json")

	schema := `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "number"}
		},
		"required": ["name"]
	}`

	data := `{"age": 25}`

	err := os.WriteFile(schemaFile, []byte(schema), 0644)
	if err != nil {
		t.Fatalf("Failed to create schema file: %v", err)
	}

	valid, err := ValidateJSONSchema(schemaFile, []byte(data))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v", err)
	}

	if valid {
		t.Error("JSON missing required field should fail validation")
	}
}

func TestValidateJSONSchemaSchemaNotFound(t *testing.T) {
	data := `{"test": "value"}`

	valid, err := ValidateJSONSchema("/nonexistent/schema.json", []byte(data))

	if err == nil {
		t.Error("ValidateJSONSchema() should return error for nonexistent schema")
	}

	if valid {
		t.Error("validation should not be valid when schema is missing")
	}
}

func TestValidateJSONSchemaTypeMismatch(t *testing.T) {
	tmpDir := t.TempDir()
	schemaFile := filepath.Join(tmpDir, "schema.json")

	schema := `{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"type": "object",
		"properties": {
			"count": {"type": "integer"}
		}
	}`

	data := `{"count": "not a number"}`

	err := os.WriteFile(schemaFile, []byte(schema), 0644)
	if err != nil {
		t.Fatalf("Failed to create schema file: %v", err)
	}

	valid, err := ValidateJSONSchema(schemaFile, []byte(data))
	if err != nil {
		t.Fatalf("ValidateJSONSchema() error = %v", err)
	}

	if valid {
		t.Error("type mismatch should fail validation")
	}
}
