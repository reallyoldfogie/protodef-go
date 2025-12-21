package datatypes

import (
	"github.com/tidwall/gjson"
)

// MapperType represents a ProtoDef 'mapper' utility type.
// Maps string keys to values.
// Arguments:
//   - type: the type of the input
//   - mappings: a mappings object
//
// Example: Maps a byte to a string, 1 to "byte", 2 to "short", 3 to "int", 4 to "long".
//
//	["mapper", {
//	  "type": "i8",
//	  "mappings": {
//	    "1": "byte",
//	    "2": "short",
//	    "3": "int",
//	    "4": "long"
//	  }
//	}]
//
// Example of value: "int"
type MapperType struct {
	Type     *Type
	Mappings map[string]any
}

func (m *MapperType) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}

	// Parse type field
	if d.Get("type").Exists() {
		m.Type = GetTypeFromJSON("mapper_type", d.Get("type"))
	}

	// Parse mappings field
	if d.Get("mappings").Exists() {
		m.Mappings = make(map[string]any)
		mappingsData := d.Get("mappings")
		if mappingsData.IsObject() {
			mappingsData.ForEach(func(key, value gjson.Result) bool {
				m.Mappings[key.String()] = value.Value()
				return true
			})
		}
	}

	return nil
}
