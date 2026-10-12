package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: creator-alert PHONE ASSET_ID URL")
		os.Exit(2)
	}
	client, err := NewSMSClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	delivery := Delivery{AssetID: os.Args[2], Title: "your new creator asset", URL: os.Args[3]}
	id, err := NotifySubscriber(context.Background(), client, os.Args[1], delivery)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("sent message:", id)
}
