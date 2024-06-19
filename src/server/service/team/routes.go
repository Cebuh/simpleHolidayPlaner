package team

import (
	"fmt"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.TeamStore
	userStore types.UserStore
}

func NewHandler(store types.TeamStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/teams", h.handleGetAllTeams).Methods(http.MethodGet)
	router.HandleFunc("/teams/{teamId}", h.handleGetTeam).Methods(http.MethodGet)
	router.HandleFunc("/teams", h.handleAddTeam).Methods(http.MethodPost)
	router.HandleFunc("/teams/addUser", h.handleAddUserToTeam).Methods(http.MethodPost)
	router.HandleFunc("/teams/removeUser", h.handleRemoveUserFromTeam).Methods(http.MethodPost)
	router.HandleFunc("/teams/{teamId}/getUsers", h.handleGetUsersFromTeam).Methods(http.MethodGet)
	router.HandleFunc("/teams/{teamId}", h.handleRenameTeam).Methods(http.MethodPatch)

	// router.HandleFunc("/teams", auth.Require(h.handleAddTeam, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleRemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	var userTeamPayload types.UserToTeamPayload
	if err := utils.ParseJson(r, &userTeamPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(userTeamPayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if err := h.store.RemoveUserFromTeam(userTeamPayload.UserId, userTeamPayload.TeamId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error while remove user from team. error: %e", err))
		return
	}

	utils.WriteJson(w, http.StatusOK, nil)
}

func (h *Handler) handleRenameTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["teamId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing team id"))
		return
	}

	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	var renamePayload types.RenameTeamPayload
	if err := utils.ParseJson(r, &renamePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(renamePayload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if err := h.store.RenameTeam(renamePayload.Name, id); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, nil)
}

func (h *Handler) handleAddUserToTeam(w http.ResponseWriter, r *http.Request) {
	var payload types.UserToTeamPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if _, err := h.store.GetTeamById(payload.TeamId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("team with id %s doe not exists", payload.TeamId))
		return
	}

	if _, err := h.userStore.GetUserById(payload.UserId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with id %s doe not exists", payload.UserId))
		return
	}

	users, err := h.userStore.GetUsersFromTeam(payload.TeamId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("while loading users from team a error occured %e", err))
		return
	}

	userExistsInTeam := false
	for _, v := range users {
		if v.Id == payload.UserId {
			userExistsInTeam = true
			break
		}
	}

	if userExistsInTeam {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user is already a part of the team"))
		return
	}

	if err := h.store.AddUserToTeam(payload.UserId, payload.TeamId, payload.RoleType); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, nil)
}

func (h *Handler) handleGetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["teamId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing team id"))
		return
	}

	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	team, err := h.store.GetTeamById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, team)
}

func (h *Handler) handleGetUsersFromTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["teamId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing team id"))
		return
	}

	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	users, err := h.userStore.GetUsersFromTeam(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, users)
}

func (h *Handler) handleGetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.store.GetAllTeams()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, teams)
}

func (h *Handler) handleAddTeam(w http.ResponseWriter, r *http.Request) {
	var payload types.AddTeamPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if _, err := h.store.GetTeamByName(payload.Name); err == nil {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("team %s already exists", payload.Name))
		return
	}

	team := types.Team{
		Id:   uuid.NewString(),
		Name: payload.Name,
	}

	err := h.store.CreateTeam(team)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, team)
}
