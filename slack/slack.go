package slack

import (
	"fmt"
	"github.com/slack-go/slack"
	"strings"
)

const helpMessage = "The help command is still under development"

type ClientWrapper struct {
	client *slack.RTM
}

func NewClientWrapper(apiKey string) *ClientWrapper {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	return &ClientWrapper{client: rtm}
}

func EventsHandler(SlackClientWrapper *ClientWrapper) {
	slackClient := SlackClientWrapper.client

	for msg := range slackClient.IncomingEvents {
		fmt.Println("Event Received: ", msg.Type)
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			botTagString := fmt.Sprintf("<@%s> ", slackClient.GetInfo().User.ID)
			if !strings.Contains(event.Msg.Text, botTagString) {
				continue
			}

			processMessageEvent(slackClient, botTagString, event)

		case *slack.ConnectedEvent:
			fmt.Printf("Connection counter %v\n", event.ConnectionCount)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %s\n", event)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %s\n", event.Value)

		case *slack.DesktopNotificationEvent:
			fmt.Printf("Desktop Notification: %v\n", event)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", event.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials\n")
			return
		default:
		}
	}
}

func processMessageEvent(slackClient *slack.RTM, botTagString string, slackMessageEvent *slack.MessageEvent){
	message := strings.Replace(slackMessageEvent.Msg.Text, botTagString, "", -1)

	if message == "" {
		sendHelp(slackClient, slackMessageEvent.Channel)
	}

	splitMessage := strings.Fields(message)

	switch strings.ToLower(splitMessage[0]) {
	case "help":
		sendHelp(slackClient, slackMessageEvent.Channel)
	case "echo":
		sendEchoMessage(slackClient, strings.Join(splitMessage[:], " "), slackMessageEvent.Channel)
	case "video":
		sendVideoMessage(slackClient, slackMessageEvent.Channel)
	}
}

// sendHelp is a working help message, for reference.
func sendHelp(slackClient *slack.RTM, slackChannel string) {
	slackClient.SendMessage(slackClient.NewOutgoingMessage(helpMessage, slackChannel))
}

// sendEchoMessage will just echo anything after the echo keyword.
func sendEchoMessage(slackClient *slack.RTM, message, slackChannel string) {
	splitMessage := strings.Fields(strings.ToLower(message))
	slackClient.SendMessage(slackClient.NewOutgoingMessage(strings.Join(splitMessage[1:], " "), slackChannel))
}

// TODO make a request to youtube and return the first video obtained
// sendVideoMessage will return a fix video after the video keyword.
func sendVideoMessage(slackClient *slack.RTM, slackChannel string) {
	postMessageParameters:= slack.PostMessageParameters{UnfurlMedia: true, UnfurlLinks: true}
	slackClient.PostMessage(slackChannel, slack.MsgOptionText("https://www.youtube.com/watch?v=Rh64GkNRDNU&ab_channel=MarySpender", true), slack.MsgOptionPostMessageParameters(postMessageParameters))
}

