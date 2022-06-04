package http

import (
	"bytes"
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
	getAllAudiosURL    = "/audios"
	getAudioByUUIDURL  = "/audio/:uuid"
	createAudioURL     = "/audio"
	deleteAudioURL     = "/audio/:uuid"
	updateAudioURL     = "/audio"
	loadAudioFileURL   = "/file/audio"
	updateAudioFileURL = "/file/audio"

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
	router.GET(loadAudioFileURL, ah.LoadFile)
	router.PUT(updateAudioFileURL, ah.Middleware.updateAudioFile(ah.UpdateFile()))
}

func (ah *audioHandler) GetAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortingOptions := r.Context().Value(CtxKeySortAndFilters).(domain.SortFilterPagination)

		audios, pagesCount, err := ah.Service.GetAll(&sortingOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all audios. err: %v", err)})
			return
		}

		if pagesCount > 0 {
			w.Header().Set("pages", strconv.Itoa(pagesCount))
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

	audio, err := ah.Service.GetByUUID(uuid)
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

		createAudioDTO.Title = strings.Replace(data["title"].(string), " ", "_", -1)
		createAudioDTO.DirectionUUID = data["direction_uuid"].(string)
		createAudioDTO.Difficulty = data["difficulty"].(string)
		createAudioDTO.Language = data["language"].(string)
		createAudioDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
		fileName := data["fileName"].(string)
		createAudioDTO.LocalURL = fmt.Sprintf("%s?file=%s&uuid=", loadAudioFileURL, fileName)

		UUID, err := ah.Service.Create(&createAudioDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating audio into DB. err: %v", err)})
			return
		}

		path := fmt.Sprintf("%s%s/", audioLocalStoragePath, UUID)

		file := data["file"].(*bytes.Buffer)
		err = ah.Service.SaveFile(path, fileName, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while saving audio into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving audio into local database. err: %v", err)})
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

	err = ah.Service.Delete(uuid, fmt.Sprintf("%s%s/", audioLocalStoragePath, uuid))
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

	// block changing URL
	updateAudioDTO.LocalURL = ""

	err := ah.Service.Update(updateAudioDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating audio into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Audio updated successfully"})
}

func (ah *audioHandler) LoadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file := r.URL.Query().Get("file")
	if file == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("file query can't be empty")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid query can't be empty")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := fmt.Sprintf("%s%s/%s", audioLocalStoragePath, uuid, file)

	fileBytes, err := ah.Service.LoadFile(path)
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

func (ah *audioHandler) UpdateFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := r.Context().Value(CtxKeyUpdateAudioFile).(domain.UpdateAudioFileDTO)

		dto.LocalPath = fmt.Sprintf("%s%s/", audioLocalStoragePath, dto.UUID)
		dto.LocalURL = fmt.Sprintf("%s?file=%s&uuid=%s", loadAudioFileURL, dto.NewFileName, dto.UUID)

		err := ah.Service.UpdateFile(&dto)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while saving audio into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving audio into local store. err: %v.", err)})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSON.Info{Msg: "File updated successfully"})
	})
}
