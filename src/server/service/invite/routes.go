package invite

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	db        *sql.DB
	store     types.InviteStore
	userStore types.UserStore
	teamStore types.TeamStore
}

func NewHandler(db *sql.DB, store types.InviteStore, userStore types.UserStore, teamStore types.TeamStore) *Handler {
	return &Handler{db: db, store: store, teamStore: teamStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/invites/from/{userId}", h.GetInvitesFromUser).Methods(http.MethodGet)
	router.HandleFunc("/invites/to/{userId}", h.GetInvitesToUser).Methods(http.MethodGet)
	router.HandleFunc("/invites", h.CreateInvite).Methods(http.MethodPost)
	router.HandleFunc("/invites/{id}/approve", h.ApproveInvite).Methods(http.MethodPost)
	router.HandleFunc("/invites/{id}/decline", h.DeclineInvite).Methods(http.MethodPost)
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

	ctx := r.Context()
	utils.WithTransaction(ctx, h.db, w, func(tx *sql.Tx) error {

		if err := h.store.UpdateInviteStatus(tx, id, types.DECLINED); err != nil {
			return err
		}

		// TODO - delete invite when declined? its useful to hold it to show the user the state
		// if err := h.store.DeleteInvite(tx, id); err != nil {
		// 	return err
		// }

		utils.WriteJson(w, http.StatusOK, nil)
		return nil
	})

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

	inv, err := h.store.GetInvite(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	users, err := h.userStore.GetUsersFromTeam(inv.TeamId)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("while loading users from team a error occured %e", err))
		return
	}

	userExistsInTeam := false
	for _, v := range users {
		if v.Id == inv.ToUserId {
			userExistsInTeam = true
			break
		}
	}

	if userExistsInTeam {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user is already a part of the team"))
		return
	}

	ctx := r.Context()
	utils.WithTransaction(ctx, h.db, w, func(tx *sql.Tx) error {
		if err := h.store.UpdateInviteStatus(tx, id, types.ACCEPTED); err != nil {
			return err
		}

		if err := h.teamStore.AddUserToTeam(tx, inv.ToUserId, inv.TeamId, types.Member); err != nil {
			return err
		}

		utils.WriteJson(w, http.StatusOK, nil)
		return nil
	})

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
