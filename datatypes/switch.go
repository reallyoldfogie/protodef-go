package datatypes

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type Switch struct {
	name           string
	CompareTo      string
	CompareToValue any
	Fields         map[string]*Type
	Default        *Type
}

func (s *Switch) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	s.CompareTo = d.Get("compareTo").String()
	if d.Get("compareToValue").Exists() {
		s.CompareToValue = d.Get("compareToValue").Value()
	}
	
	// DEBUG: Log switch compareTo field
	fmt.Printf("DEBUG [switch.ReadJSON]: Switch '%s' has compareTo='%s'\n", s.name, s.CompareTo)
	
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
		// Handle both string defaults (e.g., "void") and complex defaults (e.g., ["container", [...]])
		if defVal.Type == gjson.String {
			s.Default = GetTypeFromJSON(defVal.String(), defVal)
		} else {
			// For array/object defaults, pass to GetTypeFromJSON with generic name
			s.Default = GetTypeFromJSON("default", defVal)
		}
	}
	return nil
}

func (s *Switch) SetName(name string) {
	s.name = name
}

func (s *Switch) GetName() string {
	return s.name
}

func (s *Switch) Clone() TypeExtras {
	cloned := &Switch{
		name:           s.name,
		CompareTo:      s.CompareTo,
		CompareToValue: s.CompareToValue,
		Fields:         make(map[string]*Type, len(s.Fields)),
	}
	for key, t := range s.Fields {
		if t != nil {
			clonedType := *t
			if t.Extras != nil {
				clonedType.Extras = t.Extras.Clone()
			}
			cloned.Fields[key] = &clonedType
		}
	}
	if s.Default != nil {
		clonedDefault := *s.Default
		if s.Default.Extras != nil {
			clonedDefault.Extras = s.Default.Extras.Clone()
		}
		cloned.Default = &clonedDefault
	}
	return cloned
}

func (s *Switch) UpdateContainedNames(updatedNames map[string]string) {
	for _, t := range s.Fields {
		if t != nil && t.Extras != nil {
			t.Extras.UpdateContainedNames(updatedNames)
		}
	}
	if s.Default != nil && s.Default.Extras != nil {
		s.Default.Extras.UpdateContainedNames(updatedNames)
	}
}
