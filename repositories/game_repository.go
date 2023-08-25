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

func (repo *GameRepository) AddGame(game *models.Game) error {
	query := "insert into games (title, description, archived) VALUES ($1, $2, $3) returning id"
	result, err := repo.connector.QueryRow(query, game.Title, game.Description, game.Archived)
	if err != nil {
		return err
	}
	if err = result.Scan(&game.Id); err != nil {
		return err
	}
	return nil
}

func (repo *GameRepository) UpdateGame(id uuid.UUID, updatedGame *models.Game) error {
	query := "update games set title = $2, description = $3, archived = $4 where id = $1"
	result, err := repo.connector.Exec(query, id, updatedGame.Title, updatedGame.Description, updatedGame.Archived)

	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return DataNotFoundErr
	}
	return nil
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
