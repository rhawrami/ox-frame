package vector

import (
	"github.com/rhawrami/ox-frame/ox/dtype"
)

type Vector interface {
	Type() dtype.DataType
	Len() int
	NullCount() int
	IsNull(i int) bool
	IsNullBinary(i int) uint8
	DeepCopy() Vector
}
