package brute

import (
	_ "embed"
	"strings"
)

type UserPass struct {
	username string
	password string
}

var (
	tomcatuserpass   = []UserPass{}
	jbossuserpass    = []UserPass{}
	top100pass       = []string{}
	weblogicuserpass = []UserPass{}
	filedic          = []string{}
)

//go:embed dicts/tomcat.txt
var szTomcatuserpass string

//go:embed dicts/jboss.txt
var szJbossuserpass string

//go:embed dicts/weblogic.txt
var szWeblogicuserpass string

//go:embed dicts/file.txt
var szFiledic string

func CvtUps(s string) []UserPass {
	a := strings.Split(s, "\n")
	var aRst []UserPass
	for _, x := range a {
		x = strings.TrimSpace(x)
		if "" == x {
			continue
		}
		j := strings.Split(x, ",")
		if 1 < len(j) {
			aRst = append(aRst, UserPass{username: j[0], password: j[1]})
		}
	}
	return aRst
}
func CvtLines(s string) []string {
	return strings.Split(s, "\n")
}
func init() {
	tomcatuserpass = CvtUps(szTomcatuserpass)
	jbossuserpass = CvtUps(szJbossuserpass)
	weblogicuserpass = CvtUps(szWeblogicuserpass)
	filedic = append(filedic, CvtLines(szFiledic)...)
}
