package envconfig

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type HonorDecodeInStruct struct {
	Value string
}

func (h *HonorDecodeInStruct) Decode(env string) error {
	h.Value = "decoded"
	return nil
}

type CustomURL struct {
	Value *url.URL
}

func (cu *CustomURL) UnmarshalBinary(data []byte) error {
	u, err := url.Parse(string(data))
	cu.Value = u
	return err
}

type Specification struct {
	Embedded
	Debug                        bool           `env:"ENV_CONFIG_DEBUG"`
	Port                         int            `env:"ENV_CONFIG_PORT"`
	Rate                         float32        `env:"ENV_CONFIG_RATE"`
	User                         string         `env:"ENV_CONFIG_USER"`
	TTL                          uint32         `env:"ENV_CONFIG_TTL"`
	Timeout                      time.Duration  `env:"ENV_CONFIG_TIMEOUT"`
	AdminUsers                   []string       `env:"ENV_CONFIG_ADMIN_USERS"`
	MagicNumbers                 []int          `env:"ENV_CONFIG_MAGIC_NUMBERS"`
	EmptyNumbers                 []int          `env:"ENV_CONFIG_EMPTY_NUMBERS"`
	ByteSlice                    []byte         `env:"ENV_CONFIG_BYTE_SLICE"`
	ColorCodes                   map[string]int `env:"ENV_CONFIG_COLOR_CODES"`
	MultiWordVar                 string         `env:"ENV_CONFIG_MULTI_WORD_VAR"`
	MultiWordVarWithAutoSplit    uint32         `env:"ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT"`
	MultiWordACRWithAutoSplit    uint32         `env:"ENV_CONFIG_MULTI_WORD_ACR_WITH_AUTO_SPLIT"`
	SomePointer                  *string        `env:"ENV_CONFIG_SOME_POINTER"`
	SomePointerWithDefault       *string        `env:"ENV_CONFIG_SOME_POINTER_WITH_DEFAULT" default:"foo2baz" comment:"foorbar is the word"`
	MultiWordVarWithAlt          string         `env:"ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT" comment:"what alt"`
	MultiWordVarWithLowerCaseAlt string         `env:"ENV_CONFIG_multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string         `env:"ENV_CONFIG_SERVICE_HOST"`
	DefaultVar                   string         `env:"ENV_CONFIG_DEFAULT_VAR" default:"foobar"`
	RequiredVar                  string         `env:"ENV_CONFIG_REQUIRED_VAR"`
	NoPrefixDefault              string         `env:"ENV_CONFIG_BROKER" default:"127.0.0.1"`
	RequiredDefault              string         `env:"ENV_CONFIG_REQUIRED_DEFAULT" default:"foo2bar"`
	NestedSpecification          struct {
		Property            string `env:"inner"`
		PropertyWithDefault string `default:"fuzzybydefault"`
	}
	AfterNested  string              `env:"ENV_CONFIG_AFTER_NESTED"`
	DecodeStruct HonorDecodeInStruct `env:"honor"`
	Datetime     time.Time           `env:"ENV_CONFIG_DATETIME"`
	MapField     map[string]string   `env:"ENV_CONFIG_MAPFIELD" default:"one:two,three:four"`
	UrlValue     CustomURL           `env:"ENV_CONFIG_URL_VALUE"`
	UrlPointer   *CustomURL          `env:"ENV_CONFIG_URL_POINTER"`
}

type Embedded struct {
	Enabled             bool   `env:"ENV_CONFIG_ENABLED" comment:"some embedded value"`
	EmbeddedPort        int    `env:"ENV_CONFIG_EMBEDDED_PORT"`
	MultiWordVar        string `env:"ENV_CONFIG_MULTI_WORD_VAR"`
	MultiWordVarWithAlt string `env:"ENV_CONFIG_MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `env:"ENV_CONFIG_EMBEDDED_WITH_ALT"`
}

func TestProcess(t *testing.T) {
	var s struct {
		Debug                  bool           `env:"ENV_CONFIG_DEBUG"`
		Port                   int            `env:"ENV_CONFIG_PORT"`
		Rate                   float32        `env:"ENV_CONFIG_RATE"`
		User                   string         `env:"ENV_CONFIG_USER"`
		TTL                    uint32         `env:"ENV_CONFIG_TTL"`
		Timeout                time.Duration  `env:"ENV_CONFIG_TIMEOUT"`
		AdminUsers             []string       `env:"ENV_CONFIG_ADMIN_USERS"`
		MagicNumbers           []int          `env:"ENV_CONFIG_MAGIC_NUMBERS"`
		EmptyNumbers           []int          `env:"ENV_CONFIG_EMPTY_NUMBERS"`
		ByteSlice              []byte         `env:"ENV_CONFIG_BYTE_SLICE"`
		ColorCodes             map[string]int `env:"ENV_CONFIG_COLOR_CODES"`
		SomePointer            *string        `env:"ENV_CONFIG_SOME_POINTER"`
		SomePointerWithDefault *string        `env:"ENV_CONFIG_SOME_POINTER_WITH_DEFAULT" default:"foo2baz" comment:"foorbar is the word"`
		DefaultVar             string         `env:"ENV_CONFIG_DEFAULT_VAR" default:"foobar"`
		NestedSpecification    struct {
			Property string `env:"inner"`
		}
		AfterNested  string              `env:"ENV_CONFIG_AFTER_NESTED"`
		DecodeStruct HonorDecodeInStruct `env:"ENV_CONFIG_HONOR"`
		Datetime     time.Time           `env:"ENV_CONFIG_DATETIME"`
		MapField     map[string]string   `env:"ENV_CONFIG_MAPFIELD" default:"one:two,three:four"`
		UrlValue     CustomURL           `env:"ENV_CONFIG_URL_VALUE"`
		UrlPointer   *CustomURL          `env:"ENV_CONFIG_URL_POINTER"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMIN_USERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGIC_NUMBERS", "5,10,20")
	os.Setenv("ENV_CONFIG_EMPTY_NUMBERS", "")
	os.Setenv("ENV_CONFIG_BYTE_SLICE", "this is a test value")
	os.Setenv("ENV_CONFIG_COLOR_CODES", "red:1,green:2,blue:3")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_REQUIRED_VAR", "foo")
	os.Setenv("ENV_CONFIG_SOME_POINTER", "some-pointer")
	os.Setenv("ENV_CONFIG_OUTER_INNER", "iamnested")
	os.Setenv("ENV_CONFIG_AFTER_NESTED", "after")
	os.Setenv("ENV_CONFIG_HONOR", "honor")
	os.Setenv("ENV_CONFIG_DATETIME", "2016-08-16T18:57:05Z")
	os.Setenv("ENV_CONFIG_URL_VALUE", "https://example.com/foo")
	os.Setenv("ENV_CONFIG_URL_POINTER", "https://example.com/foo")
	os.Setenv("inner", "inner")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

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

	assert.Equal(t, []string{
		"John",
		"Adam",
		"Will",
	}, s.AdminUsers)

	assert.Equal(t, []int{5, 10, 20}, s.MagicNumbers)
	assert.Empty(t, s.EmptyNumbers)

	assert.Equal(t, "this is a test value", string(s.ByteSlice))

	assert.Equal(t, map[string]int{
		"red":   1,
		"green": 2,
		"blue":  3,
	}, s.ColorCodes)

	assert.Equal(t, "inner", s.NestedSpecification.Property)
	assert.Equal(t, "after", s.AfterNested)
	assert.Equal(t, "decoded", s.DecodeStruct.Value)

	if expected := time.Date(2016, 8, 16, 18, 57, 05, 0, time.UTC); !s.Datetime.Equal(expected) {
		t.Errorf("expected %s, got %s", expected.Format(time.RFC3339), s.Datetime.Format(time.RFC3339))
	}

	u, err := url.Parse("https://example.com/foo")
	require.NoError(t, err)

	assert.Equal(t, *u, *s.UrlValue.Value)
	assert.Equal(t, *u, *s.UrlPointer.Value)
}

func TestParseErrorBool(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "string")
	_, err := Unmarshal(&s)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

func TestParseErrorFloat32(t *testing.T) {
	var s struct {
		Rate float32 `env:"ENV_CONFIG_RATE"`
		Port int     `env:"ENV_CONFIG_PORT"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_PORT", "string")
	_, err := Unmarshal(&s)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}
}

func TestParseErrorUint(t *testing.T) {
	var s struct {
		TTL uint32 `env:"ENV_CONFIG_TTL"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_TTL", "-30")
	_, err := Unmarshal(&s)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if s.TTL != 0 {
		t.Errorf("expected %v, got %v", 0, s.TTL)
	}
}

func TestParseErrorSplitWords(t *testing.T) {
	var s struct {
		MultiWordVarWithAutoSplit uint32 `env:"ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT", "shakespeare")
	_, err := Unmarshal(&s)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if s.MultiWordVarWithAutoSplit != 0 {
		t.Errorf("expected %v, got %v", 0, s.MultiWordVarWithAutoSplit)
	}
}

func TestErrInvalidSpecification(t *testing.T) {
	m := make(map[string]string)
	_, err := Unmarshal(&m)

	if !errors.Is(err, ErrInvalidSpecification) {
		t.Errorf("expected %v, got %v", ErrInvalidSpecification, err)
	}
}

func TestUnsetVars(t *testing.T) {
	var s struct {
		User string `env:"ENV_CONFIG_USER"`
	}
	os.Clearenv()

	_, err := Unmarshal(&s)
	assert.Error(t, err)
	assert.Equal(t, "", s.User)
}

func TestRequiredVar(t *testing.T) {
	var s struct {
		RequiredVar string `env:"ENV_CONFIG_REQUIRED_VAR"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIRED_VAR", "foobar")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.RequiredVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.RequiredVar)
	}
}

func TestRequiredMissing(t *testing.T) {
	var s struct {
		RequiredVar string `env:"ENV_CONFIG_REQUIRED_VAR" required:"True"`
	}
	os.Clearenv()

	_, err := Unmarshal(&s)
	require.Error(t, err)
}

func TestBlankDefaultVar(t *testing.T) {
	var s struct {
		SomePointerWithDefault *string `env:"ENV_CONFIG_SOME_POINTER_WITH_DEFAULT" default:"foo2baz" comment:"foorbar is the word"`
		DefaultVar             string  `env:"ENV_CONFIG_DEFAULT_VAR" default:"foobar"`
	}
	os.Clearenv()

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.DefaultVar != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.DefaultVar)
	}

	if *s.SomePointerWithDefault != "foo2baz" {
		t.Errorf("expected %s, got %s", "foo2baz", *s.SomePointerWithDefault)
	}
}

func TestNonBlankDefaultVar(t *testing.T) {
	var s struct {
		DefaultVar string `env:"ENV_CONFIG_DEFAULT_VAR" default:"foobar"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULT_VAR", "nondefaultval")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.DefaultVar != "nondefaultval" {
		t.Errorf("expected %s, got %s", "nondefaultval", s.DefaultVar)
	}
}

func TestExplicitBlankDefaultVar(t *testing.T) {
	var s struct {
		DefaultVar string `env:"ENV_CONFIG_DEFAULT_VAR" default:"foobar"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEFAULT_VAR", "")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.DefaultVar != "" {
		t.Errorf("expected %s, got %s", "\"\"", s.DefaultVar)
	}
}

func TestRequiredDefault(t *testing.T) {
	var s struct {
		RequiredDefault string `env:"ENV_CONFIG_REQUIRED_DEFAULT"  default:"foo2bar"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIRED_VAR", "foo")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.RequiredDefault != "foo2bar" {
		t.Errorf("expected %q, got %q", "foo2bar", s.RequiredDefault)
	}
}

func TestPointerFieldBlank(t *testing.T) {
	var s struct {
		SomePointer *string `env:"ENV_CONFIG_SOME_POINTER"`
	}
	os.Clearenv()

	_, err := Unmarshal(&s)
	assert.Error(t, err)
	assert.Nil(t, s.SomePointer)
}

func TestEnvVarCollision(t *testing.T) {
	var s struct {
		A string `env:"A"`
		B string `env:"A"`
	}
	os.Clearenv()
	os.Setenv("A", "true")

	_, err := Unmarshal(&s)
	assert.Error(t, err)
	assert.Equal(t, "", s.A)
	assert.Equal(t, "", s.B)

}

func TestEmbeddedStruct(t *testing.T) {
	type Embedded struct {
		Bool   bool   `env:"ENV_CONFIG_ENABLED" comment:"some embedded value"`
		Int    int    `env:"ENV_CONFIG_EMBEDDED_PORT"`
		String string `env:"ENV_CONFIG_MULTI_WORD_VAR"`
	}

	var s struct {
		Embedded `env:"ENV_CONFIG_EMBEDDED" comment:"can we document a struct"`
		String   string `env:"ENV_CONFIG_STRING"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_ENABLED", "true")
	os.Setenv("ENV_CONFIG_EMBEDDED", "!foo!")
	os.Setenv("ENV_CONFIG_EMBEDDED_PORT", "1234")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR", "foo")
	os.Setenv("ENV_CONFIG_STRING", "foobaz")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	assert.True(t, s.Embedded.Bool)
	assert.Equal(t, 1234, s.Embedded.Int)
	assert.Equal(t, "foo", s.Embedded.String)
	assert.Equal(t, "foobaz", s.String)
}

func TestNonPointerFailsProperly(t *testing.T) {
	var s struct{}
	os.Clearenv()

	_, err := Unmarshal(s)
	if !errors.Is(err, ErrInvalidSpecification) {
		t.Errorf("non-pointer should fail with ErrInvalidSpecification, was instead %s", err)
	}
}

func TestCustomValueFields(t *testing.T) {
	var s struct {
		Foo    string       `env:"ENV_CONFIG_FOO"`
		Bar    bracketed    `env:"ENV_CONFIG_BAR"`
		Baz    quoted       `env:"ENV_CONFIG_BAZ"`
		Struct setterStruct `env:"ENV_CONFIG_STRUCT"`
	}

	// Set would panic when the receiver is nil,
	// so make sure it has an initial value to replace.
	s.Baz = quoted{new(bracketed)}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")
	os.Setenv("ENV_CONFIG_BAZ", "baz")
	os.Setenv("ENV_CONFIG_STRUCT", "inner")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `["baz"]`; s.Baz.String() != want {
		t.Errorf(`baz: got %#q, want %#q`, s.Baz, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestCustomPointerFields(t *testing.T) {
	var s struct {
		Foo    string        `env:"ENV_CONFIG_FOO"`
		Bar    *bracketed    `env:"ENV_CONFIG_BAR"`
		Baz    *quoted       `env:"ENV_CONFIG_BAZ"`
		Struct *setterStruct `env:"ENV_CONFIG_STRUCT"`
	}

	// Set would panic when the receiver is nil,
	// so make sure they have initial values to replace.
	s.Bar = new(bracketed)
	s.Baz = &quoted{new(bracketed)}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")
	os.Setenv("ENV_CONFIG_BAZ", "baz")
	os.Setenv("ENV_CONFIG_STRUCT", "inner")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if want := "foo"; s.Foo != want {
		t.Errorf("foo: got %#q, want %#q", s.Foo, want)
	}

	if want := "[bar]"; s.Bar.String() != want {
		t.Errorf("bar: got %#q, want %#q", s.Bar, want)
	}

	if want := `["baz"]`; s.Baz.String() != want {
		t.Errorf(`baz: got %#q, want %#q`, s.Baz, want)
	}

	if want := `setterstruct{"inner"}`; s.Struct.Inner != want {
		t.Errorf(`Struct.Inner: got %#q, want %#q`, s.Struct.Inner, want)
	}
}

func TestEmptyPrefixUsesFieldNames(t *testing.T) {
	var s struct {
		RequiredVar string `env:"ENV_CONFIG_REQUIRED_VAR" required:"True"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_REQUIRED_VAR", "foo")

	_, err := Unmarshal(&s)
	require.NoError(t, err)

	if s.RequiredVar != "foo" {
		t.Errorf(
			`RequiredVar not populated correctly: expected "foo", got %q`,
			s.RequiredVar,
		)
	}
}

func TestTextUnmarshalerError(t *testing.T) {
	var s struct {
		Datetime time.Time `env:"ENV_CONFIG_DATETIME"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DATETIME", "I'M NOT A DATE")

	_, err := Unmarshal(&s)
	var actualLowLevelError *time.ParseError
	if !errors.As(err, &actualLowLevelError) {
		t.Error("error must be as *time.ParseError")
	}
}

func TestBinaryUnmarshalerError(t *testing.T) {
	var s struct {
		UrlPointer *CustomURL `env:"ENV_CONFIG_URL_POINTER"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_URL_POINTER", "http://%41:8080/")

	_, err := Unmarshal(&s)

	// To be compatible with go 1.5 and lower we should do a very basic check,
	// because underlying error message varies in go 1.5 and go 1.6+.

	var ue *url.Error
	if !errors.As(err, &ue) {
		t.Fatal("error must be as *url.Error")
	}

	if ue.Op != "parse" {
		t.Errorf("expected error op to be \"parse\", got %q", ue.Op)
	}
}

func TestCheckUnknownEmptyPrefix(t *testing.T) {
	var s struct{}
	os.Clearenv()
	_, err := FindUnknownEnvVariablesByPrefix("", &s)
	if !errors.Is(err, ErrEmptyPrefix) {
		t.Errorf("expected ErrEmptyPrefix error, got %s", err)
	}
}

func TestCheckDisallowedOnlyAllowed(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_UNRELATED_ENV_VAR", "true")
	os.Setenv("ENV_CONFIG_UNRELATED_ENV_VAR2", "true")
	unknownVars, err := FindUnknownEnvVariablesByPrefix("ENV_CONFIG_", &s)
	require.NoError(t, err)
	require.Equal(t, []string{"ENV_CONFIG_UNRELATED_ENV_VAR", "ENV_CONFIG_UNRELATED_ENV_VAR2"}, unknownVars)
}

func TestCheckDisallowedMispelled(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_ZEBUG", "false")

	unknownVars, err := FindUnknownEnvVariablesByPrefix("ENV_CONFIG_", &s)
	require.NoError(t, err)
	require.Equal(t, []string{"ENV_CONFIG_ZEBUG"}, unknownVars)
}

func TestCheckDisallowedIgnored(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG"`
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_IGNORED", "false")

	unknownVars, err := FindUnknownEnvVariablesByPrefix("ENV_CONFIG_", &s)
	require.NoError(t, err)
	require.Equal(t, []string{"ENV_CONFIG_IGNORED"}, unknownVars)
}

func TestErrorMessageForRequired(t *testing.T) {
	var s struct {
		Foo string `env:"BAR" `
	}

	os.Clearenv()
	_, err := Unmarshal(&s)
	require.Error(t, err)
	require.Contains(t, err.Error(), "BAR")
}

type bracketed string

func (b *bracketed) Set(value string) error {
	*b = bracketed("[" + value + "]")
	return nil
}

func (b bracketed) String() string {
	return string(b)
}

// quoted is used to test the precedence of Decode over Set.
// The sole field is a flag.Value rather than a setter to validate that
// all flag.Value implementations are also Setter implementations.
type quoted struct{ flag.Value }

func (d quoted) Decode(value string) error {
	return d.Set(`"` + value + `"`)
}

type setterStruct struct {
	Inner string
}

func (ss *setterStruct) Set(value string) error {
	ss.Inner = fmt.Sprintf("setterstruct{%q}", value)
	return nil
}

func BenchmarkGatherInfo(b *testing.B) {
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMIN_USERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGIC_NUMBERS", "5,10,20")
	os.Setenv("ENV_CONFIG_COLOR_CODES", "red:1,green:2,blue:3")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_REQUIRED_VAR", "foo")
	os.Setenv("ENV_CONFIG_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_OUTER_INNER", "iamnested")
	os.Setenv("ENV_CONFIG_AFTER_NESTED", "after")
	os.Setenv("ENV_CONFIG_HONOR", "honor")
	os.Setenv("ENV_CONFIG_DATETIME", "2016-08-16T18:57:05Z")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_AUTO_SPLIT", "24")
	for i := 0; i < b.N; i++ {
		var s Specification
		gatherInfo(&s)
	}
}
