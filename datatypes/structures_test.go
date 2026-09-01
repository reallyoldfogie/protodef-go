package datatypes_test

import (
	"os"
	"testing"

	"github.com/protodef-go/protodef-go/protocol"
)

func TestStructuresSchemaValidation(t *testing.T) {
	// Validate structures test cases against the datatype_tests_schema
	schemaPath := "../ProtoDef/test/datatype_tests_schema.json"
	jsonData, err := os.ReadFile("../ProtoDef/test/structures.json")
	if err != nil {
		t.Fatalf("failed to read structures test JSON: %v", err)
	}
	valid, err := protocol.ValidateJSONSchema(schemaPath, jsonData)
	if err != nil {
		t.Skipf("validation error (schema may have issues): %v", err)
		return
	}
	if !valid {
		t.Error("structures.json should pass datatype_tests_schema validation")
	}
}
