package database_test

import (
	"testing"

	"github.com/champon1020/mgorm"
	"github.com/champon1020/mgorm/database"
	"github.com/google/go-cmp/cmp"
	"github.com/morikuni/failure"
)

func TestMock_Expectation(t *testing.T) {
	expectedReturn := []int{10, 20, 30}

	// Test phase.
	mock := database.NewMockDB("")
	mock.Expect(mgorm.Insert(nil, "table", "column1", "column2").Values(10, "str"))
	mock.ExpectWithReturn(mgorm.Select(nil, "column1").From("table"), expectedReturn)

	// Actual process.
	if err := mgorm.Insert(mock, "table", "column1", "column2").Values(10, "str").Exec(); err != nil {
		t.Errorf("Error was occurred: %+v", err)
		return
	}
	model := new([]int)
	if err := mgorm.Select(mock, "column1").From("table").Query(model); err != nil {
		t.Errorf("Error was occurred: %+v", err)
		return
	}

	// Test phase.
	if err := mock.Complete(); err != nil {
		t.Errorf("Error was occurred: %+v", err)
		return
	}

	// Validate model values.
	if diff := cmp.Diff(*model, expectedReturn); diff != "" {
		t.Errorf("Differs: (-want +got)\n%s", diff)
	}
}
func TestMockDB_Begin_Fail(t *testing.T) {
	{
		expectedErr := database.ErrInvalidMockExpectation

		// Test phase.
		mock := new(database.ExportedMockDB)

		// Actual process.
		_, err := mock.Begin()

		// Validate if the expected error was occurred.
		if !failure.Is(err, expectedErr) {
			t.Errorf("Different error was occurred")
			t.Errorf("  Expected: %+v", expectedErr)
			t.Errorf("  Actual:   %+v", err)
		}
	}
	{
		expectedErr := database.ErrInvalidMockExpectation

		// Test phase.
		mock := new(database.ExportedMockDB)
		mock.Expect(mgorm.Select(nil, "column1").From("table"))
		_ = mock.ExpectBegin()

		// Actual process.
		_, err := mock.Begin()

		// Validate if the expected error was occurred.
		if !failure.Is(err, expectedErr) {
			t.Errorf("Different error was occurred")
			t.Errorf("  Expected: %+v", expectedErr)
			t.Errorf("  Actual:   %+v", err)
		}
	}
}

func TestMockDB_Complete_Fail(t *testing.T) {
	expectedErr := database.ErrInvalidMockExpectation

	// Test phase.
	mock := database.NewMockDB("")
	mock.Expect(mgorm.Insert(nil, "table1", "column1", "column2").Values(10, "str"))
	mock.Expect(mgorm.Insert(nil, "table2", "column1", "column2").Values(10, "str"))

	// Actual process.
	if err := mgorm.Insert(mock, "table1", "column1", "column2").Values(10, "str").Exec(); err != nil {
		t.Errorf("Error was occurred: %+v", err)
		return
	}

	// Test phase.
	err := mock.Complete()

	// Validate if the expected error was occurred.
	if !failure.Is(err, expectedErr) {
		t.Errorf("Different error was occurred")
		t.Errorf("  Expected: %+v", expectedErr)
		t.Errorf("  Actual:   %+v", err)
	}
}

func TestMockDB_Complete_Transaction_Fail(t *testing.T) {
	expectedErr := database.ErrInvalidMockExpectation

	// Test phase.
	mock := database.NewMockDB("")
	mocktx := mock.ExpectBegin()
	mocktx.Expect(mgorm.Insert(nil, "table1", "column1", "column2").Values(10, "str"))
	mocktx.Expect(mgorm.Insert(nil, "table2", "column1", "column2").Values(10, "str"))

	// Actual process.
	tx, err := mock.Begin()
	if err != nil {
		t.Errorf("Error was occured: %+v", err)
		return
	}
	if err = mgorm.Insert(tx, "table1", "column1", "column2").Values(10, "str").Exec(); err != nil {
		t.Errorf("Error was occurred: %+v", err)
		return
	}

	// Test phase.
	err = mock.Complete()

	// Validate if the expected error was occurred.
	if !failure.Is(err, expectedErr) {
		t.Errorf("Different error was occurred")
		t.Errorf("  Expected: %+v", expectedErr)
		t.Errorf("  Actual:   %+v", err)
	}
}

func TestMockDB_CompareWith(t *testing.T) {
	{
		expectedErr := database.ErrInvalidMockExpectation

		// Test phase.
		mock := database.NewMockDB("")

		// Actual process.
		model := new([]int)
		err := mgorm.Select(mock, "column1").From("table").Query(model)

		// Validate if the expected error was occurred.
		if !failure.Is(err, expectedErr) {
			t.Errorf("Different error was occurred")
			t.Errorf("  Expected: %+v", expectedErr)
			t.Errorf("  Actual:   %+v", err)
		}
	}
	{
		expectedErr := database.ErrInvalidMockExpectation

		// Test phase.
		mock := database.NewMockDB("")
		_ = mock.ExpectBegin()

		// Actual process.
		model := new([]int)
		err := mgorm.Select(mock, "column1").From("table").Query(model)

		// Validate if the expected error was occurred.
		if !failure.Is(err, expectedErr) {
			t.Errorf("Different error was occurred")
			t.Errorf("  Expected: %+v", expectedErr)
			t.Errorf("  Actual:   %+v", err)
		}
	}
}