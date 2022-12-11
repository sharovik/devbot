package container

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	"github.com/sharovik/devbot/internal/dto"

	"github.com/sharovik/devbot/internal/database"

	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/log"
)

//DefinedEvent the interface for events
//@todo: move to the better place.
type DefinedEvent interface {
	//Execute The main execution method, which will run the actual functionality for the event
	Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error)

	//Install The installation method, which will executes the installation parts of the event
	Install() error

	//Update The update method, which will update the application to use new version of this event
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
func Init() (Main, error) {
	C := new(Main)
	cfg, err := config.Init()
	if err != nil {
		return Main{}, err
	}

	C.Config = cfg

	err = log.Init(C.Config.LogConfig)
	if err != nil {
		return Main{}, err
	}

	netTransport := &http.Transport{
		TLSHandshakeTimeout: time.Duration(C.Config.HTTPClient.TLSHandshakeTimeout) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: C.Config.HTTPClient.InsecureSkipVerify,
		},
	}

	httpClient := http.Client{
		Timeout:   time.Duration(C.Config.HTTPClient.RequestTimeout) * time.Second,
		Transport: netTransport,
	}

	bitBucketClient := client.BitBucketClient{}
	bitBucketClient.Init(&client.HTTPClient{
		Client:       &httpClient,
		BaseURL:      client.DefaultBitBucketBaseAPIUrl,
		ClientID:     C.Config.BitBucketConfig.ClientID,
		ClientSecret: C.Config.BitBucketConfig.ClientSecret,
	})
	C.BibBucketClient = &bitBucketClient
	C.HTTPClient = &client.HTTPClient{
		Client: &httpClient,
	}

	C.MessageClient = C.initMessageClient()

	if err := C.loadDictionary(); err != nil {
		panic(err)
	}

	C.MigrationService = database.MigrationService{
		Logger:     *log.Logger(),
		Dictionary: C.Dictionary,
	}

	return *C, nil
}

//Terminate terminates the properly connections
func (c *Main) Terminate() {
	if err := c.Dictionary.CloseDatabaseConnection(); err != nil {
		panic(err)
	}
}

func (c *Main) loadDictionary() error {
	dictionary := new(database.Dictionary)
	if err := dictionary.InitDatabaseConnection(c.Config.Database); err != nil {
		return err
	}

	c.Dictionary = dictionary

	return nil
}

func (c *Main) initMessageClient() client.MessageClientInterface {
	switch c.Config.MessagesAPIConfig.Type {
	case config.MessagesAPITypeSlack:
		h := c.HTTPClient

		h.SetOauthToken(c.Config.MessagesAPIConfig.WebAPIOAuthToken)
		h.SetBaseURL(c.Config.MessagesAPIConfig.BaseURL)

		sc := client.SlackClient{}
		sc.HTTPClient = h

		return sc
	default:
		panic(errors.New("Unknown messages API type"))
	}
}
