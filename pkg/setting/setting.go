package setting

import (
	"github.com/go-ini/ini"
	"log"
)

type App struct {
	LogSavePath  string
	LogSaveName  string
	LogLevel     string
	ReportStatus bool `json:"ReportStatus"` //是否上报服务状态
	ReportResourceStatus bool `json:"ReportResourceStatus"` //是否上报资源状态
	ReportPeriod int64 `json:"ReportPeriod"`
	LongitudinalScalingFactor float64 `json:"LongitudinalScalingFactor"`  //缩放因子
	HorizontalScalingFactor float64 `json:"HorizontalScalingFactor"`  //缩放因子
}
type Server struct {
	Ip   string
	Port int
}
type  Touch struct{
	Bind string
	Secret string
	Cert string
	Key string
	Invert bool
	Interval int
}
var (
	cfg              *ini.File
	AppSetting       = &App{}
	ServerSetting    = &Server{}
	TouchSetting = &Touch{}
)

func InitSetting(path string) (err error) {
	if cfg, err = ini.Load(path); err != nil {
		return
	}
	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("Touch",TouchSetting)
	return
}
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
