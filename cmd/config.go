package cmd

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"
)

var App *AppConfig

type AppConfig struct {
	Environment   string
	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	GithubToken   string
}

func NewAppConfig() (*AppConfig, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &AppConfig{
		Environment:   viper.GetString("ENVIRONMENT"),
		RedisHost:     viper.GetString("REDIS_HOST"),
		RedisPort:     viper.GetInt("REDIS_PORT"),
		RedisUsername: viper.GetString("REDIS_USERNAME"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		GithubToken:   viper.GetString("GITHUB_TOKEN"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *AppConfig) Validate() error {
	return v.ValidateStruct(c,
		v.Field(&c.Environment, v.Required),
		v.Field(&c.RedisHost, v.Required),
		v.Field(&c.RedisPort, v.Required, v.Min(1), v.Max(65535)),
		v.Field(&c.RedisUsername, v.Required),
		v.Field(&c.RedisPassword, v.Required),
		v.Field(&c.GithubToken, v.Required),
	)
}
