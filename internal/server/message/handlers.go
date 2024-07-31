package message

import (
	"encoding/json"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	CreateMessage(message models.CreateRequest) (models.Message, error)
	GetMessages(req models.FilterRequest) (models.FilterResponse, error)
}

type Handler struct {
	service Service
	log     *logrus.Logger
}

func NewHandler(service Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		log:     logger,
	}
}

func (rs *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/new", rs.new)
	r.Get("/get", rs.get)

	return r
}

// @Summary Creating a new message
// @Description Handles request to create a new message and returns the message information in JSON.
// @Tags messages
// @Accept json
// @Produce json
// @Param request body models.CreateRequest true "Passport"
// @Success 200 {object} models.Message "Created message"
// @Failure 400
// @Failure 500
// @Router /message/new [post]
func (rs *Handler) new(w http.ResponseWriter, r *http.Request) {
	msg := models.CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if msg.Message == "" {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := rs.service.CreateMessage(msg)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// @Summary Get messages
// @Description Handles request to get messages by filter.
// @Tags messages
// @Produce  json
// @Param id query string false "Message ID"
// @Param message query string false "Message content"
// @Param status query string false "Status"
// @Param created_at query string false "Creation timestamp (DateTime format)"
// @Param updated_at query string false "Update timestamp (DateTime format)"
// @Param limit query uint64 false "Limit"
// @Param offset query uint64 false "Offset"
// @Success 200 {object} models.FilterResponse  "Successfully got messages"
// @Failure 400
// @Failure 500
// @Router /message/get [get]
func (rs *Handler) get(w http.ResponseWriter, r *http.Request) {
	filter := models.FilterRequest{}
	params := r.URL.Query()
	if idStr := params.Get("id"); idStr != "" {
		id, err := uuid.Parse(idStr)
		if err != nil {
			rs.log.Error(w, "Invalid UUID for id", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Fields.ID = &id
	}
	if message := params.Get("message"); message != "" {
		filter.Fields.Message = &message
	}
	if status := params.Get("status"); status != "" {
		s := models.Status(status)
		filter.Fields.Status = &s
	}
	if createdAtStr := params.Get("created_at"); createdAtStr != "" {
		createdAt, err := time.Parse(time.DateTime, createdAtStr)
		if err != nil {
			rs.log.Error(w, "Invalid start_time parameter (RFC3339 format expected)", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Fields.CreatedAt = &createdAt
	}
	if updatedAtStr := params.Get("updated_at"); updatedAtStr != "" {
		updatedAt, err := time.Parse(time.DateTime, updatedAtStr)
		if err != nil {
			rs.log.Error(w, "Invalid updated_at parameter (RFC3339 format expected)", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Fields.CreatedAt = &updatedAt
	}
	if l := params.Get("limit"); l != "" {
		limit, err := strconv.ParseUint(l, 10, 64)
		if err != nil {
			rs.log.Error(w, "Invalid limit parameter (uint64 format expected)", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Limit = limit
	}
	if l := params.Get("offset"); l != "" {
		offset, err := strconv.ParseUint(l, 10, 64)
		if err != nil {
			rs.log.Error(w, "Invalid offset parameter (uint64 format expected)", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filter.Offset = offset
	}

	resp, err := rs.service.GetMessages(filter)
	if err != nil {
		rs.log.Error(w, "Error getting messages", err, http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		rs.log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
