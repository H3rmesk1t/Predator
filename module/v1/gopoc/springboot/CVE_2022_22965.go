package springboot

import (
	"Predator/pkg/utils"
)

func CVE_2022_22965(u string) bool {
	if req, err := utils.HttpRequset(u+"/?class.module.classLoader%5b1%5d=1", "GET", "", false, nil); err == nil {
		if req.StatusCode == 500 {
			if req2, err := utils.HttpRequset(u+"/?class.module.classLoader=1", "GET", "", false, nil); err == nil {
				if req2.StatusCode == 200 {
					return true
				}
			}
		}
	}
	return false
}
