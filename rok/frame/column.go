package frame

import (
	"github.com/rhawrami/rok-frame/rok/dtype"
	"github.com/rhawrami/rok-frame/rok/vector"
)

type Column struct {
	Name  string
	DType dtype.DataType
	Vec   vector.Vector
}
