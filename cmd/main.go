package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/s1em0nk3y/webserv-go/internal/app"
	"github.com/s1em0nk3y/webserv-go/internal/authenticators/jwt"
	"github.com/s1em0nk3y/webserv-go/internal/db"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {
	time.Sleep(time.Second * 3)
	log := zerolog.New(os.Stdout)

	c := &Config{}
	err := env.Parse(c)
	if err != nil {
		log.Fatal().Err(err).Msg("cant load required env vars")
	}

	// Create  DB connection
	database, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.DBName),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cant init db connection")
	}
	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("cant init db driver")
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migration/queries", // TODO: move to env
		"postgres", driver,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cant init migration instance")
	}
	if err := m.Migrate(5); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("cant migrate due to err")
	}

	dbDriver := db.New(database, &log)
	authenticator := jwt.NewAuthenticator("HS256", c.JWT.SignKey, dbDriver)
	app := &app.App{
		ReferralStorage:   dbDriver,
		TaskCompleter:     dbDriver,
		LeaderBoardGetter: dbDriver,
		UserStatusGetter:  dbDriver,
		Authenticator:     authenticator,
	}
	app.Run(c.ListenPort)
}

type Config struct {
	Postgres struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"postgres"`
		User     string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"POSTGRES_PASSWORD,required"`
		Port     uint   `env:"POSTGRES_PORT" envDefault:"5432"`
		DBName   string `env:"POSTGRES_DB_NAME" envDefault:"postgres"`
	}
	ListenPort uint `env:"LISTEN_PORT" envDefault:"8080"`
	JWT        struct {
		SignKey string `env:"JWT_SIGN_KEY,required"`
	}
}
