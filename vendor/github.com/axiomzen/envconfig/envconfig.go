// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// ErrInvalidSpecification indicates that a specification is of the wrong type.
var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

// A ParseError occurs when an environment variable cannot be converted to
// the type required by a struct field during assignment.
type ParseError struct {
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
}

// A Decoder is a type that knows how to de-serialize environment variables
// into itself.
type Decoder interface {
	Decode(value string) error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("envconfig.Process: assigning %[1]s to %[2]s: converting '%[3]s' to type %[4]s", e.KeyName, e.FieldName, e.Value, e.TypeName)
}

// Process populates the specified struct based on environment variables
func Process(prefix string, spec interface{}) error {
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr {
		return ErrInvalidSpecification
	}
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return ErrInvalidSpecification
	}
	typeOfSpec := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanSet() || typeOfSpec.Field(i).Tag.Get("ignored") == "true" {
			continue
		}

		if typeOfSpec.Field(i).Anonymous && f.Kind() == reflect.Struct {
			embeddedPtr := f.Addr().Interface()
			if err := Process(prefix, embeddedPtr); err != nil {
				return err
			}
			f.Set(reflect.ValueOf(embeddedPtr).Elem())
		}

		alt := typeOfSpec.Field(i).Tag.Get("envconfig")
		fieldName := typeOfSpec.Field(i).Name
		if alt != "" {
			fieldName = alt
		}
		key := strings.ToUpper(fmt.Sprintf("%s_%s", prefix, fieldName))
		// `os.Getenv` cannot differentiate between an explicitly set empty value
		// and an unset value. `os.LookupEnv` is preferred to `syscall.Getenv`,
		// but it is only available in go1.5 or newer.
		//value, ok := os.LookupEnv(key)
		value, ok := syscall.Getenv(key)
		//fmt.Printf("key: %s, value: %s, ok = %v\n", key, value, ok)
		if !ok && alt != "" {
			key := strings.ToUpper(fieldName)
			value, ok = syscall.Getenv(key)
			//fmt.Printf("key: %s, value: %s, ok = %v\n", key, value, ok)
		}

		def := typeOfSpec.Field(i).Tag.Get("default")
		if def != "" && !ok {
			value = def
		}

		req := typeOfSpec.Field(i).Tag.Get("required")
		if !ok && def == "" {
			if req == "true" {
				return fmt.Errorf("[Process] required key %s missing value", key)
			}
			continue
		}

		err := processField(value, f)
		if err != nil {
			return &ParseError{
				KeyName:   key,
				FieldName: fieldName,
				TypeName:  f.Type().String(),
				Value:     value,
			}
		}

	}
	return nil
}

// MustProcess is the same as Process but panics if an error occurs
func MustProcess(prefix string, spec interface{}) {
	if err := Process(prefix, spec); err != nil {
		panic(err)
	}
}

func processField(value string, field reflect.Value) error {
	typ := field.Type()

	decoder := decoderFrom(field)
	if decoder != nil {
		return decoder.Decode(value)
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}
		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}
		if err != nil {
			return err
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.Slice:
		vals := strings.Split(value, ",")
		sl := reflect.MakeSlice(typ, len(vals), len(vals))
		for i, val := range vals {
			err := processField(val, sl.Index(i))
			if err != nil {
				return err
			}
		}
		field.Set(sl)
	}

	return nil
}

func decoderFrom(field reflect.Value) Decoder {
	if field.CanInterface() {
		dec, ok := field.Interface().(Decoder)
		if ok {
			return dec
		}
	}

	// also check if pointer-to-type implements Decoder,
	// and we can get a pointer to our field
	if field.CanAddr() {
		field = field.Addr()
		dec, ok := field.Interface().(Decoder)
		if ok {
			return dec
		}
	}

	return nil
}

// exportRec recursive helper function
func exportRec(prefix string, spec interface{}, result *[]string, fillDefaults bool) error {
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr {
		return ErrInvalidSpecification
	}

	s = s.Elem()

	if s.Kind() != reflect.Struct {
		return ErrInvalidSpecification
	}

	typeOfSpec := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		if !f.CanSet() || typeOfSpec.Field(i).Tag.Get("ignored") == "true" {
			continue
		}

		if typeOfSpec.Field(i).Anonymous && f.Kind() == reflect.Struct {
			embeddedPtr := f.Addr().Interface()
			// ok to call recursivley with pointer to struct
			if err := exportRec(prefix, embeddedPtr, result, fillDefaults); err != nil {
				return err
			}

			// we don't need to do anything here with structs themselves
			continue
		}

		alt := typeOfSpec.Field(i).Tag.Get("envconfig")
		// the string needs to look like "KEY=value"
		var key string
		if alt != "" {
			key = strings.ToUpper(alt)
		} else {
			key = strings.ToUpper(fmt.Sprintf("%s_%s", prefix, typeOfSpec.Field(i).Name))
		}

		if !isZero(f) {
			// export this field, we set it manually
			// todo: using fmt.Sprintf here to get the string value, see if this is ok

			// need to customize strings as they render as [thing separed by spaces]
			if f.Kind() == reflect.Slice || f.Kind() == reflect.Array {
				r := strings.Replace(fmt.Sprintf("%s=%v", key, f.Interface()), "[", "", -1)
				r = strings.Replace(r, "]", "", -1)
				r = strings.Replace(r, " ", ",", -1)
				//fmt.Println("appending non zero array or slice: " + r)
				*result = append(*result, r)
			} else if f.Kind() == reflect.Ptr {
				if !f.IsNil() {
					elem := f.Elem()
					r := fmt.Sprintf("%s=%v", key, elem.Interface())
					*result = append(*result, r)
				}
			} else {
				r := fmt.Sprintf("%s=%v", key, f.Interface())
				//fmt.Println("appending non zero value: " + r)
				*result = append(*result, r)
			}
			// were done with anything that is not zero'd
			continue
		}

		// check required
		req := typeOfSpec.Field(i).Tag.Get("required")
		//fmt.Printf("required: %s\n", req)

		// check default
		def := typeOfSpec.Field(i).Tag.Get("default")
		if def != "" {
			// write this one
			//fmt.Println("appending default: " + fmt.Sprintf("%s=%s", key, def))
			*result = append(*result, fmt.Sprintf("%s=%s", key, def))
			// set the field if specified
			if fillDefaults {
				//fmt.Printf("Filling default: %s to %#v\n", def, f)
				err := processField(def, f)
				if err != nil {
					return &ParseError{
						KeyName:   key,
						FieldName: typeOfSpec.Field(i).Name,
						TypeName:  f.Type().String(),
						Value:     def,
					}
				}
			}

		} else if req == "true" {
			return fmt.Errorf("[Export] required key %s missing value", key)
		} else {
			// nothing to export ... ?
			fmt.Printf("[Export] warning: key %s will not be exported\n", key)
			//return fmt.Errorf("required key %s missing value", key)
		}
	}
	return nil
}

// isZero inspired from http://stackoverflow.com/questions/23555241/golang-reflection-how-to-get-zero-value-of-a-field-type
// if it is a bool or int or float value, we cannot test for zero'ed as the zero case might be
// a valid value (unfortunatley)
func isZero(v reflect.Value) bool {

	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isZero(v.Index(i))
			if !z {
				return z
			}
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				z = z && isZero(v.Field(i))
			}
			if !z {
				return z
			}
		}
		return z
	case reflect.Ptr:
		// if you use a pointer, we are assuming that it if is nil, it is zero
		// and not zero otherwise
		return v.IsNil()
		//return isZero(reflect.Indirect(v))
	case reflect.Invalid:
		//i.e. uninitialized struct, nil, 0, empty string, etc
		return true
	}

	// Compare other types that support equality operation
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

// Export takes an existing config struct and outputs
// exec.Command.Env friendly env settings
func Export(prefix string, spec interface{}, fillDefaults bool) ([]string, error) {

	var result []string
	err := exportRec(prefix, spec, &result, fillDefaults)
	return result, err
}
