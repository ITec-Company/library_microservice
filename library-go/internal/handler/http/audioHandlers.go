package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"library-go/internal/domain"
	"library-go/internal/handler"
	"library-go/internal/service"
	"library-go/pkg/JSON"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

const (
	getAllAudiosURL       = "/audios"
	getAudioByUUIDURL     = "/audio/:uuid"
	createAudioURL        = "/audio"
	deleteAudioURL        = "/audio/:uuid"
	updateAudioURL        = "/audio"
	loadAudioURL          = "/load/audio/:url"
	audioLocalStoragePath = "../store/audios"
)

type audioHandler struct {
	Service service.AudioService
	logger  *logging.Logger
}

func NewAudioHandler(service service.AudioService, logger *logging.Logger) handler.Handler {
	return &audioHandler{
		Service: service,
		logger:  logger,
	}
}

func (ah *audioHandler) Register(router *httprouter.Router) {
	router.GET(getAllAudiosURL, ah.GetAll)
	router.GET(getAudioByUUIDURL, ah.GetByUUID)
	router.POST(createAudioURL, ah.Create)
	router.DELETE(deleteAudioURL, ah.Delete)
	router.PUT(updateAudioURL, ah.Update)
	router.GET(loadAudioURL, ah.Load)
}

func (ah *audioHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	audios, err := ah.Service.GetAll(context.Background(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all audios. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(audios)
}

func (ah *audioHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	audio, err := ah.Service.GetByUUID(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting audio from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(audio)
}

func (ah *audioHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createAudioDTO := domain.CreateAudioDTO{}

	data := map[string]interface{}{
		"file":           "file",
		"title":          "text",
		"direction_uuid": "text",
		"difficulty":     "text",
		"language":       "text",
		"tags_uuids":     "text",
	}

	err := utils.ParseMultiPartFormData(r, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
		return
	}

	file := data["file"].(*bytes.Buffer)
	createAudioDTO.Title = data["title"].(string)
	createAudioDTO.DirectionUUID = data["direction_uuid"].(string)
	createAudioDTO.Difficulty = data["difficulty"].(string)
	createAudioDTO.Language = data["language"].(string)
	createAudioDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
	createAudioDTO.URL = fmt.Sprintf("%s%s-%s-%s", "load/audio/", createAudioDTO.DirectionUUID, createAudioDTO.Difficulty, data["fileName"].(string))

	path := fmt.Sprintf("%s/%s/%s", audioLocalStoragePath, createAudioDTO.DirectionUUID, createAudioDTO.Difficulty)

	err = utils.SaveFile(path, data["fileName"].(string), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while saving audio into local database. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving audio into local database. err: %v", err)})
		return
	}

	UUID, err := ah.Service.Create(context.Background(), &createAudioDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating audio into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Audio created successfully. UUID: %s", UUID)})
}

func (ah *audioHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = ah.Service.Delete(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting audio from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Audio with UUID %s was deleted", uuid)})
}

func (ah *audioHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateAudioDTO := &domain.UpdateAudioDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateAudioDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateAudioDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := ah.Service.Update(context.Background(), updateAudioDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating audio into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Audio updated successfully"})
}

func (ah *audioHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	url := ps.ByName("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("url can't be empty")
		return
	}
	urlArray := strings.Split(url, "-")

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", urlArray[len(urlArray)-1]))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := strings.Replace(url, "-", "/", -1)

	fileBytes, err := ah.Service.LoadLocalFIle(context.Background(), fmt.Sprintf("%s/%s", audioLocalStoragePath, path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)

}
