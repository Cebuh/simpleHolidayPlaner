package types

import "time"

type UserToTeamPayload struct {
	UserId   string `json:"userId" validate:"required"`
	TeamId   string `json:"teamId" validate:"required"`
	RoleType string `json:"roleType" validate:"required"`
}

const (
	Administrator = "admin"
	Member        = "member"
)

type TeamUser struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	AddedAt  time.Time `json:"addedAt"`
	RoleType string    `json:"role"`
}
