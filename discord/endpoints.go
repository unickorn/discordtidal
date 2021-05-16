package discord

import (
	"discordtidal/log"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

const (
	editApplicationData = "https://discord.com/api/v9/applications/%s"
	assets              = "https://discord.com/api/v9/oauth2/applications/%s/assets"
	deleteAsset         = "https://discord.com/api/v9/oauth2/applications/%s/assets/%s"
)

type Config struct {
	ApplicationId string
	Token         string
	UserAgent     string
}

var config Config

func LoadConfig() {
	config = Config{
		ApplicationId: "",
		Token:         "",
		UserAgent:     "",
	}

	b, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Log().Info("config not found, creating new")
		data, err := toml.Marshal(config)
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			panic(err)
		}
		return
	}
	if err := toml.Unmarshal(b, &config); err != nil {
		panic(err)
	}
	log.Log().Infof("config loaded!")
}

func GetConfig() *Config {
	return &config
}
