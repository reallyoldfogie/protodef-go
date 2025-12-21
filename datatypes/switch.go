package datatypes

import (
	"github.com/tidwall/gjson"
)

type SwitchType struct {
	CompareTo      string
	CompareToValue any
	Fields         map[string]*Type
	Default        *Type
}

func (s *SwitchType) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	s.CompareTo = d.Get("compareTo").String()
	if d.Get("compareToValue").Exists() {
		s.CompareToValue = d.Get("compareToValue").Value()
	}
	s.Fields = make(map[string]*Type)
	fields := d.Get("fields")
	if fields.IsObject() {
		fields.ForEach(func(key, value gjson.Result) bool {
			s.Fields[key.String()] = GetTypeFromJSON(key.String(), value)
			return true
		})
	}
	if d.Get("default").Exists() {
		defVal := d.Get("default")
		if defVal.Type == gjson.String {
			s.Default = GetTypeFromJSON(defVal.String(), defVal)
		}
	}
	return nil
}
