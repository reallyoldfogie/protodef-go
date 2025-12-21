package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestMapperTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"type": "i8",
		"mappings": {
			"1": "byte",
			"2": "short",
			"3": "int",
			"4": "long"
		}
	}`

	mapper := &Mapper{}
	result := gjson.Parse(jsonStr)
	err := mapper.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	// Check type was parsed
	if mapper.Type == nil {
		t.Error("mapper.Type is nil")
	} else if mapper.Type.TypeName != "i8" {
		t.Errorf("mapper.Type.TypeName = %q, want %q", mapper.Type.TypeName, "i8")
	}

	// Check mappings were parsed
	if mapper.Mappings == nil {
		t.Fatal("mapper.Mappings is nil")
	}

	expectedMappings := map[string]string{
		"1": "byte",
		"2": "short",
		"3": "int",
		"4": "long",
	}

	for key, expected := range expectedMappings {
		got, ok := mapper.Mappings[key]
		if !ok {
			t.Errorf("mapper.Mappings[%q] not found", key)
			continue
		}
		if got != expected {
			t.Errorf("mapper.Mappings[%q] = %v, want %q", key, got, expected)
		}
	}
}

func TestMapperTypeReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an object", `"string"`},
		{"missing type", `{"mappings": {"1": "a"}}`},
		{"missing mappings", `{"type": "i8"}`},
		{"empty object", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := &Mapper{}
			result := gjson.Parse(tt.jsonStr)
			err := mapper.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestMapperTypeWithComplexType(t *testing.T) {
	jsonStr := `{
		"type": "varint",
		"mappings": {
			"0": "handshake",
			"1": "status_request",
			"2": "login_start"
		}
	}`

	mapper := &Mapper{}
	result := gjson.Parse(jsonStr)
	err := mapper.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if mapper.Type == nil {
		t.Fatal("mapper.Type is nil")
	}

	if mapper.Type.TypeName != "varint" {
		t.Errorf("mapper.Type.TypeName = %q, want %q", mapper.Type.TypeName, "varint")
	}

	if len(mapper.Mappings) != 3 {
		t.Errorf("len(mapper.Mappings) = %d, want 3", len(mapper.Mappings))
	}
}

func TestGetTypeMapper(t *testing.T) {
	jsonStr := `["mapper", {
		"type": "u8",
		"mappings": {
			"0": "inactive",
			"1": "active"
		}
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetType(result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "mapper" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "mapper")
	}

	mapper, ok := typ.Extras.(*Mapper)
	if !ok {
		t.Fatalf("typ.Extras is not *MapperType, got %T", typ.Extras)
	}

	if mapper.Type == nil {
		t.Error("mapper.Type is nil")
	}

	if len(mapper.Mappings) != 2 {
		t.Errorf("len(mapper.Mappings) = %d, want 2", len(mapper.Mappings))
	}
}
