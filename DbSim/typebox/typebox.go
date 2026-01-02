package typebox

type TypeBox struct {
	store map[string]interface{}
}

func NewTypeBox() *TypeBox {
	return &TypeBox{
		store: make(map[string]interface{}),
	}
}
