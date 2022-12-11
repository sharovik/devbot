package slack

import (
	"regexp"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/log"
)

//MessageAttributes the validation message attributes
type MessageAttributes struct {
	Type    string
	Channel string
	Text    string
	User    string
	BotID   string
}

const (
	eventTypeMessage             = "message"
	eventTypeDesktopNotification = "desktop_notification"
	eventTypeFileShared          = "file_shared"
	eventTypeAppMention          = "app_mention"
)

var (
	acceptedMessageTypes = map[string]string{
		eventTypeMessage:             eventTypeMessage,
		eventTypeDesktopNotification: eventTypeDesktopNotification,
		eventTypeFileShared:          eventTypeFileShared,
		eventTypeAppMention:          eventTypeAppMention,
	}
)

func isValidMessage(msg MessageAttributes) bool {
	if acceptedMessageTypes[msg.Type] == "" {
		log.Logger().Debug().Str("type", msg.Type).Msg("Skip message check for this message type")
		return false
	}

	if msg.Channel == "" {
		log.Logger().Debug().Msg("Message channel cannot be empty")
		return false
	}

	if isGlobalAlertTriggered(msg.Text) {
		log.Logger().Debug().Msg("The global alert is triggered. Skipping.")
		return false
	}

	if msg.User == container.C.Config.MessagesAPIConfig.BotUserID || msg.BotID != "" {
		log.Logger().Debug().Msg("This message is from our bot user")
		return false
	}

	if "" == msg.Text {
		log.Logger().Debug().Msg("This message has empty text. Skipping.")
		return false
	}

	return true
}

func isGlobalAlertTriggered(text string) bool {
	re, err := regexp.Compile(`(?i)(\<\!(here|channel)\>)`)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to parse global alert text part")
		return false
	}

	return re.MatchString(text)
}
