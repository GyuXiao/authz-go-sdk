package response

import (
	"authz-go-sdk/sdk/errors"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response interface {
	ParseErrorFromHTTPResponse(body []byte) error
}

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type BaseResponse struct {
	ErrorResponse
}

func (r *BaseResponse) ParseErrorFromHTTPResponse(body []byte) error {
	if err := json.Unmarshal(body, r); err != nil {
		return err
	}
	if r.Code > 0 {
		return errors.NewAuthzSDKError(r.Code, r.Message, r.RequestID)
	}
	return nil
}

func ParseFromHttpResponse(rawResponse *http.Response, response Response) error {
	defer rawResponse.Body.Close()

	body, err := io.ReadAll(rawResponse.Body)
	if err != nil {
		return err
	}

	if rawResponse.StatusCode != 200 {
		return fmt.Errorf("request fail with status: %s, with body: %s", rawResponse.Status, body)
	}

	if err := response.ParseErrorFromHTTPResponse(body); err != nil {
		return err
	}

	return json.Unmarshal(body, &response)
}