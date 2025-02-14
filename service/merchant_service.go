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

type MerchantServiceImpl interface {
	CreateMerchantService(ctx context.Context, new *model.Merchant) error
	GetMerchantService(ctx context.Context, username string) (*model.MerchantRes, error)
	UpdateMerchantService(ctx context.Context, update *model.Merchant) error
}
type MerchantService struct {
	repo repository.MerchantRepoImpl
	zap  *zap.Logger
}

func NewMerchantService(repo repository.MerchantRepoImpl, zap *zap.Logger) *MerchantService {
	return &MerchantService{
		repo: repo,
		zap:  zap,
	}
}

func (ms *MerchantService) CreateMerchantService(ctx context.Context, new *model.Merchant) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ms.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("%w", err)
	}
	if ctxValue.Role != "merchant" {
		ms.zap.Error("invalid role", zap.String("needed", "merchant"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%w: role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	newMerchant := model.Merchant{
		MerchantID:  uuid.New(),
		Name:        new.Name,
		Rating:      0,
		Address:     new.Address,
		Category:    new.Category,
		Description: new.Description,
		UserID:      ctxValue.UserID,
		Owner:       ctxValue.Username,
	}
	if err := utils.ValidateMerchant(&newMerchant); err != nil {
		ms.zap.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return fmt.Errorf("%w", err)
	}
	return ms.repo.CreateMerchantRepo(ctx, &newMerchant)
}

func (ms *MerchantService) GetMerchantService(ctx context.Context, username string) (*model.MerchantRes, error) {
	return ms.repo.GetMerchantRepo(ctx, username)
}

func (ms *MerchantService) UpdateMerchantService(ctx context.Context, update *model.Merchant) error {
	ctxValue, err := utils.CheckContextValue(ctx)
	if err != nil {
		ms.zap.Error(utils.ErrUnauthorized.Error(), zap.Error(err))
		return fmt.Errorf("%w", err)
	}
	if ctxValue.Role != "merchant" {
		ms.zap.Error("invalid role", zap.String("needed", "merchant"), zap.String("actual", ctxValue.Role))
		return fmt.Errorf("%w: user role %s is not allowed", utils.ErrUnauthorized, ctxValue.Role)
	}
	query, args := updateMerchantQueryBuilder(update)
	return ms.repo.UpdateMerchantRepo(ctx, query, args)
}

func updateMerchantQueryBuilder(updated *model.Merchant) (string, []interface{}) {
	fields := []string{}
	argsIndex := 1
	args := []interface{}{}

	if updated.Name != "" {
		fields = append(fields, fmt.Sprintf("name = $%d", argsIndex))
		args = append(args, updated.Name)
		argsIndex++
	}
	if updated.Address != "" {
		fields = append(fields, fmt.Sprintf("address = $%d", argsIndex))
		args = append(args, updated.Address)
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
	args = append(args, updated.UserID)
	updatedQuery := fmt.Sprintf("%s WHERE user_id = $%d", strings.Join(fields, ", "), argsIndex)
	return updatedQuery, args
}
