package core

import (
	"fmt"
	"sort"
	"strings"
)

type ObjectValue struct {
	Data map[string]interface{}
}

func NewObjectValue() ObjectValue {
	return ObjectValue{
		Data: make(map[string]interface{}),
	}
}

func (o ObjectValue) ToString() string {
	keys := make([]string, 0, len(o.Data))
	for k := range o.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		valStr := FormatValue(o.Data[k])
		parts = append(parts, fmt.Sprintf("%s:%s", k, valStr))
	}

	return "{" + strings.Join(parts, ",") + "}"
}

type ListValue struct {
	Data []interface{}
}

func (l ListValue) ToString() string {
	var parts []string
	for _, item := range l.Data {
		parts = append(parts, FormatValue(item))
	}
	return "[" + strings.Join(parts, ",") + "]"
}
