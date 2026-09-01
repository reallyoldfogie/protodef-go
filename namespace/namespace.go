package namespace

import (
	"github.com/protodef-go/protodef-go/datatypes"
	"github.com/tidwall/gjson"
)

type Namespace struct {
	Name       string
	Types      []*datatypes.Type
	Namespaces map[string]*Namespace
}

func (ns *Namespace) ReadJSON(d gjson.Result) error {
	types := d.Get("types")

	if types.Exists() {
		if types.IsObject() {
			for name, option := range types.Map() {
				t := datatypes.GetTypeFromJSON(name, option)
				ns.Types = append(ns.Types, t)
			}
		}
		if types.IsArray() {
			for _, name := range types.Array() {
				ns.Types = append(ns.Types, datatypes.GetType(name))
			}
		}
	}

	ns.Namespaces = make(map[string]*Namespace)
	for name, data := range d.Map() {
		if name != "types" {
			namespace := &Namespace{Name: name}
			err := namespace.ReadJSON(data)
			if err != nil {
				return err
			}
			ns.Namespaces[name] = namespace
		}
	}
	return nil
}
