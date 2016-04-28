package main

import (
	"database/sql"
	"net/http"

	"github.com/pusher/grinder/models"

	"github.com/gocraft/web"
)

type Users struct {
	*Grinder

	User models.User
}

func (u *Users) Auth(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	token := r.Header.Get("Authorization")
	if len(token) < 64 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token = token[len(token)-64:]

	err := u.DB.Get(&u.User, `
	SELECT user_id, user_name, user_available
	FROM users
	WHERE user_token = $1
	LIMIT 1;
	`, token)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusForbidden)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	next(w, r)
}

func (u *Users) Index(w web.ResponseWriter, r *web.Request) {
	Success(w, r, u.User)
}

func (u *Users) Available(w web.ResponseWriter, r *web.Request) {
	_, err := u.DB.Exec(`
		UPDATE users
		SET
			user_available = NOT user_available
		WHERE
			user_id = $1;
	`, u.User.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.User.Available = !u.User.Available
	Success(w, r, u.User)
}
