package datatypes

import "github.com/tidwall/gjson"

// ["registryEntryHolder",
//     {
//       "baseName": "patternId",
//       "otherwise": {
//         "name": "data",
//         "type": "BannerPattern"
//       }
//     }
//   ]
//

type RegistryEntryHolder struct {
	name      string
	Type      *Type
	BaseName  string
	Otherwise RegistryEntryOthewise
}

type RegistryEntryOthewise struct {
	Name string
	Type *Type
}

func (m *RegistryEntryHolder) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}

	if d.Get("basename").Exists() {
		m.BaseName = d.Get("basename").String()
	}

	if d.Get("otherwise").Exists() {
		otherwiseData := d.Get("otherwise")
		if otherwiseData.Get("type").Exists() {
			m.Otherwise.Name = otherwiseData.Get("name").String()
			m.Otherwise.Type = GetTypeFromJSON("registryEntryHolder_type", otherwiseData.Get("type"))
		}
	}
	return nil
}

func (m *RegistryEntryHolder) SetName(name string) {
	m.name = name
}

func (m *RegistryEntryHolder) GetName() string {
	return m.name
}

func (m *RegistryEntryHolder) Clone() TypeExtras {
	cloned := &RegistryEntryHolder{
		name: m.name,
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

	if m.Otherwise.Type != nil {
		clonedType := *m.Otherwise.Type
		if m.Otherwise.Type.Extras != nil {
			clonedType.Extras = m.Otherwise.Type.Extras.Clone()
		}
		cloned.Otherwise.Type = &clonedType
	}
	return cloned
}

func (m *RegistryEntryHolder) UpdateContainedNames(updatedNames map[string]string) {
	if m.Type != nil && m.Type.Extras != nil {
		m.Type.Extras.UpdateContainedNames(updatedNames)
	}
}
