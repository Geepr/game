package utils

import (
	"github.com/Geepr/game/mocks"
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
	}

	for _, data := range testData {
		currentTestData := data
		t.Run(currentTestData.query, func(t *testing.T) {
			_, countQuery, err := Paginate(currentTestData.query, 1, 1)
			mocks.AssertEquals(t, err, nil)
			mocks.AssertEquals(t, currentTestData.countQuery, countQuery)
		})
	}
}
