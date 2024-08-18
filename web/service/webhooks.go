package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookNotification struct {
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"object"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var notification WebhookNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Process the notification
	fmt.Printf("Payment ID: %s, Status: %s\n", notification.Object.ID, notification.Object.Status)

	// Respond to the webhook
	w.WriteHeader(http.StatusOK)
}
