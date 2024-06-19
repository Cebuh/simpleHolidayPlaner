package types

import "time"

type Team struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type AddTeamPayload struct {
	Name string `json:"name" validate:"required"`
}

type RenameTeamPayload struct {
	Name string `json:"name" validate:"required"`
}
