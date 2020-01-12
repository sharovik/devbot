package container

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"

	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/log"
)

//Main container object
type Main struct {
	Config      config.Config
	SlackClient client.SlackClientInterface
	Dictionary  dto.DevBotMessageDictionary
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
	container.Dictionary = container.loadDictionary()

	return container
}

//@todo: start using sqlite database for that
func (container Main) loadDictionary() dto.DevBotMessageDictionary {
	if container.Config.GetAppEnv() == config.EnvironmentTesting {
		_, filename, _, _ := runtime.Caller(0)
		dir := path.Join(path.Dir(filename), "../../")
		if err := os.Chdir(dir); err != nil {
			log.Logger().AddError(err).Msg("Error during the root folder pointer creation")
		}
	}

	pathToDictionary := fmt.Sprintf("./internal/dictionary/%s_dictionary.json", container.Config.AppDictionary)
	bytes, err := helper.FileToBytes(pathToDictionary)
	if err != nil {
		panic("Can't load dictionary: " + pathToDictionary)
	}

	var dictionary dto.DevBotMessageDictionary
	if err := json.Unmarshal(bytes, &dictionary); err != nil {
		panic("Can't unmarshal dictionary file. Error: " + err.Error())
	}

	return dictionary
}
