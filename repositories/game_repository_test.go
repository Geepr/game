package repositories

import (
	"github.com/Geepr/game/mocks"
	"github.com/Geepr/game/models"
	"github.com/KowalskiPiotr98/gotabase"
	"strings"
	"testing"
)

type gameRepoTest struct {
	connection gotabase.Connector
	repo       *GameRepository
	mockData   *[]models.Game
	dbName     string
}

func newGameRepoTest(t *testing.T) *gameRepoTest {
	db, name := mocks.GetDatabase()
	test := &gameRepoTest{
		connection: db,
		repo:       NewGameRepository(db),
		dbName:     name,
	}
	t.Cleanup(test.cleanup)
	return test
}

func (test *gameRepoTest) cleanup() {
	mocks.DropDatabase(test.dbName)
}

func (test *gameRepoTest) insertMockData() {
	_, err := test.connection.Exec("insert into games (title, archived) values ('aaa', false), ('aab', false), ('cbb', true), ('def', false)")
	test.mockData = &[]models.Game{
		{
			Title:    "aaa",
			Archived: false,
		},
		{
			Title:    "aab",
			Archived: false,
		},
		{
			Title:    "cbb",
			Archived: true,
		},
		{
			Title:    "def",
			Archived: false,
		},
	}
	mocks.PanicOnErr(err)
}

func TestGameRepository_GetGames_NoParametersSet_ReturnsAllGames(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGames("", 0, 100, GameId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 4)
	for _, game := range *test.mockData {
		mocks.AssertArrayContains(t, *result, func(value *models.Game) bool {
			return value.Title == game.Title && value.Archived == game.Archived
		})
	}
}

func TestGameRepository_GetGames_TitleQueryDefined_ReturnsMatching(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGames("Aa", 0, 100, GameId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 2)
	for _, game := range *test.mockData {
		if !strings.Contains(game.Title, "aa") {
			continue
		}
		mocks.AssertArrayContains(t, *result, func(value *models.Game) bool {
			return value.Title == game.Title && value.Archived == game.Archived
		})
	}
}
