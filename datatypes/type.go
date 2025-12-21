package datatypes

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"

	"github.com/tidwall/gjson"
)

func getGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := bytes.Fields(buf[:n])[1]
	id, err := strconv.ParseInt(string(idField), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

type Type struct {
	Name     string
	TypeName string
	Comment  string

	Extras TypeExtras
	
	// RawDefinition contains the original protodef JSON snippet that defined this type.
	// This is useful for debugging to see what the original definition was.
	RawDefinition string
}

func GetTypeFromJSON(name string, option gjson.Result) *Type {
	if option.Type == gjson.String && option.String() == "native" {
		return GetNativeType(name)
	}

	t := GetType(option)

	fmt.Printf("DEBUG [GetTypeFromJSON]: got type %#v for %#v\n", t, option)

	if t != nil {
		// Always create a copy to avoid mutating shared type objects
		// (especially important for native types which are global singletons)
		if t.Name != name {
			tCopy := &Type{
				Name:          name,
				TypeName:      t.TypeName,
				Comment:       t.Comment,
				Extras:        t.Extras,
				RawDefinition: option.Raw,
			}
			t = tCopy
		} else {
			// Store raw definition even when name matches
			t.RawDefinition = option.Raw
		}
		// DEBUG: Log successful type creation for simple string aliases
		if option.Type == gjson.String && t.Name != t.TypeName {
			fmt.Printf("DEBUG [GetTypeFromJSON %d]: Created simple type alias: '%s' -> '%s'\n", getGoroutineID(), t.Name, t.TypeName)
			fmt.Printf("DEBUG [GetTypeFromJSON %d]: Created simple type alias: %#v (%p)\n", getGoroutineID(), t, t)
			if t.Name == "ContainerID" {
				var buf [1500]byte
				byteSize := runtime.Stack(buf[:], false)
				fmt.Printf("%s\n\n", string(buf[:byteSize]))
			}
		}
		return t
	}

	// DEBUG: Log when GetTypeFromJSON returns nil
	fmt.Printf("DEBUG [GetTypeFromJSON]: Returning nil for name='%s', option.Type=%v, option.Raw='%s'\n",
		name, option.Type, option.Raw)
	return nil
}

func GetType(d gjson.Result) *Type {
	fmt.Printf("DEBUG [GetType] %#v\n", d)
	var t *Type
	if d.Type == gjson.String {
		t = GetNativeType(d.String())
		if t != nil {
			fmt.Printf("DEBUG [GetType] returning native simple type %s -> %#v\n", d.String(), d)
			return t
		}

		fmt.Printf("DEBUG [GetType] returning non-native simple type %s -> %#v\n", d.String(), d)

		return &Type{
			Name:          d.String(),
			TypeName:      d.String(),
			RawDefinition: d.Raw,
		}
	}

	if d.IsArray() {
		t = &Type{
			RawDefinition: d.Raw,
		}
		arr := d.Array()
		arr_len := len(arr)
		if arr_len == 2 {
			arr_type := arr[0]
			if arr_type.Type == gjson.String {
				t.TypeName = arr_type.String()
				switch t.TypeName {
				case "container":
					t.Extras = &Container{}
					// t.Extras.(*Container).ReadJSON(arr[1])
				case "switch":
					t.Extras = &Switch{}
					// t.Extras.(*Switch).ReadJSON(arr[1])
				case "option":
					t.Extras = &Option{}
					// t.Extras.(*Option).ReadJSON(arr[1])
				case "array":
					t.Extras = &Array{}
					// t.Extras.(*Array).ReadJSON(arr[1])
				case "buffer":
					t.Extras = &Buffer{}
					// t.Extras.(*Buffer).ReadJSON(arr[1])
				case "bitfield":
					t.Extras = &Bitfield{}
					// t.Extras.(*Bitfield).ReadJSON(arr[1])
				case "bitflags":
					t.Extras = &Bitflags{}
					// t.Extras.(*Bitflags).ReadJSON(arr[1])
				case "pstring":
					t.Extras = &PString{}
					// t.Extras.(*PString).ReadJSON(arr[1])
				case "int":
					// Generic integer with configurable size: ["int", {"size": N}]
					t.Extras = &IntExtras{}
					// t.Extras.(*IntExtras).ReadJSON(arr[1])
				case "count":
					t.Extras = &Count{}
					// t.Extras.(*Count).ReadJSON(arr[1])
				case "mapper":
					t.Extras = &Mapper{}
					// t.Extras.(*Mapper).ReadJSON(arr[1])
				case "lint":
					// Little-endian integer with configurable size: ["lint", {"size": N}]
					t.Extras = &IntExtras{}
					// t.Extras.ReadJSON(arr[1])
				case "registryEntryHolder":
					t.Extras = &RegistryEntryHolder{}
					// t.Extras.ReadJSON(arr[1])
				case "registryEntryHolderSet":
					t.Extras = &RegistryEntryHolderSet{}
					// t.Extras.ReadJSON(arr[1])
				case "entityMetadataLoop":
					t.Extras = &EntityMetadataLoop{}
					// t.Extras.ReadJSON(arr[1])
				}
			}
			if t.Extras != nil {
				t.Extras.ReadJSON(arr[1])
			}
		}
		return t
	}
	return nil
}
