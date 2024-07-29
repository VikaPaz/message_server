package queue

import (
	"context"
	"encoding/json"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type QueueReader struct {
	service service

	reader *kafka.Reader
	logger *logrus.Logger
}

type service interface {
	UpdateMassage(request models.UpdateRequest) error
}

func (q *QueueReader) Listen() {
	for {
		msg, err := q.reader.ReadMessage(context.Background())
		if err != nil {
			q.logger.Errorf("Error reading message: %s\n", err)
			continue
		}

		req := models.UpdateRequest{}
		err = json.Unmarshal(msg.Value, &req)
		if err != nil {
			q.logger.Errorf("Error unmarshalling message: %s\n", err)
			continue
		}
		err = q.service.UpdateMassage(req)
		if err != nil {
			q.logger.Errorf("Error updating massage: %s\n", err)
			continue
		}
	}
}

func NewReader(cfg kafka.ReaderConfig, log *logrus.Logger, svc service) *QueueReader {
	return &QueueReader{
		service: svc,
		reader:  kafka.NewReader(cfg),
		logger:  log,
	}
}
