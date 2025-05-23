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

func IsAccountExists(ctx context.Context, discordUsername, accountName, platform string) (bool, error) {
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
