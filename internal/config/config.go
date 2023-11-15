package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/orm/clients"
)

// MessagesAPIConfig struct object
type MessagesAPIConfig struct {
	BaseURL          string
	OAuthToken       string
	WebAPIOAuthToken string
	BotUserID        string
	BotName          string
	MainChannelAlias string
	MainChannelID    string
	Type             string
}

// BitBucketConfig struct for bitbucket config
type BitBucketConfig struct {
	ClientID                     string
	ClientSecret                 string
	ReleaseChannelMessageEnabled bool
	ReleaseChannel               string
	CurrentUserUUID              string
	DefaultWorkspace             string
	DefaultMainBranch            string
	RequiredReviewers            []BitBucketReviewer
}

// BitBucketReviewer is used for identifying of the reviewer user
type BitBucketReviewer struct {
	UUID     string
	SlackUID string
}

// HTTPClient the configuration for the http client
type HTTPClient struct {
	RequestTimeout      int64
	TLSHandshakeTimeout int64
	InsecureSkipVerify  bool
}

// Config configuration object
type Config struct {
	appEnv            string
	Timezone          *time.Location
	LearningEnabled   bool
	MessagesAPIConfig MessagesAPIConfig
	BitBucketConfig   BitBucketConfig
	initialised       bool
	Database          clients.DatabaseConfig
	HTTPClient        HTTPClient
	LogConfig         log.Config
}

// cfg variable which contains initialised Config
var (
	cfg     Config
	envPath string
)

const (
	//EnvironmentTesting constant for testing environment
	EnvironmentTesting = "testing"

	envAppEnv          = "APP_ENV"
	envMessagesAPIType = "MESSAGES_API_TYPE"

	//EnvUserID env variable for user ID
	EnvUserID = "MESSAGES_API_USER_ID"

	//EnvMainChannelID env variable for main channel ID
	EnvMainChannelID = "MESSAGES_API_MAIN_CHANNEL_ID"

	//EnvMainChannelAlias env variable for main channel alias
	EnvMainChannelAlias = "MESSAGES_API_MAIN_CHANNEL_ALIAS"

	//EnvBotName env variable for bot name
	EnvBotName = "MESSAGES_API_BOT_NAME"

	//EnvBaseURL env variable for base url
	EnvBaseURL = "MESSAGES_API_BASE_URL"

	//EnvOAuthToken env variable for message oauth token
	EnvOAuthToken = "MESSAGES_API_OAUTH_TOKEN"

	//EnvWebAPIOAuthToken env variable for message web api oauth token.
	EnvWebAPIOAuthToken = "MESSAGES_API_WEB_API_OAUTH_TOKEN"

	//DatabaseConnection env variable for database connection type
	DatabaseConnection = "DATABASE_CONNECTION"

	//DatabaseHost env variable for database host
	DatabaseHost = "DATABASE_HOST"

	//DatabaseUsername env variable for database username
	DatabaseUsername = "DATABASE_USERNAME"

	//DatabasePassword env variable for database password
	DatabasePassword = "DATABASE_PASSWORD"

	//DatabaseName env variable for database name
	DatabaseName = "DATABASE_NAME"

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

	//BitBucketDefaultWorkspace the default workspace which can be used in the functionality, once you don't have PR link, from where to get this information
	BitBucketDefaultWorkspace = "BITBUCKET_DEFAULT_WORKSPACE"

	//BitBucketDefaultMainBranch the default main branch which can be used in cases, when you can't get the information from the PR link
	BitBucketDefaultMainBranch = "BITBUCKET_DEFAULT_MAIN_BRANCH"

	httpClientRequestTimeout      = "HTTP_CLIENT_REQUEST_TIMEOUT"
	httpClientTLSHandshakeTimeout = "HTTP_CLIENT_TLS_HANDSHAKE_TIMEOUT"
	httpClientInsecureSkipVerify  = "HTTP_CLIENT_INSECURE_SKIP_VERIFY"

	//learningEnabled enables or disables the learning mode. If enabled, the bot will try to ask in the main channel, how to react on that message.
	learningEnabled = "LEARNING_MODE_ENABLED"

	envLogOutput            = "LOG_OUTPUT"
	envLogLevel             = "LOG_LEVEL"
	envLogFieldContext      = "LOG_FIELD_CONTEXT"
	envLogFieldLevelName    = "LOG_FIELD_LEVEL_NAME"
	envLogFieldErrorMessage = "LOG_FIELD_ERROR_MESSAGE"

	//MessagesAPITypeSlack message messages API type
	MessagesAPITypeSlack = "slack"

	defaultMainChannelAlias       = "general"
	defaultBotName                = "devbot"
	defaultMessagesAPIType        = MessagesAPITypeSlack
	defaultDatabaseConnection     = "sqlite"
	defaultEnvFilePath            = "./.env"
	defaultEnvFileRootProjectPath = "./../../.env"

	//AWSSecretsBucket the bucket for secrets manager. If it's specified then the secret values will be loaded
	AWSSecretsBucket = "AWS_SECRETS_BUCKET"

	//AWSSecretsRegion the region selected for the aws session
	AWSSecretsRegion = "AWS_REGION"

	envTimezone = "APP_TIMEZONE"
)

var DefaultTimezone = time.UTC

// Init initialise configuration for this project
func Init() (Config, error) {
	if !cfg.IsInitialised() {

		envPath = defaultEnvFilePath
		if _, err := os.Stat(envPath); err != nil {
			envPath = defaultEnvFileRootProjectPath
		}

		if err := godotenv.Load(envPath); err != nil {
			return Config{}, err
		}

		httpC, err := initHTTPClientConfig()
		if err != nil {
			return Config{}, err
		}

		cfg = Config{
			appEnv:            os.Getenv(envAppEnv),
			LearningEnabled:   getBoolValue(learningEnabled),
			MessagesAPIConfig: initMessagesAPIConfig(),
			BitBucketConfig:   initBitbucketConfig(),
			initialised:       true,
			HTTPClient:        httpC,
			Database:          initDatabaseConfig(),
			LogConfig:         initLogConfig(),
		}

		cfg.loadTimezone()

		cfg = loadSecrets(cfg)

		return cfg, nil
	}

	return cfg, nil
}

func loadSecrets(cfg Config) Config {
	if os.Getenv(AWSSecretsBucket) != "" {
		secrets, _ := GetSecret(os.Getenv(AWSSecretsBucket), os.Getenv(AWSSecretsRegion))
		cfg.BitBucketConfig.ClientID = secrets.BitBucketClientID
		cfg.BitBucketConfig.ClientSecret = secrets.BitBucketClientSecret
		cfg.MessagesAPIConfig.OAuthToken = secrets.MessagesAPIOAuthToken
		cfg.MessagesAPIConfig.WebAPIOAuthToken = secrets.MessagesAPIWebAPIOAuthToken
	}

	return cfg
}

func initHTTPClientConfig() (c HTTPClient, err error) {
	var requestTimeout, tLSHandshakeTimeout int64
	if os.Getenv(httpClientRequestTimeout) != "" {
		requestTimeout, err = strconv.ParseInt(os.Getenv(httpClientRequestTimeout), 10, 64)
		if err != nil {
			return HTTPClient{}, err
		}
	}

	if os.Getenv(httpClientTLSHandshakeTimeout) != "" {
		tLSHandshakeTimeout, _ = strconv.ParseInt(os.Getenv(httpClientTLSHandshakeTimeout), 10, 64)
		if err != nil {
			return HTTPClient{}, err
		}
	}

	return HTTPClient{
		RequestTimeout:      requestTimeout,
		TLSHandshakeTimeout: tLSHandshakeTimeout,
		InsecureSkipVerify:  getBoolValue(httpClientInsecureSkipVerify),
	}, nil
}

func initBitbucketConfig() BitBucketConfig {
	return BitBucketConfig{
		ClientID:                     os.Getenv(BitBucketClientID),
		ClientSecret:                 os.Getenv(BitBucketClientSecret),
		ReleaseChannel:               os.Getenv(BitBucketReleaseChannel),
		CurrentUserUUID:              os.Getenv(BitBucketCurrentUserUUID),
		DefaultWorkspace:             os.Getenv(BitBucketDefaultWorkspace),
		DefaultMainBranch:            os.Getenv(BitBucketDefaultMainBranch),
		ReleaseChannelMessageEnabled: getBoolValue(BitBucketReleaseChannelMessageEnabled),
		RequiredReviewers:            PrepareBitBucketReviewers(os.Getenv(BitBucketRequiredReviewers)),
	}
}

func initMessagesAPIConfig() MessagesAPIConfig {
	oAuthToken := os.Getenv(EnvOAuthToken)
	webAPIOAuthToken := os.Getenv(EnvWebAPIOAuthToken)
	mainChannelAlias := defaultMainChannelAlias
	if os.Getenv(EnvMainChannelAlias) != "" {
		mainChannelAlias = os.Getenv(EnvMainChannelAlias)
	}

	botName := defaultBotName
	if os.Getenv(EnvBotName) != "" {
		botName = os.Getenv(EnvBotName)
	}

	messagesAPIType := defaultMessagesAPIType
	if os.Getenv(envMessagesAPIType) != "" {
		messagesAPIType = os.Getenv(envMessagesAPIType)
	}

	return MessagesAPIConfig{
		BaseURL:          os.Getenv(EnvBaseURL),
		OAuthToken:       oAuthToken,
		WebAPIOAuthToken: webAPIOAuthToken,
		MainChannelAlias: mainChannelAlias,
		MainChannelID:    os.Getenv(EnvMainChannelID),
		BotUserID:        os.Getenv(EnvUserID),
		BotName:          botName,
		Type:             messagesAPIType,
	}
}

func initDatabaseConfig() clients.DatabaseConfig {
	dbConnection := defaultDatabaseConnection
	if os.Getenv(DatabaseConnection) != "" {
		dbConnection = os.Getenv(DatabaseConnection)
	}

	return clients.DatabaseConfig{
		Type:     dbConnection,
		Host:     os.Getenv(DatabaseHost),
		Username: os.Getenv(DatabaseUsername),
		Password: os.Getenv(DatabasePassword),
		Database: os.Getenv(DatabaseName),
	}
}

func initLogConfig() log.Config {
	fieldContext := log.FieldContext
	if os.Getenv(envLogFieldContext) != "" {
		fieldContext = os.Getenv(envLogFieldContext)
	}

	fieldLevelName := log.FieldLevelName
	if os.Getenv(envLogFieldLevelName) != "" {
		fieldLevelName = os.Getenv(envLogFieldLevelName)
	}

	fieldErrorMessage := log.FieldErrorMessage
	if os.Getenv(envLogFieldErrorMessage) != "" {
		fieldErrorMessage = os.Getenv(envLogFieldErrorMessage)
	}

	return log.Config{
		Env:               os.Getenv(envAppEnv),
		Level:             os.Getenv(envLogLevel),
		Output:            os.Getenv(envLogOutput),
		FieldContext:      fieldContext,
		FieldLevelName:    fieldLevelName,
		FieldErrorMessage: fieldErrorMessage,
	}
}

// IsInitialised method which retrieves current status of object
func (c Config) IsInitialised() bool {
	return c.initialised
}

// GetAppEnv retrieve current environment
func (c Config) GetAppEnv() string {
	if flag.Lookup("test.v") != nil {
		return EnvironmentTesting
	}

	return c.appEnv
}

// SetToEnv method for saving env variable into memory + into .env file
func (c Config) SetToEnv(field string, value string, writeToEnvFile bool) error {

	if writeToEnvFile {
		f, err := os.OpenFile(envPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer f.Close()
		if _, err = f.WriteString(fmt.Sprintf("\n%s=%s", field, value)); err != nil {
			return err
		}
	}

	if err := os.Setenv(field, value); err != nil {
		return err
	}

	return nil
}

// PrepareBitBucketReviewers method retrieves the list of bitbucket reviewers
func PrepareBitBucketReviewers(reviewers string) []BitBucketReviewer {
	entries := strings.Split(reviewers, ",")

	var result []BitBucketReviewer
	for _, value := range entries {
		userInfo := strings.Split(value, ":")

		if len(userInfo) != 0 && userInfo[0] != "" {
			result = append(result, BitBucketReviewer{
				SlackUID: userInfo[0],
				UUID:     userInfo[1],
			})
		}
	}

	return result
}

func getBoolValue(field string) bool {
	res := false
	if os.Getenv(field) == "true" || os.Getenv(field) == "1" {
		res = true
	}

	return res
}

func (c *Config) loadTimezone() {
	if os.Getenv(envTimezone) != "" {
		c.Timezone, _ = time.LoadLocation(os.Getenv(envTimezone))
	}

	if c.Timezone == nil {
		c.Timezone = DefaultTimezone
	}
}

func (c Config) GetTimezone() (timeZone *time.Location) {
	if c.Timezone == nil {
		c.Timezone = DefaultTimezone
	}

	return c.Timezone
}
