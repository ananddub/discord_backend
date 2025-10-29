package repository

import (
	"context"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		queries: repo.New(db),
		db:      db,
	}
}

// CreateUser creates a new user in database
func (r *AuthRepository) CreateUser(ctx context.Context, username, email, password string) error {
	_, err := r.queries.CreateUser(ctx, repo.CreateUserParams{
		Username: username,
		Email:    email,
		Password: password,
	})
	return err
}

// GetUserByID retrieves user by ID
func (r *AuthRepository) GetUserByID(ctx context.Context, id int32) (*repo.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves user by email
func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*repo.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves user by username
func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*repo.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus updates user online status
func (r *AuthRepository) UpdateUserStatus(ctx context.Context, userID int32, status string) error {
	return r.queries.UpdateUserStatus(ctx, repo.UpdateUserStatusParams{
		ID:     userID,
		Status: status,
	})
}

// UpdatePassword updates user password
func (r *AuthRepository) UpdatePassword(ctx context.Context, userID int32, password string) error {
	return r.queries.UpdateUserPassword(ctx, repo.UpdateUserPasswordParams{
		ID:       userID,
		Password: password,
	})
}

// Enable2FA enables or disables 2FA for user
func (r *AuthRepository) Enable2FA(ctx context.Context, userID int32, enabled bool) error {
	return r.queries.Enable2FA(ctx, repo.Enable2FAParams{
		ID:           userID,
		Is2faEnabled: pgtype.Bool{Bool: enabled, Valid: true},
	})
}

// DeleteUser deletes a user account
func (r *AuthRepository) DeleteUser(ctx context.Context, userID int32) error {
	return r.queries.SoftDeleteUser(ctx, userID)
}
