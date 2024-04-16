package jenkins

import (
	"Predator/pkg/utils"
	"strings"
)

func Unauthorized(u string) bool {
	if req, err := utils.HttpRequset(u, "GET", "", false, nil); err == nil {
		if req.Header.Get("X-Jenkins-Session") != "" {
			if req2, err := utils.HttpRequset(u+"/script", "GET", "", false, nil); err == nil {
				if req2.StatusCode == 200 && strings.Contains(req2.Body, "Groovy script") {
					return true
				}
			}
			if req2, err := utils.HttpRequset(u+"/computer/(master)/scripts", "GET", "", false, nil); err == nil {
				if req2.StatusCode == 200 && strings.Contains(req2.Body, "Groovy script") {
					return true
				}
			}
		}
	}
	return false
}
