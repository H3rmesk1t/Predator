package springboot

import (
	"Predator/pkg/utils"
	"strings"
)

func JolokiaCheck(u string) bool {
	if req, err := utils.HttpRequset(u+"/jolokia", "POST", "test", false, nil); err == nil {
		if req.StatusCode == 200 {
			if req1, err := utils.HttpRequset(u+"/jolokia/list", "GET", "", false, nil); err == nil {
				if req1.StatusCode == 200 && (strings.Contains(req.Body, "reloadByURL") || strings.Contains(req.Body, "createJNDIRealm")) {
					return true
				}
			}
		}
	}
	if req, err := utils.HttpRequset(u+"/actuator/jolokia", "POST", "test", false, nil); err == nil {
		if req.StatusCode == 200 {
			if req1, err := utils.HttpRequset(u+"/actuator/jolokia/list", "GET", "", false, nil); err == nil {
				if req1.StatusCode == 200 && (strings.Contains(req.Body, "reloadByURL") || strings.Contains(req.Body, "createJNDIRealm")) {
					return true
				}
			}
		}
	}
	if req, err := utils.HttpRequset(u+"/api/actuator/jolokia", "POST", "test", false, nil); err == nil {
		if req.StatusCode == 200 {
			if req1, err := utils.HttpRequset(u+"/api/actuator/jolokia/list", "GET", "", false, nil); err == nil {
				if req1.StatusCode == 200 && (strings.Contains(req.Body, "reloadByURL") || strings.Contains(req.Body, "createJNDIRealm")) {
					return true
				}
			}
		}
	}
	return false
}
