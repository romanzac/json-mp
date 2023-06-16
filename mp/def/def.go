package def

import "reflect"

// Message pack format without extensions
const (
	Nil = 0xc0

	False = 0xc2
	True  = 0xc3

	FixIntMin = 0x00
	FixIntMax = 0x7f

	Uint8  = 0xcc
	Uint16 = 0xcd
	Uint32 = 0xce
	Uint64 = 0xcf

	Int8  = 0xd0
	Int16 = 0xd1
	Int32 = 0xd2
	Int64 = 0xd3

	Float32 = 0xca
	Float64 = 0xcb

	FixStr = 0xa0
	Str8   = 0xd9
	Str16  = 0xda
	Str32  = 0xdb

	FixMap = 0x80
	Map16  = 0xde
	Map32  = 0xdf

	FixArray = 0x90
	Array16  = 0xdc
	Array32  = 0xdd

	NegativeFixIntMin = 0xe0 - 0xff // -31
	NegativeFixIntMax = -0x01       //  -1
)

// Bytes
const (
	Byte1 = 1
	Byte2 = 2
	Byte4 = 4
	Byte8 = 8
)

// CheckStructField returns flag when to encode/decode or not and a field name
func CheckStructField(field reflect.StructField) (bool, string) {
	// Is the name with the first capital letter (public field)
	if 0x41 <= field.Name[0] && field.Name[0] <= 0x5a {
		if tag := field.Tag.Get("json"); tag == "-" {
			return false, ""
		} else if len(tag) > 0 {
			return true, tag
		}
		return true, field.Name
	}
	return false, ""
}
