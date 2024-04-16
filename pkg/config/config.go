package config

type HostInfo struct {
	Host    string
	Ports   string
	Url     string
	InfoStr []string
}

type PocInfoStruct struct {
	Target  string
	PocName string
}

var IsSave = true
var OutputFile = "output.txt"
var DefaultPorts = "21,22,23,80,81,88,135,139,143,389,443,445,1433,1521,2171,2181,2375,3306,3389,4444,5432,5632,5900,5984,6379,6443,7001,8000,8061,8080,8081,8086,8088,8089,8161,8443,8848,8888,9000,9043,9080,9090,9200,9300,10051,10250,11211,15672,27018,50070"
var DefaultWebPorts = "80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880"
var DefaultPasswords = []string{"", "123456~a", "Admin@123", "Aa12345", "a123456.", "root", "a12345", "1qaz@WSX", "qwer1234", "Aa123123", "pass123", "{user}111", "sa123456", "abc123456", "000000", "Charge123", "123qwe", "Passw0rd", "123123", "test", "!QAZ2wsx", "!@#$%^&*()_+", "123456789", "love123", "P@ssw0rd!", "123", "1q2w3e", "1", "123456!a", "{user}@2019", "{user}@123", "2wsx@WSX", "111111", "password", "abc123", "{user}1", "{user}_123", "pass@123", "Aa123456!", "1234qwer", "qwe123", "Aa123456789", "123321", "{user}", "1qaz2wsx", "admin123!@#", "{user}@111", "Aa123456", "1qaz!QAZ", "{user}123", "123qwe!@#", "!QAZ1qaz", "Aa1234", "123456", "Aa1234.", "system", "a123456", "sysadmin", "{user}#123", "654321", "test123", "1234567890", "a123123", "admin123", "1234567890-=", "{user}@123#4", "8888888", "a1b2c3d4", "root1234", "A123456s!", "P@ssw0rd", "admin", "Aa12345.", "666666", "12345678", "admin@123", "qwe123!@#", "a11111"}

var DefaultUsers = map[string][]string{
	"ftp":        {"root", "admin", "ftp", "ftpuser"},
	"telnet":     {"root", "admin", "telnet", "telnetuser", "cisco", "huawei"},
	"mysql":      {"root", "mysql", "admin", "test", "user", "guest"},
	"mssql":      {"sa", "SA", "root", "admin"},
	"smb":        {"administrator", "admin", "guest"},
	"rdp":        {"administrator", "admin", "guest"},
	"postgresql": {"postgres", "dbuser", "admin", "user", "test"},
	"ssh":        {"root", "admin"},
	"mongodb":    {"root", "admin", "mongouser"},
	"oracle":     {"sys", "system", "admin", "orcl"},
}

var PORTList = map[string]int{
	"ftp":      21,
	"ssh":      22,
	"telnet":   23,
	"findnet":  135,
	"netbios":  139,
	"smb":      445,
	"mssql":    1433,
	"oracle":   1521,
	"mysql":    3306,
	"rdp":      3389,
	"psql":     5432,
	"redis":    6379,
	"mem":      11211,
	"mgo":      27017,
	"ms17010":  100001,
	"smbghost": 100002,
	"web":      100003,
	"webscan":  100003,
	"smb2":     100004,
	"all":      0,
	"portscan": 0,
	"icmp":     0,
	"main":     0,
}

var PortGroup = map[string]string{
	"db":      "1433,1434,1435,14333,1521,1522,11521,3306,3307,3308,33060,5432,6379,11211,27017",
	"cloud":   "2375,2376,2377,2378,2379,2380,4222,6443,6650,8222,8300,8301,8302,8500,8848,8849,8850,9092,9700,10250,10255,10256,15010,15011,15012,15001,14268,16686",
	"service": "21,22,2222,135,139,445,3389,13389,33389,1433,1434,1435,14333,1521,1522,11521,3306,3307,3308,33060,5432,6379,11211,27017",
	"web":     "80,81,82,83,84,85,86,87,88,89,90,91,92,98,99,443,800,801,808,880,888,889,1000,1010,1080,1081,1082,1099,1118,1888,2008,2020,2100,2375,2379,3000,3008,3128,3505,5555,6080,6648,6868,7000,7001,7002,7003,7004,7005,7007,7008,7070,7071,7074,7078,7080,7088,7200,7680,7687,7688,7777,7890,8000,8001,8002,8003,8004,8006,8008,8009,8010,8011,8012,8016,8018,8020,8028,8030,8038,8042,8044,8046,8048,8053,8060,8069,8070,8080,8081,8082,8083,8084,8085,8086,8087,8088,8089,8090,8091,8092,8093,8094,8095,8096,8097,8098,8099,8100,8101,8108,8118,8161,8172,8180,8181,8200,8222,8244,8258,8280,8288,8300,8360,8443,8448,8484,8800,8834,8838,8848,8858,8868,8879,8880,8881,8888,8899,8983,8989,9000,9001,9002,9008,9010,9043,9060,9080,9081,9082,9083,9084,9085,9086,9087,9088,9089,9090,9091,9092,9093,9094,9095,9096,9097,9098,9099,9100,9200,9443,9448,9800,9981,9986,9988,9998,9999,10000,10001,10002,10004,10008,10010,10250,12018,12443,14000,16080,18000,18001,18002,18004,18008,18080,18082,18088,18090,18098,19001,20000,20720,21000,21501,21502,28018,20880",
	"main":    "21,22,80,81,82,111,135,139,389,443,445,873,888,1099,1433,1521,2049,2181,2222,2375,2379,2888,3128,3306,3389,3690,3888,4000,4040,4440,4848,4899,5000,5005,5432,5601,5631,5632,5900,5984,3,6379,7001,7051,7077,7180,7182,7848,8019,8020,8042,8048,8051,8069,8080,8081,8083,8086,8088,8161,8443,8649,8848,8880,8888,9000,9001,9042,9043,9083,9092,9100,9200,9300,9990,10000,11000,11111,11211,18080,19888,20880,25000,25010,27017,50000,50030,50070,50090,60000,60010,60030,27017,27018",
	"all":     "1-65535",
}

// scan
var (
	HostFile     string
	NoHosts      string
	PortFile     string
	UserFile     string
	PasswordFile string
	Ports        string
	NoPorts      string
	AddPorts     string
	Username     string
	Password     string
	AddUsers     string
	AddPasswords string
	ScanType     string
	SshKey       string
	Domain       string
	PocFile      string
	Url          string
	UrlFile      string
	Cookie       string
	Hash         string

	Thread      int
	LiveTop     int
	BruteThread int
	PocNum      int
	Timeout     int64
	WaitTime    int64
	WebTimeout  int64

	NoPoc             bool
	IsBrutePass       bool
	Ping              bool
	NoPing            bool
	NoSave            bool
	Silent            bool
	NoColor           bool
	DnsLog            bool
	IsCheckFastjson   bool
	IsCheckLog4j2     bool
	IsCheckSpringBoot bool

	HashBytes []byte
	Urls      []string
	HostPort  []string

	PocInfo PocInfoStruct
)

// spy
var (
	SpyModule  string
	SpyOutput  string
	SpyCidr    string
	SpyEnd     string
	SpyTcpPort string
	SpyUdpPort string

	SpyThread  int
	SpyTimes   int
	SpyTimeout int
	SpyRandom  int

	SpySpecial bool
	SpySilent  bool
	SpyDebug   bool
	SpyRapid   bool
	SpyForce   bool
)

var ModuleFlag bool

var (
	CeyeApi    = "d2891e603448e28dcce5aa5d6680311e"
	CeyeDomain = "fgtee4.ceye.io"
	Accept     = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
)
