package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"notes-service-go/internal/config"
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery/handlers"
	"notes-service-go/internal/service"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/hash"
	"notes-service-go/pkg/spell"
)

const (
	errLoadingConfig  = "error loading config"
	errConnectingToDb = "error connecting to db"

	successfulConfigLoad   = "config has been loaded successfully"
	successfulDBConnection = "successful connection to db"
	serverStart            = "server starting on port"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(errLoadingConfig+": %s\n", err)
	}
	log.Println(successfulConfigLoad)

	conn, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		log.Fatalf(errConnectingToDb+": %s\n", err)
	}
	log.Println(successfulDBConnection)
	queries := database.New(conn)

	hasher := hash.NewBcryptHasher()
	speller := spell.NewYandexSpeller(cfg.SpellerURL)
	tokenManager := auth.NewManager(cfg.AccessSigningKey, hasher)
	services := service.NewServices(service.Deps{
		Repo:           queries,
		Hasher:         hasher,
		Speller:        speller,
		TokenManager:   tokenManager,
		AccessTokenTTL: cfg.AccessTTL,
	})

	r := chi.NewRouter()
	h := handlers.NewHandler(services, validator.New(), cfg.RefreshTTL)
	h.RegisterRoutes(r)

	log.Printf(serverStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
