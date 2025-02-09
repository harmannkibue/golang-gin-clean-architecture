package mratiba_usecase

import (
	"context"
	"github.com/harmannkibue/golang-mpesa-sdk/pkg/daraja"
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/internal/entity/intfaces"
	"log"
)

func (u MRatibaUseCase) CreateMpesaStandingOrder(ctx context.Context, args intfaces.MpesaRatibaRequestBody) (*intfaces.MpesaRatibaRequestResponseBody, error) {
	darajaInstance, err := u.darajaFactory.GetDarajaInstance()

	if err != nil {
		return nil, err
	}

	var receiverType daraja.ReceiverPartyIdentifierType
	var transType daraja.TransactionType
	if args.ReceiverPartyIdentifierType == "2" {
		receiverType = daraja.ReceiverPartyMerchantTill
	} else {
		receiverType = daraja.ReceiverPartyMerchantTill
	}

	if args.TransactionType == "Standing Order Customer Pay Merchant" {
		transType = daraja.TransactionTypePayMerchant
	} else {
		transType = daraja.TransactionTypePayBill
	}

	ratibaResponse, err := darajaInstance.InitiateMpesaRatibaRequest(daraja.MpesaRatibaRequestBody{
		StandingOrderName:           args.StandingOrderName,
		ReceiverPartyIdentifierType: receiverType,
		TransactionType:             transType,
		BusinessShortCode:           args.BusinessShortCode,
		PartyA:                      args.PartyA,
		Amount:                      args.Amount,
		StartDate:                   args.StartDate,
		EndDate:                     args.EndDate,
		Frequency:                   args.Frequency,
		AccountReference:            args.AccountReference,
		TransactionDesc:             args.TransactionDesc,
		CallBackURL:                 args.CallBackURL,
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Mpesa Ratiba Response %+v", ratibaResponse)

	return &intfaces.MpesaRatibaRequestResponseBody{
		ResponseHeader: ratibaResponse.ResponseHeader,
		ResponseBody:   ratibaResponse.ResponseBody,
	}, nil
}
