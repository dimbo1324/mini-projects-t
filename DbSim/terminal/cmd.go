package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Terminal struct {
}

func (t Terminal) Parser(cmd string) {
	words := strings.Fields(cmd)

	if len(words) < 2 {
		return
	}

	token := words[0]

	switch token {
	case "SET":
		fmt.Println("SET")
	case "OBJECT":
		fmt.Println("OBJECT")
	case "PUSH":
		fmt.Println("PUSH")
	case "MERGE":
		fmt.Println("MERGE")
	case "PRINT":
		fmt.Println("PRINT")
	default:
		fmt.Println("Нет такого токена!")
		return
	}
}

func (t Terminal) InputCmd() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}

	line := scanner.Text()

	q, err := strconv.Atoi(line)

	if err != nil || q <= 0 {
		return
	}

	for i := 0; i < q; i++ {
		if !scanner.Scan() {
			break
		}
	}

}
