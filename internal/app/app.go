package app

import (
	"github.com/VikaPaz/message_server/internal/client/queue"
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/VikaPaz/message_server/internal/repository"
	"github.com/VikaPaz/message_server/internal/server"
	queue2 "github.com/VikaPaz/message_server/internal/server/queue"
	messageService "github.com/VikaPaz/message_server/internal/service"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Run(logger *logrus.Logger) error {
	logger.Infoln("Starting server...")
	logger.Debugf("Loading environment variables")
	if err := godotenv.Overload("env/.env"); err != nil {
		logger.Errorf("Error loading .env file")
		return models.ErrLoadEnvFailed
	}

	confPostgres := message.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	logger.Infof("Connecting to PostgreSQL...")
	dbConn, err := message.Connection(confPostgres)
	if err != nil {
		logger.Errorf("Error connecting to database")
		return err
	}

	repo := message.NewRepository(dbConn, logger)

	logger.Infof("Connecting to queue...")
	confKafka := queue.Config{
		Topic:     "topic1",
		Partition: 0,
		Host:      os.Getenv("KAFKA_ADDRESS"),
		Network:   "tcp",
	}
	kafkaConn, err := queue.Connection(confKafka)
	if err != nil {
		logger.Fatalf("Error connecting to kafka: %v", err)
	}

	confRead := kafka.ReaderConfig{
		Topic:     "topic2",
		Partition: 0,
		GroupID:   "g2",
		Brokers:   []string{os.Getenv("KAFKA_ADDRESS")},
	}
	que := queue.NewQueue(kafkaConn, logger)

	service := messageService.NewService(repo, que, logger)

	reader := queue2.NewReader(confRead, logger, service)
	go reader.Listen()

	srv := server.NewServer(service, logger)

	logger.Infof("Running server on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), srv.Handlers())
	if err != nil {
		logger.Errorf("Error starting server: %s", err.Error())
		return models.ErrServerFailed
	}

	return err
}
