package repositories

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	unorderedQueryErr = errors.New("query must contain the order by clause to be paginated correctly")
	DataNotFoundErr   = errors.New("requested data was not found in the database")
	DuplicateDataErr  = errors.New("this data already exists")

	DefaultUuid uuid.UUID
)

// paginate modifies completeQuery by appending required sql code to it to make pagination happen
// note that pageIndex is in "user" understandable format, as in it starts with 1
func paginate(completeQuery string, pageIndex int, pageSize int) (string, error) {
	if !strings.Contains(completeQuery, "order by") {
		log.Warnf("Unable to add pagination to query %s as it does not contain an order by clause", completeQuery)
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

func isUuidNotEmpty(value uuid.UUID) bool {
	var defaultUuid uuid.UUID
	return value != defaultUuid
}

func makeLikeQuery(value string) string {
	return fmt.Sprintf("%%%s%%", value)
}

func convertIfNotFoundErr(err error) error {
	if err.Error() == "sql: no rows in result set" {
		return DataNotFoundErr
	}
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return DataNotFoundErr
	}
	log.Warnf("Expected error to be data not found, but failed to validate it as such: %s", err.Error())
	return err
}

func convertIfDuplicateErr(err error) error {
	var pgErr *pq.Error
	if err == nil || !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		log.Warnf("Expected error to be duplicate data, but failed to validate it as such: %s", err.Error())
		return err
	}
	return DuplicateDataErr
}
