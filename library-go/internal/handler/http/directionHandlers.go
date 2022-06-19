package http

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"library-go/internal/domain"
	"library-go/internal/service"
	"library-go/pkg/JSON"
	"library-go/pkg/logging"
	"net/http"
	"strconv"
)

const (
	getAllDirectionsURL   = "/directions"
	getDirectionByUUIDURL = "/direction/:uuid"
	createDirectionURL    = "/direction"
	deleteDirectionURL    = "/direction/:uuid"
	updateDirectionURL    = "/direction"
)

type DirectionHandler struct {
	Service service.DirectionService
	logger  *logging.Logger
}

func NewDirectionHandler(service service.DirectionService, logger *logging.Logger) DirectionHandler {
	return DirectionHandler{
		Service: service,
		logger:  logger,
	}
}

func (dh *DirectionHandler) Register(router *httprouter.Router) {
	router.GET(getAllDirectionsURL, dh.GetAll)
	router.GET(getDirectionByUUIDURL, dh.GetByUUID)
	router.POST(createDirectionURL, dh.Create)
	router.DELETE(deleteDirectionURL, dh.Delete)
	router.PUT(updateDirectionURL, dh.Update)
}

func (dh *DirectionHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		dh.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		dh.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	directions, err := dh.Service.GetAll(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all directions. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(directions)
}

func (dh *DirectionHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	direction, err := dh.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting direction from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(direction)
}

func (dh *DirectionHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createDirectionDTO := &domain.CreateDirectionDTO{}
	if err := json.NewDecoder(r.Body).Decode(createDirectionDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}

	UUID, err := dh.Service.Create(createDirectionDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating direction into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Direction created successfully. UUID: %s", UUID)})
}

func (dh *DirectionHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = dh.Service.Delete(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting direction from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Direction with UUID %s was deleted", uuid)})
}

func (dh *DirectionHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateDirectionDTO := &domain.UpdateDirectionDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateDirectionDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateDirectionDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		dh.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := dh.Service.Update(updateDirectionDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating direction into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Direction updated successfully"})
}
