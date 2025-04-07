package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type postgresCfg struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"postgres"`
	User     string `env:"POSTGRES_USER" envDefault:"postgres"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Port     uint   `env:"POSTGRES_PORT" envDefault:"5432"`
	DBName   string `env:"POSTGRES_DB_NAME" envDefault:"postgres"`
}

var database *DB

func TestMain(m *testing.M) {
	var err error
	godotenv.Load()
	c := &postgresCfg{}
	err = env.Parse(c)
	if err != nil {
		log.Fatal(err)
	}
	db1, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DBName),
	)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})

	database = New(db1, &logger)
	if err != nil {
		log.Fatal(err)
	}
	m.Run()
}
