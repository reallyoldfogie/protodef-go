package datatypes

import (
	"github.com/tidwall/gjson"
)

type Type struct {
	Name     string
	TypeName string

	Extras any
}

func GetTypeFromJSON(name string, option gjson.Result) *Type {
	if option.Type == gjson.String && option.String() == "native" {
		return GetNativeType(name)
	}
	t := GetType(name, option)
	if t != nil {
		t.Name = name
		return t
	}
	return nil
}

func GetType(name string, d gjson.Result) *Type {
	var t *Type
	if d.Type == gjson.String {
		t = GetNativeType(d.String())
		if t != nil {
			return t
		}
	}

	if d.IsArray() {
		t = &Type{}
		arr := d.Array()
		arr_len := len(arr)
		if arr_len == 2 {
			arr_type := arr[0]
			if arr_type.Type == gjson.String {
				t.TypeName = arr_type.String()
				switch t.TypeName {
				case "container":
					t.Extras = &Container{}
					t.Extras.(*Container).ReadJSON(arr[1])
				case "switch":
					t.Extras = &SwitchType{}
					t.Extras.(*SwitchType).ReadJSON(arr[1])
				case "option":
					t.Extras = &OptionType{}
					t.Extras.(*OptionType).ReadJSON(arr[1])
				case "array":
					t.Extras = &ArrayType{}
					t.Extras.(*ArrayType).ReadJSON(arr[1])
				case "buffer":
					t.Extras = &BufferType{}
					t.Extras.(*BufferType).ReadJSON(arr[1])
				case "bitfield":
					t.Extras = &BitField{}
					t.Extras.(*BitField).ReadJSON(arr[1])
				case "pstring":
					t.Extras = &PStringType{}
					t.Extras.(*PStringType).ReadJSON(arr[1])
				case "int":
					// Generic integer with configurable size: ["int", {"size": N}]
					if arr[1].IsObject() && arr[1].Get("size").Exists() {
						size := int(arr[1].Get("size").Int())
						t.Extras = size
					}
				case "count":
					t.Extras = &CountType{}
					t.Extras.(*CountType).ReadJSON(arr[1])
				case "mapper":
					t.Extras = &MapperType{}
					t.Extras.(*MapperType).ReadJSON(arr[1])
				case "lint":
					// Little-endian integer with configurable size: ["lint", {"size": N}]
					if arr[1].IsObject() && arr[1].Get("size").Exists() {
						size := int(arr[1].Get("size").Int())
						t.Extras = size
					}
				}
			}
		}
		return t
	}
	return nil
}
