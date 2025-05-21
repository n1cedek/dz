package auth

import (
	"context"
	"database/sql"
	"dz/auth/internal/model"
	repo2 "dz/auth/internal/repo"
	"dz/auth/internal/repo/auth/converter"
	modelRepo "dz/auth/internal/repo/auth/model"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName = "auth"
	idCol     = "id"
	nameCol   = "name"
	emailCol  = "email"
	pasCol    = "password"
	upCol     = "updated_at"
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
func (r *repo) Delete(ctx context.Context, id int64) error {
	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idCol: id})

	q, a, err := builderDelete.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, q, a...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		log.Printf("No user found with id: %d", id)
		return sql.ErrNoRows
	}

	return nil
}
func (r *repo) Update(ctx context.Context, id int64, info *model.User) error {
	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameCol, info.Name).
		Set(emailCol, info.Email).
		Set(upCol, sq.Expr("NOW()")).
		Where(sq.Eq{idCol: id})

	q, a, err := builderUpdate.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.Exec(ctx, q, a...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()

	if rowsAffected == 0 {
		log.Printf("No user found with id: %d", id)
		return sql.ErrNoRows
	}
	return nil
}
