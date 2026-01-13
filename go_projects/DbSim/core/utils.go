package core

import "strconv"

func ParsePrimitive(typeStr, rawVal string) interface{} {
	switch typeStr {
	case "INT":
		v, _ := strconv.Atoi(rawVal)
		return v
	case "FLOAT":
		v, _ := strconv.ParseFloat(rawVal, 64)
		return v
	case "STRING":
		return rawVal
	default:
		return nil
	}
}
