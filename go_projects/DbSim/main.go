package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/dim4d/DbSim/core"
	"github.com/dim4d/DbSim/storage"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return
	}
	qStr := scanner.Text()
	q, _ := strconv.Atoi(qStr)

	tb := storage.NewTypeBox()

	for i := 0; i < q; i++ {
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]

		switch cmd {
		case "SET":
			if len(parts) < 4 {
				continue
			}
			key, typ, val := parts[1], parts[2], parts[3]
			tb.SetScalar(key, typ, val)

		case "OBJECT":
			if len(parts) < 3 {
				continue
			}
			key := parts[1]
			n, _ := strconv.Atoi(parts[2])

			newObj := core.NewObjectValue()

			for j := 0; j < n; j++ {
				scanner.Scan()
				fieldLine := scanner.Text()
				fParts := strings.Fields(fieldLine)
				if len(fParts) < 3 {
					continue
				}
				fName, fType, fVal := fParts[0], fParts[1], fParts[2]
				newObj.Data[fName] = core.ParsePrimitive(fType, fVal)
			}
			tb.SaveObject(key, newObj)

		case "PUSH":
			if len(parts) < 4 {
				continue
			}
			key, typ, val := parts[1], parts[2], parts[3]
			tb.PushValue(key, typ, val)

		case "MERGE":
			if len(parts) < 3 {
				continue
			}
			target, source := parts[1], parts[2]
			tb.MergeObjects(target, source)

		case "PRINT":
			if len(parts) < 2 {
				continue
			}
			key := parts[1]
			tb.PrintKey(key)
		}
	}
}
