package main

import "testing"

func TestDeliveryValidation(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		delivery Delivery
		wantErr  bool
	}{
		{"complete delivery", "+15551234567", Delivery{"drop-7", "Studio notes", "https://creator.example/drop-7"}, false},
		{"missing URL", "+15551234567", Delivery{"drop-7", "Studio notes", ""}, true},
		{"missing subscriber", "", Delivery{"drop-7", "Studio notes", "https://creator.example/drop-7"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDelivery(tt.phone, tt.delivery)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
