package firebase

import (
	"clipping-bot/internal/models"
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"log/slog"
)

func AddCampaign(ctx context.Context, campaign models.Campaign) (*firestore.DocumentRef, error) {
	if !IsInitialized() {
		slog.Error("Firebase is not initialized")
		return nil, errGeneric
	}

	doc, _, err := FirestoreClient.Collection("campaigns").Add(ctx, campaign)
	if err != nil {
		slog.Error("Error adding campaign", "error", err)
		return nil, errGeneric
	}
	return doc, nil
}

func GetCampaigns(ctx context.Context) ([]*models.Campaign, []string, error) {
	if !IsInitialized() {
		slog.Error("Firebase is not initialized")
		return nil, nil, errGeneric
	}
	iter := FirestoreClient.Collection("campaigns").Documents(ctx)
	var campaigns []*models.Campaign
	var campaignIDs []string
	for {
		doc, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				break
			}
		}
		var campaign models.Campaign
		convertErr := doc.DataTo(&campaign)
		if convertErr != nil {
			slog.Error("Error converting campaign to models.Campaign", "error", convertErr)
			return nil, nil, errGeneric
		}
		campaigns = append(campaigns, &campaign)
		campaignIDs = append(campaignIDs, doc.Ref.ID)
	}
	return campaigns, campaignIDs, nil
}
