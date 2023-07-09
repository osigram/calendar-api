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
	Prod: "config/prod.yml",
	Dev:  "config/dev.yml",
}

type Config struct {
	BuildMode        string `config:"buildMode" default:"dev"`
	Url              string `config:"url"`
	ConnectionString string `config:"connectionString"`
	LogToConsole     bool   `config:"logToConsole"`
	LogFilePath      string `config:"logFilePath"`
}

func NewConfig() Config {
	c := config.NewWithOptions("main", config.ParseEnv)
	c.WithOptions(func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
	})

	c.AddDriver(yaml.Driver)

	err := c.LoadFiles("build.yml")
	if err != nil {
		panic(err)
	}

	c.LoadOSEnvs(map[string]string{"CONNECTION_STRING": "connectionString"})

	buildMode := c.String("buildMode", "dev")
	pathToConfig, ok := configPath[buildMode]
	if !ok {
		panic("Error to get config path")
	}
	err = c.LoadFiles(pathToConfig)
	if err != nil {
		panic(err)
	}

	configStruct := Config{}
	if c.Decode(&configStruct) != nil {
		panic("Error to load config to struct")
	}

	return configStruct
}
