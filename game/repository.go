package game

import (
	"fmt"
	"github.com/Geepr/game/utils"
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SortOrder uint8

const (
	SortById SortOrder = iota
	SortByTitle
)

func getGames(titleQuery string, pageIndex int, pageSize int, order SortOrder) ([]*Game, int, error) {
	query := "select id, title, description, archived from games"
	query, args := utils.AppendWhereClause(query, "title_normalised", "like", utils.MakeLikeQuery(strings.ToUpper(titleQuery)), utils.IsStringNotEmpty, []any{})
	query += fmt.Sprintf(" order by %s", order.getSqlColumnName())
	query, countQuery, err := utils.Paginate(query, pageIndex, pageSize)
	if err != nil {
		return nil, 0, err
	}
	countResults, err := utils.ScanCountQuery(getConnector(), countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	results, err := scanGames(query, args...)
	return results, countResults, err
}

func getGameById(id uuid.UUID) (*Game, error) {
	query := "select id, title, description, archived from games where id = $1"
	return scanGame(query, id)
}

func addGame(game *Game) error {
	query := "insert into games (title, description, archived) VALUES ($1, $2, $3) returning id"
	result, err := getConnector().QueryRow(query, game.Title, game.Description, game.Archived)
	if err != nil {
		log.Warnf("Failed to execute insert query on games table: %s", err.Error())
		return err
	}
	if err = result.Scan(&game.Id); err != nil {
		return err
	}
	return nil
}

func updateGame(id uuid.UUID, updatedGame *Game) error {
	query := "update games set title = $2, description = $3, archived = $4 where id = $1"
	result, err := getConnector().Exec(query, id, updatedGame.Title, updatedGame.Description, updatedGame.Archived)

	if err != nil {
		log.Warnf("Failed to execute update query on games table: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to get affected rows count when running update query on games table: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func deleteGame(id uuid.UUID) error {
	query := "delete from games where id = $1"
	result, err := getConnector().Exec(query, id)
	if err != nil {
		log.Warnf("Failed to execute delete query on games table: %s", err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		log.Warnf("Failed to get affected rows count when running deleteRoute query on games table: %s", err.Error())
		return err
	}
	if affected != 1 {
		return utils.DataNotFoundErr
	}
	return nil
}

func scanGames(sql string, args ...interface{}) ([]*Game, error) {
	result, err := getConnector().QueryRows(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on games table: %s", err.Error())
		return nil, err
	}
	defer result.Close()

	games := make([]*Game, 0)
	for result.Next() {
		game, err := scanRow(result)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}

	return games, nil
}

func scanGame(sql string, args ...interface{}) (*Game, error) {
	result, err := getConnector().QueryRow(sql, args...)
	if err != nil {
		log.Warnf("Failed to run query on games table: %s", err.Error())
		return nil, err
	}
	return scanRow(result)
}

func scanRow(row gotabase.Row) (*Game, error) {
	game := Game{}
	if err := row.Scan(&game.Id, &game.Title, &game.Description, &game.Archived); err != nil {
		return nil, utils.ConvertIfNotFoundErr(err)
	}
	return &game, nil
}

func (o SortOrder) getSqlColumnName() string {
	switch o {
	case SortById:
		return "id"
	case SortByTitle:
		return "title"
	}
	return "id"
}
