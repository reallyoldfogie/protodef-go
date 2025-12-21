package datatypes

import (
	"github.com/tidwall/gjson"
)

// ArrayType represents a ProtoDef 'array' structure.
// Represents a list of values with the same type.
// Arguments:
//   - type: the type of the elements
//   - countType: the type of the length prefix
//   - count: optional (either count or countType), a reference to the field counting the elements, or a fixed size (an integer)
//
// Example: An array of int prefixed by a short length.
//
//	["array", {"countType": "i16", "type": "i32"}]
//
// Example of value: [1, 2, 3, 4]
type ArrayType struct {
	Type      *Type
	Count     int
	CountType *Type
}

func (a *ArrayType) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("type").Exists() {
		a.Type = GetTypeFromJSON("array_type", d.Get("type"))
	}
	if d.Get("count").Exists() {
		a.Count = int(d.Get("count").Int())
	}
	if d.Get("countType").Exists() {
		a.CountType = GetTypeFromJSON("countType", d.Get("countType"))
	}
	return nil
}
