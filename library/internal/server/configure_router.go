package server

import (
	"library/domain/store"
	"library/internal/handlers/book"
)

// ConfigureRouter ...
func (s *Server) ConfigureRouter() {
	s.Router.Handle("POST", "/save", book.SaveBookHandle(store.New(s.Config)))

}
