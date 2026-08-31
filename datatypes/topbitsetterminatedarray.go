package datatypes

import (
	"github.com/tidwall/gjson"
)

// TopBitSetTerminatedArray represents a ProtoDef 'topBitSetTerminatedArray' or 'topbitsetalternative' structure.
// This is an array-like structure where entries are indexed by 7-bit slot indices,
// terminated by a byte with MSB=1 (high bit set).
//
// Arguments:
//   - type: the type of elements in the array (typically a container with slot index and data)
//
// Example: Equipment slots that are indexed and terminated by MSB:
//
//	["topBitSetTerminatedArray", {"type": ["container", [...]]}]
//
// The nested type defines the structure of each entry (e.g., slot index + item data).
type TopBitSetTerminatedArray struct {
	name string
	Type *Type
}

func (t *TopBitSetTerminatedArray) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("type").Exists() {
		t.Type = GetTypeFromJSON("topBitSetTerminatedArray_type", d.Get("type"))
	}
	return nil
}

func (t *TopBitSetTerminatedArray) SetName(name string) {
	t.name = name
}

func (t *TopBitSetTerminatedArray) GetName() string {
	return t.name
}

func (t *TopBitSetTerminatedArray) Clone() TypeExtras {
	cloned := &TopBitSetTerminatedArray{
		name: t.name,
	}
	if t.Type != nil {
		clonedType := *t.Type
		if t.Type.Extras != nil {
			clonedType.Extras = t.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}
	return cloned
}

func (t *TopBitSetTerminatedArray) UpdateContainedNames(updatedNames map[string]string) {
	if t.Type != nil && t.Type.Extras != nil {
		t.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
