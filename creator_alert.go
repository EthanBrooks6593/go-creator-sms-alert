package main

import (
	"context"
	"fmt"
)

type Delivery struct {
	AssetID string
	Title   string
	URL     string
}

// NotifySubscriber sends the one state change a subscriber needs: a published asset is ready.
func NotifySubscriber(ctx context.Context, client *SMSClient, phone string, delivery Delivery) (string, error) {
	if err := ValidateDelivery(phone, delivery); err != nil {
		return "", err
	}
	message := fmt.Sprintf("New creator drop: %s is ready. Download: %s", delivery.Title, delivery.URL)
	return client.Send(ctx, phone, message, "asset-"+delivery.AssetID)
}

func ValidateDelivery(phone string, delivery Delivery) error {
	if phone == "" || delivery.AssetID == "" || delivery.Title == "" || delivery.URL == "" {
		return fmt.Errorf("phone, asset id, title, and URL are required")
	}
	return nil
}
