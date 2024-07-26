package message

import (
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MessageService struct {
	repo Repository
	log  *logrus.Logger
}

type Repository interface {
	Create(m models.CreateRequest, status models.Status) (models.Message, error)
	Get(m models.Message, limit uint64, offset uint64) (models.FilterResponse, error)
	Update(id uuid.UUID, status models.Status) (models.Message, error)
}

func NewService(repo Repository, logger *logrus.Logger) *MessageService {
	return &MessageService{repo: repo, log: logger}
}

func (s *MessageService) CreateMessage(message models.CreateRequest) (models.Message, error) {
	s.log.Debugf("Write message to repository")
	resp, err := s.repo.Create(message, models.StatusProcessing)
	if err != nil {
		s.log.Error(err)
		return models.Message{}, err
	}

	return resp, nil
}

func (s *MessageService) GetMessages(req models.FilterRequest) (models.FilterResponse, error) {
	s.log.Debugf("Read messages from repository")
	resp, err := s.repo.Get(req.Fields, req.Limit, req.Offset)
	if err != nil {
		s.log.Error(err)
		return models.FilterResponse{}, err
	}
	return resp, nil
}
