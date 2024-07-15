package types

import "time"

type InviteType int

const (
	Team_Invite InviteType = iota
	Group_Invite
)

type InviteStatus int

const (
	INVITE_OPEN InviteStatus = iota
	INVITE_ACCEPTED
	INVITE_DECLINED
)

type Invite struct {
	Id         string       `json:"id"`
	InviteType InviteType   `json:"inviteType"`
	FromUserId string       `json:"fromUserId"`
	ToUserId   string       `json:"toUserId"`
	TeamId     string       `json:"teamId"`
	Status     InviteStatus `json:"status"`
	CreatedAt  time.Time    `json:"createdAt"`
	ChangedAt  *time.Time   `json:"changedAt"`
}

type InviteInfo struct {
	Id           string       `json:"id"`
	FromUserName string       `json:"fromUserName"`
	ToUserName   string       `json:"toUserName"`
	TeamName     string       `json:"teamName"`
	Status       InviteStatus `json:"status"`
	CreatedAt    time.Time    `json:"createdAt"`
}

type CreateInvitePayload struct {
	InviteType InviteType   `json:"inviteType"`
	FromUserId string       `json:"fromUserId" validate:"required"`
	ToUserId   string       `json:"toUserId" validate:"required"`
	TeamId     string       `json:"teamId" validate:"required"`
	Status     InviteStatus `json:"status"`
}
