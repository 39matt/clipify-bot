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
)

func GetVerificationByDiscordUsername(ctx context.Context, discordUsername string) (*firestore.DocumentSnapshot, error) {
	if !IsInitialized() {
		slog.Error("firebase instance not initialized")
		return nil, errGeneric
	}
	if discordUsername == "" {
		slog.Error("missing discord id")
		return nil, fmt.Errorf("discordID cannot be empty")
	}

	unverifiedDataIter := FirestoreClient.Collection("users").Doc(discordUsername).Collection("verifications").Limit(1).Documents(ctx)
	defer unverifiedDataIter.Stop()

	dataSnapshot, err := unverifiedDataIter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		slog.Error("failed to get verification document", "error", err)
		return nil, errGeneric
	}

	return dataSnapshot, nil
}

func AddVerification(ctx context.Context, discordUsername string, accountName string, platform models.Platform, code int) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return nil, errGeneric
	}

	if accountName == "" {
		slog.Error("missing account name")
		return nil, fmt.Errorf("account name cannot be empty")
	}
	if platform == "" {
		slog.Error("missing platform")
		return nil, fmt.Errorf("platform cannot be empty")
	}
	if code < 0 {
		slog.Error("invalid verification code", "code", code)
		return nil, fmt.Errorf("invalid verification code")
	}

	verificationExists, verificationErr := GetVerificationByDiscordUsername(ctx, discordUsername)

	if verificationErr != nil {
		slog.Error("Error getting verification document", "error", verificationErr)
		return nil, verificationErr
	}

	if verificationExists != nil {
		return nil, fmt.Errorf("verification already exists for **%s** on **%s**. Please put **%s** in your bio", accountName, platform, verificationExists.Data()["code"])
	}

	data := models.Verification{
		Code:     strconv.Itoa(code),
		Platform: platform,
		Username: accountName,
	}

	subcollection := FirestoreClient.Collection("users").Doc(discordUsername).Collection("verifications")
	newDoc := subcollection.NewDoc()

	_, err := newDoc.Set(ctx, data)
	if err != nil {
		slog.Error("Failed to add verification", "error", err)
		return nil, errGeneric
	}
	slog.Info("Verification added to firestore successfully!")
	return newDoc, nil
}

func RemoveVerification(ctx context.Context, discordUsername string) error {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return errGeneric
	}
	if discordUsername == "" {
		slog.Error("discordID cannot be empty")
		return errGeneric
	}

	exists, err := GetVerificationByDiscordUsername(ctx, discordUsername)
	if err != nil {
		slog.Error("Failed to get verification", "error", err)
		return errGeneric
	}
	if exists == nil {
		return fmt.Errorf("**%s** does not have any pending verifications", discordUsername)
	}

	subcollection := FirestoreClient.Collection("users").Doc(discordUsername).Collection("verifications")
	docs, err := subcollection.Documents(ctx).GetAll()
	for _, d := range docs {
		_, deleteErr := d.Ref.Delete(ctx)
		if deleteErr != nil {
			slog.Error("Failed to delete a document from unverified_data", "docID", d.Ref.ID, "error", deleteErr)
			return errGeneric
		}
	}
	return nil
}
