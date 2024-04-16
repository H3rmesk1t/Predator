package lib

import (
	"Predator/module/v1/ymlpoc/structs"
	"Predator/pkg/xhttp"
	"strings"
)

var (
	ReversePlatformType      structs.ReverseType
	DnslogCNGetDomainRequest *xhttp.Request
	DnslogCNGetRecordRequest *xhttp.Request
)

func InitReversePlatform(api, domain string) {
	if api != "" && domain != "" && strings.HasSuffix(domain, ".ceye.io") {
		ReversePlatformType = structs.ReverseType_Ceye
	} else {
		ReversePlatformType = structs.ReverseType_DnslogCN

		// 设置请求相关参数
		DnslogCNGetDomainRequest, _ = xhttp.NewRequest("GET", "http://dnslog.cn/getdomain.php", nil)
		DnslogCNGetRecordRequest, _ = xhttp.NewRequest("GET", "http://dnslog.cn/getrecords.php", nil)

	}
}
