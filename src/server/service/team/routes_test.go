package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestTeamServiceHandlers(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByNameMock = func(name string) (*types.Team, error) { return nil, fmt.Errorf("Not found") }
	teamStore.CreateTeamMock = func(t types.Team) error { return nil }

	userStore := &mockUser{}
	handler := NewHandler(db, teamStore, userStore)

	t.Run("should run if team is created",
		func(t *testing.T) {
			payload := types.AddTeamPayload{
				Name: "Team A",
			}

			marshalled, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatal(err)
			}

			testHttp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/teams", handler.handleAddTeam).Methods(http.MethodPost)
			router.ServeHTTP(testHttp, req)

			require.Equal(t, http.StatusCreated, testHttp.Code)
		})
}

func Test_CreateTeam_Should_Fail_IfTeamAlreadyExists(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByNameMock = func(name string) (*types.Team, error) { return &types.Team{}, nil }
	userStore := &mockUser{}
	userStore.GetUserByIdMock = func(id string) (*types.User, error) { return nil, fmt.Errorf("user does not exists") }
	handler := NewHandler(db, teamStore, userStore)
	payload := types.AddTeamPayload{
		Name: "Team A",
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/teams", bytes.NewBuffer(marshalled))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/teams", handler.handleAddTeam).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusConflict, testHttp.Code)

}

func Test_AddUserToTeam_Should_Fail_IfUserDontExists(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByIdMock = func(id string) (*types.Team, error) { return &types.Team{}, nil }
	userStore := &mockUser{}
	userStore.GetUserByIdMock = func(id string) (*types.User, error) { return nil, fmt.Errorf("user does not exists") }
	handler := NewHandler(db, teamStore, userStore)
	payload := types.UserToTeamPayload{
		UserId:   uuid.NewString(),
		TeamId:   uuid.NewString(),
		RoleType: types.Member,
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/teams/addUser", bytes.NewBuffer(marshalled))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/teams/addUser", handler.handleAddUserToTeam).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusBadRequest, testHttp.Code)
}

func Test_AddUserToTeam_Should_Fail_IfTeamDontExists(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByIdMock = func(id string) (*types.Team, error) { return nil, fmt.Errorf("team does not exists") }
	userStore := &mockUser{}
	userStore.GetUserByIdMock = func(id string) (*types.User, error) { return &types.User{}, nil }
	handler := NewHandler(db, teamStore, userStore)
	payload := types.UserToTeamPayload{
		UserId:   uuid.NewString(),
		TeamId:   uuid.NewString(),
		RoleType: types.Member,
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/teams/addUser", bytes.NewBuffer(marshalled))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/teams/addUser", handler.handleAddUserToTeam).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusBadRequest, testHttp.Code)
}

type mockUser struct {
	GetUserByEmailMock   func(email string) (*types.User, error)
	GetUserByIdMock      func(id string) (*types.User, error)
	CreateUserMock       func(types.User) error
	GetUsersFromTeamMock func(teamId string) ([]types.TeamUser, error)
}

func (m *mockUser) GetUserByEmail(email string) (*types.User, error) {
	return m.GetUserByEmailMock(email)
}
func (m *mockUser) GetUserById(id string) (*types.User, error) {
	return m.GetUserByIdMock(id)
}
func (m *mockUser) CreateUser(u types.User) error {
	return m.CreateUserMock(u)
}

func (m *mockUser) GetUsersFromTeam(teamId string) ([]types.TeamUser, error) {
	return m.GetUsersFromTeamMock(teamId)
}

type mockTeam struct {
	GetAllTeamsMock        func() ([]types.Team, error)
	CreateTeamMock         func(types.Team) error
	RenameTeamMock         func(name, teamId string) error
	GetTeamByIdMock        func(id string) (*types.Team, error)
	GetTeamByNameMock      func(name string) (*types.Team, error)
	AddUserToTeamMock      func(execable interface{}, userId, teamId string, role types.UserRole) error
	RemoveUserFromTeamMock func(userId, teamId string) error
}

func (m *mockTeam) GetAllTeams() ([]types.Team, error) {
	return m.GetAllTeamsMock()
}

func (m *mockTeam) GetTeamById(id string) (*types.Team, error) {
	return m.GetTeamByIdMock(id)
}

func (m *mockTeam) CreateTeam(t types.Team) error {
	return m.CreateTeamMock(t)
}

func (m *mockTeam) GetTeamByName(name string) (*types.Team, error) {
	return m.GetTeamByNameMock(name)
}

func (m *mockTeam) AddUserToTeam(execable interface{}, userId, teamId string, role types.UserRole) error {
	return m.AddUserToTeamMock(execable, userId, teamId, role)
}

func (m *mockTeam) RemoveUserFromTeam(userId, teamId string) error {
	return m.RemoveUserFromTeamMock(userId, teamId)
}

func (m *mockTeam) RenameTeam(name, teamId string) error {
	return nil
}
