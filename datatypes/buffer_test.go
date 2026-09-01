package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestBufferTypeReadJSON(t *testing.T) {
	jsonStr := `{
		"countType": "varint"
	}`

	buffer := &Buffer{}
	result := gjson.Parse(jsonStr)
	err := buffer.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if buffer.CountType == nil {
		t.Fatal("buffer.CountType is nil")
	}

	if buffer.CountType.TypeName != "varint" {
		t.Errorf("buffer.CountType.TypeName = %q, want %q", buffer.CountType.TypeName, "varint")
	}
}

func TestBufferTypeReadJSONWithCount(t *testing.T) {
	jsonStr := `{
		"count": 256
	}`

	buffer := &Buffer{}
	result := gjson.Parse(jsonStr)
	err := buffer.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if buffer.Count != 256 {
		t.Errorf("buffer.Count = %d, want 256", buffer.Count)
	}
}

func TestBufferTypeReadJSONWithRest(t *testing.T) {
	jsonStr := `{
		"rest": true
	}`

	buffer := &Buffer{}
	result := gjson.Parse(jsonStr)
	err := buffer.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if !buffer.Rest {
		t.Error("buffer.Rest should be true")
	}
}

func TestBufferTypeReadJSONAllFields(t *testing.T) {
	jsonStr := `{
		"count": 128,
		"countType": "i16",
		"rest": false
	}`

	buffer := &Buffer{}
	result := gjson.Parse(jsonStr)
	err := buffer.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if buffer.Count != 128 {
		t.Errorf("buffer.Count = %d, want 128", buffer.Count)
	}

	if buffer.CountType == nil {
		t.Fatal("buffer.CountType is nil")
	}

	if buffer.CountType.TypeName != "i16" {
		t.Errorf("buffer.CountType.TypeName = %q, want %q", buffer.CountType.TypeName, "i16")
	}

	if buffer.Rest {
		t.Error("buffer.Rest should be false")
	}
}

func TestBufferTypeReadJSONInvalid(t *testing.T) {
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
			buffer := &Buffer{}
			result := gjson.Parse(tt.jsonStr)
			err := buffer.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestGetTypeBuffer(t *testing.T) {
	jsonStr := `["buffer", {"countType": "varint"}]`

	result := gjson.Parse(jsonStr)
	typ := GetType(result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "buffer" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "buffer")
	}

	buffer, ok := typ.Extras.(*Buffer)
	if !ok {
		t.Fatalf("typ.Extras is not *BufferType, got %T", typ.Extras)
	}

	if buffer.CountType == nil {
		t.Error("buffer.CountType is nil")
	}
}

func TestGetTypeFromJSONBuffer(t *testing.T) {
	jsonStr := `["buffer", {"count": 512, "rest": true}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myBuffer", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myBuffer" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myBuffer")
	}

	if typ.TypeName != "buffer" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "buffer")
	}

	buffer, ok := typ.Extras.(*Buffer)
	if !ok {
		t.Fatalf("typ.Extras is not *BufferType")
	}

	if buffer.Count != 512 {
		t.Errorf("buffer.Count = %d, want 512", buffer.Count)
	}

	if !buffer.Rest {
		t.Error("buffer.Rest should be true")
	}
}
