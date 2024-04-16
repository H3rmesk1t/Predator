package postgres

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func PostgresScan(info *config.HostInfo) (tmperr error) {
	if config.IsBrutePass {
		return
	}
	for _, user := range config.DefaultUsers["postgresql"] {
		for _, pass := range config.DefaultPasswords {
			pass = strings.Replace(pass, "{user}", string(user), -1)
			flag, err := postgresConn(info, user, pass)
			if flag == true && err == nil {
				return err
			}
		}
	}
	return tmperr
}

func postgresConn(info *config.HostInfo, user string, pass string) (flag bool, err error) {
	flag = false
	Host, Port, Username, Password := info.Host, info.Ports, user, pass
	dataSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", Username, Password, Host, Port, "postgres", "disable")
	db, err := sql.Open("postgres", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(time.Duration(config.Timeout) * time.Second)
		defer db.Close()
		err = db.Ping()
		if err == nil {
			result := fmt.Sprintf("[!] Postgres:%v:%v:%v %v", Host, Port, Username, Password)
			utils.LogSuccess(result)
			flag = true
		}
	}
	return flag, err
}
