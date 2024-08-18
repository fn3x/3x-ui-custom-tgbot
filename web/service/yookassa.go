package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"x-ui/database/model"

	"github.com/savsgio/gotils/uuid"
)

var (
	SHOP_ID = os.Getenv("YOO_SHOP_ID")
	API_KEY = os.Getenv("YOO_API_KEY")
)

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Item struct {
	Description string `json:"description"`
	Amount      Amount `json:"amount"`
	VatCode     int    `json:"vat_code"`
	Quantity    uint   `json:"quantity"`
}

type PaymentRequest struct {
	PaymentMethodId string `json:"payment_method_id"`
	Amount          Amount `json:"amount"`
	Confirmation    struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Receipt struct {
		Items [1]Item `json:"items"`
	} `json:"receipt"`
	Capture           bool   `json:"capture"`
	Description       string `json:"description"`
	SavePaymentMethod bool   `json:"save_payment_method"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Status        model.PaymentStatus `json:"status"`
	Amount        Amount              `json:"amount"`
	CreatedAt     string              `json:"created_at"`
	PaymentMethod struct {
		Id    string `json:"id"`
		Saved bool   `json:"saved"`
	} `json:"payment_method"`
	Type        string `json:"type"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Parameter   string `json:"parameter"`
}

func createPayment(payment PaymentRequest) (PaymentResponse, error) {
	data, _ := json.Marshal(payment)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.V4())
	req.SetBasicAuth(SHOP_ID, API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return PaymentResponse{}, err
	}
	defer resp.Body.Close()

	var paymentResponse PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return PaymentResponse{}, err
	}

	return paymentResponse, nil
}
