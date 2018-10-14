package setting

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"time"
)

type App struct {
	JwtSecret       string
	PasswordSalt    string
	PageSize        int
	RuntimeRootPath string

	PrefixUrl            string
	ImageSavePath        string
	ImageMaxSize         int
	ImageAllowExtensions []string

	ExportSavePath string

	LogSavePath    string
	LogSaveName    string
	LogFileExt     string
	TimeFormat     string
	QrCodeSavePath string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Alipay struct {
	AppId           string
	PrivateKey      string
	AlipayPublicKey string
	NotifyUrl       string
	ReturnUrl       string
	IsProduct		bool
}

var AlipaySetting = &Alipay{}

func init()  {
	Setup()
}

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}

	err = Cfg.Section("alipay").MapTo(AlipaySetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AlipaySetting err: %v", err)
	}

	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}
