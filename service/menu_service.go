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

type MenuServiceImpl interface {
	CreateMenuService(ctx context.Context, input *model.Menu, username string) error
	GetMenuService(ctx context.Context, id uuid.UUID) (*model.MenuRes, error)
	UpdateMenuService(ctx context.Context, data *model.Menu, username string) error
	DeleteMenuService(ctx context.Context, menuID uuid.UUID) error
}
type MenuService struct {
	repo repository.MenuRepoImpl
	zap  *zap.Logger
}

func NewMenuService(repo repository.MenuRepoImpl, zap *zap.Logger) *MenuService {
	return &MenuService{
		repo: repo,
		zap:  zap,
	}
}

func (ms *MenuService) CreateMenuService(ctx context.Context, input *model.Menu, username string) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ms.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("missing authorization: %w", utils.ErrUnauthorized)
	}
	if ctxValue.Username != username {
		ms.zap.Error(utils.ErrForbidden.Error(), zap.String("forbidden", username))
		return fmt.Errorf("not allowed to access: %w", utils.ErrForbidden)
	}
	if ctxValue.Role != "merchant" {
		ms.zap.Error("invalid role", zap.String("needed", "merchant"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%w: role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	merchantID, err := ms.repo.GetMerchantID(ctx, ctxValue.UserID)
	if err != nil {
		ms.zap.Error(utils.ErrForbidden.Error(), zap.String("forbidden", "invalid merchant and user id"))
		return fmt.Errorf("not allowed to access: %w", utils.ErrForbidden)
	}

	newMenu := model.Menu{
		MenuID:      uuid.New(),
		Price:       input.Price,
		Description: input.Description,
		Category:    input.Category,
		Rating:      0,
		Stock:       input.Stock,
		MerchantID:  merchantID,
	}
	return ms.repo.CreateMenuRepo(ctx, &newMenu, ctxValue.UserID)
}

func (ms *MenuService) GetMenuService(ctx context.Context, id uuid.UUID) (*model.MenuRes, error) {
	return ms.repo.GetMenuRepo(ctx, id)
}

func (ms *MenuService) UpdateMenuService(ctx context.Context, data *model.Menu, username string) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ms.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("missing authorization: %w", utils.ErrUnauthorized)
	}
	if ctxValue.Username != username {
		ms.zap.Error(utils.ErrForbidden.Error(), zap.String("forbidden", username))
		return fmt.Errorf("not allowed to access: %w", utils.ErrForbidden)
	}
	if ctxValue.Role != "merchant" {
		ms.zap.Error("invalid role", zap.String("needed", "merchant"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%w: role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	merchantID, err := ms.repo.GetMerchantID(ctx, ctxValue.UserID)
	if err != nil {
		ms.zap.Error(utils.ErrForbidden.Error(), zap.String("forbidden", "invalid merchant and user id"))
		return fmt.Errorf("not allowed to access: %w", utils.ErrForbidden)
	}
	query, args := updateMenuBuilder(data, merchantID)
	return ms.repo.UpdateMenuRepo(ctx, query, args)
}

func (ms *MenuService) DeleteMenuService(ctx context.Context, menuID uuid.UUID) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ms.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("missing authorization: %w", utils.ErrUnauthorized)
	}
	if ctxValue.Role != "merchant" {
		ms.zap.Error("invalid role", zap.String("needed", "merchant"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%w: role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	merchantID, err := ms.repo.GetMerchantID(ctx, ctxValue.UserID)
	if err != nil {
		ms.zap.Error(utils.ErrForbidden.Error(), zap.String("forbidden", "invalid merchant and user id"))
		return fmt.Errorf("not allowed to access: %w", utils.ErrForbidden)
	}
	return ms.repo.DeleteMenuRepo(ctx, menuID, merchantID)
}

func updateMenuBuilder(updated *model.Menu, merchantID uuid.UUID) (string, []interface{}) {
	fields := []string{}
	argsIndex := 1
	args := []interface{}{}

	if updated.Name != "" {
		fields = append(fields, fmt.Sprintf("name = $%d", argsIndex))
		args = append(args, updated.Name)
		argsIndex++
	}
	if updated.Price != 0 {
		fields = append(fields, fmt.Sprintf("price = $%d", argsIndex))
		args = append(args, updated.Price)
		argsIndex++
	}
	if updated.Category != "" {
		fields = append(fields, fmt.Sprintf("category = $%d", argsIndex))
		args = append(args, updated.Category)
		argsIndex++
	}
	if updated.Description != "" {
		fields = append(fields, fmt.Sprintf("description = $%d", argsIndex))
		args = append(args, updated.Description)
		argsIndex++
	}
	if updated.Stock != 0 {
		fields = append(fields, fmt.Sprintf("stock = $%d", argsIndex))
		args = append(args, updated.Stock)
		argsIndex++
	}

	args = append(args, merchantID)
	updatedQuery := fmt.Sprintf("%s WHERE merchant_id = $%d", strings.Join(fields, ", "), argsIndex)
	return updatedQuery, args
}
