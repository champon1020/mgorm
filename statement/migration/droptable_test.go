package migration_test

import (
	"testing"

	"github.com/champon1020/mgorm"
	"github.com/champon1020/mgorm/statement/migration"
	"github.com/stretchr/testify/assert"
)

func TestDropTable_String(t *testing.T) {
	testCases := []struct {
		Stmt     *migration.DropTableStmt
		Expected string
	}{
		{
			mgorm.DropTable(migration.ExportedMySQLDB, "person").(*migration.DropTableStmt),
			`DROP TABLE person`,
		},
		{
			mgorm.DropTable(migration.ExportedPSQLDB, "person").(*migration.DropTableStmt),
			`DROP TABLE person`,
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
