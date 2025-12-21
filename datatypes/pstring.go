package datatypes

import (
	"github.com/tidwall/gjson"
)

type PStringType struct {
	CountType *Type
	Count     any // Field reference or fixed size
	Encoding  string
}

func (p *PStringType) ReadJSON(d gjson.Result) error {
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
