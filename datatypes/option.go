package datatypes

import (
	"github.com/tidwall/gjson"
)

type Option struct {
	name string
	Type *Type
}

func (o *Option) ReadJSON(d gjson.Result) error {
	// Option type expects a single type as argument
	o.Type = GetTypeFromJSON("option", d)
	return nil
}

func (o *Option) SetName(name string) {
	o.name = name
}

func (o *Option) GetName() string {
	return o.name
}

func (o *Option) Clone() TypeExtras {
	cloned := &Option{
		name: o.name,
	}
	if o.Type != nil {
		clonedType := *o.Type
		if o.Type.Extras != nil {
			clonedType.Extras = o.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}
	return cloned
}

func (o *Option) UpdateContainedNames(updatedNames map[string]string) {
	if o.Type != nil && o.Type.Extras != nil {
		o.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
