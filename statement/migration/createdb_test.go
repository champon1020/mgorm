package migration_test

import (
	"testing"

	"github.com/champon1020/mgorm"
	"github.com/champon1020/mgorm/statement/migration"
	"github.com/stretchr/testify/assert"
)

func TestCreateDB_String(t *testing.T) {
	testCases := []struct {
		Stmt     *migration.CreateDBStmt
		Expected string
	}{
		{
			mgorm.CreateDB(migration.ExportedMySQLDB, "sample").(*migration.CreateDBStmt),
			`CREATE DATABASE sample`,
		},
		{
			mgorm.CreateDB(migration.ExportedPSQLDB, "sample").(*migration.CreateDBStmt),
			`CREATE DATABASE sample`,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.Stmt.String()
		errs := testCase.Stmt.ExportedGetErrors()
		if len(errs) > 0 {
			t.Errorf("Error was occurred: %v", errs[0])
			continue
		}
		assert.Equal(t, testCase.Expected, actual)
	}
}
