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
			
			// DEBUG: Log parsed field
			typeInfo := "nil"
			if field.Type != nil {
				typeInfo = fmt.Sprintf("%s (typename=%s)", field.Type.Name, field.Type.TypeName)
			}
			DebugPrintf("DEBUG [container.ReadJSON]: Container '%s' parsed field '%s' Type=%s Anon=%v\n",
				cf.Name, field.Name, typeInfo, field.Anon)
		}

	cf.Fields = append(cf.Fields, field)
	}
	return nil
}

func (cf *Container) SetName(name string) {
	cf.Name = name
}

func (cf *Container) GetName() string {
	return cf.Name
}

func (cf *Container) Clone() TypeExtras {
	cloned := &Container{
		Name:   cf.Name,
		Fields: make([]*ContainerField, len(cf.Fields)),
	}
	for i, field := range cf.Fields {
		clonedField := &ContainerField{
			Name: field.Name,
			Anon: field.Anon,
		}
		if field.Type != nil {
			clonedType := *field.Type
			if field.Type.Extras != nil {
				clonedType.Extras = field.Type.Extras.Clone()
			}
			// RawDefinition is copied automatically by the struct copy above
			clonedField.Type = &clonedType
		}
		cloned.Fields[i] = clonedField
	}
	return cloned
}

func (cf *Container) UpdateContainedNames(updatedNames map[string]string) {
	for _, field := range cf.Fields {
		if newName, exists := updatedNames[field.Name]; exists {
			field.Name = newName
		}
		if field.Type != nil && field.Type.Extras != nil {
			field.Type.Extras.UpdateContainedNames(updatedNames)
		}
	}
}
