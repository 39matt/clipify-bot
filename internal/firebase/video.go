package firebase

import (
	"clipify-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
	"strconv"
	"time"
)

var videoAgeLimitHours = 48.

func AddVideo(ctx context.Context, discordUsername string, author string, platform models.Platform, video models.Video) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, fmt.Errorf("firebase instance not initialized")
	}

	accountSnapshot, accountErr := GetAccountSnapshotByNameAndPlatform(ctx, discordUsername, author, platform)
	if accountErr != nil {
		slog.Error("error getting account", "error", accountErr)
		return nil, errGeneric
	}

	existingVideo, err := getAccountVideoByLink(ctx, discordUsername, accountSnapshot, video.Link)
	if err != nil {
		slog.Error("error getting existing video", "error", err)
		return nil, errGeneric
	}
	if existingVideo != nil {
		slog.Error("Video already exists")
		return nil, fmt.Errorf("video is already added")
	}

	createdAtUnix, parseErr := strconv.ParseInt(video.CreatedAt, 10, 64)
	if parseErr != nil {
		slog.Error("error parsing video created_at", "error", parseErr)
		return nil, errGeneric
	}
	if time.Now().Sub(time.Unix(createdAtUnix, 0)).Hours() > videoAgeLimitHours {
		slog.Error(fmt.Sprintf("video is older than %.0f hours", videoAgeLimitHours))
		return nil, fmt.Errorf("video is older than %.0f hours", videoAgeLimitHours)
	}

	ref, _, err := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Doc(accountSnapshot.Ref.ID).Collection("videos").Add(ctx, video)
	if err != nil {
		slog.Error("Failed to add video", "error", err)
		return nil, errGeneric
	}
	return ref, nil
}

func getAccountVideoByLink(ctx context.Context, discordUsername string, accountSnapshot *firestore.DocumentSnapshot, videoLink string) (*models.Video, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, errGeneric
	}
	if videoLink == "" {
		slog.Error("Username must be set")
		return nil, errGeneric
	}

	var account models.Account
	parseErr := accountSnapshot.DataTo(&account)
	if parseErr != nil {
		slog.Error("Failed to parse account", "error", parseErr)
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
