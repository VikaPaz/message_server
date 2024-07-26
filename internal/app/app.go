package app

import (
	"github.com/VikaPaz/message_server/internal/models"
	"github.com/VikaPaz/message_server/internal/repository"
	"github.com/VikaPaz/message_server/internal/server"
	messageService "github.com/VikaPaz/message_server/internal/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Run(logger *logrus.Logger) error {
	logger.Infoln("Starting server...")
	logger.Debugf("Loading environment variables")
	if err := godotenv.Overload(); err != nil {
		logger.Errorf("Error loading .env file")
		return models.ErrLoadEnvFailed
	}

	confPostgres := message.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	logger.Infof("Connecting to PostgreSQL...")
	dbConn, err := message.Connection(confPostgres)
	if err != nil {
		logger.Errorf("Error connecting to database")
		return err
	}

	repo := message.NewRepository(dbConn, logger)

	service := messageService.NewService(repo, logger)

	srv := server.NewServer(service, logger)

	logger.Infof("Running server on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), srv.Handlers())
	if err != nil {
		logger.Errorf("Error starting server")
		return models.ErrServerFailed
	}

	return err
}
