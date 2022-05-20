package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/julienschmidt/httprouter"
	"library-go/internal/domain"
	"library-go/pkg/JSON"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"net/http"
	"strconv"
	"strings"
)

type CtxKey int8

const (
	defaultLimit = 20

	CtxKeyCreateArticle  CtxKey = 1
	CtxKeyCreateBook     CtxKey = 2
	CtxKeyCreateVideo    CtxKey = 3
	CtxKeyCreateAudio    CtxKey = 4
	CtxKeySortAndFilters CtxKey = 5
)

type Middleware struct {
	logger *logging.Logger
}

func NewMiddlewares(logger *logging.Logger) Middleware {
	return Middleware{
		logger: logger,
	}
}

func (m *Middleware) createArticle(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data := map[string]interface{}{
			"file":           "file",
			"title":          "text",
			"direction_uuid": "text",
			"author_uuid":    "text",
			"difficulty":     "text",
			"edition_date":   "text",
			"description":    "text",
			"language":       "text",
			"tags_uuids":     "text",
		}

		err := utils.ParseMultiPartFormData(r, data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
			return
		}

		file := data["file"].(*bytes.Buffer)

		if !filetype.IsDocument(file.Bytes()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("file is not a document, allowed extensions: Doc, Docx, Xls, Xlsx, Ppt, Pptx")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "file is not a document, allowed extensions: Doc, Docx, Xls, Xlsx, Ppt, Pptx"})
			return
		}

		kind, err := filetype.Match(file.Bytes())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while ckecking file extension. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while ckecking file extension. err: %v", err)})
			return
		}

		data["title"] = strings.Replace(data["title"].(string), " ", "_", -1)
		data["fileName"] = fmt.Sprintf("author(%s)-title(%s).%s", data["author_uuid"].(string), data["title"].(string), kind.Extension)

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), CtxKeyCreateArticle, data)))
	}
}

func (m *Middleware) createBook(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		data := map[string]interface{}{
			"file":           "file",
			"title":          "text",
			"direction_uuid": "text",
			"author_uuid":    "text",
			"difficulty":     "text",
			"edition_date":   "text",
			"description":    "text",
			"language":       "text",
			"tags_uuids":     "text",
		}

		err := utils.ParseMultiPartFormData(r, data)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
			return
		}

		file := data["file"].(*bytes.Buffer)

		if !filetype.IsArchive(file.Bytes()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("file is not a pdf, allowed extensions: pdf")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "file is not a pdf, allowed extensions: pdf"})
			return
		}

		kind, err := filetype.Match(file.Bytes())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while ckecking file extension. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while ckecking file extension. err: %v", err)})
			return
		}

		if kind.Extension != "pdf" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("file is not a document, allowed extensions: pdf")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "file is not a document, allowed extensions: pdf"})
			return
		}

		data["title"] = strings.Replace(data["title"].(string), " ", "_", -1)
		data["fileName"] = fmt.Sprintf("author(%s)-title(%s).%s", data["author_uuid"].(string), data["title"].(string), kind.Extension)

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), CtxKeyCreateBook, data)))
	}
}

func (m *Middleware) createVideo(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
			return
		}

		file := data["file"].(*bytes.Buffer)

		if !filetype.IsVideo(file.Bytes()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("file is not a video, allowed extensions: mp4, m4v, mkv, webm, mov, avi, wmv, mpg, flv, 3gp")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "file is not a video, allowed extensions: mp4, m4v, mkv, webm, mov, avi, wmv, mpg, flv, 3gp"})
			return
		}

		kind, err := filetype.Match(file.Bytes())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while ckecking file extension. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while ckecking file extension. err: %v", err)})
			return
		}

		data["title"] = strings.Replace(data["title"].(string), " ", "_", -1)
		data["fileName"] = fmt.Sprintf("title(%s).%s", data["title"].(string), kind.Extension)

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), CtxKeyCreateVideo, data)))
	}
}

func (m *Middleware) createAudio(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
			return
		}

		file := data["file"].(*bytes.Buffer)
		if file != nil {
			m.logger.Errorf("file is not nil %v", file)
		}
		if !filetype.IsAudio(file.Bytes()) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("file is not an audio, allowed extensions: mid, mp3, m4a, ogg, flac, wav, amr, aac, aiff")
			json.NewEncoder(w).Encode(JSON.Error{Msg: "file is not an audio, allowed extensions: mid, mp3, m4a, ogg, flac, wav, amr, aac, aiff"})
			return
		}

		kind, err := filetype.Match(file.Bytes())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			m.logger.Errorf("error occurred while ckecking file extension. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while ckecking file extension. err: %v", err)})
			return
		}

		data["title"] = strings.Replace(data["title"].(string), " ", "_", -1)
		data["fileName"] = fmt.Sprintf("title(%s).%s", data["title"].(string), kind.Extension)

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), CtxKeyCreateAudio, data)))
	}
}

func (m *Middleware) sortAndFilters(next http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var sortAndFilters domain.SortFilterPagination

		// Sorting

		if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
			if err := utils.IsSQL(sortBy); err != nil {
				m.logger.Errorf("sql injectin detected. Assign nil to value sort_by: %s", sortBy)
				sortBy = ""
			} else {
				sortAndFilters.SortBy = sortBy

				if so := r.URL.Query().Get("sort_order"); so != "" {
					if err := utils.IsSQL(so); err != nil {
						m.logger.Errorf("sql injectin detected. Assign nil to value sort_order: %s", so)
						so = ""
					} else {
						sortOrder := domain.Order(strings.ToLower(so))
						if sortOrder == domain.OrderDESC || sortOrder == domain.OrderASC {
							sortAndFilters.Order = sortOrder
						} else {
							m.logger.Errorf("sort oeder must be only 'asc' or 'desc'")
						}
					}

				}
			}
		}

		// Filters filer:1,2,3,|filter2:2,4,5
		if filters := r.URL.Query().Get("filter"); filters != "" {
			if err := utils.IsSQL(filters); err != nil {
				m.logger.Errorf("sql injectin detected. Assign nil to value filters: %s", filters)
				filters = ""
			} else {
				filtersMap := make(map[string]interface{})
				filtersSplit := strings.Split(filters, "|")
				for _, s := range filtersSplit {
					filterAndArgs := strings.Split(s, ":")
					filter := filterAndArgs[0]
					args := strings.Split(filterAndArgs[1], ",")
					filtersMap[filter] = args
				}
				sortAndFilters.FiltersAndArgs = filtersMap
			}
		}

		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if err := utils.IsSQL(limitStr); err != nil {
				m.logger.Errorf("sql injectin detected. Assign defailt value to limit: %s", limitStr)
				sortAndFilters.Limit = defaultLimit
			} else {
				limit, err := strconv.ParseUint(limitStr, 10, 64)
				if err != nil {
					sortAndFilters.Limit = defaultLimit
					m.logger.Errorf("limit was assigned with default value %d cause error occured while converting input value to string. err: %v", defaultLimit, err)
				} else {
					sortAndFilters.Limit = limit
				}
				if pageStr := r.URL.Query().Get("page"); pageStr != "" {
					if err := utils.IsSQL(pageStr); err != nil {
						m.logger.Errorf("sql injectin detected. Assign 1 to value page: %s", pageStr)
						sortAndFilters.Page = 1
					} else {
						page, err := strconv.ParseUint(pageStr, 10, 64)
						if err != nil {
							sortAndFilters.Page = 1
							m.logger.Errorf("page was assigned with value 0 cause error occured while converting input value to string. err: %v", err)
						} else {
							sortAndFilters.Page = page
						}
					}
				}
			}

		}

		r = r.WithContext(context.WithValue(r.Context(), CtxKeySortAndFilters, sortAndFilters))
		next(w, r)
	}
}
