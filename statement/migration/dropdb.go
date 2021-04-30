package migration

import (
	"github.com/champon1020/mgorm/domain"
	"github.com/champon1020/mgorm/internal"
	"github.com/champon1020/mgorm/syntax/mig"
)

// DropDBStmt is DROP DATABASE statement.
type DropDBStmt struct {
	migStmt
	cmd *mig.DropDB
}

func NewDropDBStmt(conn domain.Conn, dbName string) *DropDBStmt {
	stmt := &DropDBStmt{cmd: &mig.DropDB{DBName: dbName}}
	stmt.conn = conn
	return stmt
}

func (s *DropDBStmt) String() string {
	return s.string(s.buildSQL)
}

// Migrate executes database migration.
func (s *DropDBStmt) Migrate() error {
	return s.migration(s.buildSQL)
}

func (s *DropDBStmt) buildSQL(sql *internal.SQL) error {
	ss, err := s.cmd.Build()
	if err != nil {
		return err
	}
	sql.Write(ss.Build())
	return nil
}