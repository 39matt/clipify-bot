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

var accountLimit = 3

func IsAccountAlreadyAdded(ctx context.Context, discordUsername string, accountName string, platform models.Platform) (bool, error) {
	if !IsInitialized() {
		slog.Error("Firebase not initialized")
		return false, errGeneric
	}

	query := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").
		Where("username", "==", accountName).
		Where("platform", "==", platform).
		Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	_, iterErr := iter.Next()
	if iterErr != nil {
		if errors.Is(iterErr, iterator.Done) {
			return false, nil
		}
		slog.Error("Error checking if account exists", "error", iterErr)
		return false, errGeneric
	}

	return true, nil
}

func IsAccountExists(ctx context.Context, platform models.Platform, accountName string) (bool, error) {
	query := FirestoreClient.CollectionGroup("accounts").
		Where("platform", "==", platform).
		Where("username", "==", accountName).
		Limit(1)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false, err
	}
	return len(docs) > 0, nil
}

func GetAccountSnapshotByNameAndPlatform(ctx context.Context, discordUsername string, accountName string, platform models.Platform) (*firestore.DocumentSnapshot, error) {
	if !IsInitialized() {
		slog.Error("Firebase not initialized")
		return nil, errGeneric
	}

	query := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").
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

func AddAccount(ctx context.Context, discordUsername string, account models.Account) (*firestore.DocumentRef, error) {
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

	ref, _, err := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Add(ctx, account)
	if err != nil {
		slog.Error("Failed to add verified account", "error", err)
		return nil, errGeneric
	}

	return ref, nil
}

func GetAllAccountVideos(ctx context.Context, discordUsername string, accountName string, platform models.Platform) ([]models.Video, error) {
	if !IsInitialized() {
		slog.Error("Firebase instance not initialized")
		return nil, errGeneric
	}
	if discordUsername == "" {
		slog.Error("Username must be set")
		return nil, errGeneric
	}

	accountSnapshot, accountErr := GetAccountSnapshotByNameAndPlatform(ctx, discordUsername, accountName, platform)
	if accountErr != nil {
		slog.Error("Error getting account", "error", accountErr)
		return nil, errGeneric
	}

	accountVideos, videoErr := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Doc(accountSnapshot.Ref.ID).Collection("videos").Documents(ctx).GetAll()
	if videoErr != nil {
		slog.Error("Failed to get all videos", "error", videoErr)
		return nil, errGeneric
	}

	videos := make([]models.Video, 0, len(accountVideos))

	for i, userVideo := range accountVideos {
		var video models.Video
		if parseErr := userVideo.DataTo(&video); parseErr != nil {
			slog.Error("Failed to parse video data", "error", parseErr, "index", i)
			continue
		}
		videos = append(videos, video)
	}

	return videos, nil
}

func GetUserAccounts(ctx context.Context, discordUsername string) ([]models.Account, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return nil, errGeneric
	}
	if discordUsername == "" {
		slog.Error("discordID cannot be empty")
		return nil, errGeneric
	}

	docIter := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Documents(ctx)
	defer docIter.Stop()

	var accounts []models.Account
	for {
		doc, iterErr := docIter.Next()
		if errors.Is(iterErr, iterator.Done) {
			break
		}
		if iterErr != nil {
			slog.Error("Error iterating through accounts", "error", iterErr)
			return nil, errGeneric
		}

		var account models.Account
		if parseErr := doc.DataTo(&account); parseErr != nil {
			slog.Error("Error parsing account data", "error", parseErr, "docID", doc.Ref.ID)
			continue
		}

		accounts = append(accounts, account)
	}

	if len(accounts) == 0 {
		return nil, nil
	}

	return accounts, nil
}

func CanUserAddAccount(ctx context.Context, discordUsername string, platform models.Platform) (bool, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return false, errGeneric
	}
	if discordUsername == "" {
		slog.Error("discordID cannot be empty")
		return false, errGeneric
	}

	docIter := FirestoreClient.Collection("users").Doc(discordUsername).Collection("accounts").Documents(ctx)
	defer docIter.Stop()

	numberOfAccounts := 0
	for {
		doc, err := docIter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			slog.Error("Error iterating through accounts", "error", err)
			return false, errGeneric
		}

		var account models.Account
		if err := doc.DataTo(&account); err != nil {
			slog.Error("Error parsing account data", "error", err, "docID", doc.Ref.ID)
			continue
		}
		if account.Platform == platform {
			numberOfAccounts++
		}
	}

	if numberOfAccounts < accountLimit {
		return true, nil
	}
	return false, nil
}
