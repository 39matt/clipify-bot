package firebase

import (
	"clipify-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

//func GetUserSnapshotByUsername(ctx context.Context, discordUsername string) (*firestore.DocumentSnapshot, error) {
//	if !IsInitialized() {
//		slog.Error("firebase not initialized")
//		return nil, errGeneric
//	}
//
//	query := FirestoreClient.Collection("users").Where("discord_username", "==", discordUsername).Limit(1)
//
//	iter := query.Documents(ctx)
//	defer iter.Stop()
//
//	userSnapshot, iterErr := iter.Next()
//	if iterErr != nil {
//		if errors.Is(iterErr, iterator.Done) {
//			slog.Info("user doesn't exist", "error", iterErr)
//			return nil, errNotRegistered
//		}
//		slog.Error("failed to get user snapshot", "error", iterErr)
//		return nil, errGeneric
//	}
//
//	return userSnapshot, nil
//}

//func GetUserDataByDiscordUsername(ctx context.Context, discordUsername string) (*models.User, error) {
//
//	userData := models.User{}
//	if parseErr := userSnapshot.DataTo(&userData); parseErr != nil {
//		slog.Error("Error parsing user data", "error", parseErr, "user", userSnapshot)
//		return nil, errGeneric
//	}
//
//	return &userData, nil
//}

func AddUser(ctx context.Context, discordUsername string) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("firebase is not initialized")
		return nil, errGeneric
	}

	docRef := FirestoreClient.Collection("users").Doc(discordUsername)
	_, err := docRef.Create(ctx, models.User{})
	if status.Code(err) == codes.AlreadyExists {
		slog.Info("User already exists", "user", docRef.ID)
		return nil, fmt.Errorf("user **%s** is already registered", discordUsername)
	}

	slog.Info("User added to Firestore successfully", "username", discordUsername)
	return docRef, nil

}

func UserExists(ctx context.Context, discordUsername string) (bool, error) {
	if !IsInitialized() {
		slog.Error("firebase is not initialized")
		return false, errGeneric
	}

	_, err := FirestoreClient.Collection("users").Doc(discordUsername).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, errGeneric
	}

	return true, nil
}
