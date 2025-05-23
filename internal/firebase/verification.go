package firebase

import (
	"clipping-bot/internal/models"
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
	userQuery := FirestoreClient.Collection("users").Where("discord_username", "==", discordUsername).Limit(1).Documents(ctx)
	iter := userQuery
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		slog.Error("failed to get verification document", "error", err)
		return nil, errGeneric
	}

	unverifiedDataIter := FirestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("unverified_data").Limit(1).Documents(ctx)
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

func AddVerification(ctx context.Context, discordUsername string, accountName string, platform string, code int) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("firebase not initialized")
	}

	if accountName == "" {
		return nil, fmt.Errorf("account name cannot be empty")
	}
	if platform == "" {
		return nil, fmt.Errorf("platform cannot be empty")
	}
	if code < 0 {
		return nil, fmt.Errorf("code cannot be negative")
	}

	verificationExists, verificationErr := GetVerificationByDiscordUsername(ctx, discordUsername)

	if verificationErr != nil {
		return nil, verificationErr
	}

	if verificationExists != nil {
		return nil, fmt.Errorf("verification already exists for **%s** on **%s**. Please put **%s** in your bio", accountName, platform, verificationExists.Data()["code"])
	}
	doc, userErr := GetUserSnapshotByUsername(ctx, discordUsername)
	data := models.UnverifiedAccount{
		Code:     strconv.Itoa(code),
		Platform: platform,
		Username: accountName,
	}
	if userErr != nil {
		return nil, userErr
	}
	subcollection := FirestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("unverified_data")
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

	doc, userErr := GetUserSnapshotByUsername(ctx, discordUsername)

	if userErr != nil {
		slog.Error("Failed to get user", "error", err)
		return userErr
	}
	subcollection := FirestoreClient.Collection("users").Doc(doc.Ref.ID).Collection("unverified_data")
	docs, err := subcollection.Documents(ctx).GetAll()
	for _, d := range docs {
		_, err = d.Ref.Delete(ctx)
		if err != nil {
			slog.Error("Failed to delete a document from unverified_data", "docID", d.Ref.ID, "error", err)
			return errGeneric
		}
	}
	return nil
}
