package datatypes

import (
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

type TypeExtras interface {
	ReadJSON(d gjson.Result) error
	SetName(name string)
	GetName() string
	Clone() TypeExtras
	UpdateContainedNames(updatedNames map[string]string)
}

func ReplaceAllSymbol(in, name, newName string) string {
	if strings.Contains(in, newName) {
		return in
	}
	var re = regexp.MustCompile(`(?m)([\w\W])?(` + name + `)([\w\W])?$`)
	var substitution = "$1 " + newName + " $3"
	return re.ReplaceAllString(in, substitution)
}
