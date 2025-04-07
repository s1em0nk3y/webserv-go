package db

import (
	"testing"

	_ "github.com/lib/pq"
)

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
