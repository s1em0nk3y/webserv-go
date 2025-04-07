package db

import "testing"

func TestAddReferral(t *testing.T) {
	users := []struct {
		user     string
		passhash string
		id       int
	}{
		{"testuser1", "passhash", -1},
		{"testuser2", "passhash", -1},
	}

	for i, user := range users {
		var id int
		err := database.db.QueryRow(
			`INSERT INTO Credentials(username, passhash) values ($1, $2) returning (id)`,
			user.user, user.passhash,
		).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		users[i].id = id
		defer func() {
			_, err := database.db.Exec("DELETE FROM Credentials where username = $1", user.user)
			if err != nil {
				t.Fatal("unable to delete test data, remove it manually")
			}
		}()
	}

	// Valid
	if err := database.AddNewReferral(users[0].user, users[1].user); err != nil {
		t.Fatal(err)
	}

	// Testing ids
	row := database.db.QueryRow("SELECT * FROM Referrals where id = $1", users[0].id)
	userID, referID := -1, -1
	err := row.Scan(&userID, &referID)
	if err != nil {
		t.Fatal("unable to read rows", err)
	}
	if users[0].id != userID {
		t.Fatal("user ids not eq")
	}
	if users[1].id != referID {
		t.Fatal("refer ids not eq")
	}
}
