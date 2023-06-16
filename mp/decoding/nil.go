package decoding

import "github.com/romanzac/json-mp/mp/def"

func (d *decoder) isCodeNil(v byte) bool {
	return def.Nil == v
}
