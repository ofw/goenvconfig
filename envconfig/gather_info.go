package envconfig

import (
	"fmt"
	"os"
	"reflect"
	"sort"
)

// envVarInfo maintains information about the configuration variable
type envVarInfo struct {
	Name         string
	Key          string
	Field        reflect.Value
	Comment      string
	Default      string
	IsDefaultSet bool
}

func (info *envVarInfo) GetValueFromEnv() (string, error) {
	if info.Key == "" {
		return "", fmt.Errorf(`"env" tag is empty on struct field: %s`, info.Name)
	}

	// `os.Getenv` cannot differentiate between an explicitly set empty value
	// and an unset value. `os.LookupEnv` is preferred to `syscall.Getenv`.
	value, ok := os.LookupEnv(info.Key)
	if !ok {
		if info.IsDefaultSet {
			value = info.Default
		} else {
			return "", fmt.Errorf("env variable is not set: %q", info.Key)
		}
	}
	return value, nil
}

// gatherInfo gathers information about the specified struct
func gatherInfo(spec interface{}) ([]envVarInfo, error) {
	s := reflect.ValueOf(spec)

	if s.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("spec must be a pointer: %w", ErrInvalidSpecification)
	}

	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return nil, fmt.Errorf("spec must be a pointer to a struct: %w", ErrInvalidSpecification)
	}

	typeOfSpec := s.Type()

	// over allocate an info array, we will extend if needed later
	infos := make([]envVarInfo, 0, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		ftype := typeOfSpec.Field(i)

		if !f.CanSet() {
			return nil, fmt.Errorf("spec must not contain unexported fields: %w", ErrInvalidSpecification)
		}

		f = followPointerChain(f)

		// handle embedded and referenced structs
		if f.Kind() == reflect.Struct && !implementsInterface(ftype.Type) {
			embeddedPtr := f.Addr().Interface()
			embeddedInfos, err := gatherInfo(embeddedPtr)
			if err != nil {
				return nil, err
			}
			infos = append(infos, embeddedInfos...)
		} else {
			// Capture information about the config variable
			infos = append(infos, createEnvVarInfo(f, ftype))
		}
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Key < infos[j].Key
	})
	return infos, nil
}

func createEnvVarInfo(f reflect.Value, ftype reflect.StructField) envVarInfo {
	defaultValue, isDefaultSet := ftype.Tag.Lookup("default")
	return envVarInfo{
		Name:         ftype.Name,
		Field:        f,
		Comment:      ftype.Tag.Get("comment"),
		Key:          ftype.Tag.Get("env"),
		Default:      defaultValue,
		IsDefaultSet: isDefaultSet,
	}
}

func followPointerChain(f reflect.Value) reflect.Value {
	for f.Kind() == reflect.Ptr {
		if f.IsNil() {
			if f.Type().Elem().Kind() != reflect.Struct {
				// nil pointer to a non-struct: leave it alone
				break
			}
			// nil pointer to struct: create a zero instance
			f.Set(reflect.New(f.Type().Elem()))
		}
		f = f.Elem()
	}
	return f
}
