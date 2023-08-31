package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	Parsetime bool `toml:"parse_time"`
	Loc       string
}

type Redis struct {
	IP     string
	Port     int
	Database int
}

type Server struct {
	IP   string
	Port int
}

type Path struct {
	StaticSourcePath string `toml:"static_source_path"`
	FfmpegPath       string `toml:"ffmpeg_path"`
}

type ConfigInfo struct {
	DB     Mysql  `toml:"mysql"`
	RDB    Redis  `toml:"redis"`
	server Server `toml:"server"`
	path   Path   `toml:"path"`
}

var configInfo ConfigInfo

func init() {
	if _, err := toml.Decode("E:/go/learn-titok-project/config/config.toml", &configInfo); err != nil {
		panic(err)
	}
	strings.Trim(configInfo.DB.Host, " ")
	strings.Trim(configInfo.RDB.IP, " ")
	strings.Trim(configInfo.server.IP, " ")
}

func GetDBConStr() string {
	dbConStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parsetime=%s&loc=%s", configInfo.DB.Username, configInfo.DB.Password, configInfo.DB.Host, configInfo.DB.Port, configInfo.DB.Database, configInfo.DB.Charset, configInfo.DB.Parsetime, configInfo.DB.Loc)
	log.Println(dbConStr)
	return dbConStr

}
