package storage

import (
	"context"
	"log"

	"bank/model"

	"github.com/google/uuid"
)

func (q *Queries) SaveNewRefreshToken(ctx context.Context, arg model.CreateSessionParams) (model.CreateSessionParams, error) {
	query := `INSERT INTO sessions (id,username,refresh_token,user_agent,client_ip,is_blocked,expires_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
	`
	args := []any{arg.ID, arg.Username, arg.RefreshToken, arg.UserAgent, arg.ClientIp, arg.ExpiresAt}
	var i model.CreateSessionParams
	if err := q.db.QueryRowContext(ctx, query, args).Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt); err != nil {
		log.Println("saveNewRefreshToken err")
		return model.CreateSessionParams{}, err
	}
	return i, nil
}

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (model.CreateSessionParams, error) {
	query := `Select * from sessions
		Where id =$1 Limit 1`
	var i model.CreateSessionParams
	if err := q.db.QueryRowContext(ctx, query, id).Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt); err != nil {
		return model.CreateSessionParams{}, err
	}
	return i, nil
}
