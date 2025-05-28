package firebase

import (
	"clipping-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
	"strings"
)

func AddVideo(ctx context.Context, discordUsername string, video models.Video) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, fmt.Errorf("firebase instance not initialized")
	}
	if video.Comments < 0 {
		slog.Error("Comments must be greater than 0")
		return nil, fmt.Errorf("comments must be greater than 0")
	}
	if video.Name == "" {
		slog.Error("Name must be set")
		return nil, fmt.Errorf("name must be set")
	}
	if video.Link == "" {
		slog.Error("Link must be set")
		return nil, fmt.Errorf("link must be set")
	}
	if video.Shares < 0 {
		slog.Error("Shares must be greater than 0")
		return nil, fmt.Errorf("shares must be greater than 0")
	}
	if video.Views < 0 {
		slog.Error("Views must be greater than 0")
		return nil, fmt.Errorf("views must be greater than 0")
	}

	var platform string
	if strings.Contains(video.Link, "tiktok") {
		platform = "TikTok"
	} else if strings.Contains(video.Link, "instagram") {
		platform = "Instagram"
	} else {
		slog.Error("Link is not TT or IG")
		return nil, fmt.Errorf("link is not Tiktok or Instagram link")
	}

	accountName := strings.Split(video.Link, "/")[3][1:]
	accountSnapshot, accountErr := GetAccountSnapshotByNameAndPlatform(ctx, discordUsername, accountName, platform)
	if accountErr != nil {
		slog.Error("error getting account", "error", accountErr)
		return nil, errGeneric
	}

	existingVideo, err := GetAccountVideoByLink(ctx, discordUsername, video.Link)
	if err != nil {
		slog.Error("error getting existing video", "error", err)
		return nil, errGeneric
	}
	if existingVideo != nil {
		slog.Error("Video already exists")
		return nil, fmt.Errorf("video is already added")
	}
	ref, _, err := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Doc(accountSnapshot.Ref.ID).Collection("videos").Add(ctx, video)
	if err != nil {
		slog.Error("Failed to add video", "error", err)
		return nil, errGeneric
	}
	return ref, nil
}

func GetAccountVideoByLink(ctx context.Context, discordUsername string, videoLink string) (*models.Video, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, errGeneric
	}
	if videoLink == "" {
		slog.Error("Username must be set")
		return nil, errGeneric
	}

	platform := "TikTok"
	if strings.Contains(videoLink, "instagram") {
		platform = "Instagram"
	}
	accountName := strings.Split(videoLink, "/")[3][1:]
	accountSnapshot, accountErr := GetAccountSnapshotByNameAndPlatform(ctx, discordUsername, accountName, platform)
	if accountErr != nil {
		slog.Error("error getting account", "error", accountErr)
		return nil, errGeneric
	}

	videoIter := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Doc(accountSnapshot.Ref.ID).Collection("videos").Where("link", "==", videoLink).Documents(ctx)
	defer videoIter.Stop()
	videoSnap, iterErr := videoIter.Next()
	if iterErr != nil {
		if errors.Is(iterErr, iterator.Done) {
			return nil, nil
		}
		slog.Error("Error getting video", "error", iterErr)
		return nil, errGeneric
	}

	video := &models.Video{}
	err := videoSnap.DataTo(video)
	if err != nil {
		return nil, errGeneric
	}

	return video, nil
}

//func GetAllUserVideos(ctx context.Context, discordUsername string) ([]models.Video, error) {
//	if !IsInitialized() {
//		slog.Error("Firebase instance not initialized")
//		return nil, errGeneric
//	}
//	if discordUsername == "" {
//		slog.Error("Username must be set")
//		return nil, errGeneric
//	}
//
//	userVideos, videoErr := FirestoreClient.Collection("users").Doc(discordUsername).Collection("videos").Documents(ctx).GetAll()
//	if videoErr != nil {
//		slog.Error("Failed to get all videos", "error", videoErr)
//		return nil, errGeneric
//	}
//
//	// Initialize with correct capacity (not length)
//	videos := make([]models.Video, 0, len(userVideos))
//
//	for i, userVideo := range userVideos {
//		var video models.Video
//		if err := userVideo.DataTo(&video); err != nil {
//			slog.Error("Failed to parse video data", "error", err, "index", i)
//			continue
//		}
//		videos = append(videos, video)
//	}
//
//	return videos, nil
//}
