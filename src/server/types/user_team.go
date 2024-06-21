package types

import "time"

type UserToTeamPayload struct {
	UserId   string   `json:"userId" validate:"required"`
	TeamId   string   `json:"teamId" validate:"required"`
	RoleType UserRole `json:"userRole" validate:"required"`
}

type UserRole int

const (
	Administrator UserRole = iota
	Member
)

type TeamUser struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	AddedAt  time.Time `json:"addedAt"`
	RoleType UserRole  `json:"userRole"`
}
