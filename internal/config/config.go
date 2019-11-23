package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

//SlackConfig struct object
type SlackConfig struct {
	BaseURL          string
	OAuthToken       string
	BotUserID        string
	BotName          string
	MainChannelAlias string
	MainChannelID    string
}

//Config configuration object
type Config struct {
	appEnv      string
	SlackConfig SlackConfig
	initialised bool
}

//cfg variable which contains initialised Config
var (
	cfg     Config
	envPath string
)

const (
	envFilePath = "./../../.env"

	environmentTesting = "testing"

	defaultMainChannelAlias = "general"
	defaultBotName          = "devbot"
)

//Init initialise configuration for this project
func Init() Config {
	if !cfg.IsInitialised() {

		envPath = envFilePath
		if _, err := os.Stat(envPath); err != nil {
			//In tests directory cursor is equal to config dir path.
			//When in main package directory cursor is always in root dir of project
			envPath = "./.env"
		}

		if err := godotenv.Load(envPath); err != nil {
			panic(err)
		}

		mainChannelAlias := defaultMainChannelAlias
		if os.Getenv("SLACK_MAIN_CHANNEL_ALIAS") != "" {
			mainChannelAlias = os.Getenv("SLACK_MAIN_CHANNEL_ALIAS")
		}

		BotName := defaultBotName
		if os.Getenv("SLACK_BOT_NAME") != "" {
			mainChannelAlias = os.Getenv("SLACK_BOT_NAME")
		}

		cfg = Config{
			appEnv: os.Getenv("APP_ENV"),
			SlackConfig: SlackConfig{
				BaseURL:          os.Getenv("SLACK_BASE_URL"),
				OAuthToken:       os.Getenv("SLACK_OAUTH_TOKEN"),
				MainChannelAlias: mainChannelAlias,
				MainChannelID:    os.Getenv("SLACK_MAIN_CHANNEL_ID"),
				BotUserID:        os.Getenv("SLACK_USER_ID"),
				BotName:          BotName,
			},
			initialised: true,
		}

		return cfg
	}

	return Config{}
}

//IsInitialised method which retrieves current status of object
func (c Config) IsInitialised() bool {
	return c.initialised
}

//GetAppEnv retrieve current environment
func (c Config) GetAppEnv() string {
	if c.appEnv == "" && (flag.Lookup("test.v") != nil) {
		return environmentTesting
	}

	return c.appEnv
}

func (c Config) GetSlackConfiguration() SlackConfig {
	return c.SlackConfig
}

func (c Config) SetToEnv(field string, value string, writeToEnvFile bool) error {

	if writeToEnvFile {
		f, err := os.OpenFile(envPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer f.Close()
		if _, err := f.WriteString(fmt.Sprintf("\n%s=%s", field, value)); err != nil {
			return err
		}
	}

	if err := os.Setenv(field, value); err != nil {
		return err
	}

	return nil
}
