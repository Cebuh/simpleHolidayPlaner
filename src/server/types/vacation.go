package types

import "time"

type RequestStatus int

const (
	REQUEST_OPEN RequestStatus = iota
	REQUEST_SUBSTITUTED_MEMBER
	REQUEST_SUBSTITUTED_TEAMLEAD
	REQUEST_APPROVED
	REQUEST_DECLINED
)

// the internal data to handle logic
type VacationRequest struct {
	Id            string        `json:"id"`
	RequestedFrom string        `json:"requestedFrom"`
	ToUserId      string        `json:"toUserId"`
	TeamId        string        `json:"teamId"`
	Info          string        `json:"info"`
	Status        RequestStatus `json:"status"`
	FromDate      time.Time     `json:"fromDate"`
	ToDate        time.Time     `json:"toDate"`
	ChangedAt     time.Time     `json:"changedAt"`
	CreatedAt     time.Time     `json:"createdAt"`
}

// The vacation request for display data
type VacationRequestInfo struct {
	Id           string        `json:"id"`
	FromUserName string        `json:"fromUserName"`
	ToUserName   string        `json:"toUsername"`
	TeamName     string        `json:"teamName"`
	Info         string        `json:"info"`
	Status       RequestStatus `json:"status"`
	FromDate     time.Time     `json:"fromDate"`
	ToDate       time.Time     `json:"toDate"`
	ChangedAt    time.Time     `json:"changedAt"`
}

type CreateVacationRequestPayload struct {
	RequestedFrom string    `json:"requestedFrom" validate:"required,uuid4"`
	ToUserId      string    `json:"toUserId" validate:"required,uuid4"`
	TeamId        string    `json:"teamId" validate:"required,uuid4"`
	Info          string    `json:"info" validate:"required"`
	FromDate      time.Time `json:"fromDate" validate:"required"`
	ToDate        time.Time `json:"toDate" validate:"required"`
}

type ApprovalStatus int

const (
	APPROVAL_OPEN ApprovalStatus = iota
	APPROVAL_APPROVED
	APPROVAL_DECLINED
)

type VacationApproval struct {
	RequestId  string         `json:"requestId"`
	ApproverId string         `json:"approverId"`
	Status     ApprovalStatus `json:"status"`
	ChangedAt  time.Time      `json:"changedAt"`
}

type VacationApprovalPayload struct {
	RequestId  string         `json:"requestId"`
	ApproverId string         `json:"approverId"`
	Status     ApprovalStatus `json:"status"`
}
