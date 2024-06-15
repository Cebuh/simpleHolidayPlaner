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
	userStore := &mockUserStore{}
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

func (m *mockUserStore) GetUsersFromTeam(teamId string) ([]types.TeamUser, error) {
	return nil, nil
}
