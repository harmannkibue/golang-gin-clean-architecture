package mratiba_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/internal/entity/intfaces"
	"github.com/harmannkibue/spectabill_psp_connector_clean_architecture/internal/usecase/utils/httprequest"
	"math/rand"
	"time"
)

var (
	BackOffStrategy = []time.Duration{
		1 * time.Second,
		3 * time.Second,
		4 * time.Second,
	}
)

func (u MRatibaUseCase) MockMpesaRatibaCallBack(ctx context.Context, args intfaces.MockMratibaCallbackRequest) (string, error) {
	// Check if it was a successes or a failure then do the http request to the server -.
	//	This is a success callback -.
	resMsg, err := u.performHttpMockCallBack(ctx, args.CallBackUrl, generateMpesaRatibaCallbackBody(args))
	return resMsg, err

}

func (u MRatibaUseCase) performHttpMockCallBack(ctx context.Context, callBackEndpoint string, args intfaces.MpesaRatibaCallback) (string, error) {
	// Send the callback with data to the specified endpoint -.
	body, err := json.Marshal(args)

	if err != nil {
		return "", err
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Cache-Control"] = "no-cache"

	response, err := u.httpRequest.PerformPost(httprequest.RequestDataParams{
		Endpoint: callBackEndpoint,
		Data:     body,
		Params:   make(map[string]string),
	}, BackOffStrategy,
		headers)

	if err != nil {
		return "", err
	}

	fmt.Printf("--------MOCK CALLBACK RESPONSE----- %+v\n", response)

	return "Callback send successfully", nil
}

// GenerateMpesaRatibaCallback dynamically creates a successful callback response
func generateMpesaRatibaCallbackBody(args intfaces.MockMratibaCallbackRequest) intfaces.MpesaRatibaCallback {
	if args.ResponseCode == "200" {
		return intfaces.MpesaRatibaCallback{
			ResponseHeader: intfaces.MpesaRatibaResponseHeader{
				ResponseRefID:       args.ResponseRefId,
				RequestRefID:        args.ResponseRefId, // Usually the same as ResponseRefID
				ResponseCode:        "0",
				ResponseDescription: "The service request is processed successfully",
			},
			ResponseBody: intfaces.MpesaRatibaResponseBody{
				ResponseData: []intfaces.MpesaRatibaResponseData{
					{Name: "TransactionID", Value: generateTransactionID()},
					{Name: "responseCode", Value: "0"},
					{Name: "Status", Value: "OKAY"},
					{Name: "Msisdn", Value: maskPhoneNumber(args.Msisdn)},
				},
			},
		}
	} else {
		return intfaces.MpesaRatibaCallback{
			ResponseHeader: intfaces.MpesaRatibaResponseHeader{
				ResponseRefID:       args.ResponseRefId,
				RequestRefID:        args.ResponseRefId, // Usually the same as ResponseRefID
				ResponseCode:        "1037",
				ResponseDescription: "Error",
			},
			ResponseBody: intfaces.MpesaRatibaResponseBody{
				ResponseData: []intfaces.MpesaRatibaResponseData{
					{Name: "TransactionID", Value: "0000000000"}, // Failed transaction ID
					{Name: "responseCode", Value: "1037"},
					{Name: "Status", Value: "ERROR"},
					{Name: "Msisdn", Value: maskPhoneNumber(args.Msisdn)},
				},
			},
		}
	}
}

// generateTransactionID creates a random TransactionID (10 alphanumeric characters)
func generateTransactionID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// maskPhoneNumber partially hides an MSISDN for privacy
func maskPhoneNumber(phone string) string {
	if len(phone) > 4 {
		return phone[:3] + "******" + phone[len(phone)-2:]
	}
	return "******"
}
