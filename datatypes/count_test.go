package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestCountTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"type": "i16",
		"countFor": "records"
	}`

	count := &CountType{}
	result := gjson.Parse(jsonStr)
	err := count.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if count.Type == nil {
		t.Fatal("count.Type is nil")
	}

	if count.Type.TypeName != "i16" {
		t.Errorf("count.Type.TypeName = %q, want %q", count.Type.TypeName, "i16")
	}

	if count.CountFor != "records" {
		t.Errorf("count.CountFor = %q, want %q", count.CountFor, "records")
	}
}

func TestCountTypeReadJSONVarint(t *testing.T) {
	jsonStr := `{
		"type": "varint",
		"countFor": "items"
	}`

	count := &CountType{}
	result := gjson.Parse(jsonStr)
	err := count.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if count.Type == nil {
		t.Fatal("count.Type is nil")
	}

	if count.Type.TypeName != "varint" {
		t.Errorf("count.Type.TypeName = %q, want %q", count.Type.TypeName, "varint")
	}

	if count.CountFor != "items" {
		t.Errorf("count.CountFor = %q, want %q", count.CountFor, "items")
	}
}

func TestCountTypeReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an object", `"string"`},
		{"missing type", `{"countFor": "items"}`},
		{"missing countFor", `{"type": "i16"}`},
		{"empty object", `{}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := &CountType{}
			result := gjson.Parse(tt.jsonStr)
			err := count.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestGetTypeCount(t *testing.T) {
	jsonStr := `["count", {
		"type": "u32",
		"countFor": "packets"
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_count", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "count" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "count")
	}

	count, ok := typ.Extras.(*CountType)
	if !ok {
		t.Fatalf("typ.Extras is not *CountType, got %T", typ.Extras)
	}

	if count.Type == nil {
		t.Fatal("count.Type is nil")
	}

	if count.Type.TypeName != "u32" {
		t.Errorf("count.Type.TypeName = %q, want %q", count.Type.TypeName, "u32")
	}

	if count.CountFor != "packets" {
		t.Errorf("count.CountFor = %q, want %q", count.CountFor, "packets")
	}
}

func TestGetTypeFromJSONCount(t *testing.T) {
	jsonStr := `["count", {"type": "i16", "countFor": "data"}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myCount", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myCount" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myCount")
	}

	if typ.TypeName != "count" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "count")
	}

	count, ok := typ.Extras.(*CountType)
	if !ok {
		t.Fatalf("typ.Extras is not *CountType")
	}

	if count.CountFor != "data" {
		t.Errorf("count.CountFor = %q, want %q", count.CountFor, "data")
	}
}
