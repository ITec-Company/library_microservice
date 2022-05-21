package http

import (
	"bytes"
	"context"
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
	"time"
)

const (
	getAllBooksURL       = "/books"
	getBookByUUIDURL     = "/book/:uuid"
	createBookURL        = "/book"
	deleteBookURL        = "/book/:uuid"
	updateBookURL        = "/book"
	loadBookURL          = "/load/book"
	loadBookImageURL     = "/load/image/book"
	bookLocalStoragePath = "../store/books/"
	bookImageStoragePath = "../store/images/books/"
)

type bookHandler struct {
	Service    service.BookService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewBookHandler(service service.BookService, logger *logging.Logger, middleware *Middleware) handler.Handler {
	return &bookHandler{
		Service:    service,
		logger:     logger,
		Middleware: middleware,
	}
}

func (bh *bookHandler) Register(router *httprouter.Router) {
	router.GET(getAllBooksURL, bh.Middleware.sortAndFilters(bh.GetAll()))
	router.GET(getBookByUUIDURL, bh.GetByUUID)
	router.POST(createBookURL, bh.Middleware.createBook(bh.Create()))
	router.DELETE(deleteBookURL, bh.Delete)
	router.PUT(updateBookURL, bh.Update)
	router.GET(loadBookURL, bh.Load)
	router.GET(loadBookImageURL, bh.LoadImage)
}

func (bh *bookHandler) GetAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortingOptions := r.Context().Value(CtxKeySortAndFilters).(domain.SortFilterPagination)

		books, pagesCount, err := bh.Service.GetAll(&sortingOptions)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all books. err: %v", err)})
			return
		}

		if pagesCount > 0 {
			w.Header().Set("pages", strconv.Itoa(pagesCount))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(books)
	})
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

func (bh *bookHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := r.Context().Value(CtxKeyCreateBook).(map[string]interface{})

		createBookDTO := domain.CreateBookDTO{}

		t, _ := time.Parse("2006-01-02", data["edition_date"].(string))
		createAt := time.Now().UTC()

		file := data["file"].(*bytes.Buffer)

		createBookDTO.Title = data["title"].(string)
		createBookDTO.DirectionUUID = data["direction_uuid"].(string)
		createBookDTO.AuthorUUID = data["author_uuid"].(string)
		createBookDTO.Difficulty = data["difficulty"].(string)
		createBookDTO.EditionDate = t
		createBookDTO.Description = data["description"].(string)
		createBookDTO.Language = data["language"].(string)
		createBookDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")

		createBookDTO.CreatedAt = createAt
		createdStr := createAt.Format("2006-01-02-15-04-05")

		img := data["image"].(image.Image)
		imagePath := fmt.Sprintf("%s/%s", bookImageStoragePath, createdStr)
		imageFileName, err := bh.Service.SaveImage(context.Background(), imagePath, &img)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			bh.logger.Errorf("error occurred while saving book image into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book image into local store. err: %v.", err)})
			return
		}
		bh.logger.Errorf("%v", createAt)

		createBookDTO.ImageURL = fmt.Sprintf("%s?url=%s%%2F%s", loadBookImageURL, createdStr, imageFileName)

		fileName := data["fileName"].(string)
		createBookDTO.URL = fmt.Sprintf("%s?url=%s", loadBookURL, fileName)

		err = bh.Service.Save(context.Background(), bookLocalStoragePath, fileName, file)
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
	})
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
	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("url can't be empty")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url))
	w.Header().Set("Content-Type", "application/octet-stream")

	path := fmt.Sprintf("%s%s", bookLocalStoragePath, url)

	fileBytes, err := bh.Service.Load(context.Background(), path)
	_, pathError := err.(*os.PathError)
	if pathError {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while searching file: invalid path. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while searching file: invalid path. err: %v", err)})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while reading file. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while reading file. err: %v", err)})
		return
	}

	w.Write(fileBytes)

}

func (bh *bookHandler) LoadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/jpeg")

	url := r.URL.Query().Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("url can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s", bookImageStoragePath, url)

	img, err := utils.GetImageFromLocalStore(path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while saving image to local store: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving image to local store. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	jpeg.Encode(w, *img, nil)

}
