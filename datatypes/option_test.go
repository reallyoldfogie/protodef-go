package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestOptionTypeReadJSON(t *testing.T) {
	jsonStr := `"i32"`

	option := &OptionType{}
	result := gjson.Parse(jsonStr)
	err := option.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if option.Type == nil {
		t.Fatal("option.Type is nil")
	}

	if option.Type.TypeName != "i32" {
		t.Errorf("option.Type.TypeName = %q, want %q", option.Type.TypeName, "i32")
	}
}

func TestOptionTypeReadJSONComplexType(t *testing.T) {
	jsonStr := `"varint"`

	option := &OptionType{}
	result := gjson.Parse(jsonStr)
	err := option.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if option.Type == nil {
		t.Fatal("option.Type is nil")
	}

	if option.Type.TypeName != "varint" {
		t.Errorf("option.Type.TypeName = %q, want %q", option.Type.TypeName, "varint")
	}
}

func TestOptionTypeReadJSONWithArray(t *testing.T) {
	jsonStr := `["array", {"type": "u8", "count": 10}]`

	option := &OptionType{}
	result := gjson.Parse(jsonStr)
	err := option.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if option.Type == nil {
		t.Fatal("option.Type is nil")
	}

	if option.Type.TypeName != "array" {
		t.Errorf("option.Type.TypeName = %q, want %q", option.Type.TypeName, "array")
	}
}

func TestGetTypeOption(t *testing.T) {
	jsonStr := `["option", "i64"]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_option", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "option" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "option")
	}

	option, ok := typ.Extras.(*OptionType)
	if !ok {
		t.Fatalf("typ.Extras is not *OptionType, got %T", typ.Extras)
	}

	if option.Type == nil {
		t.Fatal("option.Type is nil")
	}

	if option.Type.TypeName != "i64" {
		t.Errorf("option.Type.TypeName = %q, want %q", option.Type.TypeName, "i64")
	}
}

func TestGetTypeFromJSONOption(t *testing.T) {
	jsonStr := `["option", "i32"]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myOption", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myOption" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myOption")
	}

	if typ.TypeName != "option" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "option")
	}

	option, ok := typ.Extras.(*OptionType)
	if !ok {
		t.Fatalf("typ.Extras is not *OptionType")
	}

	if option.Type == nil {
		t.Fatal("option.Type is nil")
	} else if option.Type.TypeName != "i32" {
		t.Errorf("option.Type.TypeName = %q, want %q", option.Type.TypeName, "i32")
	}
}

func TestOptionTypeReadJSONNumber(t *testing.T) {
	// Test that number values are handled properly
	jsonStr := `123`

	option := &OptionType{}
	result := gjson.Parse(jsonStr)
	err := option.ReadJSON(result)

	// Should not error
	if err != nil {
		t.Errorf("ReadJSON() error = %v, want nil", err)
	}
}

func TestOptionTypeReadJSONObject(t *testing.T) {
	jsonStr := `{"type": "test"}`

	option := &OptionType{}
	result := gjson.Parse(jsonStr)
	err := option.ReadJSON(result)

	// Should not error
	if err != nil {
		t.Errorf("ReadJSON() error = %v, want nil", err)
	}
}
