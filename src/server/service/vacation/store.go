package vacation

import (
	"database/sql"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateVacationRequest(execable interface{}, request types.VacationRequest) error {
	_, err := utils.Exec(execable, "INSERT INTO vacation_requests (id, requestedFrom, toUserId, teamId, fromDate, toDate, info, requestStatus) VALUES (?,?,?,?,?,?,?,?)",
		request.Id, request.RequestedFrom, request.ToUserId, request.TeamId, request.FromDate, request.ToDate, request.Info, request.Status)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetVacationRequestsForUser(toUserId string) ([]types.VacationRequest, error) {
	return nil, nil
}

func (s *Store) GetVacationRequestsFromUserId(requestedFromId string) ([]types.VacationRequest, error) {
	return nil, nil
}

func (s *Store) UpdateVacationStatus(execable interface{}, requestId string, approverId string, status types.ApprovalStatus) error {
	return nil
}

func (s *Store) GetApprovalsForRequest(requestId string) ([]types.VacationApproval, error) {
	return nil, nil
}
