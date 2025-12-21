package datatypes

import (
	"github.com/tidwall/gjson"
)

// CountType represents a ProtoDef 'count' structure type.
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
type CountType struct {
	Type     *Type
	CountFor string
}

func (c *CountType) ReadJSON(d gjson.Result) error {
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
