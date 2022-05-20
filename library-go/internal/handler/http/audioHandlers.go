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
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	getAllAudiosURL       = "/audios"
	getAudioByUUIDURL     = "/audio/:uuid"
	createAudioURL        = "/audio"
	deleteAudioURL        = "/audio/:uuid"
	updateAudioURL        = "/audio"
	loadAudioURL          = "/load/audio"
	audioLocalStoragePath = "../store/audios/"
)

type audioHandler struct {
	Service    service.AudioService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewAudioHandler(service service.AudioService, logger *logging.Logger, middleware *Middleware) handler.Handler {
	return &audioHandler{
		Service:    service,
		logger:     logger,
		Middleware: middleware,
	}
}

func (ah *audioHandler) Register(router *httprouter.Router) {
	router.GET(getAllAudiosURL, ah.Middleware.sortAndFilters(ah.GetAll()))
	router.GET(getAudioByUUIDURL, ah.GetByUUID)
	router.POST(createAudioURL, ah.Middleware.createAudio(ah.Create()))
	router.DELETE(deleteAudioURL, ah.Delete)
	router.PUT(updateAudioURL, ah.Update)
	router.GET(loadAudioURL, ah.Load)
}

func (ah *audioHandler) GetAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortingOptions := r.Context().Value(CtxKeySortAndFilters).(domain.SortFilterPagination)

		audios, err := ah.Service.GetAll(&sortingOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all audios. err: %v", err)})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(audios)
	})
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

func (ah *audioHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := r.Context().Value(CtxKeyCreateAudio).(map[string]interface{})

		createAudioDTO := domain.CreateAudioDTO{}

		file := data["file"].(*bytes.Buffer)

		createAudioDTO.Title = strings.Replace(data["title"].(string), " ", "_", -1)
		createAudioDTO.DirectionUUID = data["direction_uuid"].(string)
		createAudioDTO.Difficulty = data["difficulty"].(string)
		createAudioDTO.Language = data["language"].(string)
		createAudioDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")

		fileName := data["fileName"].(string)
		createAudioDTO.URL = fmt.Sprintf("%s?url=%s", loadBookURL, fileName)

		err := ah.Service.Save(context.Background(), audioLocalStoragePath, fileName, file)
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
	})
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
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("url can't be empty")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := fmt.Sprintf("%s%s", audioLocalStoragePath, url)

	fileBytes, err := ah.Service.Load(context.Background(), path)
	_, pathError := err.(*os.PathError)
	if pathError {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while searching file: invalid path. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while searching file: invalid path. err: %v", err)})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)

}
