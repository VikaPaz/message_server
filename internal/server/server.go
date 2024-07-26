package server

import (
	_ "github.com/VikaPaz/message_server/docs"
	"github.com/VikaPaz/message_server/internal/server/message"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ImplServer struct {
	service message.Service
	log     *logrus.Logger
}

func NewServer(s message.Service, logger *logrus.Logger) *ImplServer {
	return &ImplServer{
		service: s,
		log:     logger,
	}
}

func (i *ImplServer) Handlers() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m := message.NewHandler(i.service, i.log)

	r.Mount("/message", m.Router())

	return r
}
