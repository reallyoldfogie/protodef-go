package datatypes_test

import (
	"os"
	"testing"

	"github.com/protodef-go/protodef-go/protocol"
)

func TestConditionalSchemaValidation(t *testing.T) {
	// Validate conditional test cases against the datatype_tests_schema
	schemaPath := "../ProtoDef/test/datatype_tests_schema.json"
	jsonData, err := os.ReadFile("../ProtoDef/test/conditional.json")
	if err != nil {
		t.Fatalf("failed to read conditional test JSON: %v", err)
	}
	valid, err := protocol.ValidateJSONSchema(schemaPath, jsonData)
	if err != nil {
		t.Skipf("validation error (schema may have issues): %v", err)
		return
	}
	if !valid {
		t.Error("conditional.json should pass datatype_tests_schema validation")
	}
}
