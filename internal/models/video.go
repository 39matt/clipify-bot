package models

import "time"

type Video struct {
	Name       string    `json:"name" firestore:"name"`
	CampaignId string    `json:"campaign_id" firestore:"campaignId"`
	Link       string    `json:"link" firestore:"link"`
	Views      int       `json:"views" firestore:"views"`
	Likes      int       `json:"likes" firestore:"likes"`
	Shares     int       `json:"shares" firestore:"shares"`
	Comments   int       `json:"comments" firestore:"comments"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
}
