package mp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/romanzac/json-mp/mp/def"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestIntFixMinMax(t *testing.T) {
	var r int
	if err := encodeDecode(-16, &r, func(code byte) bool {
		return def.NegativeFixIntMin <= int8(code) && int8(code) <= def.NegativeFixIntMax
	}); err != nil {
		t.Error(err)
	}
}

func TestIntNeg8(t *testing.T) {
	var r int
	if err := encodeDecode(-124, &r, func(code byte) bool {
		return code == def.Int8
	}); err != nil {
		t.Error(err)
	}
}

func TestIntNeg16(t *testing.T) {
	var r int
	if err := encodeDecode(-30109, &r, func(code byte) bool {
		return code == def.Int16
	}); err != nil {
		t.Error(err)
	}
}
func TestIntNeg32(t *testing.T) {
	var r int
	if err := encodeDecode(-1030106, &r, func(code byte) bool {
		return code == def.Int32
	}); err != nil {
		t.Error(err)
	}
}
func TestIntPos64(t *testing.T) {
	var r int64
	if err := encodeDecode(int64(math.MinInt64+12345), &r, func(code byte) bool {
		return code == def.Int64
	}); err != nil {
		t.Error(err)
	}
}

func TestIntErr8(t *testing.T) {
	var r uint8
	if err := encodeDecode(-8, &r, func(code byte) bool {
		return def.NegativeFixIntMin <= int8(code) && int8(code) <= def.NegativeFixIntMax
	}); err == nil || !strings.Contains(err.Error(), "different value") {
		t.Error("error")
	}
}
func TestIntErr64(t *testing.T) {
	var r int32
	if err := encodeDecode(int64(math.MinInt64+12345), &r, func(code byte) bool {
		return code == def.Int64
	}); err == nil || !strings.Contains(err.Error(), "different value") {
		t.Error("error")
	}
}

func TestUintFixMinMax(t *testing.T) {
	var v, r uint
	v = 8
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixIntMin <= uint8(code) && uint8(code) <= def.FixIntMax
	}); err != nil {
		t.Error(err)
	}
}
func TestUint8(t *testing.T) {
	var v, r uint
	v = 130
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Uint8
	}); err != nil {
		t.Error(err)
	}
}
func TestUint16(t *testing.T) {
	var v, r uint
	v = 30130
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Uint16
	}); err != nil {
		t.Error(err)
	}
}
func TestUint32(t *testing.T) {
	var v, r uint
	v = 1030130
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Uint32
	}); err != nil {
		t.Error(err)
	}
}
func TestUint64(t *testing.T) {
	var r uint64
	if err := encodeDecode(uint64(math.MaxUint64-12345), &r, func(code byte) bool {
		return code == def.Uint64
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatZero32(t *testing.T) {

	var v, r float32
	v = 0
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float32
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatNeg32(t *testing.T) {
	var v, r float32
	v = -5
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float32
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatMin32(t *testing.T) {
	var v, r float32
	v = math.SmallestNonzeroFloat32
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float32
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatMax32(t *testing.T) {
	var v, r float32
	v = math.MaxFloat32
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float32
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatZero64(t *testing.T) {
	var v, r float64
	v = 0
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatNeg64(t *testing.T) {
	var v, r float64
	v = -3
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatMin64(t *testing.T) {
	var v, r float64
	v = math.SmallestNonzeroFloat64
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatMax64(t *testing.T) {
	var v, r float64
	v = math.MaxFloat64
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err != nil {
		t.Error(err)
	}
}

func TestFloatErr32(t *testing.T) {
	var r int
	v := float32(2.345)
	b, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}

	err = Unmarshal(b, &r)
	if err != nil {
		t.Error(err)
	}

	if r != 2 {
		t.Error("different value", r)
	}
}

func TestFloatErrMax32(t *testing.T) {
	var v float32
	var r float64
	v = math.MaxFloat32
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float32
	}); err == nil || !strings.Contains(err.Error(), "different value") {
		t.Error(err)
	}
}

func TestFloatErrMax64(t *testing.T) {
	var v float64
	var r float32
	v = math.MaxFloat64
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err == nil || !strings.Contains(err.Error(), "invalid code cb decoding") {
		t.Error("error")
	}
}

func TestFloatErrStr64(t *testing.T) {
	var v float64
	var r string
	v = math.MaxFloat64
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Float64
	}); err == nil || !strings.Contains(err.Error(), "invalid code cb decoding") {
		t.Error("error")
	}
}

func TestBoolTrue(t *testing.T) {

	var v, r bool
	v = true
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.True
	}); err != nil {
		t.Error(err)
	}
}

func TestBoolFalse(t *testing.T) {
	var v, r bool
	v = false
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.False
	}); err != nil {
		t.Error(err)
	}
}

func TestBoolErr(t *testing.T) {
	var v bool
	var r uint8
	v = true
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.True
	}); err == nil || !strings.Contains(err.Error(), "invalid code c3 decoding") {
		t.Error("error")
	}
}

func TestNil(t *testing.T) {
	{
		var r *map[interface{}]interface{}
		d, err := Marshal(nil)
		if err != nil {
			t.Error(err)
		}
		if d[0] != def.Nil {
			t.Error("not nil type")
		}
		err = Unmarshal(d, &r)
		if err != nil {
			t.Error(err)
		}
		if r != nil {
			t.Error("not nil")
		}
	}
}

func TestStringFix(t *testing.T) {
	var v, r string
	v = ""
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixStr <= code && code < def.FixStr+32
	}); err != nil {
		t.Error(err)
	}
}

func TestStringStr8(t *testing.T) {
	var v, r string
	v = strings.Repeat("FZcF1c4e7htNU9vX3llpXg0GUwYGy59", 8)
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Str8
	}); err != nil {
		t.Error(err)
	}
}

func TestStringStr16(t *testing.T) {
	var v, r string
	v = strings.Repeat("FZcF1c4e7htNU9vX3llpXg0GUwYGy59",
		(math.MaxUint16/len("FZcF1c4e7htNU9vX3llpXg0GUwYGy59"))-1)
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Str16
	}); err != nil {
		t.Error(err)
	}
}

func TestStringStr32(t *testing.T) {
	var v, r string
	v = strings.Repeat("FZcF1c4e7htNU9vX3llpXg0GUwYGy59",
		(math.MaxUint16/len("FZcF1c4e7htNU9vX3llpXg0GUwYGy59"))+1)
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Str32
	}); err != nil {
		t.Error(err)
	}
}

func TestInterface(t *testing.T) {
	f := func(v interface{}) error {
		b, err := Marshal(v)
		if err != nil {
			return err
		}
		var r interface{}
		err = Unmarshal(b, &r)
		if err != nil {
			return err
		}
		if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", r) {
			return fmt.Errorf("different value %v, %v", v, r)
		}
		return err
	}

	a1 := make([]int, math.MaxUint16)
	a2 := make([]int, math.MaxUint16+1)
	m1 := map[string]int{}
	m2 := map[string]int{}

	for i := range a1 {
		a1[i] = i
		m1[fmt.Sprint(i)] = 1
	}
	for i := range a2 {
		a2[i] = i
		m2[fmt.Sprint(i)] = 1
	}

	vars := []interface{}{
		true, false,
		1,
		math.MaxUint8, math.MaxUint16,
		math.MaxUint32, math.MaxUint32 + 1,
		math.MaxFloat32, math.MaxFloat64,
		"z",
		strings.Repeat("n", math.MaxUint8),
		strings.Repeat("o", math.MaxUint16),
		strings.Repeat("i", math.MaxUint16+1),
		[]int{1, 5, 3},
		a1, a2,
		map[string]interface{}{"one": 1, "twelve": "a"}, m1, m2,
	}

	for i, v := range vars {
		if err := f(v); err != nil {
			t.Error(i, err)
		}
	}

	// Error case
	var r interface{}
	err := Unmarshal([]byte{def.Float32}, &r)
	if err == nil {
		t.Error("error must occur")
	}
	if err != nil && !strings.Contains(err.Error(), "too short bytes") {
		t.Error(err)
	}
}

func TestArrayNil(t *testing.T) {
	var v, r []int
	v = nil
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.Nil == code
	}); err != nil {
		t.Error(err)
	}
}

func TestArrayFix(t *testing.T) {
	var v, r []int
	v = make([]int, 15)
	for i := range v {
		v[i] = rand.Intn(math.MaxInt32)
	}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestArray16(t *testing.T) {
	var v, r []int
	v = make([]int, 30015)
	for i := range v {
		v[i] = rand.Intn(math.MaxInt32)
	}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Array16
	}); err != nil {
		t.Error(err)
	}
}

func TestArray32(t *testing.T) {
	var v, r []int
	v = make([]int, 1030015)
	for i := range v {
		v[i] = rand.Intn(math.MaxInt32)
	}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Array32
	}); err != nil {
		t.Error(err)
	}
}

func TestArrayString(t *testing.T) {
	v := "ab42e"
	var r [5]byte
	b, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}
	err = Unmarshal(b, &r)
	if err != nil {
		t.Error(err)
	}
	if v != string(r[:]) {
		t.Errorf("different value %v, %v", v, string(r[:]))
	}
}

func TestFixSliceArrayNeg(t *testing.T) {

	var v, r []int
	v = []int{-1, 1}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayPos(t *testing.T) {
	var v, r []uint
	v = []uint{0, 100}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMinMaxInt8(t *testing.T) {
	var v, r []int8
	v = []int8{math.MinInt8, math.MaxInt8}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMinMaxInt16(t *testing.T) {
	var v, r []int16
	v = []int16{math.MinInt16, math.MaxInt16}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMinMaxInt32(t *testing.T) {
	var v, r []int32
	v = []int32{math.MinInt32, math.MaxInt32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMinMaxInt64(t *testing.T) {
	var v, r []int64
	v = []int64{math.MinInt64, math.MaxInt64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMaxUInt16(t *testing.T) {
	var v, r []uint16
	v = []uint16{0, math.MaxUint16}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMaxUInt32(t *testing.T) {
	var v, r []uint32
	v = []uint32{0, math.MaxUint32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMaxUInt64(t *testing.T) {
	var v, r []uint64
	v = []uint64{0, math.MaxUint64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMinMaxUInt16(t *testing.T) {
	var v, r []float32
	v = []float32{math.SmallestNonzeroFloat32, math.MaxFloat32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayMaxFloat64(t *testing.T) {
	var v, r []float64
	v = []float64{math.SmallestNonzeroFloat64, math.MaxFloat64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayStr(t *testing.T) {
	var v, r []string
	v = []string{"423csf", "r23fs23a"}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixSliceArrayBool(t *testing.T) {
	var v, r []bool
	v = []bool{true, false}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixArray <= code && code <= def.FixArray+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringInt(t *testing.T) {

	var v, r map[string]int
	v = map[string]int{"a": 1, "b": 2}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringUInt(t *testing.T) {
	var v, r map[string]uint
	v = map[string]uint{"n": math.MaxUint32, "c": 0}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringString(t *testing.T) {
	var v, r map[string]string
	v = map[string]string{"n": "34523650", "FZcF1c4e7htNU9vX3llpXg0GUwYGy59": "FZcF1c4e7htN33333333333"}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringFloat32(t *testing.T) {
	var v, r map[string]float32
	v = map[string]float32{"a": math.MaxFloat32, "b": math.SmallestNonzeroFloat32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringFloat64(t *testing.T) {
	var v, r map[string]float64
	v = map[string]float64{"a": math.MaxFloat64, "b": math.SmallestNonzeroFloat64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringBool(t *testing.T) {
	var v, r map[string]bool
	v = map[string]bool{"x": true, "y": false}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringInt8(t *testing.T) {
	var v, r map[string]int8
	v = map[string]int8{"x": math.MinInt8, "y": math.MaxInt8}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringInt16(t *testing.T) {
	var v, r map[string]int16
	v = map[string]int16{"x": math.MaxInt16, "y": math.MinInt16}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringInt32(t *testing.T) {
	var v, r map[string]int32
	v = map[string]int32{"x": math.MaxInt32, "y": math.MinInt32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringInt64(t *testing.T) {
	var v, r map[string]int64
	v = map[string]int64{"x": math.MinInt64, "y": math.MaxInt64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringUInt8(t *testing.T) {
	var v, r map[string]uint8
	v = map[string]uint8{"x": 0, "y": math.MaxUint8}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringUInt16(t *testing.T) {
	var v, r map[string]uint16
	v = map[string]uint16{"x": 0, "y": math.MaxUint16}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringUInt32(t *testing.T) {
	var v, r map[string]uint32
	v = map[string]uint32{"x": 0, "y": math.MaxUint32}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestFixMapStringUInt64(t *testing.T) {
	var v, r map[string]uint64
	v = map[string]uint64{"x": 0, "y": math.MaxUint64}
	if err := encodeDecode(v, &r, func(code byte) bool {
		return def.FixMap <= code && code <= def.FixMap+0x0f
	}); err != nil {
		t.Error(err)
	}
}

func TestMapStringInt(t *testing.T) {
	var v map[string]int
	var r map[int]int
	v = make(map[string]int, 100)
	for i := 0; i < 100; i++ {
		v[fmt.Sprintf("%03d", i)] = i
	}
	d, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.Map16 {
		t.Error("code different")
	}
	err = Unmarshal(d, &r)
	if err == nil || !strings.Contains(err.Error(), "invalid code a3 decoding") {
		t.Error("error")
	}
}

func TestPointerUint8(t *testing.T) {
	var v, r *int
	vv := 250
	v = &vv
	if err := encodeDecode(v, &r, func(code byte) bool {
		return code == def.Uint8
	}); err != nil {
		t.Error(err)
	}
}

func TestPointerNil(t *testing.T) {
	var v, r *int
	d, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.Nil {
		t.Error("code different")
	}
	err = Unmarshal(d, &r)
	if err != nil {
		t.Error(err)
	}
	if v != r {
		t.Error("different value")
	}
}

func TestPointerErr(t *testing.T) {
	var v *int
	var r int
	if err := encodeDecode(v, r, func(code byte) bool {
		return code == def.Nil
	}); err == nil || !strings.Contains(err.Error(), "initialized pointer value is expected, but got:") {
		t.Error(err)
	}
}

func TestUnsupportedUintPtr(t *testing.T) {
	var v, r uintptr
	_, err := Marshal(v)
	if !strings.Contains(err.Error(), "unsupported type(uintptr)") {
		t.Error("test error", err)
	}
	err = Unmarshal([]byte{0xc0}, &r)
	if !strings.Contains(err.Error(), "unsupported type(uintptr)") {
		t.Error("test error", err)
	}
}

func TestUnsupportedChan(t *testing.T) {
	var v, r chan string
	_, err := Marshal(v)
	if !strings.Contains(err.Error(), "unsupported type(chan)") {
		t.Error("test error", err)
	}
	err = Unmarshal([]byte{0xc0}, &r)
	if !strings.Contains(err.Error(), "unsupported type(chan)") {
		t.Error("test error", err)
	}
}

func TestUnsupportedFunc(t *testing.T) {
	var v, r func()
	_, err := Marshal(v)
	if !strings.Contains(err.Error(), "unsupported type(func)") {
		t.Error("test error", err)
	}
	err = Unmarshal([]byte{0xc0}, &r)
	if !strings.Contains(err.Error(), "unsupported type(func)") {
		t.Error("test error", err)
	}
}

func TestUnsupportedErr(t *testing.T) {
	var v, r error
	bb, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}
	if bb[0] != def.Nil {
		t.Errorf("code is different %d, %d", bb[0], def.Nil)
	}
	err = Unmarshal([]byte{0xc0}, &r)
	if err != nil {
		t.Error(err)
	}
	if r != nil {
		t.Error("error should be nil")
	}
}

func TestStructEmbedded(t *testing.T) {
	type Emb struct {
		Int int
	}
	type A struct {
		Emb
	}
	v := A{Emb: Emb{Int: 2}}
	b, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}

	var vv A
	err = Unmarshal(b, &vv)
	if err != nil {
		t.Error(err)
	}
	if v.Int != vv.Int {
		t.Errorf("value is different %v, %v", v, vv)
	}
}

func TestStructTag(t *testing.T) {
	type vSt struct {
		One int    `json:"Three"`
		Two string `json:"four"`
		Hfn bool   `json:"-"`
	}
	type rSt struct {
		Three int
		Four  string `json:"four"`
		Hfn   bool
	}

	v := vSt{One: 1, Two: "2", Hfn: true}
	r := rSt{}

	d, err := Marshal(v)
	if err != nil {
		t.Error(err)
	}
	if d[0] != def.FixMap+0x02 {
		t.Error("code different")
	}
	err = Unmarshal(d, &r)
	if err != nil {
		t.Error(err)
	}
	if v.One != r.Three || v.Two != r.Four || r.Hfn != false {
		t.Error("error:", v, r)
	}
}

func TestStructJump(t *testing.T) {
	type v1 struct{ A interface{} }
	type r1 struct{ B interface{} }

	f := func(v v1) error {
		b, err := Marshal(v)
		if err != nil {
			return err
		}
		var r r1
		err = Unmarshal(b, &r)
		if err != nil {
			return err
		}
		if fmt.Sprint(v.A) == fmt.Sprint(r.B) {
			return fmt.Errorf("value equal %v, %v", v, r)
		}
		return nil
	}

	a1 := make([]int, math.MaxUint16)
	a2 := make([]int, math.MaxUint16+1)
	m1 := map[string]int{}
	m2 := map[string]int{}

	for i := range a1 {
		a1[i] = i
		m1[fmt.Sprint(i)] = 1
	}
	for i := range a2 {
		a2[i] = i
		m2[fmt.Sprint(i)] = 1
	}

	vs := []v1{
		{A: true},
		{A: 1}, {A: -1},
		{A: []int{1}}, {A: a1}, {A: a2},
		{A: math.MaxUint8}, {A: math.MinInt8},
		{A: math.MaxUint16}, {A: math.MinInt16},
		{A: math.MaxUint32 + 1}, {A: math.MinInt32 - 1}, {A: math.MaxFloat64},
		{A: "a"},
		{A: strings.Repeat("b", math.MaxUint8)}, {A: []byte(strings.Repeat("c", math.MaxUint8))},
		{A: strings.Repeat("e", math.MaxUint16)}, {A: []byte(strings.Repeat("d", math.MaxUint16))},
		{A: strings.Repeat("f", math.MaxUint16+1)}, {A: []byte(strings.Repeat("g", math.MaxUint16+1))},
		{A: map[string]int{"z": 1}}, {A: m1}, {A: m2},
	}

	for i, v := range vs {
		if err := f(v); err != nil {
			t.Error(i, err)
		}
	}

}

func encodeDecode(v, r interface{}, j func(d byte) bool) error {
	d, err := Marshal(v)
	if err != nil {
		return err
	}
	if !j(d[0]) {
		return fmt.Errorf("different %s", hex.Dump(d))
	}
	if err := Unmarshal(d, r); err != nil {
		return err
	}
	if err := equalCheck(v, r); err != nil {
		return err
	}
	return nil
}

func getVal(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	return rv.Interface()
}

func equalCheck(in, out interface{}) error {
	i := getVal(in)
	o := getVal(out)
	if !reflect.DeepEqual(i, o) {
		return errors.New(fmt.Sprint("different value \n[in]:", i, " \n[out]:", o))
	}
	return nil
}
