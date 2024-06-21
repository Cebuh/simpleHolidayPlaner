package invite

import (
	"fmt"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.InviteStore
	userStore types.UserStore
	teamStore types.TeamStore
}

func NewHandler(store types.InviteStore, userStore types.UserStore, teamStore types.TeamStore) *Handler {
	return &Handler{store: store, teamStore: teamStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/invites/from/{userId}", h.GetInvitesFromUser).Methods(http.MethodGet)
	router.HandleFunc("/invites/to/{userId}", h.GetInvitesToUser).Methods(http.MethodGet)
	router.HandleFunc("/invites", h.CreateInvite).Methods(http.MethodPost)
	router.HandleFunc("/invites/{id}/approve", h.ApproveInvite).Methods(http.MethodPost)
	router.HandleFunc("/invites/{id}/decline", h.DeclineInvite).Methods(http.MethodPost)
}

func (h *Handler) ApproveInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing invite id"))
		return
	}

	// TODO: validator can check this too
	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	_, err := h.store.GetInvite(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.UpdateInviteStatus(id, types.ACCEPTED); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, nil)

}
func (h *Handler) DeclineInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing invite id"))
		return
	}

	// TODO: validator can check this too
	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	_, err := h.store.GetInvite(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.UpdateInviteStatus(id, types.DECLINED); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, nil)
}

func (h *Handler) GetInvitesFromUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user id"))
		return
	}

	// TODO: validator can check this too
	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	invites, err := h.store.GetInviteInfosFrom(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, invites)
}
func (h *Handler) GetInvitesToUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user id"))
		return
	}

	// TODO: validator can check this too
	if !utils.IsValidUUID(id) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("id is not valid"))
		return
	}

	invites, err := h.store.GetInviteInfosTo(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, invites)
}
func (h *Handler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateInvitePayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if !utils.ValidatePayload(w, payload) {
		return
	}

	if _, err := h.teamStore.GetTeamById(payload.TeamId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("team does not exists"))
		return
	}

	if _, err := h.userStore.GetUserById(payload.FromUserId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("from user does not exists"))
		return
	}

	if _, err := h.userStore.GetUserById(payload.ToUserId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("to user does not exists"))
		return
	}

	invite := types.Invite{
		Id:         uuid.NewString(),
		InviteType: types.Group_Invite,
		FromUserId: payload.FromUserId,
		ToUserId:   payload.ToUserId,
		TeamId:     payload.TeamId,
		Status:     types.OPEN,
	}

	err := h.store.CreateInvite(invite)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, invite)

}
