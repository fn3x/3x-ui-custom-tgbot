package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"x-ui/database"
	"x-ui/database/model"
	"x-ui/logger"

	"gorm.io/gorm"
)

type WebhookEvent = string

var (
	WaitingForCapture WebhookEvent = "payment.waiting_for_capture"
	Pending           WebhookEvent = "payment.pending"
	Succeeded         WebhookEvent = "payment.succeeded"
	Canceled          WebhookEvent = "payment.canceled"
)

type Webhook struct {
	Event WebhookEvent `json:"event"`
	Url   string       `json:"url"`
}

type WebhookRegistered struct {
	Id    string       `json:"id"`
	Event WebhookEvent `json:"event"`
	Url   string       `json:"url"`
}

type WebhookNotification struct {
	Type   string          `json:"type"`
	Event  WebhookEvent    `json:"event"`
	Object PaymentResponse `json:"object"`
}

type WebhookService struct {
	database       *gorm.DB
	settingService SettingService
}

func (w *WebhookService) NewWebhookService() *WebhookService {
	service := new(WebhookService)
	service.database = database.GetDB()
	return service
}

func (w *WebhookService) registerWebhook(webhook Webhook, idempotenceKey string) (WebhookRegistered, error) {
	data, _ := json.Marshal(webhook)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/webhook", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)

	shopId, err := w.settingService.GetYookassaShopId()
	if err != nil {
		logger.Errorf("Couldn't get shop id from settings. Reason: %s", err.Error())
		return WebhookRegistered{}, err
	}

	apiKey, err := w.settingService.GetYookassaApiKey()
	if err != nil {
		logger.Errorf("Couldn't get api key from settings. Reason: %s", err.Error())
		return WebhookRegistered{}, err
	}

	req.SetBasicAuth(strconv.Itoa(shopId), apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WebhookRegistered{}, err
	}
	defer resp.Body.Close()

	var webhookResponse WebhookRegistered
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		logger.Errorf("Couldn't decode response body. Reason: %s", err.Error())
		return WebhookRegistered{}, err
	}

	return webhookResponse, nil
}

func (w *WebhookService) removeWebhook(webhookId string) {
	req, _ := http.NewRequest("DELETE", "https://api.yookassa.ru/v3/webhooks/"+webhookId, nil)
	req.Header.Set("Content-Type", "application/json")

	shopId, err := w.settingService.GetYookassaShopId()
	if err != nil {
		logger.Errorf("Couldn't get shop id from settings. Reason: %s", err.Error())
		return
	}

	apiKey, err := w.settingService.GetYookassaApiKey()
	if err != nil {
		logger.Errorf("Couldn't get api key from settings. Reason: %s", err.Error())
		return
	}

	req.SetBasicAuth(strconv.Itoa(shopId), apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var webhookResponse WebhookRegistered
	if err := json.NewDecoder(resp.Body).Decode(&webhookResponse); err != nil {
		logger.Errorf("Couldn't decode response body. Reason: %s", err.Error())
		return
	}
}

func (w *WebhookService) WebhookHandler(wr http.ResponseWriter, r *http.Request) {
	var notification WebhookNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(wr, "Bad Request", http.StatusBadRequest)
		return
	}

	jsonWebhook, err := json.MarshalIndent(notification, "", "  ")
	logger.Infof("Webhook notification: %s\r\n", jsonWebhook)

	err = w.database.
		Model(&model.Payment{}).
		Where("payment_id = ?", notification.Object.Id).
		Updates(model.Payment{
			Status: notification.Object.Status,
			PaymentMethodId: notification.Object.PaymentMethod.Id,
			PaymentMethodType: notification.Object.PaymentMethod.Type,
			Saved: notification.Object.PaymentMethod.Saved,
		}).Error
	if err != nil {
		logger.Errorf("Couldn't update payment on webhook notification. Reason: %s", err.Error())
		return
	}

	wr.WriteHeader(http.StatusOK)
}
