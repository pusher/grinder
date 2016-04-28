package main

import (
	"net/http"

	"github.com/gocraft/web"

	"github.com/pusher/grinder/models"
)

type Matches struct {
	*Users
}

func (m *Matches) Index(w web.ResponseWriter, r *web.Request) {
	var matches []models.Match
	err := m.DB.Select(&matches, `
	SELECT users.user_id, users.user_name, users.user_available
	FROM users, matches
	WHERE
		matches.from_user = $1 AND
		users.user_id = matches.to_user
	LIMIT 1000;
	`, m.User.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Success(w, r, matches)
}

func (m *Matches) Match(w web.ResponseWriter, r *web.Request) {
	_, err := m.DB.Exec(`
		UPDATE matches
		SET match = true
		WHERE
			from_user = $1 AND
			to_user = $2;
	`, m.User.Id, r.PathParams["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
