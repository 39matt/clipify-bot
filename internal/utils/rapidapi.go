package utils

import (
	"clipify-bot/internal/config"
	"clipify-bot/internal/models"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

var (
	apiKey = config.GetRapidApiKey()
	client = &http.Client{Timeout: time.Second * 10}
)

func GetTiktokVideoInfo(videoId string) (models.Video, error) {
	url := fmt.Sprintf("https://tiktok-api23.p.rapidapi.com/api/post/detail?videoId=%s", videoId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "tiktok-api23.p.rapidapi.com")

	client := &http.Client{Timeout: time.Second * 10}

	res, httpErr := client.Do(req)

	if httpErr != nil {
		slog.Error("API request error", "error", httpErr)
		return models.Video{}, httpErr
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		slog.Error("Failed to read response", "error", readErr)
		return models.Video{}, fmt.Errorf("failed to process API response")
	}

	var result models.TikTokVideoResponse
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		slog.Error("JSON unmarshal error", "error", jsonErr)
	}
	if result.VideoInfo.VideoStructure.Description == "" {
		return models.Video{}, fmt.Errorf("video not found or private")
	}

	return models.Video{
		Author:    result.VideoInfo.VideoStructure.Author.Username,
		Name:      result.VideoInfo.VideoStructure.Description,
		Link:      fmt.Sprintf("https://www.tiktok.com/@%s/video/%s", result.VideoInfo.VideoStructure.Author.Username, videoId),
		Views:     result.VideoInfo.VideoStructure.Stats.Views,
		Shares:    result.VideoInfo.VideoStructure.Stats.Shares,
		Comments:  result.VideoInfo.VideoStructure.Stats.Comments,
		Likes:     result.VideoInfo.VideoStructure.Stats.Likes,
		CreatedAt: result.VideoInfo.VideoStructure.CreateTime,
	}, nil
}

func GetTiktokAccountBio(username string) (string, error) {
	url := fmt.Sprintf("https://tiktok-api23.p.rapidapi.com/api/user/info?uniqueId=%s", username)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "tiktok-api23.p.rapidapi.com")

	res, httpErr := client.Do(req)

	if httpErr != nil {
		slog.Error("API request error", "error", httpErr)
		return "", httpErr
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		slog.Error("Failed to read response", "error", readErr)
		return "", fmt.Errorf("failed to process API response")
	}

	var result models.TikTokUserResponse
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		slog.Error("JSON unmarshal error", "error", jsonErr)
	}
	return result.UserInfo.User.Signature, nil
}

func GetInstagramAccountBio(username string) (string, error) {
	url := fmt.Sprintf("https://instagram-looter2.p.rapidapi.com/profile2?username=%s", username)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "instagram-looter2.p.rapidapi.com")

	res, httpErr := client.Do(req)

	if httpErr != nil {
		slog.Error("API request error", "error", httpErr)
		return "", httpErr
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		slog.Error("Failed to read response", "error", readErr)
		return "", fmt.Errorf("failed to process API response")
	}

	var result models.InstagramUserResponse
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		slog.Error("JSON unmarshal error", "error", jsonErr)
	}

	return result.Biography, nil
}

func GetInstagramVideoInfo(videoId string) (models.Video, error) {
	url := fmt.Sprintf("https://instagram-looter2.p.rapidapi.com/post?url=https://www.instagram.com/reels/%s", videoId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", "instagram-looter2.p.rapidapi.com")

	client := &http.Client{Timeout: time.Second * 10}

	res, httpErr := client.Do(req)

	if httpErr != nil {
		slog.Error("API request error", "error", httpErr)
		return models.Video{}, httpErr
	}

	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		slog.Error("Failed to read response", "error", readErr)
		return models.Video{}, fmt.Errorf("failed to process API response")
	}

	var result models.InstagramVideoResponse
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		slog.Error("JSON unmarshal error", "error", jsonErr)
	}

	return models.Video{
		Author:    result.Owner.Username,
		Name:      result.CaptionContainer.Edges[0].Node.Description,
		Link:      fmt.Sprintf("https://www.instagram.com/reel/%s", result.Shortcode),
		Views:     result.ViewCount,
		Shares:    0,
		Comments:  result.Comments.Count,
		Likes:     result.Likes.Count,
		CreatedAt: result.CaptionContainer.Edges[0].Node.CreatedAt,
	}, nil
}
