package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type scanner struct {
	data []byte
	i    int
}

func (s *scanner) current() byte {
	if s.i < len(s.data) {
		return s.data[s.i]
	}
	return 0
}

func (s *scanner) next() byte {
	s.i += 1
	return s.current()
}

func (s *scanner) skipTo(c byte) {
	for s.i < len(s.data) && s.data[s.i] != c {
		s.i += 1
	}
}

func (s *scanner) skipString() {
	start := s.i
	s.skipTo(':')
	end := s.i
	s.next()
	length, _ := strconv.ParseInt(string(s.data[start:end]), 10, 64)
	s.i += int(length)
}

func (s *scanner) skipInt() {
	s.skipTo('e')
	s.next()
}

func (s *scanner) skip() {
	switch s.current() {
	case 'd', 'l':
		s.next()
		for s.current() != 'e' {
			s.skip()
		}
		s.next()
	case 'i':
		s.skipInt()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.skipString()
	}
}

func (s *scanner) int64(rv reflect.Value) error {
	if !rv.CanSet() {
		return fmt.Errorf("can not set int64 in %s", rv.Type().Name())
	}
	if s.current() != 'i' {
		return fmt.Errorf("invalid int64 at %d", s.i)
	}
	start := s.i + 1
	s.skipTo('e')
	end := s.i
	s.next()
	value, err := strconv.ParseInt(string(s.data[start:end]), 10, 64)
	if err != nil {
		return err
	}
	rv.SetInt(value)
	return nil
}

func (s *scanner) string(rv reflect.Value) error {
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.CanSet() {
		return fmt.Errorf("can not set string in %s", rv.Type().Name())
	}
	if s.current() < '0' || s.current() > '9' {
		return fmt.Errorf("invalid string format at %d", s.i)
	}
	start := s.i
	s.skipTo(':')
	end := s.i
	s.next()
	length, err := strconv.ParseInt(string(s.data[start:end]), 10, 64)
	if err != nil {
		return err
	}
	start, end = s.i, s.i+int(length)
	s.i += int(length)
	if len(s.data) < end {
		return fmt.Errorf("invalid string length %d against %d", len(s.data), length)
	}
	str := string(s.data[start:end])
	rv.SetString(str)
	return nil
}

func (s *scanner) list(rv reflect.Value) error {
	if s.current() != 'l' {
		return fmt.Errorf("invalid list at %d", s.i)
	}
	s.next()
	i := 0
	defer func() {
		if rv.Len() > i+1 {
			rv.SetLen(i + 1)
			rv.SetCap(i + 1)
		}
	}()
	for s.current() != 'e' {
		if rv.Len() < i+1 {
			rv.Grow((i + 1) * 2)
			rv.SetLen((i + 1) * 2)
		}
		e := rv.Index(i)
		if e.Kind() == reflect.Pointer && e.IsNil() {
			e.Set(reflect.New(e.Type().Elem()))
		}
		err := s.unmarshal(e)
		if err != nil {
			return err
		}
		i += 1
	}
	s.next()
	return nil
}

func findField(rv reflect.Value, key string) reflect.Value {
	rt := reflect.TypeOf(rv.Interface())
	for i := range rt.NumField() {
		f := rt.Field(i)
		tag, _ := f.Tag.Lookup("bencode")
		name, _, _ := strings.Cut(tag, ",")
		if (f.Name == key && f.IsExported() && name != "-") || name == key {
			f := rv.Field(i)
			if f.Kind() == reflect.Pointer && f.IsNil() {
				f.Set(reflect.New(f.Type().Elem()))
			}
			return f
		}
	}
	return reflect.Value{}
}

func (s *scanner) dictionary(rv reflect.Value) error {
	if s.current() != 'd' {
		return fmt.Errorf("invalid dictionary at %d", s.i)
	}
	s.next()
	for s.current() != 'e' {
		key := ""
		err := s.string(reflect.ValueOf(&key))
		if err != nil {
			return err
		}
		fv := findField(rv, key)
		if !fv.IsValid() {
			s.skip()
			continue
		}
		err = s.unmarshal(fv)
		if err != nil {
			return err
		}
	}
	s.next()
	return nil
}

func (s *scanner) unmarshal(rv reflect.Value) error {
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	var err error
	switch rv.Kind() {
	case reflect.String:
		err = s.string(rv)
	case reflect.Int64:
		err = s.int64(rv)
	case reflect.Slice:
		err = s.list(rv)
	case reflect.Struct:
		err = s.dictionary(rv)
	default:
		return fmt.Errorf("unsupported type %s", rv.Kind().String())
	}
	return err
}

func Unmarshal(data []byte, v any) error {
	if len(data) == 0 {
		return nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("invalid type: %s", reflect.TypeOf(v).String())
	}
	s := scanner{data: data, i: 0}
	return s.unmarshal(rv)
}
