package firebase

import (
	"clipping-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"log/slog"
)

func AddVerifiedAccount(ctx context.Context, discordUsername string, account models.Account) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
	}
	if account.Username == "" {
		slog.Error("Account name cannot be empty")
		return nil, ErrGeneric
	}
	if account.Platform == "" {
		slog.Error("Platform cannot be empty")
		return nil, ErrGeneric
	}

	userSnapshot := GetUserReferenceByUsername(ctx, discordUsername)
	ref, _, err := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("accounts").Add(ctx, account)
	if err != nil {
		slog.Error("Failed to add verification account", "error", err)
		return nil, ErrGeneric
	}
	return ref, nil
}
