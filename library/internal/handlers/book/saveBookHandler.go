package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"library/domain/model"
	"library/domain/store"
	"library/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// SaveBookHandle ...
func SaveBookHandle(s *store.Store) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		book := &model.Book{
			ID:            1,
			SubDirection:  *model.TestDevSubDirection(),
			Title:         "Golang_for_beginners",
			EditionDate:   time.Date(2015, 05, 05, 0, 0, 0, 0, time.UTC),
			Diffuculty:    model.JuniorLevel,
			Rating:        8.2,
			Description:   "Description",
			Language:      "eng",
			URL:           "URL...",
			DownloadCount: 829,
		}

		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Errorf("error occured while getting file. Err msg: %v.", err)
			json.NewEncoder(w).Encode(response.Error{Messsage: fmt.Sprintf("error occured while getting file. Err msg: %v.", err)})
			return
		}
		author := model.TestAuthor()
		bookInfo := strings.Split(header.Filename, "#")
		book.Title = bookInfo[1]

		author.FullName = bookInfo[0]
		book.Author = *author

		defer file.Close()
		_, err = io.Copy(&buf, file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			s.Logger.Errorf("error occured while copy file. Err msg: %v.", err)
			json.NewEncoder(w).Encode(response.Error{Messsage: fmt.Sprintf("error occured while coping file. Err msg: %v.", err)})
			return
		}

		err = s.Open()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response.Error{Messsage: fmt.Sprintf("Error occured while openong DB. Err msg: %v", err)})
			return
		}

		s.Book().SaveBook(book, &buf)
		buf.Reset()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(book)
	}
}
