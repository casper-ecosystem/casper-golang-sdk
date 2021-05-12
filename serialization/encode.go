package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
)

// MustMarshal is the panic-on-fail version of Marshal
func MustMarshal(value interface{}) []byte {
	res, err := Marshal(value)
	if err != nil {
		panic(err)
	}
	return res
}

// Marshal returns the bytes representation of v.
func Marshal(value interface{}) ([]byte, error) {
	var w bytes.Buffer
	enc := Encoder{w: &w}
	_, err := enc.Encode(value)
	return w.Bytes(), err
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

type Encoder struct {
	w io.Writer
}

func (enc *Encoder) Encode(v interface{}) (int, error) {
	val := reflect.ValueOf(v)
	return enc.encode(val)
}

// EncodeBool encodes bool value
func (enc *Encoder) EncodeBool(v bool) (int, error) {
	if v {
		return enc.w.Write([]byte{1})
	} else {
		return enc.w.Write([]byte{0})
	}
}

// EncodeInt32 encodes int32 value
func (enc *Encoder) EncodeInt32(v int32) (int, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	return enc.w.Write(data)
}

// EncodeInt64 encodes int64 value
func (enc *Encoder) EncodeInt64(v int64) (int, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(v))
	return enc.w.Write(data)
}

// EncodeByte encodes byte value
func (enc *Encoder) EncodeByte(v byte) (int, error) {
	return enc.w.Write([]byte{v})
}

// EncodeUIn32 encodes uint32 value
func (enc *Encoder) EncodeUInt32(v uint32) (int, error) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	return enc.w.Write(data)
}

// EncodeUIn64 encodes uint64 value
func (enc *Encoder) EncodeUInt64(v uint64) (int, error) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, v)
	return enc.w.Write(data)
}

// EncodeBigInt encodes big number
func (enc *Encoder) EncodeBigInt(v big.Int) (int, error) {
	data := v.Bytes()
	for i := 0; i < len(data)/2; i++ {
		data[i], data[len(data)-i-1] = data[len(data)-i-1], data[i]
	}
	numberLen := byte(len(data))
	data = append([]byte{numberLen}, data...)
	return enc.w.Write(data)
}

// EncodeString encodes string value
func (enc *Encoder) EncodeString(v string) (int, error) {
	data := []byte(v)
	inputLen, err := enc.EncodeUInt32(uint32(len(data)))
	if err != nil {
		return 0, err
	}

	strLen, err := enc.w.Write(data)
	inputLen += strLen

	return inputLen, err
}

func (enc *Encoder) EncodeOptional(v reflect.Value) (int, error) {
	if v.IsNil() {
		return enc.EncodeBool(false)
	}

	inputLen, err := enc.EncodeBool(true)
	if err != nil {
		return 0, err
	}

	dataLen, err := enc.encode(v.Elem())
	inputLen += dataLen

	return inputLen, err
}

// EncodeArray encodes array
func (enc *Encoder) EncodeArray(v reflect.Value) (int, error) {
	inputLen, err := enc.EncodeUInt32(uint32(v.Len()))
	if err != nil {
		return 0, err
	}

	dataLen, err := enc.EncodeFixedArray(v)
	inputLen += dataLen

	return inputLen, err
}

// EncodeFixedArray encodes array with fixed size
func (enc *Encoder) EncodeFixedArray(v reflect.Value) (int, error) {
	if v.Type().Elem().Kind() == reflect.Uint8 {
		// Create a slice of the underlying array for better efficiency
		// when possible.  Can't create a slice of an unaddressable
		// value.
		if v.CanAddr() {
			return enc.EncodeFixedByteArray(v.Slice(0, v.Len()).Bytes())
		}

		// When the underlying array isn't addressable fall back to
		// copying the array into a new slice.  This is rather ugly, but
		// the inability to create a constant slice from an
		// unaddressable array is a limitation of Go.
		slice := make([]byte, v.Len(), v.Len())
		reflect.Copy(reflect.ValueOf(slice), v)
		return enc.EncodeFixedByteArray(slice)
	}

	var n int
	for i := 0; i < v.Len(); i++ {
		n2, err := enc.encode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// EncodeByteArray encodes byte array
func (enc *Encoder) EncodeByteArray(v []byte) (int, error) {
	inputLen, err := enc.EncodeUInt32(uint32(len(v)))
	if err != nil {
		return 0, err
	}

	dataLen, err := enc.EncodeFixedByteArray(v)
	inputLen += dataLen

	return inputLen, err
}

// EncodeFixedByteArray encodes byte array with fixed size
func (enc *Encoder) EncodeFixedByteArray(v []byte) (int, error) {
	return enc.w.Write(v)
}

func (enc *Encoder) EncodeResult(v reflect.Value) (int, error) {
	val := v.Interface().(Result)
	result := v.FieldByName(val.ResultFieldName()).Interface().(bool)
	n := 0
	n2, err := enc.EncodeBool(result)
	n += n2
	if err != nil {
		return n, err
	}
	if result {
		vv := v.FieldByName(val.SuccessFieldName())

		if !vv.IsValid() {
			return n, errors.New(fmt.Sprintf("invalid element: %s", val.SuccessFieldName()))
		}

		n2, err = enc.encode(vv.Elem())
	} else {
		vv := v.FieldByName(val.ErrorFieldName())

		if !vv.IsValid() {
			return n, errors.New(fmt.Sprintf("invalid element: %s", val.ErrorFieldName()))
		}

		n2, err = enc.encode(vv.Elem())
	}
	n += n2

	return n, err
}

// EncodeTuple encodes tuple
func (enc *Encoder) EncodeTuple(v reflect.Value) (int, error) {
	val := v.Interface().(Tuple)
	fields := val.TupleFields()
	n := 0
	for _, fieldName := range fields {
		field := v.FieldByName(fieldName)
		n2, err := enc.encode(field)
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// EncodeMap encodes map
func (enc *Encoder) EncodeMap(v reflect.Value) (int, error) {
	n, err := enc.EncodeUInt32(uint32(v.Len()))
	if err != nil {
		return 0, err
	}
	for _, key := range v.MapKeys() {
		n2, err := enc.encode(key)
		n += n2
		if err != nil {
			return n, err
		}

		n2, err = enc.encode(v.MapIndex(key))
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (enc *Encoder) EncodeMarshaler(v reflect.Value) (int, error) {
	marshaler := v.Interface().(Marshaler)
	return marshaler.Marshal(enc.w)
}

// EncodeStruct encodes struct
func (enc *Encoder) EncodeStruct(v reflect.Value) (int, error) {
	var n int
	vt := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// Skip unexported fields and indirect through pointers.
		vtf := vt.Field(i)
		if vtf.PkgPath != "" {
			continue
		}
		vf := v.Field(i)

		// Encode each struct field.
		n2, err := enc.encode(vf)
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// EncodeUnion encodes union
func (enc *Encoder) EncodeUnion(v reflect.Value) (int, error) {
	u := v.Interface().(Union)

	vs := byte(v.FieldByName(u.SwitchFieldName()).Uint())
	n, err := enc.EncodeByte(vs)

	if err != nil {
		return n, err
	}

	arm, ok := u.ArmForSwitch(vs)

	// void arm, we're done
	if arm == "" {
		return n, nil
	}

	vv := v.FieldByName(arm)

	if !vv.IsValid() || !ok {
		return n, errors.New(fmt.Sprintf("invalid union switch: %d", vs))
	}

	n2, err := enc.encode(vv.Elem())
	n += n2
	return n, err
}

// EncodeInterface encodes interface
func (enc *Encoder) EncodeInterface(v reflect.Value) (int, error) {
	if v.IsNil() || !v.CanInterface() {
		return 0, errors.New("can't encode nil interface")
	}

	// Extract underlying value from the interface and indirect through pointers.
	ve := reflect.ValueOf(v.Interface())
	ve = enc.indirect(ve)
	return enc.encode(ve)
}

// indirect dereferences pointers until it reaches a non-pointer.  This allows
// transparent encoding through arbitrary levels of indirection.
func (enc *Encoder) indirect(v reflect.Value) reflect.Value {
	rv := v
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}

func (enc *Encoder) encode(v reflect.Value) (int, error) {
	if _, ok := v.Interface().(Marshaler); ok {
		return enc.EncodeMarshaler(v)
	}

	switch v.Kind() {
	case reflect.Bool:
		return enc.EncodeBool(v.Bool())
	case reflect.Int32:
		return enc.EncodeInt32(int32(v.Int()))
	case reflect.Int64:
		return enc.EncodeInt64(v.Int())
	case reflect.Uint8:
		return enc.EncodeByte(byte(v.Uint()))
	case reflect.Uint32:
		return enc.EncodeUInt32(uint32(v.Uint()))
	case reflect.Uint64:
		return enc.EncodeUInt64(v.Uint())
	case reflect.String:
		return enc.EncodeString(v.String())
	case reflect.Ptr:
		return enc.EncodeOptional(v)
	case reflect.Slice:
		return enc.EncodeArray(v)
	case reflect.Array:
		return enc.EncodeFixedArray(v)
	case reflect.Map:
		return enc.EncodeMap(v)
	case reflect.Struct:
		if val, ok := v.Interface().(U128); ok {
			return enc.EncodeBigInt(val.Int)
		}
		if val, ok := v.Interface().(U256); ok {
			return enc.EncodeBigInt(val.Int)
		}
		if val, ok := v.Interface().(U512); ok {
			return enc.EncodeBigInt(val.Int)
		}
		if val, ok := v.Interface().(big.Int); ok {
			return enc.EncodeBigInt(val)
		}
		if _, ok := v.Interface().(Result); ok {
			return enc.EncodeResult(v)
		}
		if _, ok := v.Interface().(Tuple); ok {
			return enc.EncodeTuple(v)
		}
		if _, ok := v.Interface().(Union); ok {
			return enc.EncodeUnion(v)
		}

		return enc.EncodeStruct(v)
	case reflect.Interface:
		return enc.EncodeInterface(v)
	default:
		return 0, errors.New(fmt.Sprintf("unsupported type - %d", v.Kind()))
	}
}
