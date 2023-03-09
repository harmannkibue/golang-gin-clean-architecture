package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
	ErrNotFound            = errors.New("NOT_FOUND")
	ErrBadRequest          = errors.New("BAD_REQUEST")
	ErrConflict            = errors.New("CONFLICT")
	ErrInsufficientFund    = errors.New("INSUFFICIENT_FUND")
	ErrUnauthorized        = errors.New("UNAUTHORIZED")
	errStruct              ErrorCodesStruct
)

// ErrorCodesStruct This is the struct for the error codes -.
type ErrorCodesStruct struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// CreateError Creates the error code -.
func CreateError(code string, message string) error {

	err := fmt.Sprintf(`{"error_code": "%s", "error_message": "%s"}`, code, message)

	return errors.New(err)
}

// GetStatusCode Fetched the status code from the error -.
func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	errCode := extractErroCode(err)

	logrus.Error(err)
	switch errCode {
	case ErrInternalServerError.Error():
		return http.StatusInternalServerError
	case ErrNotFound.Error():
		return http.StatusNotFound
	case ErrConflict.Error():
		return http.StatusConflict
	case ErrInsufficientFund.Error():
		return http.StatusBadRequest
	case ErrUnauthorized.Error():
		return http.StatusUnauthorized
	case ErrBadRequest.Error():
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// ErrorResponse This is deprecated and should not be used -.
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Gets the error and extracts the error code from its json string -.
func extractErroCode(err error) string {
	s := err.Error()

	_ = json.Unmarshal([]byte(s), &errStruct)

	return strings.ToUpper(errStruct.ErrorCode)
}

// ErrorCodeResponse The response message for the error -.
func ErrorCodeResponse(err error) ErrorCodesStruct {
	s := err.Error()

	_ = json.Unmarshal([]byte(s), &errStruct)

	return errStruct
}
