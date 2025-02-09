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
