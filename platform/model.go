package platform

import "github.com/gofrs/uuid"

// Platform defines a hardware that can be used to play a game.
// Think PCs, PlayStation consoles, ETC.
type Platform struct {
	Id uuid.UUID `json:"id"`
	// Name is a full name of the platform.
	Name string `json:"name"`
	// ShortName is a shortened Name, useful for display when there's less available space (IE: Sony PlayStation 5 == PS5).
	ShortName string `json:"shortName"`
}
