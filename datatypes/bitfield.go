package datatypes

import (
	"github.com/tidwall/gjson"
)

// Bitfield represents a ProtoDef 'bitfield' utility type.
// Represents a list of values with sizes that are not a multiple of 8 bits.
// The sum of the sizes must be a multiple of 8.
// Arguments:
//   - name: the name of the field
//   - size: the size in bits
//   - signed: is value signed
//
// Example: ["bitfield", [
//
//	{"name": "x", "size": 26, "signed": true},
//	{"name": "y", "size": 12, "signed": true},
//	{"name": "z", "size": 26, "signed": true}
//
// ]]
// Example of value: {"x": 10, "y": 10, "z": 10}
type Bitfield struct {
	name   string
	Fields []BitFieldField
}

type BitFieldField struct {
	Name   string
	Size   int
	Signed bool
}

func (bf *Bitfield) ReadJSON(d gjson.Result) error {
	if !d.IsArray() {
		return nil
	}
	fields := d.Array()
	for _, f := range fields {
		bf.Fields = append(bf.Fields, BitFieldField{
			Name:   f.Get("name").String(),
			Size:   int(f.Get("size").Int()),
			Signed: f.Get("signed").Bool(),
		})
	}
	return nil
}

func (bf *Bitfield) SetName(name string) {
	bf.name = name
}

func (bf *Bitfield) GetName() string {
	return bf.name
}

func (bf *Bitfield) Clone() TypeExtras {
	cloned := &Bitfield{
		name:   bf.name,
		Fields: make([]BitFieldField, len(bf.Fields)),
	}
	copy(cloned.Fields, bf.Fields)
	return cloned
}

func (bf *Bitfield) UpdateContainedNames(updatedNames map[string]string) {
	for i, field := range bf.Fields {
		if newName, exists := updatedNames[field.Name]; exists {
			bf.Fields[i].Name = newName
		}
	}
}
