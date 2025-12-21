package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestArrayTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"type": "i32",
		"countType": "i16"
	}`

	array := &ArrayType{}
	result := gjson.Parse(jsonStr)
	err := array.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if array.Type == nil {
		t.Error("array.Type is nil")
	} else if array.Type.TypeName != "i32" {
		t.Errorf("array.Type.TypeName = %q, want %q", array.Type.TypeName, "i32")
	}

	if array.CountType == nil {
		t.Error("array.CountType is nil")
	} else if array.CountType.TypeName != "i16" {
		t.Errorf("array.CountType.TypeName = %q, want %q", array.CountType.TypeName, "i16")
	}
}

func TestArrayTypeReadJSONWithCount(t *testing.T) {
	jsonStr := `{
		"type": "u8",
		"count": 10
	}`

	array := &ArrayType{}
	result := gjson.Parse(jsonStr)
	err := array.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if array.Type == nil {
		t.Fatal("array.Type is nil")
	}

	if array.Type.TypeName != "u8" {
		t.Errorf("array.Type.TypeName = %q, want %q", array.Type.TypeName, "u8")
	}

	if array.Count != 10 {
		t.Errorf("array.Count = %d, want 10", array.Count)
	}
}

func TestArrayTypeReadJSONWithVarint(t *testing.T) {
	jsonStr := `{
		"type": "i32",
		"countType": "varint"
	}`

	array := &ArrayType{}
	result := gjson.Parse(jsonStr)
	err := array.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if array.Type == nil {
		t.Fatal("array.Type is nil")
	}

	if array.Type.TypeName != "i32" {
		t.Errorf("array.Type.TypeName = %q, want %q", array.Type.TypeName, "i32")
	}

	if array.CountType == nil {
		t.Fatal("array.CountType is nil")
	}

	if array.CountType.TypeName != "varint" {
		t.Errorf("array.CountType.TypeName = %q, want %q", array.CountType.TypeName, "varint")
	}
}

func TestArrayTypeReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an object", `"string"`},
		{"missing type", `{"countType": "i16"}`},
		{"missing count info", `{"type": "i32"}`},
		{"empty object", `{}`},
		{"number instead of object", `123`},
		{"array instead of object", `[]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			array := &ArrayType{}
			result := gjson.Parse(tt.jsonStr)
			err := array.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestArrayTypeReadJSONAllFields(t *testing.T) {
	jsonStr := `{
		"type": "i64",
		"count": 5,
		"countType": "u32"
	}`

	array := &ArrayType{}
	result := gjson.Parse(jsonStr)
	err := array.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if array.Type == nil {
		t.Fatal("array.Type is nil")
	}

	if array.Type.TypeName != "i64" {
		t.Errorf("array.Type.TypeName = %q, want %q", array.Type.TypeName, "i64")
	}

	if array.Count != 5 {
		t.Errorf("array.Count = %d, want 5", array.Count)
	}

	if array.CountType == nil {
		t.Fatal("array.CountType is nil")
	}

	if array.CountType.TypeName != "u32" {
		t.Errorf("array.CountType.TypeName = %q, want %q", array.CountType.TypeName, "u32")
	}
}

func TestGetTypeArray(t *testing.T) {
	jsonStr := `["array", {
		"countType": "i16",
		"type": "i32"
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_array", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "array" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "array")
	}

	array, ok := typ.Extras.(*ArrayType)
	if !ok {
		t.Fatalf("typ.Extras is not *ArrayType, got %T", typ.Extras)
	}

	if array.Type == nil {
		t.Error("array.Type is nil")
	}

	if array.CountType == nil {
		t.Error("array.CountType is nil")
	}
}

func TestGetTypeFromJSONArray(t *testing.T) {
	jsonStr := `["array", {"type": "u16", "count": 8}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myArray", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myArray" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myArray")
	}

	if typ.TypeName != "array" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "array")
	}

	array, ok := typ.Extras.(*ArrayType)
	if !ok {
		t.Fatalf("typ.Extras is not *ArrayType")
	}

	if array.Count != 8 {
		t.Errorf("array.Count = %d, want 8", array.Count)
	}
}
