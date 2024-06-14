package team

import (
	"database/sql"
	"fmt"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetAllTeams() ([]types.Team, error) {
	rows, err := s.db.Query("SELECT * FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	teams := make([]types.Team, 0)
	for rows.Next() {
		t, err := readTeamData(rows)
		if err != nil {
			return nil, err
		}
		teams = append(teams, *t)
	}

	return teams, nil
}

func (s *Store) CreateTeam(team types.Team) error {
	_, err := s.db.Exec("INSERT INTO teams (Id, name) VALUES (?, ?)",
		team.Id, team.Name)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTeamByName(name string) (*types.Team, error) {
	rows, err := s.db.Query("SELECT * FROM teams WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	t := new(types.Team)
	for rows.Next() {
		t, err = readTeamData(rows)
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsValidUUID(t.Id) {
		return nil, fmt.Errorf("team not found")
	}

	return t, nil

}

func (s *Store) GetTeamById(id string) (*types.Team, error) {
	rows, err := s.db.Query("SELECT * FROM teams WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	t := new(types.Team)
	for rows.Next() {
		t, err = readTeamData(rows)
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsValidUUID(t.Id) {
		return nil, fmt.Errorf("team not found")
	}

	return t, nil
}

func readTeamData(rows *sql.Rows) (*types.Team, error) {
	team := new(types.Team)
	err := rows.Scan(
		&team.Id,
		&team.Name,
		&team.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return team, nil
}
