package repositories

import (
	"fmt"
	"github.com/Geepr/game/models"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	"strings"
)

type GameReleaseRepository struct {
	connector gotabase.Connector
}

func NewGameReleaseRepository(connector gotabase.Connector) *GameReleaseRepository {
	return &GameReleaseRepository{connector: connector}
}

type GameReleaseSortOrder uint8

const (
	GameReleaseId GameReleaseSortOrder = iota
	GameReleaseTitle
	GameReleaseDate
)

func (repo *GameReleaseRepository) GetGameReleases(titleQuery string, gameIdQuery uuid.UUID, pageIndex int, pageSize int, order GameReleaseSortOrder) (*[]*models.GameRelease, error) {
	query := "select id, game_id, title_override, description, release_date, release_date_unknown from game_releases"
	query, args := appendWhereClause(query, "title_override_normalised", "like", makeLikeQuery(strings.ToUpper(titleQuery)), isStringNotEmpty, []any{})
	query, args = appendWhereClause(query, "game_id", "=", gameIdQuery, isUuidNotEmpty, args)
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, err := paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return repo.scanGameReleases(query, args...)
}

func (repo *GameReleaseRepository) GetGameReleaseById(id uuid.UUID) (*models.GameRelease, error) {
	query := "select id, game_id, title_override, description, release_date, release_date_unknown from game_releases where id = $1"
	return repo.scanGameRelease(query, id)
}

func (repo *GameReleaseRepository) AddGameRelease(gameRelease *models.GameRelease) error {
	query := "insert into game_releases (game_id, title_override, description, release_date, release_date_unknown) VALUES  ($1, $2, $3, $4, $5) returning id"
	result, err := repo.connector.QueryRow(query, gameRelease.GameId, gameRelease.TitleOverride, gameRelease.Description, gameRelease.ReleaseDate, gameRelease.ReleaseDateUnknown)
	if err != nil {
		return convertIfNotFoundErr(err)
	}
	if err = result.Scan(&gameRelease.Id); err != nil {
		return err
	}
	return nil
}

func (repo *GameReleaseRepository) UpdateGameRelease(id uuid.UUID, updatedGameRelease *models.GameRelease) error {
	query := "update game_releases set title_override = $2, description = $3, release_date = $4, release_date_unknown = $5 where id = $1"
	result, err := repo.connector.Exec(query, id, updatedGameRelease.TitleOverride, updatedGameRelease.Description, updatedGameRelease.ReleaseDate, updatedGameRelease.ReleaseDateUnknown)

	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return DataNotFoundErr
	}
	return nil
}

func (repo *GameReleaseRepository) DeleteGameRelease(id uuid.UUID) error {
	query := "delete from game_releases where id = $1"
	result, err := repo.connector.Exec(query, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return DataNotFoundErr
	}
	return nil
}

func (repo *GameReleaseRepository) scanGameReleases(sql string, args ...interface{}) (*[]*models.GameRelease, error) {
	result, err := repo.connector.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	releases := make([]*models.GameRelease, 0)
	for result.Next() {
		release, err := repo.scanRow(result)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}

	return &releases, nil
}

func (repo *GameReleaseRepository) scanGameRelease(sql string, args ...interface{}) (*models.GameRelease, error) {
	result, err := repo.connector.QueryRow(sql, args...)
	if err != nil {
		return nil, err
	}
	return repo.scanRow(result)
}

func (repo *GameReleaseRepository) scanRow(row gotabase.Row) (*models.GameRelease, error) {
	release := models.GameRelease{}
	if err := row.Scan(&release.Id, &release.GameId, &release.TitleOverride, &release.Description, &release.ReleaseDate, &release.ReleaseDateUnknown); err != nil {
		return nil, convertIfNotFoundErr(err)
	}
	return &release, nil
}

func (o GameReleaseSortOrder) getSqlColumnName() string {
	switch o {
	case GameReleaseId:
		return "id"
	case GameReleaseTitle:
		return "title_override"
	case GameReleaseDate:
		return "release_date"
	}
	return "id"
}