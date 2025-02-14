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

type MenuHandler struct {
	service service.MenuServiceImpl
	zap     *zap.Logger
}

func NewMenuHandler(service service.MenuServiceImpl, zap *zap.Logger) *MenuHandler {
	return &MenuHandler{
		service: service,
		zap:     zap,
	}
}

func (mh *MenuHandler) CreateMenuHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Menu
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		mh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	vars := mux.Vars(r)
	username := vars["username"]
	if err := mh.service.CreateMenuService(r.Context(), &input, username); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	mh.zap.Info("menu created", zap.Any("menu", &input))
	utils.JSONResponse(w, http.StatusCreated, &input)
}
