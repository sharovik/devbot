package container

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"

	_time "github.com/sharovik/devbot/internal/service/time"

	"github.com/sharovik/devbot/internal/dto/event"

	"github.com/sharovik/devbot/internal/database"

	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/log"
)

// Main container object
type Main struct {
	Config           config.Config
	MessageClient    client.MessageClientInterface
	BibBucketClient  client.GitClientInterface
	Dictionary       database.BaseDatabaseInterface
	HTTPClient       client.BaseHTTPClientInterface
	MigrationService database.MigrationService
	DefinedEvents    map[string]event.DefinedEventInterface
}

// C container variable
var C Main

// Init initialise container
func Init() (Main, error) {
	C = Main{}
	cfg, err := config.Init()
	if err != nil {
		return Main{}, err
	}

	C.Config = cfg

	err = log.Init(C.Config.LogConfig)
	if err != nil {
		return Main{}, err
	}

	bitBucketClient := client.BitBucketClient{}

	bH := initHttpClient()
	bH.BaseURL = client.DefaultBitBucketBaseAPIUrl
	bH.ClientID = C.Config.BitBucketConfig.ClientID
	bH.ClientSecret = C.Config.BitBucketConfig.ClientSecret

	bitBucketClient.Init(&bH)

	_time.InitNOW(C.Config.GetTimezone())

	C.BibBucketClient = &bitBucketClient

	h := initHttpClient()
	C.HTTPClient = &h

	C.MessageClient = C.initMessageClient()

	if err = C.loadDictionary(); err != nil {
		panic(err)
	}

	C.MigrationService = database.MigrationService{
		Logger:     *log.Logger(),
		Dictionary: C.Dictionary,
	}

	return C, nil
}

// Terminate terminates the properly connections
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
		h := initHttpClient()

		h.SetOauthToken(c.Config.MessagesAPIConfig.WebAPIOAuthToken)
		h.SetBaseURL(c.Config.MessagesAPIConfig.BaseURL)

		sc := client.SlackClient{}
		sc.HTTPClient = &h

		return sc
	default:
		panic(errors.New("unknown messages API type"))
	}
}

func initHttpClient() client.HTTPClient {
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

	return client.HTTPClient{
		Client: &httpClient,
	}
}
