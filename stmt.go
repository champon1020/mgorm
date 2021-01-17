package mgorm

import (
	"errors"

	"github.com/champon1020/mgorm/internal"
	"github.com/champon1020/mgorm/syntax"
)

// Op values for error handling.
const (
	OpStmtProcessQuerySQL internal.Op = "mgorm.Stmt.processQuerySQL"
	OpStmtProcessExecSQL  internal.Op = "mgorm.Stmt.processExecSQL"
)

// Stmt keeps the sql statement.
type Stmt struct {
	DB         syntax.DB
	Cmd        syntax.Cmd
	FromExpr   syntax.Expr
	ValuesExpr syntax.Expr
	SetExpr    syntax.Expr
	WhereExpr  syntax.Expr
	AndOr      []syntax.Expr
	Errors     []error
}

func (s *Stmt) addError(err error) {
	s.Errors = append(s.Errors, err)
}

// Query executes a query that returns some results.
func (s *Stmt) Query(model interface{}) error {
	sql, err := s.processQuerySQL()
	if err != nil {
		return err
	}
	if err := sql.doQuery(s.DB, model); err != nil {
		return err
	}
	return nil
}

func (s *Stmt) processQuerySQL() (SQL, error) {
	var sql SQL

	// Build SELECT.
	sel, ok := s.Cmd.(*syntax.Select)
	if !ok {
		err := errors.New("command must be SELECT")
		return "", internal.NewError(OpStmtProcessQuerySQL, internal.KindBasic, err)
	}
	sql.write(sel.Build().Build())

	// Build FROM.
	if s.FromExpr != nil {
		from, err := s.FromExpr.Build()
		if err != nil {
			return "", err
		}
		sql.write(from.Build())
	}

	// Build WHERE.
	if s.WhereExpr != nil {
		w, err := s.WhereExpr.Build()
		if err != nil {
			return "", err
		}
		sql.write(w.Build())
	}

	// Build AND or OR.
	if len(s.AndOr) > 0 {
		for _, e := range s.AndOr {
			ao, err := e.Build()
			if err != nil {
				return "", err
			}
			sql.write(ao.Build())
		}
	}
	return sql, nil
}

// Exec executes a query without returning any results.
func (s *Stmt) Exec() error {
	sql, err := s.processExecSQL()
	if err != nil {
		return err
	}
	if err := sql.doExec(s.DB); err != nil {
		return err
	}
	return nil
}

func (s *Stmt) processExecSQL() (SQL, error) {
	var sql SQL
	switch cmd := s.Cmd.(type) {
	case *syntax.Insert:
		sql.write(cmd.Build().Build())
		if s.ValuesExpr != nil {
			values, err := s.ValuesExpr.Build()
			if err != nil {
				return "", err
			}
			sql.write(values.Build())
		}
	case *syntax.Update:
		sql.write(cmd.Build().Build())
		if s.SetExpr != nil {
			set, err := s.SetExpr.Build()
			if err != nil {
				return "", err
			}
			sql.write(set.Build())
		}
	case *syntax.Delete:
		sql.write(cmd.Build().Build())
		if s.FromExpr != nil {
			from, err := s.FromExpr.Build()
			if err != nil {
				return "", err
			}
			sql.write(from.Build())
		}
	default:
		err := errors.New("command must be INSERT, UPDATE or DELETE")
		return "", internal.NewError(OpStmtProcessExecSQL, internal.KindType, err)

	}

	// Build WHERE.
	if s.WhereExpr != nil {
		w, err := s.WhereExpr.Build()
		if err != nil {
			return "", err
		}
		sql.write(w.Build())
	}

	// Build AND or OR.
	if len(s.AndOr) > 0 {
		for _, e := range s.AndOr {
			ao, err := e.Build()
			if err != nil {
				return "", err
			}
			sql.write(ao.Build())
		}
	}
	return sql, nil
}

// From calls FROM statement.
func (s *Stmt) From(tables ...string) *Stmt {
	s.FromExpr = syntax.NewFrom(tables)
	return s
}

// Values calls VALUES statement.
func (s *Stmt) Values(vals ...interface{}) *Stmt {
	s.ValuesExpr = syntax.NewValues(vals)
	return s
}

// Set calls SET statement.
func (s *Stmt) Set(vals ...interface{}) *Stmt {
	u, ok := s.Cmd.(*syntax.Update)
	if !ok {
		/* handle error */
	}
	set, err := syntax.NewSet(u.Columns, vals)
	if err != nil {
		s.addError(err)
		return s
	}
	s.SetExpr = set
	return s
}

// Where calls WHERE statement.
func (s *Stmt) Where(expr string, vals ...interface{}) *Stmt {
	w := syntax.NewWhere(expr, vals...)
	s.WhereExpr = w
	return s
}

// And calls AND statement.
func (s *Stmt) And(expr string, vals ...interface{}) *Stmt {
	w := syntax.NewAnd(expr, vals...)
	s.AndOr = append(s.AndOr, w)
	return s
}

// Or calls OR statement.
func (s *Stmt) Or(expr string, vals ...interface{}) *Stmt {
	w := syntax.NewOr(expr, vals...)
	s.AndOr = append(s.AndOr, w)
	return s
}
