package main

import (
	"authz-go-sdk/sdk"
	iamv1 "authz-go-sdk/services/iam/v1"
	"fmt"
	"github.com/ory/ladon"
)

func main() {
	// 创建客户端
	client, _ := iamv1.NewClientWithSecret("Kdz9hhmZq6qmhnOBcsMJ1aVvWgibiWioosUJ", "sCSBRpkTVClJaoe8zFnWFionNN0WPEGG")

	// 构建请求
	req := iamv1.NewAuthzRequest()
	//req.AddHeader("Authorization", "Bearer") 在 sign 函数的逻辑会加上 Authorization:Bearer 的
	req.Resource = sdk.String("resources:articles:ladon-introduction")
	req.Action = sdk.String("delete")
	req.Subject = sdk.String("users:peter")
	ctx := ladon.Context(map[string]interface{}{"remoteIP": "192.168.0.5"})
	req.Context = &ctx

	resp, err := client.Authz(req)
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Printf("get response body: `%s`\n", resp.String())
	fmt.Printf("allowed: %v\n", resp.Allowed)
}
