package frame

import (
	"github.com/rhawrami/uz-frame/uz/dtype"
	"github.com/rhawrami/uz-frame/uz/vector"
)

type Column struct {
	Name  string
	DType dtype.DataType
	Vec   vector.Vector
}
