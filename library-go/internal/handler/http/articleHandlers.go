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
	"time"
)

const (
	getAllArticlesURL       = "/articles"
	getArticleByUUIDURL     = "/article/:uuid"
	createArticleURL        = "/article"
	deleteArticleURL        = "/article/:uuid"
	updateArticleURL        = "/article"
	loadArticleURL          = "/load/article"
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
	router.GET(loadArticleURL, ah.Load)
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

	article, err := ah.Service.GetByUUID(context.Background(), uuid)
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

		t, _ := time.Parse("2006-01-02", data["edition_date"].(string))

		file := data["file"].(*bytes.Buffer)

		createArticleDTO.Title = data["title"].(string)
		createArticleDTO.DirectionUUID = data["direction_uuid"].(string)
		createArticleDTO.AuthorUUID = data["author_uuid"].(string)
		createArticleDTO.Difficulty = data["difficulty"].(string)
		createArticleDTO.EditionDate = t
		createArticleDTO.Description = data["description"].(string)
		createArticleDTO.Language = data["language"].(string)
		createArticleDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")

		fileName := data["fileName"].(string)
		createArticleDTO.URL = fmt.Sprintf("%s?url=%s", loadArticleURL, fileName)

		err := ah.Service.Save(context.Background(), articleLocalStoragePath, fileName, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ah.logger.Errorf("error occurred while saving article into local database. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local database. err: %v", err)})
			return
		}

		UUID, err := ah.Service.Create(context.Background(), &createArticleDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating article into DB. err: %v", err)})
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

	err = ah.Service.Delete(context.Background(), uuid)
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

	err := ah.Service.Update(context.Background(), updateArticleDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating article into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Article updated successfully"})
}

func (ah *articleHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Errorf("url can't be empty")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := fmt.Sprintf("%s%s", articleLocalStoragePath, url)

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
