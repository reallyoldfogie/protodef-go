package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestSwitchTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"compareTo": "packetId",
		"fields": {
			"0": "u8",
			"1": "i32",
			"2": "varint"
		}
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if sw.CompareTo != "packetId" {
		t.Errorf("sw.CompareTo = %q, want %q", sw.CompareTo, "packetId")
	}

	if sw.Fields == nil {
		t.Fatal("sw.Fields is nil")
	}

	if len(sw.Fields) != 3 {
		t.Fatalf("len(sw.Fields) = %d, want 3", len(sw.Fields))
	}

	if sw.Fields["0"] == nil {
		t.Error("sw.Fields[\"0\"] is nil")
	} else if sw.Fields["0"].TypeName != "u8" {
		t.Errorf("sw.Fields[\"0\"].TypeName = %q, want %q", sw.Fields["0"].TypeName, "u8")
	}

	if sw.Fields["1"] == nil {
		t.Error("sw.Fields[\"1\"] is nil")
	} else if sw.Fields["1"].TypeName != "i32" {
		t.Errorf("sw.Fields[\"1\"].TypeName = %q, want %q", sw.Fields["1"].TypeName, "i32")
	}

	if sw.Fields["2"] == nil {
		t.Error("sw.Fields[\"2\"] is nil")
	} else if sw.Fields["2"].TypeName != "varint" {
		t.Errorf("sw.Fields[\"2\"].TypeName = %q, want %q", sw.Fields["2"].TypeName, "varint")
	}
}

func TestSwitchTypeReadJSONWithCompareToValue(t *testing.T) {
	jsonStr := `{
		"compareTo": "type",
		"compareToValue": "default_value",
		"fields": {
			"a": "type_a",
			"b": "type_b"
		}
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if sw.CompareTo != "type" {
		t.Errorf("sw.CompareTo = %q, want %q", sw.CompareTo, "type")
	}

	if sw.CompareToValue == nil {
		t.Fatal("sw.CompareToValue is nil")
	}

	compareToVal, ok := sw.CompareToValue.(string)
	if !ok {
		t.Fatalf("sw.CompareToValue is not string, got %T", sw.CompareToValue)
	}

	if compareToVal != "default_value" {
		t.Errorf("sw.CompareToValue = %q, want %q", compareToVal, "default_value")
	}

	if len(sw.Fields) != 2 {
		t.Fatalf("len(sw.Fields) = %d, want 2", len(sw.Fields))
	}
}

func TestSwitchTypeReadJSONWithDefault(t *testing.T) {
	jsonStr := `{
		"compareTo": "mode",
		"fields": {
			"1": "u8",
			"2": "i16"
		},
		"default": "i32"
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if sw.Default == nil {
		t.Fatal("sw.Default is nil")
	} else if sw.Default.TypeName != "i32" {
		t.Errorf("sw.Default.TypeName = %q, want %q", sw.Default.TypeName, "i32")
	}
}

func TestSwitchTypeReadJSONEmptyFields(t *testing.T) {
	jsonStr := `{
		"compareTo": "field",
		"fields": {}
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if sw.CompareTo != "field" {
		t.Errorf("sw.CompareTo = %q, want %q", sw.CompareTo, "field")
	}

	if sw.Fields == nil {
		t.Fatal("sw.Fields should not be nil")
	}

	if len(sw.Fields) != 0 {
		t.Errorf("len(sw.Fields) = %d, want 0", len(sw.Fields))
	}
}

func TestSwitchTypeReadJSONAllFields(t *testing.T) {
	jsonStr := `{
		"compareTo": "version",
		"compareToValue": 1,
		"fields": {
			"47": "u8",
			"107": "u16",
			"393": "u32"
		},
		"default": "varint"
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if sw.CompareTo != "version" {
		t.Errorf("sw.CompareTo = %q, want %q", sw.CompareTo, "version")
	}

	if sw.CompareToValue == nil {
		t.Fatal("sw.CompareToValue is nil")
	}

	if sw.Fields == nil {
		t.Fatal("sw.Fields is nil")
	}

	if len(sw.Fields) != 3 {
		t.Errorf("len(sw.Fields) = %d, want 3", len(sw.Fields))
	}

	if sw.Default == nil {
		t.Fatal("sw.Default is nil")
	} else if sw.Default.TypeName != "varint" {
		t.Errorf("sw.Default.TypeName = %q, want %q", sw.Default.TypeName, "varint")
	}
}

func TestSwitchTypeReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an object", `"string"`},
		{"array instead of object", `[]`},
		{"number", `123`},
		{"empty object", `{}`},
		{"missing compareTo", `{"fields": {"1": "a"}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sw := &SwitchType{}
			result := gjson.Parse(tt.jsonStr)
			err := sw.ReadJSON(result)

			// Should not error, just leave fields uninitialized or partially initialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestGetTypeSwitch(t *testing.T) {
	jsonStr := `["switch", {
		"compareTo": "type",
		"fields": {
			"0": "type_zero",
			"1": "type_one"
		}
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_switch", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "switch" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "switch")
	}

	sw, ok := typ.Extras.(*SwitchType)
	if !ok {
		t.Fatalf("typ.Extras is not *SwitchType, got %T", typ.Extras)
	}

	if sw.CompareTo != "type" {
		t.Errorf("sw.CompareTo = %q, want %q", sw.CompareTo, "type")
	}

	if len(sw.Fields) != 2 {
		t.Errorf("len(sw.Fields) = %d, want 2", len(sw.Fields))
	}
}

func TestGetTypeFromJSONSwitch(t *testing.T) {
	jsonStr := `["switch", {
		"compareTo": "state",
		"fields": {
			"active": "u8",
			"inactive": "i32"
		},
		"default": "varint"
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("mySwitch", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "mySwitch" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "mySwitch")
	}

	if typ.TypeName != "switch" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "switch")
	}

	sw, ok := typ.Extras.(*SwitchType)
	if !ok {
		t.Fatalf("typ.Extras is not *SwitchType")
	}

	if sw.Default == nil {
		t.Error("sw.Default is nil")
	}
}

func TestSwitchTypeWithNonStringDefault(t *testing.T) {
	// Test that non-string default values are handled properly
	jsonStr := `{
		"compareTo": "mode",
		"fields": {"1": "one"},
		"default": 123
	}`

	sw := &SwitchType{}
	result := gjson.Parse(jsonStr)
	err := sw.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	// When default is not a string, it should be nil
	if sw.Default != nil {
		t.Errorf("sw.Default should be nil for non-string default, got %v", sw.Default)
	}
}
