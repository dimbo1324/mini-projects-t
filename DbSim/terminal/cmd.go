package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Terminal struct {
	SetAtr map[string]interface{}
}

/*
* Функция-метод считается завершенной
 */
func (t *Terminal) setCmd(atr []string) {
	if len(atr) != 3 {
		return
	}

	name, varType, valueStr := atr[0], strings.ToUpper(atr[1]), atr[2]

	t.SetAtr["name"] = name
	switch varType {
	case "INT":
		val, err := strconv.ParseInt(valueStr, 10, 32)
		if err != nil {
			return
		}
		t.SetAtr["type"] = "INT"
		t.SetAtr["value"] = val

	case "FLOAT":
		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return
		}
		t.SetAtr["type"] = "FLOAT"
		t.SetAtr["value"] = val

	case "STRING":
		val := valueStr
		t.SetAtr["type"] = "STRING"
		t.SetAtr["value"] = val

	default:
		return
	}

	fmt.Printf("По результатам команды SET получены следующие данные: %s = %v (%s)\n",
		t.SetAtr["name"],
		t.SetAtr["value"],
		t.SetAtr["type"],
	)
}

func (t Terminal) stringPreparator(str string) string {
	sep, words := " ", strings.Fields(str)
	return strings.Join(words, sep)
}

func (t Terminal) Parser(cmd string) {
	words := strings.Fields(t.stringPreparator(cmd))

	if len(words) < 2 {
		return
	}

	token := strings.ToUpper(words[0])

	switch token {
	case "SET":
		t.setCmd(words[1:])
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
