package conf

import (
	"github.com/spf13/viper"
)

// ViperLoader ...
type ViperLoader struct {
	*viper.Viper
}

// LoadFile ...
func (v *ViperLoader) LoadFile(filename string) (*Config, error) {
	c := &Config{}
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(c)
	return c, err
}
