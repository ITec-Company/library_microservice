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
	"time"
)

const (
	getAllBooksURL     = "/books"
	getBookByUUIDURL   = "/book/:uuid"
	createBookURL      = "/book"
	deleteBookURL      = "/book/:uuid"
	updateBookURL      = "/book"
	loadBookFileURL    = "/file/book"
	updateBookFileURL  = "/file/book"
	loadBookImageURL   = "/image/book"
	updateBookImageURL = "/image/book"

	bookLocalStoragePath = "../store/books/"
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
	router.GET(loadBookFileURL, bh.LoadFile)
	router.PUT(updateBookFileURL, bh.Middleware.updateBookFile(bh.UpdateFile()))
	router.GET(loadBookImageURL, bh.LoadImage)
	router.PUT(updateBookImageURL, bh.UpdateImage)
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

	book, err := bh.Service.GetByUUID(uuid)
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

		createBookDTO.Title = data["title"].(string)
		createBookDTO.DirectionUUID = data["direction_uuid"].(string)
		createBookDTO.AuthorUUID = data["author_uuid"].(string)
		createBookDTO.Difficulty = data["difficulty"].(string)
		t, _ := time.Parse("2006-01-02", data["edition_date"].(string))
		createBookDTO.EditionDate = t
		createBookDTO.Description = data["description"].(string)
		createBookDTO.Language = data["language"].(string)
		createBookDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
		fileName := data["fileName"].(string)
		createBookDTO.LocalURL = fmt.Sprintf("%s?file=%s&uuid=", loadBookFileURL, fileName)
		createBookDTO.ImageURL = fmt.Sprintf("%s?format=%s&uuid=", loadBookImageURL, string(utils.FormatOriginal))

		UUID, err := bh.Service.Create(&createBookDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating book into DB. err: %v", err)})
			return
		}

		path := fmt.Sprintf("%s%s/", bookLocalStoragePath, UUID)

		file := data["file"].(*bytes.Buffer)
		err = bh.Service.SaveFile(path, fileName, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			bh.logger.Errorf("error occurred while saving book into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local store. err: %v.", err)})
			return
		}

		img := data["image"].(image.Image)
		err = bh.Service.SaveImage(path, &img)
		if err != nil {
			os.Remove(fmt.Sprintf("%s%s", path, fileName))
			w.WriteHeader(http.StatusInternalServerError)
			bh.logger.Errorf("error occurred while updating book image. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local store. err: %v.", err)})
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

	err = bh.Service.Delete(uuid, fmt.Sprintf("%s%s/", bookLocalStoragePath, uuid))
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

	// block changing URL
	updateBookDTO.LocalURL = ""
	updateBookDTO.ImageURL = ""

	err := bh.Service.Update(updateBookDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating book into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Book updated successfully"})
}

func (bh *bookHandler) LoadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file := r.URL.Query().Get("file")
	if file == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("file query can't be empty")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid query can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s/%s", bookLocalStoragePath, uuid, file)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))
	w.Header().Set("Content-Type", "application/octet-stream")

	fileBytes, err := bh.Service.LoadFile(path)
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

func (bh *bookHandler) UpdateFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := r.Context().Value(CtxKeyUpdateBookFile).(domain.UpdateBookFileDTO)

		dto.LocalPath = fmt.Sprintf("%s%s/", bookLocalStoragePath, dto.UUID)
		dto.LocalURL = fmt.Sprintf("%s?file=%s&uuid=%s", loadBookFileURL, dto.NewFileName, dto.UUID)

		err := bh.Service.UpdateFile(&dto)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			bh.logger.Errorf("error occurred while saving book into local store. err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local store. err: %v.", err)})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JSON.Info{Msg: "File updated successfully"})
	})
}

func (bh *bookHandler) LoadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/jpeg")

	format := r.URL.Query().Get("format")
	if format == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("format query can't be empty")
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s/", bookLocalStoragePath, uuid)

	img, err := utils.GetImageFromLocalStore(path, utils.Format(format), utils.JPG)
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

func (bh *bookHandler) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid can't be empty")
		return
	}

	path := fmt.Sprintf("%s%s/", bookLocalStoragePath, uuid)

	img, err := jpeg.Decode(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while updating book image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local store. err: %v.", err)})
		return
	}

	err = bh.Service.SaveImage(path, &img)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while updating book image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving book into local store. err: %v.", err)})
		return
	}

	json.NewEncoder(w).Encode(JSON.Info{Msg: "Book image updated successfully"})
}
