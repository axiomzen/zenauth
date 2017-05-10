// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	//"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type Specification struct {
	Embedded
	EmbeddedButIgnored           `ignored:"true"`
	Debug                        bool //5
	Port                         int
	Rate                         float32
	User                         string
	TTL                          uint32
	Timeout                      time.Duration //10
	AdminUsers                   []string
	MagicNumbers                 []int
	MultiWordVar                 string
	SomePointer                  *string
	SomePointerWithDefault       *string       `default:"foo2baz"` //15
	MultiWordVarWithAlt          string        `envconfig:"MULTI_WORD_VAR_WITH_ALT"`
	MultiWordVarWithLowerCaseAlt string        `envconfig:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string        `envconfig:"SERVICE_HOST"`
	DefaultVar                   string        `default:"foobar"`
	RequiredVar                  string        `required:"true"` //20
	NoPrefixDefault              string        `envconfig:"BROKER" default:"127.0.0.1"`
	RequiredDefault              string        `required:"true" default:"foo2bar"` //22
	Ignored                      string        `ignored:"true"`
	NotRequiredVar               string        `required:"false"`
	TimeoutWithDefault           time.Duration `default:"32h"`
	DefaultInt                   int           `default:"7"`
	DefaultBoolPointer           *bool         `default:"true"`
	DefaultBoolPointerOveridden  *bool         `default:"true"`
}

type Embedded struct {
	Enabled             bool
	EmbeddedPort        int
	MultiWordVar        string
	MultiWordVarWithAlt string `envconfig:"MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `envconfig:"EMBEDDED_WITH_ALT"`
	EmbeddedIgnored     string `ignored:"true"`
}

type EmbeddedButIgnored struct {
	FirstEmbeddedButIgnored  string
	SecondEmbeddedButIgnored string
}

func TestProcess(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMINUSERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGICNUMBERS", "5,10,20")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	os.Setenv("ENV_CONFIG_IGNORED", "was-not-ignored")

	defer os.Unsetenv("ENV_CONFIG_DEBUG")
	defer os.Unsetenv("ENV_CONFIG_PORT")
	defer os.Unsetenv("ENV_CONFIG_RATE")
	defer os.Unsetenv("ENV_CONFIG_USER")
	defer os.Unsetenv("ENV_CONFIG_TIMEOUT")
	defer os.Unsetenv("ENV_CONFIG_ADMINUSERS")
	defer os.Unsetenv("ENV_CONFIG_MAGICNUMBERS")
	defer os.Unsetenv("SERVICE_HOST")
	defer os.Unsetenv("ENV_CONFIG_TTL")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	defer os.Unsetenv("ENV_CONFIG_IGNORED")

	err := Process("env_config", &s)
	if err != nil {
		t.Error(err.Error())
	}
	if s.NoPrefixWithAlt != "127.0.0.1" {
		t.Errorf("expected %v, got %v", "127.0.0.1", s.NoPrefixWithAlt)
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.TTL != 30 {
		t.Errorf("expected %d, got %v", 30, s.TTL)
	}
	if s.User != "Kelsey" {
		t.Errorf("expected %s, got %s", "Kelsey", s.User)
	}
	if s.Timeout != 2*time.Minute {
		t.Errorf("expected %s, got %s", 2*time.Minute, s.Timeout)
	}
	if s.RequiredVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.RequiredVar)
	}
	if len(s.AdminUsers) != 3 ||
		s.AdminUsers[0] != "John" ||
		s.AdminUsers[1] != "Adam" ||
		s.AdminUsers[2] != "Will" {
		t.Errorf("expected %#v, got %#v", []string{"John", "Adam", "Will"}, s.AdminUsers)
	}
	if len(s.MagicNumbers) != 3 ||
		s.MagicNumbers[0] != 5 ||
		s.MagicNumbers[1] != 10 ||
		s.MagicNumbers[2] != 20 {
		t.Errorf("expected %#v, got %#v", []int{5, 10, 20}, s.MagicNumbers)
	}
	if s.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}

	if s.DefaultInt != 7 {
		t.Errorf("expected %d, got %d", 7, s.DefaultInt)
	}

	if *s.DefaultBoolPointer != true {
		t.Errorf("expected %v, got %v", true, *s.DefaultBoolPointer)
	}

	if *s.DefaultBoolPointerOveridden != true {
		t.Errorf("expected %v, got %v", true, *s.DefaultBoolPointerOveridden)
	}

}

func TestExport(t *testing.T) {
	// TODO: brittle test
	var s Specification
	os.Clearenv()
	s.Embedded = Embedded{}
	s.Embedded.Enabled = true
	s.Embedded.EmbeddedPort = 5000
	s.Embedded.MultiWordVar = "fooembedded"
	s.Embedded.MultiWordVarWithAlt = "bazembedded"
	s.Embedded.EmbeddedAlt = "embeddedalt"

	s.Debug = true
	s.Port = 8080
	s.Rate = 0.5
	s.User = "Kelsey"
	s.TTL = 30
	s.Timeout = 2 * time.Minute
	s.AdminUsers = []string{"John", "Adam", "Will"}
	s.MagicNumbers = []int{5, 10, 20}
	s.MultiWordVar = "foo bar"
	s.SomePointer = &s.MultiWordVar
	s.MultiWordVarWithAlt = "ALT"
	s.MultiWordVarWithLowerCaseAlt = "lower_casE"
	s.NoPrefixWithAlt = "127.0.0.1"
	s.RequiredVar = "foo"
	s.Ignored = "was-not-ignored"
	f := false
	//s.DefaultBoolPointer = nil
	s.DefaultBoolPointerOveridden = &f

	res, err := Export("env_config", &s, true)

	if err != nil {
		t.Error(err.Error())
	}

	// test default filling
	if s.SomePointerWithDefault == nil {
		t.Errorf("expected SomePointerWithDefault to not be nil")
	}
	if *s.SomePointerWithDefault != "foo2baz" {
		t.Errorf("expected %s, got %s", "foo2baz", *s.SomePointerWithDefault)
	}
	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}
	if res[0] != "ENV_CONFIG_ENABLED=true" {
		t.Errorf("expected %v, got %s", "ENV_CONFIG_ENABLED=true", res[0])
	}
	if res[1] != "ENV_CONFIG_EMBEDDEDPORT=5000" {
		t.Errorf("expected %v, got %s", "ENV_CONFIG_EMBEDDEDPORT=5000", res[1])
	}
	if res[2] != "ENV_CONFIG_MULTIWORDVAR=fooembedded" {
		t.Errorf("expected %v, got %s", "ENV_CONFIG_MULTIWORDVAR=fooembedded", res[2])
	}
	if res[3] != "MULTI_WITH_DIFFERENT_ALT=bazembedded" {
		t.Errorf("expected %v, got %s", "MULTI_WITH_DIFFERENT_ALT=bazembedded", res[3])
	}
	if res[4] != "EMBEDDED_WITH_ALT=embeddedalt" {
		t.Errorf("expected %v, got %s", "EMBEDDED_WITH_ALT=embeddedalt", res[4])
	}
	if res[5] != "ENV_CONFIG_DEBUG=true" {
		t.Errorf("expected %v, got %s", "ENV_CONFIG_DEBUG=true", res[5])
	}
	if res[6] != "ENV_CONFIG_PORT=8080" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_PORT=8080", res[6])
	}
	if res[7] != "ENV_CONFIG_RATE=0.5" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_RATE=0.5", res[7])
	}
	if res[9] != "ENV_CONFIG_TTL=30" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_TTL=30", res[9])
	}
	if res[8] != "ENV_CONFIG_USER=Kelsey" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_USER=Kelsey", res[8])
	}
	if res[10] != "ENV_CONFIG_TIMEOUT=2m0s" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_TIMEOUT=2m0s", res[10])
	}
	{
		admins := strings.Split(res[11], "=")
		if len(admins) != 2 || admins[0] != "ENV_CONFIG_ADMINUSERS" {
			t.Errorf("expected %s, got %#v", "ENV_CONFIG_ADMINUSERS", admins[0])
		} else {
			admins = strings.Split(admins[1], ",")
			if len(admins) != 3 ||
				admins[0] != "John" ||
				admins[1] != "Adam" ||
				admins[2] != "Will" {
				t.Errorf("expected %#v, got %#v", []string{"John", "Adam", "Will"}, admins)
			}
		}
	}

	{
		magic := strings.Split(res[12], "=")

		if len(magic) != 2 || magic[0] != "ENV_CONFIG_MAGICNUMBERS" {
			t.Errorf("expected %s, got %#v", "ENV_CONFIG_MAGICNUMBERS", magic[0])
		} else {
			magic = strings.Split(magic[1], ",")
			if len(magic) != 3 ||
				magic[0] != "5" ||
				magic[1] != "10" ||
				magic[2] != "20" {
				t.Errorf("expected %#v, got %#v", []string{"5", "10", "20"}, magic)
			}
		}
	}

	if res[13] != "ENV_CONFIG_MULTIWORDVAR=foo bar" {
		t.Errorf("expected %v, got %v", "ENV_CONFIG_MULTIWORDVAR=foo bar", res[13])
	}
	if res[14] != "ENV_CONFIG_SOMEPOINTER=foo bar" {
		t.Errorf("expected %v, got %v", "ENV_CONFIG_SOMEPOINTER=foo bar", res[14])
	}
	if res[15] != "ENV_CONFIG_SOMEPOINTERWITHDEFAULT=foo2baz" {
		t.Errorf("expected %v, got %v", "ENV_CONFIG_SOMEPOINTERWITHDEFAULT=foo2baz", res[15])
	}
	if res[16] != "MULTI_WORD_VAR_WITH_ALT=ALT" {
		t.Errorf("expected %v, got %v", "MULTI_WORD_VAR_WITH_ALT=ALT", res[16])
	}
	// wrong?
	if res[17] != "MULTI_WORD_VAR_WITH_LOWER_CASE_ALT=lower_casE" {
		t.Errorf("expected %v, got %v", "MULTI_WORD_VAR_WITH_LOWER_CASE_ALT=lower_casE", res[17])
	}
	if res[18] != "SERVICE_HOST=127.0.0.1" {
		t.Errorf("expected %v, got %v", "SERVICE_HOST=127.0.0.1", res[20])
	}
	if res[20] != "ENV_CONFIG_REQUIREDVAR=foo" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_REQUIREDVAR=foo", res[20])
	}
	if res[21] != "BROKER=127.0.0.1" {
		t.Errorf("expected %s, got %s", "BROKER=127.0.0.1", res[21])
	}
	if res[22] != "ENV_CONFIG_REQUIREDDEFAULT=foo2bar" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_REQUIREDDEFAULT=foo2bar", res[22])
	}
	if res[23] != "ENV_CONFIG_TIMEOUTWITHDEFAULT=32h" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_TIMEOUTWITHDEFAULT=32h", res[23])
	}
	if res[24] != "ENV_CONFIG_DEFAULTINT=7" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_DEFAULTINT=7", res[24])
	}
	if res[25] != "ENV_CONFIG_DEFAULTBOOLPOINTER=true" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_DEFAULTBOOLPOINTER=true", res[25])
	}
	if res[26] != "ENV_CONFIG_DEFAULTBOOLPOINTEROVERIDDEN=false" {
		t.Errorf("expected %s, got %s", "ENV_CONFIG_DEFAULTBOOLPOINTEROVERIDDEN=false", res[26])
	}
	// expect ignored & NotRequiredVar to not be there
	if len(res) != 27 {
		t.Errorf("expected length to be %d, got %d", 27, len(res))
	}
}

func TestParseErrorBool(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_DEBUG")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Debug" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}

}

func TestParseErrorFloat32(t *testing.T) {
	var s Specification

	os.Setenv("ENV_CONFIG_RATE", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_RATE")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Rate" {
		t.Errorf("expected %s, got %v", "Rate", v.FieldName)
	}
	if s.Rate != 0 {
		t.Errorf("expected %v, got %v", 0, s.Rate)
	}

}

func TestParseErrorInt(t *testing.T) {
	var s Specification
	//os.Clearenv()

	os.Setenv("ENV_CONFIG_PORT", "string")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")

	defer os.Unsetenv("ENV_CONFIG_PORT")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Port" {
		t.Errorf("expected %s, got %v", "Port", v.FieldName)
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}

}

func TestParseErrorUint(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_TTL", "-30")
	defer os.Unsetenv("ENV_CONFIG_TTL")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "TTL" {
		t.Errorf("expected %s, got %v", "TTL", v.FieldName)
	}
	if s.TTL != 0 {
		t.Errorf("expected %v, got %v", 0, s.TTL)
	}

}

func TestErrInvalidSpecification(t *testing.T) {
	//os.Clearenv()
	m := make(map[string]string)
	err := Process("env_config", &m)
	if err != ErrInvalidSpecification {
		t.Errorf("expected %v, got %v", ErrInvalidSpecification, err)
	}
}

func TestUnsetVars(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("USER", "foo")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("USER")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// If the var is not defined the non-prefixed version should not be used
	// unless the struct tag says so
	if s.User != "" {
		t.Errorf("expected %q, got %q", "", s.User)
	}

}

func TestAlternateVarNames(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_LOWER_CASE_ALT", "baz")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_MULTI_WORD_VAR")
	defer os.Unsetenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT")
	defer os.Unsetenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_LOWER_CASE_ALT")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// Setting the alt version of the var in the environment has no effect if
	// the struct tag is not supplied
	if s.MultiWordVar != "" {
		t.Errorf("expected %q, got %q", "", s.MultiWordVar)
	}

	// Setting the alt version of the var in the environment correctly sets
	// the value if the struct tag IS supplied
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %q, got %q", "bar", s.MultiWordVarWithAlt)
	}

	// Alt value is not case sensitive and is treated as all uppercase
	if s.MultiWordVarWithLowerCaseAlt != "baz" {
		t.Errorf("expected %q, got %q", "baz", s.MultiWordVarWithLowerCaseAlt)
	}
}

func TestRequiredVar(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foobar")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.RequiredVar)
	}

}

func TestBlankDefaultVar(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "requiredvalue")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}

	if *s.SomePointerWithDefault != "foo2baz" {
		t.Errorf("expected %s, got %s", "foo2baz", *s.SomePointerWithDefault)
	}

}

func TestNonBlankDefaultVar(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULTVAR", "nondefaultval")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "requiredvalue")

	defer os.Unsetenv("ENV_CONFIG_DEFAULTVAR")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "nondefaultval" {
		t.Errorf("expected %s, got %s", "nondefaultval", s.DefaultVar)
	}

}

func TestExplicitBlankDefaultVar(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULTVAR", "")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "")

	defer os.Unsetenv("ENV_CONFIG_DEFAULTVAR")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.DefaultVar != "" {
		t.Errorf("expected %s, got %s", "\"\"", s.DefaultVar)
	}

}

func TestAlternateNameDefaultVar(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("BROKER", "betterbroker")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.NoPrefixDefault != "betterbroker" {
		t.Errorf("expected %q, got %q", "betterbroker", s.NoPrefixDefault)
	}

	os.Unsetenv("BROKER")
	os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.NoPrefixDefault != "127.0.0.1" {
		t.Errorf("expected %q, got %q", "127.0.0.1", s.NoPrefixDefault)
	}

}

func TestRequiredDefault(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.RequiredDefault != "foo2bar" {
		t.Errorf("expected %q, got %q", "foo2bar", s.RequiredDefault)
	}
}

func TestPointerFieldBlank(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.SomePointer != nil {
		t.Errorf("expected <nil>, got %s", *s.SomePointer)
	}
}

func TestMustProcess(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "foo")

	defer os.Unsetenv("ENV_CONFIG_DEBUG")
	defer os.Unsetenv("ENV_CONFIG_PORT")
	defer os.Unsetenv("ENV_CONFIG_RATE")
	defer os.Unsetenv("ENV_CONFIG_USER")
	defer os.Unsetenv("SERVICE_HOST")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	MustProcess("env_config", &s)

	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Error("expected panic")
	}()
	m := make(map[string]string)
	MustProcess("env_config", &m)
}

func TestEmbeddedStruct(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "required")
	os.Setenv("ENV_CONFIG_ENABLED", "true")
	os.Setenv("ENV_CONFIG_EMBEDDEDPORT", "1234")
	os.Setenv("ENV_CONFIG_MULTIWORDVAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WITH_DIFFERENT_ALT", "baz")
	os.Setenv("ENV_CONFIG_EMBEDDED_WITH_ALT", "foobar")
	os.Setenv("ENV_CONFIG_SOMEPOINTER", "foobaz")
	os.Setenv("ENV_CONFIG_EMBEDDED_IGNORED", "was-not-ignored")

	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	defer os.Unsetenv("ENV_CONFIG_ENABLED")
	defer os.Unsetenv("ENV_CONFIG_EMBEDDEDPORT")
	defer os.Unsetenv("ENV_CONFIG_MULTIWORDVAR")
	defer os.Unsetenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT")
	defer os.Unsetenv("ENV_CONFIG_MULTI_WITH_DIFFERENT_ALT")
	defer os.Unsetenv("ENV_CONFIG_EMBEDDED_WITH_ALT")
	defer os.Unsetenv("ENV_CONFIG_SOMEPOINTER")
	defer os.Unsetenv("ENV_CONFIG_EMBEDDED_IGNORED")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if !s.Enabled {
		t.Errorf("expected %v, got %v", true, s.Enabled)
	}
	if s.EmbeddedPort != 1234 {
		t.Errorf("expected %d, got %v", 1234, s.EmbeddedPort)
	}
	if s.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.MultiWordVar)
	}
	if s.Embedded.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.Embedded.MultiWordVar)
	}
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %s, got %s", "bar", s.MultiWordVarWithAlt)
	}
	if s.Embedded.MultiWordVarWithAlt != "baz" {
		t.Errorf("expected %s, got %s", "baz", s.Embedded.MultiWordVarWithAlt)
	}
	if s.EmbeddedAlt != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.EmbeddedAlt)
	}
	if *s.SomePointer != "foobaz" {
		t.Errorf("expected %s, got %s", "foobaz", *s.SomePointer)
	}
	if s.EmbeddedIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestEmbeddedButIgnoredStruct(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "required")
	os.Setenv("ENV_CONFIG_FIRSTEMBEDDEDBUTIGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_SECONDEMBEDDEDBUTIGNORED", "was-not-ignored")

	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")
	defer os.Unsetenv("ENV_CONFIG_FIRSTEMBEDDEDBUTIGNORED")
	defer os.Unsetenv("ENV_CONFIG_SECONDEMBEDDEDBUTIGNORED")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if s.FirstEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
	if s.SecondEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestNonPointerFailsProperly(t *testing.T) {
	var s Specification
	//os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIREDVAR", "snap")
	defer os.Unsetenv("ENV_CONFIG_REQUIREDVAR")

	err := Process("env_config", s)
	if err != ErrInvalidSpecification {
		t.Errorf("non-pointer should fail with ErrInvalidSpecification, was instead %s", err)
	}
}

func TestCustomDecoder(t *testing.T) {
	s := struct {
		Foo string
		Bar bracketed
	}{}

	//os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")

	defer os.Unsetenv("ENV_CONFIG_FOO")
	defer os.Unsetenv("ENV_CONFIG_BAR")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.Foo != "foo" {
		t.Errorf("foo: expected 'foo', got %q", s.Foo)
	}

	if string(s.Bar) != "[bar]" {
		t.Errorf("bar: expected '[bar]', got %q", string(s.Bar))
	}
}

func TestCustomDecoderWithPointer(t *testing.T) {
	s := struct {
		Foo string
		Bar *bracketed
	}{}

	// Decode would panic when b is nil, so make sure it
	// has an initial value to replace.
	var b bracketed = "initial_value"
	s.Bar = &b

	//os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")

	defer os.Unsetenv("ENV_CONFIG_FOO")
	defer os.Unsetenv("ENV_CONFIG_BAR")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.Foo != "foo" {
		t.Errorf("foo: expected 'foo', got %q", s.Foo)
	}

	if string(*s.Bar) != "[bar]" {
		t.Errorf("bar: expected '[bar]', got %q", string(*s.Bar))
	}
}

type bracketed string

func (b *bracketed) Decode(value string) error {
	*b = bracketed("[" + value + "]")
	return nil
}
