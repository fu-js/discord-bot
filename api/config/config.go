package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Viblo VibloBotConfig
	Port  string `envconfig:"PORT"`
}

type VibloBotConfig struct {
	Token     string `envconfig:"VIBLO_BOT_TOKEN"`
	ChannelID string `envconfig:"VIBLO_BOT_CHANNEL"`
}

var Global Config

func init() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}
	Global = config
}
