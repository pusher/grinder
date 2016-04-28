package models

type Match struct {
	User

	Match bool `db:"match" json:"match"`
}
