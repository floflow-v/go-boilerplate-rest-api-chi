package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	scalargo "github.com/bdpiprava/scalar-go"
	"github.com/bdpiprava/scalar-go/model"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/database"
	internalValidator "go-boilerplate-rest-api-chi/internal/validator"
)

func buildAPI(cfg config.Config, logger zerolog.Logger, db *database.Database) *http.Server {
	validator := internalValidator.New()

	// -------- Repos / Services / Handlers --------

	bookService := book.NewBookService(db.Queries, logger)
	authorService := author.NewAuthorService(db.Queries, logger)

	bookHandler := book.NewBookHandler(bookService, validator, logger)
	authorHandler := author.NewAuthorHandler(authorService, validator, logger)

	// Router
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.GetHead,
		middleware.Timeout(10*time.Second),
		middleware.Throttle(100),
		httprate.LimitByRealIP(100, 1*time.Minute),
	)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * int(time.Hour),
	}))

	r.Route("/api", func(api chi.Router) {
		api.Use(middleware.Heartbeat("/api/alive"))

		api.Mount("/books", bookHandler.Routes())
		api.Mount("/authors", authorHandler.Routes())

		// SCALAR DOCUMENTATION
		if cfg.Api.Environment == "development" {
			api.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
				raw, err := os.ReadFile("./docs/swagger.json")
				if err != nil {
					http.Error(w, "Documentation non disponible", http.StatusInternalServerError)
					return
				}

				var swagger openapi2.T
				_ = json.Unmarshal(raw, &swagger)
				openapi3, _ := openapi2conv.ToV3(&swagger)
				content, _ := openapi3.MarshalJSON()

				html, err := scalargo.NewV2(
					scalargo.WithSpecBytes(content),
					scalargo.WithSpecModifier(func(spec *model.Spec) *model.Spec {
						spec.Servers = []model.Server{
							{URL: fmt.Sprintf("http://%s:%d/api", cfg.Api.Host, cfg.Api.Port)},
						}
						return spec
					}),
					scalargo.WithTheme(scalargo.ThemeBluePlanet),
					scalargo.WithShowToolbar(scalargo.ShowToolbarNever),
					scalargo.WithHideModels(),
				)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				_, _ = fmt.Fprint(w, html)
			})
		}
	})

	addr := fmt.Sprintf("%s:%d", cfg.Api.Host, cfg.Api.Port)

	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return server
}
