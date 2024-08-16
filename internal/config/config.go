package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

const (
	Prod = "prod"
	Dev  = "dev"
)

var configPath = map[string]string{
	Prod: "configs/prod.yml",
	Dev:  "configs/dev.yml",
}

var secrets = map[string]string{
	"CONNECTION_STRING": "connectionString",
	"AUTH_SECRET":       "authSecret",
}

type Config struct {
	BuildMode            string `config:"buildMode" default:"dev"`
	URL                  string `config:"url"`
	ConnectionString     string `config:"connectionString"`
	GoogleClientId       string `config:"googleClientId"`
	AuthSecret           string `config:"authSecret"`
	EnableConsoleLogging bool   `config:"enableConsoleLogging"`
	LogFilePath          string `config:"logFilePath"`
}

func MustNewConfig() *Config {
	c := config.NewWithOptions("main", config.ParseEnv)
	c.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
	})

	c.AddDriver(yaml.Driver)

	err := c.LoadFiles("configs/build.yml")
	if err != nil {
		panic(err)
	}

	c.LoadOSEnvs(secrets)

	buildMode := c.String("buildMode", Dev)
	pathToConfig, ok := configPath[buildMode]
	if !ok {
		panic("unable to get configs path")
	}
	err = c.LoadFiles(pathToConfig)
	if err != nil {
		panic(err)
	}

	configStruct := Config{}
	if c.Decode(&configStruct) != nil {
		panic("unable to load configs to struct")
	}

	return &configStruct
}
