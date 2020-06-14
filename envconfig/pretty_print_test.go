package envconfig

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFieldParseResults_PrettyPrint(t *testing.T) {
	type Embedded struct {
		Enabled      bool   `env:"ENV_CONFIG_ENABLED" default:"false" comment:"some embedded value"`
		EmbeddedPort int    `env:"ENV_CONFIG_EMBEDDED_PORT" default:"0"`
		MultiWordVar string `env:"ENV_CONFIG_MULTI_WORD_VAR" default:""`
	}

	var s struct {
		Embedded
		Debug        bool           `env:"ENV_CONFIG_DEBUG" default:"false"`
		Port         int            `env:"ENV_CONFIG_PORT" default:"0"`
		Rate         float32        `env:"ENV_CONFIG_RATE" default:"0"`
		User         string         `env:"ENV_CONFIG_USER" default:""`
		TTL          uint32         `env:"ENV_CONFIG_TTL"  default:"0"`
		Timeout      time.Duration  `env:"ENV_CONFIG_TIMEOUT"  default:"0"`
		AdminUsers   []string       `env:"ENV_CONFIG_ADMIN_USERS" default:""`
		MagicNumbers []int          `env:"ENV_CONFIG_MAGIC_NUMBERS" default:""`
		EmptyNumbers []int          `env:"ENV_CONFIG_EMPTY_NUMBERS" default:""`
		ByteSlice    []byte         `env:"ENV_CONFIG_BYTE_SLICE" default:""`
		ColorCodes   map[string]int `env:"ENV_CONFIG_COLOR_CODES" default:""`
		UrlValue     CustomURL      `env:"ENV_CONFIG_URL_VALUE" default:""`
		UrlPointer   *CustomURL     `env:"ENV_CONFIG_URL_POINTER" default:""`
	}

	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_EMBEDDED_PORT", "invalid")

	results, err := Unmarshal(&s)
	require.Error(t, err)

	buf := bytes.NewBuffer(nil)
	output = buf
	defer func() {
		output = os.Stdout // restore stdout
	}()
	results.PrettyPrint()

	expected := ` Env Variable                 Type                   OK                                                                                                                            
 ----                         ----                   ----                                                                                                                          
 ENV_CONFIG_ADMIN_USERS       []string               v                                                                                                                             
 ENV_CONFIG_BYTE_SLICE        []uint8                v                                                                                                                             
 ENV_CONFIG_COLOR_CODES       map[string]int         v                                                                                                                             
 ENV_CONFIG_DEBUG             bool                   v                                                                                                                             
 ENV_CONFIG_EMBEDDED_PORT     int                    assigning ENV_CONFIG_EMBEDDED_PORT="invalid" to EmbeddedPort type int: strconv.ParseInt: parsing "invalid": invalid syntax    
 ENV_CONFIG_EMPTY_NUMBERS     []int                  v                                                                                                                             
 ENV_CONFIG_ENABLED           bool                   v                                                                                                                             
 ENV_CONFIG_MAGIC_NUMBERS     []int                  v                                                                                                                             
 ENV_CONFIG_MULTI_WORD_VAR    string                 v                                                                                                                             
 ENV_CONFIG_PORT              int                    v                                                                                                                             
 ENV_CONFIG_RATE              float32                v                                                                                                                             
 ENV_CONFIG_TIMEOUT           time.Duration          v                                                                                                                             
 ENV_CONFIG_TTL               uint32                 v                                                                                                                             
 ENV_CONFIG_URL_POINTER       envconfig.CustomURL    v                                                                                                                             
 ENV_CONFIG_URL_VALUE         envconfig.CustomURL    v                                                                                                                             
 ENV_CONFIG_USER              string                 v                                                                                                                             
`

	ioutil.WriteFile("/tmp/foo", buf.Bytes(), 777)
	require.Equal(t, expected, buf.String())
}
