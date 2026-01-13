package core

import (
	"strconv"
)

type Printable interface {
	ToString() string
}

func FormatValue(v interface{}) string {
	if p, ok := v.(Printable); ok {
		return p.ToString()
	}

	switch val := v.(type) {
	case int:
		return strconv.Itoa(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case string:
		return val
	default:
		return ""
	}
}
