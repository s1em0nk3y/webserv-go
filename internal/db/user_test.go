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

func TestCreateUser(t *testing.T) {
	user := "Username"
	passhash := "somehash"
	if err := database.CreateUser(user, passhash); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err := database.db.Exec("DELETE FROM Credentials where username = $1", user)
		if err != nil {
			t.Fatal("unable to delete test data, remove it manually")
		}
	}()
	if err := database.CreateUser(user, passhash); err == nil {
		t.Fatal("duplicate assign allowed")
	}
}

func TestCheckCredents(t *testing.T) {
	user := "Username"
	passhash := "somehash"
	id := -1
	err := database.db.QueryRow(
		`INSERT INTO Credentials(username, passhash) values ($1, $2) returning id`,
		user, passhash,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err := database.db.Exec("DELETE FROM Credentials where username = $1", user)
		if err != nil {
			t.Fatal("unable to delete test data, remove it manually")
		}
	}()
	if id == -1 {
		t.Fatal("an id was not returned")
	}

	ok, _ := database.CheckCredents(user, passhash)
	if !ok {
		t.Fatal("user not found")
	}
}

func TestValidateUsername(t *testing.T) {
	user := "Username"
	passhash := "somehash"
	id := -1
	err := database.db.QueryRow(
		`INSERT INTO Credentials(username, passhash) values ($1, $2) returning id`,
		user, passhash,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err := database.db.Exec("DELETE FROM Credentials where username = $1", user)
		if err != nil {
			t.Fatal("unable to delete test data, remove it manually")
		}
	}()

	if err := database.ValidateUsername(user); err != nil {
		t.Fatal(err)
	}

	if err := database.ValidateUsername("non-exist-user"); err == nil {
		t.Fatal("non exist user validated")
	}
}
