package models

type TikTokUserResponse struct {
	UserInfo struct {
		User struct {
			Signature string `json:"signature"`
		} `json:"user"`
	} `json:"userInfo"`
}

type TikTokVideoResponse struct {
	VideoInfo struct {
		VideoStructure struct {
			Description string `json:"desc"`
			Stats       struct {
				Shares   int `json:"shareCount"`
				Views    int `json:"playCount"`
				Comments int `json:"commentCount"`
			} `json:"stats"`
			Author struct {
				Username string `json:"uniqueId"`
			} `json:"author"`
		} `json:"itemStruct"`
	} `json:"itemInfo"`
}
