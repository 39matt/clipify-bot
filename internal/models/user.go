package models

type User struct {
	Username          string            `json:"username" firestore:"username"`
	UnverifiedAccount UnverifiedAccount `json:"unverified_account" firestore:"unverified_account"`
	Accounts          []Account         `json:"accounts" firestore:"accounts"`
}
