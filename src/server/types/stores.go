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
	AddUserToTeam(userId, teamId string, role UserRole) error
	RemoveUserFromTeam(userId, teamId string) error
}

type InviteStore interface {
	CreateInvite(invite Invite) error
	GetInvite(id string) (*InviteInfo, error)
	GetInviteInfosFrom(from string) ([]InviteInfo, error)
	GetInviteInfosTo(to string) ([]InviteInfo, error)
	UpdateInviteStatus(id string, status InviteStatus) error
}
