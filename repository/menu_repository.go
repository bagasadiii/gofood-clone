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

type MenuRepoImpl interface {
	CreateMenuRepo(ctx context.Context, new *model.Menu, userID uuid.UUID) error
	UpdateMenuRepo(ctx context.Context, query string, args []interface{}) error
	GetMenuRepo(ctx context.Context, id uuid.UUID) (*model.MenuRes, error)
	DeleteMenuRepo(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error
	GetMerchantID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
}
type MenuRepo struct {
	db  *pgxpool.Pool
	zap *zap.Logger
}

func NewMenuRepo(db *pgxpool.Pool, zap *zap.Logger) *MenuRepo {
	return &MenuRepo{
		db:  db,
		zap: zap,
	}
}

func (mr *MenuRepo) CreateMenuRepo(ctx context.Context, new *model.Menu, userID uuid.UUID) error {
	_, err := mr.db.Exec(ctx, `
    INSERT INTO menus (menu_id, name, description, price, category, rating, stock, merchant_id)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, new.MenuID, new.Name, new.Description, new.Price, new.Category, new.Rating, new.Stock, new.MerchantID)
	if err != nil {
		mr.zap.Warn(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("failed to create menu: %w", utils.ErrDatabase)
	}
	return nil
}

func (mr *MenuRepo) GetMenuRepo(ctx context.Context, id uuid.UUID) (*model.MenuRes, error) {
	var res model.MenuRes
	err := mr.db.QueryRow(ctx, `
    SELECT name, price, description, category, rating, stock FROM menus WHERE menu_id = $1
    `, id).Scan(&res.Name, &res.Price, &res.Description, &res.Category, &res.Rating, &res.Stock)
	if err == pgx.ErrNoRows {
		mr.zap.Warn(utils.ErrNotFound.Error(), zap.String("menu_id", id.String()))
		return nil, fmt.Errorf("menu not found: %w", utils.ErrNotFound)
	} else if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("error while fetching menu: %w", utils.ErrNotFound)
	}
	return &res, nil
}

func (mr *MenuRepo) UpdateMenuRepo(ctx context.Context, query string, args []interface{}) error {
	_, err := mr.db.Exec(ctx, query, args...)
	if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("failed to update menu: %w", utils.ErrDatabase)
	}
	return nil
}

func (mr *MenuRepo) DeleteMenuRepo(ctx context.Context, id uuid.UUID, merchantID uuid.UUID) error {
	_, err := mr.db.Exec(ctx, `
    DELETE FROM menus WHERE menu_id = $1 AND merchant_id = $2
    `, id)
	if err != nil {
		mr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("failed to delete menu: %w", utils.ErrDatabase)
	}
	return nil
}

func (mr *MenuRepo) GetMerchantID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var merchantID uuid.UUID
	err := mr.db.QueryRow(ctx, `
    SELECT merchant_id FROM merchants WHERE user_id = $1
    `, userID).Scan(&merchantID)
	if err == pgx.ErrNoRows {
		mr.zap.Warn(utils.ErrNotFound.Error(), zap.String("no merchant_id found", userID.String()))
		return uuid.Nil, fmt.Errorf("merchant not found: %w", utils.ErrNotFound)
	} else if err != nil {
		mr.zap.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return uuid.Nil, fmt.Errorf("%w", utils.ErrDatabase)
	}
	return merchantID, nil
}
