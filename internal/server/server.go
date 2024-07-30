package server

import (
	_ "github.com/VikaPaz/message_server/docs"
	"github.com/VikaPaz/message_server/internal/server/message"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	m := message.NewHandler(i.service, i.log)

	r.Mount("/message", m.Router())

	return r
}
