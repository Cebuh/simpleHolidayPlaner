package types

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
	AddUserToTeam(userId, teamId, roleType string) error
	RemoveUserFromTeam(userId, teamId string) error
}
