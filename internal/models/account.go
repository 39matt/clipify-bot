package models

type Account struct {
	Username string  `json:"username" firestore:"username"`
	Platform string  `json:"platform" firestore:"platform"`
	Link     string  `json:"link" firestore:"link"`
	Videos   []Video `json:"videos" firestore:"videos"`
}
