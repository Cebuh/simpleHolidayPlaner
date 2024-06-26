package invite

import (
	"database/sql"
	"fmt"

	"github.com/cebuh/simpleHolidayPlaner/types"
	"github.com/cebuh/simpleHolidayPlaner/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateInvite(i types.Invite) error {
	_, err := s.db.Exec("INSERT INTO invites (id, fromUserId, toUserId, teamId, InviteType, status) VALUES (?, ?, ?, ?, ?, ?)",
		i.Id, i.FromUserId, i.ToUserId, i.TeamId, i.InviteType, i.Status)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetInvite(id string) (*types.Invite, error) {
	rows, err := s.db.Query(`SELECT * from invites i where i.Id = ?`, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	invite := new(types.Invite)
	for rows.Next() {
		invite, err = readInviteData(rows)
		if err != nil {
			return nil, err
		}
	}

	if !utils.IsValidUUID(invite.Id) {
		return nil, fmt.Errorf("invite not found")
	}

	return invite, nil
}

func (s *Store) GetInviteInfosFrom(fromUserId string) ([]types.InviteInfo, error) {
	rows, err := s.db.Query(`SELECT i.Id, ufrom.name as 'FromUserName', uto.name as 'ToUserName', t.name as 'TeamName', i.status, i.createdAt From invites i  
							inner join users ufrom  on ufrom.id = i.fromUserId 
							inner join users uto  on uto.id = i.toUserId  
							inner join teams t  on t.Id = i.teamId 
							where i.fromUserId = ?`, fromUserId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inviteInfos := make([]types.InviteInfo, 0)
	for rows.Next() {
		info, err := readInviteInfoData(rows)
		if err != nil {
			return nil, err
		}
		inviteInfos = append(inviteInfos, *info)
	}

	return inviteInfos, nil
}

func (s *Store) UpdateInviteStatus(execable interface{}, id string, status types.InviteStatus) error {

	_, err := utils.Exec(execable, `UPDATE invites set status = ?, changedAt = UTC_TIMESTAMP where id = ?`, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetInviteInfosTo(toUserId string) ([]types.InviteInfo, error) {
	rows, err := s.db.Query(`SELECT i.Id, ufrom.name as 'FromUserName', uto.name as 'ToUserName', t.name as 'TeamName', i.status, i.createdAt From invites i
	inner join users ufrom  on ufrom.id = i.fromUserId 
	inner join users uto  on uto.id = i.toUserId  
	inner join teams t  on t.Id = i.teamId 
	where i.toUserId = ?`, toUserId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	inviteInfos := make([]types.InviteInfo, 0)
	for rows.Next() {
		info, err := readInviteInfoData(rows)
		if err != nil {
			return nil, err
		}
		inviteInfos = append(inviteInfos, *info)
	}

	return inviteInfos, nil
}

func readInviteInfoData(rows *sql.Rows) (*types.InviteInfo, error) {
	inv := new(types.InviteInfo)
	err := rows.Scan(
		&inv.Id,
		&inv.FromUserName,
		&inv.ToUserName,
		&inv.TeamName,
		&inv.Status,
		&inv.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func readInviteData(rows *sql.Rows) (*types.Invite, error) {
	inv := new(types.Invite)
	err := rows.Scan(
		&inv.Id,
		&inv.FromUserId,
		&inv.ToUserId,
		&inv.TeamId,
		&inv.Status,
		&inv.InviteType,
		&inv.CreatedAt,
		&inv.ChangedAt,
	)
	if err != nil {
		return nil, err
	}
	return inv, nil
}
