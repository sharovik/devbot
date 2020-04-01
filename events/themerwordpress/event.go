package themerwordpress

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/log"
)

const (
	//EventName the name of the event
	EventName    = "themer_wordpress_event"
	EventVersion = "1.0.0"

	zipFileType           = "zip"
	defaultResultFilename = "result.zip"
)

var supportedFileTypes = map[string]string{
	zipFileType: zipFileType,
}

//ThemerEvent the struct for the event object
type ThemerEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = ThemerEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e ThemerEvent) Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error) {
	var answer = message
	if message.OriginalMessage.Files != nil {
		file, err := processFiles(message.OriginalMessage)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to process file")

			answer = fileErrorMessage(message.Channel, file, err)
			return answer, nil
		}

		message.OriginalMessage.Files = nil
	}

	answer.Text = prepareThemeInstructions()
	return answer, nil
}

func (e ThemerEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Start event Install")
	eventId, err := container.C.Dictionary.FindEventByAlias(EventName)
	if err != nil {
		log.Logger().AddError(err).Msg("Error during FindEventBy method execution")
		return err
	}

	if eventId == 0 {
		log.Logger().Info().
			Str("event_name", EventName).
			Str("event_version", EventVersion).
			Msg("Event wasn't installed. Trying to install it")

		eventId, err := container.C.Dictionary.InsertEvent(EventName)
		if err != nil {
			log.Logger().AddError(err).Msg("Error during FindEventBy method execution")
			return err
		}

		log.Logger().Debug().
			Str("event_name", EventName).
			Str("event_version", EventVersion).
			Int64("event_id", eventId).
			Msg("Event installed")
	}

	return nil
}

func (e ThemerEvent) Update() error {
	return nil
}

func isValidFile(fileType string) bool {
	return supportedFileTypes[fileType] != ""
}

func validateFiles(files []dto.File) (dto.File, error) {
	for _, file := range files {
		if !isValidFile(file.Filetype) {
			err := fmt.Errorf("Wrong file type ")
			log.Logger().AddError(err).Interface("file", file).Msg("Wrong file type")
			return file, err
		}
	}

	return dto.File{}, nil
}

func processFile(channel string, file dto.File) (dto.File, error) {
	log.Logger().Debug().
		Str("url", file.URLPrivate).
		Msg("Start processing file")

	//First we need to download the file
	tmpFile, err := downloadFile(file.URLPrivate)
	if err != nil {
		return file, err
	}

	//Now we need to unzip the file and save the destination folder path
	var (
		src         = filepath.Join(os.TempDir(), file.ID)
		pathToFiles = src + "/downloaded_template"
	)

	log.Logger().Debug().
		Str("src", src).
		Str("path_to_files", pathToFiles).
		Msg("Start unzip")

	_, err = helper.Unzip(tmpFile.Name(), pathToFiles)
	if err != nil {
		return file, err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return file, err
	}

	log.Logger().Debug().
		Str("template_dir", pathToFiles).
		Str("current_dir", currentDir).
		Msg("Template dir generated")

	//We run the command which compiles the template.
	//This will create in src 2 directories: one is for template html preview and second one for template

	cmd := exec.Command(filepath.Join(currentDir, "./scripts/themer/themer.phar"), fmt.Sprintf("--path=%s", pathToFiles))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Logger().AddError(err).
			Interface("file", file).
			Msg("Failed generate template")
		return file, err
	}

	//Now we need to remove the downloaded dir and zip the contains of src directory
	if err := deleteSrc(pathToFiles); err != nil {
		return file, err
	}

	resultFilePath := src + fmt.Sprintf("/%s", defaultResultFilename)
	if err := helper.Zip(src, resultFilePath); err != nil {
		return file, err
	}

	log.Logger().Debug().Str("result_zip_path", src+"/result.zip").Msg("Zip file created")

	if _, _, err := container.C.SlackClient.AttachFileTo(channel, resultFilePath, defaultResultFilename); err != nil {
		return file, err
	}

	if err := deleteSrc(src); err != nil {
		return file, err
	}

	return file, nil
}

func downloadFile(url string) (*os.File, error) {
	log.Logger().StartMessage("Download file")
	// Get the data
	resp, _, err := container.C.SlackClient.Request(http.MethodGet, url, []byte(``))
	if err != nil {
		return nil, err
	}

	// Create the file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "devbot-*.zip")
	if err != nil {
		return nil, err
	}

	if _, err = tmpFile.Write(resp); err != nil {
		return nil, err
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	log.Logger().FinishMessage("Download file")
	return tmpFile, nil
}

func deleteSrc(src string) error {
	return os.RemoveAll(src)
}

//processFiles method which processes the received files
func processFiles(message dto.SlackResponseEventMessage) (dto.File, error) {
	log.Logger().Debug().
		Interface("files", message.Files).
		Msg("Files received")

	file, err := validateFiles(message.Files)
	if err != nil {
		return file, err
	}

	for _, fileReceived := range message.Files {
		file, err := processFile(message.Channel, fileReceived)
		if err != nil {
			return file, err
		}
	}

	return dto.File{}, nil
}

func prepareThemeInstructions() string {
	return "In that archive you can find 2 directories - preview(which contains the html preview of your design) and wordpress(directory contains the wordpress template)\n\n Installation guide:\n - copy wordpress directory into wp-content/themes directory\n - go to admin dashboard of your wordpress site and install your theme"
}

func fileErrorMessage(channelID string, file dto.File, err error) dto.SlackRequestChatPostMessage {
	return dto.SlackRequestChatPostMessage{
		Text:    fmt.Sprintf("Can't process the file. \nReason: %s\nFile name: %s\nFile type: %s", err.Error(), file.Name, file.Filetype),
		Channel: channelID,
		AsUser:  true,
	}
}
