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

type MerchantHandlerImpl interface {
	CreateMerchantHandler(w http.ResponseWriter, r *http.Request)
	UpdateMerchantHandler(w http.ResponseWriter, r *http.Request)
	GetMerchantHandler(w http.ResponseWriter, r *http.Request)
}
type MerchantHandler struct {
	service service.MerchantServiceImpl
	zap     *zap.Logger
}

func NewMerchantHandler(service service.MerchantServiceImpl, zap *zap.Logger) *MerchantHandler {
	return &MerchantHandler{
		service: service,
		zap:     zap,
	}
}

func (mh *MerchantHandler) CreateMerchantHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Merchant
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		mh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := mh.service.CreateMerchantService(r.Context(), &input); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	resInput := map[string]string{
		"name":        input.Name,
		"address":     input.Address,
		"category":    input.Category,
		"description": input.Description,
	}
	utils.JSONResponse(w, http.StatusOK, resInput)
	mh.zap.Info("Merchant Created", zap.String("merchant", input.Name))
}

func (mh *MerchantHandler) GetMerchantHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	res, err := mh.service.GetMerchantService(r.Context(), username)
	if err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	mh.zap.Info("User fetched", zap.String("Merchant", res.Name))
	utils.JSONResponse(w, http.StatusOK, res)
}

func (mh *MerchantHandler) UpdateMerchantHandler(w http.ResponseWriter, r *http.Request) {
	var input model.Merchant
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || r.Body == nil {
		mh.zap.Error(utils.ErrBadRequest.Error(), zap.Error(utils.ErrBadRequest))
		utils.JSONResponse(w, http.StatusBadRequest, err)
		return
	}
	if err := mh.service.UpdateMerchantService(r.Context(), &input); err != nil {
		status, errIs := utils.ErrCheck(err)
		utils.JSONResponse(w, status, errIs)
		return
	}
	utils.JSONResponse(w, http.StatusOK, &input)
	mh.zap.Info("Merchant updated", zap.String("merchant", input.Name))
}
