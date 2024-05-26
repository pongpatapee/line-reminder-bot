package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"

	"github.com/joho/godotenv"
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

	// NOTE: might want to change returns -> continue
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

		timeTilReminder, reminderMessage, err := ParseReminderCommand(message.Text)
		if err != nil {
			if err.Error() == "not a reminder" {
				return
			}

			sendReplyMessage(reminderBot, e, err.Error())
			return
		}

		log.Printf("Will send reminder message in %v\n", timeTilReminder)

		reminderTime := time.Now().Add(timeTilReminder)
		reminderConfirmation := fmt.Sprintf("Reminder set!\nReminder Time: %v\nReminder message: %v", reminderTime, reminderMessage)
		sendReplyMessage(reminderBot, e, reminderConfirmation)

		time.AfterFunc(timeTilReminder, func() {
			log.Printf("Sending reminder message '%v'\n", reminderMessage)
			sendReplyMessage(reminderBot, e, reminderMessage)
		})
	}
}

func sendReplyMessage(reminderBot *ReminderBot, e webhook.MessageEvent, message string) {
	if _, err := reminderBot.bot.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: e.ReplyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: message,
				},
			},
		},
	); err != nil {
		log.Println("An error occured while sending a reply message")
		log.Println(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "From /test")

		duration, err := time.ParseDuration("5s")
		if err != nil {
			return
		}
		time.AfterFunc(duration, func() {
			fmt.Printf("yoyo brother what's good\n")
		})
	})

	http.HandleFunc("/callback", reminderBot.callbackHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
