package discord

import (
	"github.com/pelletier/go-toml"
	"github.com/unickorn/discordtidal/log"
	"io/ioutil"
)

var config *Config

// Config is a struct that holds the configuration for discordtidal.
type Config struct {
	ApplicationId string
	Token         string
	UserAgent     string
	LogLevel      string
}

// LoadConfig loads the config.
func LoadConfig() {
	c := Config{}

	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Log().Infoln("Config not found, creating new")
		data, err := toml.Marshal(c)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			panic(err)
		}
		return
	}
	if err := toml.Unmarshal(b, &c); err != nil {
		panic(err)
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	log.SetLevel(c.LogLevel)
	log.Log().Infoln("Config loaded!")
	config = &c
}

// GetConfig returns the config.
func GetConfig() *Config {
	return config
}
