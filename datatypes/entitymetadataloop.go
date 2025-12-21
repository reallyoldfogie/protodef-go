package datatypes

import (
	"github.com/tidwall/gjson"
)

// EntityMetadataLoop represents a ProtoDef 'entityMetadataLoop' structure.
// This is a special loop that reads entries until a terminator value (endVal) is encountered.
// Arguments:
//   - type: the type of the elements in the loop
//   - endVal: the terminator value that signals end of loop (typically 255 for entity metadata)
//
// Example: EntityMetadata that reads entries until 0xFF
//
//	["entityMetadataLoop", {"endVal": 255, "type": "entityMetadataEntry"}]
type EntityMetadataLoop struct {
	name   string
	Type   *Type
	EndVal int
}

func (e *EntityMetadataLoop) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("type").Exists() {
		e.Type = GetTypeFromJSON("entityMetadataLoop_type", d.Get("type"))
	}
	if d.Get("endVal").Exists() {
		e.EndVal = int(d.Get("endVal").Int())
	}
	return nil
}

func (e *EntityMetadataLoop) SetName(name string) {
	e.name = name
}

func (e *EntityMetadataLoop) GetName() string {
	return e.name
}

func (e *EntityMetadataLoop) Clone() TypeExtras {
	cloned := &EntityMetadataLoop{
		name:   e.name,
		EndVal: e.EndVal,
	}
	if e.Type != nil {
		clonedType := *e.Type
		if e.Type.Extras != nil {
			clonedType.Extras = e.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}
	return cloned
}

func (e *EntityMetadataLoop) UpdateContainedNames(updatedNames map[string]string) {
	if e.Type != nil && e.Type.Extras != nil {
		e.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
