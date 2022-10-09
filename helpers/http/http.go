package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/litestack-hq/lgst-tracker/helpers/config"

	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type HttpMiddlewareType func(http.Handler) http.Handler

type HandlerOpts struct {
	Conf   *config.Config
	Logger zerolog.Logger
}

const GRAPHQL_PATH = "/graph"
const GRAPHQL_PLAYGROUND_PATH = "/gql"

func Handler(opts HandlerOpts) http.Handler {
	// h := NewHttpModule(opts)
	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(httplog.RequestLogger(opts.Logger)) // already chains in request ID and recoverer middlwares
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	router.Use(middleware.Timeout(60 * time.Second))
	// router.Use(auth.HttpMiddleware(opts.DataStore, opts.Conf.APP_KEY))

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]interface{}{
			"message": "Resource not found",
		})
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusMethodNotAllowed)
		render.JSON(w, r, map[string]interface{}{
			"message": "Method not allowed",
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		// r.Post("/login", h.Login)
	})

	chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	})

	return router
}

type CorsLogger struct {
	logger *zerolog.Logger
}

func (l *CorsLogger) Printf(format string, args ...any) {
	l.logger.Info().Msg(fmt.Sprintf(format, args...))
}

func corsMiddleware(opts HandlerOpts) func(h http.Handler) http.Handler {
	allowedOrigins := []string{"*"}
	if opts.Conf.APP_ENV == "production" {
		allowedOrigins = []string{} // TODO: Configure cors list from env
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: !(opts.Conf.APP_ENV == "production"),
		Debug:            !(opts.Conf.APP_ENV == "production"),
	})

	c.Log = &CorsLogger{
		logger: &opts.Logger,
	}

	return c.Handler
}
