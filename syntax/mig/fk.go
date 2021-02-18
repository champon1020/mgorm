package mig

import "github.com/champon1020/mgorm/syntax"

// FK is FOREIGN KEY clasue.
type FK struct {
	Columns []string
}

// Name returns clause keyword.
func (f *FK) Name() string {
	return "FOREIGN KEY"
}

// Build makes FOREIGN KEY clasue with syntax.StmtSet.
func (f *FK) Build() (*syntax.StmtSet, error) {
	ss := new(syntax.StmtSet)
	ss.WriteKeyword(f.Name())
	ss.WriteValue("(")
	for i, c := range f.Columns {
		if i > 0 {
			ss.WriteValue(",")
		}
		ss.WriteValue(c)
	}
	ss.WriteValue(")")
	return ss, nil
}
