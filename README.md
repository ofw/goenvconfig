# envconfig

[![Go Report Card](https://goreportcard.com/badge/github.com/ofw/goenvconfig)](https://goreportcard.com/report/github.com/ofw/goenvconfig)
![Go](https://github.com/ofw/goenvconfig/workflows/Go/badge.svg)
[![GoDoc](https://godoc.org/github.com/ofw/goenvconfig?status.svg)](https://godoc.org/github.com/ofw/goenvconfig/envconfig)
[![go.dev](https://img.shields.io/badge/go.dev-pkg-007d9c.svg?style=flat)](https://pkg.go.dev/github.com/ofw/goenvconfig/envconfig)

```Go
import "github.com/ofw/goenvconfig/envconfig"
```

## Goals

This library was inspired by [envconfig](https://github.com/kelseyhightower/envconfig) with some goals in mind:
1. Simplicity: 
    - only flat structure of config. so mapping env variables to fields is obvious.
    - no guessing of env variable names. just case sensitive match.
    - all env variables are required to be set. support default values
2. Observability: 
    - supports pretty printing results of parsing env variables to struct
3. User-friendliness:
    - supports ability to pretty print results of unmarshaling env variables to config

## Usage

Set some environment variables:

```Bash
export MYAPP_DEBUG=false
export MYAPP_PORT=8080
export MYAPP_USER=Kelsey
export MYAPP_RATE="0.5"
export MYAPP_TIMEOUT="3m"
export MYAPP_USERS="rob,ken,robert"
export MYAPP_COLOR_CODES="red:1,green:2,blue:3"
```

Write some code:

```Go
package main

import (
	"fmt"
	"time"

	"github.com/ofw/goenvconfig/envconfig"
)

func main() {
	var s struct {
        Debug      bool           `env:"MYAPP_DEBUG"`
        Port       int            `env:"MYAPP_PORT"`
        User       string         `env:"MYAPP_USER"`
        Users      []string       `env:"MYAPP_USERS"`
        Rate       float32        `env:"MYAPP_RATE"`
        Timeout    time.Duration  `env:"MYAPP_TIMEOUT"`
        ColorCodes map[string]int `env:"MYAPP_COLOR_CODES"`
    }

	results, err := envconfig.Unmarshal(&s)
	results.PrettyPrint() // it is nil safe
	// Env Variable         Type              OK
	// ----                 ----              ----
	// MYAPP_COLOR_CODES    map[string]int    v
	// MYAPP_DEBUG          bool              v
	// MYAPP_PORT           int               v
	// MYAPP_RATE           float32           v
	// MYAPP_TIMEOUT        time.Duration     v
	// MYAPP_USER           string            v
	// MYAPP_USERS          []string          v

    fmt.Printf("%+v", s)
    // Result:
    // {
    //  Debug:false 
    //  Port:8080 User:Kelsey 
    //  Users:[rob ken robert] 
    //  Rate:0.5 
    //  Timeout:3m0s 
    //  ColorCodes:map[blue:3 green:2 red:1]
    // }
}
```

## Struct Tag Support

Env variable name must be specified using `env` tag.
If struct contains field without `env` tag that will result in error.

Environment variables are required by default. 
If this is not desired one can use a `default` tag to specify value 
in case environment variable is not set.

For example, consider the following struct:

```Go
type Specification struct {
    Default         string `env:"MYAPP_DEFAULT" default:"foobar"`
    Foo             string `env:"MYAPP_FOO"`
}
```

If envconfig can't find an environment variable `MYAPP_DEFAULT`
it will populate it with "foobar" as a default value.

If envconfig can't find an environment variable `MYAPP_FOO` it will return an error.

## Supported Struct Field Types

envconfig supports these struct field types:

  * string
  * int8, int16, int32, int64
  * bool
  * float32, float64
  * slices of any supported type
  * maps (keys and values of any supported type)
  * [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler)
  * [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler)
  * [time.Duration](https://golang.org/pkg/time/#Duration)

Embedded structs using these fields are also supported.

## Custom Decoders

Any field whose type (or pointer-to-type) implements `envconfig.Decoder` can
control its own deserialization:

```Bash
export DNS_SERVER=8.8.8.8
```

```Go
type IPDecoder net.IP

func (ipd *IPDecoder) Decode(value string) error {
    *ipd = IPDecoder(net.ParseIP(value))
    return nil
}

type DNSConfig struct {
    Address IPDecoder `env:"DNS_SERVER"`
}
```

Also, envconfig will use a `Set(string) error` method like from the
[flag.Value](https://godoc.org/flag#Value) interface if implemented.
