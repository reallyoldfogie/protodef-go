package datatypes

var (
	// Big-endian numeric types (default)
	Int8   = &Type{Name: "i8", TypeName: "i8"}
	UInt8  = &Type{Name: "u8", TypeName: "u8"}
	Int16  = &Type{Name: "i16", TypeName: "i16"}
	UInt16 = &Type{Name: "u16", TypeName: "u16"}
	Int32  = &Type{Name: "i32", TypeName: "i32"}
	UInt32 = &Type{Name: "u32", TypeName: "u32"}
	fInt32 = &Type{Name: "f32", TypeName: "f32"}
	fInt64 = &Type{Name: "f64", TypeName: "f64"}
	Int64  = &Type{Name: "i64", TypeName: "i64"}
	UInt64 = &Type{Name: "u64", TypeName: "u64"}

	// Little-endian numeric types (prefixed with 'l')
	LInt8   = &Type{Name: "li8", TypeName: "li8"}
	LUInt8  = &Type{Name: "lu8", TypeName: "lu8"}
	LInt16  = &Type{Name: "li16", TypeName: "li16"}
	LUInt16 = &Type{Name: "lu16", TypeName: "lu16"}
	LInt32  = &Type{Name: "li32", TypeName: "li32"}
	LUInt32 = &Type{Name: "lu32", TypeName: "lu32"}
	LfInt32 = &Type{Name: "lf32", TypeName: "lf32"}
	LfInt64 = &Type{Name: "lf64", TypeName: "lf64"}
	LInt64  = &Type{Name: "li64", TypeName: "li64"}
	LUInt64 = &Type{Name: "lu64", TypeName: "lu64"}

	// Variable-length integers
	VarInt = &Type{Name: "varint", TypeName: "varint"}

	// Non-standard extensions
	VarInt64  = &Type{Name: "varint64", TypeName: "varint64"}
	VarInt128 = &Type{Name: "varint128", TypeName: "varint128"}
	ZigZag32  = &Type{Name: "zigzag32", TypeName: "zigzag32"}
	ZigZag64  = &Type{Name: "zigzag64", TypeName: "zigzag64"}

	// Generic integer types with configurable size
	GenericInt  = func(size int) *Type { return &Type{Name: "int", TypeName: "int", Extras: size} }
	GenericLInt = func(size int) *Type { return &Type{Name: "lint", TypeName: "lint", Extras: size} }
)
