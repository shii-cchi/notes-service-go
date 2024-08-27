package config

import (
	"errors"
	"github.com/joho/godotenv"
	"notes-service-go/internal/domain"
	"os"
	"time"
)

type Config struct {
	Port             string
	DbUser           string
	DbPassword       string
	DbHost           string
	DbPort           string
	DbName           string
	AccessTTL        time.Duration
	AccessSigningKey string
	RefreshTTL       time.Duration
	SpellerURL       string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")

	if port == "" {
		return nil, errors.New("PORT " + domain.ErrUndefinedEnvParam)
	}

	dbUser := os.Getenv("DB_USER")

	if dbUser == "" {
		return nil, errors.New("DB_USER " + domain.ErrUndefinedEnvParam)
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	if dbPassword == "" {
		return nil, errors.New("DB_PASSWORD " + domain.ErrUndefinedEnvParam)
	}

	dbHost := os.Getenv("DB_HOST")

	if dbHost == "" {
		return nil, errors.New("DB_HOST " + domain.ErrUndefinedEnvParam)
	}

	dbPort := os.Getenv("DB_PORT")

	if dbPort == "" {
		return nil, errors.New("DB_PORT " + domain.ErrUndefinedEnvParam)
	}

	dbName := os.Getenv("DB_NAME")

	if dbName == "" {
		return nil, errors.New("DB_NAME " + domain.ErrUndefinedEnvParam)
	}

	accessTTLStr := os.Getenv("ACCESS_TTL")

	if accessTTLStr == "" {
		return nil, errors.New("ACCESS_TTL " + domain.ErrUndefinedEnvParam)
	}

	accessTTL, err := time.ParseDuration(accessTTLStr)

	if err != nil {
		return nil, errors.New(domain.ErrParsingAccessTTL)
	}

	accessSigningKey := os.Getenv("ACCESS_SIGNING_KEY")

	if accessSigningKey == "" {
		return nil, errors.New("ACCESS_SIGNING_KEY " + domain.ErrUndefinedEnvParam)
	}

	refreshTTLStr := os.Getenv("REFRESH_TTL")

	if refreshTTLStr == "" {
		return nil, errors.New("REFRESH_TTL " + domain.ErrUndefinedEnvParam)
	}

	refreshTTL, err := time.ParseDuration(refreshTTLStr)

	if err != nil {
		return nil, errors.New(domain.ErrParsingRefreshTTL)
	}

	spellerURL := os.Getenv("SPELLER_URL")

	if spellerURL == "" {
		return nil, errors.New("SPELLER_URL" + domain.ErrUndefinedEnvParam)
	}

	return &Config{
		Port:             port,
		DbUser:           dbUser,
		DbPassword:       dbPassword,
		DbHost:           dbHost,
		DbPort:           dbPort,
		DbName:           dbName,
		AccessTTL:        accessTTL,
		AccessSigningKey: accessSigningKey,
		RefreshTTL:       refreshTTL,
		SpellerURL:       spellerURL,
	}, nil
}
