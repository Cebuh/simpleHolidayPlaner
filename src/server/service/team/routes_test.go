package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/gorilla/mux"
)

func TestTeamServiceHandlers(t *testing.T) {
	teamStore := &mockTeam{}
	teamStore.GetTeamByNameMock = func(name string) (*types.Team, error) { return nil, fmt.Errorf("Not found") }
	teamStore.CreateTeamMock = func(t types.Team) error { return nil }

	userStore := &mockUser{}
	handler := NewHandler(teamStore, userStore)

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

			if testHttp.Code != http.StatusCreated {
				t.Errorf("expected status code %d, but got %d", http.StatusCreated, testHttp.Code)
			}
		})
}

func TestTeamServiceHandlers2(t *testing.T) {
	teamStore := &mockTeam{}
	teamStore.GetTeamByNameMock = func(name string) (*types.Team, error) { return &types.Team{}, nil }
	userStore := &mockUser{}
	handler := NewHandler(teamStore, userStore)

	t.Run("should fail if team exists",
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

			if testHttp.Code != http.StatusConflict {
				t.Errorf("expected status code %d, but got %d", http.StatusConflict, testHttp.Code)
			}
		})
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
	AddUserToTeamMock      func(userId, teamId, roleType string) error
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

func (m *mockTeam) AddUserToTeam(userId, teamId, roleType string) error {
	return m.AddUserToTeamMock(userId, teamId, roleType)
}

func (m *mockTeam) RemoveUserFromTeam(userId, teamId string) error {
	return m.RemoveUserFromTeamMock(userId, teamId)
}

func (m *mockTeam) RenameTeam(name, teamId string) error {
	return nil
}
