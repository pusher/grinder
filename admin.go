package main

import (
	"html/template"
	"net/http"

	"github.com/pusher/grinder/models"

	"github.com/gocraft/web"
)

type Admin struct {
	*Grinder
}

var tplAdmin = template.Must(template.New("").Parse(`
{{ define "header" }}
<html>
	<head>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />

		<title>Grinder Admin</title>

		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous" />
	</head>
	<body>
		<nav class="navbar navbar-default navbar-static-top navbar-inverse">
			<div class="container">
				<a class="navbar-brand" href="/admin">Grinder</a>

				<ul class="nav navbar-nav">
					<li><a href="/admin/users">Users</a></li>
					<li><a href="/admin/matches">Matches</a></li>
				</ul>
			</div>
		</nav>

		<div class="container">
{{ end }}

{{ define "footer" }}
		</div>
	</body>
</html>
{{ end }}

{{ define "index" }}
{{ template "header" }}
			<a href="/admin/reset" class="btn btn-primary">Reset Availability/Matches</a>
			<a href="/admin/match" class="btn btn-primary">Create Matches</a>
{{ template "footer" }}
{{ end }}

{{ define "users" }}
{{ template "header" }}
			<table class="table">
				<thead>
					<tr>
						<th>ID</th>
						<th>Name</th>
						<th>Status</th>
						<th>&nbsp;</th>
					</tr>
				</thead>
				<tbody>
{{ range . }}
					<tr>
						<td>{{ .Id }}</td>
						<td>{{ .Name }}</td>
{{ if .Available }}
						<td><span class="label label-success">Available</span></td>
{{ else }}
						<td><span class="label label-danger">Unavailable</span></td>
{{ end }}
						<td class="text-right">
							<a href="/admin/user/{{ .Id }}/toggle" class="btn btn-primary btn-xs">Toggle Availability</a>
						</td>
					</tr>
{{ end }}
				</tbody>
			</table>
{{ template "footer" }}
{{ end }}

{{ define "matches" }}
{{ template "header" }}
			<table class="table">
				<thead>
					<tr>
						<th>From</th>
						<th>To</th>
						<th>Match</th>
						<th>&nbsp;</th>
					</tr>
				</thead>
				<tbody>
{{ range . }}
					<tr>
						<td>{{ .FromName }} ({{ .FromId }})</td>
						<td>{{ .ToName }} ({{ .ToId }})</td>
{{ if .Match }}
						<td><span class="label label-success">Match</span></td>
{{ else }}
						<td><span class="label label-danger">No Match</span></td>
{{ end }}
						<td><a href="/admin/matches/{{ .FromId }}/{{ .ToId }}/toggle" class="btn btn-primary btn-xs">Toggle Match</a></td>
					</tr>
{{ end }}
				</tbody>
			</table>
{{ template "footer" }}
{{ end }}
`))

func (a *Admin) Index(w web.ResponseWriter, r *web.Request) {
	tplAdmin.ExecuteTemplate(w, "index", nil)
}

func (a *Admin) Users(w web.ResponseWriter, r *web.Request) {
	var users []models.User

	err := a.DB.Select(&users, `
		SELECT user_id, user_name, user_available
		FROM users
		ORDER BY user_name ASC
		LIMIT 1000;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tplAdmin.ExecuteTemplate(w, "users", users)
}

func (a *Admin) Matches(w web.ResponseWriter, r *web.Request) {
	var matches []struct {
		FromId   string `db:"from_user_id"`
		FromName string `db:"from_user_name"`
		ToId     string `db:"to_user_id"`
		ToName   string `db:"to_user_name"`
		Match    bool   `db:"match"`
	}

	err := a.DB.Select(&matches, `
		SELECT
			from_user.user_id AS from_user_id,
			from_user.user_name AS from_user_name,
			to_user.user_id AS to_user_id,
			to_user.user_name AS to_user_name,
			matches.match
		FROM
			users from_user,
			users to_user,
			matches
		WHERE
			from_user.user_id = matches.from_user AND
			to_user.user_id = matches.to_user;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tplAdmin.ExecuteTemplate(w, "matches", matches)
}

func (a *Admin) Reset(w web.ResponseWriter, r *web.Request) {
	tx, err := a.DB.Beginx()
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

	w.Header().Set("Location", "/admin")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Admin) Toggle(w web.ResponseWriter, r *web.Request) {
	_, err := a.DB.Exec(`
	UPDATE users
	SET
		user_available = NOT user_available
	WHERE
		user_id = $1;
	`, r.PathParams["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/admin/users")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Admin) ToggleMatch(w web.ResponseWriter, r *web.Request) {
	_, err := a.DB.Exec(`
	UPDATE matches
	SET
		match = NOT match
	WHERE
		from_user = $1 AND
		to_user = $2
	`, r.PathParams["from_user"], r.PathParams["to_user"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/admin/matches")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (a *Admin) Match(w web.ResponseWriter, r *web.Request) {
	tx, err := a.DB.Beginx()
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

	w.Header().Set("Location", "/admin")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
