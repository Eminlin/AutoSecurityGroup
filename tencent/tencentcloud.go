package autocloudgroup

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

// VPC 中添加 IP组的方式管理安全组

func main() {

	credential := common.NewCredential(
		"<SecretId>",  // 请替换为你的腾讯云 SecretId
		"<SecretKey>", // 请替换为你的腾讯云 SecretKey
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, "<Region>", cpf) // 请替换为可用区ID，例如 ap-guangzhou

	request := vpc.NewModifyAddressTemplateAttributeRequest()

	responseClient, errClient := http.Get("https://ipw.cn/api/ip/myip") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		panic(errClient)
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)

	clientIP := string(body)
	params := fmt.Sprintf("{\"AddressTemplateId\":\"<AddressTemplateId>\",\"Addresses\":[\"%s\"]}", clientIP) // 请将 AddressTemplateId 替换为参数模板 ID
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.ModifyAddressTemplateAttribute(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", response.ToJsonString())
	fmt.Printf("\nOuter IP : %s\n", clientIP)

}
