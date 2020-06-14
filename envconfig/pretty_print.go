package envconfig

import (
	"fmt"
	"text/tabwriter"
)

// PrettyPrint prints verbose information about env variables specification, current values and
// whether it was parsed correctly.
func (f FieldUnmarshalResults) PrettyPrint() {
	w := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(output, 8, 8, 4, ' ', 0)
	defer w.Flush()

	write := func(field, env, ok string) {
		fmt.Fprintf(w, " %s\t%s\t%s\t\n", field, env, ok)
	}

	write("Env Variable", "Type", "OK")
	write("----", "----", "----")

	for _, result := range f {
		var err string
		if result.Err != nil {
			err = result.Err.Error()
		} else {
			err = "v"
		}
		write(result.KeyName, result.TypeName, err)
	}
}
