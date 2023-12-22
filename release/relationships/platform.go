package relationships

import "github.com/gofrs/uuid"

type GameReleasePlatform struct {
	PlatformId    uuid.UUID
	GameReleaseId uuid.UUID
}
