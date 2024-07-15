package vacation

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	db            *sql.DB
	userStore     types.UserStore
	teamStore     types.TeamStore
	vacationStore types.VacationStore
}

func NewHandler(db *sql.DB, userStore types.UserStore, teamStore types.TeamStore, vacationStore types.VacationStore) *Handler {
	return &Handler{db: db, userStore: userStore, teamStore: teamStore, vacationStore: vacationStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/vacations/request", h.CreateVacationRequest).Methods(http.MethodPost)
}

func (h *Handler) CreateVacationRequest(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateVacationRequestPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	if _, err := h.teamStore.GetTeamById(payload.TeamId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("team with id %s does not exists", payload.TeamId))
		return
	}

	if _, err := h.userStore.GetUserById(payload.ToUserId); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("userToId with id %s does not exists", payload.ToUserId))
		return
	}

	if _, err := h.userStore.GetUserById(payload.RequestedFrom); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("requested from user with id %s does not exists", payload.RequestedFrom))
		return
	}

	ctx := r.Context()
	utils.WithTransaction(ctx, h.db, w, func(tx *sql.Tx) error {

		request := types.VacationRequest{
			Id:            uuid.NewString(),
			RequestedFrom: payload.RequestedFrom,
			ToUserId:      payload.ToUserId,
			TeamId:        payload.TeamId,
			Info:          payload.Info,
			Status:        types.REQUEST_OPEN,
			FromDate:      payload.FromDate,
			ToDate:        payload.ToDate,
		}

		if err := h.vacationStore.CreateVacationRequest(tx, request); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return err
		}

		utils.WriteJson(w, http.StatusOK, nil)
		return nil
	})
}
