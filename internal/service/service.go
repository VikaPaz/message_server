package message

import (
	"errors"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MessageService struct {
	repo Repository
	que  Queue
	log  *logrus.Logger
}

type Repository interface {
	Create(m models.CreateRequest, status models.Status) (models.Message, error)
	Get(m models.Message, limit uint64, offset uint64) (models.FilterResponse, error)
	Update(id uuid.UUID, status models.Status) (models.Message, error)
}

func (s *MessageService) UpdateMassage(request models.UpdateRequest) error {
	resp, err := s.repo.Update(request.ID, *request.Status)
	if err != nil {
		return err
	}
	s.log.Infof("updated message id: %s; status: %s", *resp.ID, *resp.Status)
	return nil
}

type Queue interface {
	Write(m models.WriteRequest) error
}

func NewService(repo Repository, que Queue, logger *logrus.Logger) *MessageService {
	return &MessageService{
		repo: repo,
		que:  que,
		log:  logger,
	}
}

func (s *MessageService) CreateMessage(message models.CreateRequest) (models.Message, error) {
	s.log.Debugf("Write message to repository")
	resp, err := s.repo.Create(message, models.StatusProcessing)
	if err != nil {
		return models.Message{}, errors.Join(models.ErrRequestDBFailed, err)
	}

	s.log.Debugf("Write message to queue")
	w := models.WriteRequest{
		ID:      *resp.ID,
		Message: resp.Message,
	}
	err = s.que.Write(w)
	if err != nil {
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
