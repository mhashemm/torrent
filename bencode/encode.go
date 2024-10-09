package bencode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func marshalByteSlice(v []byte, buf *bytes.Buffer) {
	count := base64.StdEncoding.EncodedLen(len(v))
	b := bytes.NewBuffer(make([]byte, 0, count))
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(v)
	e.Close()
	buf.WriteString(strconv.Itoa(count))
	buf.WriteByte(':')
	buf.Write(b.Bytes())
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
		if v, ok := v.Interface().([]byte); ok {
			marshalByteSlice(v, buf)
			return
		}
		marshalSlice(v, buf)
	case reflect.Map:
		marshalMap(v, buf)
	case reflect.Struct:
		marshalStruct(v.Interface(), buf)
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		marshalValue(v.Elem(), buf)
	}
}

func marshalStruct(v any, buf *bytes.Buffer) error {
	buf.WriteByte('d')
	typ, val := reflect.TypeOf(v), reflect.ValueOf(v)
	for i := 0; i < typ.NumField(); i++ {
		if !typ.Field(i).IsExported() || (val.Field(i).Kind() == reflect.Pointer && val.Field(i).IsNil()) {
			continue
		}
		tag, _ := typ.Field(i).Tag.Lookup("bencode")
		if strings.Contains(tag, "omitempty") {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = typ.Field(i).Name
		}
		fv := val.Field(i)
		marshalString(name, buf)
		marshalValue(fv, buf)
	}
	buf.WriteByte('e')
	return nil
}
