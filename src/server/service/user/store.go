package user

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

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u := new(types.User)
	for rows.Next() {
		u, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsValidUUID(u.Id) {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil

}

func scanUserRow(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserById(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u := new(types.User)
	for rows.Next() {
		u, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsValidUUID(u.Id) {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (Id, name, email, password) VALUES(?, ?, ?, ?)",
		user.Id, user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	return nil
}
