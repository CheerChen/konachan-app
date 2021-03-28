package conf

import "github.com/spf13/viper"

// Loader ...
type Loader interface {
	LoadFile(filename string) (*Config, error)
}

// NewLoader ...
func NewLoader() Loader {
	return &ViperLoader{
		viper.New(),
	}
}
