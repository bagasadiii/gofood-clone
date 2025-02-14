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

type MerchantRepoImpl interface {
	CreateMerchantRepo(ctx context.Context, new *model.Merchant) error
	GetMerchantRepo(ctx context.Context, username string) (*model.MerchantRes, error)
	UpdateMerchantRepo(ctx context.Context, query string, args []interface{}) error
}
type MerchantRepo struct {
	db  *pgxpool.Pool
	zap *zap.Logger
}

func NewMerchantRepo(db *pgxpool.Pool, zap *zap.Logger) *MerchantRepo {
	return &MerchantRepo{
		db:  db,
		zap: zap,
	}
}

func (mr *MerchantRepo) CreateMerchantRepo(ctx context.Context, new *model.Merchant) error {
	var exists bool
	err := mr.db.QueryRow(ctx, `
    SELECT EXISTS (SELECT 1 FROM merchants WHERE user_id = $1 OR owner = $2)
    `, new.MerchantID, new.Owner).Scan(&exists)
	if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%w:%w", utils.ErrUnexpected, utils.ErrDatabase)
	}
	if exists {
		mr.zap.Warn(utils.ErrUniqueConstraint.Error(), zap.String("merchant exists", new.Name))
		return fmt.Errorf("merchant already exists: %w", utils.ErrUniqueConstraint)
	}
	_, err = mr.db.Exec(ctx, `
    INSERT INTO merchants (merchant_id, name, rating, address, category, user_id, owner)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, new.MerchantID, new.Name, new.Rating, new.Address, new.Category, new.UserID, new.Owner)
	if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("failed to create merchant: %w", utils.ErrDatabase)
	}
	return nil
}

func (mr *MerchantRepo) GetMerchantRepo(ctx context.Context, username string) (*model.MerchantRes, error) {
	var id uuid.UUID
	err := mr.db.QueryRow(ctx, `
    SELECT merchant_id FROM merchants WHERE owner = $1
    `, username).Scan(&id)
	if err == pgx.ErrNoRows {
		mr.zap.Warn(utils.ErrNotFound.Error(), zap.String("Username", username))
		return nil, fmt.Errorf("merchant not exists: %war", utils.ErrNotFound)
	} else if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("%w: %w", utils.ErrUnexpected, utils.ErrDatabase)
	}
	var res model.MerchantRes
	err = mr.db.QueryRow(ctx, `
    SELECT (name, rating, address, category) FROM merchants WHERE merchant_id = $1
    `, id).Scan(&res.Name, &res.Rating, &res.Address, &res.Category)
	if err == pgx.ErrNoRows {
		mr.zap.Warn(utils.ErrNotFound.Error(), zap.String("MerchantID", id.String()))
		return nil, fmt.Errorf("merchant not exists: %w", utils.ErrNotFound)
	} else if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch merchant: %w", utils.ErrDatabase)
	}
	return &res, nil
}

func (mr *MerchantRepo) UpdateMerchantRepo(ctx context.Context, query string, args []interface{}) error {
	_, err := mr.db.Exec(ctx, fmt.Sprintf(`UPDATE merchants SET %s`, query), args)
	if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%w", utils.ErrDatabase)
	}
	return nil
}
