package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type ReminderBot struct {
	bot           *messaging_api.MessagingApiAPI
	channelSecret string
}

func (reminderBot *ReminderBot) callbackHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("/callback called...")

	callback, err := webhook.ParseRequest(reminderBot.channelSecret, req)
	if err != nil {
		log.Printf("Couldn't parse request, error: %+v\n", err)
		if errors.Is(err, webhook.ErrInvalidSignature) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	log.Println("Handling events: ")

	for _, event := range callback.Events {
		log.Printf("/callback called%+v...\n", event)

		// Only hanlding Text message events (early return otherwise)
		e, ok := event.(webhook.MessageEvent)
		if !ok {
			log.Printf("Unsupported message: %T\n", event)
			return
		}

		message, ok := e.Message.(webhook.TextMessageContent)
		if !ok {
			log.Printf("Unsupported message content: %T\n", e.Message)
			return
		}

		log.Println("Message received from user: " + message.Text)

		if _, err = reminderBot.bot.ReplyMessage(
			&messaging_api.ReplyMessageRequest{
				ReplyToken: e.ReplyToken,
				Messages: []messaging_api.MessageInterface{
					messaging_api.TextMessage{
						Text: message.Text,
					},
				},
			},
		); err != nil {
			log.Print(err)
		} else {
			// log.Println("Sent text reply.")
			log.Println("Message sending back to user: " + message.Text)
		}

	}
}

func main() {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	bot, err := messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	reminderBot := ReminderBot{
		channelSecret: channelSecret,
		bot:           bot,
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Wh1fty line reminder bot")
	})

	http.HandleFunc("/callback", reminderBot.callbackHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
