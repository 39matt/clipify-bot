package models

type Verification struct {
	Username string `json:"username" firestore:"username"`
	Code     string `json:"code" firestore:"code"`
	Platform string `json:"platform" firestore:"platform"`
}
