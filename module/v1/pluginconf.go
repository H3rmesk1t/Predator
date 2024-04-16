package v1

import (
	"Predator/api/plugin/findnet"
	"Predator/api/plugin/ftp"
	"Predator/api/plugin/memcached"
	"Predator/api/plugin/mongodb"
	"Predator/api/plugin/ms"
	"Predator/api/plugin/mssql"
	"Predator/api/plugin/mysql"
	"Predator/api/plugin/netbios"
	"Predator/api/plugin/oracle"
	"Predator/api/plugin/postgres"
	"Predator/api/plugin/rdp"
	"Predator/api/plugin/redis"
	"Predator/api/plugin/smb"
	"Predator/api/plugin/ssh"
	"Predator/api/plugin/telnet"
)

var PluginList = map[string]interface{}{
	"21":     ftp.FtpScan,
	"22":     ssh.SshScan,
	"23":     telnet.TelnetScan,
	"135":    findnet.FindNet,
	"139":    netbios.NetBIOS,
	"445":    smb.SmbScan,
	"1433":   mssql.MssqlScan,
	"1521":   oracle.OracleScan,
	"3306":   mysql.MysqlScan,
	"3389":   rdp.RdpScan,
	"5432":   postgres.PostgresScan,
	"6379":   redis.RedisScan,
	"11211":  memcached.MemcachedScan,
	"27017":  mongodb.MongodbScan,
	"100001": ms.MS,
	"100002": smb.SmbGhost,
	"100003": WebTitle,
	"100004": smb.SmbScanTwo,
}
