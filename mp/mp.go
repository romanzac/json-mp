package mp

import (
	"github.com/romanzac/json-mp/mp/decoding"
	"github.com/romanzac/json-mp/mp/encoding"
)

// Marshal returns the MessagePack byte array of data in v with shape defined in JSONData
func Marshal(v interface{}) ([]byte, error) {
	return encoding.Encode(v)
}

// Unmarshal reads the MessagePack-encoded data and interprets them according to
// shape object stored in JSONData (v)
func Unmarshal(data []byte, v interface{}) error {
	return decoding.Decode(data, v)
}
