package models

import "github.com/gofrs/uuid"

// Platform defines a hardware that can be used to play a game.
// Think PCs, PlayStation consoles, ETC.
type Platform struct {
	Id uuid.UUID
	// Name is a full name of the platform.
	Name string
	// ShortName is a shortened Name, useful for display when there's less available space (IE: Sony PlayStation 5 == PS5).
	ShortName string
}
