package repositories

import (
	"github.com/Geepr/game/mocks"
	"github.com/Geepr/game/models"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	"testing"
	"time"
)

type gameReleaseRepoTest struct {
	connection gotabase.Connector
	repo       *GameReleaseRepository
	mockData   *[]models.GameRelease
	dbName     string
}

func newGameReleaseRepoTest(t *testing.T) *gameReleaseRepoTest {
	db, name := mocks.GetDatabase()
	test := &gameReleaseRepoTest{
		connection: db,
		repo:       NewGameReleaseRepository(db),
		dbName:     name,
	}
	t.Cleanup(test.cleanup)
	return test
}

func (test *gameReleaseRepoTest) cleanup() {
	mocks.DropDatabase(test.dbName)
}

func (test *gameReleaseRepoTest) insertMockData() {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	id4, _ := uuid.NewV4()
	_, err := test.connection.Exec("insert into games (id, title, archived) values ($1, 'aaa', false), ($2, 'aab', false), ($3, 'cbb', true), ($4, 'def', false)", id1, id2, id3, id4)
	mocks.PanicOnErr(err)
	_, err = test.connection.Exec("insert into game_releases (id, game_id, title_override, description, release_date, release_date_unknown) values"+
		"($1, $1, null, null, null, true),"+
		"($2, $1, 'other title', 'some description', '2023-01-01', true),"+
		"($3, $3, 'other title', null, '2023-01-01', true),"+
		"($4, $4, null, null, '2023-01-01', false)",
		id1, id2, id3, id4)
	mocks.PanicOnErr(err)
	title2, desc2 := "other title", "some description"
	release, _ := time.Parse(time.DateOnly, "2023-01-01")
	test.mockData = &[]models.GameRelease{
		{
			Id:                 id1,
			GameId:             id1,
			ReleaseDateUnknown: true,
		},
		{
			Id:                 id2,
			GameId:             id1,
			TitleOverride:      &title2,
			Description:        &desc2,
			ReleaseDate:        &release,
			ReleaseDateUnknown: true,
		},
		{
			Id:                 id3,
			GameId:             id3,
			TitleOverride:      &title2,
			ReleaseDate:        &release,
			ReleaseDateUnknown: true,
		},
		{
			Id:                 id4,
			GameId:             id4,
			ReleaseDate:        &release,
			ReleaseDateUnknown: false,
		},
	}
}

func (test *gameReleaseRepoTest) compareDates(date1 *time.Time, date2 *time.Time) bool {
	if date1 == nil {
		return date2 == nil
	}
	if date2 == nil {
		return date1 == nil
	}
	return date1.Year() == date2.Year() && date1.Month() == date2.Month() && date1.Day() == date2.Day()
}

func TestGameReleaseRepository_GetReleases_NoParametersSet_ReturnsAllReleases(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGameReleases("", DefaultUuid, 0, 100, GameReleaseId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 4)
	for _, release := range *test.mockData {
		mocks.AssertArrayContains(t, *result, func(value *models.GameRelease) bool {
			return release.GameId == value.GameId &&
				mocks.CompareNillable(release.TitleOverride, value.TitleOverride) &&
				mocks.CompareNillable(release.Description, value.Description) &&
				test.compareDates(release.ReleaseDate, value.ReleaseDate) &&
				release.ReleaseDateUnknown == value.ReleaseDateUnknown
		})
	}
}

func TestGameReleaseRepository_GetReleases_TitleAndGameIdQueryDefined_ReturnsMatching(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGameReleases("other", (*(test.mockData))[1].GameId, 0, 100, GameReleaseId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 1)
	single := (*(result))[0]
	mocks.AssertEquals(t, single.Id, (*(test.mockData))[1].Id)
}

func TestGameReleaseRepository_GetReleases_QueryDefinedAndNotFound_ReturnsEmpty(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()

	result, err := test.repo.GetGameReleases("definitely not found", DefaultUuid, 0, 100, GameReleaseId)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, *result, 0)
}

func TestGameReleaseRepository_GetGameReleaseById_ReleaseIdValid_ReleaseReturned(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()

	for _, testCaseGlobal := range *test.mockData {
		testCase := testCaseGlobal
		t.Run(testCase.Id.String(), func(t *testing.T) {
			result, err := test.repo.GetGameReleaseById(testCase.Id)

			mocks.AssertDefault(t, err)
			mocks.AssertEquals(t, result.Id, testCase.Id)
			mocks.AssertEquals(t, result.GameId, testCase.GameId)
			mocks.AssertEqualsNillable(t, result.TitleOverride, testCase.TitleOverride)
		})
	}
}

func TestGameReleaseRepository_GetGameReleaseById_IdNotFound_ReturnsSpecificError(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	testId, _ := uuid.NewV4()

	_, err := test.repo.GetGameReleaseById(testId)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameReleaseRepository_AddRelease_New_ReleaseAdded(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	release, _ := time.Parse(time.DateOnly, "2020-05-05")
	newRelease := models.GameRelease{
		GameId:             (*(test.mockData))[0].GameId,
		ReleaseDate:        &release,
		ReleaseDateUnknown: false,
	}

	err := test.repo.AddGameRelease(&newRelease)

	mocks.AssertDefault(t, err)
	mocks.AssertNotDefault(t, newRelease.Id)
}

func TestGameReleaseRepository_AddRelease_MissingGameId_NotFoundReturned(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	testId, _ := uuid.NewV4()
	newRelease := models.GameRelease{
		GameId:             testId,
		ReleaseDateUnknown: false,
	}

	err := test.repo.AddGameRelease(&newRelease)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameReleaseRepository_UpdateRelease_Exists_Updates(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	modified := (*(test.mockData))[1]
	modified.TitleOverride = nil
	desc := "new description"
	modified.Description = &desc
	modified.ReleaseDateUnknown = true
	modified.ReleaseDate = nil

	err := test.repo.UpdateGameRelease(modified.Id, &modified)

	mocks.AssertDefault(t, err)
	loaded, _ := test.repo.GetGameReleaseById(modified.Id)
	mocks.AssertEquals(t, loaded.Id, modified.Id)
	mocks.AssertEquals(t, loaded.TitleOverride, modified.TitleOverride)
	mocks.AssertEquals(t, *loaded.Description, *modified.Description)
	mocks.AssertEquals(t, loaded.ReleaseDate, modified.ReleaseDate)
	mocks.AssertEquals(t, loaded.ReleaseDateUnknown, modified.ReleaseDateUnknown)
}

func TestGameReleaseRepository_UpdateRelease_Missing_ReturnsNotFound(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()
	modified := (*(test.mockData))[0]

	err := test.repo.UpdateGameRelease(fakeId, &modified)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameReleaseRepository_DeleteRelease_Exists_Removes(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	toDelete := (*(test.mockData))[2]

	err := test.repo.DeleteGameRelease(toDelete.Id)

	mocks.AssertDefault(t, err)
	_, err = test.repo.GetGameReleaseById(toDelete.Id)
	mocks.AssertEquals(t, err, DataNotFoundErr)
}

func TestGameReleaseRepository_DeleteGameRelease_MissingId_ReturnsNotFound(t *testing.T) {
	test := newGameReleaseRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()

	err := test.repo.DeleteGameRelease(fakeId)

	mocks.AssertEquals(t, err, DataNotFoundErr)
}
