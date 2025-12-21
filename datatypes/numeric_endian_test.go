package datatypes

import (
	"testing"
)

func TestLittleEndianTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		{"li8", "li8", LInt8},
		{"lu8", "lu8", LUInt8},
		{"li16", "li16", LInt16},
		{"lu16", "lu16", LUInt16},
		{"li32", "li32", LInt32},
		{"lu32", "lu32", LUInt32},
		{"lf32", "lf32", LfInt32},
		{"lf64", "lf64", LfInt64},
		{"li64", "li64", LInt64},
		{"lu64", "lu64", LUInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNativeType(tt.typeName)
			if got == nil {
				t.Errorf("GetNativeType(%q) returned nil", tt.typeName)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetNativeType(%q).Name = %q, want %q", tt.typeName, got.Name, tt.want.Name)
			}
			if got.TypeName != tt.want.TypeName {
				t.Errorf("GetNativeType(%q).TypeName = %q, want %q", tt.typeName, got.TypeName, tt.want.TypeName)
			}
		})
	}
}

func TestBigEndianTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		{"i8", "i8", Int8},
		{"u8", "u8", UInt8},
		{"i16", "i16", Int16},
		{"u16", "u16", UInt16},
		{"i32", "i32", Int32},
		{"u32", "u32", UInt32},
		{"f32", "f32", fInt32},
		{"f64", "f64", fInt64},
		{"i64", "i64", Int64},
		{"u64", "u64", UInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNativeType(tt.typeName)
			if got == nil {
				t.Errorf("GetNativeType(%q) returned nil", tt.typeName)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetNativeType(%q).Name = %q, want %q", tt.typeName, got.Name, tt.want.Name)
			}
			if got.TypeName != tt.want.TypeName {
				t.Errorf("GetNativeType(%q).TypeName = %q, want %q", tt.typeName, got.TypeName, tt.want.TypeName)
			}
		})
	}
}

func TestVarintTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		{"varint", "varint", VarInt},
		{"varint64", "varint64", VarInt64},
		{"varint128", "varint128", VarInt128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNativeType(tt.typeName)
			if got == nil {
				t.Errorf("GetNativeType(%q) returned nil", tt.typeName)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetNativeType(%q).Name = %q, want %q", tt.typeName, got.Name, tt.want.Name)
			}
		})
	}
}

func TestZigZagTypes(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		{"zigzag32", "zigzag32", ZigZag32},
		{"zigzag64", "zigzag64", ZigZag64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNativeType(tt.typeName)
			if got == nil {
				t.Errorf("GetNativeType(%q) returned nil", tt.typeName)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetNativeType(%q).Name = %q, want %q", tt.typeName, got.Name, tt.want.Name)
			}
		})
	}
}

func TestGenericIntFunction(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 1", 1},
		{"size 2", 2},
		{"size 3", 3},
		{"size 4", 4},
		{"size 8", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenericInt(tt.size)
			if got == nil {
				t.Errorf("GenericInt(%d) returned nil", tt.size)
				return
			}
			if got.Name != "int" {
				t.Errorf("GenericInt(%d).Name = %q, want %q", tt.size, got.Name, "int")
			}
			if got.TypeName != "int" {
				t.Errorf("GenericInt(%d).TypeName = %q, want %q", tt.size, got.TypeName, "int")
			}
			if gotSize, ok := got.Extras.(int); !ok || gotSize != tt.size {
				t.Errorf("GenericInt(%d).Extras = %v, want %d", tt.size, got.Extras, tt.size)
			}
		})
	}
}

func TestGenericLIntFunction(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"size 1", 1},
		{"size 2", 2},
		{"size 3", 3},
		{"size 5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenericLInt(tt.size)
			if got == nil {
				t.Errorf("GenericLInt(%d) returned nil", tt.size)
				return
			}
			if got.Name != "lint" {
				t.Errorf("GenericLInt(%d).Name = %q, want %q", tt.size, got.Name, "lint")
			}
			if got.TypeName != "lint" {
				t.Errorf("GenericLInt(%d).TypeName = %q, want %q", tt.size, got.TypeName, "lint")
			}
			if gotSize, ok := got.Extras.(int); !ok || gotSize != tt.size {
				t.Errorf("GenericLInt(%d).Extras = %v, want %d", tt.size, got.Extras, tt.size)
			}
		})
	}
}

func TestGetNativeTypePrimitives(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		want     *Type
	}{
		{"bool", "bool", Bool},
		{"cstring", "cstring", CString},
		{"void", "void", Void},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetNativeType(tt.typeName)
			if got == nil {
				t.Errorf("GetNativeType(%q) returned nil", tt.typeName)
				return
			}
			if got.Name != tt.want.Name {
				t.Errorf("GetNativeType(%q).Name = %q, want %q", tt.typeName, got.Name, tt.want.Name)
			}
		})
	}
}

func TestGetNativeTypeUnknown(t *testing.T) {
	got := GetNativeType("unknown_type")
	if got != nil {
		t.Errorf("GetNativeType(\"unknown_type\") = %v, want nil", got)
	}
}
