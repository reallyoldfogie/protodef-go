package datatypes_test

import (
	"testing"
)

func TestDatatypeTestsSchemaValidation(t *testing.T) {
	// This would validate the datatype_tests_schema.json itself against a JSON Schema meta-schema
	// Skipping for now as the schema has strict mode issues with ajv
	t.Skip("Schema self-validation not yet implemented")
}
