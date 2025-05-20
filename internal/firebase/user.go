package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
)

func GetUser(ctx context.Context, discordUsername string) (*firestore.DocumentSnapshot, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("firebase instance not initialized")
	}

	if discordUsername == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	iter := FirestoreClient.Collection("users").Where("discord_username", "==", discordUsername).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()

	if err != nil {
		if errors.Is(err, iterator.Done) {
			slog.Error("user doesnt exist", "error", err)
			return nil, nil
		}
		slog.Error("failed to get user", "error", err)
		return nil, ErrGeneric
	}
	return doc, nil
}

func GetUserReferenceByUsername(ctx context.Context, discordUsername string) *firestore.DocumentSnapshot {
	query := FirestoreClient.Collection("users").Where("discord_username", "==", discordUsername).Limit(1)
	iter := query.Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err != nil {
		slog.Error("failed to get verification document", "error", err)
	}
	return doc
}

func AddUser(ctx context.Context, discordUsername string) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("firebase instance not initialized")
	}

	userDoc, getErr := GetUser(ctx, discordUsername)
	if getErr != nil {
		slog.Error("failed to get user", "error", getErr)
		return nil, ErrGeneric
	}
	if userDoc != nil {
		slog.Error("user already exists", "error", getErr)
		return nil, fmt.Errorf("User **%s** already exists", discordUsername)
	}

	doc := FirestoreClient.Collection("users").NewDoc()
	data := map[string]interface{}{
		"discord_username": discordUsername,
	}
	_, err := doc.Set(ctx, data)
	if err != nil {
		slog.Error("failed to add user", "error", err)
		return nil, ErrGeneric
	}

	slog.Info("User added to firestore successfully!")
	return doc, nil
}
