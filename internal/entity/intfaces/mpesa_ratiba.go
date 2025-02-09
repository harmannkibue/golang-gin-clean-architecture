package intfaces

import (
	"context"
)

type ReceiverPartyIdentifierType string
type TransactionType string

const (
	ReceiverPartyMerchantTill      ReceiverPartyIdentifierType = "2"
	ReceiverPartyBusinessShortCode ReceiverPartyIdentifierType = "4"
)

const (
	TransactionTypePayBill     TransactionType = "Standing Order Customer Pay Bill"
	TransactionTypePayMerchant TransactionType = "Standing Order Customer Pay Merchant"
)

type IntMRatibaUsecase interface {
	CreateMpesaStandingOrder(ctx context.Context, args MpesaRatibaRequestBody) (*MpesaRatibaRequestResponseBody, error)
}

type MpesaRatibaRequestBody struct {
	StandingOrderName           string                      `json:"StandingOrderName"`
	ReceiverPartyIdentifierType ReceiverPartyIdentifierType `json:"ReceiverPartyIdentifierType"`
	TransactionType             TransactionType             `json:"TransactionType"`
	BusinessShortCode           string                      `json:"BusinessShortCode"`
	PartyA                      string                      `json:"PartyA"`
	Amount                      string                      `json:"Amount"`
	StartDate                   string                      `json:"StartDate"`
	EndDate                     string                      `json:"EndDate"`
	Frequency                   string                      `json:"Frequency"`
	AccountReference            string                      `json:"AccountReference"`
	TransactionDesc             string                      `json:"TransactionDesc"`
	CallBackURL                 string                      `json:"CallBackURL"`
}

type MpesaRatibaRequestResponseBody struct {
	ResponseHeader struct {
		ResponseRefID       string `json:"responseRefID"`
		ResponseCode        string `json:"responseCode"`
		ResponseDescription string `json:"responseDescription"`
	} `json:"ResponseHeader"`
	ResponseBody struct {
		ResponseDescription string `json:"responseDescription"`
		ResponseCode        string `json:"responseCode"`
	} `json:"ResponseBody"`
}

// MockMratibaCallbackRequest This and anything else below are used for mocking and should be removed and fixed on going live since live callbacks from mpesa daraja would work ideally
type MockMratibaCallbackRequest struct {
	ResponseCode  string `json:"response_code"`
	ResponseRefId string `json:"response_ref_id"`
	Msisdn        string `json:"msisdn"`
	CallBackUrl   string `json:"call_back_url"`
}

// MpesaRatibaCallback represents the full callback response from M-Pesa Ratiba
type MpesaRatibaCallback struct {
	ResponseHeader MpesaRatibaResponseHeader `json:"responseHeader"`
	ResponseBody   MpesaRatibaResponseBody   `json:"responseBody"`
}

// MpesaRatibaResponseHeader contains metadata of the callback response
type MpesaRatibaResponseHeader struct {
	ResponseRefID       string `json:"responseRefID"`
	RequestRefID        string `json:"requestRefID"`
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
}

// MpesaRatibaResponseBody contains transaction-specific response data
type MpesaRatibaResponseBody struct {
	ResponseData []MpesaRatibaResponseData `json:"responseData"`
}

// MpesaRatibaResponseData represents an individual response field
type MpesaRatibaResponseData struct {
	Name  string `json:"name,omitempty" json:"Name,omitempty"` // Supports both lowercase and uppercase keys
	Value string `json:"value,omitempty" json:"Value,omitempty"`
}
