package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bagasadiii/gofood-clone/model"
	"github.com/bagasadiii/gofood-clone/service"
	"github.com/bagasadiii/gofood-clone/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type DriverHandlerImpl interface {
	CreateDriverHandler(w http.ResponseWriter, r *http.Request)
	GetDriverHandler(w http.ResponseWriter, r *http.Request)
	UpdateDriverHandler(w http.ResponseWriter, r *http.Request)
}

type DriverHandler struct {
	service service.DriverServiceImpl
	zap     *zap.Logger
}

func NewDriverHandler(service service.DriverServiceImpl, zap *zap.Logger) *DriverHandler {
	return &DriverHandler{
		service: service,
		zap:     zap,
	}
}

func (dh *DriverHandler) CreateDriverHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Driver
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		dh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := dh.service.CreateDriverService(r.Context(), &input); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	dh.zap.Info("Driver created", zap.String("name", input.Name))
	utils.JSONResponse(w, http.StatusCreated, map[string]string{
		"name":    input.Name,
		"license": input.License,
		"area":    input.Area,
		"created": time.Now().String(),
	})
}

func (dh *DriverHandler) GetDriverHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	res, err := dh.service.GetDriverService(r.Context(), username)
	if err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
	}
	dh.zap.Info("user fetched", zap.String("username", username))
	utils.JSONResponse(w, http.StatusOK, res)
}

func (dh *DriverHandler) UpdateDriverHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Driver
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		dh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := dh.service.UpdateDriverService(r.Context(), &input); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	dh.zap.Info("Driver updated", zap.String("Driver", input.Name))
	utils.JSONResponse(w, http.StatusOK, &input)
}
