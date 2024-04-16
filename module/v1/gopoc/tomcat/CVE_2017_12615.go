package tomcat

import (
	"Predator/pkg/utils"
)

func CVE_2017_12615(url string) bool {
	if req, err := utils.HttpRequset(url+"/vtset.txt", "PUT", "test", false, nil); err == nil {
		if req.StatusCode == 204 || req.StatusCode == 201 {
			return true
		}
	}
	return false
}
