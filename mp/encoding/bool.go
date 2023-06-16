package encoding

import "github.com/romanzac/json-mp/mp/def"

func (e *encoder) writeBool(v bool, offset int) int {
	if v {
		offset = e.setByte1Int(def.True, offset)
	} else {
		offset = e.setByte1Int(def.False, offset)
	}
	return offset
}
