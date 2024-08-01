package queue

import (
	"encoding/json"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type MessageQueue struct {
	conn *kafka.Conn
	log  *logrus.Logger
}

func NewQueue(conn *kafka.Conn, log *logrus.Logger) *MessageQueue {
	return &MessageQueue{
		conn: conn,
		log:  log,
	}
}

func (q *MessageQueue) Write(m models.WriteRequest) error {
	value, err := json.Marshal(m)
	if err != nil {
		q.log.Error("failed to serialize structure: %v", err)
		return err
	}

	q.log.Debugf("Writing message to kafka: %v", string(value))
	_, err = q.conn.WriteMessages(
		kafka.Message{Value: value},
	)
	if err != nil {
		q.log.Errorf("failed to write messages: %v", err)
		return err
	}
	return nil
}
