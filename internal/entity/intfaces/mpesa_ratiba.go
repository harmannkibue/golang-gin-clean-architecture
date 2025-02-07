package intfaces

import (
	"context"
	"github.com/harmannkibue/golang_gin_clean_architecture/internal/usecase/repository/sqlc"
)

type IntMRatibaUsecase interface {
	CreateMpesaStandingOrder(ctx context.Context, args MpesaRatibaRequestBody) (*sqlc.Blog, error)
}

type MpesaRatibaRequestBody struct {
	StandingOrderName           string `json:"StandingOrderName"`
	ReceiverPartyIdentifierType string `json:"ReceiverPartyIdentifierType"`
	TransactionType             string `json:"TransactionType"`
	BusinessShortCode           string `json:"BusinessShortCode"`
	PartyA                      string `json:"PartyA"`
	Amount                      string `json:"Amount"`
	StartDate                   string `json:"StartDate"`
	EndDate                     string `json:"EndDate"`
	Frequency                   string `json:"Frequency"`
	AccountReference            string `json:"AccountReference"`
	TransactionDesc             string `json:"TransactionDesc"`
	CallBackURL                 string `json:"CallBackURL"`
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
