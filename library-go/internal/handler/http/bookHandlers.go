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
	"time"
)

const (
	getAllBooksURL       = "/books"
	getBookByUUIDURL     = "/book/:uuid"
	createBookURL        = "/book"
	deleteBookURL        = "/book/:uuid"
	updateBookURL        = "/book"
	loadBookURL          = "/load/book/:url"
	bookLocalStoragePath = "../store/books"
)

type bookHandler struct {
	Service service.BookService
	logger  *logging.Logger
}

func NewBookHandler(service service.BookService, logger *logging.Logger) handler.Handler {
	return &bookHandler{
		Service: service,
		logger:  logger,
	}
}

func (bh *bookHandler) Register(router *httprouter.Router) {
	router.GET(getAllBooksURL, bh.GetAll)
	router.GET(getBookByUUIDURL, bh.GetByUUID)
	router.POST(createBookURL, bh.Create)
	router.DELETE(deleteBookURL, bh.Delete)
	router.PUT(updateBookURL, bh.Update)
	router.GET(loadBookURL, bh.Load)
}

func (bh *bookHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		bh.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		bh.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	books, err := bh.Service.GetAll(context.Background(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all books. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func (bh *bookHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid can't be empty")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid can't be empty"})
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	book, err := bh.Service.GetByUUID(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting book from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (bh *bookHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createBookDTO := domain.CreateBookDTO{}

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
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("error occurred while parsing multiform data. err msg: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing multiform data. err msg: %v.", err)})
		return
	}

	t, _ := time.Parse("2006-01-02", data["edition_date"].(string))

	file := data["file"].(*bytes.Buffer)
	createBookDTO.Title = data["title"].(string)
	createBookDTO.DirectionUUID = data["direction_uuid"].(string)
	createBookDTO.AuthorUUID = data["author_uuid"].(string)
	createBookDTO.Difficulty = data["difficulty"].(string)
	createBookDTO.EditionDate = t
	createBookDTO.Description = data["description"].(string)
	createBookDTO.Language = data["language"].(string)
	createBookDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
	createBookDTO.URL = fmt.Sprintf("%s%s-%s-%s", "load/book/", createBookDTO.DirectionUUID, createBookDTO.Difficulty, data["fileName"].(string))

	path := fmt.Sprintf("%s/%s/%s", bookLocalStoragePath, createBookDTO.DirectionUUID, createBookDTO.Difficulty)

	err = utils.SaveFile(path, data["fileName"].(string), file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while saving book into local database. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local database. err: %v", err)})
		return
	}

	UUID, err := bh.Service.Create(context.Background(), &createBookDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating book into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Book created successfully. UUID: %s", UUID)})
}

func (bh *bookHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = bh.Service.Delete(context.Background(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting book from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Book with UUID %s was deleted", uuid)})
}

func (bh *bookHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateBookDTO := &domain.UpdateBookDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateBookDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateBookDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := bh.Service.Update(context.Background(), updateBookDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating book into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Book updated successfully"})
}

func (bh *bookHandler) Load(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	url := ps.ByName("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("url can't be empty")
		return
	}
	urlArray := strings.Split(url, "-")

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", urlArray[len(urlArray)-1]))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := strings.Replace(url, "-", "/", -1)

	fileBytes, err := bh.Service.LoadLocalFIle(context.Background(), fmt.Sprintf("%s/%s", bookLocalStoragePath, path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)

}
