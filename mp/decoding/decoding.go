package decoding

import (
	"fmt"
	"reflect"
)

type decoder struct {
	data []byte
}

func Decode(data []byte, v interface{}) error {
	d := decoder{data: data}

	if d.data == nil || len(d.data) < 1 {
		return fmt.Errorf("empty data - nothing to unmarshall")
	}
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("initialized pointer value is expected, but got: %t", v)
	}

	rv = rv.Elem()

	last, err := d.decode(rv, 0)
	if err != nil {
		return err
	}
	if len(data) != last {
		return fmt.Errorf("unmarshall failed at size=%d, last=%d", len(data), last)
	}
	return err
}

func (d *decoder) decode(rv reflect.Value, offset int) (int, error) {
	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, o, err := d.asInt(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetInt(v)
		offset = o

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, o, err := d.asUint(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetUint(v)
		offset = o

	case reflect.Float32:
		v, o, err := d.asFloat32(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(float64(v))
		offset = o

	case reflect.Float64:
		v, o, err := d.asFloat64(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetFloat(v)
		offset = o

	case reflect.String:
		v, o, err := d.asString(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetString(v)
		offset = o

	case reflect.Bool:
		v, o, err := d.asBool(offset, k)
		if err != nil {
			return 0, err
		}
		rv.SetBool(v)
		offset = o

	case reflect.Slice:
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}
		// Decode string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			bs, offset, err := d.asStringByteByLength(offset, l)
			if err != nil {
				return 0, err
			}
			rv.SetBytes(bs)
			return offset, nil
		}

		l, o, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
			return 0, err
		}

		// Check fix type for slice
		fixOffset, found, err := d.asFixSlice(rv, o, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixOffset, nil
		}

		// Add slice
		tmpSlice := reflect.MakeSlice(rv.Type(), l, l)
		for i := 0; i < l; i++ {
			v := tmpSlice.Index(i)
			if v.Kind() == reflect.Struct {
				o, err = d.setStruct(v, o, k)
			} else {
				o, err = d.decode(v, o)
			}
			if err != nil {
				return 0, err
			}
		}
		rv.Set(tmpSlice)
		offset = o

	case reflect.Array:
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}

		// Decode string to bytes
		if d.isCodeString(d.data[offset]) {
			l, offset, err := d.stringByteLength(offset, k)
			if err != nil {
				return 0, err
			}
			if l > rv.Len() {
				return 0, fmt.Errorf("%v len is %d, but messagepack has %d elements", rv.Type(), rv.Len(), l)
			}
			bs, offset, err := d.asStringByteByLength(offset, l)
			if err != nil {
				return 0, err
			}
			for i, b := range bs {
				rv.Index(i).SetUint(uint64(b))
			}
			return offset, nil
		}

		l, o, err := d.sliceLength(offset, k)
		if err != nil {
			return 0, err
		}

		if l > rv.Len() {
			return 0, fmt.Errorf("%v len is %d, but messagepack has %d elements", rv.Type(), rv.Len(), l)
		}

		if err = d.hasRequiredLeastSliceSize(o, l); err != nil {
			return 0, err
		}

		// Add array
		for i := 0; i < l; i++ {
			o, err = d.decode(rv.Index(i), o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	case reflect.Map:
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}

		l, o, err := d.mapLength(offset, k)
		if err != nil {
			return 0, err
		}

		if err = d.hasRequiredLeastMapSize(o, l); err != nil {
			return 0, err
		}

		// Check fix map type
		fixOffset, found, err := d.asFixMap(rv, o, l)
		if err != nil {
			return 0, err
		}
		if found {
			return fixOffset, nil
		}

		// Add elements dynamically
		key := rv.Type().Key()
		value := rv.Type().Elem()
		if rv.IsNil() {
			rv.Set(reflect.MakeMapWithSize(rv.Type(), l))
		}
		for i := 0; i < l; i++ {
			k := reflect.New(key).Elem()
			v := reflect.New(value).Elem()
			o, err = d.decode(k, o)
			if err != nil {
				return 0, err
			}
			o, err = d.decode(v, o)
			if err != nil {
				return 0, err
			}

			rv.SetMapIndex(k, v)
		}
		offset = o

	case reflect.Struct:
		o, err := d.setStruct(rv, offset, k)
		if err != nil {
			return 0, err
		}
		offset = o

	case reflect.Pointer:
		if d.isCodeNil(d.data[offset]) {
			offset++
			return offset, nil
		}

		if rv.Elem().Kind() == reflect.Invalid {
			n := reflect.New(rv.Type().Elem())
			rv.Set(n)
		}

		o, err := d.decode(rv.Elem(), offset)
		if err != nil {
			return 0, err
		}
		offset = o

	case reflect.Interface:
		if rv.Elem().Kind() == reflect.Pointer {
			o, err := d.decode(rv.Elem(), offset)
			if err != nil {
				return 0, err
			}
			offset = o
		} else {
			v, o, err := d.asInterface(offset, k)
			if err != nil {
				return 0, err
			}
			if v != nil {
				rv.Set(reflect.ValueOf(v))
			}
			offset = o
		}

	default:
		return 0, fmt.Errorf("unsupported type(%v)", rv.Kind())
	}
	return offset, nil
}

func (d *decoder) errorTemplate(code byte, k reflect.Kind) error {
	return fmt.Errorf("invalid code %x decoding %v", code, k)
}
