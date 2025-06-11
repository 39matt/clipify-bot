package firebase

import (
	"clipping-bot/internal/globalctx"
	"cloud.google.com/go/firestore"
	"errors"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log/slog"
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
	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load environment variables", "error", err)
		return
	}

	ctx, cancel := globalctx.ForRequest()
	defer cancel()

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
