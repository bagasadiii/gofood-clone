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
	LoginRepo(ctx context.Context, username string) (*model.User, error)
	GetUserRepo(ctx context.Context, username string) (*model.UserResp, error)
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
	err := ur.db.QueryRow(ctx, `
    SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 OR email = $2)
    `, new.Username, new.Email).
		Scan(&exists)
	if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.String("failed to register", new.Username), zap.Error(err))
		return fmt.Errorf("failed to register user: %w", utils.ErrDatabase)
	}
	if exists {
		ur.zap.Warn(utils.ErrUniqueConstraint.Error(), zap.String("username", new.Username))
		return fmt.Errorf("username or email already exist: %w", utils.ErrUniqueConstraint)
	}
	_, err = ur.db.Exec(ctx, `
    INSERT INTO users (user_id, username, email, password, role, created_at, phone, balance, name)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, new.UserID, new.Username, new.Email, new.Password, new.Role, new.CreatedAt, new.Phone, new.Balance, new.Name)
	if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.String("failed to register", new.Username), zap.Error(err))
		return fmt.Errorf("failed to create user: %w", utils.ErrDatabase)
	}
	return nil
}

func (ur *UserRepo) GetUserRepo(ctx context.Context, username string) (*model.UserResp, error) {
	var userID uuid.UUID
	var res model.UserResp
	err := ur.db.QueryRow(ctx, `
    SELECT user_id FROM users WHERE username = $1
    `, username).Scan(&userID)
	if err == pgx.ErrNoRows {
		ur.zap.Warn(utils.ErrNotFound.Error(), zap.String("username", username))
		return nil, fmt.Errorf("no row found: %w", utils.ErrNotFound)
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", utils.ErrDatabase)
	}

	err = ur.db.QueryRow(ctx, `
    SELECT username, email, role, created_at, phone, name FROM users WHERE user_id = $1
    `, &userID).Scan(&res.Username, &res.Email, &res.Role, &res.CreatedAt, &res.Phone)
	if err == pgx.ErrNoRows {
		ur.zap.Warn(utils.ErrNotFound.Error(), zap.String("user_id", userID.String()))
		return nil, fmt.Errorf("no row found: %w", utils.ErrNotFound)
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", utils.ErrDatabase)
	}

	return &res, nil
}

func (ur *UserRepo) LoginRepo(ctx context.Context, username string) (*model.User, error) {
	var res model.User
	err := ur.db.QueryRow(ctx, `
    SELECT user_id, username, password, role FROM users WHERE username = $1
    `, username).Scan(&res.UserID, &res.Username, &res.Password, &res.Role)
	if err == pgx.ErrNoRows {
		ur.zap.Warn(utils.ErrNotFound.Error(), zap.String("Username", username))
		return nil, fmt.Errorf("no username found: %w", utils.ErrNotFound)
	} else if err != nil {
		ur.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %w", utils.ErrDatabase)
	}
	return &res, nil
}
