package decoding

import (
	"encoding/binary"
	"reflect"
	"sync"

	"github.com/romanzac/json-mp/mp/def"
)

type structCache struct {
	keys    [][]byte
	indexes []int
}

// Struct cache is stored as Map
var mapSC = sync.Map{}

func (d *decoder) setStruct(rv reflect.Value, offset int, k reflect.Kind) (int, error) {
	l, o, err := d.mapLength(offset, k)
	if err != nil {
		return 0, err
	}

	if err = d.hasRequiredLeastMapSize(o, l); err != nil {
		return 0, err
	}

	var sc *structCache
	cache, cacheFind := mapSC.Load(rv.Type())
	if !cacheFind {
		sc = &structCache{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, name := def.CheckStructField(rv.Type().Field(i)); ok {
				sc.keys = append(sc.keys, []byte(name))
				sc.indexes = append(sc.indexes, i)
			}
		}
		mapSC.Store(rv.Type(), sc)
	} else {
		sc = cache.(*structCache)
	}

	for i := 0; i < l; i++ {
		dataKey, o2, err := d.asStringByte(o, k)
		if err != nil {
			return 0, err
		}

		fieldIndex := -1
		for keyIndex, keyBytes := range sc.keys {
			if len(keyBytes) != len(dataKey) {
				continue
			}

			fieldIndex = sc.indexes[keyIndex]
			for dataIndex := range dataKey {
				if dataKey[dataIndex] != keyBytes[dataIndex] {
					fieldIndex = -1
					break
				}
			}
			if fieldIndex >= 0 {
				break
			}
		}

		if fieldIndex >= 0 {
			o2, err = d.decode(rv.Field(fieldIndex), o2)
			if err != nil {
				return 0, err
			}
		} else {
			o2, err = d.jumpOffset(o2)
			if err != nil {
				return 0, err
			}
		}
		o = o2
	}
	return o, nil
}

func (d *decoder) jumpOffset(offset int) (int, error) {
	code, offset, err := d.readSize1(offset)
	if err != nil {
		return 0, err
	}

	switch {
	case code == def.True, code == def.False, code == def.Nil:
		// No change in offset - do nothing

	case d.isPositiveFixNum(code) || d.isNegativeFixNum(code):
		// No change in offset - do nothing
	case code == def.Uint8, code == def.Int8:
		offset += def.Byte1
	case code == def.Uint16, code == def.Int16:
		offset += def.Byte2
	case code == def.Uint32, code == def.Int32, code == def.Float32:
		offset += def.Byte4
	case code == def.Uint64, code == def.Int64, code == def.Float64:
		offset += def.Byte8

	case d.isFixString(code):
		offset += int(code - def.FixStr)
	case code == def.Str8:
		b, o, err := d.readSize1(offset)
		if err != nil {
			return 0, err
		}
		o += int(b)
		offset = o
	case code == def.Str16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		o += int(binary.BigEndian.Uint16(bs))
		offset = o
	case code == def.Str32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		o += int(binary.BigEndian.Uint32(bs))
		offset = o

	case d.isFixSlice(code):
		l := int(code - def.FixArray)
		for i := 0; i < l; i++ {
			offset, err = d.jumpOffset(offset)
			if err != nil {
				return 0, err
			}
		}
	case code == def.Array16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o
	case code == def.Array32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	case d.isFixMap(code):
		l := int(code - def.FixMap)
		for i := 0; i < l*2; i++ {
			offset, err = d.jumpOffset(offset)
			if err != nil {
				return 0, err
			}
		}
	case code == def.Map16:
		bs, o, err := d.readSize2(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint16(bs))
		for i := 0; i < l*2; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o
	case code == def.Map32:
		bs, o, err := d.readSize4(offset)
		if err != nil {
			return 0, err
		}
		l := int(binary.BigEndian.Uint32(bs))
		for i := 0; i < l*2; i++ {
			o, err = d.jumpOffset(o)
			if err != nil {
				return 0, err
			}
		}
		offset = o

	}
	return offset, nil
}
