package repositories

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
)

type GameReleasePlatformRepository struct {
	connector gotabase.Connector
}

func NewGameReleasePlatformRepository(connector gotabase.Connector) *GameReleasePlatformRepository {
	return &GameReleasePlatformRepository{connector: connector}
}

func (repo *GameReleasePlatformRepository) GetPlatformIdsForRelease(releaseId uuid.UUID) (*[]uuid.UUID, error) {
	query := "select platform_id from game_release_platforms where game_release_id = $1"
	return repo.scanUuids(query, releaseId)
}

func (repo *GameReleasePlatformRepository) GetReleaseIdsForPlatforms(platformId uuid.UUID) (*[]uuid.UUID, error) {
	query := "select game_release_id from game_release_platforms where platform_id = $1"
	return repo.scanUuids(query, platformId)
}

func (repo *GameReleasePlatformRepository) AddReleasePlatform(releaseId uuid.UUID, platformId uuid.UUID) error {
	_, err := repo.connector.Exec("insert into game_release_platforms (platform_id, game_release_id) VALUES ($1, $2)", platformId, releaseId)
	if err != nil {
		return convertIfNotFoundErr(convertIfDuplicateErr(err))
	}
	return nil
}

func (repo *GameReleaseRepository) RemoveReleasePlatform(releaseId uuid.UUID, platformId uuid.UUID) error {
	result, err := repo.connector.Exec("delete from game_release_platforms where platform_id = $1 and game_release_id = $2", platformId, releaseId)
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

func (repo *GameReleasePlatformRepository) scanUuids(sql string, args ...interface{}) (*[]uuid.UUID, error) {
	result, err := repo.connector.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	uuids := make([]uuid.UUID, 0)
	for result.Next() {
		var loadedUuid uuid.UUID
		if err = result.Scan(&loadedUuid); err != nil {
			return nil, err
		}
		uuids = append(uuids, loadedUuid)
	}
	return &uuids, nil
}