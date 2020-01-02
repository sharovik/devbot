package themerwordpress

import (
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service"
)

//EventName the name of the event
const EventName = "themer_wordpress_event"

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
	return "In that archive you can find 2 directories - preview(which contains the html preview of your design) and wordpress(directory contains the wordpress template)\n\n Installation guide:\n - copy wordpress directory into wp-content/themes directory\n - go to admin dashboard of your wordpress site and install your theme"
}

func fileErrorMessage(channelID string, file dto.File, err error) dto.SlackRequestChatPostMessage {
	return dto.SlackRequestChatPostMessage{
		Text:    fmt.Sprintf("Can't process the file. \nReason: %s\nFile name: %s\nFile type: %s", err.Error(), file.Name, file.Filetype),
		Channel: channelID,
		AsUser:  true,
	}
}
