package datatypes

import (
	"github.com/tidwall/gjson"
)

type PString struct {
	name      string
	CountType *Type
	Count     any // Field reference or fixed size
	Encoding  string
}

func (p *PString) ReadJSON(d gjson.Result) error {
	if !d.IsObject() {
		return nil
	}
	if d.Get("countType").Exists() {
		p.CountType = GetTypeFromJSON("countType", d.Get("countType"))
	}
	if d.Get("count").Exists() {
		countVal := d.Get("count")
		if countVal.Type == gjson.Number {
			p.Count = int(countVal.Int())
		} else {
			p.Count = countVal.Value()
		}
	}
	if d.Get("encoding").Exists() {
		p.Encoding = d.Get("encoding").String()
	}
	return nil
}

func (p *PString) SetName(name string) {
	p.name = name
}

func (p *PString) GetName() string {
	return p.name
}

func (p *PString) Clone() TypeExtras {
	cloned := &PString{
		name:     p.name,
		Count:    p.Count,
		Encoding: p.Encoding,
	}
	if p.CountType != nil {
		clonedCountType := *p.CountType
		if p.CountType.Extras != nil {
			clonedCountType.Extras = p.CountType.Extras.Clone()
		}
		cloned.CountType = &clonedCountType
	}
	return cloned
}

func (p *PString) UpdateContainedNames(updatedNames map[string]string) {
	if p.CountType != nil && p.CountType.Extras != nil {
		p.CountType.Extras.UpdateContainedNames(updatedNames)
	}
}
