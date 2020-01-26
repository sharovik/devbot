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
	appEnv             string
	AppDictionary      string
	SlackConfig        SlackConfig
	initialised        bool
	DatabaseConnection string
	DatabaseHost       string
	DatabaseUsername   string
	DatabasePassword   string
}

//cfg variable which contains initialised Config
var (
	cfg     Config
	envPath string
)

const (
	envFilePath = "./../../.env"

	//EnvironmentTesting constant for testing environment
	EnvironmentTesting = "testing"

	appEnv        = "APP_ENV"
	appDictionary = "APP_DICTIONARY"

	//SlackEnvUserID env variable for slack user ID
	SlackEnvUserID = "SLACK_USER_ID"

	//SlackEnvMainChannelID env variable for slack main channel ID
	SlackEnvMainChannelID = "SLACK_MAIN_CHANNEL_ID"

	//SlackEnvMainChannelAlias env variable for slack main channel alias
	SlackEnvMainChannelAlias = "SLACK_MAIN_CHANNEL_ALIAS"

	//SlackEnvBotName env variable for slack bot name
	SlackEnvBotName = "SLACK_BOT_NAME"

	//SlackEnvBaseURL env variable for slack base url
	SlackEnvBaseURL = "SLACK_BASE_URL"

	//SlackEnvOAuthToken env variable for slack oauth token
	SlackEnvOAuthToken = "SLACK_OAUTH_TOKEN"

	//DatabaseConnection env variable for database connection type
	DatabaseConnection = "DATABASE_CONNECTION"

	//DatabaseHost env variable for database host
	DatabaseHost = "DATABASE_HOST"

	//DatabaseHost env variable for database username
	DatabaseUsername = "DATABASE_USERNAME"

	//DatabaseHost env variable for database password
	DatabasePassword = "DATABASE_PASSWORD"

	defaultMainChannelAlias = "general"
	defaultBotName          = "devbot"
	defaultAppDictionary    = "slack"
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
		if os.Getenv(SlackEnvMainChannelAlias) != "" {
			mainChannelAlias = os.Getenv(SlackEnvMainChannelAlias)
		}

		BotName := defaultBotName
		if os.Getenv(SlackEnvBotName) != "" {
			mainChannelAlias = os.Getenv(SlackEnvBotName)
		}

		AppDictionary := defaultAppDictionary
		if os.Getenv(appDictionary) != "" {
			AppDictionary = os.Getenv(appDictionary)
		}

		cfg = Config{
			appEnv:        os.Getenv(appEnv),
			AppDictionary: AppDictionary,
			SlackConfig: SlackConfig{
				BaseURL:          os.Getenv(SlackEnvBaseURL),
				OAuthToken:       os.Getenv(SlackEnvOAuthToken),
				MainChannelAlias: mainChannelAlias,
				MainChannelID:    os.Getenv(SlackEnvMainChannelID),
				BotUserID:        os.Getenv(SlackEnvUserID),
				BotName:          BotName,
			},
			initialised:        true,
			DatabaseConnection: os.Getenv(DatabaseConnection),
			DatabaseHost:       os.Getenv(DatabaseHost),
			DatabaseUsername:   os.Getenv(DatabaseUsername),
			DatabasePassword:   os.Getenv(DatabasePassword),
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
	if flag.Lookup("test.v") != nil {
		return EnvironmentTesting
	}

	return c.appEnv
}

//SetToEnv method for saving env variable into memory + into .env file
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
