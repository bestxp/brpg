package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/rs/zerolog/log"
)

var config *Config = &Config{}
var mx sync.Once

const cfgFile = "config.json"

func init() {
	mx.Do(func() {
		f, err := os.ReadFile(cfgFile)
		if err != nil {
			log.Fatal().Err(err).Msg("failed load config")
		}

		if err = json.Unmarshal(f, config); err != nil {
			log.Fatal().Err(err).Msg("failed unmarshal config")
		}

	})
}

func GetConfig() *Config {
	return config
}

type Config struct {
	Resolution   Resolution `json:"resolution"`
	VsyncEnabled bool       `json:vsyncEnabled`

	Host string `json:"host"`
}
