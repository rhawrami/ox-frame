package frame

import (
	"github.com/rhawrami/ox-frame/ox/dtype"
	"github.com/rhawrami/ox-frame/ox/vector"
)

type Column struct {
	Name  string
	DType dtype.DataType
	Vec   vector.Vector
}
