package relationships

import (
	"github.com/Geepr/game/utils"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

func getPlatformIdsForRelease(releaseId uuid.UUID) ([]uuid.UUID, error) {
	query := "select platform_id from game_release_platforms where game_release_id = $1"
	return scanUuids(query, releaseId)
}

func getReleaseIdsForPlatforms(platformId uuid.UUID) ([]uuid.UUID, error) {
	query := "select game_release_id from game_release_platforms where platform_id = $1"
	return scanUuids(query, platformId)
}

func addReleasePlatform(releaseId uuid.UUID, platformId uuid.UUID) error {
	_, err := getConnector().Exec("insert into game_release_platforms (platform_id, game_release_id) VALUES ($1, $2)", platformId, releaseId)
	if err != nil {
		log.Infof("Failed to add new release platform to database: %s", err.Error())
		return utils.ConvertIfNotFoundErr(utils.ConvertIfDuplicateErr(err))
	}
	return nil
}

func removeReleasePlatform(releaseId uuid.UUID, platformId uuid.UUID) error {
	result, err := getConnector().Exec("delete from game_release_platforms where platform_id = $1 and game_release_id = $2", platformId, releaseId)
	if err != nil {
		log.Warnf("Failed to execute query to remove release platform from the database: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to read affected rows after removing release platform: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func scanUuids(sql string, args ...interface{}) ([]uuid.UUID, error) {
	result, err := getConnector().QueryRows(sql, args...)
	if err != nil {
		log.Warnf("Failed to scan uuid when reading from release platforms table: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	uuids := make([]uuid.UUID, 0)
	for result.Next() {
		var loadedUuid uuid.UUID
		if err = result.Scan(&loadedUuid); err != nil {
			log.Warnf("Failed to scan uuid from release platforms table into a variable: %s", err.Error())
			return nil, err
		}
		uuids = append(uuids, loadedUuid)
	}
	return uuids, nil
}
