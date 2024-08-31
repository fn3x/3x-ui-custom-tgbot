package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"x-ui/database/model"
)

type (
	CancellationParty  = string
	CancellationReason = string
)

const (
	YOO_MONEY                     CancellationParty  = "yoo_money"
	PAYMENT_NETWORK               CancellationParty  = "payment_network"
	MERCHANT                      CancellationParty  = "merchant"
	SECURE_FAILED                 CancellationReason = "3d_secure_failed"
	CALL_ISSUER                   CancellationReason = "call_issuer"
	CANCELED_BY_MERCHANT          CancellationReason = "canceled_by_merchant"
	CARD_EXPIRED                  CancellationReason = "card_expired"
	COUNTRY_FORBIDDEN             CancellationReason = "country_forbidden"
	DEAL_EXPIRED                  CancellationReason = "deal_expired"
	EXPIRED_ON_CAPTURE            CancellationReason = "expired_on_capture"
	EXPIRED_ON_CONFIRMATION       CancellationReason = "expired_on_confirmation"
	FRAUD_SUSPECTED               CancellationReason = "fraud_suspected"
	GENERAL_DECLINE               CancellationReason = "general_decline"
	IDENTIFICATION_REQUIRED       CancellationReason = "identification_required"
	INSUFFICIENT_FUNDS            CancellationReason = "insufficient_funds"
	INTERNAL_TIMEOUT              CancellationReason = "internal_timeout"
	INVALID_CARD_NUMBER           CancellationReason = "invalid_card_number"
	INVALID_CSC                   CancellationReason = "invalid_csc"
	ISSUER_UNAVAILABLE            CancellationReason = "issuer_unavailable"
	PAYMENT_METHOD_LIMIT_EXCEEDED CancellationReason = "payment_method_limit_exceeded"
	PAYMENT_METHOD_RESTRICTED     CancellationReason = "payment_method_restricted"
	PERMISSION_REVOKED            CancellationReason = "permission_revoked"
	UNSUPPORTED_MOBILE_OPERATOR   CancellationReason = "unsupported_mobile_operator"
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
	Test        bool   `json:"test"`
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
	Id           string `json:"id"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Status        model.PaymentStatus `json:"status"`
	Amount        Amount              `json:"amount"`
	CreatedAt     string              `json:"created_at"`
	PaymentMethod struct {
		Id    string `json:"id"`
		Saved bool   `json:"saved"`
		Type  string `json:"type"`
	} `json:"payment_method"`
	CancellationDetails struct {
		Party  CancellationParty  `json:"party"`
		Reason CancellationReason `json:"reason"`
	} `json:"cancellation_details"`
	Test        bool   `json:"test"`
	Type        string `json:"type"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Parameter   string `json:"parameter"`
}

func createPayment(payment any, idempotenceKey string, shopId int, apiKey string) (PaymentResponse, error) {
	data, _ := json.Marshal(payment)
	req, _ := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)
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
