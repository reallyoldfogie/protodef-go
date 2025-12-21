package datatypes

import (
	"github.com/tidwall/gjson"
)

// BufferType represents a ProtoDef 'buffer' utility type.
// Represents a raw bytes with count prefix/field or without it.
// Arguments:
//   - countType: the type of the length prefix
//   - count: optional (either count or countType), a reference to the field counting the elements, or a fixed size (an integer)
//   - rest: optional, represent rest bytes as-is
//
// Example: ["buffer", {"countType": "varint"}]
// Example of value: Buffer <01 02 03>
type BufferType struct {
	Count     int
	CountType *Type
	Rest      bool
}

func (b *BufferType) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("count").Exists() {
		b.Count = int(d.Get("count").Int())
	}
	if d.Get("countType").Exists() {
		b.CountType = GetTypeFromJSON("countType", d.Get("countType"))
	}
	if d.Get("rest").Exists() {
		b.Rest = d.Get("rest").Bool()
	}
	return nil
}
