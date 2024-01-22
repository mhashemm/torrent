package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func Marshal(v any) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("nil pointer you little bitch")
	}
	buf := bytes.NewBuffer(nil)
	err := marshalStruct(v, buf)
	return buf.Bytes(), err
}

func marshalInt(v int64, buf *bytes.Buffer) {
	buf.WriteByte('i')
	buf.WriteString(strconv.FormatInt(v, 10))
	buf.WriteByte('e')
}

func marshalUint(v uint64, buf *bytes.Buffer) {
	buf.WriteByte('i')
	buf.WriteString(strconv.FormatUint(v, 10))
	buf.WriteByte('e')
}

func marshalString(v string, buf *bytes.Buffer) {
	buf.WriteString(strconv.Itoa(len(v)))
	buf.WriteByte(':')
	buf.WriteString(v)
}

func marshalSlice(v reflect.Value, buf *bytes.Buffer) {
	buf.WriteByte('l')
	for i := 0; i < v.Len(); i++ {
		marshalValue(v.Index(i), buf)
	}
	buf.WriteByte('e')
}

func marshalMap(v reflect.Value, buf *bytes.Buffer) {
	buf.WriteByte('d')
	mi := v.MapRange()
	for mi.Next() {
		marshalString(mi.Key().String(), buf)
		marshalValue(mi.Value(), buf)
	}
	buf.WriteByte('e')
}

func marshalValue(v reflect.Value, buf *bytes.Buffer) {
	switch v.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		marshalInt(v.Int(), buf)
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		marshalUint(v.Uint(), buf)
	case reflect.String:
		marshalString(v.String(), buf)
	case reflect.Slice, reflect.Array:
		marshalSlice(v, buf)
	case reflect.Map:
		marshalMap(v, buf)
	case reflect.Struct:
		marshalStruct(v, buf)
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		marshalValue(v.Elem(), buf)
	}
}

func marshalStruct(v any, buf *bytes.Buffer) error {
	buf.WriteByte('d')
	t, val := reflect.TypeOf(v), reflect.ValueOf(v)
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		name, ok := t.Field(i).Tag.Lookup("bencode")
		if !ok {
			name = t.Field(i).Name
		}
		fmt.Println(name)
		marshalValue(val.Field(i), buf)
		// fmt.Println(fv)
	}
	buf.WriteByte('e')
	return nil
}

func Unmarshal(data []byte, v any) error {
	return nil
}
