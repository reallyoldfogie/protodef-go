package datatypes

import (
	"github.com/tidwall/gjson"
)

// Count represents a ProtoDef 'count' structure type.
// Represents a count field for an array or a buffer.
// Arguments:
//   - type: the type of the count
//   - countFor: a field to count for
//
// Example: A count for a field name records, of type short.
//
//	["count", {"type": "i16", "countFor": "records"}]
//
// Example of value: 4
type Count struct {
	name     string
	Type     *Type
	CountFor string
}

func (c *Count) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("type").Exists() {
		c.Type = GetTypeFromJSON("count_type", d.Get("type"))
	}
	if d.Get("countFor").Exists() {
		c.CountFor = d.Get("countFor").String()
	}
	return nil
}

func (c *Count) SetName(name string) {
	c.name = name
}

func (c *Count) GetName() string {
	return c.name
}

func (c *Count) Clone() TypeExtras {
	cloned := &Count{
		name:     c.name,
		CountFor: c.CountFor,
	}
	if c.Type != nil {
		clonedType := *c.Type
		if c.Type.Extras != nil {
			clonedType.Extras = c.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}
	return cloned
}

func (c *Count) UpdateContainedNames(updatedNames map[string]string) {
	if newName, exists := updatedNames[c.CountFor]; exists {
		c.CountFor = newName
	}
	if c.Type != nil && c.Type.Extras != nil {
		c.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
