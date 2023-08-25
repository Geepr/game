package repositories

import (
	"fmt"
	"github.com/Geepr/game/models"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	"strings"
)

type GameRepository struct {
	connector gotabase.Connector
}

func NewGameRepository(connector gotabase.Connector) *GameRepository {
	return &GameRepository{connector: connector}
}

type GameSortOrder uint8

const (
	GameId GameSortOrder = iota
	GameTitle
)

func (repo *GameRepository) GetGames(titleQuery string, pageIndex int, pageSize int, order GameSortOrder) (*[]*models.Game, error) {
	query := "select id, title, description, archived from games"
	query, args := appendWhereClause(query, "title_normalised", "like", makeLikeQuery(strings.ToUpper(titleQuery)), isStringNotEmpty, []any{})
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, err := paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, err
	}
	return repo.scanGames(query, args...)
}

func (repo *GameRepository) GetGameById(id uuid.UUID) (*models.Game, error) {
	query := "select id, title, description, archived from games where id = $1"
	return repo.scanGame(query, id)
}

func (repo *GameRepository) scanGames(sql string, args ...interface{}) (*[]*models.Game, error) {
	result, err := repo.connector.QueryRows(sql, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	games := make([]*models.Game, 0)
	for result.Next() {
		game, err := repo.scanRow(result)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return &games, nil
}

func (repo *GameRepository) scanGame(sql string, args ...interface{}) (*models.Game, error) {
	result, err := repo.connector.QueryRow(sql, args...)
	if err != nil {
		return nil, err
	}
	return repo.scanRow(result)
}

func (repo *GameRepository) scanRow(row gotabase.Row) (*models.Game, error) {
	game := models.Game{}
	if err := row.Scan(&game.Id, &game.Title, &game.Description, &game.Archived); err != nil {
		return nil, convertIfNotFoundErr(err)
	}
	return &game, nil
}

func (o GameSortOrder) getSqlColumnName() string {
	switch o {
	case GameId:
		return "id"
	case GameTitle:
		return "title"
	}
	return "id"
}
