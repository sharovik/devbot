package container

import (
	"crypto/tls"
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
	"net/http"
	"time"

	"github.com/sharovik/devbot/internal/database"

	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/log"
)

//DefinedEvent the interface for events
//@todo: move to the better place.
type DefinedEvent interface {
	//The main execution method, which will run the actual functionality for the event
	Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error)

	//The installation method, which will executes the installation parts of the event
	Install() error

	//The update method, which will update the application to use new version of this event
	Update() error
}

//Main container object
type Main struct {
	Config           config.Config
	MessageClient    client.MessageClientInterface
	BibBucketClient  client.GitClientInterface
	Dictionary       database.BaseDatabaseInterface
	HTTPClient       client.BaseHTTPClientInterface
	MigrationService database.MigrationService
	DefinedEvents    map[string]DefinedEvent
}

//C container variable
var C Main

//Init initialise container
func (container Main) Init() Main {
	container.Config = config.Init()

	_ = log.Init(container.Config.LogConfig)

	netTransport := &http.Transport{
		TLSHandshakeTimeout: time.Duration(container.Config.HTTPClient.TLSHandshakeTimeout) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: container.Config.HTTPClient.InsecureSkipVerify,
		},
	}

	httpClient := http.Client{
		Timeout:   time.Duration(container.Config.HTTPClient.RequestTimeout) * time.Second,
		Transport: netTransport,
	}

	bitBucketClient := client.BitBucketClient{}
	bitBucketClient.Init(&client.HTTPClient{
		Client:       &httpClient,
		BaseURL:      client.DefaultBitBucketBaseAPIUrl,
		ClientID:     container.Config.BitBucketConfig.ClientID,
		ClientSecret: container.Config.BitBucketConfig.ClientSecret,
	})
	container.BibBucketClient = &bitBucketClient

	slackClient := client.SlackClient{
		Client:     &httpClient,
		BaseURL:    container.Config.SlackConfig.BaseURL,
		OAuthToken: container.Config.SlackConfig.WebAPIOAuthToken,
	}

	container.HTTPClient = &client.HTTPClient{
		Client: &httpClient,
	}

	container.MessageClient = slackClient
	if err := container.loadDictionary(); err != nil {
		panic(err)
	}

	container.MigrationService = database.MigrationService{
		Logger:     *log.Logger(),
		Dictionary: container.Dictionary,
	}

	return container
}

//Terminate terminates the properly connections
func (container *Main) Terminate() {
	if err := container.Dictionary.CloseDatabaseConnection(); err != nil {
		panic(err)
	}
}

func (container *Main) loadDictionary() error {
	switch container.Config.DatabaseConnection {
	case database.ConnectionSQLite:
		dictionary := database.SQLiteDictionary{
			Cfg: container.Config,
		}

		if err := dictionary.InitDatabaseConnection(); err != nil {
			return err
		}

		container.Dictionary = &dictionary
		return nil
	case database.ConnectionMySQL:
		dictionary := database.MySQLDictionary{
			Cfg: container.Config,
		}

		if err := dictionary.InitDatabaseConnection(); err != nil {
			return err
		}

		container.Dictionary = &dictionary
		return nil
	default:
		return fmt.Errorf("Unknown dictionary database used: %s ", container.Config.DatabaseConnection)
	}
}
