package datatypes

import (
	"github.com/tidwall/gjson"
)

// Array represents a ProtoDef 'array' structure.
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
type Array struct {
	name           string
	Type           *Type
	Count          int          // Fixed count or 0 if count is a field reference
	CountFieldName string       // Name of field to use as count (if Count == 0 and this is set)
	CountType      *Type
}

func (a *Array) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("type").Exists() {
		a.Type = GetTypeFromJSON("array_type", d.Get("type"))
	}
	if d.Get("count").Exists() {
		countVal := d.Get("count")
		if countVal.Type == gjson.String {
			// String value means field reference
			a.CountFieldName = countVal.String()
			a.Count = 0
		} else {
			// Numeric value means fixed count
			a.Count = int(countVal.Int())
			a.CountFieldName = ""
		}
	}
	if d.Get("countType").Exists() {
		a.CountType = GetTypeFromJSON("countType", d.Get("countType"))
	}
	return nil
}

func (a *Array) SetName(name string) {
	// Array doesn't have a direct name field, but we store it for reference
	a.name = name
}

func (a *Array) GetName() string {
	return a.name
}

func (a *Array) Clone() TypeExtras {
	cloned := &Array{
		Count:          a.Count,
		CountFieldName: a.CountFieldName,
	}
	if a.Type != nil {
		clonedType := *a.Type
		if a.Type.Extras != nil {
			clonedType.Extras = a.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}
	if a.CountType != nil {
		clonedCountType := *a.CountType
		if a.CountType.Extras != nil {
			clonedCountType.Extras = a.CountType.Extras.Clone()
		}
		cloned.CountType = &clonedCountType
	}
	return cloned
}

func (a *Array) UpdateContainedNames(updatedNames map[string]string) {
	if a.Type != nil && a.Type.Extras != nil {
		a.Type.Extras.UpdateContainedNames(updatedNames)
	}
	if a.CountType != nil && a.CountType.Extras != nil {
		a.CountType.Extras.UpdateContainedNames(updatedNames)
	}
}
