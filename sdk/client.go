package sdk

import (
	"authz-go-sdk/sdk/log"
	"authz-go-sdk/sdk/request"
	"authz-go-sdk/sdk/response"
	"fmt"
	"github.com/marmotedu/component-base/pkg/json"
	"net/http"
	"runtime"
	"strings"
)

var defaultUserAgent = fmt.Sprintf("MEDUSDKGo/%s (%s; %s) Golang/%s", Version, runtime.GOOS, runtime.GOARCH, strings.Trim(runtime.Version(), "go"))

type Client struct {
	signMethod  string
	Credential  *Credential
	Config      *Config
	ServiceName string
	Logger      log.Logger
}

type SignFunc func(r *http.Request) error

type ListOptions struct {
	Offset *int64 `json:"offset,omitempty"`
	Limit  *int64 `json:"limit,omitempty"`
}

func (c *Client) Init(serviceName string) *Client {
	c.signMethod = "jwt"
	c.Logger = log.New()
	c.ServiceName = serviceName
	return c
}

func (c *Client) WithCredential(cred *Credential) *Client {
	c.Credential = cred
	return c
}

func (c *Client) WithSecret(secretID, secretKey string) *Client {
	c.Credential = NewCredentials(secretID, secretKey)
	return c
}

func (c *Client) WithConfig(config *Config) *Client {
	c.Config = config
	c.Logger.SetLevel(config.LogLevel)
	return c
}

// Send 关键方法
func (c *Client) Send(req request.Request, resp response.Response) error {
	method := req.GetMethod()
	builder := GetParameterBuilder(method, c.Logger)
	jsonReq, _ := json.Marshal(req)
	encodedUrl, err := builder.BuildURL(req.GetURL(), jsonReq)
	if err != nil {
		return err
	}

	endPoint := c.Config.Endpoint
	if endPoint == "" {
		endPoint = fmt.Sprintf("%s/%s", defaultEndpoint, c.ServiceName)
	}
	// 1，根据传入的 AuthzRequest 和客户端配置 config，构建 HTTP 请求参数，包括请求路径和请求 body
	// 请求路径：http://iam.api.marmotedu.com:8080/v1/authz
	reqUrl := fmt.Sprintf("%s://%s/%s%s", c.Config.Scheme, endPoint, req.GetVersion(), encodedUrl)
	body, err := builder.BuildBody(jsonReq)
	if err != nil {
		return err
	}

	// 2，签发并添加认证头
	sign := func(r *http.Request) error {
		signer := NewSigner(c.signMethod, c.Credential, c.Logger)
		_ = signer.Sign(c.ServiceName, r, strings.NewReader(body))
		return err
	}

	// 3，调用 doSend，指定 HTTP 请求
	rawResponse, err := c.doSend(method, reqUrl, body, req.GetHeaders(), sign)
	if err != nil {
		return err
	}

	// 4，处理 HTTP 请求返回的结果
	return response.ParseFromHttpResponse(rawResponse, resp)
}

func (c *Client) doSend(method, url, data string, header map[string]string, sign SignFunc) (*http.Response, error) {
	client := &http.Client{Timeout: c.Config.Timeout}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		c.Logger.Errorf("%s", err.Error())
		return nil, err
	}

	c.setHeader(req, header)

	err = sign(req)
	if err != nil {
		return nil, err
	}

	return client.Do(req)
}

func (c *Client) setHeader(req *http.Request, header map[string]string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	for k, v := range req.Header {
		c.Logger.Infof("header key: %s, header value: %s", k, v)
	}
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }