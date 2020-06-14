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
	results.PrettyPrint() // safe if results are nil
	// Env Variable         Type              OK
	// ----                 ----              ----
	// MYAPP_COLOR_CODES    map[string]int    v
	// MYAPP_DEBUG          bool              v
	// MYAPP_PORT           int               v
	// MYAPP_RATE           float32           v
	// MYAPP_TIMEOUT        time.Duration     v
	// MYAPP_USER           string            v
	// MYAPP_USERS          []string          v

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", s)
}
