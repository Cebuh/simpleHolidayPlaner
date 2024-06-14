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
	teamStore := &mockTeamStore{}
	userStore := &mockUserStore{}
	handler := NewHandler(teamStore, userStore)

	t.Run("should run if team is created",
		func(t *testing.T) {
			payload := types.AddTeamPayload{
				Name:     "Team A",
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



type mockTeamStore struct{}
func (m *mockTeamStore) GetAllTeams() ([]types.Team, error) {
	return nil, nil
}

func (m *mockTeamStore) GetTeamById(id string) (*types.Team, error) {
	return nil, nil
}

func (m *mockTeamStore) CreateTeam(types.Team) error {
	return nil
}

func (m *mockTeamStore) GetTeamByName(name string) (*types.Team, error) {
	return nil, fmt.Errorf("Team not exists")
}

// Interface returns team
func TestTeamServiceHandlers2(t *testing.T) {
	teamStore := &mockTeamStore2{}
	userStore := &mockUserStore{}
	handler := NewHandler(teamStore, userStore)

	t.Run("should fail if team exists",
		func(t *testing.T) {
			payload := types.AddTeamPayload{
				Name:     "Team A",
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

type mockTeamStore2 struct{}
func (m *mockTeamStore2) GetAllTeams() ([]types.Team, error) {
	return nil, nil
}

func (m *mockTeamStore2) GetTeamById(id string) (*types.Team, error) {
	return nil, nil
}

func (m *mockTeamStore2) CreateTeam(types.Team) error {
	return nil
}

func (m *mockTeamStore2) GetTeamByName(name string) (*types.Team, error) {
	return &types.Team{}, nil
}


type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) GetUserById(id string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
