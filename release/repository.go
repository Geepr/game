package release

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
	SortByTitle
	SortByDate
)

func getGameReleases(titleQuery string, gameIdQuery uuid.UUID, pageIndex int, pageSize int, order SortOrder) ([]*GameRelease, error) {
	query := "select id, game_id, title_override, description, release_date, release_date_unknown from game_releases"
	//todo: this should probably fallback to the original game title query if override is null? - a view of some manner would be helpful here
	query, args := utils.AppendWhereClause(query, "title_override_normalised", "like", utils.MakeLikeQuery(strings.ToUpper(titleQuery)), utils.IsStringNotEmpty, []any{})
	query, args = utils.AppendWhereClause(query, "game_id", "=", gameIdQuery, utils.IsUuidNotEmpty, args)
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, _, err := utils.Paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return scanGameReleases(query, args...)
}

func getGameReleaseById(id uuid.UUID) (*GameRelease, error) {
	query := "select id, game_id, title_override, description, release_date, release_date_unknown from game_releases where id = $1"
	return scanGameRelease(query, id)
}

func addGameRelease(gameRelease *GameRelease) error {
	query := "insert into game_releases (game_id, title_override, description, release_date, release_date_unknown) VALUES  ($1, $2, $3, $4, $5) returning id"
	result, err := getConnector().QueryRow(query, gameRelease.GameId, gameRelease.TitleOverride, gameRelease.Description, gameRelease.ReleaseDate, gameRelease.ReleaseDateUnknown)
	if err != nil {
		return utils.ConvertIfNotFoundErr(err)
	}
	if err = result.Scan(&gameRelease.Id); err != nil {
		return err
	}
	return nil
}

func updateGameRelease(id uuid.UUID, updatedGameRelease *GameRelease) error {
	query := "update game_releases set title_override = $2, description = $3, release_date = $4, release_date_unknown = $5 where id = $1"
	result, err := getConnector().Exec(query, id, updatedGameRelease.TitleOverride, updatedGameRelease.Description, updatedGameRelease.ReleaseDate, updatedGameRelease.ReleaseDateUnknown)

	if err != nil {
		log.Warnf("Failed to execute update query on game releases: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to read affected rows when running update query on game releases: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func deleteGameRelease(id uuid.UUID) error {
	query := "delete from game_releases where id = $1"
	result, err := getConnector().Exec(query, id)
	if err != nil {
		log.Warnf("Failed to execute delete query on game releases: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to read affected rows when running delete query on game releases: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func scanGameReleases(sql string, args ...interface{}) ([]*GameRelease, error) {
	result, err := getConnector().QueryRows(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on game releases: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	releases := make([]*GameRelease, 0)
	for result.Next() {
		release, err := scanRow(result)
		if err != nil {
			return nil, err
		}
		releases = append(releases, release)
	}

	return releases, nil
}

func scanGameRelease(sql string, args ...interface{}) (*GameRelease, error) {
	result, err := getConnector().QueryRow(sql, args...)
	if err != nil {
		log.Warnf("Failed to run row query on game releases: %s", err.Error())
		return nil, err
	}
	return scanRow(result)
}

func scanRow(row gotabase.Row) (*GameRelease, error) {
	release := GameRelease{}
	if err := row.Scan(&release.Id, &release.GameId, &release.TitleOverride, &release.Description, &release.ReleaseDate, &release.ReleaseDateUnknown); err != nil {
		return nil, utils.ConvertIfNotFoundErr(err)
	}
	return &release, nil
}

func (o SortOrder) getSqlColumnName() string {
	switch o {
	case SortById:
		return "id"
	case SortByTitle:
		return "title_override"
	case SortByDate:
		return "release_date"
	}
	return "id"
}
