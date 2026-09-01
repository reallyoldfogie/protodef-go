package datatypes

var (
	Bool    = &Type{Name: "bool", TypeName: "bool"}
	CString = &Type{Name: "cstring", TypeName: "cstring"}
	Void    = &Type{Name: "void", TypeName: "void"}
)

type BoolType struct{}
type CStringType struct{}
type VoidType struct{}

// TODO: Implement encoding/decoding logic for each primitive type
