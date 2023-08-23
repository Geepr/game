package models

import (
	"errors"
	"github.com/gofrs/uuid"
)

var (
	InvalidGradeError = errors.New("this grade is invalid: it must be within [0, 100] range")
)

// Game represents an existing or future game, regardless of play status, platform, release date, etc.
type Game struct {
	Id uuid.UUID
	// Title is a common name of the game.
	// In cases when different platforms or releases differ in naming, those changes can be set per-release.
	Title       string
	Description *string
	// PersonalGrade is user-assignable grade for the game, ranging [0, 100]. Can be null when no grade has been set.
	// Grade is set per entire game, not per platform (as some outlets do), as this should represent the game in general, not per-platform flaws.
	PersonalGrade *uint8
	// Archived games are generally hidden from most views, but not removed outright.
	// This allows users to hide certain titles but keep the data for future reference.
	Archived bool
}

func (g *Game) SetGrade(grade uint8) error {
	if grade < 0 || grade > 100 {
		return InvalidGradeError
	}

	g.PersonalGrade = &grade
	return nil
}

func (g *Game) RemoveGrade() {
	g.PersonalGrade = nil
}
