package main

import (
	"net/http"

	"github.com/pusher/grinder/models"

	"github.com/gocraft/web"
)

type Internal struct {
	*Grinder
}

func (i *Internal) Reset(w web.ResponseWriter, r *web.Request) {
	tx, err := i.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE users SET user_available = false;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		DELETE FROM matches;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (i *Internal) Match(w web.ResponseWriter, r *web.Request) {
	tx, err := i.DB.Beginx()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []models.User
	err = tx.Select(&users, `
		SELECT user_id
		FROM users
		WHERE
			user_available = true
		LIMIT 1000;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		_, err = tx.Exec(`
			INSERT INTO matches (from_user, to_user)
				SELECT $1, user_id FROM users
				WHERE
					user_available IS true AND
					user_id != $1
				ORDER BY RANDOM()
				LIMIT 5;
		`, user.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
