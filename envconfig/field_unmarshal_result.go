package envconfig

import (
	"io"
	"os"
)

var output io.Writer = os.Stdout

// A FieldUnmarshalResult. Err is filled when an environment variable is required but missing or
// cannot be converted to the type required by a struct field during assignment.
type FieldUnmarshalResult struct {
	KeyName   string
	FieldName string
	TypeName  string
	Value     string
	Err       error
}

type FieldUnmarshalResults []FieldUnmarshalResult

func (f FieldUnmarshalResults) firstError() error {
	for _, result := range f {
		if result.Err != nil {
			return result.Err
		}
	}
	return nil
}
