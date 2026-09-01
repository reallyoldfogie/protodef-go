package datatypes

func GetNativeType(name string) *Type {
	switch name {
	// Big-endian numeric types
	case "i8":
		return Int8
	case "u8":
		return UInt8
	case "i16":
		return Int16
	case "u16":
		return UInt16
	case "i32":
		return Int32
	case "u32":
		return UInt32
	case "f32":
		return fInt32
	case "f64":
		return fInt64
	case "i64":
		return Int64
	case "u64":
		return UInt64

	// Little-endian numeric types
	case "li8":
		return LInt8
	case "lu8":
		return LUInt8
	case "li16":
		return LInt16
	case "lu16":
		return LUInt16
	case "li32":
		return LInt32
	case "lu32":
		return LUInt32
	case "lf32":
		return LfInt32
	case "lf64":
		return LfInt64
	case "li64":
		return LInt64
	case "lu64":
		return LUInt64

	// Variable-length integers
	case "varint":
		return VarInt

	// Non-standard extensions
	case "varint64":
		return VarInt64
	case "varint128":
		return VarInt128
	case "zigzag32":
		return ZigZag32
	case "zigzag64":
		return ZigZag64

	// Primitives
	case "bool":
		return Bool
	case "cstring":
		return CString
	case "void":
		return Void
	}

	return nil
}
