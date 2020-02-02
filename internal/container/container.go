package container

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/sharovik/devbot/internal/database"

	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/log"
)

//Main container object
type Main struct {
	Config      config.Config
	SlackClient client.SlackClientInterface
	Dictionary  database.BaseDatabaseInterface
}

//C container variable
var C Main

//Init initialise container
func (container Main) Init() Main {
	container.Config = config.Init()

	_ = log.Init(log.Config(container.Config))

	netTransport := &http.Transport{
		TLSHandshakeTimeout: 7 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	slackClient := client.SlackClient{
		Client: &http.Client{
			Timeout:   time.Duration(15) * time.Second,
			Transport: netTransport,
		},
		BaseURL:    container.Config.SlackConfig.BaseURL,
		OAuthToken: container.Config.SlackConfig.OAuthToken,
	}

	container.SlackClient = slackClient
	if err := container.loadDictionary(); err != nil {
		panic(err)
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
	default:
		return fmt.Errorf("Unknown dictionary database used: %s ", container.Config.DatabaseConnection)
	}
}
