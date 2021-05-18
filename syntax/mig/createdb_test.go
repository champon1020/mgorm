package mig_test

import (
	"testing"

	"github.com/champon1020/mgorm/syntax"
	"github.com/champon1020/mgorm/syntax/mig"
	"github.com/google/go-cmp/cmp"
)

func TestCreateDB_Build(t *testing.T) {
	testCases := []struct {
		CreateDB *mig.CreateDB
		Expected *syntax.StmtSet
	}{
		{
			&mig.CreateDB{DBName: "database"},
			&syntax.StmtSet{Keyword: "CREATE DATABASE", Value: "database"},
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.CreateDB.Build()
		if err != nil {
			t.Errorf("Error was occurred: %v", err)
			continue
		}
		if diff := cmp.Diff(testCase.Expected, actual); diff != "" {
			t.Errorf("Differs: (-want +got)\n%s", diff)
		}
	}
}