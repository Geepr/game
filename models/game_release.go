package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type GameRelease struct {
	Id     uuid.UUID
	GameId uuid.UUID
	// TitleOverride can be used if the title of the game is somehow modified for a specific release.
	TitleOverride string
	Description   *string
	// ReleaseDate indicates when this release was (or will be) published. Nil if not known.
	// Keep in mind that setting nil assumes that this is already released.
	// Use ReleaseDateUnknown to change that.
	ReleaseDate *time.Time
	// ReleaseDateUnknown indicates that no release date is currently known or available, but the game is not published yet and will be in the future.
	// Setting this field to true automatically assumes that the release is not public yet, even if ReleaseDate is set in the past.
	// If ReleaseDate is set when this field is true, it should be treated as an estimate instead.
	ReleaseDateUnknown bool
}

type GameReleasePlatform struct {
	PlatformId    uuid.UUID
	GameReleaseId uuid.UUID
}
