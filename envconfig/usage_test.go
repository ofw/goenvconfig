package envconfig

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/stretchr/testify/require"
)

var testUsageTableResult, testUsageListResult, testUsageCustomResult string

func TestMain(m *testing.M) {

	// Load the expected test results from a text file
	data, err := ioutil.ReadFile("testdata/default_table.txt")
	if err != nil {
		log.Fatal(err)
	}
	testUsageTableResult = string(data)

	data, err = ioutil.ReadFile("testdata/default_list.txt")
	if err != nil {
		log.Fatal(err)
	}
	testUsageListResult = string(data)

	data, err = ioutil.ReadFile("testdata/custom.txt")
	if err != nil {
		log.Fatal(err)
	}
	testUsageCustomResult = string(data)

	retCode := m.Run()
	os.Exit(retCode)
}

func TestUsageDefault(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG" default:"true" comment:"foo"`
		Port  int  `env:"ENV_CONFIG_PORT"`
	}

	buf := bytes.NewBuffer(nil)
	output = buf
	defer func() {
		output = os.Stdout // restore stdout
	}()

	err := PrintUsage(&s)
	out := buf.String()

	require.NoError(t, err)
	require.Equal(t, testUsageTableResult, out)
}

func TestUsageTable(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG" default:"true" comment:"foo"`
		Port  int  `env:"ENV_CONFIG_PORT"`
	}

	buf := new(bytes.Buffer)
	tabs := tabwriter.NewWriter(buf, 1, 0, 4, ' ', 0)
	err := Usagef(&s, tabs, DefaultTableFormat)
	tabs.Flush()

	require.NoError(t, err)
	require.Equal(t, testUsageTableResult, buf.String())
}

func TestUsageList(t *testing.T) {
	var s struct {
		Debug bool `env:"ENV_CONFIG_DEBUG" default:"true" comment:"foo"`
		Port  int  `env:"ENV_CONFIG_PORT"`
	}

	buf := new(bytes.Buffer)
	err := Usagef(&s, buf, DefaultListFormat)

	require.NoError(t, err)
	require.Equal(t, testUsageListResult, buf.String())
}

func TestUsageCustomFormat(t *testing.T) {

	type Embedded struct {
		Enabled bool `env:"ENV_CONFIG_ENABLED" comment:"some embedded value"`
	}

	var s struct {
		Embedded
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
		SomePointerWithDefault *string        `env:"ENV_CONFIG_SOME_POINTER_WITH_DEFAULT" default:"foo2baz" comment:"foorbar is the word"`
		NestedSpecification    struct {
			Property            string `env:"inner"`
			PropertyWithDefault string `default:"fuzzybydefault"`
		}
		DecodeStruct HonorDecodeInStruct `env:"honor"`
		Datetime     time.Time           `env:"ENV_CONFIG_DATETIME"`
		MapField     map[string]string   `env:"ENV_CONFIG_MAPFIELD" default:"one:two,three:four"`
		UrlValue     CustomURL           `env:"ENV_CONFIG_URL_VALUE"`
		UrlPointer   *CustomURL          `env:"ENV_CONFIG_URL_POINTER"`
	}

	buf := new(bytes.Buffer)
	err := Usagef(&s, buf, "{{range .}}{{usage_key .}}|{{usage_type .}}|{{usage_default .}}|{{usage_description .}}\n{{end}}")

	require.NoError(t, err)
	require.Equal(t, testUsageCustomResult, buf.String())
}

func TestUsageUnknownKeyFormat(t *testing.T) {
	var s Specification
	unknownError := "template: envconfig:1:2: executing \"envconfig\" at <.UnknownKey>"
	os.Clearenv()

	buf := new(bytes.Buffer)
	err := Usagef(&s, buf, "{{.UnknownKey}}")

	if err == nil {
		t.Errorf("expected 'unknown key' error, but got no error")
	}
	if strings.Index(err.Error(), unknownError) == -1 {
		t.Errorf("expected '%s', but got '%s'", unknownError, err.Error())
	}
}
