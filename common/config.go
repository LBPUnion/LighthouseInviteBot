package common

import "github.com/BurntSushi/toml"

var configPath = "./lighthouseBot.toml"

type Config struct {
	Bot        BotConfig
	Lighthouse LighthouseConfig
}

type BotConfig struct {
	Token        string
	DisplayStats bool
}

type LighthouseConfig struct {
	APIKey       string
	ServerURL    string
	InstanceName string
}

func GetAPIURL() string {
	return LoadConfig().Lighthouse.ServerURL + "/api/v1"
}

func LoadConfig() Config {
	return LoadConfigPath(configPath)
}

func LoadConfigPath(path string) Config {
	var out Config
	toml.DecodeFile(path, &out)
	return out
}
