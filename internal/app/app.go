package app

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"notes-service-go/internal/config"
	"notes-service-go/internal/constants"
	"notes-service-go/internal/database"
	"notes-service-go/internal/delivery"
	"notes-service-go/internal/service"
	"notes-service-go/pkg/auth"
	"notes-service-go/pkg/hash"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(constants.ErrLoadingConfig+": %s\n", err)
	}
	log.Println(constants.SuccessfulConfigLoad)

	conn, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName))
	if err != nil {
		log.Fatalf(constants.ErrConnectingToDb+": %s\n", err)
	}
	log.Println(constants.SuccessfulDBConnection)

	queries := database.New(conn)
	hasher := hash.NewBcryptHasher()
	tokenManager := auth.NewManager(cfg.AccessSigningKey, hasher)
	services := service.NewServices(service.Deps{
		Repo:           queries,
		Hasher:         hasher,
		TokenManager:   tokenManager,
		AccessTokenTTL: cfg.AccessTTL,
	})

	r := chi.NewRouter()
	h := delivery.NewHandler(services, cfg.RefreshTTL)
	h.RegisterRoutes(r)

	log.Printf(constants.ServerStart+" %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
