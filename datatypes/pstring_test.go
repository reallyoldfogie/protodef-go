package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestPStringTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"countType": "varint"
	}`

	pstring := &PStringType{}
	result := gjson.Parse(jsonStr)
	err := pstring.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if pstring.CountType == nil {
		t.Fatal("pstring.CountType is nil")
	}

	if pstring.CountType.TypeName != "varint" {
		t.Errorf("pstring.CountType.TypeName = %q, want %q", pstring.CountType.TypeName, "varint")
	}
}

func TestPStringTypeReadJSONWithCountInt(t *testing.T) {
	jsonStr := `{
		"count": 32
	}`

	pstring := &PStringType{}
	result := gjson.Parse(jsonStr)
	err := pstring.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if pstring.Count == nil {
		t.Fatal("pstring.Count is nil")
	}

	countVal, ok := pstring.Count.(int)
	if !ok {
		t.Fatalf("pstring.Count is not int, got %T", pstring.Count)
	}

	if countVal != 32 {
		t.Errorf("pstring.Count = %d, want 32", countVal)
	}
}

func TestPStringTypeReadJSONWithCountField(t *testing.T) {
	jsonStr := `{
		"count": "length_field"
	}`

	pstring := &PStringType{}
	result := gjson.Parse(jsonStr)
	err := pstring.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if pstring.Count == nil {
		t.Fatal("pstring.Count is nil")
	}

	// When count is not a number, it's stored as the raw value
	countStr, ok := pstring.Count.(string)
	if !ok {
		t.Fatalf("pstring.Count is not string, got %T", pstring.Count)
	}

	if countStr != "length_field" {
		t.Errorf("pstring.Count = %q, want %q", countStr, "length_field")
	}
}

func TestPStringTypeReadJSONWithEncoding(t *testing.T) {
	jsonStr := `{
		"countType": "i16",
		"encoding": "utf-8"
	}`

	pstring := &PStringType{}
	result := gjson.Parse(jsonStr)
	err := pstring.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if pstring.CountType == nil {
		t.Fatal("pstring.CountType is nil")
	}

	if pstring.Encoding != "utf-8" {
		t.Errorf("pstring.Encoding = %q, want %q", pstring.Encoding, "utf-8")
	}
}

func TestPStringTypeReadJSONAllFields(t *testing.T) {
	jsonStr := `{
		"countType": "u8",
		"count": 16,
		"encoding": "ascii"
	}`

	pstring := &PStringType{}
	result := gjson.Parse(jsonStr)
	err := pstring.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if pstring.CountType == nil {
		t.Fatal("pstring.CountType is nil")
	}

	if pstring.CountType.TypeName != "u8" {
		t.Errorf("pstring.CountType.TypeName = %q, want %q", pstring.CountType.TypeName, "u8")
	}

	countVal, ok := pstring.Count.(int)
	if !ok {
		t.Fatalf("pstring.Count is not int, got %T", pstring.Count)
	}

	if countVal != 16 {
		t.Errorf("pstring.Count = %d, want 16", countVal)
	}

	if pstring.Encoding != "ascii" {
		t.Errorf("pstring.Encoding = %q, want %q", pstring.Encoding, "ascii")
	}
}

func TestPStringTypeReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an object", `"string"`},
		{"array instead of object", `[]`},
		{"number", `123`},
		{"empty object", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pstring := &PStringType{}
			result := gjson.Parse(tt.jsonStr)
			err := pstring.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestGetTypePString(t *testing.T) {
	jsonStr := `["pstring", {"countType": "varint"}]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_pstring", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "pstring" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "pstring")
	}

	pstring, ok := typ.Extras.(*PStringType)
	if !ok {
		t.Fatalf("typ.Extras is not *PStringType, got %T", typ.Extras)
	}

	if pstring.CountType == nil {
		t.Error("pstring.CountType is nil")
	}
}

func TestGetTypeFromJSONPString(t *testing.T) {
	jsonStr := `["pstring", {"count": 128, "encoding": "utf-16"}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myPString", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myPString" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myPString")
	}

	if typ.TypeName != "pstring" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "pstring")
	}

	pstring, ok := typ.Extras.(*PStringType)
	if !ok {
		t.Fatalf("typ.Extras is not *PStringType")
	}

	if pstring.Encoding != "utf-16" {
		t.Errorf("pstring.Encoding = %q, want %q", pstring.Encoding, "utf-16")
	}
}
