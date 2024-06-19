package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUser{}
	userStore.GetUserByEmailMock = func(email string) (*types.User, error) { return &types.User{}, nil }
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is not valid",
		func(t *testing.T) {
			payload := types.RegisterUserPayload{
				Name:     "Chris",
				Email:    "invalid",
				Password: "test124",
			}

			marshalled, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatal(err)
			}

			testHttp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
			router.ServeHTTP(testHttp, req)

			if testHttp.Code != http.StatusBadRequest {
				t.Errorf("expected status code %d, but got %d", http.StatusBadRequest, testHttp.Code)
			}
		})
	t.Run("should run if email is valid",
		func(t *testing.T) {
			payload := types.RegisterUserPayload{
				Name:     "Chris",
				Email:    "valid@email.com",
				Password: "test124",
			}

			marshalled, _ := json.Marshal(payload)
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatal(err)
			}

			testHttp := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
			router.ServeHTTP(testHttp, req)

			// Throw Status conflict, when user already exists
			if testHttp.Code != http.StatusConflict {
				t.Errorf("expected status code %d, but got %d", http.StatusCreated, testHttp.Code)
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
