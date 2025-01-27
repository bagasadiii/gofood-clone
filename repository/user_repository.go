package repository

import (
	"context"
	"fmt"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepoImpl interface {
	RegisterUserRepo(ctx context.Context, new *model.User) error
	GetUserRepo(ctx context.Context, id uuid.UUID) (*model.UserResp, error)
	GetIDRepo(ctx context.Context, username string) (uuid.UUID, error)
	LoginRepo(ctx context.Context, username string) (*model.User, error)
}
type UserRepo struct {
	db  *pgxpool.Pool
	zap *zap.Logger
}

func NewUserRepo(db *pgxpool.Pool, zap *zap.Logger) *UserRepo {
	return &UserRepo{
		db:  db,
		zap: zap,
	}
}

func (ur *UserRepo) RegisterUserRepo(ctx context.Context, new *model.User) error {
	var exists bool
	err := ur.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2)`, new.Username, new.Email).
		Scan(&exists)
	if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%v", utils.ErrDatabase)
	}
	if exists {
		return fmt.Errorf("%v", utils.ErrUniqueConstraint)
	}
	_, err = ur.db.Exec(ctx, `INSERT INTO users (user_id, username, email, password, role, created_at, phone, balance, name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		new.UserID, new.Username, new.Email, new.Password, new.Role, new.CreatedAt, new.Phone, new.Balance, new.Name)
	if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("failed to create user: %w", utils.ErrDatabase)
	}
	return nil
}

func (ur *UserRepo) GetUserRepo(ctx context.Context, id uuid.UUID) (*model.UserResp, error) {
	resp := &model.UserResp{}
	row := ur.db.QueryRow(ctx, `SELECT username, email, role, created_at, phone, name FROM users WHERE user_id = $1`, id)

	err := row.Scan(&resp.Username, &resp.Email, &resp.Role, &resp.CreatedAt, &resp.Phone, &resp.Name)
	if err == pgx.ErrNoRows {
		ur.zap.Error(utils.ErrNotFound.Error(), zap.String("user_id", id.String()))
		return nil, fmt.Errorf("no row found: %v", utils.ErrNotFound)
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", utils.ErrDatabase)
	}

	ur.zap.Info("User fetched", zap.String("Username", resp.Username))
	return resp, nil
}

func (ur *UserRepo) GetIDRepo(ctx context.Context, username string) (uuid.UUID, error) {
	var id uuid.UUID
	err := ur.db.QueryRow(ctx, `SELECT user_id FROM users WHERE username = $1`, username).Scan(&id)
	if err == pgx.ErrNoRows {
		ur.zap.Warn(utils.ErrNotFound.Error(), zap.String("Username", username))
		return uuid.Nil, err
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return uuid.Nil, fmt.Errorf("failed to fetch id: %w", utils.ErrDatabase)
	}
	return id, nil
}

func (ur *UserRepo) LoginRepo(ctx context.Context, username string) (*model.User, error) {
	var res model.User
	row := ur.db.QueryRow(ctx, `SELECT user_id, username, password, role FROM users WHERE username = $1`, username)
	err := row.Scan(
		&res.UserID,
		&res.Username,
		&res.Password,
		&res.Role,
	)
	if err == pgx.ErrNoRows {
		ur.zap.Warn(utils.ErrNotFound.Error(), zap.String("Username", username))
		return nil, err
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", utils.ErrDatabase)
	}
	return &res, nil
}

