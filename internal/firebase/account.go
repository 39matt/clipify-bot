package firebase

import (
	"clipping-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
)

func IsAccountExists(ctx context.Context, discordUsername string, accountName string, platform string) (bool, error) {
	userSnapshot, err := GetUserSnapshotByUsername(ctx, discordUsername)
	if err != nil {
		return false, err
	}

	query := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("accounts").
		Where("username", "==", accountName).
		Where("platform", "==", platform).
		Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	_, err = iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return false, nil
		}
		slog.Error("Error checking if account exists", "error", err)
		return false, errGeneric
	}

	return true, nil
}

func GetAccountSnapshotByNameAndPlatform(ctx context.Context, discordUsername string, accountName string, platform string) (*firestore.DocumentSnapshot, error) {
	userSnapshot, err := GetUserSnapshotByUsername(ctx, discordUsername)
	if err != nil {
		return nil, err
	}

	query := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("accounts").
		Where("username", "==", accountName).
		Where("platform", "==", platform).
		Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	accountSnapshot, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		slog.Error("Error checking if account exists", "error", err)
		return nil, errGeneric
	}
	return accountSnapshot, nil
}

func AddVerifiedAccount(ctx context.Context, discordUsername string, account models.Account) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("Firebase not initialized")
		return nil, errGeneric
	}
	if account.Username == "" {
		slog.Error("Account name cannot be empty")
		return nil, fmt.Errorf("account name cannot be empty")
	}
	if account.Platform == "" {
		slog.Error("Platform cannot be empty")
		return nil, fmt.Errorf("platform cannot be empty")
	}

	userSnapshot, err := GetUserSnapshotByUsername(ctx, discordUsername)
	if err != nil {
		slog.Error("Failed to get user snapshot", "error", err)
		return nil, err
	}

	ref, _, err := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("accounts").Add(ctx, account)
	if err != nil {
		slog.Error("Failed to add verified account", "error", err)
		return nil, errGeneric
	}

	return ref, nil
}

func GetAllAccountVideos(ctx context.Context, discordUsername string, accountName string, platform string) ([]models.Video, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, errGeneric
	}
	if discordUsername == "" {
		slog.Error("Username must be set")
		return nil, errGeneric
	}

	userSnapshot, userErr := GetUserSnapshotByUsername(ctx, discordUsername)
	if userErr != nil {
		slog.Error("Error getting user", "error", userErr)
		return nil, errGeneric
	}
	accountSnapshot, accountErr := GetAccountSnapshotByNameAndPlatform(ctx, discordUsername, accountName, platform)
	if accountErr != nil {
		slog.Error("Error getting account", "error", accountErr)
		return nil, errGeneric
	}

	accountVideos, videoErr := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("accounts").Doc(accountSnapshot.Ref.ID).Collection("videos").Documents(ctx).GetAll()
	if videoErr != nil {
		slog.Error("Failed to get all videos", "error", videoErr)
		return nil, errGeneric
	}

	videos := make([]models.Video, 0, len(accountVideos))

	for i, userVideo := range accountVideos {
		var video models.Video
		if err := userVideo.DataTo(&video); err != nil {
			slog.Error("Failed to parse video data", "error", err, "index", i)
			continue
		}
		videos = append(videos, video)
	}

	return videos, nil
}
