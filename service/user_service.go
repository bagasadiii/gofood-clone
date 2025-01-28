package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/bagasadiii/gofood-clone/middleware"
	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/repository"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl interface {
	RegisterService(ctx context.Context, input *model.RegisterReq) error
	GetUserService(ctx context.Context, username string) (*model.UserResp, error)
	LoginService(ctx context.Context, input *model.LoginReq) (string, error)
}
type UserService struct {
	userRepo   repository.UserRepoImpl
	zap        *zap.Logger
	jwtService middleware.JWTServiceImpl
}

func NewUserService(repo repository.UserRepoImpl, zap *zap.Logger, jwt middleware.JWTServiceImpl) *UserService {
	return &UserService{
		userRepo:   repo,
		zap:        zap,
		jwtService: jwt,
	}
}

func (us *UserService) RegisterService(ctx context.Context, input *model.RegisterReq) error {
	re := regexp.MustCompile("^[a-z0-9_]+$")
	if !re.MatchString(input.Username) {
		us.zap.Warn(utils.ErrBadRequest.Error(), zap.String("Invalid username", input.Username))
		return fmt.Errorf("invalid username:%v", utils.ErrBadRequest)
	}
	if err := utils.ValidateUser(input); err != nil {
		us.zap.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return fmt.Errorf("%v: %v", utils.ErrBadRequest, err)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		us.zap.Error(utils.ErrInternal.Error(), zap.Error(err))
		return fmt.Errorf("%v", utils.ErrInternal)
	}
	newUser := &model.User{
		UserID:    uuid.New(),
		Username:  input.Username,
		Email:     input.Email,
		Password:  string(hashed),
		Role:      input.Role,
		CreatedAt: time.Now(),
		Phone:     input.Phone,
		Balance:   0,
		Name:      input.Username,
	}
	if err := us.userRepo.RegisterUserRepo(ctx, newUser); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserService(ctx context.Context, username string) (*model.UserResp, error) {
	id, err := us.userRepo.GetIDRepo(ctx, username)
	if err != nil {
		return nil, err
	}
	resp, err := us.userRepo.GetUserRepo(ctx, id)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (us *UserService) LoginService(ctx context.Context, input *model.LoginReq) (string, error) {
	if err := utils.ValidateLogin(input); err != nil {
		us.zap.Error(utils.ErrBadRequest.Error(), zap.Error(err))
		return "", fmt.Errorf("%v: %v", utils.ErrBadRequest, err)
	}
	res, err := us.userRepo.LoginRepo(ctx, input.Username)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(input.Password)); err != nil {
		us.zap.Warn(utils.ErrInvalidPassword.Error())
		return "", utils.ErrInvalidPassword
	}
	newClaims := &middleware.TokenClaims{
		UserID:   res.UserID,
		Username: res.Username,
		Role:     res.Role,
	}
	token, err := us.jwtService.CreateToken(newClaims)
	if err != nil {
		return "", err
	}
	return token, nil
}
