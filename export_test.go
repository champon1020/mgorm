package mgorm

import "github.com/champon1020/mgorm/internal"

// Exported values which is declared in db.go.
func (db *DB) ExportedSetConn(conn sqlDB) {
	db.conn = conn
}

func (db *DB) ExportedSetDriver(driver internal.SQLDriver) {
	db.driver = driver
}

func (tx *Tx) ExportedSetConn(conn sqlTx) {
	tx.conn = conn
}

// Exported values which is declared in migration.go.
func (m *MigStmt) ExportedGetErrors() []error {
	return m.errors
}

// Exported values which is declared in mockdb.go.
var (
	CompareStmts = compareStmts
)

// Exported values which is declared in stmt.go.
var (
	SelectStmtProcessSQL = (*SelectStmt).processSQL
	InsertStmtProcessSQL = (*InsertStmt).processSQL
	UpdateStmtProcessSQL = (*UpdateStmt).processSQL
	DeleteStmtProcessSQL = (*DeleteStmt).processSQL
)

func (s *stmt) ExportedGetErrors() []error {
	return s.errors
}
