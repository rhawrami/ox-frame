package frame

type ColExpr struct {
	Name string
}

func Col(name string) ColExpr { return ColExpr{Name: name} }

func (c ColExpr) Eq(x any) ColExpr {
	return c
}

func (c ColExpr) Gt(x any) ColExpr {
	return c
}

func (c ColExpr) Lt(x any) ColExpr {
	return c
}

func (c ColExpr) Ge(x any) ColExpr {
	return c
}

func (c ColExpr) Le(x any) ColExpr {
	return c
}

func (c ColExpr) And(x any) ColExpr {
	return c
}

func (c ColExpr) Or(x any) ColExpr {
	return c
}

func (c ColExpr) Not(x any) ColExpr {
	return c
}
