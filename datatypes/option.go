package datatypes

import (
	"github.com/tidwall/gjson"
)

type OptionType struct {
	Type *Type
}

func (o *OptionType) ReadJSON(d gjson.Result) error {
	// Option type expects a single type as argument
	o.Type = GetTypeFromJSON("option", d)
	return nil
}
