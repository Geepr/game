package repositories

import (
	"errors"
	"fmt"
	"strings"
)

var (
	unorderedQueryErr = errors.New("query must contain the order by clause to be paginated correctly")
)

// paginate modifies completeQuery by appending required sql code to it to make pagination happen
// note that pageIndex is in "user" understandable format, as in it starts with 1
func paginate(completeQuery string, pageIndex int, pageSize int) (string, error) {
	if !strings.Contains(completeQuery, "order by") {
		return "", unorderedQueryErr
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := pageSize * (pageIndex - 1)
	return fmt.Sprintf("%s offset %d limit %d", completeQuery, offset, pageSize), nil
}

func appendWhereClause[T any](currentQuery string, columnName string, operand string, value T, isSet func(T) bool, positionalValues []any) (newQuery string, newPositional []any) {
	if !isSet(value) {
		return currentQuery, positionalValues
	}

	combiningWord := "where"
	if strings.Contains(currentQuery, "where") {
		combiningWord = "and"
	}

	newQuery = fmt.Sprintf("%s %s %s %s $%d", currentQuery, combiningWord, columnName, operand, len(positionalValues)+1)
	newPositional = append(positionalValues, value)
	return
}

func isStringNotEmpty(value string) bool {
	return value != "" && value != "%%"
}

func makeLikeQuery(value string) string {
	return fmt.Sprintf("%%%s%%", value)
}
