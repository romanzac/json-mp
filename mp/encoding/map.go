package encoding

import (
	"math"
	"reflect"

	"github.com/romanzac/json-mp/mp/def"
)

func (e *encoder) computeFixMap(rv reflect.Value) (int, bool) {
	size := 0

	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeInt(int64(v))
		}
		return size, true

	case map[string]uint:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeUint(uint64(v))
		}
		return size, true

	case map[string]string:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeString(v)
		}
		return size, true

	case map[string]float32:
		for k := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeFloat32()
		}
		return size, true

	case map[string]float64:
		for k := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeFloat64()
		}
		return size, true

	case map[string]bool:
		for k := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 /*+ e.computeBool()*/
		}
		return size, true

	case map[string]int8:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeInt(int64(v))
		}
		return size, true
	case map[string]int16:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeInt(int64(v))
		}
		return size, true
	case map[string]int32:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeInt(int64(v))
		}
		return size, true
	case map[string]int64:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeInt(v)
		}
		return size, true
	case map[string]uint8:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeUint(uint64(v))
		}
		return size, true
	case map[string]uint16:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeUint(uint64(v))
		}
		return size, true
	case map[string]uint32:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeUint(uint64(v))
		}
		return size, true
	case map[string]uint64:
		for k, v := range m {
			size += def.Byte1 + e.computeString(k)
			size += def.Byte1 + e.computeUint(v)
		}
		return size, true
	}
	return size, false
}

func (e *encoder) writeMapLength(l int, offset int) int {
	if l <= 0x0f {
		offset = e.setByte1Int(def.FixMap+l, offset)
	} else if l <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(l, offset)
	} else if uint(l) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(l, offset)
	}
	return offset
}

func (e *encoder) writeFixMap(rv reflect.Value, offset int) (int, bool) {
	switch m := rv.Interface().(type) {
	case map[string]int:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true

	case map[string]float32:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeFloat32(float64(v), offset)
		}
		return offset, true

	case map[string]float64:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeFloat64(v, offset)
		}
		return offset, true

	case map[string]bool:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeBool(v, offset)
		}
		return offset, true

	case map[string]string:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeString(v, offset)
		}
		return offset, true

	case map[string]int8:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true
	case map[string]int16:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true
	case map[string]int32:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true
	case map[string]int64:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeInt(int64(v), offset)
		}
		return offset, true

	case map[string]uint8:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint16:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint32:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true
	case map[string]uint64:
		for k, v := range m {
			offset = e.writeString(k, offset)
			offset = e.writeUint(uint64(v), offset)
		}
		return offset, true
	}
	return offset, false
}
