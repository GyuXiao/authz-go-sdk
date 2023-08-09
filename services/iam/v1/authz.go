package v1

import (
	"authz-go-sdk/sdk/request"
	"authz-go-sdk/sdk/response"
	"encoding/json"
	"github.com/ory/ladon"
)

type AuthzRequest struct {
	*request.BaseRequest

	// Resource is the resource that access is requested to.
	Resource *string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action *string `json:"action"`

	// Subject is the subject that is requesting access.
	Subject *string `json:"subject"`
	Context *ladon.Context
}

type AuthzResponse struct {
	*response.BaseResponse
	Allowed bool   `json:"allowed,omitempty"`
	Denied  bool   `json:"denied,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewAuthzRequest() (req *AuthzRequest) {
	req = &AuthzRequest{
		BaseRequest: &request.BaseRequest{
			URL:     "/authz",
			Method:  "POST",
			Header:  nil,
			Version: "v1",
		},
	}
	return
}

func NewAuthzResponse() *AuthzResponse {
	return &AuthzResponse{
		BaseResponse: &response.BaseResponse{},
	}
}

func (r *AuthzResponse) String() string {
	data, _ := json.Marshal(r)
	return string(data)
}

// Authz 客户端可以有多个 Sender 方法，其中之一可以是 Authz 执行 REST API 调用
func (c *Client) Authz(req *AuthzRequest) (resp *AuthzResponse, err error) {
	if req == nil {
		req = NewAuthzRequest()
	}

	resp = NewAuthzResponse()
	// 下面可以添加业务逻辑
	err = c.Send(req, resp)
	return
}