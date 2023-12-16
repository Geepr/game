package repositories

import (
	"fmt"
	"github.com/Geepr/game/models"
	"github.com/Geepr/game/utils"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

type PlatformRepository struct {
	connector gotabase.Connector
}

func NewPlatformRepository(connector gotabase.Connector) *PlatformRepository {
	return &PlatformRepository{connector: connector}
}

type PlatformSortOrder uint8

const (
	PlatformId PlatformSortOrder = iota
	PlatformName
	PlatformShortName
)

func (repo *PlatformRepository) GetPlatforms(nameQuery string, pageIndex int, pageSize int, order PlatformSortOrder) (*[]*models.Platform, error) {
	query := "select id, name, short_name from platforms"
	query, args := utils.AppendWhereClause(query, "name_normalised", "like", utils.MakeLikeQuery(strings.ToUpper(nameQuery)), utils.IsStringNotEmpty, []any{})
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, _, err := utils.Paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return repo.scanPlatforms(query, args...)
}

func (repo *PlatformRepository) GetPlatformById(id uuid.UUID) (*models.Platform, error) {
	query := "select id, name, short_name from platforms where id = $1"
	return repo.scanPlatform(query, id)
}

func (repo *PlatformRepository) AddPlatform(platform *models.Platform) error {
	query := "insert into platforms (name, short_name) VALUES ($1, $2) returning id"
	result, err := repo.connector.QueryRow(query, platform.Name, platform.ShortName)
	if err != nil {
		log.Warnf("Failed to execute insert query on platforms table: %s", err.Error())
		return utils.ConvertIfDuplicateErr(err)
	}
	if err = result.Scan(&platform.Id); err != nil {
		return err
	}
	return nil
}

func (repo *PlatformRepository) UpdatePlatform(id uuid.UUID, updatedPlatform *models.Platform) error {
	query := "update platforms set name = $2, short_name = $3 where id = $1"
	result, err := repo.connector.Exec(query, id, updatedPlatform.Name, updatedPlatform.ShortName)

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

func (repo *PlatformRepository) DeletePlatform(id uuid.UUID) error {
	query := "delete from platforms where id = $1"
	result, err := repo.connector.Exec(query, id)
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

func (repo *PlatformRepository) scanPlatforms(sql string, args ...interface{}) (*[]*models.Platform, error) {
	result, err := repo.connector.QueryRows(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on platforms table: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	platforms := make([]*models.Platform, 0)
	for result.Next() {
		platform, err := repo.scanRow(result)
		if err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}

	return &platforms, nil
}

func (repo *PlatformRepository) scanPlatform(sql string, args ...interface{}) (*models.Platform, error) {
	result, err := repo.connector.QueryRow(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on platforms table: %s", err.Error())
		return nil, err
	}
	return repo.scanRow(result)
}

func (repo *PlatformRepository) scanRow(row gotabase.Row) (*models.Platform, error) {
	platform := models.Platform{}
	if err := row.Scan(&platform.Id, &platform.Name, &platform.ShortName); err != nil {
		return nil, utils.ConvertIfNotFoundErr(err)
	}
	return &platform, nil
}

func (o PlatformSortOrder) getSqlColumnName() string {
	switch o {
	case PlatformId:
		return "id"
	case PlatformName:
		return "name"
	case PlatformShortName:
		return "short_name"
	}
	return "id"
}
