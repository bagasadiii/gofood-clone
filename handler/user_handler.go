package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/service"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type UserHandlerImpl interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	GetUserHandler(w http.ResponseWriter, r *http.Request)
}
type UserHandler struct {
	service service.UserServiceImpl
	zap     *zap.Logger
}

func NewUserHandler(service service.UserServiceImpl, zap *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		zap:     zap,
	}
}

func (uh *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.zap.Error(http.StatusText(http.StatusMethodNotAllowed))
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}
	var input model.RegisterReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		uh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := uh.service.RegisterService(r.Context(), &input); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	uh.zap.Info(http.StatusText(http.StatusCreated), zap.String("user_registered", input.Username))
	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"username": input.Username,
		"email":    input.Email,
		"role":     input.Role,
		"Phone":    input.Phone,
		"name":     input.Username,
	})
}

func (uh *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.zap.Error(http.StatusText(http.StatusMethodNotAllowed))
		utils.JSONResponse(w, http.StatusMethodNotAllowed, nil)
		return
	}
	var input model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		uh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	token, err := uh.service.LoginService(r.Context(), &input)
	if err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	resp := map[string]string{
		"username": input.Username,
		"token":    token,
	}
	uh.zap.Info(http.StatusText(http.StatusOK), zap.String("User logged in", input.Username))
	utils.JSONResponse(w, http.StatusOK, resp)
}

func (uh *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	resp, err := uh.service.GetUserService(r.Context(), username)
	if err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	uh.zap.Info("User fetched", zap.String("Username", username))
	utils.JSONResponse(w, http.StatusOK, resp)
}
