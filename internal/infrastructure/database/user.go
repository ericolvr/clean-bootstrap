package database

import (
	"context"
	"database/sql"

	"github.com/ericolvr/sec-backend/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			name, mobile, user_type, password, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Mobile,
		user.UserType,
		user.Password,
		user.Status,
	).Scan(&user.ID)

	return err
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, name, mobile, user_type, password, status
		FROM users
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Mobile,
			&user.UserType,
			&user.Password,
			&user.Status,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, name, mobile, user_type, password, status
		FROM users
		WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.UserType,
		&user.Password,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	query := `
		SELECT id, name, mobile, user_type, password, status
		FROM users
		WHERE mobile = $1`

	row := r.db.QueryRowContext(ctx, query, mobile)
	user := &domain.User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Mobile,
		&user.UserType,
		&user.Password,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $2, mobile = $3, user_type = $4, status = $5, password = $6
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Mobile,
		user.UserType,
		user.Status,
		user.Password,
	)

	return err
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
