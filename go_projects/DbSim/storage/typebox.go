package storage

import (
	"fmt"

	"github.com/dim4d/DbSim/core"
)

type TypeBox struct {
	store map[string]interface{}
}

func NewTypeBox() *TypeBox {
	return &TypeBox{
		store: make(map[string]interface{}),
	}
}

func (tb *TypeBox) SetScalar(key, typ, raw string) {
	val := core.ParsePrimitive(typ, raw)
	tb.store[key] = val
}

func (tb *TypeBox) SaveObject(key string, obj core.ObjectValue) {
	tb.store[key] = obj
}

func (tb *TypeBox) PushValue(key, typ, raw string) {
	newVal := core.ParsePrimitive(typ, raw)
	existingVal, exists := tb.store[key]

	if !exists {
		tb.store[key] = core.ListValue{Data: []interface{}{newVal}}
	} else {
		if listVal, ok := existingVal.(core.ListValue); ok {
			listVal.Data = append(listVal.Data, newVal)
			tb.store[key] = listVal
		} else {
			newList := core.ListValue{
				Data: []interface{}{existingVal, newVal},
			}
			tb.store[key] = newList
		}
	}
}

func (tb *TypeBox) MergeObjects(targetKey, sourceKey string) {
	targetRaw, tExists := tb.store[targetKey]
	sourceRaw, sExists := tb.store[sourceKey]

	if !tExists || !sExists {
		return
	}

	targetObj, tOk := targetRaw.(core.ObjectValue)
	sourceObj, sOk := sourceRaw.(core.ObjectValue)

	if tOk && sOk {
		for k, v := range sourceObj.Data {
			targetObj.Data[k] = v
		}
		tb.store[targetKey] = targetObj
	}
}

func (tb *TypeBox) PrintKey(key string) {
	val, exists := tb.store[key]
	if !exists {
		fmt.Println("null")
		return
	}
	fmt.Println(core.FormatValue(val))
}
