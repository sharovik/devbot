package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

//BitBucketConfig struct for bitbucket config
type BitBucketConfig struct {
	ClientID                     string
	ClientSecret                 string
	ReleaseChannelMessageEnabled bool
	ReleaseChannel               string
	CurrentUserUUID              string
	RequiredReviewers            []BitBucketReviewer
}

//BitBucketReviewer is used for identifying of the reviewer user
type BitBucketReviewer struct {
	UUID     string
	SlackUID string
}

//Config configuration object
type Config struct {
	appEnv             string
	AppDictionary      string
	SlackConfig        SlackConfig
	BitBucketConfig    BitBucketConfig
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

	//DatabaseUsername env variable for database username
	DatabaseUsername = "DATABASE_USERNAME"

	//DatabasePassword env variable for database password
	DatabasePassword = "DATABASE_PASSWORD"

	//BitBucketClientID the client id which is used fo oauth token generation
	BitBucketClientID = "BITBUCKET_CLIENT_ID"

	//BitBucketClientSecret the client secret which is used fo oauth token generation
	BitBucketClientSecret = "BITBUCKET_CLIENT_SECRET"

	//BitBucketRequiredReviewers the required reviewers list separated by comma
	BitBucketRequiredReviewers = "BITBUCKET_REQUIRED_REVIEWERS"

	//BitBucketReleaseChannel the release channel ID. To that channel bot will publish the result of the release
	BitBucketReleaseChannel = "BITBUCKET_RELEASE_CHANNEL"

	//BitBucketReleaseChannelMessageEnabled the release channel ID. To that channel bot will publish the result of the release
	BitBucketReleaseChannelMessageEnabled = "BITBUCKET_RELEASE_CHANNEL_MESSAGE_ENABLE"

	//BitBucketCurrentUserUUID the current BitBucket user UUID the client credentials of which we are using in BITBUCKET_CLIENT_ID and BITBUCKET_CLIENT_SECRET
	BitBucketCurrentUserUUID = "BITBUCKET_USER_UUID"

	defaultMainChannelAlias       = "general"
	defaultBotName                = "devbot"
	defaultAppDictionary          = "slack"
	defaultDatabaseConnection     = "sqlite"
	defaultEnvFilePath            = "./.env"
	defaultEnvFileRootProjectPath = "./../../.env"
)

//Init initialise configuration for this project
func Init() Config {
	if !cfg.IsInitialised() {

		envPath = defaultEnvFilePath
		if _, err := os.Stat(envPath); err != nil {
			envPath = defaultEnvFileRootProjectPath
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

		dbConnection := defaultDatabaseConnection
		if os.Getenv(DatabaseConnection) != "" {
			dbConnection = os.Getenv(DatabaseConnection)
		}

		bitBucketReleaseChannelMessageEnabled := false
		bitBucketReleaseChannelMessageEnabledValue := os.Getenv(BitBucketReleaseChannelMessageEnabled)
		if bitBucketReleaseChannelMessageEnabledValue == "true" || bitBucketReleaseChannelMessageEnabledValue == "1" {
			bitBucketReleaseChannelMessageEnabled = true
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
			BitBucketConfig: BitBucketConfig{
				ClientID:                     os.Getenv(BitBucketClientID),
				ClientSecret:                 os.Getenv(BitBucketClientSecret),
				ReleaseChannel:               os.Getenv(BitBucketReleaseChannel),
				CurrentUserUUID:              os.Getenv(BitBucketCurrentUserUUID),
				ReleaseChannelMessageEnabled: bitBucketReleaseChannelMessageEnabled,
				RequiredReviewers:            prepareBitBucketReviewers(os.Getenv(BitBucketRequiredReviewers)),
			},
			initialised:        true,
			DatabaseConnection: dbConnection,
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

func prepareBitBucketReviewers(reviewers string) []BitBucketReviewer {
	entries := strings.Split(reviewers, ",")

	var result []BitBucketReviewer
	for _, value := range entries {
		userInfo := strings.Split(value, ":")

		if len(userInfo) != 0 {
			result = append(result, BitBucketReviewer{
				SlackUID: userInfo[0],
				UUID:     userInfo[1],
			})
		}
	}

	return result
}
