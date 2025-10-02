package repository

import (
	"context"
	"discord/config"
	"discord/gen/repo"
)

type AuthRepository struct {
	readdb  *repo.Queries
	writedb *repo.Queries
}

func NewAuthRepository() (*AuthRepository, error) {
	queries, err := config.RepoQuieries()
	if err != nil {
		return nil, err
	}
	return &AuthRepository{
		readdb:  repo.New(queries.ReadDb),
		writedb: repo.New(queries.WriteDb),
	}, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, value repo.CreateUserParams) (*repo.User, error) {
	user, err := r.writedb.CreateUser(ctx, value)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, username string) (*repo.User, error) {
	user, err := r.readdb.GetUserByEmail(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*repo.User, error) {
	user, err := r.readdb.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
