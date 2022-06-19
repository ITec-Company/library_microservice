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
	"strings"
)

const (
	getAllTagsURL   = "/tags"
	getTagByUUIDURL = "/tag/:uuid"
	GetManyByUUIDs  = "/tags/:uuids"
	createTagURL    = "/tag"
	deleteTagURL    = "/tag/:uuid"
	updateTagURL    = "/tag"
)

type TagHandler struct {
	Service    service.TagService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewTagHandler(service service.TagService, logger *logging.Logger) TagHandler {
	return TagHandler{
		Service: service,
		logger:  logger,
	}
}

func (th *TagHandler) Register(router *httprouter.Router) {
	router.GET(getAllTagsURL, th.GetAll)
	router.GET(GetManyByUUIDs, th.GetManyByUUIDs)
	router.GET(getTagByUUIDURL, th.GetByUUID)
	router.POST(createTagURL, th.Create)
	router.DELETE(deleteTagURL, th.Delete)
	router.PUT(updateTagURL, th.Update)
}

func (th *TagHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		th.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		th.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	tags, err := th.Service.GetAll(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all tags. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tags)
}

func (th *TagHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	tag, err := th.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting tag from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tag)
}

func (th *TagHandler) GetManyByUUIDs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuids")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("uuids can't be empty")
		return
	}
	UUIDs := strings.Split(uuid, ",")
	for _, d := range UUIDs {
		uuidInt, err := strconv.Atoi(d)
		if err != nil || uuidInt < 0 {
			w.WriteHeader(http.StatusBadRequest)
			th.logger.Errorf("uuid need to be uint")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
			return
		}
	}

	tag, err := th.Service.GetManyByUUIDs(UUIDs)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting tags from DB by UUIDs. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tag)
}

func (th *TagHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createTagDTO := &domain.CreateTagDTO{}
	if err := json.NewDecoder(r.Body).Decode(createTagDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}

	UUID, err := th.Service.Create(createTagDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating tag into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Tag created successfully. UUID: %s", UUID)})
}

func (th *TagHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = th.Service.Delete(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting tag from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Tag with UUID %s was deleted", uuid)})
}

func (th *TagHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateTagDTO := &domain.UpdateTagDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateTagDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateTagDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		th.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := th.Service.Update(updateTagDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating tag into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Tag updated successfully"})
}
