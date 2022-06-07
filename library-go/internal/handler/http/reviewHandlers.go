package http

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"library-go/internal/domain"
	"library-go/internal/handler"
	"library-go/internal/service"
	"library-go/pkg/JSON"
	"library-go/pkg/logging"
	"net/http"
	"strconv"
)

const (
	getAllReviewsURL   = "/reviews"
	getReviewByUUIDURL = "/review/:uuid"
	createReviewURL    = "/review"
	deleteReviewURL    = "/review/:uuid"
	updateReviewURL    = "/review"
	rateReviewUrl      = "/rate/review"
)

type reviewHandler struct {
	Service    service.ReviewService
	logger     *logging.Logger
	Middleware *Middleware
}

func NewReviewHandler(service service.ReviewService, logger *logging.Logger, middleware *Middleware) handler.Handler {
	return &reviewHandler{
		Service:    service,
		logger:     logger,
		Middleware: middleware,
	}
}

func (rh *reviewHandler) Register(router *httprouter.Router) {
	router.GET(getAllReviewsURL, rh.GetAll)
	router.GET(getReviewByUUIDURL, rh.GetByUUID)
	router.POST(createReviewURL, rh.Create)
	router.DELETE(deleteReviewURL, rh.Delete)
	router.PUT(updateReviewURL, rh.Update)
	router.PUT(rateReviewUrl, rh.Rate)
}

func (rh *reviewHandler) GetAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 0 {
		limit = 0
		rh.logger.Debugf("error occurred while parsing limit. err: %v. Assigning '0' to limit", err)
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
		rh.logger.Debugf("error occurred while parsing offset. err: %v. Assigning '0' to offset", err)
	}

	reviews, err := rh.Service.GetAll(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting all reviews. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}

func (rh *reviewHandler) GetByUUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	review, err := rh.Service.GetByUUID(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while getting review from DB by UUID. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

func (rh *reviewHandler) Create(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	createReviewDTO := &domain.CreateReviewDTO{}
	if err := json.NewDecoder(r.Body).Decode(createReviewDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}

	UUID, err := rh.Service.Create(createReviewDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while creating review into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Review created successfully. UUID: %s", UUID)})
}

func (rh *reviewHandler) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	uuid := ps.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("uuid can't be empty")
		return
	}
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil || uuidInt < 0 {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("uuid need to be uint")
		json.NewEncoder(w).Encode(JSON.Error{Msg: "uuid need to be uint"})
		return
	}

	err = rh.Service.Delete(uuid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while deleting review from DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Review with UUID %s was deleted", uuid)})
}

func (rh *reviewHandler) Update(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	updateReviewDTO := &domain.UpdateReviewDTO{}
	if err := json.NewDecoder(r.Body).Decode(updateReviewDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("error occurred while decoding JSON request. err: %v", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: %v", err)})
		return
	}
	if updateReviewDTO.UUID == "" {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("error occurred while decoding JSON request. err: nil UUID")
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while decoding JSON request. err: nil UUID")})
		return
	}

	err := rh.Service.Update(updateReviewDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while updating review into DB. err: %v", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: "Review updated successfully"})
}

func (rh *reviewHandler) Rate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	ratingStr := r.URL.Query().Get("rating")
	if ratingStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("rating query can't be empty")
		return
	}
	rating, err := strconv.ParseFloat(ratingStr, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("error occurred while parsing rating. Should be float32. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while parsing rating. Should be float32. err: %v.", err)})
		return
	}

	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		rh.logger.Errorf("uuid can't be empty")
		return
	}

	err = rh.Service.Rate(uuid, float32(rating))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rh.logger.Errorf("error occurred while rating review image. err: %v.", err)
		json.NewEncoder(w).Encode(JSON.Error{Msg: fmt.Sprintf("error occurred while rating review into local store. err: %v.", err)})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JSON.Info{Msg: fmt.Sprintf("Review rated successfully. UUID: %s", uuid)})
}
