package compare

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	c := New()
	b1 := true
	b2 := true
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = false
	b2 = false
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = false
	b2 = true
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for bools b1: %v and b2:%v", b1, b2)
	}
}

func TestBoolPointer(t *testing.T) {
	c := New()
	b1 := true
	b2 := true
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = false
	b2 = false
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = false
	b2 = true
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for bools b1: %v and b2:%v", b1, b2)
	}
}

func TestString(t *testing.T) {
	c := New()
	b1 := "true"
	b2 := "true"
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = "false"
	b2 = "true"
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for strings b1: %s and b2: %s", b1, b2)
	}
}

func TestStringPointers(t *testing.T) {
	c := New()
	b1 := "true"
	b2 := "true"
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = "false"
	b2 = "true"
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for strings b1: %s and b2: %s", b1, b2)
	}
}

func TestInt(t *testing.T) {
	c := New()
	b1 := 1
	b2 := 1
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0
	b2 = 0
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 1
	b2 = 0
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for ints b1: %d and b2: %d", b1, b2)
	}
}

func TestIntPointers(t *testing.T) {
	c := New()
	b1 := 1
	b2 := 1
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0
	b2 = 0
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 1
	b2 = 0
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for ints b1: %d and b2: %d", b1, b2)
	}
}

func TestFloat(t *testing.T) {
	c := New()

	b1 := 1.0
	b2 := 1.0
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 1.0
	b2 = 0.0
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}

	b1 = 1.0
	b2 = 1.0 + c.epsilon
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0 + c.epsilon/(1+c.epsilon)
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0 + 2*c.epsilon
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}

	b1 = 1.0
	b2 = 1.0 + 2*c.epsilon
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}
}

func TestFloatPointers(t *testing.T) {
	c := New()

	b1 := 1.0
	b2 := 1.0
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 1.0
	b2 = 0.0
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}

	b1 = 1.0
	b2 = 1.0 + c.epsilon
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0 + c.epsilon/(1+c.epsilon)
	if err := c.DeepEquals(&b1, &b2, ""); err != nil {
		t.Error(err)
	}

	b1 = 0.0
	b2 = 0.0 + 2*c.epsilon
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}

	b1 = 1.0
	b2 = 1.0 + 2*c.epsilon
	if err := c.DeepEquals(&b1, &b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for floats b1: %f and b2: %f", b1, b2)
	}
}

func TestTimeLocation(t *testing.T) {
	// two times should be equal if they are the same instant
	// (ignoring nanoseconds)
	t1 := time.Now()
	t2 := t1.UTC()
	if err := New().DeepEquals(t1, t2, "times with different time zones"); err != nil {
		t.Error(err)
	}
}

func TestTimeEpsilon(t *testing.T) {
	c := New()
	a1 := time.Now().Nanosecond()
	b1 := time.Time{}
	b2 := time.Time{}
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = time.Time{}.Add(time.Duration(a1))
	b2 = time.Time{}.Add(time.Duration(a1))
	if err := c.DeepEquals(b1, b2, ""); err != nil {
		t.Error(err)
	}

	b1 = time.Time{}.Add(time.Duration(a1))
	b2 = time.Time{}.Add(time.Duration(a1 + c.timePrecision))
	if err := c.DeepEquals(b1, b2, ""); err == nil {
		t.Errorf("DeepEquals should have failed for times b1: %v and b2: %v", b1, b2)
	}
}

// TODO: test structs, slices of all of the above, arrays of all of the above, slices and arrays of structs and struct pointers
// sub-structures, etc

type Outer struct {
	InnerStruct Inner
	Name        string
	Other       int
}

type Inner struct {
	Name string
}

func TestIgnoreRootFields(t *testing.T) {
	o1 := Outer{InnerStruct: Inner{Name: "inner1"}, Name: "outer1", Other: 2}
	o2 := Outer{InnerStruct: Inner{Name: "inner2"}, Name: "outer2", Other: 2}
	c := New().IgnoreFields([]string{".InnerStruct", ".Name"})
	//c := New().IgnoreFields([]string{".InnerStruct"})
	if err := c.DeepEquals(o1, o2, "whatever"); err != nil {
		t.Errorf("DeepEquals should have passed for types %v and %v, err: %s", o1, o2, err.Error())
	}
}

type SInner struct {
	Name      string
	Other     int
	SubSinner []SInner
}

// test going into slices
type OuterWithSlice struct {
	InnerFirst  []SInner
	InnerSecond []SInner
}

func TestIngoreSubSliceFields(t *testing.T) {
	sub1 := []SInner{SInner{Name: "a", Other: 10}, SInner{Name: "b", Other: 20}}
	sub2 := []SInner{SInner{Name: "c", Other: 11}, SInner{Name: "d", Other: 21}}

	s11 := []SInner{SInner{Name: "inner111", Other: 1, SubSinner: sub1}, SInner{Name: "inner121", Other: 2}}
	s12 := []SInner{SInner{Name: "inner211", Other: 3}, SInner{Name: "inner221", Other: 4}}

	o1 := OuterWithSlice{InnerFirst: s11, InnerSecond: s12}

	s21 := []SInner{SInner{Name: "inner112", Other: 1, SubSinner: sub2}, SInner{Name: "inner122", Other: 2}}
	o2 := OuterWithSlice{InnerFirst: s21, InnerSecond: s12}

	c := New().IgnoreFields([]string{".InnerFirst.Name", ".InnerFirst.SubSinner.Name", ".InnerFirst.SubSinner.Other"})

	if err := c.DeepEquals(o1, o2, "whatever"); err != nil {
		t.Errorf("DeepEquals should have passed for types %v and %v, err: %s", o1, o2, err.Error())
	}

	c2 := New().IgnoreFields([]string{".InnerFirst.Name", ".InnerFirst.SubSinner.Name"})
	if err := c2.DeepEquals(o1, o2, "whatever"); err == nil {
		t.Errorf("DeepEquals should have failed for types %v and %v", o1, o2)
	}

}

// Ignorify TestString
type Embedded struct {
	Hi string
}

type Sub struct {
	FieldOne string
	FieldTwo *string
}

type Bla struct {
	Embedded
	FieldOne   string
	FieldTwo   *string
	FieldThree Sub
	FieldFour  *Sub
}

func TestIgnorify(t *testing.T) {
	str := "hi"
	sub := Sub{FieldOne: "one_sub", FieldTwo: &str}
	bla := Bla{Embedded: Embedded{Hi: "hi"}, FieldOne: "one_bla", FieldTwo: &str, FieldThree: sub, FieldFour: &sub}

	expected := []string{
		".Embedded.Hi",
		".FieldOne",
		".FieldTwo",
		".FieldThree.FieldOne",
		".FieldThree.FieldTwo",
		".FieldFour.FieldOne",
		".FieldFour.FieldTwo"}

	result := Ignorify(bla)
	for i, s := range expected {
		if s != result[i] {
			t.Errorf("Expected %s, got %s", s, result[i])
		}
	}

	result = Ignorify(&bla)
	for i, s := range expected {
		if s != result[i] {
			t.Errorf("Expected %s, got %s", s, result[i])
		}
	}
}
