package invite

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

func Test_CreateInvite_Should_Fail_IfUserDontExists(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByIdMock = func(id string) (*types.Team, error) { return &types.Team{}, nil }
	userStore := &mockUser{}
	userStore.GetUserByIdMock = func(id string) (*types.User, error) { return nil, fmt.Errorf("user does not exists") }
	inviteStore := &mockInvite{}
	inviteStore.CreateInviteMock = func(inv types.Invite) error { return nil }
	handler := NewHandler(db, inviteStore, userStore, teamStore)
	payload := types.CreateInvitePayload{
		FromUserId: uuid.NewString(),
		ToUserId:   uuid.NewString(),
		TeamId:     uuid.NewString(),
		InviteType: types.Team_Invite,
		Status:     types.OPEN,
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/invites", bytes.NewBuffer(marshalled))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites", handler.CreateInvite).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusBadRequest, testHttp.Code)
}

func Test_CreateInvite_Should_Pass_IfInviteIsCreated(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.GetTeamByIdMock = func(id string) (*types.Team, error) { return &types.Team{}, nil }
	userStore := &mockUser{}
	userStore.GetUserByIdMock = func(id string) (*types.User, error) { return &types.User{}, nil }
	inviteStore := &mockInvite{}
	inviteStore.CreateInviteMock = func(inv types.Invite) error { return nil }
	handler := NewHandler(db, inviteStore, userStore, teamStore)
	payload := types.CreateInvitePayload{
		FromUserId: uuid.NewString(),
		ToUserId:   uuid.NewString(),
		TeamId:     uuid.NewString(),
		InviteType: types.Team_Invite,
		Status:     types.OPEN,
	}

	marshalled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/invites", bytes.NewBuffer(marshalled))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites", handler.CreateInvite).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusCreated, testHttp.Code)
}

func Test_GetInvites_Should_Pass_ForFrom(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	userStore := &mockUser{}
	inviteStore := &mockInvite{}
	inviteStore.GetInviteInfosFromMock = func(from string) ([]types.InviteInfo, error) { return make([]types.InviteInfo, 0), nil }
	handler := NewHandler(db, inviteStore, userStore, teamStore)

	testGuid := uuid.NewString()
	req, err := http.NewRequest(http.MethodGet, "/invites/from/"+testGuid, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites/from/{userId}", handler.GetInvitesFromUser).Methods(http.MethodGet)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusOK, testHttp.Code)
}

func Test_GetInvites_Should_Pass_ForTo(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	userStore := &mockUser{}
	inviteStore := &mockInvite{}
	inviteStore.GetInviteInfosToMock = func(to string) ([]types.InviteInfo, error) { return make([]types.InviteInfo, 0), nil }
	handler := NewHandler(db, inviteStore, userStore, teamStore)

	testGuid := uuid.NewString()
	req, err := http.NewRequest(http.MethodGet, "/invites/to/"+testGuid, nil)
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites/to/{userId}", handler.GetInvitesToUser).Methods(http.MethodGet)
	router.ServeHTTP(testHttp, req)

	require.Equal(t, http.StatusOK, testHttp.Code)
}

func Test_ApproveInvite_Should_Fail_IfUserIsAlreadyInTeam(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	userStore := &mockUser{}
	testGuid := uuid.NewString()
	testUser := types.TeamUser{
		Id: testGuid,
	}

	userSlice := make([]types.TeamUser, 0)
	userSlice = append(userSlice, testUser)
	userStore.GetUsersFromTeamMock = func(id string) ([]types.TeamUser, error) { return userSlice, nil }
	inviteStore := &mockInvite{}
	inviteStore.GetInviteMock = func(id string) (*types.Invite, error) {
		return &types.Invite{ToUserId: testGuid}, nil
	}
	handler := NewHandler(db, inviteStore, userStore, teamStore)

	req, err := http.NewRequest(http.MethodPost, "/invites/"+testGuid+"/approve", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites/{id}/approve", handler.ApproveInvite).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)
	require.Equal(t, http.StatusBadRequest, testHttp.Code)
}

func Test_ApproveInvite_Should_Pass_IfUserIsNotInTeam(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.ExpectBegin()
	mock.ExpectCommit()
	require.NoError(t, err)
	defer db.Close()
	teamStore := &mockTeam{}
	teamStore.AddUserToTeamMock = func(execable interface{}, userId, teamId string, role types.UserRole) error { return nil }
	userStore := &mockUser{}
	testGuid := uuid.NewString()
	testUser := types.TeamUser{
		Id: testGuid,
	}

	userSlice := make([]types.TeamUser, 0)
	userSlice = append(userSlice, testUser)
	userStore.GetUsersFromTeamMock = func(id string) ([]types.TeamUser, error) { return userSlice, nil }
	inviteStore := &mockInvite{}
	inviteStore.GetInviteMock = func(id string) (*types.Invite, error) {
		return &types.Invite{ToUserId: uuid.NewString()}, nil
	}
	inviteStore.UpdateInviteStatusMock = func(execable interface{}, id string, status types.InviteStatus) error { return nil }

	handler := NewHandler(db, inviteStore, userStore, teamStore)

	req, err := http.NewRequest(http.MethodPost, "/invites/"+testGuid+"/approve", nil)
	if err != nil {
		t.Fatal(err)
	}

	testHttp := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/invites/{id}/approve", handler.ApproveInvite).Methods(http.MethodPost)
	router.ServeHTTP(testHttp, req)
	require.Equal(t, http.StatusOK, testHttp.Code)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

type mockInvite struct {
	CreateInviteMock       func(types.Invite) error
	GetInviteInfosFromMock func(from string) ([]types.InviteInfo, error)
	GetInviteInfosToMock   func(to string) ([]types.InviteInfo, error)
	GetInviteMock          func(id string) (*types.Invite, error)
	DeleteInviteMock       func(execable interface{}, id string) error
	UpdateInviteStatusMock func(execable interface{}, id string, status types.InviteStatus) error
}

func (m *mockInvite) DeleteInvite(execable interface{}, id string) error {
	return m.DeleteInviteMock(execable, id)
}

func (m *mockInvite) GetInvite(id string) (*types.Invite, error) {
	return m.GetInviteMock(id)
}
func (m *mockInvite) UpdateInviteStatus(execable interface{}, id string, status types.InviteStatus) error {
	return m.UpdateInviteStatusMock(execable, id, status)
}

func (m *mockInvite) CreateInvite(inv types.Invite) error {
	return m.CreateInviteMock(inv)
}

func (m *mockInvite) GetInviteInfosFrom(from string) ([]types.InviteInfo, error) {
	return m.GetInviteInfosFromMock(from)
}

func (m *mockInvite) GetInviteInfosTo(to string) ([]types.InviteInfo, error) {
	return m.GetInviteInfosToMock(to)
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
