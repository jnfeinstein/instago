// +build !heroku

package config

import (
	"fmt"
	"github.com/go-martini/martini"
)

func IsHeroku() bool {
	return false
}

func Initialize(m *martini.ClassicMartini) {
	fmt.Println("Running in debug environment")

	m.Use(martini.Static("private"))
	m.Use(martini.Static("tmp"))
}
