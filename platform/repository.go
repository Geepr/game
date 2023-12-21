package platform

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SortOrder uint8

const (
	SortById SortOrder = iota
	SortByName
	SortByShortName
)

func getPlatforms(nameQuery string, pageIndex int, pageSize int, order SortOrder) ([]*Platform, error) {
	query := "select id, name, short_name from platforms"
	query, args := utils.AppendWhereClause(query, "name_normalised", "like", utils.MakeLikeQuery(strings.ToUpper(nameQuery)), utils.IsStringNotEmpty, []any{})
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, _, err := utils.Paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return scanPlatforms(query, args...)
}

func getPlatformById(id uuid.UUID) (*Platform, error) {
	query := "select id, name, short_name from platforms where id = $1"
	return scanPlatform(query, id)
}

func addPlatform(platform *Platform) error {
	query := "insert into platforms (name, short_name) VALUES ($1, $2) returning id"
	result, err := getConnector().QueryRow(query, platform.Name, platform.ShortName)
	if err != nil {
		log.Warnf("Failed to execute insert query on platforms table: %s", err.Error())
		return utils.ConvertIfDuplicateErr(err)
	}
	if err = result.Scan(&platform.Id); err != nil {
		return err
	}
	return nil
}

func updatePlatform(id uuid.UUID, updatedPlatform *Platform) error {
	query := "update platforms set name = $2, short_name = $3 where id = $1"
	result, err := getConnector().Exec(query, id, updatedPlatform.Name, updatedPlatform.ShortName)

	if err != nil {
		return utils.ConvertIfDuplicateErr(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to get affected rows count when running update query on platgorms table: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func deletePlatform(id uuid.UUID) error {
	query := "delete from platforms where id = $1"
	result, err := getConnector().Exec(query, id)
	if err != nil {
		log.Warnf("Failed to execute delete query on platforms table: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to get affected rows count when running delete query on platforms table: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func scanPlatforms(sql string, args ...interface{}) ([]*Platform, error) {
	result, err := getConnector().QueryRows(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on platforms table: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	platforms := make([]*Platform, 0)
	for result.Next() {
		platform, err := scanRow(result)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}

	return platforms, nil
}

func scanPlatform(sql string, args ...interface{}) (*Platform, error) {
	result, err := getConnector().QueryRow(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on platforms table: %s", err.Error())
		return nil, err
	}
	return scanRow(result)
}

func scanRow(row gotabase.Row) (*Platform, error) {
	platform := Platform{}
	if err := row.Scan(&platform.Id, &platform.Name, &platform.ShortName); err != nil {
		return nil, utils.ConvertIfNotFoundErr(err)
	}
	return &platform, nil
}

func (o SortOrder) getSqlColumnName() string {
	switch o {
	case SortById:
		return "id"
	case SortByName:
		return "name"
	case SortByShortName:
		return "short_name"
	}
	return "id"
}
