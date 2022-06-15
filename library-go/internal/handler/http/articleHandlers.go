package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/jpeg"
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
	getAllArticlesURL     = "/articles"
	getArticleByUUIDURL   = "/article/:uuid"
	createArticleURL      = "/article"
	deleteArticleURL      = "/article/:uuid"
	updateArticleURL      = "/article"
	updateArticleFileURL  = "/file/article"
	updateArticleImageURL = "/image/article"
	rateArticleUrl        = "/rate/article"

	loadArticleFileURL  = "/articles/"
	loadArticleImageURL = "/articles/"

	articleLocalStoragePath = "../store/articles/"
)

type articleHandler struct {
	Service    service.ArticleService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewArticleHandler(service service.ArticleService, logger *logging.Logger, middleware *Middleware) handler.Handler {
	return &articleHandler{
		Service:    service,
		logger:     logger,
		Middleware: middleware,
	}
}

func (ah *articleHandler) Register(router *httprouter.Router) {
	router.GET(getAllArticlesURL, ah.Middleware.sortAndFilters(ah.GetAll()))
	router.GET(getArticleByUUIDURL, ah.GetByUUID)
	router.POST(createArticleURL, ah.Middleware.createArticle(ah.Create()))
	router.DELETE(deleteArticleURL, ah.Delete)
	router.PUT(updateArticleURL, ah.Update)
	router.PUT(updateArticleFileURL, ah.Middleware.updateArticleFile(ah.UpdateFile()))
	router.PUT(updateArticleImageURL, ah.UpdateImage)
	router.PUT(rateArticleUrl, ah.Rate)

	//router.GET(loadArticleImageURL, ah.LoadImage)
	//router.GET(loadArticleFileURL, ah.LoadFile)
}

func (ah *articleHandler) GetAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortingOptions := r.Context().Value(CtxKeySortAndFilters).(domain.SortFilterPagination)

		articles, pagesCount, err := ah.Service.GetAll(&sortingOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all articles. err: %v", err)})
			return
		}

		if pagesCount > 0 {
			w.Header().Set("pages", strconv.Itoa(pagesCount))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(articles)
	})
}

func (ah *articleHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	article, err := ah.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting article from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

func (ah *articleHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := r.Context().Value(CtxKeyCreateArticle).(map[string]interface{})

		createArticleDTO := domain.CreateArticleDTO{}

		createArticleDTO.Title = data["title"].(string)
		createArticleDTO.DirectionUUID = data["direction_uuid"].(string)
		createArticleDTO.AuthorUUID = data["author_uuid"].(string)
		createArticleDTO.Difficulty = data["difficulty"].(string)
		editionDate, err := strconv.Atoi(data["edition_date"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while creating article (converting edition date string to int). err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
			return
		}
		createArticleDTO.EditionDate = uint(editionDate)
		createArticleDTO.Description = data["description"].(string)
		createArticleDTO.Text = data["text"].(string)
		createArticleDTO.Language = data["language"].(string)
		createArticleDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
		fileName, ok := data["fileName"].(string)
		if ok {
			fileName = data["fileName"].(string)
			createArticleDTO.LocalURL = fmt.Sprintf("%s|split|/%s", loadArticleFileURL, fileName)
		} else {
			createArticleDTO.LocalURL = "file wasn't added"
		}
		createArticleDTO.WebURL = data["web_url"].(string)
		createArticleDTO.ImageURL = fmt.Sprintf("%s|split|/%s.jpg", loadArticleImageURL, string(utils.FormatOriginal))

		UUID, err := ah.Service.Create(&createArticleDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating article into DB. err: %v", err)})
			return
		}

		path := fmt.Sprintf("%s%s/", articleLocalStoragePath, UUID)

		file, ok := data["file"].(*bytes.Buffer)
		if ok && file != nil {
			err = ah.Service.SaveFile(path, fileName, file)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ah.logger.Errorf("error occurred while saving article into local store. err: %v.", err)
				json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
				return
			}
		}

		img := data["image"].(image.Image)
		err = ah.Service.SaveImage(path, &img)
		if err != nil {
			os.Remove(fmt.Sprintf("%s%s", path, fileName))
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while saving article image. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Article created successfully. UUID: %s", UUID)})
	})
}

func (ah *articleHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = ah.Service.Delete(uuid, fmt.Sprintf("%s%s/", articleLocalStoragePath, uuid))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting article from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Article with UUID %s was deleted", uuid)})
}

func (ah *articleHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateArticleDTO := &domain.UpdateArticleDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateArticleDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateArticleDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	// block changing URL
	updateArticleDTO.LocalURL = ""
	updateArticleDTO.ImageURL = ""

	err := ah.Service.Update(updateArticleDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating article into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Article updated successfully"})
}

func (ah *articleHandler) LoadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err := ah.Service.DownloadCountUp(uuid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)})
		return
	}

	path := fmt.Sprintf("%s%s/%s", articleLocalStoragePath, uuid, file)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	w.Header().Set("Content-Type", "application/octet-stream")

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

func (ah *articleHandler) UpdateFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := r.Context().Value(CtxKeyUpdateArticleFile).(domain.UpdateArticleFileDTO)

		dto.LocalPath = fmt.Sprintf("%s%s/", articleLocalStoragePath, dto.UUID)
		dto.LocalURL = fmt.Sprintf("%s?file=%s&uuid=%s", loadArticleFileURL, dto.NewFileName, dto.UUID)

		err := ah.Service.UpdateFile(&dto)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while saving article into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSON.Info{Msg: "File updated successfully"})
	})
}

func (ah *articleHandler) LoadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/jpeg")

	format := r.URL.Query().Get("format")
	if format == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("format query can't be empty")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s/", articleLocalStoragePath, uuid)

	img, err := utils.GetImageFromLocalStore(path, utils.Format(format), utils.JPG)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while saving image to local store. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving image to local store. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	jpeg.Encode(w, *img, nil)
}

func (ah *articleHandler) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s/", articleLocalStoragePath, uuid)

	img, err := jpeg.Decode(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while updating article image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
		return
	}

	err = ah.Service.SaveImage(path, &img)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while updating article image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
		return
	}

	json.NewEncoder(w).Encode(JSON.Info{Msg: "Article image updated successfully"})
}

func (ah *articleHandler) Rate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingStr := r.URL.Query().Get("rating")
	if ratingStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("rating query can't be empty")
		return
	}
	rating, err := strconv.ParseFloat(ratingStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("error occurred while parsing rating. Should be float32. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing rating. Should be float32. err: %v.", err)})
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("uuid can't be empty")
		return
	}

	err = ah.Service.Rate(uuid, float32(rating))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Errorf("error occurred while rating article image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while rating article into local store. err: %v.", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Article rated successfully. UUID: %s", uuid)})
}
