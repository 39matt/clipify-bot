package firebase

import (
	"clipping-bot/internal/globalctx"
	"cloud.google.com/go/firestore"
	"encoding/base64"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log/slog"
	"os"
	"sync"
)

var (
	firebaseApp      *firebase.App
	FirestoreClient  *firestore.Client
	once             sync.Once
	mu               sync.RWMutex
	errGeneric       = errors.New("an unexpected error occurred")
	errNotRegistered = errors.New("user not registered")
)

func Initialize() {
	once.Do(initializeOnce)
}

func initializeOnce() {
	// Get the base64-encoded Firebase credentials JSON string from the environment variable
	credentialsBase64 := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsBase64 == "" {
		slog.Error("FIREBASE_CREDENTIALS not set in .env file")
		return
	}

	// Decode the base64 string
	credentialsJSON, err := base64.StdEncoding.DecodeString(credentialsBase64)
	if err != nil {
		slog.Error("Failed to decode FIREBASE_CREDENTIALS from base64", "error", err)
		return
	}

	// Validate the JSON structure (optional, but recommended)
	var credentialsMap map[string]interface{}
	err = json.Unmarshal(credentialsJSON, &credentialsMap)
	if err != nil {
		slog.Error("Invalid Firebase credentials JSON", "error", err)
		return
	}

	// Create a context for the Firebase initialization
	ctx, cancel := globalctx.ForRequest()
	defer cancel()

	// Use the decoded JSON to initialize Firebase
	opt := option.WithCredentialsJSON(credentialsJSON)

	mu.Lock()
	defer mu.Unlock()

	// Initialize the Firebase app
	firebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		slog.Error("Failed to initialize Firebase app", "error", err)
		return
	}
	slog.Info("Initiated Firebase instance!")

	// Initialize the Firestore client
	FirestoreClient, err = firebaseApp.Firestore(ctx)
	if err != nil {
		slog.Error("Failed to initialize Firestore client", "error", err)
		return
	}
	slog.Info("Initiated Firestore client!")
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
	return FirestoreClient != nil
}
