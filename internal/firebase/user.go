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

func GetUserAccountNames(ctx context.Context, discordUsername string) ([]string, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return nil, errGeneric
	}
	if discordUsername == "" {
		slog.Error("discordID cannot be empty")
		return nil, errGeneric
	}

	doc, userErr := GetUserSnapshotByUsername(ctx, discordUsername)
	if userErr != nil {
		return nil, userErr
	}

	docIter := FirestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("accounts").Documents(ctx)
	defer docIter.Stop()

	var accountNames []string
	for {
		doc, err := docIter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			slog.Error("Error iterating through accounts", "error", err)
			return nil, fmt.Errorf("failed to retrieve accounts")
		}

		var account models.Account
		if err := doc.DataTo(&account); err != nil {
			slog.Error("Error parsing account data", "error", err, "docID", doc.Ref.ID)
			continue
		}

		accountNames = append(accountNames, account.Username)
	}

	if len(accountNames) == 0 {
		return nil, fmt.Errorf("no accounts found for user %s", discordUsername)
	}

	return accountNames, nil
}

func CanUserAddTikTokAccount(ctx context.Context, discordUsername string) (bool, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return false, errGeneric
	}
	if discordUsername == "" {
		slog.Error("discordID cannot be empty")
		return false, errGeneric
	}

	doc, userErr := GetUserSnapshotByUsername(ctx, discordUsername)
	if userErr != nil {
		return false, userErr
	}

	docIter := FirestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("accounts").Documents(ctx)
	defer docIter.Stop()

	numberOfAccounts := 0
	for {
		doc, err := docIter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			slog.Error("Error iterating through accounts", "error", err)
			return false, fmt.Errorf("failed to retrieve accounts")
		}

		var account models.Account
		if err := doc.DataTo(&account); err != nil {
			slog.Error("Error parsing account data", "error", err, "docID", doc.Ref.ID)
			continue
		}
		if account.Platform == "TikTok" {
			numberOfAccounts++
		}
	}

	if numberOfAccounts < 3 {
		return true, nil
	}
	return false, nil
}

func GetUserSnapshotByUsername(ctx context.Context, discordUsername string) (*firestore.DocumentSnapshot, error) {
	query := FirestoreClient.Collection("users").Where("discord_username", "==", discordUsername).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			slog.Info("user doesnt exist", "error", err)
			return nil, fmt.Errorf("user isn't registered")
		}
		slog.Error("failed to get user snapshot", "error", err)
		return nil, errGeneric
	}
	return doc, nil
}

func GetUserDataByDiscordUsername(ctx context.Context, discordUsername string) (*models.User, error) {
	userSnapshot, userErr := GetUserSnapshotByUsername(ctx, discordUsername)
	if userErr != nil {
		slog.Error("failed to get user snapshot", "error", userErr)
		return nil, errGeneric
	}
	userData := models.User{}
	if err := userSnapshot.DataTo(&userData); err != nil {
		slog.Error("Error parsing user data", "error", err, "user", userSnapshot)
		return nil, errGeneric
	}
	return &userData, nil
}

func AddUser(ctx context.Context, discordUsername string) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("firebase is not initialized")
		return nil, errGeneric
	}

	userDoc, err := GetUserSnapshotByUsername(ctx, discordUsername)

	if userDoc != nil {
		slog.Info("User already exists", "username", discordUsername)
		return nil, fmt.Errorf("user **%s** is already in use", discordUsername)
	}

	if err != nil && err.Error() == "user isn't registered" {
		doc := FirestoreClient.Collection("users").NewDoc()
		data := map[string]interface{}{
			"discord_username": discordUsername,
		}

		_, err := doc.Set(ctx, data)
		if err != nil {
			slog.Error("Failed to add user", "error", err, "username", discordUsername)
			return nil, errGeneric
		}

		slog.Info("User added to Firestore successfully", "username", discordUsername)
		return doc, nil
	}

	slog.Error("Failed to check if user exists", "error", err, "username", discordUsername)
	return nil, errGeneric
}
