package encoding

import (
	"fmt"
	"math"
	"reflect"

	"github.com/romanzac/json-mp/mp/def"
)

type encoder struct {
	d  []byte
	mk map[uintptr][]reflect.Value
	mv map[uintptr][]reflect.Value
}

func Encode(v interface{}) ([]byte, error) {
	e := encoder{}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem()
		}
	}
	size, err := e.computeSize(rv)
	if err != nil {
		return nil, err
	}
	e.d = make([]byte, size)
	last := e.add(rv, 0)
	if size != last {
		return nil, fmt.Errorf("failed serialization size=%d, lastIdx=%d", size, last)
	}

	return e.d, err
}

func (e *encoder) computeSize(rv reflect.Value) (int, error) {
	ret := def.Byte1

	switch rv.Kind() {
	case reflect.Bool:
		// Single byte size - do nothing

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		ret += e.computeUint(v)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		ret += e.computeInt(int64(v))

	case reflect.Float32:
		ret += e.computeFloat32()

	case reflect.Float64:
		ret += e.computeFloat64()

	case reflect.String:
		ret += e.computeString(rv.String())

	case reflect.Slice:
		if rv.IsNil() {
			return ret, nil
		}
		l := rv.Len()

		// Check format size
		if l <= 0x0f {
			// Do nothing - format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}

		if size, find := e.computeFixSlice(rv); find {
			ret += size
			return ret, nil
		}

	case reflect.Array:
		l := rv.Len()

		// Check format size
		if l <= 0x0f {
			// Do nothing - format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this array length : %d", l)
		}

	case reflect.Map:
		if rv.IsNil() {
			return ret, nil
		}

		l := rv.Len()
		// Check format size
		if l <= 0x0f {
			// Do nothing - format code only
		} else if l <= math.MaxUint16 {
			ret += def.Byte2
		} else if uint(l) <= math.MaxUint32 {
			ret += def.Byte4
		} else {
			// not supported error
			return 0, fmt.Errorf("not support this map length : %d", l)
		}

		if size, find := e.computeFixMap(rv); find {
			ret += size
			return ret, nil
		}

		if e.mk == nil {
			e.mk = map[uintptr][]reflect.Value{}
			e.mv = map[uintptr][]reflect.Value{}
		}

		// Fill in keys and values
		keys := rv.MapKeys()
		mv := make([]reflect.Value, len(keys))
		i := 0
		for _, k := range keys {
			keySize, err := e.computeSize(k)
			if err != nil {
				return 0, err
			}
			value := rv.MapIndex(k)
			valueSize, err := e.computeSize(value)
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
			mv[i] = value
			i++
		}
		e.mk[rv.Pointer()], e.mv[rv.Pointer()] = keys, mv

	case reflect.Struct:
		size, err := e.computeStruct(rv)
		if err != nil {
			return 0, err
		}
		ret += size

	case reflect.Pointer:
		if rv.IsNil() {
			return ret, nil
		}
		size, err := e.computeSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = size

	case reflect.Interface:
		size, err := e.computeSize(rv.Elem())
		if err != nil {
			return 0, err
		}
		ret = size

	case reflect.Invalid:
		// Return nil

	default:
		return 0, fmt.Errorf("unsupported type(%v)", rv.Kind())
	}

	return ret, nil
}

func (e *encoder) add(rv reflect.Value, offset int) int {

	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := rv.Uint()
		offset = e.writeUint(v, offset)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		v := rv.Int()
		offset = e.writeInt(v, offset)

	case reflect.Float32:
		offset = e.writeFloat32(rv.Float(), offset)

	case reflect.Float64:
		offset = e.writeFloat64(rv.Float(), offset)

	case reflect.Bool:
		offset = e.writeBool(rv.Bool(), offset)

	case reflect.String:
		offset = e.writeString(rv.String(), offset)

	case reflect.Slice:
		if rv.IsNil() {
			return e.writeNil(offset)
		}
		l := rv.Len()

		// Format slice
		offset = e.writeSliceLength(l, offset)

		if offset, find := e.writeFixSlice(rv, offset); find {
			return offset
		}

		// Encode func
		elem := rv.Type().Elem()
		var f structWriteFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructWriter()
		} else {
			f = e.add
		}

		// Encode objects
		for i := 0; i < l; i++ {
			offset = f(rv.Index(i), offset)
		}

	case reflect.Array:
		l := rv.Len()

		// Format array same as slice
		offset = e.writeSliceLength(l, offset)

		// Encode func
		elem := rv.Type().Elem()
		var f structWriteFunc
		if elem.Kind() == reflect.Struct {
			f = e.getStructWriter()
		} else {
			f = e.add
		}

		// Encode objects
		for i := 0; i < l; i++ {
			offset = f(rv.Index(i), offset)
		}

	case reflect.Map:
		if rv.IsNil() {
			return e.writeNil(offset)
		}

		l := rv.Len()
		offset = e.writeMapLength(l, offset)

		if offset, find := e.writeFixMap(rv, offset); find {
			return offset
		}

		// Fill in keys and values
		p := rv.Pointer()
		for i := range e.mk[p] {
			offset = e.add(e.mk[p][i], offset)
			offset = e.add(e.mv[p][i], offset)
		}

	case reflect.Struct:
		offset = e.writeStruct(rv, offset)

	case reflect.Pointer:
		if rv.IsNil() {
			return e.writeNil(offset)
		}
		offset = e.add(rv.Elem(), offset)

	case reflect.Interface:
		offset = e.add(rv.Elem(), offset)

	case reflect.Invalid:
		return e.writeNil(offset)

	}
	return offset
}
