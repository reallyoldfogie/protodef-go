package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestBitFieldReadJSON(t *testing.T) {
	jsonStr := `[
		{"name": "x", "size": 26, "signed": true},
		{"name": "y", "size": 12, "signed": true},
		{"name": "z", "size": 26, "signed": true}
	]`

	bf := &Bitfield{}
	result := gjson.Parse(jsonStr)
	err := bf.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(bf.Fields) != 3 {
		t.Fatalf("len(bf.Fields) = %d, want 3", len(bf.Fields))
	}

	// Check first field
	if bf.Fields[0].Name != "x" {
		t.Errorf("bf.Fields[0].Name = %q, want %q", bf.Fields[0].Name, "x")
	}
	if bf.Fields[0].Size != 26 {
		t.Errorf("bf.Fields[0].Size = %d, want 26", bf.Fields[0].Size)
	}
	if !bf.Fields[0].Signed {
		t.Error("bf.Fields[0].Signed should be true")
	}

	// Check second field
	if bf.Fields[1].Name != "y" {
		t.Errorf("bf.Fields[1].Name = %q, want %q", bf.Fields[1].Name, "y")
	}
	if bf.Fields[1].Size != 12 {
		t.Errorf("bf.Fields[1].Size = %d, want 12", bf.Fields[1].Size)
	}
	if !bf.Fields[1].Signed {
		t.Error("bf.Fields[1].Signed should be true")
	}

	// Check third field
	if bf.Fields[2].Name != "z" {
		t.Errorf("bf.Fields[2].Name = %q, want %q", bf.Fields[2].Name, "z")
	}
	if bf.Fields[2].Size != 26 {
		t.Errorf("bf.Fields[2].Size = %d, want 26", bf.Fields[2].Size)
	}
	if !bf.Fields[2].Signed {
		t.Error("bf.Fields[2].Signed should be true")
	}
}

func TestBitFieldReadJSONUnsigned(t *testing.T) {
	jsonStr := `[
		{"name": "flags", "size": 8, "signed": false},
		{"name": "data", "size": 16, "signed": false}
	]`

	bf := &Bitfield{}
	result := gjson.Parse(jsonStr)
	err := bf.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(bf.Fields) != 2 {
		t.Fatalf("len(bf.Fields) = %d, want 2", len(bf.Fields))
	}

	if bf.Fields[0].Name != "flags" {
		t.Errorf("bf.Fields[0].Name = %q, want %q", bf.Fields[0].Name, "flags")
	}
	if bf.Fields[0].Size != 8 {
		t.Errorf("bf.Fields[0].Size = %d, want 8", bf.Fields[0].Size)
	}
	if bf.Fields[0].Signed {
		t.Error("bf.Fields[0].Signed should be false")
	}

	if bf.Fields[1].Name != "data" {
		t.Errorf("bf.Fields[1].Name = %q, want %q", bf.Fields[1].Name, "data")
	}
	if bf.Fields[1].Size != 16 {
		t.Errorf("bf.Fields[1].Size = %d, want 16", bf.Fields[1].Size)
	}
	if bf.Fields[1].Signed {
		t.Error("bf.Fields[1].Signed should be false")
	}
}

func TestBitFieldReadJSONSingleField(t *testing.T) {
	jsonStr := `[{"name": "value", "size": 32, "signed": false}]`

	bf := &Bitfield{}
	result := gjson.Parse(jsonStr)
	err := bf.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(bf.Fields) != 1 {
		t.Fatalf("len(bf.Fields) = %d, want 1", len(bf.Fields))
	}

	if bf.Fields[0].Name != "value" {
		t.Errorf("bf.Fields[0].Name = %q, want %q", bf.Fields[0].Name, "value")
	}
	if bf.Fields[0].Size != 32 {
		t.Errorf("bf.Fields[0].Size = %d, want 32", bf.Fields[0].Size)
	}
}

func TestBitFieldReadJSONEmpty(t *testing.T) {
	jsonStr := `[]`

	bf := &Bitfield{}
	result := gjson.Parse(jsonStr)
	err := bf.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(bf.Fields) != 0 {
		t.Errorf("len(bf.Fields) = %d, want 0", len(bf.Fields))
	}
}

func TestBitFieldReadJSONInvalid(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
	}{
		{"not an array", `"string"`},
		{"object instead of array", `{"name": "x", "size": 8}`},
		{"number", `123`},
		{"null", `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &Bitfield{}
			result := gjson.Parse(tt.jsonStr)
			err := bf.ReadJSON(result)

			// Should not error, just leave fields uninitialized
			if err != nil {
				t.Errorf("ReadJSON() error = %v, want nil", err)
			}
		})
	}
}

func TestGetTypeBitField(t *testing.T) {
	jsonStr := `["bitfield", [
		{"name": "a", "size": 4, "signed": false},
		{"name": "b", "size": 4, "signed": false}
	]]`

	result := gjson.Parse(jsonStr)
	typ := GetType(result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "bitfield" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "bitfield")
	}

	bf, ok := typ.Extras.(*Bitfield)
	if !ok {
		t.Fatalf("typ.Extras is not *BitField, got %T", typ.Extras)
	}

	if len(bf.Fields) != 2 {
		t.Errorf("len(bf.Fields) = %d, want 2", len(bf.Fields))
	}
}

func TestGetTypeFromJSONBitField(t *testing.T) {
	jsonStr := `["bitfield", [
		{"name": "field1", "size": 16, "signed": true}
	]]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myBitField", result)

	if typ == nil {
		t.Fatal("GetTypeFromJSON() returned nil")
	}

	if typ.Name != "myBitField" {
		t.Errorf("typ.Name = %q, want %q", typ.Name, "myBitField")
	}

	if typ.TypeName != "bitfield" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "bitfield")
	}

	bf, ok := typ.Extras.(*Bitfield)
	if !ok {
		t.Fatalf("typ.Extras is not *BitField")
	}

	if len(bf.Fields) != 1 {
		t.Errorf("len(bf.Fields) = %d, want 1", len(bf.Fields))
	}

	if bf.Fields[0].Name != "field1" {
		t.Errorf("bf.Fields[0].Name = %q, want %q", bf.Fields[0].Name, "field1")
	}
}
