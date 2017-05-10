package compare

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// Going to leave compare here instead of its own package until we get more test cases / usage

// Valueable allows an implementer to return another value to compare to
// useful for comparing nullabletypes to simpler ones
type Valueable interface {
	GetValue() reflect.Value
}

// reuse the same empty struct instance to save memory
var keyExists = struct{}{}

// ignoreMap type def
type ignoreMap map[interface{}]struct{}

// TODO: perhaps we should fix round trip time precision instead?
const defaultTimePrecision = 1000000

// TODO: perhaps we should fix round trip float precision instead?
const defaultEpsilon = 0.000001

// Comparison is the helper struct
type Comparison struct {
	ignoredFields ignoreMap
	timePrecision int
	epsilon       float64
}

// New will create a new comparison
func New() *Comparison {
	return &Comparison{make(ignoreMap), defaultTimePrecision, defaultEpsilon}
}

// WithTimePrecision sets the time precision to use
func (c *Comparison) WithTimePrecision(timePrec int) *Comparison {
	c.timePrecision = timePrec
	return c
}

// WithFloatEpsilon sets the time precision to use
func (c *Comparison) WithFloatEpsilon(fEpsilon float64) *Comparison {
	c.epsilon = fEpsilon
	return c
}

// FailHandler is called when there is a comparison fail, so you know what
// exactly failed
//type FailHandler func(msg string)

// Time will compare two times
func (c *Comparison) Time(t1, t2 time.Time, name string) error {
	return compareTime(t1, t2, name, c.timePrecision)
}

// time will compare two times to the milisecond?
func compareTime(t1, t2 time.Time, name string, timePrec int) error {

	// if t1.Year() != t2.Year() {
	// 	return fmt.Errorf("%s: Years should match: %d, %d", name, t1.Year(), t2.Year())
	// }

	// if t1.Month() != t2.Month() {
	// 	return fmt.Errorf("%s: Month should match: %d, %d", name, t1.Month(), t2.Month())
	// }

	// if t1.Day() != t2.Day() {
	// 	return fmt.Errorf("%s: Days should match: %d, %d", name, t1.Day(), t2.Day())
	// }

	// if t1.Hour() != t2.Hour() {
	// 	return fmt.Errorf("%s: Hours should match: %d, %d", name, t1.Hour(), t2.Hour())
	// }

	// if t1.Minute() != t2.Minute() {
	// 	return fmt.Errorf("%s: Minutes should match: %d, %d", name, t1.Minute(), t2.Minute())
	// }

	// if t1.Second() != t2.Second() {
	// 	return fmt.Errorf("%s: Seconds should match: %d, %d", name, t1.Second(), t2.Second())
	// }

	if t1.Unix() != t2.Unix() {
		return fmt.Errorf("%s: Seconds from epoch should match: %d, %d", name, t1.Unix(), t2.Unix())
	}

	if (t1.Nanosecond() / timePrec) != (t2.Nanosecond() / timePrec) {
		return fmt.Errorf("%s: Nanoseconds should match: %d, %d", name, t1.Nanosecond()/timePrec, t2.Nanosecond()/timePrec)
	}

	return nil
}

// Float32 will compare two floats to see if they are approximatley equal
func (c *Comparison) Float32(f1, f2 float32, name string) error {
	return compareFloat64(float64(f1), float64(f2), name, c.epsilon)
}

//func (c *Compare)Float32()

// Float64 will compare two floats to see if they are approximatley equal
func (c *Comparison) Float64(f1, f2 float64, name string) error {
	return compareFloat64(f1, f2, name, c.epsilon)
}

// compareFloat64 inner function
func compareFloat64(f1, f2 float64, name string, epsilon float64) error {
	if math.Abs(f1-f2) >= epsilon {
		return fmt.Errorf("%s: Floats should match: %f, %f", name, f1, f2)
	}
	return nil
}

// Ignore will ignore the field name given
// will prepend "." in front if you don't have it already
// should support sub fields, like .FieldName.ID for example
// will need to do that even for embedded structs????? yes it looks like it
// so .UserBase.ID
func (c *Comparison) Ignore(field string) *Comparison {
	if !strings.HasPrefix(field, ".") {
		field = "." + field
	}
	c.ignoredFields[field] = keyExists
	return c
}

// IgnoreFields will ignore all these fields
func (c *Comparison) IgnoreFields(fields []string) *Comparison {
	for _, field := range fields {
		c.Ignore(field)
	}
	return c
}

// DeepEquals will compare these two things and return an error
// if they are not equal
func (c *Comparison) DeepEquals(x, y interface{}, info string) error {
	if x == nil || y == nil {
		if x != y {
			return fmt.Errorf("%s: should both be nil or not nil: %t, %t", info, x == nil, y == nil)
		}
	}

	v1 := getValue(reflect.ValueOf(x))
	v2 := getValue(reflect.ValueOf(y))
	// ignoring initial types equal for now, as this makes it more flexible
	// statement := v1.Type() == v2.Type()

	// if !statement {
	// 	return false
	// }

	return sensibleDeepValueEqual(v1, v2, make(map[visit]bool), 0, "", c.timePrecision, c.epsilon, c.ignoredFields)
}

// this will detect if this value implements Valueable and if so return the new value
// otherwise it returns the same value
func getValue(field reflect.Value) reflect.Value {

	//fmt.Printf("GET VALUE CALLED\n")
	if field.IsValid() && field.CanInterface() {
		val, ok := field.Interface().(Valueable)
		if ok {
			//fmt.Printf("RETURNED VALUE[0]: %v\n", val.GetValue())
			return val.GetValue()
		}
	}

	// also check if pointer-to-type implements Valueable,
	// and we can get a pointer to our field
	if field.CanAddr() {
		newField := field.Addr()
		if newField.IsValid() && newField.CanInterface() {
			val, ok := newField.Interface().(Valueable)
			if ok {
				//fmt.Printf("RETURNED VALUE[0]: %v\n", val.GetValue())
				return val.GetValue()
			}
		}
	}

	return field
}

// "Inspired" from golang DeepEquals
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

var hard = func(k reflect.Kind) bool {
	switch k {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		return true
	}
	return false
}

func sensibleDeepValueEqual(v1, v2 reflect.Value, visited map[visit]bool, depth int, name string, timePrec int, epsilon float64, ignoredFields ignoreMap) error {

	// if !v1.IsValid() || !v2.IsValid() {
	// 	if v1.IsValid() != v2.IsValid() {
	// 		//gomega.Ω(statement).Should(gomega.BeTrue(), "%s: should both be valid or not valid: %t, %t", name, v1.IsValid(), v2.IsValid())
	// 		return fmt.Errorf("%s: should both be valid or not valid: %t, %t", name, v1.IsValid(), v2.IsValid())
	// 	} else if !v1.IsValid() && !v2.IsValid() {
	// 		// they are both not valid, we can't call type
	// 		fmt.Printf("BOTH NOT VALID: %v, %v\n", v1, v2)
	// 		return nil
	// 	}
	// }

	//fmt.Printf("COMPARING %s, at depth %d: %v, %v\n", name, depth, v1, v2)

	if depth > 0 {

		// check for ignored fields (by field name)
		if _, has := ignoredFields[name]; has {
			//fmt.Printf("IGNORING %s\n", name)
			return nil
		}

		v1 = getValue(v1)
		v2 = getValue(v2)

		if !v1.IsValid() || !v2.IsValid() {
			if v1.IsValid() != v2.IsValid() {
				return fmt.Errorf("%s: should both be valid or not valid: %t, %t", name, v1.IsValid(), v2.IsValid())
			} else if !v1.IsValid() && !v2.IsValid() {
				// they are both not valid, we can't call type
				//fmt.Printf("BOTH NOT VALID: %v, %v\n", v1, v2)
				return nil
			}
		}

		if v1.Type() != v2.Type() {
			return fmt.Errorf("%s: should both have the same type: %s, %s", name, v1.Type().Name(), v2.Type().Name())
		}
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return nil
		}

		// Remember for later.
		visited[v] = true
	}

	// TODO: if they both can get our interface from them, then
	// do that, and then interface.Compare(v ourinterface){ check type, if ptr, then get value, else compare}
	// or do we just get the reflect.Value from it, then compare

	switch v1.Kind() {
	case reflect.Array:
		// TODO: maybe we don't care about array being equal length
		// {
		// 	lengthEquals := v1.Len() == v2.Len()
		// 	gomega.Ω(lengthEquals).Should(gomega.BeTrue(), "%s: should both have the same length: %d, %d", name, v1.Len(), v2.Len())
		// 	if !lengthEquals {
		// 		return false
		// 	}
		// }
		for i := 0; i < v1.Len(); i++ {

			if err := sensibleDeepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, name, timePrec, epsilon, ignoredFields); err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:

		if v1.IsNil() != v2.IsNil() {
			return fmt.Errorf("%s: should both be nil or not nil: %t, %t", name, v1.IsNil(), v2.IsNil())
		}

		if v1.Len() != v2.Len() {
			return fmt.Errorf("%s: should both have the same length: %d, %d", name, v1.Len(), v2.Len())
		}

		if v1.Pointer() == v2.Pointer() {
			return nil
		}
		for i := 0; i < v1.Len(); i++ {
			if err := sensibleDeepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, name, timePrec, epsilon, ignoredFields); err != nil {
				return err
			}
		}
		return nil
	case reflect.Interface:
		if v1.IsNil() != v2.IsNil() {
			return fmt.Errorf("%s: should both be nil or not nil: %t, %t", name, v1.IsNil(), v2.IsNil())
		}

		return sensibleDeepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, name, timePrec, epsilon, ignoredFields)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return nil
		}
		return sensibleDeepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, name, timePrec, epsilon, ignoredFields)
	case reflect.Struct:
		// check for time.Time type
		if v1.Type() == timeType {
			// use comparetime
			t1 := v1.Interface().(time.Time)
			t2 := v2.Interface().(time.Time)
			return compareTime(t1, t2, name, timePrec)
		}
		v1type := v1.Type()
		for i, n := 0, v1.NumField(); i < n; i++ {
			// get field name
			f := v1type.Field(i)
			//var nn = f.Name
			//if len(name) > 0 {
			//    nn = name + "." + nn
			//}
			if err := sensibleDeepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1, name+"."+f.Name, timePrec, epsilon, ignoredFields); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return fmt.Errorf("%s: should both be nil or not nil: %t, %t", name, v1.IsNil(), v2.IsNil())
		}

		if v1.Len() != v2.Len() {
			return fmt.Errorf("%s: should both have the same length: %d, %d", name, v1.Len(), v2.Len())
		}

		if v1.Pointer() == v2.Pointer() {
			return nil
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)

			if !val1.IsValid() {
				return fmt.Errorf("%s[%v]: val1 should be valid", name, k)
			}

			if !val2.IsValid() {
				return fmt.Errorf("%s[%v]: val2 should be valid", name, k)
			}

			if err := sensibleDeepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1, fmt.Sprintf("%s[%v]", name, k), timePrec, epsilon, ignoredFields); err != nil {
				return err
			}
		}
		return nil
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v1.Float(), v2.Float(), name, epsilon)
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		if v1.Uint() != v2.Uint() {
			return fmt.Errorf("%s: should both have the same value: %d, %d", name, v1.Uint(), v2.Uint())
		}
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		if v1.Int() != v2.Int() {
			return fmt.Errorf("%s: should both have the same value: %d, %d", name, v1.Int(), v2.Int())
		}
	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			return fmt.Errorf("%s: should both have the same value: %t, %t", name, v1.Bool(), v2.Bool())
		}
	case reflect.String:
		if v1.String() != v2.String() {
			return fmt.Errorf("%s: should both have the same value: %s, %s", name, v1.String(), v2.String())
		}
	default:
		return nil
	}
	return nil
}

// Ignorify converts all struct fields into a ".FieldName" list
// will recurse into embedded structs and struct fields
// does not go into arrays or slices of structs
// or nil pointers of structs
func Ignorify(v interface{}) []string {
	// check if we are a struct
	// if not, then return error
	// otherwise, recurse through the struct, and print all fields
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		var values = []string{}
		vn := val.Elem()
		// only take in pointers to structs
		if vn.Kind() == reflect.Struct {
			ignorify(vn, &values, "")
			return values
		}
	case reflect.Struct:
		var values = []string{}
		ignorify(val, &values, "")
		return values
	}
	return []string{}
}

func ignorify(val reflect.Value, fields *[]string, prefix string) {
	for i, n := 0, val.NumField(); i < n; i++ {

		var fieldVal reflect.Value
		// TODO: this won't go into nil pointers of sub-structs
		// also; what about slices or arrays of structs and/or struct pointers?
		if val.Field(i).Kind() == reflect.Ptr {
			fieldVal = val.Field(i).Elem()
			//fmt.Printf("I am a pointer: %v\n", fieldVal)
		} else {
			fieldVal = val.Field(i)
		}
		if fieldVal.Kind() == reflect.Struct {
			ignorify(fieldVal, fields, prefix+"."+val.Type().Field(i).Name)
		} else {
			*fields = append(*fields, prefix+"."+val.Type().Field(i).Name)
		}
	}
}
