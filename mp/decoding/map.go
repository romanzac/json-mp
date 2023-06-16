package decoding

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/romanzac/json-mp/mp/def"
)

var (
	typeMapStringInt   = reflect.TypeOf(map[string]int{})
	typeMapStringInt8  = reflect.TypeOf(map[string]int8{})
	typeMapStringInt16 = reflect.TypeOf(map[string]int16{})
	typeMapStringInt32 = reflect.TypeOf(map[string]int32{})
	typeMapStringInt64 = reflect.TypeOf(map[string]int64{})

	typeMapStringUint   = reflect.TypeOf(map[string]uint{})
	typeMapStringUint8  = reflect.TypeOf(map[string]uint8{})
	typeMapStringUint16 = reflect.TypeOf(map[string]uint16{})
	typeMapStringUint32 = reflect.TypeOf(map[string]uint32{})
	typeMapStringUint64 = reflect.TypeOf(map[string]uint64{})

	typeMapStringFloat32 = reflect.TypeOf(map[string]float32{})
	typeMapStringFloat64 = reflect.TypeOf(map[string]float64{})

	typeMapStringBool   = reflect.TypeOf(map[string]bool{})
	typeMapStringString = reflect.TypeOf(map[string]string{})
)

func (d *decoder) isFixMap(v byte) bool {
	return def.FixMap <= v && v <= def.FixMap+0x0f
}

func (d *decoder) mapLength(offset int, k reflect.Kind) (int, int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, 0, err
	}

	switch {
	case d.isFixMap(code):
		return int(code - def.FixMap), offset, nil
	case code == def.Map16:
		bs, offset, err := d.readSize2(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint16(bs)), offset, nil
	case code == def.Map32:
		bs, offset, err := d.readSize4(offset)
		if err != nil {
			return 0, 0, err
		}
		return int(binary.BigEndian.Uint32(bs)), offset, nil
	}

	return 0, 0, d.errorTemplate(code, k)
}

func (d *decoder) hasRequiredLeastMapSize(offset, length int) error {
	if len(d.data[offset:]) < length*2 {
		return errors.New("data length lacks to add map")
	}
	return nil
}

func (d *decoder) asFixMap(rv reflect.Value, offset int, l int) (int, bool, error) {
	t := rv.Type()

	keyKind := t.Key().Kind()
	valueKind := t.Elem().Kind()

	switch t {
	case typeMapStringInt:
		m := make(map[string]int, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringUint:
		m := make(map[string]uint, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asUint(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = uint(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringFloat32:
		m := make(map[string]float32, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asFloat32(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringFloat64:
		m := make(map[string]float64, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asFloat64(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringBool:
		m := make(map[string]bool, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asBool(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringString:
		m := make(map[string]string, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asString(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringInt8:
		m := make(map[string]int8, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int8(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringInt16:
		m := make(map[string]int16, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int16(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringInt32:
		m := make(map[string]int32, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = int32(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringInt64:
		m := make(map[string]int64, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asInt(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringUint8:
		m := make(map[string]uint8, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asUint(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = uint8(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil
	case typeMapStringUint16:
		m := make(map[string]uint16, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asUint(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = uint16(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringUint32:
		m := make(map[string]uint32, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asUint(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = uint32(v)
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil

	case typeMapStringUint64:
		m := make(map[string]uint64, l)
		for i := 0; i < l; i++ {
			k, o, err := d.asString(offset, keyKind)
			if err != nil {
				return 0, false, err
			}
			v, o, err := d.asUint(o, valueKind)
			if err != nil {
				return 0, false, err
			}
			m[k] = v
			offset = o
		}
		rv.Set(reflect.ValueOf(m))
		return offset, true, nil
	}

	return offset, false, nil
}
