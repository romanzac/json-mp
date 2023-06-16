package encoding

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	"github.com/romanzac/json-mp/mp/def"
)

type structCache struct {
	indexes []int
	names   []string
}

var mapSC = sync.Map{}

type structWriteFunc func(rv reflect.Value, offset int) int

func (e *encoder) computeStruct(rv reflect.Value) (int, error) {
	ret := 0
	t := rv.Type()
	cache, find := mapSC.Load(t)
	var c *structCache
	if !find {
		c = &structCache{}
		for i := 0; i < rv.NumField(); i++ {
			if ok, name := def.CheckStructField(rv.Type().Field(i)); ok {
				keySize := def.Byte1 + e.computeString(name)
				valueSize, err := e.computeSize(rv.Field(i))
				if err != nil {
					return 0, err
				}
				ret += keySize + valueSize
				c.indexes = append(c.indexes, i)
				c.names = append(c.names, name)
			}
		}
		mapSC.Store(t, c)
	} else {
		c = cache.(*structCache)
		for i := 0; i < len(c.indexes); i++ {
			keySize := def.Byte1 + e.computeString(c.names[i])
			valueSize, err := e.computeSize(rv.Field(c.indexes[i]))
			if err != nil {
				return 0, err
			}
			ret += keySize + valueSize
		}
	}

	// Check format size
	l := len(c.indexes)
	if l <= 0x0f {
		// Do nothing - format code only
	} else if l <= math.MaxUint16 {
		ret += def.Byte2
	} else if uint(l) <= math.MaxUint32 {
		ret += def.Byte4
	} else {
		return 0, fmt.Errorf("not support this array length : %d", l)
	}
	return ret, nil
}

func (e *encoder) getStructWriter() structWriteFunc {
	return e.writeStruct
}

func (e *encoder) writeStruct(rv reflect.Value, offset int) int {

	cache, _ := mapSC.Load(rv.Type())
	c := cache.(*structCache)

	// Check format size
	num := len(c.indexes)
	if num <= 0x0f {
		offset = e.setByte1Int(def.FixMap+num, offset)
	} else if num <= math.MaxUint16 {
		offset = e.setByte1Int(def.Map16, offset)
		offset = e.setByte2Int(num, offset)
	} else if uint(num) <= math.MaxUint32 {
		offset = e.setByte1Int(def.Map32, offset)
		offset = e.setByte4Int(num, offset)
	}

	for i := 0; i < num; i++ {
		offset = e.writeString(c.names[i], offset)
		offset = e.add(rv.Field(c.indexes[i]), offset)
	}
	return offset
}
