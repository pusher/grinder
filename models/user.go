package models

type User struct {
	Id   string `db:"user_id" json:"id"`
	Name string `db:"user_name" json:"name"`

	Available bool   `db:"user_available" json:"available"`
	Token     string `db:"user_token" json:"token,omitempty"`
}
