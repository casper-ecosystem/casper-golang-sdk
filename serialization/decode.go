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

func MustUnmarshal(src []byte, dest interface{}) {
	err := Unmarshal(src, dest)
	if err != nil {
		panic(err)
	}
}

func Unmarshal(src []byte, dest interface{}) error {
	r := bytes.NewReader(src)
	dec := NewDecoder(r)
	_, err := dec.Decode(dest)
	return err
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

type Decoder struct {
	r io.Reader
}

func (dec *Decoder) Decode(v interface{}) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("can't unmarshal to nil interface")
	}

	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("can't unmarshal to non-pointer '%v' - use & operator", vv.Type().String())
	}
	if vv.IsNil() && !vv.CanSet() {
		return 0, fmt.Errorf("can't unmarshal to unsettable '%v' - use & operator", vv.Type().String())
	}
	return dec.decode(vv.Elem())
}

func (dec *Decoder) DecodeBool() (bool, int, error) {
	val, n, err := dec.DecodeByte()
	if err != nil {
		return false, 0, err
	}

	switch val {
	case 0:
		return false, n, nil
	case 1:
		return true, n, nil
	}

	return false, 0, errors.New("bool not 0 or 1")
}

func (dec *Decoder) DecodeInt32() (int32, int, error) {
	var buf [4]byte
	n, err := io.ReadFull(dec.r, buf[:])

	val := binary.LittleEndian.Uint32(buf[:])
	return int32(val), n, err
}

func (dec *Decoder) DecodeInt64() (int64, int, error) {
	var buf [8]byte
	n, err := io.ReadFull(dec.r, buf[:])

	val := binary.LittleEndian.Uint64(buf[:])
	return int64(val), n, err
}

func (dec *Decoder) DecodeByte() (byte, int, error) {
	var buf [1]byte
	n, err := io.ReadFull(dec.r, buf[:])
	return buf[0], n, err
}

func (dec *Decoder) DecodeUInt32() (uint32, int, error) {
	var buf [4]byte
	n, err := io.ReadFull(dec.r, buf[:])

	val := binary.LittleEndian.Uint32(buf[:])
	return val, n, err
}

func (dec *Decoder) DecodeUInt64() (uint64, int, error) {
	var buf [8]byte
	n, err := io.ReadFull(dec.r, buf[:])

	val := binary.LittleEndian.Uint64(buf[:])
	return val, n, err
}

func (dec *Decoder) DecodeBigNumber(v reflect.Value, expectedLength int) (int, error) {
	var numLen [1]byte
	n, err := io.ReadFull(dec.r, numLen[:])
	if err != nil {
		return n, err
	}

	buf := make([]byte, numLen[0])
	n2, err := io.ReadFull(dec.r, buf)
	n += n2
	if err != nil {
		return n, err
	}

	numBytes := make([]byte, expectedLength)
	for i, b := range buf {
		numBytes[len(numBytes)-i-1] = b
	}

	var val big.Int
	val.SetBytes(numBytes[:])
	v.FieldByName("Int").Set(reflect.ValueOf(val))

	return n, nil
}

func (dec *Decoder) DecodeU128(v reflect.Value) (int, error) {
	return dec.DecodeBigNumber(v, 16)
}

func (dec *Decoder) DecodeU256(v reflect.Value) (int, error) {
	return dec.DecodeBigNumber(v, 32)
}

func (dec *Decoder) DecodeU512(v reflect.Value) (int, error) {
	return dec.DecodeBigNumber(v, 64)
}

func (dec *Decoder) DecodeString() (string, int, error) {
	length, n, err := dec.DecodeUInt32()
	if err != nil {
		return "", n, err
	}
	buf := make([]byte, length)
	n2, err := io.ReadFull(dec.r, buf)
	n += n2

	return string(buf), n, err
}

func (dec *Decoder) DecodeOptional(v reflect.Value) (int, error) {
	present, n, err := dec.DecodeBool()
	if err != nil || !present {
		return n, err
	}

	if err = allocPtrIfNil(&v); err != nil {
		return n, err
	}

	n2, err := dec.decode(v.Elem())
	n += n2

	return n, err
}

func (dec *Decoder) DecodeArray(v reflect.Value) (int, error) {
	length, n, err := dec.DecodeUInt32()
	if err != nil {
		return n, err
	}

	sliceLen := int(length)
	if v.Cap() < sliceLen {
		v.Set(reflect.MakeSlice(v.Type(), sliceLen, sliceLen))
	}
	if v.Len() < sliceLen {
		v.SetLen(sliceLen)
	}

	if v.Type().Elem().Kind() == reflect.Uint8 {
		data, n2, err := dec.DecodeFixedByteArray(sliceLen)
		n += n2
		if err != nil {
			return n, err
		}
		v.SetBytes(data)
		return n, nil
	}

	for i := 0; i < sliceLen; i++ {
		n2, err := dec.decode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (dec *Decoder) DecodeFixedArray(v reflect.Value) (int, error) {
	n := 0
	if v.Type().Elem().Kind() == reflect.Uint8 {
		data, n2, err := dec.DecodeFixedByteArray(v.Len())
		n += n2
		if err != nil {
			return n, err
		}
		v.SetBytes(data)
		return n, nil
	}

	for i := 0; i < v.Len(); i++ {
		n2, err := dec.decode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (dec *Decoder) DecodeFixedByteArray(length int) ([]byte, int, error) {
	buf := make([]byte, length)
	n, err := io.ReadFull(dec.r, buf)
	return buf, n, err
}

func (dec *Decoder) DecodeResult(v reflect.Value) (int, error) {
	val := v.Interface().(Result)
	result := v.FieldByName(val.ResultFieldName())
	isSuccess, n, err := dec.DecodeBool()
	if err != nil {
		return n, err
	}
	result.SetBool(isSuccess)

	var n2 int
	if isSuccess {
		successVal := v.FieldByName(val.SuccessFieldName())
		successVal.Set(reflect.New(successVal.Type().Elem()))
		n2, err = dec.decode(successVal.Elem())
	} else {
		errorVal := v.FieldByName(val.ErrorFieldName())
		errorVal.Set(reflect.New(errorVal.Type().Elem()))
		n2, err = dec.decode(errorVal.Elem())
	}
	n += n2
	return n, err
}

func (dec *Decoder) DecodeTuple(v reflect.Value) (int, error) {
	val := v.Interface().(Tuple)
	fields := val.TupleFields()
	n := 0
	for _, fieldName := range fields {
		field := v.FieldByName(fieldName)
		n2, err := dec.decode(field)
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (dec *Decoder) DecodeMap(v reflect.Value) (int, error) {
	mapLen, n, err := dec.DecodeUInt32()
	if err != nil {
		return n, err
	}

	vt := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(vt))
	}

	keyType := vt.Key()
	elemType := vt.Elem()
	for i := uint32(0); i < mapLen; i++ {
		key := reflect.New(keyType).Elem()
		n2, err := dec.decode(key)
		n += n2
		if err != nil {
			return n, err
		}

		val := reflect.New(elemType).Elem()
		n2, err = dec.decode(val)
		n += n2
		if err != nil {
			return n, err
		}
		v.SetMapIndex(key, val)
	}
	return n, nil
}

func (dec *Decoder) DecodeUnmarshaler(v reflect.Value) (int, error) {
	unmarshaler := v.Interface().(Unmarshaler)
	return unmarshaler.Unmarshal(dec.r)
}

func (dec *Decoder) DecodeStruct(v reflect.Value) (int, error) {
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
		n2, err := dec.decode(vf)
		n += n2
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func (dec *Decoder) DecodeUnion(v reflect.Value) (int, error) {
	union := v.Interface().(Union)

	discriminant, n, err := dec.DecodeByte()
	if err != nil {
		return n, err
	}

	vs := v.FieldByName(union.SwitchFieldName())

	vs.SetUint(uint64(discriminant))

	arm, ok := union.ArmForSwitch(discriminant)

	if !ok {
		return n, fmt.Errorf("switch '%d' is not valid for union", discriminant)
	}

	if arm == "" {
		return n, nil
	}

	vv := v.FieldByName(arm)
	vv.Set(reflect.New(vv.Type().Elem()))
	n2, err := dec.decode(vv.Elem())
	n += n2

	return n, err
}

func (dec *Decoder) decodeInterface(v reflect.Value) (int, error) {
	if v.IsNil() || !v.CanInterface() {
		return 0, fmt.Errorf("can't decode to nil interface")
	}

	// Extract underlying value from the interface and indirect through
	// any pointer, allocating as needed.
	ve := reflect.ValueOf(v.Interface())
	ve, err := dec.indirectIfPtr(ve)
	if err != nil {
		return 0, err
	}
	if !ve.CanSet() {
		return 0, fmt.Errorf("can't decode to unsettable '%v'", ve.Type().String())
	}
	return dec.decode(ve)
}

func (dec *Decoder) indirectIfPtr(v reflect.Value) (reflect.Value, error) {
	if v.Kind() == reflect.Ptr {
		err := allocPtrIfNil(&v)
		return v.Elem(), err
	}
	return v, nil
}

func allocPtrIfNil(v *reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("value is not a pointer: '%v'", v.Type().String())
	}
	isNil := v.IsNil()
	if isNil && !v.CanSet() {
		return fmt.Errorf("unable to allocate pointer for '%v'", v.Type().String())
	}
	if isNil {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return nil
}

func (dec *Decoder) decode(v reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.Bool:
		val, n, err := dec.DecodeBool()
		if err != nil {
			return n, err
		}
		v.SetBool(val)
		return n, nil
	case reflect.Int32:
		val, n, err := dec.DecodeInt32()
		if err != nil {
			return n, err
		}
		v.SetInt(int64(val))
		return n, nil
	case reflect.Int64:
		val, n, err := dec.DecodeInt64()
		if err != nil {
			return n, err
		}
		v.SetInt(val)
		return n, nil
	case reflect.Uint8:
		val, n, err := dec.DecodeByte()
		if err != nil {
			return n, err
		}
		v.SetUint(uint64(val))
		return n, nil
	case reflect.Uint32:
		val, n, err := dec.DecodeUInt32()
		if err != nil {
			return n, err
		}
		v.SetUint(uint64(val))
		return n, nil
	case reflect.Uint64:
		val, n, err := dec.DecodeUInt64()
		if err != nil {
			return n, err
		}
		v.SetUint(val)
		return n, nil
	case reflect.String:
		val, n, err := dec.DecodeString()
		if err != nil {
			return n, err
		}
		v.SetString(val)
		return n, nil
	case reflect.Ptr:
		return dec.DecodeOptional(v)
	case reflect.Slice:
		return dec.DecodeArray(v)
	case reflect.Array:
		return dec.DecodeFixedArray(v)
	case reflect.Map:
		return dec.DecodeMap(v)
	case reflect.Struct:
		if _, ok := v.Interface().(U128); ok {
			return dec.DecodeU128(v)
		}
		if _, ok := v.Interface().(U256); ok {
			return dec.DecodeU256(v)
		}
		if _, ok := v.Interface().(U512); ok {
			return dec.DecodeU512(v)
		}
		if _, ok := v.Interface().(Result); ok {
			return dec.DecodeResult(v)
		}
		if _, ok := v.Interface().(Tuple); ok {
			return dec.DecodeTuple(v)
		}
		if _, ok := v.Interface().(Union); ok {
			return dec.DecodeUnion(v)
		}

		return dec.DecodeStruct(v)
	case reflect.Interface:
		return dec.decodeInterface(v)
	default:
		return 0, errors.New(fmt.Sprintf("unsupported type - %d", v.Kind()))
	}
}
