package auth

import (
	"context"
	"dz/auth/model"
	repo2 "dz/auth/repo"
	"dz/auth/repo/auth/converter"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	modelRepo "dz/auth/repo/auth/model"
)

const (
	tableName = "auth"

	idCol    = "id"
	nameCol  = "name"
	emailCol = "email"
	pasCol   = "password"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) repo2.AuthRepo {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *modelRepo.User) (int64, error) {
	builderCreate := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameCol, emailCol, pasCol).
		Values(info.Name, info.Email, info.Password).
		Suffix("RETURNING id")

	q, a, err := builderCreate.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build insert query: %w", err)

	}
	var id int64
	err = r.db.QueryRow(ctx, q, a...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.PublicInfo, error) {
	builderGet := sq.Select(idCol, nameCol, emailCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idCol: id}).
		Limit(1)

	q, a, err := builderGet.ToSql()
	if err != nil {
		return nil, err
	}
	var user modelRepo.User
	err = r.db.QueryRow(ctx, q, a...).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return converter.ToPublicInfo(&user), nil
}
