package models

type Config struct {
	Port                int    `yaml:"port"`
	LogPath             string `yaml:"logpath"`
	IsRequestsToFiles   bool   `yaml:"isrequeststofile"`
	RequestsToFilesPath string `yaml:"requeststofilespath"`
}

// NewConfig generates a base config
// Also known as config defaults
func NewConfig() Config {
	return Config{
		8080,
		"./api-logs.txt",
		false,
		"./http-requests",
	}
}
