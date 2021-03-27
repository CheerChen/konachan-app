package conf

import (
	"log"

	"github.com/spf13/viper"
)

type Download struct {
	Host string
	Path string
}

type Proxy struct {
	Enable bool
	Socket string
}

// Config of konachan app
type Config struct {
	Addr     string
	Dbfile   string
	Download *Download
	Proxy    *Proxy
}

var (
	Cfg Config
)

func OpenCfgfile(filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
