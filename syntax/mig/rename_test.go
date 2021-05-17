package mig_test

import (
	"testing"

	"github.com/champon1020/mgorm/syntax"
	"github.com/champon1020/mgorm/syntax/mig"
	"github.com/google/go-cmp/cmp"
)

func TestRename_Build(t *testing.T) {
	testCases := []struct {
		Rename   *mig.Rename
		Expected *syntax.StmtSet
	}{
		{
			&mig.Rename{Table: "table"},
			&syntax.StmtSet{Keyword: "RENAME TO", Value: "table"},
		},
	}

	for _, testCase := range testCases {
		actual, err := testCase.Rename.Build()
		if err != nil {
			t.Errorf("Error was occurred: %v", err)
			continue
		}
		if diff := cmp.Diff(testCase.Expected, actual); diff != "" {
			t.Errorf("Differs: (-want +got)\n%s", diff)
		}
	}
}
