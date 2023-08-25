package repositories

import (
	"github.com/Geepr/game/mocks"
	"github.com/Geepr/game/models"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
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
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	id4, _ := uuid.NewV4()
	_, err := test.connection.Exec("insert into games (id, title, archived) values ($1, 'aaa', false), ($2, 'aab', false), ($3, 'cbb', true), ($4, 'def', false)", id1, id2, id3, id4)
	test.mockData = &[]models.Game{
		{
			Id:       id1,
			Title:    "aaa",
			Archived: false,
		},
		{
			Id:       id2,
			Title:    "aab",
			Archived: false,
		},
		{
			Id:       id3,
			Title:    "cbb",
			Archived: true,
		},
		{
			Id:       id4,
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

func TestGameRepository_GetGames_TitleQueryDefinedAndNotFound_ReturnsEmpty(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGames("definitely not found", 0, 100, GameId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 0)
}

func TestGameRepository_GetGameById_GameIdValid_GameReturned(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()

	for _, testCaseGlobal := range *test.mockData {
		testCase := testCaseGlobal
		t.Run(testCase.Id.String(), func(t *testing.T) {
			result, err := test.repo.GetGameById(testCase.Id)

			mocks.AssertDefault(t, err)
			mocks.AssertEquals(t, result.Id, testCase.Id)
			mocks.AssertEquals(t, result.Title, testCase.Title)
		})
	}
}

func TestGameRepository_GetGameById_GameIdNotFound_ReturnsSpecificError(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	testId, _ := uuid.NewV4()

	_, err := test.repo.GetGameById(testId)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}
