package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/repository"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DriverServiceImpl interface {
	CreateDriverService(ctx context.Context, new *model.Driver) error
	GetDriverService(ctx context.Context, username string) (*model.DriverRes, error)
	UpdateDriverService(ctx context.Context, update *model.Driver) error
}

type DriverService struct {
	repo repository.DriverRepoImpl
	zap  *zap.Logger
}

func NewDriverService(repo repository.DriverRepoImpl, zap *zap.Logger) *DriverService {
	return &DriverService{
		repo: repo,
		zap:  zap,
	}
}

func (ds *DriverService) CreateDriverService(ctx context.Context, new *model.Driver) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ds.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("%v", err)
	}
	if ctxValue.Role != "driver" {
		ds.zap.Error("invalid role", zap.String("needed", "driver"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%v: user role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	newDriver := model.Driver{
		DriverID: uuid.New(),
		Name:     new.Name,
		License:  new.License,
		Area:     new.Area,
		Income:   0,
		UserID:   ctxValue.UserID,
		Username: ctxValue.Username,
	}
	if err := utils.ValidateDriver(&newDriver); err != nil {
		ds.zap.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return fmt.Errorf("%v", err)
	}
	return ds.repo.CreateDriverRepo(ctx, &newDriver)
}

func (ds *DriverService) GetDriverService(ctx context.Context, username string) (*model.DriverRes, error) {
	return ds.repo.GetDriverRepo(ctx, username)
}

func (ds *DriverService) UpdateDriverService(ctx context.Context, update *model.Driver) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ds.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("%v", err)
	}
	if ctxValue.Role != "driver" {
		ds.zap.Error("invalid role", zap.String("needed", "driver"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%v: user role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	update.UserID = ctxValue.UserID
	update.Username = ctxValue.Username
	query, args := updateDriverQueryBuilder(update)
	return ds.repo.UpdateDriverRepo(ctx, query, args)
}

func updateDriverQueryBuilder(updated *model.Driver) (string, []interface{}) {
	fields := []string{}
	argsIndex := 1
	args := []interface{}{}

	if updated.Name != "" {
		fields = append(fields, fmt.Sprintf("name = $%d", argsIndex))
		args = append(args, updated.Name)
		argsIndex++
	}
	if updated.License != "" {
		fields = append(fields, fmt.Sprintf("license = $%d", argsIndex))
		args = append(args, updated.License)
		argsIndex++
	}
	if updated.Area != "" {
		fields = append(fields, fmt.Sprintf("area = $%d", argsIndex))
		args = append(args, updated.Area)
		argsIndex++
	}
	args = append(args, updated.UserID)
	updatedQuery := fmt.Sprintf("%s WHERE user_id = $%d", strings.Join(fields, ", "), argsIndex)
	return updatedQuery, args
}
