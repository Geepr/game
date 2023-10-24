package models

import (
	"github.com/gofrs/uuid"
)

// Game represents an existing or future game, regardless of play status, platform, release date, etc.
type Game struct {
	Id uuid.UUID `json:"id"`
	// Title is a common name of the game.
	// In cases when different platforms or releases differ in naming, those changes can be set per-release.
	Title       string  `json:"title"`
	Description *string `json:"description"`
	// Archived games are generally hidden from most views, but not removed outright.
	// This allows users to hide certain titles but keep the data for future reference.
	Archived bool `json:"archived"`
}
