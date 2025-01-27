package repository

import (
	"context"
	"fmt"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DriverRepoImpl interface {
	CreateDriverRepo(ctx context.Context, new *model.Driver) error
	GetDriverRepo(ctx context.Context, username string) (*model.DriverRes, error)
	UpdateDriverRepo(ctx context.Context, query string, args []interface{}) error
}
type DriverRepo struct {
	db  *pgxpool.Pool
	zap *zap.Logger
}

func NewDriverRepo(db *pgxpool.Pool, zap *zap.Logger) *DriverRepo {
	return &DriverRepo{
		db:  db,
		zap: zap,
	}
}

func (dr *DriverRepo) CreateDriverRepo(ctx context.Context, new *model.Driver) error {
	var exists bool
	err := dr.db.QueryRow(ctx, `
    SELECT EXISTS (SELECT 1 FROM drivers WHERE user_id = $1 OR username = $2)
    `, new.DriverID, new.Username).Scan(&exists)
	if err != nil {
		dr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%v", utils.ErrDatabase)
	}
	if exists {
		dr.zap.Warn(utils.ErrUniqueConstraint.Error(), zap.String("Merchant exists", new.Name))
		return fmt.Errorf("%v", utils.ErrUniqueConstraint)
	}

	_, err = dr.db.Exec(ctx, `
    INSERT INTO drivers (driver_id, name, rating, license, area, income, user_id, username)
    VALUES ($1, $2, $3, $4 ,$5 , $6, $7, $8)
    `, new.DriverID, new.Name, new.Rating, new.License, new.Area, new.Income, new.UserID, new.Username)
	if err != nil {
		dr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%v", utils.ErrDatabase)
	}
	return nil
}

func (dr *DriverRepo) GetDriverRepo(ctx context.Context, username string) (*model.DriverRes, error) {
	var res model.DriverRes
	row := dr.db.QueryRow(ctx, `
    SELECT name, rating, license, area, income FROM drivers
    WHERE username = $1
    `, username)
	err := row.Scan(res.Name, res.Rating, res.Area, res.Income)
	if err != pgx.ErrNoRows {
		dr.zap.Warn(utils.ErrNotFound.Error(), zap.String("Username", username))
		return nil, fmt.Errorf("%v", utils.ErrNotFound)
	} else if err != nil {
		dr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return nil, fmt.Errorf("%v", utils.ErrDatabase)
	}
	return &res, nil
}

func (dr *DriverRepo) UpdateDriverRepo(ctx context.Context, query string, args []interface{}) error {
	_, err := dr.db.Exec(ctx, fmt.Sprintf(`UPDATE drivers SET %s`, query), args)
	if err != nil {
		dr.zap.Error(utils.ErrDatabase.Error(), zap.Error(err))
		return fmt.Errorf("%v", utils.ErrDatabase)
	}
	return nil
}
