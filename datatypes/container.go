package datatypes

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// ContainerField represents a field in a ProtoDef 'container' structure.
type ContainerField struct {
	Name string
	Type *Type
	Anon bool // Anonymous field (also known as unnamed field)
}

type Container struct {
	Name   string
	Fields []*ContainerField
}

func (cf *Container) ReadJSON(d gjson.Result) error {
	if !d.IsArray() {
		// Container represents a ProtoDef 'container' structure.
		return fmt.Errorf("container %s data is not array", cf.Name)
	}

	for _, value := range d.Array() {
		field := &ContainerField{}

		// Parse anon field (anonymous/unnamed field)
		if value.Get("anon").Exists() {
			field.Anon = value.Get("anon").Bool()
		}

		// Parse name field (regular named field)
		if value.Get("name").Exists() {
			field.Name = value.Get("name").String()
		}

		// Parse type (required)
		if value.Get("type").Exists() {
			typeData := value.Get("type")
			// Use name for type lookup if available, otherwise use generic name
			typeName := field.Name
			if typeName == "" {
				typeName = "anon_field"
			}
			field.Type = GetTypeFromJSON(typeName, typeData)
		}

		cf.Fields = append(cf.Fields, field)
	}
	return nil
}
