package lorem

import (
	"strings"
	"testing"

	"github.com/twinj/uuid"
)

type SimpleStruct struct {
	Int8               int8
	Int16              int16
	Int32              int32
	Int64              int64
	Int                int
	UInt8              uint8
	UInt16             uint16
	UInt32             uint32
	UInt64             uint64
	UInt               uint
	Float32            float32
	Float64            float64
	PlainString        string
	Word               string `lorem:"word"`
	WordWithRange      string `lorem:"word,10,11"`
	Sentence           string `lorem:"sentence"`
	SentenceWithRange  string `lorem:"sentence,10,11"`
	Paragraph          string `lorem:"paragraph"`
	ParagraphWithRange string `lorem:"paragraph,10,11"`
	URL                string `lorem:"url"`
	ReadablePath       string `lorem:"readablepath"`
	Host               string `lorem:"host"`
	Email              string `lorem:"email"`
	UUID               string `lorem:"uuid"`
	Bool               bool
}

func TestSimpleStruct(t *testing.T) {
	var ss SimpleStruct

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	// check everything

	if ss.PlainString == "" {
		t.Errorf("PlainString: expected string not empty, got %s", ss.PlainString)
	}

	if ss.Word == "" {
		t.Errorf("Word: expected string not empty, got %s", ss.Word)
	}

	if len(ss.WordWithRange) < 10 || len(ss.WordWithRange) > 11 {
		t.Errorf("WordWithRange: expected 9 < len(len(ss.WordWithRange)) < 12, got %d", len(ss.WordWithRange))
	}

	if ss.Sentence == "" {
		t.Errorf("Sentence: expected string not empty, got %s", ss.Sentence)
	}
	// try for a certain number of periods?
	numspaces := strings.Count(ss.SentenceWithRange, " ")
	if numspaces < 9 || numspaces > 10 {
		t.Errorf("SentenceWithRange: expected 8 < strings.Count(ss.SentenceWithRange, \".\") < 11, got %d", numspaces)
	}

	if ss.Paragraph == "" {
		t.Errorf("Paragraph: expected string not empty, got %s", ss.Paragraph)
	}

	numperiodSpaces := strings.Count(ss.ParagraphWithRange, ". ")
	if numperiodSpaces < 9 || numperiodSpaces > 10 {
		t.Errorf("ParagraphWithRange: expected 8 < strings.Count(ss.ParagraphWithRange, \".\") < 11, got %d", numperiodSpaces)
	}

	if !strings.HasPrefix(ss.URL, "http://www.") {
		t.Errorf("URL: expected url to start with http://www., got %s", ss.URL)
	}

	// not really a url, just a readable path
	if strings.Contains(ss.ReadablePath, " ") || strings.Contains(ss.ReadablePath, ",") || strings.Contains(ss.ReadablePath, ".") {
		t.Errorf("ReadablePath: expected readable url to not contain invalid characters, got %s", ss.ReadablePath)
	}

	if !strings.Contains(ss.Host, ".") {
		t.Errorf("Host: expected host to contain '.', got %s", ss.Host)
	}

	if !strings.Contains(ss.Email, ".") || !strings.Contains(ss.Email, "@") {
		t.Errorf("Email: expected Email to contain '.' and '@', got %s", ss.Email)
	}

	if _, err := uuid.Parse(ss.UUID); err != nil {
		t.Errorf("Email: no error parsing uuid, got %s", err.Error())
	}

}

type StructWithPointers struct {
	Int8Pointer               *int8
	Int16Pointer              *int16
	Int32Pointer              *int32
	Int64Pointer              *int64
	IntPointer                *int
	UInt8Pointer              *uint8
	UInt16Pointer             *uint16
	UInt32Pointer             *uint32
	UInt64Pointer             *uint64
	UIntPointer               *uint
	Float32Pointer            *float32
	Float64Pointer            *float64
	PlainStringPointer        *string
	WordPointer               *string `lorem:"word"`
	WordWithRangePointer      *string `lorem:"word,10,11"`
	SentencePointer           *string `lorem:"sentence"`
	SentenceWithRangePointer  *string `lorem:"sentence,10,11"`
	ParagraphPointer          *string `lorem:"paragraph"`
	ParagraphWithRangePointer *string `lorem:"paragraph,10,11"`
	URLPointer                *string `lorem:"url"`
	ReadablePathPointer       *string `lorem:"readablepath"`
	HostPointer               *string `lorem:"host"`
	EmailPointer              *string `lorem:"email"`
	BoolPointer               *bool
}

func TestStructWithPointers(t *testing.T) {
	var ss StructWithPointers

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	// test to see that they are all not nil
	if ss.Int8Pointer == nil {
		t.Errorf("Int8Pointer: expected pointer to not be nil")
	}

	if ss.Int16Pointer == nil {
		t.Errorf("Int16Pointer: expected pointer to not be nil")
	}

	if ss.Int32Pointer == nil {
		t.Errorf("Int32Pointer: expected pointer to not be nil")
	}
	if ss.Int64Pointer == nil {
		t.Errorf("Int64Pointer: expected pointer to not be nil")
	}
	if ss.IntPointer == nil {
		t.Errorf("IntPointer: expected pointer to not be nil")
	}
	if ss.UInt8Pointer == nil {
		t.Errorf("UInt8Pointer: expected pointer to not be nil")
	}
	if ss.UInt16Pointer == nil {
		t.Errorf("UInt16Pointer: expected pointer to not be nil")
	}
	if ss.UInt32Pointer == nil {
		t.Errorf("UInt32Pointer: expected pointer to not be nil")
	}
	if ss.UInt64Pointer == nil {
		t.Errorf("UInt64Pointer: expected pointer to not be nil")
	}
	if ss.UIntPointer == nil {
		t.Errorf("UInt64Pointer: expected pointer to not be nil")
	}
	if ss.Float32Pointer == nil {
		t.Errorf("Float32Pointer: expected pointer to not be nil")
	}
	if ss.Float64Pointer == nil {
		t.Errorf("Float64Pointer: expected pointer to not be nil")
	}
	if ss.PlainStringPointer == nil || *ss.PlainStringPointer == "" {
		t.Errorf("PlainStringPointer: expected pointer to not be nil and not empty")
	}
	if ss.BoolPointer == nil {
		t.Errorf("BoolPointer: expected pointer to not be nil")
	}

	if *ss.WordPointer == "" {
		t.Errorf("WordPointer: expected string not empty, got %s", *ss.WordPointer)
	}

	if len(*ss.WordWithRangePointer) < 2 || len(*ss.WordWithRangePointer) > 11 {
		t.Errorf("WordWithRangePointer: expected 2 < len(len(ss.WordWithRangePointer)) < 12, got %d", len(*ss.WordWithRangePointer))
	}

	if *ss.SentencePointer == "" {
		t.Errorf("Sentence: expected string not empty, got %s", *ss.SentencePointer)
	}
	// try for a certain number of periods?
	numspaces := strings.Count(*ss.SentenceWithRangePointer, " ")
	if numspaces < 9 || numspaces > 10 {
		t.Errorf("SentenceWithRangePointer: expected 8 < strings.Count(*ss.SentenceWithRangePointer, \".\") < 11, got %d", numspaces)
	}

	if *ss.ParagraphPointer == "" {
		t.Errorf("ParagraphPointer: expected string not empty, got %s", *ss.ParagraphPointer)
	}

	numperiodSpaces := strings.Count(*ss.ParagraphWithRangePointer, ". ")
	if numperiodSpaces < 9 || numperiodSpaces > 10 {
		t.Errorf("ParagraphWithRangePointer: expected 8 < strings.Count(*ss.ParagraphWithRangePointer, \".\") < 11, got %d", numperiodSpaces)
	}

	if !strings.HasPrefix(*ss.URLPointer, "http://www.") {
		t.Errorf("URLPointer: expected url to start with http://www., got %s", *ss.URLPointer)
	}

	// not really a url, just a readable path
	if strings.Contains(*ss.ReadablePathPointer, " ") || strings.Contains(*ss.ReadablePathPointer, ",") || strings.Contains(*ss.ReadablePathPointer, ".") {
		t.Errorf("ReadablePathPointer: expected readable url to not contain invalid characters, got %s", *ss.ReadablePathPointer)
	}

	if !strings.Contains(*ss.HostPointer, ".") {
		t.Errorf("HostPointer: expected host to contain '.', got %s", *ss.HostPointer)
	}

	if !strings.Contains(*ss.EmailPointer, ".") || !strings.Contains(*ss.EmailPointer, "@") {
		t.Errorf("EmailPointer: expected Email to contain '.' and '@', got %s", *ss.EmailPointer)
	}
}

type OtherStruct struct {
	SubEmailPointer  *string `lorem:"email"`
	SubWordWithRange string  `lorem:"word,10,11"`
}

type StructWithSlices struct {
	// try a slice of things
	// tag applies to each thing, rather than the slice as a whole
	Words []string `lorem:"word"`
}

func TestStructWithSlices(t *testing.T) {
	var ss StructWithSlices

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if ss.Words == nil {
		t.Errorf("Words: expected Words not be nil")
	}

	if len(ss.Words) > 0 {
		if len(ss.Words[0]) == 0 {
			t.Errorf("Words: expected words[0] to be longer than 0")
		}
	}
}

type StructWithSlicesOfPointers struct {
	// try a slice of pointers to things
	Sentences []*string `lorem:"sentence,10,11"`
}

func TestStructWithSlicesOfPointers(t *testing.T) {
	var ss StructWithSlicesOfPointers

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if ss.Sentences == nil {
		t.Errorf("Sentences: expected Sentences not be nil")
	}

	for _, sp := range ss.Sentences {
		if sp == nil {
			t.Errorf("Sentences: expected sp not be nil")
		}

		if len(*sp) == 0 {
			t.Errorf("Sentences: expected *sp to be longer than 0")
		}

		numspaces := strings.Count(*sp, " ")
		if numspaces < 9 || numspaces > 10 {
			t.Errorf("Sentences: expected 8 < strings.Count(*ss.Sentences[i], \".\") < 11, got %d", numspaces)
		}
	}
}

// try a map
type StructWithMap struct {

	// we don't do maps at this time
	Map        map[string]string
	MapWithTag map[string]string `lorem:"word"`
}

func TestStructWithMap(t *testing.T) {
	var ss StructWithMap

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if ss.Map != nil {
		t.Errorf("Map: expected %v to be nil", ss.Map)
	}

	if ss.MapWithTag != nil {
		t.Errorf("MapWithTag: expected %v to be nil", ss.MapWithTag)
	}
}

type StructWithStruct struct {
	OtherStruct        OtherStruct
	OtherStructPointer *OtherStruct
}

func TestStructWithStruct(t *testing.T) {
	var ss StructWithStruct

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	length := len(ss.OtherStruct.SubWordWithRange)
	if length < 10 || length > 11 {
		t.Errorf("SubWordWithRange: expected 9 < len(ss.OtherStruct.SubWordWithRange) < 12, got %d", length)
	}

	if ss.OtherStruct.SubEmailPointer == nil {
		t.Error("SubEmailPointer: expected it to not be nil")
	}

	if !strings.Contains(*ss.OtherStruct.SubEmailPointer, ".") || !strings.Contains(*ss.OtherStruct.SubEmailPointer, "@") {
		t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *ss.OtherStruct.SubEmailPointer)
	}

	if ss.OtherStructPointer == nil {
		t.Errorf("OtherStructPointer: expected OtherStruct to not be nil")
	}

	length = len(ss.OtherStructPointer.SubWordWithRange)
	if length < 10 || length > 11 {
		t.Errorf("SubWordWithRange: expected 9 < len(ss.OtherStructPointer.SubWordWithRange) < 12, got %d", length)
	}

	if ss.OtherStructPointer.SubEmailPointer == nil {
		t.Error("SubEmailPointer: expected it to not be nil")
	}

	if !strings.Contains(*ss.OtherStructPointer.SubEmailPointer, ".") || !strings.Contains(*ss.OtherStructPointer.SubEmailPointer, "@") {
		t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *ss.OtherStructPointer.SubEmailPointer)
	}
}

type StructWithEmbeddedStruct struct {
	OtherStruct
	URLPointer *string `lorem:"url"`
}

func TestStructWithEmbeddedStruct(t *testing.T) {
	var ss StructWithEmbeddedStruct

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	length := len(ss.SubWordWithRange)
	if length < 10 || length > 11 {
		t.Errorf("SubWordWithRange: expected 9 < len(ss.SubWordWithRange) < 12, got %d", length)
	}

	if ss.SubEmailPointer == nil {
		t.Error("SubEmailPointer: expected it to not be nil")
	}

	if !strings.Contains(*ss.SubEmailPointer, ".") || !strings.Contains(*ss.SubEmailPointer, "@") {
		t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *ss.SubEmailPointer)
	}

	if !strings.HasPrefix(*ss.URLPointer, "http://www.") {
		t.Errorf("URLPointer: expected url to start with http://www., got %s", *ss.URLPointer)
	}

}

type StructWithEmbeddedStructPointer struct {
	*OtherStruct
	Word string `lorem:"word"`
}

func TestStructWithEmbeddedStructPointer(t *testing.T) {
	var ss StructWithEmbeddedStructPointer

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	length := len(ss.SubWordWithRange)
	if length < 10 || length > 11 {
		t.Errorf("SubWordWithRange: expected 9 < len(ss.SubWordWithRange) < 12, got %d", length)
	}

	if ss.SubEmailPointer == nil {
		t.Error("SubEmailPointer: expected it to not be nil")
	}

	if !strings.Contains(*ss.SubEmailPointer, ".") || !strings.Contains(*ss.SubEmailPointer, "@") {
		t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *ss.SubEmailPointer)
	}

	if ss.Word == "" {
		t.Errorf("Word: expected string not empty, got %s", ss.Word)
	}
}

type StructWithSliceOfStructs struct {
	OtherStructs        []OtherStruct
	OtherStructPointers []*OtherStruct
}

func TestStructWithSliceOfStructs(t *testing.T) {
	var ss StructWithSliceOfStructs

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if len(ss.OtherStructs) < 1 {
		t.Errorf("OtherStructs: expected a non empty slice")
	}

	for _, b := range ss.OtherStructs {
		length := len(b.SubWordWithRange)
		if length < 10 || length > 11 {
			t.Errorf("SubWordWithRange: expected 9 < len(ss.SubWordWithRange) < 12, got %d", length)
		}

		if b.SubEmailPointer == nil {
			t.Error("SubEmailPointer: expected it to not be nil")
		}

		if !strings.Contains(*b.SubEmailPointer, ".") || !strings.Contains(*b.SubEmailPointer, "@") {
			t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *b.SubEmailPointer)
		}
	}

	if len(ss.OtherStructPointers) < 1 {
		t.Errorf("OtherStructPointers: expected a non empty slice")
	}

	for _, b := range ss.OtherStructPointers {
		if b == nil {
			t.Errorf("OtherStructPointers: expected the element to not be nil")
		}

		length := len(b.SubWordWithRange)
		if length < 10 || length > 11 {
			t.Errorf("SubWordWithRange: expected 9 < len(ss.SubWordWithRange) < 12, got %d", length)
		}

		if b.SubEmailPointer == nil {
			t.Error("SubEmailPointer: expected it to not be nil")
		}

		if !strings.Contains(*b.SubEmailPointer, ".") || !strings.Contains(*b.SubEmailPointer, "@") {
			t.Errorf("OtherStruct: expected SubEmailPointer to contain '.' and '@', got %s", *b.SubEmailPointer)
		}
	}

}

type StructWithIgnoredFields struct {
	IgnoredInt                   int            `lorem:"-"`
	IgnoredUInt                  uint           `lorem:"-"`
	IgnoredFloat32               float32        `lorem:"-"`
	IgnoredFloat64               float64        `lorem:"-"`
	IgnoredString                string         `lorem:"-"`
	IgnoredBool                  bool           `lorem:"-"`
	IgnoredIntPointer            *int           `lorem:"-"`
	IgnoredUIntPointer           *uint          `lorem:"-"`
	IgnoredStringPointer         *string        `lorem:"-"`
	IgnoredBoolPointer           *bool          `lorem:"-"`
	IgnoredStruct                OtherStruct    `lorem:"-"`
	IgnoredStructPointer         *OtherStruct   `lorem:"-"`
	IgnoredSlice                 []string       `lorem:"-"`
	IgnoredSliceOfStructs        []OtherStruct  `lorem:"-"`
	IgnoredSliceOfStructPointers []*OtherStruct `lorem:"-"`
}

func TestStructWithIgnoredFields(t *testing.T) {
	var ss StructWithIgnoredFields

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	// use reflection to see if it is the zero value?
	// bah, for now just use known 0's

	if ss.IgnoredInt != 0 {
		t.Error("IgnoredInt: Expected to equal zero value")
	}
	if ss.IgnoredUInt != 0 {
		t.Error("IgnoredUInt: Expected to equal zero value")
	}
	if ss.IgnoredFloat32 != 0 {
		t.Error("IgnoredFloat32: Expected to equal zero value")
	}
	if ss.IgnoredFloat64 != 0 {
		t.Error("IgnoredFloat64: Expected to equal zero value")
	}
	if ss.IgnoredString != "" {
		t.Error("IgnoredString: Expected to equal empty string")
	}
	if ss.IgnoredBool != false {
		t.Error("IgnoredBool: Expected to equal false")
	}
	if ss.IgnoredIntPointer != nil {
		t.Error("IgnoredIntPointer: Expected to equal nil")
	}
	if ss.IgnoredUIntPointer != nil {
		t.Error("IgnoredUIntPointer: Expected to equal nil")
	}
	if ss.IgnoredStringPointer != nil {
		t.Error("IgnoredStringPointer: Expected to equal nil")
	}
	if ss.IgnoredBoolPointer != nil {
		t.Error("IgnoredBoolPointer: Expected to equal nil")
	}
	if ss.IgnoredStruct.SubEmailPointer != nil {
		t.Error("IgnoredStruct: Expected to equal nil")
	}
	if ss.IgnoredStruct.SubWordWithRange != "" {
		t.Error("IgnoredStruct: Expected to equal empty string")
	}
	if ss.IgnoredStructPointer != nil {
		t.Error("IgnoredStructPointer: Expected to equal nil")
	}
	if ss.IgnoredSlice != nil {
		t.Error("IgnoredSlice: Expected to equal nil")
	}
	if ss.IgnoredSliceOfStructs != nil {
		t.Error("IgnoredSliceOfStructs: Expected to equal nil")
	}
	if ss.IgnoredSliceOfStructPointers != nil {
		t.Error("IgnoredSliceOfStructPointers: Expected to equal nil")
	}
}

type StructWithIgnoredEmbeddedStruct struct {
	OtherStruct   `lorem:"-"`
	WordWithRange string `lorem:"word,10,11"`
}

func TestStructWithIgnoredEmbeddedStructs(t *testing.T) {
	var ss StructWithIgnoredEmbeddedStruct

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if ss.SubEmailPointer != nil {
		t.Error("SubEmailPointer: Expected to equal nil")
	}
	if ss.SubWordWithRange != "" {
		t.Error("SubWordWithRange: Expected to equal empty string")
	}

	if len(ss.WordWithRange) < 10 || len(ss.WordWithRange) > 11 {
		t.Errorf("WordWithRange: expected 9 < len(len(ss.WordWithRange)) < 12, got %d", len(ss.WordWithRange))
	}

}

type StructWithIgnoredEmbeddedStructPointer struct {
	*OtherStruct         `lorem:"-"`
	WordWithRangePointer *string `lorem:"word,10,11"`
}

func TestStructWithIgnoredEmbeddedStructPointers(t *testing.T) {
	var ss StructWithIgnoredEmbeddedStructPointer

	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if ss.OtherStruct != nil {
		t.Error("OtherStruct: expected it to be nil")
	}

	if len(*ss.WordWithRangePointer) < 10 || len(*ss.WordWithRangePointer) > 11 {
		t.Errorf("WordWithRangePointer: expected 9 < len(ss.WordWithRangePointer) < 12, got %d", len(*ss.WordWithRangePointer))
	}
}

type StructWithFieldThatImplementsDecode struct {
	Sub SubStructLikeWord `lorem:"word,10,11"`
}

type SubStructLikeWord struct {
	word string
}

func (s *SubStructLikeWord) LoremDecode(tag, example string) error {
	s.word = example
	return nil
}

func TestStructWithStructFieldThatImplementsDecode(t *testing.T) {
	var ss StructWithFieldThatImplementsDecode
	if err := Fill(&ss); err != nil {
		t.Error(err.Error())
	}

	if len(ss.Sub.word) < 10 || len(ss.Sub.word) > 11 {
		t.Errorf("Sub.Word: expected 9 < len(ss.Sub.Word) < 12, got %d", len(ss.Sub.word))
	}
}
