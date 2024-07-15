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
	AddUserToTeam(execable interface{}, userId, teamId string, role UserRole) error
	RemoveUserFromTeam(userId, teamId string) error
}

type InviteStore interface {
	CreateInvite(invite Invite) error
	DeleteInvite(execable interface{}, id string) error
	GetInvite(id string) (*Invite, error)
	GetInviteInfosFrom(from string) ([]InviteInfo, error)
	GetInviteInfosTo(to string) ([]InviteInfo, error)
	UpdateInviteStatus(execable interface{}, id string, status InviteStatus) error
}

type VacationStore interface {
	CreateVacationRequest(execable interface{}, request VacationRequest) error
	GetVacationRequestsForUser(toUserId string) ([]VacationRequest, error)
	GetVacationRequestsFromUserId(requestedFromId string) ([]VacationRequest, error)
	UpdateVacationStatus(execable interface{}, requestId string, approverId string, status ApprovalStatus) error
	GetApprovalsForRequest(requestId string) ([]VacationApproval, error)
}
