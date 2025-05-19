package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log/slog"
	"sync"
	"time"
)

var (
	firebaseApp     *firebase.App
	FirestoreClient *firestore.Client
	once            sync.Once
	mu              sync.RWMutex
)

func Initialize() {
	once.Do(initializeOnce)
}

func initializeOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load environment variables", "error", err)
		return
	}

	opt := option.WithCredentialsFile("keys/firebase_credentials.json")

	mu.Lock()
	defer mu.Unlock()

	firebaseApp, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		slog.Error("failed to initialize firebase app", "error", err)
		return
	}
	slog.Info("Initiated firebase instance!")

	FirestoreClient, err = firebaseApp.Firestore(ctx)
	if err != nil {
		slog.Error("failed to initialize firestore client", "error", err)
		return
	}
	slog.Info("Initiated firestore client!")

	return
}

func Close() {
	mu.Lock()
	defer mu.Unlock()

	if FirestoreClient != nil {
		err := FirestoreClient.Close()
		FirestoreClient = nil
		if err != nil {
			slog.Error("failed to close firebase client", "error", err)
			return
		}
		slog.Info("Closed firebase instance!")
	}
	return
}

func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return FirestoreClient != nil && firebaseApp != nil
}

func GetVerificationByDiscordID(context context.Context, discordId string) (*firestore.DocumentSnapshot, error) {
	if !IsInitialized() {
		return nil, fmt.Errorf("firebase not initialized")
	}
	if discordId == "" {
		return nil, fmt.Errorf("discordID cannot be empty")
	}
	query := FirestoreClient.Collection("verifications").Where("discordid", "==", discordId).Limit(1)
	iter := query.Documents(context)
	defer iter.Stop()

	doc, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve verification document: %s", err)
	}
	return doc, nil
}

func AddVerification(context context.Context, discordId string, accountName string, platform string, code int) (*firestore.DocumentRef, error) {
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

	exists, err := GetVerificationByDiscordID(context, discordId)

	if err != nil {
		return nil, err
	}

	if exists != nil {
		return nil, fmt.Errorf("verification already exists for **%s** on **%s**. Please put **%d** in your bio", accountName, platform, exists.Data()["code"])
	}

	doc := FirestoreClient.Collection("verifications").NewDoc()
	data := map[string]interface{}{
		"discordid":   discordId,
		"code":        code,
		"platform":    platform,
		"accountname": accountName,
		"createdat":   time.Now(),
	}

	_, err = doc.Set(context, data)
	if err != nil {
		slog.Error("Failed to add verification", "error", err)
		return nil, fmt.Errorf("Sorry, something went wrong. Please try again later or contact us if the error keeps occurring!")
	}
	slog.Info("Verification added to firestore successfully!")
	return doc, nil
}

func RemoveVerification(context context.Context, discordId string) error {
	if !IsInitialized() {
		slog.Error("firebase not initialized")
		return fmt.Errorf("Sorry something went wrong. Please try again or contact us if the error keeps occurring!")
	}
	if discordId == "" {
		slog.Error("discordID cannot be empty")
		return fmt.Errorf("Sorry something went wrong. Please try again or contact us if the error keeps occurring!")
	}

	exists, err := GetVerificationByDiscordID(context, discordId)
	if err != nil {
		slog.Error("Failed to get verification", "error", err)
		return fmt.Errorf("Sorry something went wrong. Please try again or contact us if the error keeps occurring!")
	}
	if exists == nil {
		return fmt.Errorf("verification does not exists for **%s**", discordId)
	}

	query := FirestoreClient.Collection("verifications").Where("discordid", "==", discordId).Limit(1)

	doc, err := query.Documents(context).Next()
	if err != nil {
		return err
	}

	_, err = doc.Ref.Delete(context)
	if err != nil {
		slog.Error("Failed to delete verification", "error", err)
		return fmt.Errorf("Sorry something went wrong. Please try again or contact us if the error keeps occurring!")
	}
	return nil
}
