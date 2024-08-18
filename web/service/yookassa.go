package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"x-ui/database/model"

	"github.com/savsgio/gotils/uuid"
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

type SinglePaymentRequest struct {
	Amount       Amount `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Receipt struct {
		Customer struct {
			Email string `json:"email"`
		} `json:"customer"`
		Items [1]Item `json:"items"`
	} `json:"receipt"`
	Capture     bool   `json:"capture"`
	Description string `json:"description"`
}

type SavePaymentRequest struct {
	SinglePaymentRequest
	PaymentMethodData struct {
		Type string `json:"type"`
	} `json:"payment_method_data"`
	SavePaymentMethod bool `json:"save_payment_method"`
}

type AutoPaymentRequest struct {
	Amount          Amount `json:"amount"`
	Capture         bool   `json:"capture"`
	Description     string `json:"description"`
	PaymentMethodId string `json:"payment_method_id"`
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

func createPayment(payment any, shopId int, apiKey string) (PaymentResponse, error) {
	data, _ := json.Marshal(payment)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", uuid.V4())
	req.SetBasicAuth(strconv.Itoa(shopId), apiKey)

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
