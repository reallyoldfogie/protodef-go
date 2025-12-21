package datatypes

import (
	"github.com/tidwall/gjson"
)

// Bitflags represents a ProtoDef 'bitflags' utility type.
// Represents a set of boolean flags stored in an integer type.
// Arguments:
//   - type: the underlying integer type (u8, u16, u32, u64, i8, i16, i32, i64)
//   - flags: array of flag names (each flag is one bit)
//
// Example: ["bitflags", {
//   "type": "u32",
//   "flags": ["x", "y", "z", "yaw", "pitch"]
// }]
type Bitflags struct {
	name  string
	Type  *Type
	Flags []string
}

func (bf *Bitflags) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}

	// Read the underlying type
	typeData := d.Get("type")
	if typeData.Exists() {
		bf.Type = GetTypeFromJSON("", typeData)
	}

	// Read the flags array
	flagsData := d.Get("flags")
	if flagsData.IsArray() {
		for _, flag := range flagsData.Array() {
			bf.Flags = append(bf.Flags, flag.String())
		}
	}

	return nil
}

func (bf *Bitflags) SetName(name string) {
	bf.name = name
}

func (bf *Bitflags) GetName() string {
	return bf.name
}

func (bf *Bitflags) Clone() TypeExtras {
	cloned := &Bitflags{
		name:  bf.name,
		Flags: make([]string, len(bf.Flags)),
	}
	copy(cloned.Flags, bf.Flags)
	if bf.Type != nil {
		typeCopy := *bf.Type
		cloned.Type = &typeCopy
	}
	return cloned
}

func (bf *Bitflags) UpdateContainedNames(updatedNames map[string]string) {
	// Bitflags don't contain named types that need updating
}
