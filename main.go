package main

import (
	"github.com/llusardo/go-simple-slack-bot/slack"
	"net/http"
	"os"
)

var port string
var botToken string
var slackClientWrapper *slack.ClientWrapper

func init() {

}

func main() {
	port := ":" + os.Getenv("PORT")
	go http.ListenAndServe(port, nil)

	slackInitialization()
}

// slackInitialization is a function that initializes the Slackbot.
func slackInitialization() {
	botToken = os.Getenv("BOT_USER_OAUTH_ACCESS_TOKEN")
	slackClientWrapper = slack.NewClientWrapper(botToken)
	slack.EventsHandler(slackClientWrapper)
}
