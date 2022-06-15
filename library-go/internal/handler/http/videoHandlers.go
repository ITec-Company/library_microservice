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
	getAllVideosURL    = "/videos"
	getVideoByUUIDURL  = "/video/:uuid"
	createVideoURL     = "/video"
	deleteVideoURL     = "/video/:uuid"
	updateVideoURL     = "/video"
	loadVideoFileURL   = "/file/video"
	updateVideoFileURL = "/file/video"
	rateVideoUrl       = "/rate/video"

	videoLocalStoragePath = "../store/videos/"
)

type videoHandler struct {
	Service    service.VideoService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewVideoHandler(service service.VideoService, logger *logging.Logger, middleware *Middleware) handler.Handler {
	return &videoHandler{
		Service:    service,
		logger:     logger,
		Middleware: middleware,
	}
}

func (vh *videoHandler) Register(router *httprouter.Router) {
	router.GET(getAllVideosURL, vh.Middleware.sortAndFilters(vh.GetAll()))
	router.GET(getVideoByUUIDURL, vh.GetByUUID)
	router.POST(createVideoURL, vh.Middleware.createVideo(vh.Create()))
	router.DELETE(deleteVideoURL, vh.Delete)
	router.PUT(updateVideoURL, vh.Update)
	//router.GET(loadVideoFileURL, vh.LoadFile)
	//router.PUT(updateVideoFileURL, vh.Middleware.updateVideoFile(vh.UpdateFile()))
	router.PUT(rateVideoUrl, vh.Rate)
}

func (vh *videoHandler) GetAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortingOptions := r.Context().Value(CtxKeySortAndFilters).(domain.SortFilterPagination)

		videos, pagesCount, err := vh.Service.GetAll(&sortingOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all videos. err: %v", err)})
			return
		}

		if pagesCount > 0 {
			w.Header().Set("pages", strconv.Itoa(pagesCount))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(videos)
	})
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

	video, err := vh.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting video from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(video)
}

func (vh *videoHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := r.Context().Value(CtxKeyCreateVideo).(map[string]interface{})

		createVideoDTO := domain.CreateVideoDTO{}

		createVideoDTO.Title = strings.Replace(data["title"].(string), " ", "_", -1)
		createVideoDTO.DirectionUUID = data["direction_uuid"].(string)
		createVideoDTO.Difficulty = data["difficulty"].(string)
		createVideoDTO.Language = data["language"].(string)
		createVideoDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
		fileName, ok := data["fileName"].(string)
		if ok {
			fileName = data["fileName"].(string)
			createVideoDTO.LocalURL = fmt.Sprintf("%s?url=%s", loadVideoFileURL, fileName)
		} else {
			createVideoDTO.LocalURL = "file wasn't added"
		}
		createVideoDTO.WebURL = data["web_url"].(string)

		file, ok := data["file"].(*bytes.Buffer)
		if ok && file != nil {
			err := vh.Service.SaveFile(videoLocalStoragePath, fileName, file)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				vh.logger.Errorf("error occurred while saving video into local database. err: %v.", err)
				json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving video into local database. err: %v", err)})
				return
			}
		}

		UUID, err := vh.Service.Create(&createVideoDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating video into DB. err: %v", err)})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Video created successfully. UUID: %s", UUID)})
	})
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

	err = vh.Service.Delete(uuid, fmt.Sprintf("%s%s/", videoLocalStoragePath, uuid))
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

	err := vh.Service.Update(updateVideoDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating video into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Video updated successfully"})
}

func (vh *videoHandler) LoadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("url can't be empty")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		vh.logger.Errorf("empty uuid")
	}

	if uuid != "" {
		err := vh.Service.DownloadCountUp(uuid)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			vh.logger.Errorf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)})
			return
		}
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := fmt.Sprintf("%s%s", videoLocalStoragePath, url)

	fileBytes, err := vh.Service.LoadFile(path)
	_, pathError := err.(*os.PathError)
	if pathError {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while searching file: invalid path. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while searching file: invalid path. err: %v", err)})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)
}

func (vh *videoHandler) UpdateFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := r.Context().Value(CtxKeyUpdateVideoFile).(domain.UpdateVideoFileDTO)

		dto.LocalPath = fmt.Sprintf("%s%s/", videoLocalStoragePath, dto.UUID)
		dto.LocalURL = fmt.Sprintf("%s?file=%s&uuid=%s", loadVideoFileURL, dto.NewFileName, dto.UUID)

		err := vh.Service.UpdateFile(&dto)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			vh.logger.Errorf("error occurred while saving video into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving video into local store. err: %v.", err)})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSON.Info{Msg: "File updated successfully"})
	})
}

func (vh *videoHandler) Rate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingStr := r.URL.Query().Get("rating")
	if ratingStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("rating query can't be empty")
		return
	}
	rating, err := strconv.ParseFloat(ratingStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("error occurred while parsing rating. Should be float32. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing rating. Should be float32. err: %v.", err)})
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		vh.logger.Errorf("uuid can't be empty")
		return
	}

	err = vh.Service.Rate(uuid, float32(rating))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		vh.logger.Errorf("error occurred while rating video image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while rating video into local store. err: %v.", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Video rated successfully. UUID: %s", uuid)})
}
