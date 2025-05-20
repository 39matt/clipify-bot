package models

type Video struct {
	Name     string `json:"name" firestore:"name"`
	Link     string `json:"link" firestore:"link"`
	Views    int    `json:"views" firestore:"views"`
	Shares   int    `json:"shares" firestore:"shares"`
	Comments int    `json:"comments" firestore:"comments"`
}
