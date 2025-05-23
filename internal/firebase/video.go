package firebase

import (
	"clipping-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log/slog"
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

	userSnapshot, err := GetUserSnapshotByUsername(ctx, discordUsername)
	if err != nil {
		slog.Error("error getting user", "error", err)
		return nil, err
	}

	ref, _, err := FirestoreClient.Collection("users").Doc(userSnapshot.Ref.ID).Collection("videos").Add(ctx, video)
	if err != nil {
		slog.Error("Failed to add video", "error", err)
		return nil, errGeneric
	}
	return ref, nil
}
