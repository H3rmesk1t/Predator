package mysql

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func MysqlScan(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return
	}
	for _, user := range config.DefaultUsers["mysql"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", user, -1)
			flag, err := mysqlConn(info, user, pass)
			if flag == true && err == nil {
				return err
			}
		}
	}
	return tmperr
}

func mysqlConn(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/mysql?charset=utf8&timeout=%v", Username, Password, Host, Port, time.Duration(config.Timeout)*time.Second)
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Second)
		db.SetConnMaxIdleTime(time.Duration(config.Timeout) * time.Second)
		db.SetMaxIdleConns(0)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[!] MySQL %v:%v:%v %v", Host, Port, Username, Password)
			utils.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}
