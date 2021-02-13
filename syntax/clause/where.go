package clause

import (
	"fmt"

	"github.com/champon1020/mgorm/internal"
	"github.com/champon1020/mgorm/syntax"
)

// Where is WHERE clause.
type Where struct {
	Expr   string
	Values []interface{}
}

// Name returns clause keyword.
func (w *Where) Name() string {
	return "WHERE"
}

// String returns function call with string.
func (w *Where) String() string {
	s := fmt.Sprintf("%q", w.Expr)
	if len(w.Values) > 0 {
		s += ", "
		s += internal.SliceToString(w.Values)
	}
	return fmt.Sprintf("%s(%s)", w.Name(), s)
}

// Build makes WHERE clause with syntax.StmtSet.
func (w *Where) Build() (*syntax.StmtSet, error) {
	ss, err := syntax.BuildStmtSetForExpression(w.Expr, w.Values...)
	if err != nil {
		return nil, err
	}
	ss.WriteKeyword(w.Name())
	return ss, nil
}