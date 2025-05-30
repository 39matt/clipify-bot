package models

type InstagramUserResponse struct {
	Biography string `json:"biography"`
}

type InstagramVideoResponse struct {
	Shortcode string `json:"shortcode"`
	ViewCount int    `json:"video_play_count"`
	Likes     struct {
		Count int `json:"count"`
	} `json:"edge_media_preview_like"`
	Comments struct {
		Count int `json:"count"`
	} `json:"edge_media_preview_comment"`
	CaptionContainer struct {
		Edges []struct {
			Node struct {
				CreatedAt   string `json:"created_at"`
				Description string `json:"text"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_media_to_caption"`
	Owner struct {
		Username string `json:"username"`
	} `json:"owner"`
}
