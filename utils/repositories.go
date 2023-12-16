package utils

import (
	"errors"
	"fmt"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var (
	unorderedQueryErr = errors.New("query must contain the order by clause to be paginated correctly")
	DataNotFoundErr   = errors.New("requested data was not found in the database")
	DuplicateDataErr  = errors.New("this data already exists")

	DefaultUuid uuid.UUID
)

// Paginate modifies completeQuery by appending required sql code to it to make pagination happen
// note that pageIndex is in "user" understandable format, as in it starts with 1
func Paginate(completeQuery string, pageIndex int, pageSize int) (query string, countQuery string, err error) {
	if !strings.Contains(completeQuery, "order by") {
		log.Warnf("Unable to add pagination to query %s as it does not contain an order by clause", completeQuery)
		return "", "", unorderedQueryErr
	}
	//todo: this breaks things, as it's expected for the page size to be as passed
	//return error or handle that on the higher level somewhere
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := pageSize * (pageIndex - 1)
	replaceRegex := regexp.MustCompile("select .* from")
	trimRegex := regexp.MustCompile("order by .*")
	return fmt.Sprintf("%s offset %d limit %d", completeQuery, offset, pageSize), trimRegex.ReplaceAllString(replaceRegex.ReplaceAllString(completeQuery, "select count(*) from"), ""), nil
}

func AppendWhereClause[T any](currentQuery string, columnName string, operand string, value T, isSet func(T) bool, positionalValues []any) (newQuery string, newPositional []any) {
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

func IsStringNotEmpty(value string) bool {
	return value != "" && value != "%%"
}

func IsUuidNotEmpty(value uuid.UUID) bool {
	var defaultUuid uuid.UUID
	return value != defaultUuid
}

func MakeLikeQuery(value string) string {
	return fmt.Sprintf("%%%s%%", value)
}

func ConvertIfNotFoundErr(err error) error {
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

func ConvertIfDuplicateErr(err error) error {
	var pgErr *pq.Error
	if err == nil || !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		log.Warnf("Expected error to be duplicate data, but failed to validate it as such: %s", err.Error())
		return err
	}
	return DuplicateDataErr
}

func ScanCountQuery(connector gotabase.Connector, query string, args ...interface{}) (int, error) {
	result, err := connector.QueryRow(query, args...)
	if err != nil {
		log.Warnf("Failed to run counting query: %s", err.Error())
		return -1, err
	}
	var count int
	if err := result.Scan(&count); err != nil {
		return -1, ConvertIfNotFoundErr(err)
	}
	return count, nil
}
