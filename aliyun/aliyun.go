package aliyun

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func main() {

	responseClient, errClient := http.Get("https://4.ipw.cn/") // 获取外网 IP
	if errClient != nil {
		fmt.Printf("获取外网 IP 失败，请检查网络\n")
		panic(errClient)
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := string(body)

	// <accessKeyId>, <accessSecret>: 前往 https://ram.console.aliyun.com/manage/ak 添加 accessKey
	// RegionId：安全组所属地域ID ，比如 `cn-guangzhou`
	// 访问 [DescribeRegions:查询可以使用的阿里云地域](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeRegions) 查阅
	// 国内一般是去掉 ECS 所在可用区的后缀，比如去掉 cn-guangzhou-b 的尾号 -b

	client, err := ecs.NewClientWithAccessKey("<RegionId>", "<accessKeyId>", "<accessSecret>")
	if err != nil {
		fmt.Print(err.Error())
	}

	request := ecs.CreateAuthorizeSecurityGroupRequest()
	request.Scheme = "https"

	request.SecurityGroupId = "<YourSecurityGroupId>" // 安全组ID
	request.IpProtocol = "tcp"                        // 协议,可选 tcp,udp, icmp, gre, all：支持所有协议
	request.PortRange = "22/22"                       // 端口范围，使用斜线（/）隔开起始端口和终止端口
	request.Priority = "1"                            // 安全组规则优先级，数字越小，代表优先级越高。取值范围：1~100
	request.Policy = "accept"                         // accept:接受访问, drop: 拒绝访问
	request.NicType = "internet"                      // internet：公网网卡, intranet：内网网卡。
	request.SourceCidrIp = clientIP                   // 源端IPv4 CIDR地址段。支持CIDR格式和IPv4格式的IP地址范围。

	response, err := client.AuthorizeSecurityGroup(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("Response: %#v\nClient IP: %s  was successfully added to the Security Group.\n", response, clientIP)
}
