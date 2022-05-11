package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"library-go/internal/domain"
	"library-go/internal/handler"
	"library-go/internal/service"
	"library-go/pkg/JSON"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	getAllVideosURL       = "/videos"
	getVideoByUUIDURL     = "/video/:uuid"
	createVideoURL        = "/video"
	deleteVideoURL        = "/video/:uuid"
	updateVideoURL        = "/video"
	loadVideoURL          = "/load/video/:url"
	videoLocalStoragePath = "../store/videos"
)

type videoHandler struct {
	Service service.VideoService
	logger  *logging.Logger
}

func NewVideoHandler(service service.VideoService, logger *logging.Logger) handler.Handler {
	return &videoHandler{
		Service: service,
		logger:  logger,
	}
}

func (vh *videoHandler) Register(router *httprouter.Router) {
	router.GET(getAllVideosURL, vh.GetAll)
	router.GET(getVideoByUUIDURL, vh.GetByUUID)
	router.POST(createVideoURL, vh.Create)
	router.DELETE(deleteVideoURL, vh.Delete)
	router.PUT(updateVideoURL, vh.Update)
	router.GET(loadVideoURL, vh.Load)
}

func (vh *videoHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		vh.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		vh.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	videos, err := vh.Service.GetAll(context.Background(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all videos. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}

func (vh *videoHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	video, err := vh.Service.GetByUUID(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting video from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(video)
}

func (vh *videoHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createVideoDTO := domain.CreateVideoDTO{}

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
		vh.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
		return
	}

	file := data["file"].(*bytes.Buffer)
	createVideoDTO.Title = data["title"].(string)
	createVideoDTO.DirectionUUID = data["direction_uuid"].(string)
	createVideoDTO.Difficulty = data["difficulty"].(string)
	createVideoDTO.Language = data["language"].(string)
	createVideoDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
	createVideoDTO.URL = fmt.Sprintf("%s%s-%s-%s", "load/video/", createVideoDTO.DirectionUUID, createVideoDTO.Difficulty, data["fileName"].(string))

	path := fmt.Sprintf("%s/%s/%s", videoLocalStoragePath, createVideoDTO.DirectionUUID, createVideoDTO.Difficulty)

	err = utils.SaveFile(path, data["fileName"].(string), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while saving video into local database. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving video into local database. err: %v", err)})
		return
	}

	UUID, err := vh.Service.Create(context.Background(), &createVideoDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating video into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Video created successfully. UUID: %s", UUID)})
}

func (vh *videoHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = vh.Service.Delete(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting video from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Video with UUID %s was deleted", uuid)})
}

func (vh *videoHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateVideoDTO := &domain.UpdateVideoDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateVideoDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateVideoDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := vh.Service.Update(context.Background(), updateVideoDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating video into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Video updated successfully"})
}

func (vh *videoHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	url := ps.ByName("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("url can't be empty")
		return
	}
	urlArray := strings.Split(url, "-")

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", urlArray[len(urlArray)-1]))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := strings.Replace(url, "-", "/", -1)

	file, err := os.Open(fmt.Sprintf("%s/%s", videoLocalStoragePath, path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)

}
