package datatypes

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestContainerAnonField(t *testing.T) {
	jsonStr := `[
		{"name": "x", "type": "i32"},
		{"name": "z", "type": "i32"},
		{"anon": true, "type": ["container", [
			{"name": "bitMap", "type": "u16"},
			{"name": "addBitMap", "type": "u16"}
		]]}
	]`

	container := &Container{Name: "test"}
	result := gjson.Parse(jsonStr)
	err := container.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(container.Fields) != 3 {
		t.Fatalf("len(container.Fields) = %d, want 3", len(container.Fields))
	}

	// Check first field (named)
	if container.Fields[0].Name != "x" {
		t.Errorf("container.Fields[0].Name = %q, want %q", container.Fields[0].Name, "x")
	}
	if container.Fields[0].Anon {
		t.Error("container.Fields[0].Anon = true, want false")
	}

	// Check second field (named)
	if container.Fields[1].Name != "z" {
		t.Errorf("container.Fields[1].Name = %q, want %q", container.Fields[1].Name, "z")
	}
	if container.Fields[1].Anon {
		t.Error("container.Fields[1].Anon = true, want false")
	}

	// Check third field (anonymous)
	if container.Fields[2].Name != "" {
		t.Errorf("container.Fields[2].Name = %q, want empty", container.Fields[2].Name)
	}
	if !container.Fields[2].Anon {
		t.Error("container.Fields[2].Anon = false, want true")
	}
	if container.Fields[2].Type == nil {
		t.Fatal("container.Fields[2].Type is nil")
	}
	if container.Fields[2].Type.TypeName != "container" {
		t.Errorf("container.Fields[2].Type.TypeName = %q, want %q", container.Fields[2].Type.TypeName, "container")
	}
}

func TestContainerOnlyNamedFields(t *testing.T) {
	jsonStr := `[
		{"name": "id", "type": "varint"},
		{"name": "data", "type": "buffer"}
	]`

	container := &Container{Name: "test"}
	result := gjson.Parse(jsonStr)
	err := container.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(container.Fields) != 2 {
		t.Fatalf("len(container.Fields) = %d, want 2", len(container.Fields))
	}

	for i, field := range container.Fields {
		if field.Anon {
			t.Errorf("container.Fields[%d].Anon = true, want false", i)
		}
		if field.Name == "" {
			t.Errorf("container.Fields[%d].Name is empty", i)
		}
	}
}

func TestContainerMixedFields(t *testing.T) {
	jsonStr := `[
		{"name": "header", "type": "u32"},
		{"anon": true, "type": "u16"},
		{"name": "footer", "type": "u32"}
	]`

	container := &Container{Name: "test"}
	result := gjson.Parse(jsonStr)
	err := container.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(container.Fields) != 3 {
		t.Fatalf("len(container.Fields) = %d, want 3", len(container.Fields))
	}

	// First field: named
	if container.Fields[0].Anon {
		t.Error("container.Fields[0].Anon = true, want false")
	}
	if container.Fields[0].Name != "header" {
		t.Errorf("container.Fields[0].Name = %q, want %q", container.Fields[0].Name, "header")
	}

	// Second field: anonymous
	if !container.Fields[1].Anon {
		t.Error("container.Fields[1].Anon = false, want true")
	}
	if container.Fields[1].Name != "" {
		t.Errorf("container.Fields[1].Name = %q, want empty", container.Fields[1].Name)
	}

	// Third field: named
	if container.Fields[2].Anon {
		t.Error("container.Fields[2].Anon = true, want false")
	}
	if container.Fields[2].Name != "footer" {
		t.Errorf("container.Fields[2].Name = %q, want %q", container.Fields[2].Name, "footer")
	}
}

func TestGetTypeContainer(t *testing.T) {
	jsonStr := `["container", [
		{"name": "x", "type": "i32"},
		{"anon": true, "type": "u16"}
	]]`

	result := gjson.Parse(jsonStr)
	typ := GetType("test_container", result)

	if typ == nil {
		t.Fatal("GetType() returned nil")
	}

	if typ.TypeName != "container" {
		t.Errorf("typ.TypeName = %q, want %q", typ.TypeName, "container")
	}

	container, ok := typ.Extras.(*Container)
	if !ok {
		t.Fatalf("typ.Extras is not *Container, got %T", typ.Extras)
	}

	if len(container.Fields) != 2 {
		t.Fatalf("len(container.Fields) = %d, want 2", len(container.Fields))
	}

	// Check anonymous field is properly parsed
	if !container.Fields[1].Anon {
		t.Error("container.Fields[1].Anon = false, want true")
	}
}

func TestContainerEmptyArray(t *testing.T) {
	jsonStr := `[]`

	container := &Container{Name: "test"}
	result := gjson.Parse(jsonStr)
	err := container.ReadJSON(result)

	if err != nil {
		t.Fatalf("ReadJSON() error = %v", err)
	}

	if len(container.Fields) != 0 {
		t.Errorf("len(container.Fields) = %d, want 0", len(container.Fields))
	}
}

func TestContainerInvalidFormat(t *testing.T) {
	jsonStr := `"not an array"`

	container := &Container{Name: "test"}
	result := gjson.Parse(jsonStr)
	err := container.ReadJSON(result)

	// Should return error for non-array
	if err == nil {
		t.Error("ReadJSON() error = nil, want error for non-array")
	}
}
