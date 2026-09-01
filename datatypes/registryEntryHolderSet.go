package datatypes

import "github.com/tidwall/gjson"

// ["registryEntryHolderSet",
//     {
//       "base": {
//         "name": "name",
//         "type": "string"
//       },
//       "otherwise": {
//         "name": "ids",
//         "type": "varint"
//       }
//     }
//   ]
//

type RegistryEntryHolderSet struct {
	name      string
	Type      *Type
	Base      RegistryEntryBase
	Otherwise RegistryEntryOthewise
}

type RegistryEntryBase struct {
	Name string
	Type *Type
}

func (m *RegistryEntryHolderSet) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}

	if d.Get("base").Exists() {
		baseData := d.Get("base")
		if baseData.Get("name").Exists() {
			m.Base.Name = baseData.Get("name").String()
		}
		if baseData.Get("type").Exists() {
			m.Base.Type = GetTypeFromJSON("registryEntryHolderSet_base_type", baseData.Get("type"))
		}
	}

	if d.Get("otherwise").Exists() {
		otherwiseData := d.Get("otherwise")
		if otherwiseData.Get("type").Exists() {
			m.Otherwise.Name = otherwiseData.Get("name").String()
			m.Otherwise.Type = GetTypeFromJSON("registryEntryHolderSet_otherwise_type", otherwiseData.Get("type"))
		}
	}
	return nil
}

func (m *RegistryEntryHolderSet) SetName(name string) {
	m.name = name
}

func (m *RegistryEntryHolderSet) GetName() string {
	return m.name
}

func (m *RegistryEntryHolderSet) Clone() TypeExtras {
	cloned := &RegistryEntryHolderSet{
		name: m.name,
		Base: RegistryEntryBase{
			Name: m.Base.Name,
		},
		Otherwise: RegistryEntryOthewise{
			Name: m.Otherwise.Name,
		},
	}

	if m.Type != nil {
		clonedType := *m.Type
		if m.Type.Extras != nil {
			clonedType.Extras = m.Type.Extras.Clone()
		}
		cloned.Type = &clonedType
	}

	if m.Base.Type != nil {
		clonedType := *m.Base.Type
		if m.Base.Type.Extras != nil {
			clonedType.Extras = m.Base.Type.Extras.Clone()
		}
		cloned.Base.Type = &clonedType
	}

	if m.Otherwise.Type != nil {
		clonedType := *m.Otherwise.Type
		if m.Otherwise.Type.Extras != nil {
			clonedType.Extras = m.Otherwise.Type.Extras.Clone()
		}
		cloned.Otherwise.Type = &clonedType
	}
	return cloned
}

func (m *RegistryEntryHolderSet) UpdateContainedNames(updatedNames map[string]string) {
	if m.Type != nil && m.Type.Extras != nil {
		m.Type.Extras.UpdateContainedNames(updatedNames)
	}
	if m.Base.Type != nil && m.Base.Type.Extras != nil {
		m.Base.Type.Extras.UpdateContainedNames(updatedNames)
	}
	if m.Otherwise.Type != nil && m.Otherwise.Type.Extras != nil {
		m.Otherwise.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
