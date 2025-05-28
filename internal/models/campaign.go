package models

type Campaign struct {
	Influencer         string  `json:"influencer" firestore:"influencer"`
	Activity           string  `json:"activity" firestore:"activity"`
	CreatedAt          string  `json:"created_at" firestore:"created_at"`
	Budget             string  `json:"budget" firestore:"budget"`
	Progress           float64 `json:"progress" firestore:"progress"`
	PerMillion         float64 `json:"per_million" firestore:"per_million"`
	MaxSubmissions     float64 `json:"max_submissions" firestore:"max_submissions"`
	MaxEarnings        float64 `json:"max_earnings" firestore:"max_earnings"`
	MaxEarningsPerPost float64 `json:"max_earnings_per_post" firestore:"max_earnings_per_post"`
	MinViewsForPayout  float64 `json:"min_views_for_payout" firestore:"min_views_for_payout"`
}
