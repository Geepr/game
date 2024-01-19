package utils

import (
	"github.com/Geepr/game/mocks"
	"strings"
	"testing"
)

func TestPaginate_CountQueryReturnsCorrect(t *testing.T) {
	testData := []struct {
		query      string
		countQuery string
	}{
		{"select * from test_table order by a", "select count(*) from test_table"},
		{"select * from test_table where test_column = 1 order by a", "select count(*) from test_table where test_column = 1"},
		{"select a, b, c from test_table order by a", "select count(*) from test_table"},
		{"select * from test_table where from_column = 5 order by a", "select count(*) from test_table where from_column = 5"},
		{"select * from from_test_table order by a", "select count(*) from from_test_table"},
		{"select a, b, (select * from another_table where a = $1) from test_table where b = $2 order by a", "select count(*) from test_table where b = $2"},
		{"select a, b, (select * from another_table where a = $1 order by u) from test_table where b = $2 order by a", "select count(*) from test_table where b = $2"},
	}

	for _, data := range testData {
		currentTestData := data
		t.Run(currentTestData.query, func(t *testing.T) {
			resultQuery, countQuery, err := Paginate(currentTestData.query, 1, 1)
			mocks.AssertEquals(t, err, nil)
			mocks.AssertEquals(t, currentTestData.countQuery, countQuery)
			mocks.AssertEquals(t, strings.HasPrefix(resultQuery, currentTestData.query), true)
		})
	}
}

func TestAppendWhereClause_AppendedCorrectly(t *testing.T) {
	testData := []struct {
		query    string
		appended string
	}{
		{"select * from test_table", "select * from test_table where a = $1"},
		{"select a, b, c from test_table", "select a, b, c from test_table where a = $1"},
		{"select a, b, (select * from another_table where a = u) from test_table", "select a, b, (select * from another_table where a = u) from test_table where a = $1"},
	}

	for _, data := range testData {
		currentData := data
		t.Run(currentData.query, func(t *testing.T) {
			resultQuery, _ := AppendWhereClause(currentData.query, "a", "=", "b", func(s string) bool { return true }, []any{})
			mocks.AssertEquals(t, resultQuery, currentData.appended)
		})
	}
}
