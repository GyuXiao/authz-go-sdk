package errors

import (
	"fmt"
)

type AuthzSDKError struct {
	Code      int
	Message   string
	RequestID string
}

func (err *AuthzSDKError) Error() string {
	return fmt.Sprintf("[AuthzSDKError] code=%d, message=%s, requestID=%s", err.Code, err.Message, err.RequestID)
}

func NewAuthzSDKError(code int, message, requestID string) error {
	return &AuthzSDKError{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}
}

func (err *AuthzSDKError) GetCode() int {
	return err.Code
}

func (err *AuthzSDKError) GetMessage() string {
	return err.Message
}

func (err *AuthzSDKError) GetRequestID() string {
	return err.RequestID
}