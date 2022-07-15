package common

import "github.com/BurntSushi/toml"

var configPath = "./lighthouseBot.toml"

type Config struct {
	Bot        BotConfig
	Lighthouse LighthouseConfig
}

type BotConfig struct {
	Token string
}

type LighthouseConfig struct {
	APIKey       string
	ServerAPIURL string
	ServerURL    string
	InstanceName string
}

func LoadConfig() Config {
	return LoadConfigPath(configPath)
}

func LoadConfigPath(path string) Config {
	var out Config
	toml.DecodeFile(path, &out)
	return out
}
