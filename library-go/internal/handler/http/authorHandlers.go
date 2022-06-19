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
	getAllAuthorsURL   = "/authors"
	getAuthorByUUIDURL = "/author/:uuid"
	createAuthorURL    = "/author"
	deleteAuthorURL    = "/author/:uuid"
	updateAuthorURL    = "/author"
)

type AuthorHandler struct {
	Service service.AuthorService
	logger  *logging.Logger
}

func NewAuthorHandler(service service.AuthorService, logger *logging.Logger) AuthorHandler {
	return AuthorHandler{
		Service: service,
		logger:  logger,
	}
}

func (ah *AuthorHandler) Register(router *httprouter.Router) {
	router.GET(getAllAuthorsURL, ah.GetAll)
	router.GET(getAuthorByUUIDURL, ah.GetByUUID)
	router.POST(createAuthorURL, ah.Create)
	router.DELETE(deleteAuthorURL, ah.Delete)
	router.PUT(updateAuthorURL, ah.Update)
}

func (ah *AuthorHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		ah.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		ah.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	authors, err := ah.Service.GetAll(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all authors. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authors)
}

func (ah *AuthorHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid can't be empty")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid can't be empty"})
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	author, err := ah.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting author from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(author)
}

func (ah *AuthorHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createAuthorDTO := &domain.CreateAuthorDTO{}
	if err := json.NewDecoder(r.Body).Decode(createAuthorDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}

	UUID, err := ah.Service.Create(createAuthorDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating author into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Author created successfully. UUID: %s", UUID)})
}

func (ah *AuthorHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = ah.Service.Delete(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting author from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Author with UUID %s was deleted", uuid)})
}

func (ah *AuthorHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateAuthorDTO := &domain.UpdateAuthorDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateAuthorDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateAuthorDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := ah.Service.Update(updateAuthorDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating author into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Author updated successfully"})
}
