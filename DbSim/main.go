package main

import (
	"fmt"

	"github.com/dim4d/DbSim/terminal"
)

// var tb *typebox.TypeBox = typebox.NewTypeBox()

var ter terminal.Terminal = terminal.Terminal{SetAtr: make(map[string]interface{})}

func main() {

	str := "SET b FLOAT 3.5"

	ter.Parser(str)

	fmt.Println(ter.SetAtr)

}
