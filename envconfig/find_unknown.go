package envconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrEmptyPrefix indicates that prefix is empty.
var ErrEmptyPrefix = errors.New("prefix must be non-empty")

// FindUnknownEnvVariablesByPrefix checks that no environment variables with the prefix are set
// that we don't know how or want to parse. This is likely only meaningful with
// a non-empty prefix.
// Returns list of unknown env variables.
func FindUnknownEnvVariablesByPrefix(prefix string, spec interface{}) ([]string, error) {
	if prefix == "" {
		return nil, fmt.Errorf("%w", ErrEmptyPrefix)
	}
	infos, err := gatherInfo(spec)
	if err != nil {
		return nil, fmt.Errorf("gather info: %w", err)
	}

	vars := make(map[string]struct{})
	for _, info := range infos {
		vars[info.Key] = struct{}{}
	}

	var unknownVars []string
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}
		v := strings.SplitN(env, "=", 2)[0]
		if _, found := vars[v]; !found {
			unknownVars = append(unknownVars, v)
		}
	}

	return unknownVars, nil
}
