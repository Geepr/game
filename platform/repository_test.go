package platform

import (
	"github.com/Geepr/game/mocks"
	"github.com/Geepr/game/utils"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	"strings"
	"testing"
)

type platformRepoTest struct {
	connection gotabase.Connector
	mockData   []*Platform
	dbName     string
}

func newPlatformRepoTest(t *testing.T) *platformRepoTest {
	db, name := mocks.GetDatabase()
	test := &platformRepoTest{
		connection: db,
		dbName:     name,
	}
	getConnector = func() gotabase.Connector { return db }
	t.Cleanup(test.cleanup)
	return test
}

func (test *platformRepoTest) cleanup() {
	mocks.DropDatabase(test.dbName)
}

func (test *platformRepoTest) insertMockData() {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	id4, _ := uuid.NewV4()
	_, err := test.connection.Exec("insert into platforms (id, name, short_name) values ($1, 'aaa', 'aa'), ($2, 'aab', 'ab'), ($3, 'cbb', 'cb'), ($4, 'def', 'de')", id1, id2, id3, id4)
	test.mockData = []*Platform{
		{
			Id:        id1,
			Name:      "aaa",
			ShortName: "aa",
		},
		{
			Id:        id2,
			Name:      "aab",
			ShortName: "ab",
		},
		{
			Id:        id3,
			Name:      "cbb",
			ShortName: "cb",
		},
		{
			Id:        id4,
			Name:      "def",
			ShortName: "de",
		},
	}
	mocks.PanicOnErr(err)
}

func TestPlatformRepository_GetPlatforms_NoParametersSet_ReturnsAllPlatforms(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()

	result, items, err := getPlatforms("", 0, 100, SortById)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, result, 4)
	for _, game := range test.mockData {
		mocks.AssertArrayContains(t, result, func(value *Platform) bool {
			return value.Name == game.Name && value.ShortName == game.ShortName
		})
	}
	mocks.AssertEquals(t, items, 4)
}

func TestPlatformRepository_GetPlatforms_NameQueryDefined_ReturnsMatching(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()

	result, items, err := getPlatforms("Aa", 0, 100, SortById)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, result, 2)
	for _, platform := range test.mockData {
		if !strings.Contains(platform.Name, "aa") {
			continue
		}
		mocks.AssertArrayContains(t, result, func(value *Platform) bool {
			return value.Name == platform.Name && value.ShortName == platform.ShortName
		})
	}
	mocks.AssertEquals(t, items, 2)
}

func TestPlatformRepository_GetPlatforms_TitleQueryDefinedAndNotFound_ReturnsEmpty(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()

	result, items, err := getPlatforms("definitely not found", 0, 100, SortById)

	mocks.AssertDefault(t, err)
	mocks.AssertCountEqual(t, result, 0)
	mocks.AssertEquals(t, items, 0)
}

func TestPlatformRepository_GetPlatformById_ValidId_FoundAndReturned(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()

	for _, testCaseGlobal := range test.mockData {
		testCase := testCaseGlobal
		t.Run(testCase.Id.String(), func(t *testing.T) {
			result, err := getPlatformById(testCase.Id)

			mocks.AssertDefault(t, err)
			mocks.AssertEquals(t, result.Id, testCase.Id)
			mocks.AssertEquals(t, result.Name, testCase.Name)
			mocks.AssertEquals(t, result.ShortName, testCase.ShortName)
		})
	}
}

func TestPlatformRepository_GetPlatformById_PlatformIdNotFound_ReturnsSpecificError(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	testId, _ := uuid.NewV4()

	_, err := getPlatformById(testId)

	mocks.AssertEquals(t, err, utils.DataNotFoundErr)
}

func TestPlatformRepository_AddPlatform_ValidNewPlatform_PlatformAdded(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	newPlatform := Platform{
		Name:      "test",
		ShortName: "test",
	}

	err := addPlatform(&newPlatform)

	mocks.AssertDefault(t, err)
	mocks.AssertNotDefault(t, newPlatform.Id)
}

func TestPlatformRepository_AddPlatform_DuplicateName_ErrorReturned(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	toDuplicate := test.mockData[1]
	duplicate := Platform{
		Name:      toDuplicate.Name,
		ShortName: toDuplicate.ShortName,
	}

	err := addPlatform(&duplicate)

	mocks.AssertEquals(t, err, utils.DuplicateDataErr)
}

func TestPlatformRepository_UpdatePlatform_PlatformExists_Updates(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	modified := test.mockData[0]
	modified.Name = "new name"
	modified.ShortName = "nn"

	err := updatePlatform(modified.Id, modified)

	mocks.AssertDefault(t, err)
	loaded, _ := getPlatformById(modified.Id)
	mocks.AssertEquals(t, loaded.Id, modified.Id)
	mocks.AssertEquals(t, loaded.Name, modified.Name)
	mocks.AssertEquals(t, loaded.ShortName, modified.ShortName)
}

func TestPlatformRepository_UpdatePlatform_NewNameDuplicate_ReturnsErr(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	modified := test.mockData[0]
	toDuplicate := test.mockData[1]
	modified.Name = toDuplicate.Name
	modified.ShortName = toDuplicate.ShortName

	err := updatePlatform(modified.Id, modified)

	mocks.AssertEquals(t, err, utils.DuplicateDataErr)
}

func TestPlatformRepository_UpdatePlatform_PlatformMissing_ReturnsNotFound(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()
	modified := test.mockData[0]

	err := updatePlatform(fakeId, modified)

	mocks.AssertEquals(t, err, utils.DataNotFoundErr)
}

func TestPlatformRepository_DeletePlatform_PlatformExists_RemovesPlatform(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	toDelete := test.mockData[2]

	err := deletePlatform(toDelete.Id)

	mocks.AssertDefault(t, err)
	_, err = getPlatformById(toDelete.Id)
	mocks.AssertEquals(t, err, utils.DataNotFoundErr)
}

func TestPlatformRepository_DeletePlatform_MissingId_ReturnsNotFound(t *testing.T) {
	test := newPlatformRepoTest(t)
	test.insertMockData()
	fakeId, _ := uuid.NewV4()

	err := deletePlatform(fakeId)

	mocks.AssertEquals(t, err, utils.DataNotFoundErr)
}
