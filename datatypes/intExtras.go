package datatypes

import (
	"github.com/tidwall/gjson"
)

// IntExtras represents a generic integer type with configurable size (int/lint)
type IntExtras struct {
	name string
	Size int
}

func (ie *IntExtras) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("size").Exists() {
		ie.Size = int(d.Get("size").Int())
	}
	return nil
}

func (ie *IntExtras) SetName(name string) {
	ie.name = name
}

func (ie *IntExtras) GetName() string {
	return ie.name
}

func (ie *IntExtras) Clone() TypeExtras {
	return &IntExtras{
		name: ie.name,
		Size: ie.Size,
	}
}

func (ie *IntExtras) UpdateContainedNames(updatedNames map[string]string) {
	// IntExtras doesn't contain any field references
}
