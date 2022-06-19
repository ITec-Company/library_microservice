package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"image"
	"image/jpeg"
	"library-go/internal/domain"
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
	getAllBooksURL     = "/books"
	getBookByUUIDURL   = "/book/:uuid"
	createBookURL      = "/book"
	deleteBookURL      = "/book/:uuid"
	updateBookURL      = "/book"
	updateBookFileURL  = "/file/book"
	updateBookImageURL = "/image/book"
	rateBookUrl        = "/rate/book"

	loadBookImageURL = "/books/"
	loadBookFileURL  = "/books/"

	bookLocalStoragePath = "../store/books/"
)

type BookHandler struct {
	Service service.BookService
	logger  *logging.Logger
}

func NewBookHandler(service service.BookService, logger *logging.Logger) BookHandler {
	return BookHandler{
		Service: service,
		logger:  logger,
	}
}

func (bh *BookHandler) Register(router *httprouter.Router, middleware *Middleware) {
	router.GET(getAllBooksURL, middleware.sortAndFilters(bh.GetAll()))
	router.GET(getBookByUUIDURL, bh.GetByUUID)
	router.POST(createBookURL, middleware.createBook(bh.Create()))
	router.DELETE(deleteBookURL, bh.Delete)
	router.PUT(updateBookURL, bh.Update)
	router.PUT(updateBookFileURL, middleware.updateBookFile(bh.UpdateFile()))
	router.PUT(updateBookImageURL, bh.UpdateImage)
	router.PUT(rateBookUrl, bh.Rate)

	//router.GET(loadBookImageURL, bh.LoadImage)
	//router.GET(loadBookFileURL, bh.LoadFile)
}

func (bh *BookHandler) GetAll() http.HandlerFunc {
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

func (bh *BookHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (bh *BookHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := r.Context().Value(CtxKeyCreateBook).(map[string]interface{})

		createBookDTO := domain.CreateBookDTO{}

		createBookDTO.Title = data["title"].(string)
		createBookDTO.DirectionUUID = data["direction_uuid"].(string)
		createBookDTO.AuthorUUID = data["author_uuid"].(string)
		createBookDTO.Difficulty = data["difficulty"].(string)
		editionDate, err := strconv.Atoi(data["edition_date"].(string))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			bh.logger.Errorf("error occurred while creating book (converting edition date string to int). err: %v.", err)
			json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while saving article into local store. err: %v.", err)})
			return
		}
		createBookDTO.EditionDate = uint(editionDate)
		createBookDTO.Description = data["description"].(string)
		createBookDTO.Language = data["language"].(string)
		createBookDTO.TagsUUIDs = strings.Split(data["tags_uuids"].(string), ",")
		fileName, ok := data["fileName"].(string)
		if ok {
			fileName = data["fileName"].(string)
			createBookDTO.LocalURL = fmt.Sprintf("%s|split|/%s", loadBookFileURL, fileName)
		} else {
			createBookDTO.LocalURL = "file wasn't added"
		}

		createBookDTO.ImageURL = fmt.Sprintf("%s|split|/%s.jpg", loadBookImageURL, string(utils.FormatOriginal))

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

func (bh *BookHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (bh *BookHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (bh *BookHandler) LoadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err := bh.Service.DownloadCountUp(uuid)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while downloading (DB ping to increase dowload count). err: %v", err)})
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

func (bh *BookHandler) UpdateFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dto := r.Context().Value(CtxKeyUpdateBookFile).(domain.UpdateBookFileDTO)

		dto.LocalPath = fmt.Sprintf("%s%s/", bookLocalStoragePath, dto.UUID)
		dto.LocalURL = fmt.Sprintf("%s|split|/%s", loadBookFileURL, dto.NewFileName)

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

func (bh *BookHandler) LoadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (bh *BookHandler) UpdateImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (bh *BookHandler) Rate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingStr := r.URL.Query().Get("rating")
	if ratingStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("rating query can't be empty")
		return
	}
	rating, err := strconv.ParseFloat(ratingStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("error occurred while parsing rating. Should be float32. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing rating. Should be float32. err: %v.", err)})
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		bh.logger.Errorf("uuid can't be empty")
		return
	}

	err = bh.Service.Rate(uuid, float32(rating))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		bh.logger.Errorf("error occurred while rating book image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while rating book into local store. err: %v.", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Book rated successfully. UUID: %s", uuid)})
}
