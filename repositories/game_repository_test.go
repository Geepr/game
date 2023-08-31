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

func TestGameRepository_AddGame_NewName_GameAdded(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	newGame := models.Game{
		Title: "totally new and unique title",
	}

	err := test.repo.AddGame(&newGame)

	mocks.AssertDefault(t, err)
	mocks.AssertNotDefault(t, newGame.Id)
}

func TestGameRepository_UpdateGame_GameExists_Updates(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	modified := (*(test.mockData))[0]
	modified.Title = "new title"
	desc := "new description"
	modified.Description = &desc
	modified.Archived = true

	err := test.repo.UpdateGame(modified.Id, &modified)

	mocks.AssertDefault(t, err)
	loaded, _ := test.repo.GetGameById(modified.Id)
	mocks.AssertEquals(t, loaded.Id, modified.Id)
	mocks.AssertEquals(t, loaded.Title, modified.Title)
	mocks.AssertEquals(t, *loaded.Description, *modified.Description)
	mocks.AssertEquals(t, loaded.Archived, modified.Archived)
}

func TestGameRepository_UpdateGame_GameMissing_ReturnsNotFound(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()
	modified := (*(test.mockData))[0]

	err := test.repo.UpdateGame(fakeId, &modified)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameRepository_DeleteGame_GameExists_RemovesGame(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	toDelete := (*(test.mockData))[2]

	err := test.repo.DeleteGame(toDelete.Id)

	mocks.AssertDefault(t, err)
	_, err = test.repo.GetGameById(toDelete.Id)
	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameRepository_DeleteGame_MissingId_ReturnsNotFound(t *testing.T) {
	test := newGameRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()

	err := test.repo.DeleteGame(fakeId)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}
