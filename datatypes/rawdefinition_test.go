package datatypes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestRawDefinitionSimpleType(t *testing.T) {
	jsonStr := `"varint"`
	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myvarint", result)

	require.NotNil(t, typ)
	require.Equal(t, `"varint"`, typ.RawDefinition)
	require.Equal(t, "myvarint", typ.Name)
	require.Equal(t, "varint", typ.TypeName)
}

func TestRawDefinitionContainer(t *testing.T) {
	jsonStr := `["container", [
		{"name": "x", "type": "i32"},
		{"name": "z", "type": "i32"}
	]]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("position", result)

	require.NotNil(t, typ)
	require.Equal(t, "container", typ.TypeName)
	require.Equal(t, "position", typ.Name)
	
	// The raw definition should contain the full container definition
	require.Contains(t, typ.RawDefinition, "container")
	require.Contains(t, typ.RawDefinition, "i32")
	
	// Check that container fields also have their raw definitions
	container, ok := typ.Extras.(*Container)
	require.True(t, ok)
	require.Len(t, container.Fields, 2)
	
	// Field types should have their raw definitions
	require.Equal(t, `"i32"`, container.Fields[0].Type.RawDefinition)
	require.Equal(t, `"i32"`, container.Fields[1].Type.RawDefinition)
}

func TestRawDefinitionNestedContainer(t *testing.T) {
	jsonStr := `["container", [
		{"name": "header", "type": "u32"},
		{"anon": true, "type": ["container", [
			{"name": "bitMap", "type": "u16"},
			{"name": "addBitMap", "type": "u16"}
		]]},
		{"name": "footer", "type": "u32"}
	]]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("packet", result)

	require.NotNil(t, typ)
	require.Equal(t, "container", typ.TypeName)
	
	container, ok := typ.Extras.(*Container)
	require.True(t, ok)
	require.Len(t, container.Fields, 3)
	
	// Check the anonymous nested container
	anonField := container.Fields[1]
	require.True(t, anonField.Anon)
	require.NotNil(t, anonField.Type)
	require.Equal(t, "container", anonField.Type.TypeName)
	
	// The nested container should have its raw definition
	require.Contains(t, anonField.Type.RawDefinition, "container")
	require.Contains(t, anonField.Type.RawDefinition, "bitMap")
	require.Contains(t, anonField.Type.RawDefinition, "addBitMap")
}

func TestRawDefinitionArrayType(t *testing.T) {
	jsonStr := `["array", {
		"countType": "varint",
		"type": "u16"
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myarray", result)

	require.NotNil(t, typ)
	require.Equal(t, "array", typ.TypeName)
	require.Equal(t, "myarray", typ.Name)
	
	// The raw definition should contain the full array definition
	require.Contains(t, typ.RawDefinition, "array")
	require.Contains(t, typ.RawDefinition, "countType")
	require.Contains(t, typ.RawDefinition, "varint")
}

func TestRawDefinitionSwitch(t *testing.T) {
	jsonStr := `["switch", {
		"compareTo": "type",
		"fields": {
			"1": "i32",
			"2": "buffer"
		}
	}]`

	result := gjson.Parse(jsonStr)
	typ := GetTypeFromJSON("myswitch", result)

	require.NotNil(t, typ)
	require.Equal(t, "switch", typ.TypeName)
	require.Equal(t, "myswitch", typ.Name)
	
	// The raw definition should contain the full switch definition
	require.Contains(t, typ.RawDefinition, "switch")
	require.Contains(t, typ.RawDefinition, "compareTo")
	require.Contains(t, typ.RawDefinition, "fields")
}
