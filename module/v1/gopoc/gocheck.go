package gopoc

import (
	"Predator/module/v1/brute"
	"Predator/module/v1/gopoc/confluence"
	"Predator/module/v1/gopoc/f5"
	"Predator/module/v1/gopoc/fastjson"
	"Predator/module/v1/gopoc/gitlab"
	"Predator/module/v1/gopoc/jboss"
	"Predator/module/v1/gopoc/jenkins"
	"Predator/module/v1/gopoc/log4j"
	"Predator/module/v1/gopoc/shiro"
	"Predator/module/v1/gopoc/spark"
	"Predator/module/v1/gopoc/springboot"
	"Predator/module/v1/gopoc/sunlogin"
	"Predator/module/v1/gopoc/thinkphp"
	"Predator/module/v1/gopoc/tomcat"
	"Predator/module/v1/gopoc/weblogic"
	"Predator/module/v1/gopoc/zabbix"
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"fmt"
	"net/url"
	"strings"
)

func CheckInfoPoc(infostr string) string {
	for _, goPoc := range GoPocDatas {
		if infostr == goPoc.Name {
			return goPoc.Alias
		}
	}
	return ""
}

func GoPocCheck(Url string, poc string) {
	var host string
	if tmp, err := url.Parse(Url); err == nil {
		host = fmt.Sprintf("%s://%s", tmp.Scheme, tmp.Host)
	}
	ip := strings.Split(strings.Split(Url, "//")[1], ":")[0]

	switch poc {
	case "Shiro":
		key := shiro.CVE_2016_4437(Url)
		if key != "" {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, key)
			utils.LogSuccess(result)
		}
	case "Tomcat":
		username, password := brute.Tomcat_brute(Url)
		if username != "" {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, username+":"+password)
			utils.LogSuccess(result)
		}
		if tomcat.CVE_2020_1938(ip) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE_2020_1938")
			utils.LogSuccess(result)
		}
		if tomcat.CVE_2017_12615(Url) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE_2017_12615")
			utils.LogSuccess(result)
		}
	case "Jboss":
		if jboss.CVE_2017_12149(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE_2017_12149")
			utils.LogSuccess(result)
		}
		username, password := brute.Jboss_brute(Url)
		if username != "" {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, username+":"+password)
			utils.LogSuccess(result)
		}
	case "ThinkPHP":
		flag, s := thinkphp.Vuln(host)
		if flag {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, s)
			utils.LogSuccess(result)
		}
	case "GitLab":
		if gitlab.CVE_2021_22205(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-22205")
			utils.LogSuccess(result)
		}
	case "SunLogin":
		if sunlogin.SunloginRCE(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "RCE")
			utils.LogSuccess(result)
		}
	case "Zabbix":
		if zabbix.CVE_2022_23131(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-23131")
			utils.LogSuccess(result)
		}
	case "Spark":
		if spark.CVE_2022_33891(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-33891")
			utils.LogSuccess(result)
		}
	case "Confluence":
		if confluence.CVE_2021_26084(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-26084")
			utils.LogSuccess(result)
		}
		if confluence.CVE_2021_26085(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-26085")
			utils.LogSuccess(result)
		}
		if confluence.CVE_2022_26134(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-26134")
			utils.LogSuccess(result)
		}
		if confluence.CVE_2022_26138(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-26138")
			utils.LogSuccess(result)
		}
		//if confluence.CVE_2023_22515(host) {
		//	result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE_2023_22515")
		//	utils.LogSuccess(result)
		//}
		if confluence.CVE_2023_22527(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2023-22627")
			utils.LogSuccess(result)
		}
	case "F5 BIG-IP":
		if f5.CVE_2020_5902(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2020-5902")
			utils.LogSuccess(result)
		}
		if f5.CVE_2022_1388(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-1388")
			utils.LogSuccess(result)
		}
		if f5.CVE_2021_22986(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-22986")
			utils.LogSuccess(result)
		}
	case "Jenkins":
		if jenkins.Unauthorized(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "Unauthorized")
			utils.LogSuccess(result)
		}
		if jenkins.CVE_2018_1000110(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2018-1000110")
			utils.LogSuccess(result)
		}
		if jenkins.CVE_2018_1000861(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2018-1000861")
			utils.LogSuccess(result)
		}
		if jenkins.CVE_2019_10003000(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2019-10003000")
			utils.LogSuccess(result)
		}
	case "SpringBoot":
		if springboot.CVE_2021_21234(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-21234")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2022_22947(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-22947")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2021_22963(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-22963")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2022_22965(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-22965")
			utils.LogSuccess(result)
		}
		if springboot.JolokiaCheck(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "Jolokia RCE")
			utils.LogSuccess(result)
		}
	case "Weblogic":
		username, password := brute.Weblogic_brute(host)
		if username != "" && username != "login_page" {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, username+":"+password)
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2014_4210(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2014-4210")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2017_3506(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2017-3506")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2017_10271(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2017-10271")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2018_2894(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2018-2894")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2019_2725(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2019-2725")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2019_2729(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2019-2729")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2020_2883(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2020-2883")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2020_14882(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2020-14882")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2020_14883(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2020-14883")
			utils.LogSuccess(result)
		}
		if weblogic.CVE_2021_2109(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-2109")
			utils.LogSuccess(result)
		}
	}
	if config.IsCheckFastjson {
		fastjsonRceType := fastjson.Check(Url)
		if fastjsonRceType != "" {
			result := fmt.Sprintf("[!] %s %s %s", Url, "Fastjson", fastjsonRceType)
			utils.LogSuccess(result)
		}
	}
	if config.IsCheckLog4j2 {
		if log4j.Check(Url) {
			result := fmt.Sprintf("[!] %s %s %s", Url, "Log4j2", "RCE")
			utils.LogSuccess(result)
		}
	}
	if config.IsCheckSpringBoot && poc != "SpringBoot" {
		poc := "SpringBoot"
		if springboot.CVE_2021_21234(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-21234")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2022_22947(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-22947")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2021_22963(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2021-22963")
			utils.LogSuccess(result)
		}
		if springboot.CVE_2022_22965(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "CVE-2022-22965")
			utils.LogSuccess(result)
		}
		if springboot.JolokiaCheck(host) {
			result := fmt.Sprintf("[!] %s %s %s", Url, poc, "Jolokia RCE")
			utils.LogSuccess(result)
		}
	}
}
