package datatypes

import (
	"github.com/tidwall/gjson"
)

// Buffer represents a ProtoDef 'buffer' utility type.
// Represents a raw bytes with count prefix/field or without it.
// Arguments:
//   - countType: the type of the length prefix
//   - count: optional (either count or countType), a reference to the field counting the elements, or a fixed size (an integer)
//   - rest: optional, represent rest bytes as-is
//
// Example: ["buffer", {"countType": "varint"}]
// Example of value: Buffer <01 02 03>
type Buffer struct {
	name      string
	Count     int
	CountType *Type
	Rest      bool
}

func (b *Buffer) ReadJSON(d gjson.Result) error {
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

func (b *Buffer) SetName(name string) {
	b.name = name
}

func (b *Buffer) GetName() string {
	return b.name
}

func (b *Buffer) Clone() TypeExtras {
	cloned := &Buffer{
		name:  b.name,
		Count: b.Count,
		Rest:  b.Rest,
	}
	if b.CountType != nil {
		clonedCountType := *b.CountType
		if b.CountType.Extras != nil {
			clonedCountType.Extras = b.CountType.Extras.Clone()
		}
		cloned.CountType = &clonedCountType
	}
	return cloned
}

func (b *Buffer) UpdateContainedNames(updatedNames map[string]string) {
	if b.CountType != nil && b.CountType.Extras != nil {
		b.CountType.Extras.UpdateContainedNames(updatedNames)
	}
}
