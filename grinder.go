package main

import (
	"encoding/json"
	"net/http"

	"github.com/pusher/grinder/models"

	"github.com/gocraft/web"
	"github.com/jmoiron/sqlx"
)

type Grinder struct {
	DB *sqlx.DB
}

func (g *Grinder) Claims(w web.ResponseWriter, r *web.Request) {
	var users []models.User
	err := g.DB.Select(&users, `
	SELECT user_id, user_name
	FROM users
	WHERE
		user_token IS NULL
	LIMIT 1000;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Success(w, r, users)
}

func (g *Grinder) Claim(w web.ResponseWriter, r *web.Request) {
	token := RandString(64)

	_, err := g.DB.Exec(`
	UPDATE users
	SET
		user_token = $1
	WHERE
		user_id = $2 AND
		user_token IS NULL;
	`, token, r.PathParams["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user models.User
	err = g.DB.Get(&user, `
		SELECT user_id, user_name, user_token
		FROM users
		WHERE
			user_id = $1
		LIMIT 1;
	`, r.PathParams["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token != user.Token {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	Success(w, r, user)
}

func Success(w web.ResponseWriter, r *web.Request, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(body)
}
