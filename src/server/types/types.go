package types

import (
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	CreateUser(User) error
	GetUsersFromTeam(teamId string) ([]TeamUser, error)
}

type TeamStore interface {
	GetAllTeams() ([]Team, error)
	CreateTeam(Team) error
	RenameTeam(name, teamId string) error
	GetTeamById(id string) (*Team, error)
	GetTeamByName(name string) (*Team, error)
	AddUserToTeam(userId, teamId string) error
	RemoveUserFromTeam(userId, teamId string) error
}

type Team struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type TeamUser struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	AddedAt time.Time `json:"addedAt"`
}

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AddTeamPayload struct {
	Name string `json:"name" validate:"required"`
}

type UserToTeamPayload struct {
	UserId string `json:"userId" validate:"required"`
	TeamId string `json:"teamId" validate:"required"`
}

type RenameTeamPayload struct {
	Name string `json:"name" validate:"required"`
}
