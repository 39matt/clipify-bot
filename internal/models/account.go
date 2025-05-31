package models

type Account struct {
	Username string   `json:"username" firestore:"username"`
	Platform Platform `json:"platform" firestore:"platform"`
	Link     string   `json:"link" firestore:"link"`
}
