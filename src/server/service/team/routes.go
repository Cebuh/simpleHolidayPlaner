package team

import (
	"fmt"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/service/auth"
	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.TeamStore
	userStore types.UserStore
}

func NewHandler(store types.TeamStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/teams", h.handleGetAllTeams).Methods(http.MethodGet)
	router.HandleFunc("/teams/{teamId}", h.handleGetAllTeams).Methods(http.MethodGet)
	// router.HandleFunc("/teams", h.handleAddTeam).Methods(http.MethodPost)
	router.HandleFunc("/teams", auth.Require(h.handleAddTeam, h.userStore)).Methods(http.MethodPost)
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

	if _, err := h.store.GetTeamByName(payload.Name); err == nil  {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("team %s already exists", payload.Name ))
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
