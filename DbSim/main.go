package main

import (
	"github.com/dim4d/DbSim/terminal"
)

// var tb *typebox.TypeBox = typebox.NewTypeBox()

var ter terminal.Terminal = terminal.Terminal{}

func main() {

	str := "SET a INT 1"

	ter.Parser(str)
}
