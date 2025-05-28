package firebase

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"log/slog"
)

func GetAllAdminNames(ctx context.Context) ([]string, error) {
	if !IsInitialized() {
		slog.Error("Not initialized")
		return nil, fmt.Errorf("not initialized")
	}

	var names []string
	iter := FirestoreClient.Collection("admins").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
			slog.Error("Error iterating", "error", err)
			return nil, errGeneric
		}
		names = append(names, doc.Ref.ID)
	}
	return names, nil
}

func IsAdmin(ctx context.Context, username string) (bool, error) {
	if !IsInitialized() {
		slog.Error("Not initialized")
		return false, errGeneric
	}

	admins, err := GetAllAdminNames(ctx)
	if err != nil {
		return false, err
	}
	for _, admin := range admins {
		if username == admin {
			return true, nil
		}
	}
	return false, nil
}
