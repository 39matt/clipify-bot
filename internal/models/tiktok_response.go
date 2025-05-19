package models

type TikTokUserResponse struct {
	UserInfo struct {
		User struct {
			Signature string `json:"signature"`
		} `json:"user"`
	} `json:"userInfo"`
}
