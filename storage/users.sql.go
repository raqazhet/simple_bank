package storage

import (
	"context"
	"log"

	"bank/model"
)

func (q *Queries) CreateUser(ctx context.Context, arg model.User) (model.User, error) {
	query := `INSERT INTO users (username,hashed_password,full_name,email)
	VALUES($1,$2,$3,$4)
	Returning username,hashed_password,full_name,email,password_changed_at,created_at
	`
	args := []any{arg.Username, arg.HashedPassword, arg.Fullname, arg.Email}
	user := model.User{}
	if err := q.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	); err != nil {
		log.Printf("CreateUser err: %v", err)
		return model.User{}, err
	}
	return user, nil
}

func (r *Queries) GetUser(ctx context.Context, username string) (model.User, error) {
	query := `SELECT * FROM users
	WHERE username = $1`
	user := model.User{}
	if err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Fullname,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
	); err != nil {
		log.Printf("GetUser err: %v", err)
		return model.User{}, err
	}
	return user, nil
}
