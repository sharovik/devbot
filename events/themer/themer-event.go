package themer

import (
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service"
)

const EventName = "themer_event"

type ThemerEvent struct {
	EventName string
}

var Event = ThemerEvent{
	EventName: EventName,
}

func (e ThemerEvent) Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error) {
	var answer = message
	go func() {
		if message.OriginalMessage.Files != nil {
			file, err := service.ProcessFiles(message.OriginalMessage)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to process file")

				answer = fileErrorMessage(message.Channel, file, err)
			}

			message.OriginalMessage.Files = nil
		}
	}()

	answer.Text = prepareThemeInstructions()
	return answer, nil
}

func prepareThemeInstructions() string {
	return "In that archive you can find 2 directories - preview and wordpress(directory contains the wordpress template)\n\n Installation guide:\n - copy wordpress directory into wp-content/themes directory\n - go to admin dashboard of your wordpress site and install your theme"
}

func fileErrorMessage(channelID string, file dto.File, err error) dto.SlackRequestChatPostMessage {
	return dto.SlackRequestChatPostMessage{
		Text:    fmt.Sprintf("Can't process the file. \nReason: %s\nFile name: %s\nFile type: %s", err.Error(), file.Name, file.Filetype),
		Channel: channelID,
		AsUser:  true,
	}
}